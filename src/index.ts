/**
 * Main entry point for the FoundryVTT REST API Relay Server.
 * 
 * This server provides WebSocket connectivity and a REST API to access Foundry VTT data remotely.
 * It facilitates communication between Foundry VTT clients and external applications through
 * WebSocket relays and HTTP endpoints.
 * 
 * @author ThreeHats
 * @since 1.8.1
 */

import express, { Request, Response, NextFunction } from "express";
import { createServer } from "http";
import { WebSocketServer } from "ws";
import { corsMiddleware } from "./middleware/cors";
import { log } from "./utils/logger";
import { wsRoutes } from "./routes/websocket";
import { apiRoutes, browserSessions, pendingSessions } from "./routes/api";
import authRoutes from "./routes/auth";
import { config } from "dotenv";
import * as path from "path";
import * as fs from "fs";
import { sequelize } from "./sequelize";
import stripeRouter from './routes/stripe';
import webhookRouter from './routes/webhook';
import { initRedis, closeRedis } from './config/redis';
import { scheduleHeadlessSessionsCheck } from './workers/headlessSessions';
import { redisSessionMiddleware } from './middleware/redisSession';
import { startHealthMonitoring, logSystemHealth, getSystemHealth } from './utils/healthCheck';
import { setupCronJobs } from './cron';
import { migrateDailyRequestTracking } from './migrations/addDailyRequestTracking';

config();

/**
 * Express application instance
 * @public
 */
const app = express();

/**
 * HTTP server instance that wraps the Express app
 * @public
 */
const httpServer = createServer(app);
// Disable timeouts to keep WebSocket connections open may want to sent a long timeout in the future instead
httpServer.setTimeout(0);
httpServer.keepAliveTimeout = 0;
httpServer.headersTimeout = 0;

// Setup CORS
app.use(corsMiddleware());

app.use('/api/webhooks/stripe', express.raw({ type: 'application/json' }));

// Special handling for /upload endpoint to preserve raw body for binary uploads
app.use('/upload', (req, res, next) => {
  const contentType = req.headers['content-type'] || '';
  
  if (!contentType.includes('application/json')) {
    express.raw({ 
      type: '*/*', 
      limit: '250mb' 
    })(req, res, next);
  } else {
    // For JSON requests to /upload, use the regular JSON parser
    express.json({ 
      limit: '250mb' 
    })(req, res, next);
  }
});

// Parse JSON bodies for all other routes with 250MB limit
// Only apply to methods that can have a body (not GET, HEAD, DELETE)
app.use((req, res, next) => {
  if (['POST', 'PUT', 'PATCH'].includes(req.method)) {
    express.json({ limit: '250mb' })(req, res, next);
  } else {
    next();
  }
});

// Add Redis session middleware
app.use(redisSessionMiddleware);

// Add global error handler
app.use((err: any, req: Request, res: Response, next: NextFunction) => {
  log.error('Unhandled error:', err);
  log.error(`Error on ${req.method} ${req.path}`);
  if (!res.headersSent) {
    // Handle JSON parsing errors
    if (err.type === 'entity.parse.failed' || err instanceof SyntaxError || 
        (err.message && err.message.includes('JSON'))) {
      return res.status(400).json({ 
        error: 'Invalid JSON format',
        message: 'The request body contains malformed JSON. Please check your JSON syntax.',
        details: err.message
      });
    }
    
    // Handle payload too large errors
    if (err.type === 'entity.too.large') {
      return res.status(413).json({ 
        error: 'Request entity too large',
        message: 'The request body is too large. Please reduce the size of your request.'
      });
    }
    
    // Default error response
    res.status(500).json({ error: 'Internal server error' });
  }
});

// Serve static files from public directory
app.use("/static", express.static(path.join(__dirname, "../public")));
app.use("/static/css", express.static(path.join(__dirname, "../public/css")));
app.use("/static/js", express.static(path.join(__dirname, "../public/js")));

