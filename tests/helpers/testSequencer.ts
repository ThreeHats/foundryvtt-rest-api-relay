/**
 * Custom Jest Test Sequencer that enforces a specific test order
 * 
 * Edit TEST_ORDER to add new tests or change execution order.
 */

import Sequencer from '@jest/test-sequencer';
import type { Test } from '@jest/test-result';
import path from 'path';

export const TEST_ORDER = [
  // Phase 0: Account setup (registers ephemeral account or validates pre-provisioned one)
  'register-account.test.ts',

  // Phase 1: Session setup
  'session-endpoints.test.ts',   // Must run first to create sessions

  // Phase 1b: Client listing (requires active session)
  'client-endpoints.test.ts',    // Lists connected clients

  // Phase 2: Entity creation (creates test data for other tests)
  'entity-endpoints.test.ts',    // Creates actors, items, etc.

  // Phase 3: Auth validation
  'auth-requirements.test.ts',
  'auth-endpoints.test.ts',          // Auth endpoint CRUD tests (register, login, user-data, etc.)

  // Phase 4: Scene + Canvas (test scene must be active for canvas and subsequent tests)
  'scene-endpoints.test.ts',     // Creates test scene, activates it
  'canvas-endpoints.test.ts',    // Token/wall/light CRUD + creates persistent token

  // Phase 5: Core functionality tests (runs on the test scene + token)
  'structure-endpoints.test.ts',
  'search-endpoints.test.ts',
  'roll-endpoints.test.ts',
  'sheet-endpoints.test.ts',
  'macro-endpoints.test.ts',
  'utility-endpoints.test.ts',   // Selects token for encounters
  'encounter-endpoints.test.ts',
  'fileSystem-endpoints.test.ts',
  'chat-endpoints.test.ts',
  'permission-filtering.test.ts', // Tests userId permission scoping
  'scoped-keys-endpoints.test.ts', // Scoped API key CRUD + auto-routing (needs clientId)
  'scopes.test.ts',               // Action scopes + multi-client enforcement

  // Phase 5b: New feature tests
  'playlist-endpoints.test.ts',  // Playlist control + play sound
  'hooks-subscribe.test.ts',     // Hooks firehose SSE/WS
  'combat-subscribe.test.ts',    // Combat event subscription
  'websocket-api.test.ts',       // Client-facing WebSocket API (/ws/api)
  'scene-image.test.ts',         // Scene screenshot + raw image
  'user-endpoints.test.ts',      // User CRUD (self-contained: creates + deletes test user)

  // Phase 6: System-agnostic feature tests
  'effects-endpoints.test.ts',   // ActiveEffect CRUD + status effects list

  // Phase 7: System-specific tests (only run on matching systems)
  'dnd5e-endpoints.test.ts',     // Only runs if client has dnd5e (includes concentration + inventory)

  // Phase 5c: Security architecture tests
  'account-security.test.ts',     // Email verification
  'connection-tokens.test.ts',    // Pairing, connection tokens, audit logs
  'credentials.test.ts',          // Credential vault + known clients
  'key-request.test.ts',          // OAuth-like key request flow
  'notifications.test.ts',        // Connection notification settings
  'transfer.test.ts',             // Relay-brokered transfers

  // Phase 5d: Admin dashboard tests (require admin user, separate cookie auth)
  'admin-auth.test.ts',           // Admin login/logout/JWT/CSRF/lockout
  'admin-users.test.ts',          // Admin user management
  'admin-keys.test.ts',           // Admin API key management
  'admin-clients.test.ts',        // Admin connected clients view
  'admin-audit.test.ts',          // Admin audit log queries
  'admin-health.test.ts',         // Admin system health endpoint
  'admin-metrics.test.ts',        // In-app metrics + Prometheus
  'admin-ops.test.ts',            // Operational tools + feature flags
  'admin-alerts.test.ts',         // Alert subscriptions
  'admin-sessions.test.ts',       // Headless session admin
  'admin-subscriptions.test.ts',  // Stripe admin overview

  // Phase 8: Cleanup (order matters: entities first, then restore scene, then end sessions, then auth)
  'cleanup-entities.test.ts',    // Deletes all created entities
  'scene-cleanup.test.ts',       // Restores original scene, deletes test scene
  'end-sessions.test.ts',        // Must run last to cleanup sessions
  'auth-cleanup.test.ts',        // Deletes throwaway auth test user
];

class OrderedTestSequencer extends Sequencer {
  sort(tests: Array<Test>): Array<Test> {
    // Sort tests based on their position in our ordered list
    const sorted = [...tests].sort((a, b) => {
      const aFilename = path.basename(a.path);
      const bFilename = path.basename(b.path);
      
      const aOrder = TEST_ORDER.indexOf(aFilename);
      const bOrder = TEST_ORDER.indexOf(bFilename);
      
      // Unknown files go last
      const aFinalOrder = aOrder === -1 ? Infinity : aOrder;
      const bFinalOrder = bOrder === -1 ? Infinity : bOrder;
      
      return aFinalOrder - bFinalOrder;
    });
    
    return sorted;
  }
}

export default OrderedTestSequencer;
