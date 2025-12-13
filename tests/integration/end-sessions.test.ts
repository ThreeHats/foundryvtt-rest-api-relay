/**
 * @file end-sessions.test.ts
 * @generated false
 * @description Session Termination Tests
 * @endpoints DELETE /end-session
 */

import { describe, test, expect, afterAll } from '@jest/globals';
import { ApiRequestConfig } from '../helpers/apiRequest';
import { testVariables } from '../helpers/testVariables';
import { captureExample, appendExamples } from '../helpers/captureExample';
import { forEachVersion } from '../helpers/multiVersion';
import { getGlobalVariable } from '../helpers/globalVariables';
import * as path from 'path';

// Skip end-session tests when using existing sessions (don't kill user's session!)
const useExistingSession = process.env.USE_EXISTING_SESSION === 'true';

// Store captured examples for documentation
const capturedExamples: any[] = [];

const describeOrSkip = useExistingSession ? describe.skip : describe;

describeOrSkip('End Session', () => {
    afterAll(() => {
        // Append to session-examples.json (same file as session-endpoints.test.ts)
        const outputPath = path.join(__dirname, '../../docs/examples/session-examples.json');
        appendExamples(capturedExamples, outputPath);
        console.log(`\nAppended ${capturedExamples.length} examples to ${outputPath}`);
    });

    forEachVersion((version, getClientId) => {
        describe(`/end-session (v${version})`, () => {
            test('DELETE /end-session', async () => {
                // Get sessionId from global variables (set by session-endpoints.test.ts)
                const sessionId = getGlobalVariable(version, 'sessionId');
                
                if (!sessionId) {
                    console.log(`⚠️  No sessionId found for v${version} - skipping end-session`);
                    console.log(`   (session-endpoints.test.ts should have set this via setGlobalVariable)`);
                    return;
                }
                
                console.log(`🔚 Ending session for v${version}: ${sessionId}`);
                
                const requestConfig: ApiRequestConfig = {
                    url: {
                        raw: '{{baseUrl}}/end-session',
                        host: ['{{baseUrl}}'],
                        path: ['end-session'],
                        query: [
                            {
                                key: 'sessionId',
                                value: sessionId,
                                description: 'The ID of the session to end'
                            }
                        ]
                    },
                    method: 'DELETE',
                    header: [
                        {
                            key: 'x-api-key',
                            value: '{{apiKey}}',
                            type: 'text'
                        }
                    ]
                };

                // Capture this example for documentation (also makes the request)
                const captured = await captureExample(
                    requestConfig,
                    testVariables,
                    '/end-session'
                );
                capturedExamples.push(captured);

                // Session should be ended successfully or already gone
                expect([200, 404]).toContain(captured.response.status);
                
                if (captured.response.status === 200) {
                    console.log(`✓ Ended headless session for v${version}: ${sessionId}`);
                } else {
                    console.log(`ℹ️  Session for v${version} already ended or not found`);
                }
            });
        });
    });
});