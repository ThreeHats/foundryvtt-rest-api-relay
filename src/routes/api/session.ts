import { Router, Request, Response } from 'express';
import express from 'express';
import { requestForwarderMiddleware } from '../../middleware/requestForwarder';
import { authMiddleware, trackApiUsage } from '../../middleware/auth';
import { ClientManager } from '../../core/ClientManager';
import { safeResponse } from '../shared';
import { log } from '../../utils/logger';
import crypto from 'crypto';
import * as puppeteer from 'puppeteer';
import { getRedisClient } from '../../config/redis';
import { registerHeadlessSession } from "../../workers/headlessSessions";

// Temporary handshake tokens
interface PendingHandshake {
    apiKey: string;
    foundryUrl: string;
    worldName?: string;
    username: string;
    publicKey: string;    // To send to client
    privateKey: string;   // To keep on server for decryption
    nonce: string;
    expires: number;
  }
  
const pendingHandshakes = new Map<string, PendingHandshake>();

export const browserSessions = new Map<string, puppeteer.Browser>();

// Track pending sessions (browsers launched but waiting for client connection)
// This is critical for cleanup - these can become orphaned if the connection never completes
interface PendingSession {
    sessionId: string;
    browser: puppeteer.Browser;
    apiKey: string;
    expectedClientId: string;
    startTime: number;
    foundryUrl: string;
    worldName?: string;
    username: string;
}
export const pendingSessions = new Map<string, PendingSession>();

// Changed to support multiple sessions per API key (especially for local mode testing)
// Maps: apiKey -> sessionId -> session data
export const apiKeyToSessions = new Map<string, Map<string, { sessionId: string, clientId: string, lastActivity: number }>>();
// Legacy export for backward compatibility
// TODO: Remove this once all clients are updated to use multiple sessions per API key
export const apiKeyToSession = {
    get: (apiKey: string) => {
        const sessions = apiKeyToSessions.get(apiKey);
        if (!sessions || sessions.size === 0) return undefined;
        // Return the most recent session for backward compatibility
        let mostRecent: { sessionId: string, clientId: string, lastActivity: number } | undefined;
        for (const session of sessions.values()) {
            if (!mostRecent || session.lastActivity > mostRecent.lastActivity) {
                mostRecent = session;
            }
        }
        return mostRecent;
    },
    set: (apiKey: string, session: { sessionId: string, clientId: string, lastActivity: number }) => {
        let sessions = apiKeyToSessions.get(apiKey);
        if (!sessions) {
            sessions = new Map();
            apiKeyToSessions.set(apiKey, sessions);
        }
        sessions.set(session.sessionId, session);
    },
    delete: (apiKey: string) => {
        apiKeyToSessions.delete(apiKey);
    },
    deleteSession: (apiKey: string, sessionId: string) => {
        const sessions = apiKeyToSessions.get(apiKey);
        if (sessions) {
            sessions.delete(sessionId);
            if (sessions.size === 0) {
                apiKeyToSessions.delete(apiKey);
            }
        }
    },
    getAll: (apiKey: string) => {
        const sessions = apiKeyToSessions.get(apiKey);
        return sessions ? Array.from(sessions.values()) : [];
    },
    entries: () => {
        const allEntries: [string, { sessionId: string, clientId: string, lastActivity: number }][] = [];
        for (const [apiKey, sessions] of apiKeyToSessions.entries()) {
            for (const session of sessions.values()) {
                allEntries.push([apiKey, session]);
            }
        }
        return allEntries[Symbol.iterator]();
    }
};
const pendingHeadlessSessionsRequests = new Map<string, string>();

export const sessionRouter = Router();

const commonMiddleware = [requestForwarderMiddleware, authMiddleware, trackApiUsage];

/**
 * Create a handshake token for the client to use for secure authentication
 * 
 * @route POST /session-handshake
 * @param {string} x-api-key - [header] API key header
 * @param {string} x-foundry-url - [header] Foundry URL header
 * @param {string} x-world-name - [header,?] World name header
 * @param {string} x-username - [header] Username header
 * @returns {object} Handshake token and encryption details
 */
