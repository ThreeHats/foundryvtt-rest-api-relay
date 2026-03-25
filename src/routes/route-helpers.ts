import { Request, Response } from 'express';
import { ClientManager } from '../core/ClientManager';
import { pendingRequests, safeResponse, PendingRequest, PendingRequestType, tryAutoStartForScopedKey } from './shared';
import { log } from '../utils/logger';

/**
 * Defines a parameter to be extracted from the request.
 */
interface ParamDef {
  name: string;
  from: 'body' | 'query' | 'params' | ('body' | 'query' | 'params')[];
  type?: 'string' | 'number' | 'boolean' | 'array' | 'object';
}

/**
 * Configuration for creating a standardized API route handler.
 */
interface ApiRouteConfig {
  type: PendingRequestType;
  requiredParams?: ParamDef[];
  optionalParams?: ParamDef[];
  timeout?: number;
  /**
   * Custom validation logic for the parameters.
   * @param params Extracted parameters from the request.
   * @param req The Express request object.
   * @returns An object with an error message and/or a how-to-use message, or null if no validation is needed.
   */
  validateParams?: (params: Record<string, any>, req: Request) => Promise<{ error?: string; howToUse?: string } | null> | { error?: string; howToUse?: string } | null;
  /**
   * Custom logic to build the payload for the client.
   * @param params Extracted parameters from the request.
   * @param req The Express request object.
   * @returns The payload object.
   */
  buildPayload?: (params: Record<string, any>, req: Request) => Promise<Record<string, any>> | Record<string, any>;
  /**
   * Custom logic to build additional properties for the pending request.
   * @param params Extracted parameters from the request.
   * @returns An object with additional properties for PendingRequest.
   */
  buildPendingRequest?: (params: Record<string, any>) => Partial<Omit<PendingRequest, 'res' | 'timestamp' | 'type' | 'clientId'>>;
}

/**
 * Creates a standardized Express route handler for API endpoints.
 * This function abstracts away the boilerplate of handling client lookups,
 * request tracking, and timeouts.
 *
 * @param config - The configuration for the API route.
 * @returns An Express route handler function.
 */
