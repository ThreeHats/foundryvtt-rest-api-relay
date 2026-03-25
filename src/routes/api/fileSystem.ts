import { Router, Request, Response } from 'express';
import express from 'express';
import { requestForwarderMiddleware } from '../../middleware/requestForwarder';
import { authMiddleware, trackApiUsage } from '../../middleware/auth';
import { ClientManager } from '../../core/ClientManager';
import { pendingRequests, safeResponse, resolveClientId, resolveScopedUserId } from '../shared';
import { log } from '../../utils/logger';

export const fileSystemRouter = Router();

const commonMiddleware = [requestForwarderMiddleware, authMiddleware, trackApiUsage];

/**
 * Get file system structure
 * 
 * @route GET /file-system
 * @param {string} clientId - [query] The ID of the Foundry client to connect to
 * @param {string} path - [query,?] The path to retrieve (relative to source)
 * @param {string} source - [query,?] The source directory to use (data, systems, modules, etc.)
 * @param {boolean} recursive - [query,?] Whether to recursively list all subdirectories
 * @param {string} userId - [query,?] Foundry user ID or username to scope permissions (omit for GM-level access)
 * @returns {object} File system structure with files and directories
 */
fileSystemRouter.get("/file-system", ...commonMiddleware, async (req: Request, res: Response) => {
    const path = req.query.path as string || "";
    const source = req.query.source as string || "data";
    const recursive = req.query.recursive === "true";
    const userId = resolveScopedUserId(req, req.query.userId as string | undefined);

    const clientId = await resolveClientId(req, res, req.query.clientId as string);
    if (!clientId) return;

    const client = await ClientManager.getClient(clientId);
    if (!client) {
      safeResponse(res, 404, { 
        error: "Invalid client ID"
      });
      return;
    }
    
    try {
      const requestId = `file_system_${Date.now()}_${Math.random().toString(36).substring(2, 9)}`;
      
      pendingRequests.set(requestId, { 
        res,
        type: 'file-system',
        clientId,
        timestamp: Date.now() 
      });
      
      const sent = client.send({
        type: "file-system",
        path,
        source,
        recursive,
        requestId,
        ...(userId && { userId })
      });

      if (!sent) {
        pendingRequests.delete(requestId);
        safeResponse(res, 500, { error: "Failed to send request to Foundry client" });
        return;
      }
      
      setTimeout(() => {
        if (pendingRequests.has(requestId)) {
          pendingRequests.delete(requestId);
          safeResponse(res, 504, { error: "Request timed out" });
        }
      }, 15000);
    } catch (error) {
      log.error(`Error processing file system request: ${error}`);
      safeResponse(res, 500, { error: "Failed to process file system request" });
      return;
    }
});

/**
 * Upload a file to Foundry's file system (handles both base64 and binary data)
 * 
 * @route POST /upload
 * @param {string} clientId - [query] The ID of the Foundry client to connect to
 * @param {string} path - [query/body] The directory path to upload to
 * @param {string} filename - [query/body] The filename to save as
 * @param {string} source - [query/body,?] The source directory to use (data, systems, modules, etc.)
 * @param {string} mimeType - [query/body,?] The MIME type of the file
 * @param {boolean} overwrite - [query/body,?] Whether to overwrite an existing file
 * @param {string} fileData - [body,?] Base64 encoded file data (if sending as JSON) 250MB limit
 * @param {string} userId - [query,?] Foundry user ID or username to scope permissions (omit for GM-level access)
 * @returns {object} Result of the file upload operation
 */