sessionRouter.post('/session-handshake', authMiddleware, async (req: Request, res: Response) => {
try {
    const apiKey = req.header('x-api-key') as string;
    const foundryUrl = req.header('x-foundry-url') as string;
    const worldName = req.header('x-world-name') as string;
    const username = req.header('x-username') as string;
    
    if (!foundryUrl || !username) {
    res.status(400).json({ error: "Missing required parameters" });
    return;
    }

    // Generate an RSA key pair for this handshake
    const { publicKey, privateKey } = crypto.generateKeyPairSync('rsa', {
    modulusLength: 2048,
    publicKeyEncoding: {
        type: 'spki',
        format: 'pem'
    },
    privateKeyEncoding: {
        type: 'pkcs8',
        format: 'pem'
    }
    });
    
    // Generate a random handshake token that will be valid for 5 minutes
    const handshakeToken = crypto.randomBytes(32).toString('hex');
    const expires = Date.now() + (5 * 60 * 1000); // 5 minutes
    
    const nonce = crypto.randomBytes(16).toString('hex');
    const instanceId = process.env.FLY_ALLOC_ID || 'local';

    // Store handshake in Redis instead of local memory
    const redis = getRedisClient();
    if (redis) {
    // Store all handshake data in Redis with an expiry
    await redis.hSet(`handshake:${handshakeToken}`, {
        apiKey,
        foundryUrl,
        worldName: worldName || '',
        username,
        publicKey,
        privateKey,
        nonce,
        expires: expires.toString(),
        instanceId
    });
    
    // Set expiry for 5 minutes
    await redis.expire(`handshake:${handshakeToken}`, 300);
    
    log.info(`Created handshake token ${handshakeToken.substring(0, 8)}... for ${foundryUrl} in Redis`);
    } else {
    // Fallback to local storage if Redis is unavailable
    pendingHandshakes.set(handshakeToken, {
        apiKey,
        foundryUrl,
        worldName,
        username,
        publicKey,
        privateKey,
        nonce,
        expires
    });
    
    // Set cleanup timeout for local storage
    setTimeout(() => {
        pendingHandshakes.delete(handshakeToken);
        log.debug(`Handshake token ${handshakeToken.substring(0, 8)}... expired and removed from local storage`);
    }, 5 * 60 * 1000);
    
    log.info(`Created handshake token ${handshakeToken.substring(0, 8)}... for ${foundryUrl} in local storage`);
    }
    
    // Return the token and public key to the client
    res.status(200).json({
    token: handshakeToken,
    publicKey: publicKey,
    nonce,
    expires
    });
    return;
} catch (error) {
    log.error(`Error creating handshake: ${error}`);
    res.status(500).json({ error: 'Failed to create handshake' });
    return;
}
});

/**
 * Start a headless Foundry session using puppeteer
 * 
 * @route POST /start-session
 * @param {string} handshakeToken - [body] The token received from session-handshake
 * @param {string} encryptedPassword - [body] Password encrypted with the public key
 * @param {string} x-api-key - [header] API key header
 * @returns {object} Session information including sessionId and clientId
 */
