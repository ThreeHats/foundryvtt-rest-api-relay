import { ClientManager } from '../../core/ClientManager';
import { Router } from 'express';
import express from 'express';
import { requestForwarderMiddleware } from '../../middleware/requestForwarder';
import { authMiddleware, trackApiUsage } from '../../middleware/auth';
import { pendingRequests, safeResponse } from '../shared';
import { log } from '../../utils/logger';

export const sheetRouter = Router();

const commonMiddleware = [requestForwarderMiddleware, authMiddleware, trackApiUsage, express.json()];

/**
 * Get actor sheet HTML
 * 
 * This endpoint retrieves the HTML for an actor sheet based on the provided UUID or selected actor. Only works on Foundry version 12.
 * 
 * @route GET /sheet
 * @param {string} uuid - [query,?] The UUID of the entity to get the sheet for
 * @param {boolean} selected - [query,?] Whether to get the sheet for the selected entity
 * @param {boolean} actor - [query,?] Whether to get the sheet for the selected token's actor if selected is true
 * @param {string} clientId - [query] The ID of the Foundry client to connect to
 * @param {string} format - [query,?] The format to return the sheet in (html, json)
 * @param {number} scale - [query,?] The initial scale of the sheet
 * @param {number} tab - [query,?] The active tab index to open
 * @param {boolean} darkMode - [query,?] Whether to use dark mode for the sheet
 * @returns {object} The sheet HTML or data depending on format requested
 */
sheetRouter.get("/sheet", ...commonMiddleware, async (req: express.Request, res: express.Response) => {
    const uuid = req.query.uuid as string;
    const selected = req.query.selected === 'true';
    const actor = req.query.actor === 'true';
    const clientId = req.query.clientId as string;
    const format = req.query.format as string || 'html';
    const initialScale = parseFloat(req.query.scale as string) || null;
    const activeTab = req.query.tab ? (isNaN(Number(req.query.tab)) ? null : Number(req.query.tab)) : null;
    const darkMode = req.query.darkMode === 'true';
    
    if (!clientId) {
      safeResponse(res, 400, { 
        error: "Client ID is required to identify the Foundry instance"
      });
			return;
    }

    if (!uuid && !selected) {
      safeResponse(res, 400, { error: "UUID or selected parameter is required" });
      return;
    }
    
    // Find a connected client with this ID
    const client = await ClientManager.getClient(clientId);
    if (!client) {
      safeResponse(res, 404, { 
        error: "Invalid client ID"
      });
			return;
    }
    
    try {
      // Generate a unique requestId for this request
      const requestId = `actor_sheet_${Date.now()}_${Math.random().toString(36).substring(2, 9)}`;
      
      // Store the response object in a request map
      pendingRequests.set(requestId, { 
        res,
        type: 'get-sheet',
        uuid,
        clientId,
        format,
        initialScale,
        activeTab,
        darkMode,
        timestamp: Date.now() 
      });
      
      // Send request to Foundry for actor sheet HTML
      const sent = client.send({
        type: "get-sheet",
        uuid,
        selected,
        actor,
        requestId,
        initialScale,
        activeTab,
        darkMode
      });

      if (!sent) {
        pendingRequests.delete(requestId);
        safeResponse(res, 500, { 
          error: "Failed to send actor sheet request to Foundry client",
          suggestion: "The client may be disconnecting or experiencing issues"  
        });
        return;
      }
      
      // Set timeout for request
      setTimeout(() => {
        if (pendingRequests.has(requestId)) {
          pendingRequests.delete(requestId);
          safeResponse(res, 408, { 
            error: "Actor sheet request timed out", 
            tip: "The Foundry client might be busy or the actor UUID might not exist."
          });
        }
      }, 20000); // 20 seconds timeout

    } catch (error) {
      log.error(`Error processing actor sheet request: ${error}`);
      safeResponse(res, 500, { error: "Failed to process actor sheet request" });
      return;
    }
});