fileSystemRouter.post("/upload", ...commonMiddleware, async (req: express.Request, res: express.Response) => {
    // Handle different content types
    const contentType = req.get('Content-Type') || '';
    let parsePromise: Promise<void>;
    
    if (contentType.includes('application/json')) {
      // Parse as JSON with size limit
      parsePromise = new Promise((resolve, reject) => {
        express.json({ limit: '250mb' })(req, res, (err) => {
          if (err) reject(err);
          else resolve();
        });
      });
    } else {
      // Parse as raw binary data
      parsePromise = new Promise((resolve, reject) => {
        express.raw({ limit: '250mb', type: '*/*' })(req, res, (err) => {
          if (err) reject(err);
          else resolve();
        });
      });
    }
    
    try {
      await parsePromise;
    } catch (error) {
      safeResponse(res, 400, {
        error: "Failed to parse request body",
        details: error instanceof Error ? error.message : String(error),
        suggestion: "Check your request size (max 250MB) and content type"
      });
      return;
    }

    const path = req.query.path || req.body?.path as string;
    const filename = req.query.filename || req.body?.filename as string;
    const source = req.query.source as string || req.body?.source || "data";
    const mimeType = req.query.mimeType as string || req.body?.mimeType || "application/octet-stream";
    const overwrite = req.query.overwrite === "true" || req.body?.overwrite === "true" || req.body?.overwrite === true;
    const fileData = req.body?.fileData as string | undefined;
    const userId = resolveScopedUserId(req, req.query.userId as string || req.body?.userId as string | undefined);

    const clientId = await resolveClientId(req, res, req.query.clientId as string);
    if (!clientId) return;

    if (!path || !filename) {
      safeResponse(res, 400, {
        error: "Required parameters missing",
        requiredParams: "path, filename",
        howToUse: "Add ?path=your/path&filename=your-file.png to your request"
      });
      return;
    }

    const client = await ClientManager.getClient(clientId);
    if (!client) {
      safeResponse(res, 404, {
        error: "Invalid client ID"
      });
      return;
    }

    try {
      let binaryData: number[] | null = null;
      let processedFileData: string | null = null;

      // Handle different types of file data
      if (fileData) {
        // Handle base64 data from JSON body
        const base64Match = fileData.match(/^data:([A-Za-z-+\/]+);base64,(.+)$/);
        if (!base64Match) {
          safeResponse(res, 400, {
            error: "Invalid file data format",
            expected: "Base64 encoded data URL (e.g., data:image/png;base64,...)",
            received: fileData.substring(0, 50) + "..."
          });
          return;
        }
        
        // Validate base64 data
        try {
          const base64Data = base64Match[2];
          const buffer = Buffer.from(base64Data, 'base64');
          if (buffer.length === 0) {
            throw new Error("Empty file data");
          }
          processedFileData = fileData;
          log.info(`Processing base64 file data: ${buffer.length} bytes`);
        } catch (error) {
          safeResponse(res, 400, {
            error: "Invalid base64 data",
            details: error instanceof Error ? error.message : String(error)
          });
          return;
        }
      } else if (contentType.includes('application/octet-stream') || !contentType.includes('application/json')) {
        // Handle binary data from raw body
        if (Buffer.isBuffer(req.body) && req.body.length > 0) {
          binaryData = Array.from(req.body);
          log.info(`Processing binary file data: ${req.body.length} bytes`);
        } else {
          safeResponse(res, 400, {
            error: "No file data received",
            tip: "Send binary file data with Content-Type: application/octet-stream, or JSON with base64 fileData field",
            contentType: contentType
          });
          return;
        }
      } else {
        safeResponse(res, 400, {
          error: "No file data provided",
          howToProvide: [
            "Option 1: Send JSON with fileData field containing base64 data URL",
            "Option 2: Send binary data with Content-Type: application/octet-stream"
          ]
        });
        return;
      }
      // Generate a unique requestId
      const requestId = `upload_file_${Date.now()}_${Math.random().toString(36).substring(2, 9)}`;

      pendingRequests.set(requestId, {
        res,
        type: 'upload-file',
        clientId,
        timestamp: Date.now()
      });

      const payload: any = {
        type: "upload-file",
        path,
        filename,
        source: source || "data",
        overwrite: overwrite || false,
        requestId,
        ...(userId && { userId })
      };

      if (processedFileData) {
        payload.fileData = processedFileData;
        payload.mimeType = mimeType;
      } else if (binaryData) {
        payload.binaryData = binaryData;
        payload.mimeType = mimeType;
      } else {
        pendingRequests.delete(requestId);
        safeResponse(res, 400, {
          error: "No valid file data to send",
          debug: { hasFileData: !!processedFileData, hasBinaryData: !!binaryData }
        });
        return;
      }

      log.info(`Sending upload request: ${JSON.stringify({ 
        requestId, 
        path, 
        filename, 
        source, 
        hasFileData: !!processedFileData, 
        hasBinaryData: !!binaryData,
        payloadSize: processedFileData ? processedFileData.length : (binaryData ? binaryData.length : 0)
      })}`);

      const sent = client.send(payload);

      if (!sent) {
        pendingRequests.delete(requestId);
        safeResponse(res, 500, { error: "Failed to send request to Foundry client" });
        return;
      }

      // Set timeout for request - file uploads may take longer
      setTimeout(() => {
        if (pendingRequests.has(requestId)) {
          pendingRequests.delete(requestId);
          safeResponse(res, 504, { 
            error: "File upload request timed out",
            suggestion: "Try uploading a smaller file or check your connection to Foundry"
          });
        }
      }, 30000); // 30 second timeout for uploads
    } catch (error) {
      log.error(`Error processing file upload request: ${error}`);
      if (error instanceof Error) {
        log.error(`Upload error stack: ${error.stack}`);
      }
      safeResponse(res, 500, { 
        error: "Failed to process file upload request",
        details: error instanceof Error ? error.message : String(error)
      });
      return;
    }
});

