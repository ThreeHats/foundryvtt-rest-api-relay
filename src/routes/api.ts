import express, { Request, Response, NextFunction } from "express";
import path from "path";
// import { log } from "../middleware/logger";
import { ClientManager } from "../core/ClientManager";
import { Client } from "../core/Client"; // Import Client type
import axios from 'axios';
import { PassThrough } from 'stream';
import { JSDOM } from 'jsdom';
import { authMiddleware, trackApiUsage } from '../middleware/auth';
import { requestForwarderMiddleware } from '../middleware/requestForwarder';
import { pendingRequests, PENDING_REQUEST_TYPES, safeResponse } from './shared';
import { dnd5eRouter } from './api/dnd5e';
import { healthCheck } from '../routes/health';
import { getRedisClient } from '../config/redis';
import { returnHtmlTemplate } from "../config/htmlResponseTemplate";
import multer from "multer";
import fs from "fs/promises";
import { searchRouter } from './api/search';
import { entityRouter } from './api/entity';
import { rollRouter } from './api/roll';
import { utilityRouter } from './api/utility';
import { fileSystemRouter } from './api/fileSystem';
import { sessionRouter } from './api/session';
import { encounterRouter } from './api/encounter';
import { sheetRouter } from './api/sheet';
import { macroRouter } from './api/macro';
import { structureRouter } from './api/structure';
import { log } from '../utils/logger';
// Import from session.ts instead of duplicating
import { browserSessions, apiKeyToSession, pendingSessions } from './api/session';

export { browserSessions, apiKeyToSession, pendingSessions };

export const VERSION = '2.1.1';

const INSTANCE_ID = process.env.INSTANCE_ID || 'default';

const HEADLESS_SESSION_TIMEOUT = 10 * 60 * 1000; // 10 minutes in milliseconds
const PENDING_SESSION_TIMEOUT = 5 * 60 * 1000; // 5 minutes - matches the clientConnectionPromise timeout

async function cleanupInactiveSessions() {
  const now = Date.now();
  
  for (const [apiKey, session] of apiKeyToSession.entries()) {
    if (now - session.lastActivity > HEADLESS_SESSION_TIMEOUT) {
      log.info(`Closing inactive headless session ${session.sessionId} for API key ${apiKey.substring(0, 8)}... (inactive for ${Math.round((now - session.lastActivity) / 60000)} minutes)`);
      
      try {
        // Close the browser if it exists
        if (browserSessions.has(session.sessionId)) {
          const browser = browserSessions.get(session.sessionId);
          browser?.close().catch(err => log.error(`Error closing browser: ${err}`));
          browserSessions.delete(session.sessionId);
        }
        
        // Also remove the associated client from ClientManager
        if (session.clientId) {
          await ClientManager.removeClient(session.clientId);
          log.info(`Removed stale client ${session.clientId} during session cleanup`);
        }
        
        // Clean up the session mapping
        apiKeyToSession.deleteSession(apiKey, session.sessionId);
      } catch (error) {
        log.error(`Error during inactive session cleanup: ${error}`);
      }
    }
  }
}

/**
 * Cleanup orphaned pending sessions - browsers that were launched but never received
 * a client connection. This can happen if:
 * - The Foundry module's "primary GM" logic prevents the headless client from connecting
 * - Network issues prevent the WebSocket connection
 * - The test/process was interrupted before the timeout handler could run
 */
