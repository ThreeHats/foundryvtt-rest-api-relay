import { Router, Request, Response, NextFunction } from 'express';
import bcrypt from 'bcryptjs';
import { User } from '../models/user';
import { PasswordResetToken } from '../models/passwordResetToken';
import crypto from 'crypto';
import { Op } from 'sequelize';
import { safeResponse } from './shared';
import { log } from '../utils/logger';
import { authRateLimiter, accountManagementRateLimiter, passwordResetRateLimiter } from '../middleware/rateLimiting';
import { stripe, isStripeDisabled } from '../config/stripe';
import { sendPasswordResetEmail } from '../services/email';

const router = Router();

// Check if registration is disabled (for staging/private environments)
const REGISTRATION_DISABLED = process.env.DISABLE_REGISTRATION === 'true';

/**
 * Validate password complexity:
 * - Minimum 8 characters
 * - At least one uppercase letter
 * - At least one lowercase letter
 * - At least one number
 */
function validatePassword(password: string): { valid: boolean; error?: string } {
  if (password.length < 8) {
    return { valid: false, error: 'Password must be at least 8 characters long' };
  }
  if (!/[A-Z]/.test(password)) {
    return { valid: false, error: 'Password must contain at least one uppercase letter' };
  }
  if (!/[a-z]/.test(password)) {
    return { valid: false, error: 'Password must contain at least one lowercase letter' };
  }
  if (!/[0-9]/.test(password)) {
    return { valid: false, error: 'Password must contain at least one number' };
  }
  return { valid: true };
}

// Register a new user
router.post('/register', authRateLimiter, async (req: Request, res: Response) => {
  // Block registration if disabled
  if (REGISTRATION_DISABLED) {
    log.warn('Registration attempt blocked - registration is disabled');
    safeResponse(res, 403, { error: 'Registration is disabled on this server' });
    return;
  }

  log.info('Registration attempt in auth.ts');
  try {
    const { email, password } = req.body;
    
    if (!email || !password) {
      log.warn('Missing email or password');
      safeResponse(res, 400, { error: 'Email and password are required' });
      return;
    }

    // Validate password complexity
    const passwordCheck = validatePassword(password);
    if (!passwordCheck.valid) {
      safeResponse(res, 400, { error: passwordCheck.error });
      return;
    }

    // Check if user already exists
    const existingUser = await User.findOne({ where: { email } });
    if (existingUser) {
      safeResponse(res, 409, { error: 'User already exists' });
      return;
    }
    
    log.info('Creating new user...');
    // Create a new user
    const user = await User.create({
      email,
      password, // Will be hashed by the beforeCreate hook
      apiKey: crypto.randomBytes(32).toString('hex'), // Explicitly generate an API key
      requestsThisMonth: 0
    });
    
    log.info(`User created with ID: ${user.getDataValue('id')}`);
    
    // Return the user (exclude password but include API key)
    res.status(201).json({
      id: user.getDataValue('id'),
      email: user.getDataValue('email'),
      apiKey: user.getDataValue('apiKey'),
      createdAt: user.getDataValue('createdAt'),
      subscriptionStatus: user.getDataValue('subscriptionStatus') || 'free'
    });
    return;
  } catch (error) {
    log.error('Registration error', { error });
    safeResponse(res, 500, { error: 'Registration failed' });
    return;
  }
});

// Login route - update the password comparison logic
router.post('/login', authRateLimiter, async (req: Request, res: Response) => {
  try {
    const { email, password } = req.body;
    
    if (!email || !password) {
      log.warn('Missing email or password');
      res.status(400).json({ error: 'Email and password are required' });
      return;
    }
    
    // Find the user
    const user = await User.findOne({ where: { email } });
    if (!user) {
      res.status(401).json({ error: 'Invalid credentials' });
      return;
    }
    
    try {
      // Get the stored hash directly from the data value
      const storedHash = user.getDataValue('password');
      
      const isPasswordValid = await bcrypt.compare(password, storedHash);
      
      if (!isPasswordValid) {
        res.status(401).json({ error: 'Invalid credentials' });
        return;
      }
      
      // Return the user (exclude password)
      res.status(200).json({
        id: user.getDataValue('id'),
        email: user.getDataValue('email'),
        apiKey: user.getDataValue('apiKey'),
        requestsThisMonth: user.getDataValue('requestsThisMonth'),
        createdAt: user.getDataValue('createdAt')
      });
      return;
    } catch (bcryptError) {
      log.error('bcrypt comparison error', { error: bcryptError });
      res.status(500).json({ error: 'Authentication error' });
      return;
    }
  } catch (error) {
    log.error('Login error', { error });
    res.status(500).json({ error: 'Login failed' });
    return;
  }
});

