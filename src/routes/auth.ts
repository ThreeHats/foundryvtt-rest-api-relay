import { Router, Request, Response, NextFunction } from 'express';
import bcrypt from 'bcryptjs';
import { User } from '../models/user';
import { ApiKey } from '../models/apiKey';
import { PasswordResetToken } from '../models/passwordResetToken';
import crypto from 'crypto';
import { Op } from 'sequelize';
import { safeResponse } from './shared';
import { log } from '../utils/logger';
import { authRateLimiter, accountManagementRateLimiter, passwordResetRateLimiter } from '../middleware/rateLimiting';
import { authMiddleware } from '../middleware/auth';
import { stripe, isStripeDisabled } from '../config/stripe';
import { sendPasswordResetEmail } from '../services/email';
import { encrypt, isEncryptionAvailable } from '../utils/encryption';

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

/**
 * Register a new user account
 *
 * Creates a new user with the provided email and password. Returns the user's API key.
 *
 * @route POST /auth/register
 * @group Auth
 * @param {string} email - [body] The user's email address
 * @param {string} password - [body] The user's password (min 8 chars, must include upper, lower, number)
 * @returns {object} Created user object with id, email, apiKey, and subscriptionStatus
 * @security none
 */
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

/**
 * Log in with email and password
 *
 * Authenticates a user and returns their account details including API key.
 *
 * @route POST /auth/login
 * @group Auth
 * @param {string} email - [body] The user's email address
 * @param {string} password - [body] The user's password
 * @returns {object} User object with id, email, apiKey, and requestsThisMonth
 * @security none
 */
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

/**
 * Regenerate API key
 *
 * Generates a new API key for the authenticated user. Requires email and password confirmation.
 *
 * @route POST /auth/regenerate-key
 * @group Auth
 * @param {string} email - [body] The user's email address
 * @param {string} password - [body] The user's password
 * @returns {object} Object containing the new apiKey
 * @security none
 */
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
    const userId = user.getDataValue ? user.getDataValue('id') : user.id;
    await user.update({ apiKey: newApiKey });

    // Cascade: delete all scoped sub-keys when master key is regenerated
    const deletedCount = await ApiKey.deleteAllByUser(userId);
    if (deletedCount > 0) {
      log.info(`Deleted ${deletedCount} scoped API keys for user ${userId} due to master key regeneration`);
    }

    // Return the new API key
    res.status(200).json({
      apiKey: newApiKey
    });
  } catch (error) {
    log.error('API key regeneration error', { error });
    res.status(500).json({ error: 'Failed to regenerate API key' });
  }
});

/**
 * Get user data
 *
 * Retrieves the authenticated user's account details and usage information.
 *
 * @route GET /auth/user-data
 * @group Auth
 * @returns {object} User data including id, email, usage stats, and subscription info
 */
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

/**
 * Export user data
 *
 * Exports all stored user data for GDPR data portability compliance.
 *
 * @route GET /auth/export-data
 * @group Auth
 * @returns {object} Complete user data export including account, subscription, and usage data
 */
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
    
    // Get scoped keys metadata (no key values or credentials)
    const userId = user.getDataValue('id');
    const scopedKeys = await ApiKey.findAllByUser(userId);
    const scopedKeysExport = scopedKeys.map((k: any) => {
      const get = (field: string) => k.getDataValue ? k.getDataValue(field) : k[field];
      return {
        id: get('id'),
        name: get('name'),
        scopedClientId: get('scopedClientId'),
        scopedUserId: get('scopedUserId'),
        dailyLimit: get('dailyLimit'),
        expiresAt: get('expiresAt'),
        enabled: get('enabled'),
        hasFoundryCredentials: !!(get('encryptedFoundryPassword')),
        createdAt: get('createdAt'),
      };
    });

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
      },
      scopedKeys: scopedKeysExport,
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