async function cleanupOrphanedPendingSessions() {
  const now = Date.now();
  
  for (const [sessionId, pending] of pendingSessions.entries()) {
    const waitingTime = now - pending.startTime;
    
    // If waiting longer than the timeout, this session is orphaned
    if (waitingTime > PENDING_SESSION_TIMEOUT) {
      const waitingMinutes = Math.round(waitingTime / 60000);
      log.warn(`Found orphaned pending session ${sessionId} - waiting for ${pending.expectedClientId} for ${waitingMinutes} minutes. Cleaning up...`);
      
      try {
        // Close the browser
        await pending.browser.close();
        log.info(`Closed orphaned browser for session ${sessionId}`);
        
        // Remove from pending sessions
        pendingSessions.delete(sessionId);
        
        // Also remove from browserSessions if it was added there
        browserSessions.delete(sessionId);
        
        log.info(`Cleaned up orphaned pending session ${sessionId} (was waiting for client ${pending.expectedClientId})`);
      } catch (error) {
        log.error(`Error cleaning up orphaned pending session ${sessionId}: ${error}`);
        // Still remove from tracking even if browser close fails
        pendingSessions.delete(sessionId);
      }
    }
  }
}

// Start the session cleanup intervals when the module is loaded
setInterval(cleanupInactiveSessions, 60000); // Check every minute for inactive established sessions
setInterval(cleanupOrphanedPendingSessions, 30000); // Check every 30 seconds for orphaned pending sessions

