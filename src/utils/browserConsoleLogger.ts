/**
 * Shared browser console logging for headless Puppeteer sessions.
 * Reads the log level from the CAPTURE_BROWSER_CONSOLE env var (or an explicit override).
 * Valid levels: "error", "warn", "debug" (each includes all levels above it).
 */
import fs from 'fs';
import path from 'path';
import type { Page, ConsoleMessage, HTTPRequest } from 'puppeteer';
import { log } from './logger';

const LEVEL_ORDER = ['error', 'warn', 'debug'] as const;
type BrowserConsoleLevel = typeof LEVEL_ORDER[number];

export interface BrowserConsoleLoggerOptions {
    /** Override for env var — typically from a request body param */
    level?: string;
    worldName?: string;
    username: string;
    foundryVersion?: string;
}

/**
 * Attach console/error/requestfailed listeners to a Puppeteer page that write
 * to a log file under `data/browser-logs/`.
 *
 * Does nothing if neither `options.level` nor `CAPTURE_BROWSER_CONSOLE` is set.
 */
export function attachBrowserConsoleLogger(page: Page, options: BrowserConsoleLoggerOptions): void {
    const browserConsoleLevel = options.level || process.env.CAPTURE_BROWSER_CONSOLE;
    if (!browserConsoleLevel) return;

    const threshold = LEVEL_ORDER.indexOf(browserConsoleLevel as BrowserConsoleLevel);
    if (threshold === -1) {
        log.warn(`Invalid CAPTURE_BROWSER_CONSOLE value: "${browserConsoleLevel}". Use "error", "warn", or "debug".`);
        return;
    }

    const logsDir = path.join(process.cwd(), 'data', 'browser-logs');
    fs.mkdirSync(logsDir, { recursive: true });

    // Remove previous log files for the same world/username
    const logPrefix = `${options.worldName || 'unknown'}_${options.username}`;
    try {
        for (const file of fs.readdirSync(logsDir)) {
            if (file.startsWith(logPrefix) && file.endsWith('.log')) {
                fs.unlinkSync(path.join(logsDir, file));
            }
        }
    } catch { /* ignore cleanup errors */ }

    const timestamp = new Date().toISOString().replace(/[:.]/g, '-');
    const logFileName = `${logPrefix}_${options.foundryVersion ? `fvtt-${options.foundryVersion}_` : ''}${timestamp}.log`;
    const logStream = fs.createWriteStream(path.join(logsDir, logFileName), { flags: 'a' });
    log.info(`Browser console logs will be written to data/browser-logs/${logFileName}`);

    page.on('pageerror', (error: unknown) => {
        const message = error instanceof Error ? error.message : String(error);
        const ts = new Date().toISOString();
        logStream.write(`[${ts}] [PAGEERROR] ${message}\n`);
    });

    page.on('requestfailed', (request: HTTPRequest) => {
        const ts = new Date().toISOString();
        logStream.write(`[${ts}] [REQUESTFAILED] ${request.url()}\n`);
    });

    page.on('console', (msg: ConsoleMessage) => {
        const type = msg.type();
        const msgLevel = type === 'error' ? 0 : type === 'warn' ? 1 : 2;
        if (msgLevel > threshold) return;

        const fallbackText = msg.text();
        const ts = new Date().toISOString();

        const writeLog = (text: string) => {
            logStream.write(`[${ts}] [${type.toUpperCase()}] ${text}\n`);
        };

        const args = msg.args();
        if (args.length === 0) {
            writeLog(fallbackText);
            return;
        }

        // Serialize each arg in the browser context to capture non-enumerable
        // properties, getters, and class instances that JSON.stringify misses
        Promise.all(
            args.map(arg =>
                arg.evaluate((obj: unknown) => {
                    if (obj === null || obj === undefined) return String(obj);
                    if (typeof obj === 'string') return obj;
                    if (obj instanceof Error) return `${obj.name}: ${obj.message}\n${obj.stack}`;
                    try {
                        const seen = new WeakSet();
                        return JSON.stringify(obj, (_key, value) => {
                            if (typeof value === 'object' && value !== null) {
                                if (seen.has(value)) return '[Circular]';
                                seen.add(value);
                            }
                            return value;
                        }, 2);
                    } catch {
                        return String(obj);
                    }
                }).catch(() => arg.toString())
            )
        ).then(resolved => {
            writeLog(resolved.join(' '));
        }).catch(() => {
            writeLog(fallbackText);
        });
    });

    page.once('close', () => logStream.end());
}
