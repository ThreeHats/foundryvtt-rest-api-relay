import type { Request, Response, NextFunction } from "express";

interface CorsOptions {
  origin?: string | string[] | ((origin: string) => boolean);
  methods?: string[];
  allowedHeaders?: string[];
  exposedHeaders?: string[];
  credentials?: boolean;
  maxAge?: number;
  preflightContinue?: boolean;
}

export const corsMiddleware = (options: CorsOptions = {}) => {
  const defaultOptions: CorsOptions = {
    origin: "*",  // Allow all origins (required for Foundry instances from anywhere)
    methods: ["GET", "HEAD", "PUT", "PATCH", "POST", "DELETE"],
    allowedHeaders: ["Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization", "x-api-key"],
    exposedHeaders: [],
    credentials: false, // Set to false since we use API key auth (x-api-key header), not cookies
    maxAge: 86400, // 24 hours
    preflightContinue: false,
  };

  // Merge provided options with defaults
  const corsOptions = { ...defaultOptions, ...options };

  return async (req: Request, res: Response, next: NextFunction) => {
    // Get the origin from request headers
    const origin = req.headers.origin;
    
    // Handle CORS headers
    // If credentials are needed, echo back the origin (can't use * with credentials)
    // Otherwise use * for maximum compatibility
    if (corsOptions.credentials && origin) {
      res.header("Access-Control-Allow-Origin", origin);
      res.header("Access-Control-Allow-Credentials", "true");
    } else {
      res.header("Access-Control-Allow-Origin", "*");
    }
    
    res.header("Access-Control-Allow-Methods", corsOptions.methods!.join(", "));
    res.header("Access-Control-Allow-Headers", corsOptions.allowedHeaders!.join(", "));

    // Handle preflight requests
    if (req.method === "OPTIONS") {
      res.status(200).end()
      return;
    }

    next();
  };
};
