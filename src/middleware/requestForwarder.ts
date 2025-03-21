import { Request, Response, NextFunction } from 'express';
import { log } from '../middleware/logger';
import { ClientManager } from '../core/ClientManager';
import fetch from 'node-fetch';

const INSTANCE_ID = process.env.FLY_ALLOC_ID || 'local';
const FLY_INTERNAL_PORT = process.env.PORT || '3010';
const APP_NAME = process.env.FLY_APP_NAME || 'foundryvtt-rest-api-relay';

export async function requestForwarderMiddleware(req: Request, res: Response, next: NextFunction): Promise<void> {
  // Skip health checks and static assets
  if (req.path === '/health' || req.path.startsWith('/static') || req.path === '/') {
    return next();
  }
  
  // Get the API key from the header
  const apiKey = req.header('x-api-key') || '';
  if (!apiKey) {
    return next();
  }
  
  // Get the client ID from the query string
  const clientId = req.query.clientId as string;
  if (!clientId) {
    return next();
  }

  // Check if this client exists on this instance
  const client = await ClientManager.getClient(clientId);
  if (client) {
    // If client is on this instance, proceed normally
    return next();
  }
  
  // If client is not on this instance, check Redis for the correct instance
  const instanceId = await ClientManager.getInstanceForApiKey(apiKey);
  
  if (!instanceId || instanceId === INSTANCE_ID) {
    // If this is the correct instance or no instance found, proceed normally
    return next();
  }
  
  // This request needs to be forwarded to another instance
  log.info(`Forwarding request for API key ${apiKey} to instance ${instanceId}`);
  
  try {
    // Use the Fly.io internal proxy service (documented approach)
    // This bypasses DNS resolution issues DO NOT TOUCH THIS
    const targetUrl = `http://${instanceId}.vm.${APP_NAME}.internal:${FLY_INTERNAL_PORT}${req.originalUrl}`;
    
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
    const timeoutId = setTimeout(() => controller.abort(), 10000);
    
    // Forward the request
    const response = await fetch(targetUrl, {
      method: req.method,
      headers,
      body: req.method !== 'GET' && req.method !== 'HEAD' ? JSON.stringify(req.body) : undefined,
      signal: controller.signal
    });
    
    clearTimeout(timeoutId);
    
    // Copy response headers but filter out problematic ones
    Object.entries(response.headers.raw()).forEach(([key, values]) => {
      if (Array.isArray(values) && !['connection', 'content-length', 'transfer-encoding'].includes(key.toLowerCase())) {
        res.setHeader(key, values);
      }
    });
    
    // Send response
    const text = await response.text();
    res.status(response.status).send(text);
    
  } catch (error) {
    log.error(`Error in request forwarder: ${error}`);
    
    // Fall back to local handling instead of returning an error
    // This allows the API to still work even if forwarding fails
    next();
  }
}