sessionRouter.post("/start-session", requestForwarderMiddleware, authMiddleware, async (req: Request, res: Response) => {
try {
    const { handshakeToken, encryptedPassword } = req.body;
    const apiKey = req.header('x-api-key') as string;
    
    const redis = getRedisClient();
    const isLocalMode = !redis;
    
    // In non-local mode (Redis available), check if an active session already exists for this API key
    // If so, return that session instead of creating a new one
    if (!isLocalMode) {
        const existingSessionId = await redis.get(`headless_apikey:${apiKey}:session`);
        if (existingSessionId) {
            const existingSessionData = await redis.hGetAll(`headless_session:${existingSessionId}`);
            if (existingSessionData && existingSessionData.clientId) {
                log.info(`Returning existing session ${existingSessionId} for API key ${apiKey.substring(0, 8)}...`);
                return safeResponse(res, 200, {
                    success: true,
                    message: "Existing session returned (each API key is limited to one headless session)",
                    sessionId: existingSessionId,
                    clientId: existingSessionData.clientId,
                    existingSession: true
                });
            }
        }
    }

    // Get handshake data from Redis or local storage
    let handshake: any = null;
    let fromRedis = false;
    if (redis) {
    // Try to get handshake from Redis
    const handshakeExists = await redis.exists(`handshake:${handshakeToken}`);
    
    if (handshakeExists) {
        const handshakeData = await redis.hGetAll(`handshake:${handshakeToken}`);
        
        // Check if this instance should handle the request
        const handshakeInstanceId = handshakeData.instanceId;
        const currentInstanceId = process.env.FLY_ALLOC_ID || 'local';
        
        if (handshakeInstanceId !== currentInstanceId) {
        // This should be handled by a different instance
        log.info(`Handshake ${handshakeToken.substring(0, 8)}... belongs to instance ${handshakeInstanceId}, current instance is ${currentInstanceId}`);
        
        // Store the client's request in Redis for the correct instance to pick up
        await redis.hSet(`pending_session:${handshakeToken}`, {
            apiKey,
            encryptedPassword: encryptedPassword,
            timestamp: Date.now().toString()
        });
        
        // Set expiry for 5 minutes
        await redis.expire(`pending_session:${handshakeToken}`, 300);
        
        // Wait for the other instance to process the request and return the result
        log.info(`Waiting for instance ${handshakeInstanceId} to process headless session request`);
        
        // Set a timeout for waiting
        const maxWaitTime = 300000; // 5 minute timeout (matches Redis expiry)
        const startTime = Date.now();
        
        // Poll Redis for the result
        const checkInterval = setInterval(async () => {
            try {
            // Check if the result has been posted back
            const resultKey = `session_result:${handshakeToken}`;
            const hasResult = await redis.exists(resultKey);
            
            if (hasResult) {
            // Get the result data
            const resultData = await redis.get(resultKey);
            await redis.del(resultKey); // Clean up the result
            clearInterval(checkInterval);
            
            // Parse and return the actual response - handle null case with default response
            const result = JSON.parse(resultData || '{"statusCode":200, "data":{"message":"Session started on another instance"}}');
            return safeResponse(res, result.statusCode || 200, result.data || {
            message: "Session started on another instance"
            });
            } else if (Date.now() - startTime > maxWaitTime) {
            // Timeout reached
            clearInterval(checkInterval);
            await redis.del(`pending_session:${handshakeToken}`);
            return safeResponse(res, 408, {
            error: "Timeout waiting for session to be processed by other instance",
            handshakeInstance: handshakeInstanceId
            });
            }
            } catch (err) {
            log.error(`Error polling for session result: ${err}`);
            clearInterval(checkInterval);
            return safeResponse(res, 500, {
            error: "Error while waiting for session to be processed"
            });
            }
        }, 2000); // Check every 2 seconds
        }
        
        // Convert Redis string data to properly typed handshake object
        handshake = {
          ...handshakeData,
          expires: parseInt(handshakeData.expires, 10)
        };
        fromRedis = true;
    }
    }
    
    // If not found in Redis, try local storage
    if (!handshake && pendingHandshakes.has(handshakeToken)) {
    handshake = pendingHandshakes.get(handshakeToken);
    }
    
    // Verify handshake token exists
    if (!handshake) {
    return safeResponse(res, 401, { error: 'Invalid or expired handshake token' });
    }
    
    // Verify API key matches
    if (handshake.apiKey !== apiKey) {
    // Clean up
    if (fromRedis && redis) {
        await redis.del(`handshake:${handshakeToken}`);
    } else {
        pendingHandshakes.delete(handshakeToken);
    }
    
    return safeResponse(res, 401, { error: 'Unauthorized' });
    }
    
    // Verify token is not expired
    if (handshake.expires < Date.now()) {
    // Clean up
    if (fromRedis && redis) {
        await redis.del(`handshake:${handshakeToken}`);
    } else {
        pendingHandshakes.delete(handshakeToken);
    }
    
    return safeResponse(res, 401, { error: 'Handshake token expired' });
    }
    
    // Decrypt the password and nonce using the handshake's private key
    let password;
    let nonce;
    try {
    const buffer = Buffer.from(encryptedPassword, 'base64');
    const decryptedData = crypto.privateDecrypt(
        {
    key: handshake.privateKey,
    padding: crypto.constants.RSA_PKCS1_OAEP_PADDING
        },
        buffer
    ).toString('utf8');

    // Parse the decrypted data as JSON which should contain password and nonce
    const parsedData = JSON.parse(decryptedData);
    password = parsedData.password;
    nonce = parsedData.nonce;
    
    // Verify the nonce matches
    if (!nonce || nonce !== handshake.nonce) {
        if (fromRedis && redis) {
        await redis.del(`handshake:${handshakeToken}`);
        } else {
        pendingHandshakes.delete(handshakeToken);
        }
        res.status(401).json({ error: 'Invalid nonce' });
        return;
    }
    } catch (error) {
    log.error(`Failed to decrypt data: ${error}`);
    if (fromRedis && redis) {
        await redis.del(`handshake:${handshakeToken}`);
    } else {
        pendingHandshakes.delete(handshakeToken);
    }
    res.status(400).json({ error: 'Invalid encrypted data' });
    return;
    }
    
    // Remove the handshake token immediately after use
    const { foundryUrl, worldName, username } = handshake;
    // Remove the handshake token from pending handshakes
    if (fromRedis && redis) {
    await redis.del(`handshake:${handshakeToken}`);
    } else {
    pendingHandshakes.delete(handshakeToken);
    }

    // Launch Puppeteer and connect to Foundry
try {
    log.info(`Starting headless Foundry session for URL: ${foundryUrl}`);
    
    const browser = await puppeteer.launch({
    headless: true,
    executablePath: process.env.PUPPETEER_EXECUTABLE_PATH || undefined,
    args: [
        '--no-sandbox',
        '--disable-setuid-sandbox',
        '--enable-gpu-rasterization', 
        '--enable-oop-rasterization',
        '--disable-dev-shm-usage',
        '--no-first-run',
        '--no-zygote',
        '--disable-extensions',
        '--disable-web-security',
        '--disable-features=site-per-process,IsolateOrigins,site-isolation-trials',
        '--disable-background-networking',
        '--disable-background-timer-throttling',
        '--disable-backgrounding-occluded-windows',
        '--disable-renderer-backgrounding',
        '--disable-sync',
        '--disable-breakpad',
        '--disable-component-extensions-with-background-pages',
        '--disable-default-apps',
        '--disable-infobars',
        '--disable-popup-blocking',
        '--disable-translate',
        '--metrics-recording-only',
        '--mute-audio',
        '--log-level=0',
        '--js-flags="--max_old_space_size=8192"',
    ],
    defaultViewport: { width: 1280, height: 720 }
    });


    const page = await browser.newPage();

    // Enable logging
    page.on('pageerror', (error: unknown) => {
      const message = error instanceof Error ? error.message : String(error);
      log.error(`Browser page error: ${message}`);
    });
    page.on('requestfailed', (request: puppeteer.HTTPRequest) => log.error(`Request failed: ${request.url()}`));
    
    // Navigate to Foundry
    log.debug(`Navigating to Foundry URL: ${foundryUrl}`);
    await page.goto(foundryUrl, { waitUntil: 'networkidle0', timeout: 180000 });
    
    // Debug: Log current URL
    log.debug(`Current page URL: ${page.url()}`);

    // First, check if there are any overlays or tours to dismiss
    log.debug("Checking for overlays or tours to dismiss");
    try {
    // Look for various types of overlays and dismiss them
    const selectors = [
        '.tour-overlay', '.tour', '.tour-fadeout',
        'a.step-button[data-action="exit"]', 'button.tour-exit'
    ];
    
    for (const selector of selectors) {
        const elements = await page.$$(selector);
        if (elements.length > 0) {
        log.debug(`Found ${elements.length} ${selector} elements, attempting to dismiss`);
        await page.click(selector).catch((e: unknown) => {
          const message = e instanceof Error ? e.message : String(e);
          log.debug(`Couldn't click ${selector}: ${message}`);
        });
        await new Promise(resolve => setTimeout(resolve, 1000)); // Wait for a second
        }
    }
    } catch (e) {
    log.info(`Overlay handling: ${(e as Error).message}`);
    }

    // Handle world selection
    if (worldName) {
    log.info(`Looking for world: ${worldName}`);
    
    try {
        // Wait for world list to load
        await page.waitForSelector('li.package.world', { timeout: 10000 })
        .catch(() => {
            log.info('Could not find world list, checking page content');
            return page.content().then((html: string) => {
            log.info(`Page HTML preview: ${html.substring(0, 1000)}...`);
            });
        });
        
        // Try to find and click on the world using multiple strategies
        log.info('Attempting to find and launch the world');
        
        // Strategy 1: Try to find the play button directly associated with the world name
        const worldLaunched = await page.evaluate((worldName: string) => {
        // Find all world titles
        const titles = Array.from(document.querySelectorAll('h3.package-title'));
        
        for (const title of titles) {
            if (title.textContent && title.textContent.trim() === worldName) {
            // Find the parent li element
            const worldLi = title.closest('li.package.world');
            if (worldLi) {
                // Find and click the play button
                const playButton = worldLi.querySelector('a.control.play');
                if (playButton) {
                (playButton as HTMLElement).click();
                return true;
                }
            }
            }
        }
        return false;
        }, worldName);

        await new Promise(resolve => setTimeout(resolve, 2000)); // Give time for action to complete
        
        if (worldLaunched) {
        log.info('World launch button clicked successfully');
        } else {
        log.info('Failed to find/click world launch button');
        
        // Strategy 2: Try using a more direct selector
        try {
            log.info('Trying alternative launch approach');
            
            // Look for all world elements and try to find a match by text content
            const worlds = await page.$$('li.package.world');
            log.info(`Found ${worlds.length} world elements`);
            
            let launched = false;
            for (const worldElement of worlds) {
            const title = await worldElement.$eval('h3.package-title', (el: Element) => el.textContent?.trim())
                .catch(() => null);
                
            log.info(`Found world with title: ${title}`);
            
            if (title === worldName) {
                log.info('Found matching world, looking for play button');
                const playButton = await worldElement.$('a.control.play');
                if (playButton) {
                log.info('Clicking play button');
                await playButton.click();
                launched = true;
                break;
                }
            }
            }
            
            if (!launched) {
            log.info('Failed to launch world using alternative approach');
            }
        } catch (error) {
            log.info(`Error in alternative launch approach: ${(error as Error).message}`);
        }
        }
        
        // Wait and check if we have navigated to a login page
        log.info('Waiting to see if we reached the login page...');
        await new Promise(resolve => setTimeout(resolve, 6000));
        
        // info: Log current URL again
        log.info(`Current URL after world selection: ${page.url()}`);
        
        // Check if we're on a login page by looking for various login elements
        const loginElements = ['select[name="userid"]', 'input[name="userid"]', 'input[name="password"]'];
        let loginFormFound = false;
        
        for (const selector of loginElements) {
        const element = await page.$(selector);
        if (element) {
            log.info(`Found login element: ${selector}`);
            loginFormFound = true;
            break;
        }
        }
        
        if (!loginFormFound) {
        // If we don't see login elements, check the HTML to see what page we're on
        const html = await page.content();
        log.info(`Page HTML after world selection (preview): ${html.substring(0, 500)}...`);
        throw new Error('Login form not found after world selection');
        }
        
    } catch (error) {
        await browser.close();
        const errorMessage = error instanceof Error ? error.message : String(error);
        pendingHeadlessSessionsRequests.delete(apiKey);
        return safeResponse(res, 404, { error: `Failed to find or launch world: ${worldName}`, details: errorMessage });
    }
    }

    // Handle the login process
    log.debug('Attempting to log in...');
    // Handle username input (could be select or input)
    let userId = username; // Default
    let userSelectFound = false;
    let retries = 0;
    const maxRetries = 10;
    const retryInterval = 10000; // 10 seconds between retries
    
    while (!userSelectFound && retries < maxRetries) {
    const hasUserSelect = await page.$('select[name="userid"]')
        .then((element: puppeteer.ElementHandle | null) => !!element)
        .catch(() => false);
    
    if (hasUserSelect) {
        log.debug('Found username dropdown, selecting user');
        userSelectFound = true;
        
        // Get all available users from dropdown
        const options = await page.$$eval('select[name="userid"] option', (options: Element[]) => 
    options.map((opt: any) => ({ value: opt.value, text: opt.textContent?.trim() }))
        );
        
        log.debug(`Available users: ${JSON.stringify(options)}`);
        
        // Find matching username
        const matchingOption = options.find((opt: any) => opt.text === username);
        if (matchingOption) {
    log.info(`Selected user ${username} with value ${matchingOption.value}`);
    await page.select('select[name="userid"]', matchingOption.value);
    userId = matchingOption.value; // Use the value attribute as userId
        } else {
    throw new Error(`Username "${username}" not found in dropdown`);
        }
    } else {
        retries++;
        log.info(`No username dropdown found yet. Attempt ${retries}/${maxRetries}, checking again in ${retryInterval / 1000} seconds...`);
        
        if (retries < maxRetries) {
    await new Promise(resolve => setTimeout(resolve, retryInterval));
        } else {
    log.info('Max retries reached. Assuming direct username input is required.');
    // Try to input username directly if there's an input field
    const hasUserInput = await page.$('input[name="userid"]')
        .then((element: puppeteer.ElementHandle | null) => !!element)
        .catch(() => false);
        
    if (hasUserInput) {
        log.info(`Found username input field, entering username: ${username}`);
        await page.type('input[name="userid"]', username);
    } else {
        log.warn('No username input field found after retries');
    }
        }
    }
    }
    
    // Enter password
    await page.type('input[name="password"]', password);
    
    // Submit form
    log.info('Submitting login form');
    await page.click('button[type="submit"]')
    .catch(() => page.evaluate(() => {
        (document.querySelector('form') as HTMLFormElement)?.submit();
    }));
    
    // Wait for the game to load
    log.info('Waiting for game to load...');
    await page.waitForSelector('#ui-left, #sidebar, .vtt, #game', { timeout: 30000 })
    .catch(async (error: unknown) => {
        const message = error instanceof Error ? error.message : String(error);
        log.error(`Error waiting for game selectors: ${message}`);
        throw error;
    });
    
    // Create a unique session ID and store it
    const sessionId = crypto.randomUUID();
    browserSessions.set(sessionId, browser);
    
    // Register this session in Redis for cross-instance support
    await registerHeadlessSession(sessionId, userId, apiKey);
    
    // The expected client ID will be in format "foundry-{userId}"
    const expectedClientId = `foundry-${userId}`;
    log.info(`Waiting for Foundry client connection with ID: ${expectedClientId}`);
    
    // Track this as a pending session immediately so it can be cleaned up if something goes wrong
    pendingSessions.set(sessionId, {
        sessionId,
        browser,
        apiKey,
        expectedClientId,
        startTime: Date.now(),
        foundryUrl,
        worldName,
        username
    });
    log.info(`Tracking pending session ${sessionId} for client ${expectedClientId}`);
    
    // Create a promise that resolves when the client connects or rejects on timeout
    const clientConnectionPromise = new Promise<string>((resolve, reject) => {
    // Initial check for existing client
    const checkExistingClient = async () => {
        const client = await ClientManager.getClient(expectedClientId);
        if (client && client.getApiKey() === apiKey) {
    return expectedClientId;
        } else if (client) {
    // If the client ID matches but the API key doesn't, log a warning
    log.warn(`Client ID ${expectedClientId} found but API key mismatch`);
    return 'invalid';
        }
        return null;
    };
    
    // Set up polling for client connection with reduced verbosity
    let logCounter = 0;
    const checkInterval = setInterval(async () => {
        try {
    const clientId = await checkExistingClient();
    if (clientId) {
        // Only log the connection once
        if (clientId === 'invalid') {
        // close the browser session
        await browser.close();
        browserSessions.delete(sessionId);
        clearInterval(checkInterval);
        clearTimeout(timeoutId);
        reject(new Error(`Unauthorized client connection attempt`));
        return;
        }
        log.info(`Client connected successfully: ${clientId}`);
        clearInterval(checkInterval);
        clearTimeout(timeoutId);
        resolve(clientId);
    } else {
        // Log less frequently to reduce noise
        if (++logCounter % 10 === 0) {
        log.debug(`Waiting for client connection: ${expectedClientId} (${logCounter} checks)`);
        }
    }
        } catch (error) {
    log.error(`Error checking for client: ${error}`);
        }
    }, 2000);
    
    // Set timeout for client connection
    const timeoutId = setTimeout(() => {
        clearInterval(checkInterval);
        reject(new Error(`Timeout waiting for client connection: ${expectedClientId}`));
    }, 300000); // Wait up to 5 minutes for the client to connect
    });
    
    try {
    // Wait for client connection
    const connectedClientId = await clientConnectionPromise;
    
    // Remove from pending sessions - connection successful
    pendingSessions.delete(sessionId);
    log.info(`Session ${sessionId} connected successfully, removed from pending`);
    
    // Store the session in our API key mapping
    apiKeyToSession.set(apiKey, { 
        sessionId, 
        clientId: connectedClientId,
        lastActivity: Date.now()
    });
    
    // Return success with the session ID and client ID
    pendingHeadlessSessionsRequests.delete(apiKey);
    return safeResponse(res, 200, {
        success: true,
        message: "Foundry session started successfully",
        sessionId,
        clientId: connectedClientId
    });
    } catch (error) {
    // Remove from pending sessions - connection failed
    pendingSessions.delete(sessionId);
    log.info(`Session ${sessionId} connection failed, removed from pending`);
    
    // Close the browser if client connection times out
    await browser.close();
    browserSessions.delete(sessionId);
    
    const errorMessage = error instanceof Error ? error.message : String(error);
    pendingHeadlessSessionsRequests.delete(apiKey);
    return safeResponse(res, 408, { 
        error: "Client connection timeout", 
        details: errorMessage,
        message: "Foundry client failed to connect to the API within the 5 minute timeout period. " +
                 "This usually happens when another GM is already logged in and is elected as the 'primary GM'. " +
                 "The Foundry REST API module only allows one GM to connect to the relay at a time. " +
                 "Try logging out other GM users before starting a headless session.",
        expectedClientId
    });
    }
    } catch (error) {
    log.error(`Error starting headless Foundry session: ${error}`);
    const errorMessage = error instanceof Error ? error.message : String(error);
    return safeResponse(res, 500, { error: "Failed to start headless Foundry session", details: errorMessage });
}
} catch (error) {
    log.error(`Error in start-session handler: ${error}`);
    return safeResponse(res, 500, { error: "Internal server error" });
}
});