/**
 * Delete user account
 *
 * Permanently deletes the user's account and cancels any active subscriptions.
 * Requires email confirmation and password verification.
 *
 * @route DELETE /auth/account
 * @group Auth
 * @param {string} confirmEmail - [body] The user's email address (must match account email)
 * @param {string} password - [body] The user's password for verification
 * @returns {object} Deletion confirmation with subscription cancellation status
 */
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
    
    // Cascade: delete all scoped API keys before deleting user
    const deletedKeys = await ApiKey.deleteAllByUser(userId);
    if (deletedKeys > 0) {
      log.info(`Deleted ${deletedKeys} scoped API keys for user ${userId} during account deletion`);
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

/**
 * Change password while logged in
 *
 * Allows an authenticated user to change their password by providing their current password and a new one.
 *
 * @route POST /auth/change-password
 * @group Auth
 * @param {string} currentPassword - [body] The user's current password
 * @param {string} newPassword - [body] The new password (min 8 chars, must include upper, lower, number)
 * @returns {object} Success message
 */
router.post('/change-password', accountManagementRateLimiter, async (req: Request, res: Response) => {
  try {
    const apiKey = req.header('x-api-key');

    if (!apiKey) {
      res.status(401).json({ error: 'API key is required' });
      return;
    }

    const { currentPassword, newPassword } = req.body;

    if (!currentPassword || !newPassword) {
      res.status(400).json({ error: 'Current password and new password are required' });
      return;
    }

    const user = await User.findOne({ where: { apiKey } });
    if (!user) {
      res.status(404).json({ error: 'User not found' });
      return;
    }

    // Verify current password
    const storedHash = user.getDataValue('password');
    const isPasswordValid = await bcrypt.compare(currentPassword, storedHash);
    if (!isPasswordValid) {
      res.status(401).json({ error: 'Current password is incorrect' });
      return;
    }

    // Validate new password complexity
    const passwordCheck = validatePassword(newPassword);
    if (!passwordCheck.valid) {
      res.status(400).json({ error: passwordCheck.error });
      return;
    }

    // Update password
    if (process.env.DB_TYPE === 'memory') {
      const salt = await bcrypt.genSalt(10);
      const hashedPassword = await bcrypt.hash(newPassword, salt);
      user.password = hashedPassword;
      user.updatedAt = new Date();
      const memoryStore = (await import('../sequelize')).sequelize as any;
      memoryStore.users.set(user.getDataValue('email'), user);
    } else {
      await user.update({ password: newPassword });
    }

    log.info('Password changed successfully', { userId: user.getDataValue('id') });
    res.status(200).json({ message: 'Password changed successfully' });
    return;
  } catch (error) {
    log.error('Change password error', { error });
    res.status(500).json({ error: 'Failed to change password' });
    return;
  }
});

/**
 * Request a password reset
 *
 * Sends a password reset email if the account exists. Always returns success to prevent email enumeration.
 *
 * @route POST /auth/forgot-password
 * @group Auth
 * @param {string} email - [body] The email address associated with the account
 * @returns {object} Generic success message
 * @security none
 */
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

    const response: any = { message: genericMessage };
    if (process.env.RETURN_RESET_TOKEN === 'true') {
      response.token = rawToken;
    }
    safeResponse(res, 200, response);
    return;
  } catch (error) {
    log.error('Forgot password error', { error });
    safeResponse(res, 500, { error: 'An error occurred processing your request' });
    return;
  }
});

/**
 * Reset password with token
 *
 * Resets the user's password using a valid reset token from the forgot-password flow.
 *
 * @route POST /auth/reset-password
 * @group Auth
 * @param {string} token - [body] The password reset token from the email
 * @param {string} password - [body] The new password (min 8 chars, must include upper, lower, number)
 * @returns {object} Success message
 * @security none
 */
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

/**
 * Validate a password reset token
 *
 * Checks whether a password reset token is still valid and unused.
 *
 * @route GET /auth/validate-reset-token/:token
 * @group Auth
 * @param {string} token - [params] The password reset token to validate
 * @returns {object} Object with a boolean `valid` field
 * @security none
 */
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

// ==================== Scoped API Key CRUD ====================

/**
 * Create a new scoped API key
 *
 * Creates a sub-key under the authenticated user's master key with optional scope restrictions.
 *
 * @route POST /auth/api-keys
 * @group Auth - Scoped Keys
 * @param {string} name - [body] Friendly name for the key
 * @param {string} scopedClientId - [body,?] Lock to specific Foundry client ID
 * @param {string} scopedUserId - [body,?] Lock to specific Foundry user ID
 * @param {number} dailyLimit - [body,?] Per-key daily request cap
 * @param {string} expiresAt - [body,?] Expiry timestamp (ISO 8601)
 * @param {string} foundryUrl - [body,?] Foundry instance URL for headless sessions
 * @param {string} foundryUsername - [body,?] Foundry login username
 * @param {string} foundryPassword - [body,?] Foundry login password (encrypted at rest)
 * @returns {object} Created key with full token (shown only once)
 */