export const apiRoutes = (app: express.Application): void => {
  // Setup handlers for storing search results and entity data from WebSocket
  setupMessageHandlers();
  
  // Create a router instead of using app directly
  const router = express.Router();

  // Define routes on the router
  router.get("/", (req: Request, res: Response) => {
    res.sendFile(path.join(__dirname, "../../_test/test-client.html"));
  });

  router.get("/health", healthCheck);

  router.get("/api/status", (req: Request, res: Response) => {
    res.json({ 
      status: "ok",
      version: VERSION,
      websocket: "/relay"
    });
  });

  // Get all connected clients
  router.get("/clients", authMiddleware, async (req: Request, res: Response) => {
    try {
      const apiKey = req.header('x-api-key') || '';
      const redis = getRedisClient();
      
      // Array to store all client details
      let allClients: any[] = [];
      
      if (redis) {
        // Step 1: Get all client IDs from Redis for this API key
        const clientIds = await redis.sMembers(`apikey:${apiKey}:clients`);
        
        if (clientIds.length > 0) {
          // Step 2: For each client ID, get details from Redis
          const clientDetailsPromises = clientIds.map(async (clientId) => {
            try {
              // Get the instance this client is connected to
              const instanceId = await redis.get(`client:${clientId}:instance`);
              
              if (!instanceId) return null;
              
              // Get the last seen timestamp if stored
              const lastSeen = await redis.get(`client:${clientId}:lastSeen`) || Date.now();
              const connectedSince = await redis.get(`client:${clientId}:connectedSince`) || lastSeen;
              
              // Return client details including its instance
              return {
                id: clientId,
                instanceId,
                lastSeen: parseInt(lastSeen.toString()),
                connectedSince: parseInt(connectedSince.toString()),
                worldId: await redis.get(`client:${clientId}:worldId`) || '',
                worldTitle: await redis.get(`client:${clientId}:worldTitle`) || '',
                foundryVersion: await redis.get(`client:${clientId}:foundryVersion`) || '',
                systemId: await redis.get(`client:${clientId}:systemId`) || '',
                systemTitle: await redis.get(`client:${clientId}:systemTitle`) || '',
                systemVersion: await redis.get(`client:${clientId}:systemVersion`) || '',
                customName: await redis.get(`client:${clientId}:customName`) || ''
              };
            } catch (err) {
              log.error(`Error getting details for client ${clientId}: ${err}`);
              return null;
            }
          });
          
          // Resolve all promises and filter out nulls
          const clientDetails = (await Promise.all(clientDetailsPromises)).filter(client => client !== null);
          allClients = clientDetails;
        }
      } else {
        // Fallback to local clients if Redis isn't available
        const localClientIds = await ClientManager.getConnectedClients(apiKey);
        
        // Use Promise.all to wait for all getClient calls to complete
        allClients = await Promise.all(localClientIds.map(async (id) => {
          const client = await ClientManager.getClient(id);
          return {
            id,
            instanceId: INSTANCE_ID,
            lastSeen: client?.getLastSeen() || Date.now(),
            connectedSince: client?.getLastSeen() || Date.now(),
            worldId: client?.getWorldId() || '',
            worldTitle: client?.getWorldTitle() || '',
            foundryVersion: client?.getFoundryVersion() || '',
            systemId: client?.getSystemId() || '',
            systemTitle: client?.getSystemTitle() || '',
            systemVersion: client?.getSystemVersion() || '',
            customName: client?.getCustomName() || ''
          };
        }));
      }
      
      // Send combined response
      safeResponse(res, 200, {
        total: allClients.length,
        clients: allClients
      });
    } catch (error) {
      log.error(`Error aggregating clients: ${error}`);
      safeResponse(res, 500, { error: "Failed to retrieve clients" });
    }
  });
  
  // Proxy asset requests to Foundry
  router.get('/proxy-asset/:path(*)', requestForwarderMiddleware, async (req: Request, res: Response) => {
    try {
      // Get the Foundry URL from the client metadata or use default
      const clientId = req.query.clientId as string;
      let foundryBaseUrl = 'http://localhost:30000'; // Default Foundry URL
      
      // If we have client info, use its URL
      if (clientId) {
        const client = await ClientManager.getClient(clientId);
        if (client && 'metadata' in client && client.metadata && (client.metadata as any).origin) {
          foundryBaseUrl = (client.metadata as any).origin;
        }
      }
      
      const assetPath = req.params.path;
      const assetUrl = `${foundryBaseUrl}/${assetPath}`;
      
      log.debug(`Proxying asset request to: ${assetUrl}`);
      
      // Check if it's a Font Awesome file - redirect to CDN if so
      if (assetPath.includes('/webfonts/fa-') || assetPath.includes('/fonts/fontawesome/') || 
          assetPath.includes('/fonts/fa-')) {
        
        // Extract the filename
        const filename = assetPath.split('/').pop() || '';
        
        // Redirect to CDN
        const cdnUrl = `https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.2/webfonts/${filename}`;
        log.debug(`Redirecting Font Awesome asset to CDN: ${cdnUrl}`);
        res.redirect(cdnUrl);
        return;
      }
      
      // Handle The Forge specific assets
      if (assetPath.includes('forgevtt-module.css') || assetPath.includes('forge-vtt.com')) {
        log.debug(`Skipping The Forge asset: ${assetPath}`);
        // Return an empty CSS file for Forge assets to prevent errors
        if (assetPath.endsWith('.css')) {
          res.type('text/css').send('/* Placeholder for The Forge CSS */');
          return;
        } else if (assetPath.endsWith('.js')) {
          res.type('application/javascript').send('// Placeholder for The Forge JS');
          return;
        } else {
          // Return a transparent 1x1 pixel for images
          res.type('image/png').send(Buffer.from('iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNkYAAAAAYAAjCB0C8AAAAASUVORK5CYII=', 'base64'));
          return;
        }
      }
      
      // Check for texture files - use GitHub raw content as fallback
      if (assetPath.includes('texture1.webp') || assetPath.includes('texture2.webp') || 
          assetPath.includes('parchment.jpg')) {
        log.debug(`Serving texture file from GitHub fallback`);
        res.redirect('https://raw.githubusercontent.com/foundryvtt/dnd5e/master/ui/parchment.jpg');
        return;
      }
      
      // Additional asset fallbacks...
      
      // Try to make the request to Foundry with better error handling
      try {
        const response = await axios({
          method: 'get',
          url: assetUrl,
          responseType: 'stream',
          timeout: 30000, // Increased timeout to 30s
          maxRedirects: 5,
          validateStatus: (status) => status < 500 // Only treat 500+ errors as errors
        });
        
        // Copy headers
        Object.keys(response.headers).forEach(header => {
          res.setHeader(header, response.headers[header]);
        });
        
        // Set CORS headers
        res.setHeader('Access-Control-Allow-Origin', '*');
        
        // Stream the response
        response.data.pipe(res);
      } catch (error) {
        log.error(`Request failed: ${assetUrl}`);
        
        // For CSS files, return an empty CSS file
        if (assetPath.endsWith('.css')) {
          res.type('text/css').send('/* CSS not available */');
        } else if (assetPath.endsWith('.js')) {
          res.type('application/javascript').send('// JavaScript not available');
        } else {
          // Return a transparent 1x1 pixel for images and other files
          res.type('image/png').send(Buffer.from('iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNkYAAAAAYAAjCB0C8AAAAASUVORK5CYII=', 'base64'));
        }
      }
    } catch (error) {
      log.error(`Error in proxy asset handler: ${error}`);
      res.status(404).send('Asset not found');
    }
  });

  // API Documentation endpoint - returns all available endpoints with their documentation
  router.get("/api/docs", async (req: Request, res: Response) => {
    try {
        // Try multiple possible paths for the API docs file
        const possiblePaths = [
          path.resolve(__dirname, '../../../public/api-docs.json'),  // Development path
          path.resolve(__dirname, '../../public/api-docs.json'),     // Alternative path
          path.resolve(process.cwd(), 'public/api-docs.json'),       // Production path from app root
          path.resolve(process.cwd(), 'dist/public/api-docs.json'),  // If public is copied to dist
        ];
        
        let docsContent: string | null = null;
        let usedPath: string | null = null;
        
        // Try each path until we find the file
        for (const docsPath of possiblePaths) {
          try {
            docsContent = await fs.readFile(docsPath, 'utf8');
            usedPath = docsPath;
            break;
          } catch (err) {
            // File not found at this path, try next one
            log.debug(`API docs not found at: ${docsPath}`);
          }
        }
        
        if (!docsContent) {
          throw new Error(`API docs file not found at any of the expected paths: ${possiblePaths.join(', ')}`);
        }
        
        log.debug(`Successfully loaded API docs from: ${usedPath}`);
        const apiDocs = JSON.parse(docsContent);

        // Dynamically set the baseUrl
        apiDocs.baseUrl = `${req.protocol}://${req.get('host')}`;

        res.json(apiDocs);
    } catch (error) {
        log.error('Failed to load API documentation:', { 
          error: error instanceof Error ? error.message : String(error),
          cwd: process.cwd(),
          __dirname: __dirname
        });
        
        // Provide a basic fallback response
        res.status(500).json({ 
          error: 'API documentation is currently unavailable.',
          message: 'The documentation file could not be loaded. Please check if the server was built correctly.',
          baseUrl: `${req.protocol}://${req.get('host')}`
        });
    }
  });

  // Mount the router
  app.use("/", router);
  app.use('/', searchRouter);
  app.use('/', entityRouter);
  app.use('/', rollRouter);
  app.use('/', utilityRouter);
  app.use('/', fileSystemRouter);
  app.use('/', sessionRouter);
  app.use('/', encounterRouter);
  app.use('/', sheetRouter);
  app.use('/', macroRouter);
  app.use('/', structureRouter);
  app.use('/dnd5e', dnd5eRouter);
};

