import rateLimit from 'express-rate-limit';
import { log } from '../utils/logger';

/**
 * Rate limiter for authentication endpoints (login, register, API key regeneration)
 * Prevents brute force attacks on authentication
 */
export const authRateLimiter = rateLimit({
  windowMs: 15 * 60 * 1000, // 15 minutes
  max: 5, // Limit each IP to 5 requests per windowMs
  standardHeaders: true, // Return rate limit info in the `RateLimit-*` headers
  legacyHeaders: false, // Disable the `X-RateLimit-*` headers
  handler: (req, res) => {
    log.warn(`Rate limit exceeded for IP ${req.ip} on ${req.path}`);
    res.status(429).json({
      error: 'Too many authentication attempts from this IP, please try again after 15 minutes'
    });
  }
});

/**
 * Rate limiter for account management endpoints (account deletion, data export)
 * More lenient than auth endpoints but still protected
 */
/**
 * Rate limiter for password reset endpoints
 * Stricter than auth to prevent spam relay abuse
 */
export const passwordResetRateLimiter = rateLimit({
  windowMs: 60 * 60 * 1000, // 1 hour
  max: 3, // Limit each IP to 3 requests per hour
  standardHeaders: true,
  legacyHeaders: false,
  handler: (req, res) => {
    log.warn(`Password reset rate limit exceeded for IP ${req.ip} on ${req.path}`);
    res.status(429).json({
      error: 'Too many password reset attempts from this IP, please try again later'
    });
  }
});

export const accountManagementRateLimiter = rateLimit({
  windowMs: 60 * 60 * 1000, // 1 hour
  max: 10, // Limit each IP to 10 requests per hour
  standardHeaders: true,
  legacyHeaders: false,
  handler: (req, res) => {
    log.warn(`Account management rate limit exceeded for IP ${req.ip} on ${req.path}`);
    res.status(429).json({
      error: 'Too many account management requests from this IP, please try again later'
    });
  }
});