// Regenerate API key (for authenticated users)
router.post('/regenerate-key', authRateLimiter, async (req: Request, res: Response) => {
  try {
    const { email, password } = req.body;
    
    if (!email || !password) {
      res.status(400).json({ error: 'Email and password are required' });
      return;
    }
    
    // Find the user
    const user = await User.findOne({ where: { email } });
    if (!user) {
      res.status(401).json({ error: 'Invalid credentials' });
      return;
    }
    
    // Check password
    const isPasswordValid = await bcrypt.compare(password, user.password);
    if (!isPasswordValid) {
      res.status(401).json({ error: 'Invalid credentials' });
      return;
    }
    
    // Generate new API key
    const newApiKey = crypto.randomBytes(32).toString('hex');
    await user.update({ apiKey: newApiKey });
    
    // Return the new API key
    res.status(200).json({
      apiKey: newApiKey
    });
  } catch (error) {
    log.error('API key regeneration error', { error });
    res.status(500).json({ error: 'Failed to regenerate API key' });
  }
});

// Get user data (for authenticated users)
router.get('/user-data', async (req: Request, res: Response) => {
  try {
    // Get API key from header
    const apiKey = req.header('x-api-key');
    
    if (!apiKey) {
      res.status(401).json({ error: 'API key is required' });
      return;
    }
    
    // Find user by API key
    const user = await User.findOne({ where: { apiKey } });
    if (!user) {
      res.status(404).json({ error: 'User not found' });
      return;
    }
    
    // Return user data (exclude sensitive information)
    res.status(200).json({
      id: user.getDataValue('id'),
      email: user.getDataValue('email'),
      apiKey: user.getDataValue('apiKey'),
      requestsThisMonth: user.getDataValue('requestsThisMonth'),
      requestsToday: user.getDataValue('requestsToday') || 0,
      subscriptionStatus: user.getDataValue('subscriptionStatus') || 'free',
      limits: {
        dailyLimit: parseInt(process.env.DAILY_REQUEST_LIMIT || '1000'),
        monthlyLimit: parseInt(process.env.FREE_API_REQUESTS_LIMIT || '100'),
        unlimitedMonthly: (user.getDataValue('subscriptionStatus') === 'active')
      }
    });
    return;
  } catch (error) {
    log.error('Error fetching user data', { error });
    res.status(500).json({ error: 'Failed to fetch user data' });
    return;
  }
});

// Export user data (GDPR data portability)
router.get('/export-data', accountManagementRateLimiter, async (req: Request, res: Response) => {
  try {
    const apiKey = req.header('x-api-key');
    
    if (!apiKey) {
      res.status(401).json({ error: 'API key is required' });
      return;
    }
    
    const user = await User.findOne({ where: { apiKey } });
    if (!user) {
      res.status(404).json({ error: 'User not found' });
      return;
    }
    
    // Export all user data (excluding password hash for security)
    const exportData = {
      exportDate: new Date().toISOString(),
      user: {
        id: user.getDataValue('id'),
        email: user.getDataValue('email'),
        createdAt: user.getDataValue('createdAt'),
        updatedAt: user.getDataValue('updatedAt'),
      },
      subscription: {
        status: user.getDataValue('subscriptionStatus') || 'free',
        stripeCustomerId: user.getDataValue('stripeCustomerId') || null,
        subscriptionId: user.getDataValue('subscriptionId') || null,
      },
      usage: {
        requestsThisMonth: user.getDataValue('requestsThisMonth') || 0,
        requestsToday: user.getDataValue('requestsToday') || 0,
        lastRequestDate: user.getDataValue('lastRequestDate') || null,
      },
      apiAccess: {
        // Note: API key is included for user reference but should be regenerated if compromised
        apiKey: user.getDataValue('apiKey'),
      }
    };
    
    log.info(`Data export requested for user ID ${user.getDataValue('id')}`);
    
    res.status(200).json(exportData);
    return;
  } catch (error) {
    log.error('Error exporting user data', { error });
    res.status(500).json({ error: 'Failed to export user data' });
    return;
  }
});