const REQUEST_TYPES_WITH_SPECIAL_RESPONSE_HANDLERS = [
  'actor-sheet', 'download-file'
] as const;

// Setup WebSocket message handlers to route responses back to API requests
function setupMessageHandlers() {
  
  for (const type of PENDING_REQUEST_TYPES) {
    if (REQUEST_TYPES_WITH_SPECIAL_RESPONSE_HANDLERS.includes(type as (typeof REQUEST_TYPES_WITH_SPECIAL_RESPONSE_HANDLERS)[number])) {
      continue;
    }

    ClientManager.onMessageType(`${type}-result`, (client: Client, data: any) => {
      log.info(`Received ${type} response for requestId: ${data.requestId}`);

      if (data.requestId && pendingRequests.has(data.requestId)) {
        const pending = pendingRequests.get(data.requestId);
        if (!pending) {
          log.warn(`Pending request ${data.requestId} was deleted before processing`);
          return;
        }
        
        const response: Record<string, any> = { 
          requestId: data.requestId, 
          clientId: pending.clientId || client.getId() 
        };
        for (const [key, value] of Object.entries(data)) {
          if (key !== 'requestId') {
            response[key] = value;
          }
        }
        if (response.error) {
          safeResponse(pending.res, 400, response);
        } else {
          safeResponse(pending.res, 200, response);
        }
        pendingRequests.delete(data.requestId);
        return;
      }
    });
  }

  // Handler for actor sheet HTML response
  ClientManager.onMessageType("get-sheet-response", (client: Client, data: any) => {
    log.info(`Received actor sheet HTML response for requestId: ${data.requestId}`);
    
    try {
      // Extract the UUID from either data.uuid or data.data.uuid
      const responseUuid = data.uuid || (data.data && data.data.uuid);
      
      // Debug what we're receiving
      log.debug(`Actor sheet response data structure:`, {
        requestId: data.requestId,
        uuid: responseUuid,
        dataKeys: data.data ? Object.keys(data.data) : [],
        html: data.data && data.data.html ? `${data.data.html.substring(0, 100)}...` : undefined,
        cssLength: data.data && data.data.css ? data.data.css.length : 0
      });
      
      if (data.requestId && pendingRequests.has(data.requestId)) {
        const pending = pendingRequests.get(data.requestId)!;
        
        // Check if this is an error response (handle even if UUID is undefined)
        const hasError = data.error || (data.data && data.data.error);
        
        // Compare with either location, but allow error responses even if UUID is undefined
        if (pending.type === 'get-sheet' && (pending.uuid === responseUuid || (hasError && responseUuid === undefined))) {
          if (hasError) {
            const errorMsg = data.error || (data.data && data.data.error) || "Unknown error";
            safeResponse(pending.res, 400, {
              requestId: data.requestId,
              clientId: pending.clientId,
              uuid: pending.uuid,
              error: errorMsg
            });
          } else {
            // Get HTML content from either data or data.data
            let html = data.html || (data.data && data.data.html) || '';
            const css = data.css || (data.data && data.data.css) || '';
            
            // Get the system ID for use in the HTML output
            const gameSystemId = (client as any).metadata?.systemId || 'unknown';
            
            if (pending.format === 'json') {
              // Send response as JSON
              safeResponse(pending.res, 200, {
                requestId: data.requestId,
                clientId: pending.clientId,
                uuid: pending.uuid,
                html: html,
                css: css
              });
            } else {
              // Get the scale and tab parameters from the pending request
              const initialScale = pending.initialScale || null;
              // Convert activeTab to a number if it exists, or keep as null
              const activeTabIndex = pending.activeTab !== null ? Number(pending.activeTab) : null;

              // If a specific tab index is requested, pre-process the HTML to activate that tab
              if (activeTabIndex !== null && !isNaN(activeTabIndex)) {
              try {
                // Create a virtual DOM to manipulate the HTML
                const dom = new JSDOM(html);
                const document = dom.window.document;
                
                // Find all tab navigation elements
                const tabsElements = document.querySelectorAll('nav.tabs, .tabs');
                
                tabsElements.forEach(tabsElement => {
                // Find all tab items and content tabs
                const tabs = Array.from(tabsElement.querySelectorAll('.item'));
                const sheet = tabsElement.closest('.sheet');
                
                if (sheet && tabs.length > 0 && activeTabIndex < tabs.length) {
                  const tabContent = sheet.querySelectorAll('.tab');
                  
                  if (tabs.length > 0 && tabContent.length > 0) {
                  // Deactivate all tabs first
                  tabs.forEach(tab => tab.classList.remove('active'));
                  tabContent.forEach(content => content.classList.remove('active'));
                  
                  // Get the tab at the specified index
                  const targetTab = tabs[activeTabIndex];
                  
                  if (targetTab) {
                    // Get the data-tab attribute from this tab
                    const tabName = targetTab.getAttribute('data-tab');
                    
                    // Find the corresponding content tab
                    let targetContent = null;
                    for (let i = 0; i < tabContent.length; i++) {
                    if (tabContent[i].getAttribute('data-tab') === tabName) {
                      targetContent = tabContent[i];
                      break;
                    }
                    }
                    
                    // Activate both the tab and its content
                    targetTab.classList.add('active');
                    if (targetContent) {
                    targetContent.classList.add('active');
                    log.debug(`Pre-activated tab index ${activeTabIndex} with data-tab: ${tabName}`);
                    }
                  }
                  }
                }
                });
                
                // Get the modified HTML
                html = document.querySelector('body')?.innerHTML || html;
                }
              catch (error) {
                log.warn(`Failed to pre-process HTML for tab selection: ${error}`);
                // Continue with the original HTML if there was an error
              }}

              // If dark mode is requested, flag it for later use in the full HTML document
              const darkModeEnabled = pending.darkMode || false;

              // Determine if we should include interactive JavaScript
              const includeInteractiveJS = initialScale === null && activeTabIndex === null;

              // Generate the full HTML document
              const fullHtml = returnHtmlTemplate(responseUuid, html, css, gameSystemId, darkModeEnabled, includeInteractiveJS, activeTabIndex || 0, initialScale || 0, pending);
              
              pending.res.send(fullHtml);
            }
          }
          
          // Remove pending request
          pendingRequests.delete(data.requestId);
        } else {
          // Log an issue if UUID doesn't match what we expect
          log.warn(`Received actor sheet response with mismatched values: expected type=${pending.type}, uuid=${pending.uuid}, got uuid=${responseUuid}`);
        }
      } else {
        log.warn(`Received actor sheet response for unknown requestId: ${data.requestId}`);
      }
    } catch (error) {
      log.error(`Error handling actor sheet HTML response:`, { error });
      log.debug(`Response data that caused the error:`, {
        requestId: data.requestId,
        hasData: !!data.data,
        dataType: typeof data.data
      });
    }
  });

  // Handler for file download result
  ClientManager.onMessageType("download-file-result", (client: Client, data: any) => {
    log.info(`Received file download result for requestId: ${data.requestId}`);
    
    if (data.requestId && pendingRequests.has(data.requestId)) {
      const request = pendingRequests.get(data.requestId)!;
      pendingRequests.delete(data.requestId);
      
      if (data.error) {
        safeResponse(request.res, 500, { 
          clientId: client.getId(),
          requestId: data.requestId,
          error: data.error
        });
        return;
      }
      
      // Check if the client wants raw binary data or JSON response
      const format = request.format || 'binary'; // Default to binary format
      
      if (format === 'binary' || format === 'raw') {
        // Extract the base64 data and send as binary
        const base64Data = data.fileData.split(',')[1];
        const buffer = Buffer.from(base64Data, 'base64');
        
        // Set the appropriate content type
        request.res.setHeader('Content-Type', data.mimeType || 'application/octet-stream');
        request.res.setHeader('Content-Disposition', `attachment; filename="${data.filename}"`);
        request.res.setHeader('Content-Length', buffer.length);
        
        // Send the binary data
        request.res.status(200).end(buffer);
      } else {
        // Send JSON response with the file data
        safeResponse(request.res, 200, {
          clientId: client.getId(),
          requestId: data.requestId,
          success: true,
          path: data.path,
          filename: data.filename,
          mimeType: data.mimeType,
          fileData: data.fileData,
          size: Buffer.from(data.fileData.split(',')[1], 'base64').length
        });
      }
    }
  });

  // Clean up old pending requests periodically
  setInterval(() => {
    const now = Date.now();
    for (const [requestId, request] of pendingRequests.entries()) {
      // Remove requests older than 30 seconds
      if (now - request.timestamp > 30000) {
        log.warn(`Request ${requestId} timed out and was never completed`);
        pendingRequests.delete(requestId);
      }
    }
  }, 10000);
}
