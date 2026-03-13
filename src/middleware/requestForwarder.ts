import { Request, Response, NextFunction } from 'express';
import { log } from '../utils/logger';
import { getRedisClient } from '../config/redis';
import fetch from 'node-fetch';

// Constants
const INSTANCE_ID = process.env.FLY_ALLOC_ID || 'local';
const FLY_INTERNAL_PORT = process.env.FLY_INTERNAL_PORT || '3010';
const APP_NAME = process.env.APP_NAME || 'foundryvtt-rest-api-relay';

export async function requestForwarderMiddleware(req: Request, res: Response, next: NextFunction): Promise<void> {
  try {
    // Skip health checks and static assets
    if (req.path === '/health' || req.path.startsWith('/static') || req.path === '/') {
      return next();
    }
    
    // Get the API key from the header
    const apiKey = req.header('x-api-key');
    if (!apiKey) {
      return next(); // No API key, continue
    }

    const redis = getRedisClient();
    if (!redis) {
      return next(); // No Redis, handle locally
    }

    // Check if we can find which instance serves the client
    let targetInstanceId: string | null = null;
    const clientId = req.query.clientId as string;
    
    if (clientId) {
      try {
        // Check if this client exists on another instance
        const clientInstance = await redis.get(`client:${clientId}:instance`);
        if (clientInstance && clientInstance !== INSTANCE_ID) {
          targetInstanceId = clientInstance;
          log.info(`Client ${clientId} is connected to instance ${targetInstanceId}, not this instance ${INSTANCE_ID}`);
        }
      } catch (error) {
        log.error(`Error checking client instance: ${error}`);
      }
    }

    // If no target found with client ID, check for API key mapping
    if (!targetInstanceId) {
      try {
        const instanceId = await redis.get(`apikey:${apiKey}:instance`);
        if (instanceId && instanceId !== INSTANCE_ID) {
          targetInstanceId = instanceId;
          log.info(`Forwarding request for API key ${apiKey.substring(0, 8)}... to instance ${targetInstanceId}`);
        }
      } catch (error) {
        log.error(`Error checking API key instance: ${error}`);
      }
    }

    // If no target instance is found, process locally
    if (!targetInstanceId) {
      return next();
    }

    // Forward the request to the target instance
    const targetUrl = `http://${targetInstanceId}.vm.${APP_NAME}.internal:${FLY_INTERNAL_PORT}${req.originalUrl}`;
    log.info(`Forwarding to proxy: ${targetUrl}`);

    // Create safe headers object, removing host to avoid conflicts
    const headers: Record<string, string> = {};
    Object.entries(req.headers).forEach(([key, value]) => {
      if (key.toLowerCase() !== 'host' && typeof value === 'string') {
        headers[key] = value;
      } else if (key.toLowerCase() !== 'host' && Array.isArray(value)) {
        headers[key] = value[0] || '';
      }
    });
    
    // Set up timeout with AbortController
    const controller = new AbortController();
    // Use longer timeout for download/upload requests
    const isFileOperation = req.path === '/upload' || req.path === '/download';
    const timeout = isFileOperation ? 45000 : 20000; // 45s for file operations, 20s for others
    const timeoutId = setTimeout(() => controller.abort(), timeout);
      // Prepare the body based on content type
    let requestBody: any = undefined;
    
    if (req.method !== 'GET' && req.method !== 'HEAD') {
      const contentType = req.headers['content-type'] || '';
        // Handle binary data for uploads and potentially binary responses for downloads
      if (contentType.includes('application/octet-stream') || 
          (req.path === '/upload' && !contentType.includes('application/json')) ||
          req.path === '/download') {
        // Pass binary data as-is without JSON stringify
        requestBody = req.body;
        
        // Enhanced logging for binary data handling
        if (req.body) {
          const bodySize = Buffer.isBuffer(req.body) ? req.body.length : 
                          (typeof req.body === 'string' ? Buffer.byteLength(req.body) : 'unknown type');
          log.info(`Forwarding binary data request to ${targetInstanceId}, path: ${req.path}, size: ${bodySize}`);
        } else {
          log.info(`Forwarding request to ${targetInstanceId}, path: ${req.path}, no body data`);
        }
      } else {
        // For JSON and other content types, use JSON stringify
        requestBody = JSON.stringify(req.body);
        log.info(`Forwarding JSON request to ${targetInstanceId}, path: ${req.path}`);
      }
    }
    
    // Forward the request
    const response = await fetch(targetUrl, {
      method: req.method,
      headers,
      body: requestBody,
      signal: controller.signal
    });
    
    clearTimeout(timeoutId);
      // Copy response headers but filter out problematic ones
    Object.entries(response.headers.raw()).forEach(([key, values]) => {
      if (Array.isArray(values) && !['connection', 'content-length', 'transfer-encoding'].includes(key.toLowerCase())) {
        res.setHeader(key, values);
        
        // Log important headers for debugging, especially for binary responses
        if (['content-type', 'content-disposition', 'content-encoding'].includes(key.toLowerCase())) {
          log.info(`Forwarding header: ${key} = ${values.join(', ')}`);
        }
      }
    });// Special handling for different response types based on endpoint and content type
    if ((req.path === '/upload' && req.method === 'POST') || 
        (req.path === '/download' && req.method === 'GET')) {
      
      const endpoint = req.path === '/upload' ? 'upload' : 'download';
      log.info(`Handling forwarded ${endpoint} response from ${targetInstanceId}`);
      
      // Get the content type from the response
      const responseContentType = response.headers.get('content-type') || '';
      const disposition = response.headers.get('content-disposition') || '';
      
      // If it's a download or binary upload response
      if (!responseContentType.includes('application/json') || 
          disposition.includes('attachment') ||
          endpoint === 'download') {
        try {
          log.info(`Processing binary ${endpoint} response with content-type: ${responseContentType}`);
          
          // For binary responses, get the buffer and send it
          const buffer = await response.buffer();
          
          if (!buffer || buffer.length === 0) {
            log.error(`Empty buffer received for ${endpoint} response`);
            res.status(500).json({
              error: `Empty data received from ${endpoint} request`,
              details: "The target instance returned an empty response"
            });
            return;
          }
          
          log.info(`Forwarding binary ${endpoint} response, size: ${buffer.length} bytes`);
          res.status(response.status).send(buffer);
        } catch (error) {
          log.error(`Error processing binary ${endpoint} response: ${error}`);
          res.status(500).json({ 
            error: `Failed to process ${endpoint} response`,
            details: error instanceof Error ? error.message : String(error)
          });
        }
      } else {
        // JSON responses can be handled normally
        const text = await response.text();
        res.status(response.status).send(text);
      }
    } else {
      // Standard response handling for other routes
      const text = await response.text();
      res.status(response.status).send(text);
    }
    
  } catch (error) {
    log.error(`Error in request forwarder: ${error}`);
    
    // Fall back to local handling instead of returning an error
    // This allows the API to still work even if forwarding fails
    next();
  }
}