router.post('/api-keys', authMiddleware, async (req: Request, res: Response) => {
  try {
    const apiKey = req.header('x-api-key') as string;

    // Only master keys can create scoped keys
    if (req.scopedKey) {
      safeResponse(res, 403, { error: 'Scoped keys cannot create other scoped keys. Use your master API key.' });
      return;
    }

    const user = req.user;
    const userId = user.getDataValue ? user.getDataValue('id') : user.id;
    const { name, scopedClientId, scopedUserId, dailyLimit, expiresAt, foundryUrl, foundryUsername, foundryPassword } = req.body;

    if (!name) {
      safeResponse(res, 400, { error: 'Name is required' });
      return;
    }

    // Prepare credential fields
    let encryptedFoundryPassword: string | null = null;
    let passwordIv: string | null = null;
    let passwordAuthTag: string | null = null;

    if (foundryPassword) {
      if (!isEncryptionAvailable()) {
        safeResponse(res, 400, { error: 'Credential storage is not available. CREDENTIALS_ENCRYPTION_KEY is not configured.' });
        return;
      }
      const encrypted = encrypt(foundryPassword);
      encryptedFoundryPassword = encrypted.ciphertext;
      passwordIv = encrypted.iv;
      passwordAuthTag = encrypted.authTag;
    }

    const key = crypto.randomBytes(32).toString('hex');

    const record = await ApiKey.create({
      userId,
      key,
      name,
      scopedClientId: scopedClientId || null,
      scopedUserId: scopedUserId || null,
      dailyLimit: dailyLimit ? parseInt(dailyLimit, 10) : null,
      expiresAt: expiresAt ? new Date(expiresAt) : null,
      foundryUrl: foundryUrl || null,
      foundryUsername: foundryUsername || null,
      encryptedFoundryPassword,
      passwordIv,
      passwordAuthTag,
    });

    const recordId = record.getDataValue ? record.getDataValue('id') : record.id;

    log.info(`Created scoped API key "${name}" (ID: ${recordId}) for user ${userId}`);

    safeResponse(res, 201, {
      id: recordId,
      key, // Full key shown only on creation
      name,
      scopedClientId: scopedClientId || null,
      scopedUserId: scopedUserId || null,
      dailyLimit: dailyLimit ? parseInt(dailyLimit, 10) : null,
      expiresAt: expiresAt || null,
      hasFoundryCredentials: !!foundryPassword,
      enabled: true,
      createdAt: record.getDataValue ? record.getDataValue('createdAt') : record.createdAt,
    });
    return;
  } catch (error) {
    log.error('Error creating scoped API key', { error });
    safeResponse(res, 500, { error: 'Failed to create API key' });
    return;
  }
});

/**
 * List all scoped API keys
 *
 * Returns all scoped keys for the authenticated user. Keys are masked.
 *
 * @route GET /auth/api-keys
 * @group Auth - Scoped Keys
 * @returns {object} Array of scoped keys with masked tokens
 */
router.get('/api-keys', authMiddleware, async (req: Request, res: Response) => {
  try {
    // Only master keys can list scoped keys
    if (req.scopedKey) {
      safeResponse(res, 403, { error: 'Use your master API key to manage scoped keys.' });
      return;
    }

    const user = req.user;
    const userId = user.getDataValue ? user.getDataValue('id') : user.id;

    const keys = await ApiKey.findAllByUser(userId);

    const result = keys.map((k: any) => {
      const get = (field: string) => k.getDataValue ? k.getDataValue(field) : k[field];
      const keyVal = get('key');
      const expiresAtVal = get('expiresAt');
      return {
        id: get('id'),
        key: keyVal ? keyVal.substring(0, 8) + '...' : null,
        name: get('name'),
        scopedClientId: get('scopedClientId'),
        scopedUserId: get('scopedUserId'),
        dailyLimit: get('dailyLimit'),
        requestsToday: get('requestsToday') || 0,
        expiresAt: expiresAtVal,
        isExpired: expiresAtVal ? new Date(expiresAtVal) < new Date() : false,
        enabled: get('enabled'),
        hasFoundryCredentials: !!(get('encryptedFoundryPassword')),
        foundryUrl: get('foundryUrl'),
        foundryUsername: get('foundryUsername'),
        createdAt: get('createdAt'),
        updatedAt: get('updatedAt'),
      };
    });

    safeResponse(res, 200, { keys: result });
    return;
  } catch (error) {
    log.error('Error listing scoped API keys', { error });
    safeResponse(res, 500, { error: 'Failed to list API keys' });
    return;
  }
});

/**
 * Update a scoped API key
 *
 * Update name, scopes, limits, enabled status, expiry, or credentials.
 *
 * @route PATCH /auth/api-keys/:id
 * @group Auth - Scoped Keys
 * @param {number} id - [params] The scoped key ID
 * @returns {object} Updated key data
 */