/**
 * Download a file from Foundry's file system
 * 
 * @route GET /download
 * @param {string} clientId - [query] The ID of the Foundry client to connect to
 * @param {string} path - [query] The full path to the file to download
 * @param {string} source - [query,?] The source directory to use (data, systems, modules, etc.)
 * @param {string} format - [query,?] The format to return the file in (binary, base64)
 * @param {string} userId - [query,?] Foundry user ID or username to scope permissions (omit for GM-level access)
 * @returns {binary|object} File contents in the requested format
 */
fileSystemRouter.get("/download", ...commonMiddleware, async (req: Request, res: Response) => {
    const path = req.query.path as string;
    const source = req.query.source as string || "data";
    const format = req.query.format as string || "binary";
    const userId = resolveScopedUserId(req, req.query.userId as string | undefined);

    const clientId = await resolveClientId(req, res, req.query.clientId as string);
    if (!clientId) return;

    if (!path) {
      safeResponse(res, 400, { 
        error: "Path parameter is required",
        howToUse: "Add &path=yourFilePath to your request" 
      });
      return;
    }
    
    const client = await ClientManager.getClient(clientId);
    if (!client) {
      safeResponse(res, 404, { 
        error: "Invalid client ID"
      });
      return;
    }
    
    try {
      // Generate a unique requestId
      const requestId = `download_file_${Date.now()}_${Math.random().toString(36).substring(2, 9)}`;
      
      pendingRequests.set(requestId, { 
        res,
        type: 'download-file',
        clientId,
        format, // Store the requested format in the pending request
        timestamp: Date.now() 
      });
      
      const sent = client.send({
        type: "download-file",
        path,
        source,
        requestId,
        ...(userId && { userId })
      });

      if (!sent) {
        pendingRequests.delete(requestId);
        safeResponse(res, 500, { error: "Failed to send request to Foundry client" });
        return;
      }
      
      // Set timeout for request
      setTimeout(() => {
        if (pendingRequests.has(requestId)) {
          pendingRequests.delete(requestId);
          safeResponse(res, 504, { error: "File download request timed out" });
        }
      }, 45000); // 45 second timeout for downloads
    } catch (error) {
      log.error(`Error processing file download request: ${error}`);
      safeResponse(res, 500, { error: "Failed to process file download request" });
      return;
    }
});