// Redirect trailing slashes in docs routes to clean URLs
app.use('/docs', (req, res, next) => {
  if (req.path !== '/' && req.path.endsWith('/')) {
    const cleanPath = req.path.slice(0, -1);
    const queryString = req.url.includes('?') ? req.url.substring(req.url.indexOf('?')) : '';
    return res.redirect(301, `/docs${cleanPath}${queryString}`);
  }
  next();
});

// Serve Docusaurus documentation from /docs route
const docsPath = path.resolve(__dirname, "../docs/build");
try {
  // Check if docs build directory exists
  if (fs.existsSync(docsPath)) {
    // Docusaurus builds static HTML files for each page, so we use redirect: false
    // to prevent Express from auto-redirecting directories (which can cause loops)
    app.use("/docs", express.static(docsPath, { 
      index: 'index.html',
      redirect: false
    }));

    // Handle /docs root path
    app.get('/docs', (req, res) => {
      res.sendFile(path.join(docsPath, 'index.html'));
    });

    // Handle directory requests by serving their index.html
    app.get('/docs/*', (req, res, next) => {
      const requestPath = req.path.replace('/docs', '');
      const filePath = path.join(docsPath, requestPath);
      const indexPath = path.join(filePath, 'index.html');
      const notFoundPath = path.join(docsPath, '404.html');
      
      // Try serving index.html first (for directory requests)
      res.sendFile(indexPath, (err) => {
        if (err) {
          // Try serving the file directly (for assets)
          res.sendFile(filePath, (err2) => {
            if (err2) {
              // Serve 404 page
              res.status(404).sendFile(notFoundPath, (err3) => {
                if (err3) next();
              });
            }
          });
        }
      });
    });
  } else {
    log.warn('Documentation build directory not found, docs will not be available');
    app.get('/docs*', (req, res) => {
      res.status(404).json({ error: 'Documentation not available' });
    });
  }
} catch (error) {
  log.error('Error setting up documentation routes:', { error: error instanceof Error ? error.message : String(error) });
  app.get('/docs*', (req, res) => {
    res.status(500).json({ error: 'Documentation setup failed' });
  });
}

// Serve the main HTML page at the root URL
app.get("/", (req: Request, res: Response) => {
  res.sendFile(path.join(__dirname, "../public/index.html"));
});

// Create WebSocket server
const wss = new WebSocketServer({ server: httpServer });

// Setup WebSocket routes
wsRoutes(wss);

// Setup API routes
apiRoutes(app);

// Setup Auth routes
app.use("/", authRoutes);
app.use('/api/subscriptions', stripeRouter);
app.use('/api/webhooks', webhookRouter);

// Add default static image for tokens
app.get("/default-token.png", (req: Request, res: Response) => {
  res.sendFile(path.join(__dirname, "../public/default-token.png"));
});

// Add health endpoint
app.get('/api/health', (req, res) => {
  try {
    const health = getSystemHealth();
    res.status(200).json(health);
  } catch (error) {
    // Always return 200 during startup
    log.warn('Health check error during startup:', { error: error instanceof Error ? error.message : String(error) });
    res.status(200).json({ 
      healthy: true,
      status: 'starting',
      timestamp: Date.now(),
      instanceId: process.env.FLY_ALLOC_ID || 'local',
      message: 'Service initializing'
    });
  }
});

/**
 * Server port number, defaults to 3010 if not specified in environment
 */
const port = process.env.PORT ? parseInt(process.env.PORT) : 3010;

/**
 * Initializes all server services in the correct order.
 * 
 * This function performs the following initialization steps:
 * 1. Starts the HTTP and WebSocket servers first
 * 2. Synchronizes the database connection in background
 * 3. Initializes Redis if configured in background
 * 4. Sets up cron jobs for scheduled tasks in background
 * 5. Starts health monitoring in background
 * 
 * @throws {Error} Exits the process if server startup fails
 * @returns {Promise<void>} Resolves when server is started
 */