/**
 * Stop a headless Foundry session
 * 
 * @route DELETE /end-session
 * @param {string} sessionId - [query] The ID of the session to end
 * @param {string} x-api-key - [header] API key header
 * @returns {object} Status of the operation
 */
sessionRouter.delete("/end-session", requestForwarderMiddleware, authMiddleware, async (req: Request, res: Response) => {
try {
    const sessionId = req.query.sessionId as string;
    const apiKey = req.header('x-api-key') as string;
    
    if (!sessionId) {
    return safeResponse(res, 400, { error: "Session ID is required" });
    }
    
    // First check if this is a pending session (browser launched but waiting for connection)
    const pendingSession = pendingSessions.get(sessionId);
    if (pendingSession) {
        if (pendingSession.apiKey !== apiKey) {
            return safeResponse(res, 403, { error: "Not authorized to terminate this session" });
        }
        
        log.info(`Terminating pending session ${sessionId} (was waiting for ${pendingSession.expectedClientId})`);
        
        try {
            await pendingSession.browser.close();
            pendingSessions.delete(sessionId);
            browserSessions.delete(sessionId);
            
            return safeResponse(res, 200, {
                success: true,
                message: "Pending Foundry session terminated",
                details: `Session was waiting for client ${pendingSession.expectedClientId} to connect`
            });
        } catch (error) {
            log.error(`Error closing pending session ${sessionId}: ${error}`);
            pendingSessions.delete(sessionId);
            return safeResponse(res, 500, { error: "Failed to close pending session browser" });
        }
    }
    
    // Check if we have this session locally (established session)
    const browser = browserSessions.get(sessionId);
    let sessionClosed = false;
    
    // Try to close browser if we have it locally
    if (browser) {
    try {
        await browser.close();
        browserSessions.delete(sessionId);
        sessionClosed = true;
        log.info(`Closed browser for session ${sessionId} locally`);
    } catch (error) {
        log.error(`Failed to close browser: ${error}`);
    }
    }
    
    // Clean up session data in Redis regardless
    try {
    const redis = getRedisClient();
    if (redis) {
        // Get session data to find associated client
        const sessionData = await redis.hGetAll(`headless_session:${sessionId}`);
        
        if (sessionData && sessionData.apiKey === apiKey) {
        // Delete all session-related keys
        if (sessionData.clientId) {
            await redis.del(`headless_client:${sessionData.clientId}`);
            // Also remove from ClientManager
            ClientManager.removeClient(sessionData.clientId);
            log.info(`Removed client ${sessionData.clientId} from ClientManager`);
        }
        await redis.del(`headless_apikey:${apiKey}`);
        await redis.del(`headless_session:${sessionId}`);
        
        log.info(`Cleaned up Redis data for session ${sessionId}`);
        return safeResponse(res, 200, { 
            success: true, 
            message: sessionClosed ? "Foundry session terminated" : "Foundry session data cleaned up" 
        });
        } else {
        return safeResponse(res, 403, { error: "Not authorized to terminate this session" });
        }
    }
    } catch (error) {
    log.error(`Error cleaning up Redis session data: ${error}`);
    }
    
    // Clean up local session storage as well
    // First check if it exists, then delete
    const localSessions = apiKeyToSession.getAll(apiKey);
    const localSession = localSessions.find(s => s.sessionId === sessionId);
    
    // Remove client from ClientManager if found
    if (localSession && localSession.clientId) {
        ClientManager.removeClient(localSession.clientId);
        log.info(`Removed client ${localSession.clientId} from ClientManager`);
    }
    
    apiKeyToSession.deleteSession(apiKey, sessionId);
    
    // If we got here with sessionClosed true, we closed the browser but failed Redis cleanup
    if (sessionClosed) {
    return safeResponse(res, 200, { success: true, message: "Foundry session terminated (partial cleanup)" });
    }
    
    // Check if session existed in local storage (for local mode without Redis)
    if (localSession) {
        return safeResponse(res, 200, { success: true, message: "Foundry session terminated" });
    }
    
    return safeResponse(res, 404, { error: "Session not found" });
} catch (error) {
    log.error(`Error in end-session handler: ${error}`);
    return safeResponse(res, 500, { error: "Internal server error" });
}
});

