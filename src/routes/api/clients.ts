import { Router, Request, Response } from 'express';
import { authMiddleware } from '../../middleware/auth';
import { getRedisClient } from '../../config/redis';
import { ClientManager } from '../../core/ClientManager';
import { safeResponse } from '../shared';
import { log } from '../../utils/logger';

const INSTANCE_ID = process.env.INSTANCE_ID || 'default';

export const clientsRouter = Router();

/**
 * Get all connected clients for the authenticated API key
 *
 * Returns a list of all currently connected Foundry VTT clients associated with
 * the provided API key, including their connection details and world information.
 *
 * @route GET /clients
 * @group Clients
 * @security apiKey
 * @returns {object} Object containing total count and array of connected client details
 */
clientsRouter.get("/clients", authMiddleware, async (req: Request, res: Response) => {
  try {
    // Use master key for scoped keys so clients are found via the parent's tokenGroups
    const apiKey = req.masterApiKey || req.header('x-api-key') || '';
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

    // Scoped keys with scopedClientId only see their bound client
    if (req.scopedKey?.scopedClientId) {
      allClients = allClients.filter(c => c.id === req.scopedKey!.scopedClientId);
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
