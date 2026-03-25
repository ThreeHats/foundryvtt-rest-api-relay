/**
 * Standalone utility for starting headless sessions from stored credentials.
 * Extracted to avoid circular dependencies (session.ts ↔ auth.ts ↔ api.ts).
 */
import crypto from 'crypto';
import { ClientManager } from '../core/ClientManager';
import { ApiKey } from '../models/apiKey';
import { decrypt, isEncryptionAvailable } from './encryption';
import { log } from './logger';
import { registerHeadlessSession } from '../workers/headlessSessions';
import { attachBrowserConsoleLogger } from './browserConsoleLogger';

// Import session tracking maps lazily to avoid circular deps
let _browserSessions: Map<string, any> | null = null;
let _pendingSessions: Map<string, any> | null = null;
let _apiKeyToSession: any | null = null;

async function getSessionMaps() {
    if (!_browserSessions) {
        const sessionModule = await import('../routes/api/session');
        _browserSessions = sessionModule.browserSessions;
        _pendingSessions = sessionModule.pendingSessions;
        _apiKeyToSession = sessionModule.apiKeyToSession;
    }
    return {
        browserSessions: _browserSessions!,
        pendingSessions: _pendingSessions!,
        apiKeyToSession: _apiKeyToSession!,
    };
}

/**
 * Start a headless session using a scoped key's stored credentials.
 * Reusable from both the HTTP /start-session endpoint and the WS /ws/api handler.
 *
 * @param scopedKeyId The scoped API key row ID
 * @param masterApiKey The parent user's master API key
 * @param worldName Optional world to select
 * @returns The connected clientId
 */