// Delete user account (GDPR right to be forgotten)
router.delete('/account', accountManagementRateLimiter, async (req: Request, res: Response) => {
  try {
    const apiKey = req.header('x-api-key');
    const { confirmEmail, password } = req.body;
    
    if (!apiKey) {
      res.status(401).json({ error: 'API key is required' });
      return;
    }
    
    const user = await User.findOne({ where: { apiKey } });
    if (!user) {
      res.status(404).json({ error: 'User not found' });
      return;
    }
    
    // Require email confirmation to prevent accidental deletion
    const userEmail = user.getDataValue('email');
    if (!confirmEmail || confirmEmail !== userEmail) {
      res.status(400).json({ 
        error: 'Email confirmation required',
        message: 'Please provide your email address in the confirmEmail field to confirm account deletion'
      });
      return;
    }
    
    // Require password verification for extra security
    if (!password) {
      res.status(400).json({ 
        error: 'Password required',
        message: 'Please provide your password to confirm account deletion'
      });
      return;
    }
    
    const passwordHash = user.getDataValue('password');
    const passwordValid = await bcrypt.compare(password, passwordHash);
    if (!passwordValid) {
      res.status(401).json({ error: 'Invalid password' });
      return;
    }
    
    const userId = user.getDataValue('id');
    
    // If user has a Stripe subscription, cancel it before deleting the account
    const subscriptionStatus = user.getDataValue('subscriptionStatus');
    const subscriptionId = user.getDataValue('subscriptionId');
    const hasActiveSubscription = subscriptionStatus === 'active';
    let subscriptionCancelled = false;
    
    if (hasActiveSubscription && subscriptionId && !isStripeDisabled && !stripe.disabled) {
      try {
        await stripe.subscriptions.cancel(subscriptionId);
        subscriptionCancelled = true;
        log.info(`Cancelled Stripe subscription ${subscriptionId} for user ID ${userId}`);
      } catch (stripeError) {
        log.error('Failed to cancel Stripe subscription', { 
          error: stripeError instanceof Error ? stripeError.message : stripeError,
          subscriptionId,
          userId
        });
        // Continue with account deletion even if subscription cancellation fails
        // The subscription will eventually fail due to no associated user
      }
    }
    
    // Delete the user
    await user.destroy();
    
    log.info(`Account deleted for user ID ${userId}`);
    
    res.status(200).json({ 
      success: true,
      message: 'Account successfully deleted',
      subscriptionCancelled: hasActiveSubscription ? subscriptionCancelled : undefined,
      note: hasActiveSubscription && subscriptionCancelled
        ? 'Your Stripe subscription has been cancelled.'
        : hasActiveSubscription && !subscriptionCancelled
        ? 'We could not automatically cancel your subscription. Please contact support if you continue to be charged.'
        : undefined
    });
    return;
  } catch (error) {
    log.error('Error deleting account', { 
      error: error instanceof Error ? { message: error.message, stack: error.stack } : error 
    });
    res.status(500).json({ error: 'Failed to delete account' });
    return;
  }
});