export function createApiRoute(config: ApiRouteConfig) {
  return async (req: Request, res: Response) => {
    // Extract parameters from request body, query or path params
    const params: Record<string, any> = {};
    const allParamDefs = [...(config.requiredParams || []), ...(config.optionalParams || [])];

    for (const p of allParamDefs) {
      const sources = Array.isArray(p.from) ? p.from : [p.from];
      let value: any;
      for (const source of sources) {
        value = req[source]?.[p.name];
        if (value !== undefined) {
          break;
        }
      }
      params[p.name] = value;
    }

    // Type validation and coercion
    for (const p of allParamDefs) {
        let value = params[p.name];
        if (value === undefined || value === null) continue;

        if (p.type) {
          let coercedValue = value;
          let validationError: string | null = null;

          switch (p.type) {
            case 'number':
              if (typeof value !== 'number') {
                coercedValue = parseFloat(value);
              }
              if (isNaN(coercedValue)) {
                validationError = `'${p.name}' must be a valid number.`;
              }
              break;
            case 'boolean':
              if (typeof value !== 'boolean') {
                if (String(value).toLowerCase() === 'true') coercedValue = true;
                else if (String(value).toLowerCase() === 'false') coercedValue = false;
                else validationError = `'${p.name}' must be a valid boolean.`;
              }
              break;
            case 'array':
              if (!Array.isArray(value)) {
                // Try to parse as array if it's a string from query params
                try {
                  if (typeof value === 'string') {
                    coercedValue = JSON.parse(value);
                    if (!Array.isArray(coercedValue)) {
                      validationError = `'${p.name}' must be an array.`;
                    }
                  } else {
                    validationError = `'${p.name}' must be an array.`;
                  }
                } catch (e) {
                  validationError = `'${p.name}' must be a valid array.`;
                }
              }
              break;
            case 'string':
              if (typeof value !== 'string') validationError = `'${p.name}' must be a string.`;
              break;
            case 'object':
              if (typeof value !== 'object' || Array.isArray(value)) validationError = `'${p.name}' must be an object.`;
              break;
          }

          if (validationError) {
            return safeResponse(res, 400, { error: validationError });
          }
          params[p.name] = coercedValue;
        }
      }

    // Validate parameters
    const validationResult = (await config.validateParams?.(params, req)) || null;
    if (validationResult) {
      return safeResponse(res, 400, validationResult);
    }

    // Validate that all required parameters are present
    for (const p of config.requiredParams || []) {
      if (params[p.name] === undefined || params[p.name] === null) {
        return safeResponse(res, 400, { error: `'${p.name}' is required` });
      }
    }

    // Scoped key enforcement — always override, never trust caller
    if (req.scopedKey) {
      if (req.scopedKey.scopedClientId) {
        params.clientId = req.scopedKey.scopedClientId;
      }
      if (req.scopedKey.scopedUserId) {
        params.userId = req.scopedKey.scopedUserId;
      }
    }

    // Auto-resolve clientId if still missing
    if (!params.clientId) {
      const masterKey = req.masterApiKey || req.user?.apiKey ||
        (req.user?.getDataValue ? req.user.getDataValue('apiKey') : undefined);
      if (masterKey) {
        const clients = await ClientManager.getConnectedClients(masterKey);
        if (clients.length === 1) {
          params.clientId = clients[0];
        } else if (clients.length === 0) {
          // No clients connected — try auto-starting a headless session
          const autoClientId = await tryAutoStartForScopedKey(req);
          if (autoClientId) {
            params.clientId = autoClientId;
          } else {
            return safeResponse(res, 404, { error: 'No connected Foundry clients found' });
          }
        } else {
          return safeResponse(res, 400, {
            error: 'Multiple clients connected. Please specify clientId.',
            connectedClients: clients
          });
        }
      } else {
        return safeResponse(res, 400, { error: "'clientId' is required" });
      }
    }

    const clientId = params.clientId as string;

    // Get the client instance
    let client = await ClientManager.getClient(clientId);

    // If scopedClientId was forced but client isn't connected, try auto-start
    if (!client && req.scopedKey) {
      const autoClientId = await tryAutoStartForScopedKey(req);
      if (autoClientId) {
        params.clientId = autoClientId;
        client = await ClientManager.getClient(autoClientId);
      }
    }

    if (!client) {
      return safeResponse(res, 404, { error: "Invalid client ID" });
    }

    try {
      const requestId = `${config.type}_${Date.now()}`;

      // Register the pending request
      const pendingRequestData: PendingRequest = {
        res,
        type: config.type,
        clientId: clientId,
        timestamp: Date.now(),
        ...(config.buildPendingRequest ? config.buildPendingRequest(params) : {}),
      };
      pendingRequests.set(requestId, pendingRequestData);

      // Build the payload, excluding clientId
      const payloadSource = config.buildPayload
        ? await config.buildPayload(params, req)
        : params;
      const { clientId: _clientId, type: userDefinedType, ...payload } = payloadSource;

      // Send the message to the Foundry client
      const sent = client.send({
        type: config.type,
        requestId,
        ...payload,
        data: {
            ...payload.data,
        }
      });

      // If sending fails, clean up and respond with an error
      if (!sent) {
        pendingRequests.delete(requestId);
        return safeResponse(res, 500, { error: "Failed to send request to Foundry client" });
      }

      // Set a timeout for the request
      const timeoutDuration = config.timeout || 10000;
      setTimeout(() => {
        if (pendingRequests.has(requestId)) {
          pendingRequests.delete(requestId);
          safeResponse(res, 408, { error: "Request timed out" });
        }
      }, timeoutDuration);
    } catch (error) {
      log.error(`Error processing ${config.type} request:`, { error });
      safeResponse(res, 500, { error: `Internal server error during ${config.type} request` });
    }
  };
}
