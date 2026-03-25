import WebSocket from 'ws';
import { log } from '../utils/logger';

export interface SheetSession {
    sessionId: string;
    clientId: string;         // Foundry client ID
    apiKey: string;
    consumerWs: WebSocket;    // The consumer's /ws/api connection
    state: 'pending' | 'active' | 'closed';
    createdAt: number;
    lastActivity: number;
    metadata: { uuid?: string; quality?: number; scale?: number; };
}

const MAX_SESSIONS_PER_API_KEY = parseInt(process.env.MAX_SHEET_SESSIONS_PER_KEY || '3', 10);
const PENDING_TIMEOUT_MS = 30000;
const INACTIVE_TIMEOUT_MS = 5 * 60 * 1000; // 5 minutes
const CLEANUP_INTERVAL_MS = 15000;

const sessions = new Map<string, SheetSession>();

// Track sessions per API key for limit enforcement
function countSessionsForKey(apiKey: string): number {
    let count = 0;
    for (const session of sessions.values()) {
        if (session.apiKey === apiKey && session.state !== 'closed') {
            count++;
        }
    }
    return count;
}

export class SheetSessionManager {
    private static cleanupInterval: ReturnType<typeof setInterval> | null = null;

    static init(): void {
        if (this.cleanupInterval) return;

        this.cleanupInterval = setInterval(() => {
            this.cleanup();
        }, CLEANUP_INTERVAL_MS);
    }

    static createSession(
        clientId: string,
        apiKey: string,
        consumerWs: WebSocket,
        metadata: SheetSession['metadata'] = {}
    ): SheetSession | { error: string } {
        // Check per-key limit
        if (countSessionsForKey(apiKey) >= MAX_SESSIONS_PER_API_KEY) {
            return { error: `Maximum concurrent sheet sessions (${MAX_SESSIONS_PER_API_KEY}) reached for this API key` };
        }

        const sessionId = `ss_${Date.now()}_${Math.random().toString(36).slice(2, 8)}`;
        const now = Date.now();

        const session: SheetSession = {
            sessionId,
            clientId,
            apiKey,
            consumerWs,
            state: 'pending',
            createdAt: now,
            lastActivity: now,
            metadata,
        };

        sessions.set(sessionId, session);
        log.info(`Sheet session ${sessionId} created for client ${clientId}`);
        return session;
    }

    static getSession(sessionId: string): SheetSession | undefined {
        return sessions.get(sessionId);
    }

    static activateSession(sessionId: string): void {
        const session = sessions.get(sessionId);
        if (session && session.state === 'pending') {
            session.state = 'active';
            session.lastActivity = Date.now();
        }
    }

    static updateActivity(sessionId: string): void {
        const session = sessions.get(sessionId);
        if (session) {
            session.lastActivity = Date.now();
        }
    }

    static endSession(sessionId: string): boolean {
        const session = sessions.get(sessionId);
        if (!session) return false;

        session.state = 'closed';
        sessions.delete(sessionId);
        log.info(`Sheet session ${sessionId} ended`);
        return true;
    }

    /**
     * Terminate all sessions for a specific Foundry client (e.g., when it disconnects).
     */
    static terminateSessionsForClient(clientId: string): void {
        for (const [sessionId, session] of sessions.entries()) {
            if (session.clientId === clientId) {
                // Notify consumer if connected
                try {
                    if (session.consumerWs.readyState === WebSocket.OPEN) {
                        session.consumerWs.send(JSON.stringify({
                            type: 'sheet-session-ended',
                            sessionId,
                            reason: 'foundry-disconnected',
                        }));
                    }
                } catch {
                    // Ignore send errors
                }
                sessions.delete(sessionId);
                log.info(`Sheet session ${sessionId} terminated (client ${clientId} disconnected)`);
            }
        }
    }

    /**
     * Terminate all sessions for a specific consumer WebSocket (e.g., when consumer disconnects).
     * Returns session IDs that need to be forwarded to Foundry for cleanup.
     */
    static terminateSessionsForConsumer(consumerWs: WebSocket): string[] {
        const sessionIds: string[] = [];
        for (const [sessionId, session] of sessions.entries()) {
            if (session.consumerWs === consumerWs) {
                sessionIds.push(sessionId);
                sessions.delete(sessionId);
                log.info(`Sheet session ${sessionId} terminated (consumer disconnected)`);
            }
        }
        return sessionIds;
    }

    private static cleanup(): void {
        const now = Date.now();
        for (const [sessionId, session] of sessions.entries()) {
            // Remove stale pending sessions
            if (session.state === 'pending' && now - session.createdAt > PENDING_TIMEOUT_MS) {
                try {
                    if (session.consumerWs.readyState === WebSocket.OPEN) {
                        session.consumerWs.send(JSON.stringify({
                            type: 'sheet-session-error',
                            sessionId,
                            error: 'Session timed out waiting for Foundry response',
                        }));
                    }
                } catch { /* ignore */ }
                sessions.delete(sessionId);
                log.info(`Sheet session ${sessionId} timed out (pending)`);
                continue;
            }

            // Remove inactive sessions
            if (session.state === 'active' && now - session.lastActivity > INACTIVE_TIMEOUT_MS) {
                try {
                    if (session.consumerWs.readyState === WebSocket.OPEN) {
                        session.consumerWs.send(JSON.stringify({
                            type: 'sheet-session-ended',
                            sessionId,
                            reason: 'inactivity-timeout',
                        }));
                    }
                } catch { /* ignore */ }
                sessions.delete(sessionId);
                log.info(`Sheet session ${sessionId} expired (inactive)`);
            }
        }
    }
}