/**
 * Get all active headless Foundry sessions
 * 
 * @route GET /session
 * @param {string} x-api-key - [header] API key header
 * @returns {object} List of active sessions for the current API key
 */
sessionRouter.get("/session", requestForwarderMiddleware, authMiddleware, async (req: Request, res: Response) => {
    try {
        const apiKey = req.header('x-api-key') as string;
        const redis = getRedisClient();
        let sessions: any[] = [];
        
        // Helper function to get client details
        const getClientDetails = async (clientId: string) => {
            if (redis) {
                // Try to get from Redis first
                const [worldId, worldTitle, foundryVersion, systemId, systemTitle, systemVersion, customName] = await Promise.all([
                    redis.get(`client:${clientId}:worldId`),
                    redis.get(`client:${clientId}:worldTitle`),
                    redis.get(`client:${clientId}:foundryVersion`),
                    redis.get(`client:${clientId}:systemId`),
                    redis.get(`client:${clientId}:systemTitle`),
                    redis.get(`client:${clientId}:systemVersion`),
                    redis.get(`client:${clientId}:customName`)
                ]);
                return { worldId, worldTitle, foundryVersion, systemId, systemTitle, systemVersion, customName };
            }
            
            // Fallback to ClientManager
            const client = await ClientManager.getClient(clientId);
            if (client) {
                return {
                    worldId: client.getWorldId() || '',
                    worldTitle: client.getWorldTitle() || '',
                    foundryVersion: client.getFoundryVersion() || '',
                    systemId: client.getSystemId() || '',
                    systemTitle: client.getSystemTitle() || '',
                    systemVersion: client.getSystemVersion() || '',
                    customName: client.getCustomName() || ''
                };
            }
            return null;
        };
        
        // Try to get session data from Redis first
        if (redis) {
        // Check if this API key has a headless session in Redis - FIX: Use correct key pattern
        const sessionId = await redis.get(`headless_apikey:${apiKey}:session`);
        
        if (sessionId) {
            // Get full session details
            const sessionData = await redis.hGetAll(`headless_session:${sessionId}`);
            
            if (sessionData) {
            // Parse timestamps
            const lastActivity = parseInt(sessionData.lastActivity || '0');
            
            // Get client details
            const clientDetails = sessionData.clientId ? await getClientDetails(sessionData.clientId) : null;
            
            sessions.push({
                id: sessionId,
                clientId: sessionData.clientId || '',
                lastActivity: lastActivity,
                idleMinutes: Math.round((Date.now() - lastActivity) / 60000),
                instanceId: sessionData.instanceId || 'unknown',
                // Include client details if available
                ...(clientDetails && {
                    worldId: clientDetails.worldId || '',
                    worldTitle: clientDetails.worldTitle || '',
                    foundryVersion: clientDetails.foundryVersion || '',
                    systemId: clientDetails.systemId || '',
                    systemTitle: clientDetails.systemTitle || '',
                    systemVersion: clientDetails.systemVersion || '',
                    customName: clientDetails.customName || ''
                })
            });
            }
        }
        }
        
        // Fall back to local storage if no Redis session found
        if (sessions.length === 0) {
        // Get ALL sessions for this API key (supports multiple sessions in local mode)
        const userSessions = apiKeyToSession.getAll(apiKey);
        
        for (const userSession of userSessions) {
            // Get client details
            const clientDetails = await getClientDetails(userSession.clientId);
            
            sessions.push({
            id: userSession.sessionId,
            clientId: userSession.clientId,
            lastActivity: userSession.lastActivity,
            idleMinutes: Math.round((Date.now() - userSession.lastActivity) / 60000),
            instanceId: process.env.FLY_ALLOC_ID || 'local',
            // Include client details if available
            ...(clientDetails && {
                worldId: clientDetails.worldId || '',
                worldTitle: clientDetails.worldTitle || '',
                foundryVersion: clientDetails.foundryVersion || '',
                systemId: clientDetails.systemId || '',
                systemTitle: clientDetails.systemTitle || '',
                systemVersion: clientDetails.systemVersion || '',
                customName: clientDetails.customName || ''
            })
            });
        }
        }
        
        // Also include pending sessions (browsers launched but waiting for client connection)
        // These are critical to track because they can become orphaned
        const pendingForApiKey: any[] = [];
        for (const [sessionId, pending] of pendingSessions.entries()) {
            if (pending.apiKey === apiKey) {
                const waitingSeconds = Math.round((Date.now() - pending.startTime) / 1000);
                pendingForApiKey.push({
                    id: sessionId,
                    status: 'pending',
                    expectedClientId: pending.expectedClientId,
                    foundryUrl: pending.foundryUrl,
                    worldName: pending.worldName || '',
                    username: pending.username,
                    waitingSeconds,
                    startTime: pending.startTime,
                    instanceId: process.env.FLY_ALLOC_ID || 'local'
                });
            }
        }
        
        safeResponse(res, 200, { 
        activeSessions: sessions,
        pendingSessions: pendingForApiKey
        });
    } catch (error) {
        log.error(`Error retrieving headless sessions: ${error}`);
        safeResponse(res, 500, { error: "Failed to retrieve session data" });
    }
});