export async function startHeadlessWithStoredCredentials(
    scopedKeyId: number,
    masterApiKey: string,
    worldName?: string,
): Promise<string> {
    const scopedKeyRecord = await ApiKey.findOne({ where: { id: scopedKeyId } });
    if (!scopedKeyRecord) throw new Error('Scoped key not found');

    const get = (f: string) => scopedKeyRecord.getDataValue ? scopedKeyRecord.getDataValue(f) : (scopedKeyRecord as any)[f];
    const encPwd = get('encryptedFoundryPassword');
    const foundryUrl = get('foundryUrl');
    const username = get('foundryUsername');

    if (!encPwd || !foundryUrl || !username) {
        throw new Error('Scoped key does not have stored Foundry credentials');
    }

    if (!isEncryptionAvailable()) {
        throw new Error('Credential decryption not available. CREDENTIALS_ENCRYPTION_KEY not configured.');
    }

    const password = decrypt(encPwd, get('passwordIv'), get('passwordAuthTag'));
    const apiKey = masterApiKey;

    // Check if a client is already connected for this key and URL
    // We only ever want one GM user connected per worldId/URL combo
    const existingClients = await ClientManager.getConnectedClients(apiKey);
    if (existingClients.length > 0) {
        const { apiKeyToSession } = await getSessionMaps();
        const allSessions: any[] = apiKeyToSession.getAll(apiKey);

        // Look for a session matching the requested foundryUrl
        const matchingSession = allSessions.find((s: any) => s.foundryUrl === foundryUrl);
        if (matchingSession) {
            log.info(`GM already connected for ${foundryUrl}, reusing clientId=${matchingSession.clientId}`);
            return matchingSession.clientId;
        }

        // If no URL match but there are sessions, check if any client lacks session metadata
        // (legacy sessions without foundryUrl stored — assume same URL for safety)
        if (allSessions.length === 0 || allSessions.every((s: any) => !s.foundryUrl)) {
            log.info(`Client already connected for key ${apiKey.substring(0, 8)}...: ${existingClients[0]}`);
            return existingClients[0];
        }

        // Sessions exist for different URLs — proceed to create new session
        log.info(`Existing sessions are for different URLs, will create new session for ${foundryUrl}`);
    }

    // Check if a session is already being launched for this URL
    const sessionMaps = await getSessionMaps();
    for (const pending of sessionMaps.pendingSessions.values()) {
        if (pending.apiKey === apiKey && pending.foundryUrl === foundryUrl) {
            log.info(`Headless session already launching for ${foundryUrl}, waiting for it`);
            // Wait for the pending session to complete by polling for a connected client
            const clientId = await new Promise<string>((resolve, reject) => {
                const check = setInterval(async () => {
                    const allSessionsNow: any[] = sessionMaps.apiKeyToSession.getAll(apiKey);
                    const ready = allSessionsNow.find((s: any) => s.foundryUrl === foundryUrl);
                    if (ready) { clearInterval(check); clearTimeout(timeout); resolve(ready.clientId); }
                }, 2000);
                const timeout = setTimeout(() => { clearInterval(check); reject(new Error('Timeout waiting for pending session')); }, 300000);
            });
            return clientId;
        }
    }

    log.info(`Starting headless session via stored credentials for ${username} at ${foundryUrl}`);

    const puppeteer = await import('puppeteer');
    const browser = await puppeteer.launch({
        headless: true,
        executablePath: process.env.PUPPETEER_EXECUTABLE_PATH || undefined,
        args: [
            '--no-sandbox', '--disable-setuid-sandbox', '--enable-gpu-rasterization',
            '--enable-oop-rasterization', '--disable-dev-shm-usage', '--no-first-run',
            '--no-zygote', '--disable-extensions', '--disable-web-security',
            '--disable-features=site-per-process,IsolateOrigins,site-isolation-trials',
            '--disable-background-networking', '--disable-background-timer-throttling',
            '--disable-backgrounding-occluded-windows', '--disable-renderer-backgrounding',
            '--disable-sync', '--disable-breakpad',
            '--disable-component-extensions-with-background-pages',
            '--disable-default-apps', '--disable-infobars', '--disable-popup-blocking',
            '--disable-translate', '--metrics-recording-only', '--mute-audio',
            '--log-level=0', '--js-flags="--max_old_space_size=8192"',
            '--enable-unsafe-swiftshader', '--use-gl=angle', '--use-angle=swiftshader',
        ],
        defaultViewport: { width: 1366, height: 768 }
    });

    const page = await browser.newPage();
    attachBrowserConsoleLogger(page, { username, worldName });
    await page.goto(foundryUrl, { waitUntil: 'networkidle0', timeout: 180000 });

    if (worldName) {
        await page.waitForSelector('li.package.world', { timeout: 10000 }).catch(() => {});
        await page.evaluate((wn: string) => {
            const titles = Array.from(document.querySelectorAll('h3.package-title'));
            for (const title of titles) {
                if (title.textContent?.trim() === wn) {
                    const li = title.closest('li.package.world');
                    const btn = li?.querySelector('a.control.play') as HTMLElement;
                    if (btn) { btn.click(); return; }
                }
            }
        }, worldName);
        await new Promise(r => setTimeout(r, 6000));
    }

    const hasUserSelect = await page.$('select[name="userid"]').then(el => !!el).catch(() => false);
    if (hasUserSelect) {
        const options = await page.$$eval('select[name="userid"] option', opts =>
            opts.map((o: any) => ({ value: o.value, text: o.textContent?.trim() })));
        const match = options.find((o: any) => o.text === username);
        if (match) await page.select('select[name="userid"]', match.value);
    }
    await page.type('input[name="password"]', password);
    await page.click('button[type="submit"]').catch(() =>
        page.evaluate(() => (document.querySelector('form') as HTMLFormElement)?.submit()));

    await page.waitForSelector('#ui-left, #sidebar, .vtt, #game', { timeout: 30000 });

    const sessionId = crypto.randomUUID();
    const { browserSessions, pendingSessions, apiKeyToSession } = sessionMaps;

    browserSessions.set(sessionId, browser);
    await registerHeadlessSession(sessionId, username, apiKey);

    // We can't predict the exact clientId format because it uses the Foundry internal
    // user ID (e.g. "foundry-r6bXhB7k9cXa3cif"), not the login username ("tester").
    // Instead, poll for any NEW client that connects with the matching master API key.
    const clientsBefore = new Set(await ClientManager.getConnectedClients(apiKey));

    pendingSessions.set(sessionId, {
        sessionId, browser, apiKey, expectedClientId: `foundry-*-${username} (polling)`,
        startTime: Date.now(), foundryUrl, worldName, username
    });

    const connectedClientId = await new Promise<string>((resolve, reject) => {
        const check = setInterval(async () => {
            const currentClients = await ClientManager.getConnectedClients(apiKey);
            const newClientIds = currentClients.filter(id => !clientsBefore.has(id));
            for (const candidateId of newClientIds) {
                const client = await ClientManager.getClient(candidateId);
                if (!client) continue;
                const clientWorldId = client.getWorldId();
                // If worldName is known, verify the client is from the right world
                if (worldName && clientWorldId && clientWorldId !== worldName) continue;
                clearInterval(check); clearTimeout(timeout); resolve(candidateId);
                return;
            }
        }, 2000);
        const timeout = setTimeout(() => { clearInterval(check); reject(new Error('Timeout waiting for Foundry client connection')); }, 300000);
    });

    // Update the pending session with the actual clientId
    pendingSessions.delete(sessionId);
    apiKeyToSession.set(apiKey, { sessionId, clientId: connectedClientId, lastActivity: Date.now(), foundryUrl });

    log.info(`Headless session ready: sessionId=${sessionId}, clientId=${connectedClientId}`);
    return connectedClientId;
}