// Request a password reset
router.post('/forgot-password', passwordResetRateLimiter, async (req: Request, res: Response) => {
  try {
    const { email } = req.body;

    // Always return the same response to prevent email enumeration
    const genericMessage = 'If an account with that email exists, a password reset link has been sent.';

    if (!email) {
      safeResponse(res, 400, { error: 'Email is required' });
      return;
    }

    const user = await User.findOne({ where: { email } });
    if (!user) {
      // Don't reveal that the email doesn't exist
      safeResponse(res, 200, { message: genericMessage });
      return;
    }

    const userId = user.getDataValue('id');

    // Invalidate any prior tokens for this user
    await PasswordResetToken.invalidateForUser(userId);

    // Generate a 32-byte random token
    const rawToken = crypto.randomBytes(32).toString('hex');
    const tokenHash = crypto.createHash('sha256').update(rawToken).digest('hex');

    // Store the hashed token with 1-hour expiry
    await PasswordResetToken.create({
      userId,
      tokenHash,
      expiresAt: new Date(Date.now() + 60 * 60 * 1000),
      used: false
    });

    // Send email (fire-and-forget)
    sendPasswordResetEmail(user.getDataValue('email'), rawToken);

    safeResponse(res, 200, { message: genericMessage });
    return;
  } catch (error) {
    log.error('Forgot password error', { error });
    safeResponse(res, 500, { error: 'An error occurred processing your request' });
    return;
  }
});

// Reset password with token
router.post('/reset-password', passwordResetRateLimiter, async (req: Request, res: Response) => {
  try {
    const { token, password } = req.body;

    if (!token || !password) {
      safeResponse(res, 400, { error: 'Token and password are required' });
      return;
    }

    // Validate password complexity
    const passwordCheck = validatePassword(password);
    if (!passwordCheck.valid) {
      safeResponse(res, 400, { error: passwordCheck.error });
      return;
    }

    // Hash the incoming token and look it up
    const tokenHash = crypto.createHash('sha256').update(token).digest('hex');
    const resetToken = await PasswordResetToken.findOne({
      where: {
        tokenHash,
        used: false,
        expiresAt: { [Op.gt]: new Date() }
      }
    });

    if (!resetToken) {
      safeResponse(res, 400, { error: 'Invalid or expired reset token' });
      return;
    }

    const userId = resetToken.getDataValue('userId');
    const user = await User.findOne({ where: { id: userId } });

    if (!user) {
      safeResponse(res, 400, { error: 'Invalid or expired reset token' });
      return;
    }

    // Update password (beforeUpdate hook handles bcrypt hashing)
    if (process.env.DB_TYPE === 'memory') {
      const salt = await bcrypt.genSalt(10);
      const hashedPassword = await bcrypt.hash(password, salt);
      user.password = hashedPassword;
      user.updatedAt = new Date();
      // Update in memory store
      const memoryStore = (await import('../sequelize')).sequelize as any;
      memoryStore.users.set(user.email, user);
    } else {
      await user.update({ password });
    }

    // Mark token as used
    if (process.env.DB_TYPE === 'memory') {
      await resetToken.update({ used: true });
    } else {
      await resetToken.update({ used: true });
    }

    // Invalidate all other tokens for this user
    await PasswordResetToken.invalidateForUser(userId);

    log.info('Password reset successful', { userId });
    safeResponse(res, 200, { message: 'Password has been reset successfully' });
    return;
  } catch (error) {
    log.error('Reset password error', { error });
    safeResponse(res, 500, { error: 'An error occurred processing your request' });
    return;
  }
});

// Validate a reset token (lets frontend check before showing form)
router.get('/validate-reset-token/:token', passwordResetRateLimiter, async (req: Request, res: Response) => {
  try {
    const token = req.params.token as string;

    if (!token) {
      res.status(400).json({ valid: false });
      return;
    }

    const tokenHash = crypto.createHash('sha256').update(token).digest('hex');
    const resetToken = await PasswordResetToken.findOne({
      where: {
        tokenHash,
        used: false,
        expiresAt: { [Op.gt]: new Date() }
      }
    });

    res.status(200).json({ valid: !!resetToken });
    return;
  } catch (error) {
    log.error('Validate reset token error', { error });
    res.status(200).json({ valid: false });
    return;
  }
});

export default router;