async function initializeServices() {
  try {
    httpServer.listen(port, () => {
      log.info(`Server running at http://localhost:${port}`);
      log.info(`WebSocket server ready at ws://localhost:${port}/relay`);
    });
    
    // Do heavy initialization in background after server is running
    setImmediate(async () => {
      try {
        log.info('Starting background initialization...');
        
        // First initialize database
        await sequelize.sync();
        log.info('Database synced');
        
        // Run migration to add daily request tracking columns
        await migrateDailyRequestTracking();
        log.info('Database migrations completed');
        
        if (process.env.REDIS_URL && process.env.REDIS_URL.length > 0) {
          // Then initialize Redis
          const redisInitialized = await initRedis();
          if (!redisInitialized) {
            log.warn('Redis initialization failed - continuing with local storage only');
          } else {
            log.info('Redis initialized successfully');
          }
        }
        
        // Set up cron jobs
        setupCronJobs();
        log.info('Cron jobs initialized');
        
        // Start health monitoring
        logSystemHealth(); // Log initial health
        startHealthMonitoring(60000); // Check every minute
        log.info('Health monitoring started');
        
        log.info('All background services initialized successfully');
      } catch (error) {
        log.error(`Error during background initialization: ${error}`);
        // Don't exit in production - let the server continue running
        if (process.env.NODE_ENV !== 'production') {
          process.exit(1);
        }
      }
    });
    
  } catch (error) {
    log.error(`Error starting server: ${error}`);
    process.exit(1);
  }
}

// Schedule the headless sessions worker
scheduleHeadlessSessionsCheck();

// Note: Cron jobs are already initialized in initServices()

// Handle graceful shutdown
async function gracefulShutdown(signal: string) {
  log.info(`${signal} received, shutting down gracefully`);
  
  // First, stop accepting new HTTP connections
  httpServer.close(() => {
    log.info('HTTP server closed');
  });
  
  // Close WebSocket server to stop accepting new WS connections
  wss.close(() => {
    log.info('WebSocket server closed');
  });
  
  // Give in-flight requests a moment to complete
  await new Promise(resolve => setTimeout(resolve, 1000));
  
  // Close all pending sessions first (browsers waiting for client connection)
  if (pendingSessions.size > 0) {
    log.info(`Closing ${pendingSessions.size} pending session(s)...`);
    const pendingClosePromises: Promise<void>[] = [];
    
    for (const [sessionId, pendingSession] of pendingSessions.entries()) {
      pendingClosePromises.push(
        (async () => {
          try {
            await pendingSession.browser.close();
            log.info(`Closed pending browser session: ${sessionId} (was waiting for ${pendingSession.expectedClientId})`);
          } catch (error) {
            log.error(`Error closing pending browser session ${sessionId}: ${error}`);
          }
        })()
      );
    }
    
    await Promise.race([
      Promise.all(pendingClosePromises),
      new Promise(resolve => setTimeout(resolve, 5000))
    ]);
    
    pendingSessions.clear();
    log.info('All pending sessions closed');
  }
  
  // Close all established browser sessions
  if (browserSessions.size > 0) {
    log.info(`Closing ${browserSessions.size} browser session(s)...`);
    const closePromises: Promise<void>[] = [];
    
    for (const [sessionId, browser] of browserSessions.entries()) {
      closePromises.push(
        (async () => {
          try {
            await browser.close();
            log.info(`Closed browser session: ${sessionId}`);
          } catch (error) {
            log.error(`Error closing browser session ${sessionId}: ${error}`);
          }
        })()
      );
    }
    
    // Wait for all browsers to close with a timeout
    await Promise.race([
      Promise.all(closePromises),
      new Promise(resolve => setTimeout(resolve, 5000)) // 5 second timeout
    ]);
    
    browserSessions.clear();
    log.info('All browser sessions closed');
  }
  
  await closeRedis();
  process.exit(0);
}

process.on('SIGTERM', () => gracefulShutdown('SIGTERM'));
process.on('SIGINT', () => gracefulShutdown('SIGINT'));

// Initialize services and start server
initializeServices().catch(err => {
  log.error(`Failed to initialize services: ${err}`);
  process.exit(1);
});
