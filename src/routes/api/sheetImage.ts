import { ClientManager } from '../../core/ClientManager';
import { Router } from 'express';
import express from 'express';
import { requestForwarderMiddleware } from '../../middleware/requestForwarder';
import { authMiddleware, trackApiUsage } from '../../middleware/auth';
import { pendingRequests, safeResponse, resolveClientId, resolveScopedUserId } from '../shared';
import { log } from '../../utils/logger';

export const sheetImageRouter = Router();

const commonMiddleware = [requestForwarderMiddleware, authMiddleware, trackApiUsage, express.json()];

/**
 * Get actor sheet as screenshot image
 *
 * Captures the rendered actor sheet using html2canvas and returns it as a PNG or JPEG image.
 * Works on both Foundry v12 and v13+.
 *
 * @route GET /sheet/image
 * @param {string} uuid - [query,?] The UUID of the entity to screenshot
 * @param {boolean} selected - [query,?] Whether to screenshot the selected entity's sheet
 * @param {boolean} actor - [query,?] Whether to use the selected token's actor if selected is true
 * @param {string} clientId - [query] The ID of the Foundry client to connect to
 * @param {string} format - [query,?] Image format: png or jpeg (default: png)
 * @param {number} quality - [query,?] Image quality 0-1 for JPEG (default: 0.9)
 * @param {number} scale - [query,?] Capture scale factor (default: 1)
 * @param {string} userId - [query,?] Foundry user ID or username to scope permissions
 * @returns {binary} The sheet screenshot as an image
 */
sheetImageRouter.get("/sheet/image", ...commonMiddleware, async (req: express.Request, res: express.Response) => {
    const uuid = req.query.uuid as string;
    const selected = req.query.selected === 'true';
    const actor = req.query.actor === 'true';
    const format = (req.query.format as string) || 'png';
    const quality = parseFloat(req.query.quality as string) || 0.9;
    const scale = parseFloat(req.query.scale as string) || 1;
    const userId = resolveScopedUserId(req, req.query.userId as string | undefined);

    const clientId = await resolveClientId(req, res, req.query.clientId as string);
    if (!clientId) return;

    if (!uuid && !selected) {
        safeResponse(res, 400, { error: "UUID or selected parameter is required" });
        return;
    }

    if (format !== 'png' && format !== 'jpeg') {
        safeResponse(res, 400, { error: "Format must be 'png' or 'jpeg'" });
        return;
    }

    const client = await ClientManager.getClient(clientId);
    if (!client) {
        safeResponse(res, 404, { error: "Invalid client ID" });
        return;
    }

    try {
        const requestId = `sheet_screenshot_${Date.now()}_${Math.random().toString(36).substring(2, 9)}`;

        pendingRequests.set(requestId, {
            res,
            type: 'sheet-screenshot',
            uuid,
            clientId,
            format,
            timestamp: Date.now()
        });

        const sent = client.send({
            type: "sheet-screenshot",
            uuid,
            selected,
            actor,
            requestId,
            format,
            quality,
            scale,
            ...(userId && { userId })
        });

        if (!sent) {
            pendingRequests.delete(requestId);
            safeResponse(res, 500, { error: "Failed to send screenshot request to Foundry client" });
            return;
        }

        setTimeout(() => {
            if (pendingRequests.has(requestId)) {
                pendingRequests.delete(requestId);
                safeResponse(res, 408, {
                    error: "Sheet screenshot request timed out",
                    tip: "The Foundry client might be busy or the actor UUID might not exist."
                });
            }
        }, 20000);

    } catch (error) {
        log.error(`Error processing sheet screenshot request: ${error}`);
        safeResponse(res, 500, { error: "Failed to process sheet screenshot request" });
    }
});
