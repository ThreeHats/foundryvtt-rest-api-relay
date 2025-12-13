/**
 * Custom Jest Test Sequencer that enforces a specific test order
 * 
 * Edit TEST_ORDER to add new tests or change execution order.
 */

import Sequencer from '@jest/test-sequencer';
import type { Test } from '@jest/test-result';

export const TEST_ORDER = [
  // Phase 1: Session setup
  'session-endpoints.test.ts',   // Must run first to create sessions
  
  // Phase 2: Entity creation (creates test data for other tests)
  'entity-endpoints.test.ts',    // Creates actors, items, etc.
  
  // Phase 3: Auth validation
  'auth-requirements.test.ts',
  
  // Phase 4: Core functionality tests
  'structure-endpoints.test.ts',
  'search-endpoints.test.ts',
  'roll-endpoints.test.ts',
  'sheet-endpoints.test.ts',
  'macro-endpoints.test.ts',
  'utility-endpoints.test.ts', // Selects tokens for encounters
  'encounter-endpoints.test.ts',
  'fileSystem-endpoints.test.ts',
  
  // TODO: implement system tests (and maybe some other systems ;)
  // Phase 5: System-specific tests (implement after adding /modules or something bc of dependencies)
  // also needs actual dnd5e items and actors to properly test. Need to think about how to properly implement this
//   'dnd5e-endpoints.test.ts',     // Only runs if client has dnd5e

  // Phase 6: Cleanup
  'cleanup-entities.test.ts',    // Deletes all created entities
  'end-sessions.test.ts',        // Must run last to cleanup sessions
];

class OrderedTestSequencer extends Sequencer {
  sort(tests: Array<Test>): Array<Test> {
    // Sort tests based on their position in our ordered list
    const sorted = [...tests].sort((a, b) => {
      const aFilename = a.path.split('/').pop() || '';
      const bFilename = b.path.split('/').pop() || '';
      
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