router.patch('/api-keys/:id', authMiddleware, async (req: Request, res: Response) => {
  try {
    if (req.scopedKey) {
      safeResponse(res, 403, { error: 'Use your master API key to manage scoped keys.' });
      return;
    }

    const user = req.user;
    const userId = user.getDataValue ? user.getDataValue('id') : user.id;
    const keyId = parseInt(req.params.id as string, 10);

    const record = await ApiKey.findOne({ where: { id: keyId, userId } });
    if (!record) {
      safeResponse(res, 404, { error: 'API key not found' });
      return;
    }

    const updates: any = {};
    const { name, scopedClientId, scopedUserId, dailyLimit, expiresAt, enabled, foundryUrl, foundryUsername, foundryPassword } = req.body;

    if (name !== undefined) updates.name = name;
    if (scopedClientId !== undefined) updates.scopedClientId = scopedClientId || null;
    if (scopedUserId !== undefined) updates.scopedUserId = scopedUserId || null;
    if (dailyLimit !== undefined) updates.dailyLimit = dailyLimit ? parseInt(dailyLimit, 10) : null;
    if (expiresAt !== undefined) updates.expiresAt = expiresAt ? new Date(expiresAt) : null;
    if (enabled !== undefined) updates.enabled = !!enabled;
    if (foundryUrl !== undefined) updates.foundryUrl = foundryUrl || null;
    if (foundryUsername !== undefined) updates.foundryUsername = foundryUsername || null;

    // Handle password update
    if (foundryPassword !== undefined) {
      if (foundryPassword === null) {
        // Clear credentials
        updates.encryptedFoundryPassword = null;
        updates.passwordIv = null;
        updates.passwordAuthTag = null;
      } else {
        if (!isEncryptionAvailable()) {
          safeResponse(res, 400, { error: 'Credential storage is not available. CREDENTIALS_ENCRYPTION_KEY is not configured.' });
          return;
        }
        const encrypted = encrypt(foundryPassword);
        updates.encryptedFoundryPassword = encrypted.ciphertext;
        updates.passwordIv = encrypted.iv;
        updates.passwordAuthTag = encrypted.authTag;
      }
    }

    if (record.update && typeof record.update === 'function') {
      await record.update(updates);
    }

    const get = (field: string) => record.getDataValue ? record.getDataValue(field) : record[field];

    log.info(`Updated scoped API key ID ${keyId} for user ${userId}`);

    safeResponse(res, 200, {
      id: get('id'),
      name: get('name'),
      scopedClientId: get('scopedClientId'),
      scopedUserId: get('scopedUserId'),
      dailyLimit: get('dailyLimit'),
      expiresAt: get('expiresAt'),
      enabled: get('enabled'),
      hasFoundryCredentials: !!(get('encryptedFoundryPassword')),
      foundryUrl: get('foundryUrl'),
      foundryUsername: get('foundryUsername'),
      updatedAt: get('updatedAt'),
    });
    return;
  } catch (error) {
    log.error('Error updating scoped API key', { error });
    safeResponse(res, 500, { error: 'Failed to update API key' });
    return;
  }
});

/**
 * Delete a scoped API key
 *
 * Permanently deletes a scoped key.
 *
 * @route DELETE /auth/api-keys/:id
 * @group Auth - Scoped Keys
 * @param {number} id - [params] The scoped key ID
 * @returns {object} Deletion confirmation
 */
router.delete('/api-keys/:id', authMiddleware, async (req: Request, res: Response) => {
  try {
    if (req.scopedKey) {
      safeResponse(res, 403, { error: 'Use your master API key to manage scoped keys.' });
      return;
    }

    const user = req.user;
    const userId = user.getDataValue ? user.getDataValue('id') : user.id;
    const keyId = parseInt(req.params.id as string, 10);

    const record = await ApiKey.findOne({ where: { id: keyId, userId } });
    if (!record) {
      safeResponse(res, 404, { error: 'API key not found' });
      return;
    }

    if (record.destroy && typeof record.destroy === 'function') {
      await record.destroy();
    }

    log.info(`Deleted scoped API key ID ${keyId} for user ${userId}`);

    safeResponse(res, 200, { success: true, message: 'API key deleted' });
    return;
  } catch (error) {
    log.error('Error deleting scoped API key', { error });
    safeResponse(res, 500, { error: 'Failed to delete API key' });
    return;
  }
});

export default router;