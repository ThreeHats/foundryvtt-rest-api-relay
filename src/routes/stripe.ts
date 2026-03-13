import express, { Request, Response } from 'express';
import { stripe, SUBSCRIPTION_PRICES, isStripeDisabled } from '../config/stripe';
import { User } from '../models/user';
import { authMiddleware } from '../middleware/auth';
import { log } from '../utils/logger';
import path from 'path';

const router = express.Router();

// Get subscription status
router.get('/status', authMiddleware, async (req: Request, res: Response) => {
  // If Stripe is disabled, return free tier status
  if (isStripeDisabled) {
    res.json({
      subscriptionStatus: 'free',
      subscriptionEndsAt: null
    });
    return;
  }
  
  try {
    const apiKey = req.headers['x-api-key'] as string;
    const user = await User.findOne({ where: { apiKey } });

    if (!user) {
      res.status(404).json({ error: 'User not found' });
      return;
    }

    res.json({
      subscriptionStatus: user.dataValues.subscriptionStatus || 'free',
      subscriptionEndsAt: user.dataValues.subscriptionEndsAt || null
    });
    return;
  } catch (error) {
    log.error(`Error getting subscription status: ${error}`);
    res.status(500).json({ error: 'Failed to get subscription status' });
    return;
  }
});

// Create checkout session
router.post('/create-checkout-session', authMiddleware, async (req: Request, res: Response) => {
  try {
    log.info('Creating checkout session');
    const apiKey = req.headers['x-api-key'] as string;
    const user = await User.findOne({ where: { apiKey } });

    if (!user) {
      res.status(404).json({ error: 'User not found' });
      return;
    }

    // Get or create Stripe customer
    let customerId = user.stripeCustomerId;
    log.info(`Creating checkout session for user ID: ${user.id}`);
    
    if (!customerId) {
      const customer = await stripe.customers.create({
        email: user.dataValues.email,
        metadata: { userId: user.dataValues.id.toString() }
      });
      
      customerId = customer.id;
      await user.update({ stripeCustomerId: customerId });
    }


    // Create checkout session
    const session = await stripe.checkout.sessions.create({
      customer: customerId,
      payment_method_types: ['card'],
      line_items: [
        {
          price: SUBSCRIPTION_PRICES.monthly,
          quantity: 1
        }
      ],
      mode: 'subscription',
      success_url: `${process.env.FRONTEND_URL}/api/subscriptions/subscription-success?session_id={CHECKOUT_SESSION_ID}`,
      cancel_url: `${process.env.FRONTEND_URL}/api/subscriptions/subscription-cancel`,
      metadata: { userId: user.dataValues.id.toString() }
    });

    res.json({ url: session.url });
  } catch (error) {
    log.error(`Error creating checkout session: ${error}`);
    res.status(500).json({ error: 'Failed to create checkout session' });
  }
});

// Update the create-portal-session route
router.post('/create-portal-session', authMiddleware, async (req: Request, res: Response) => {
  try {
    // You can still get user info if needed for analytics
    const apiKey = req.headers['x-api-key'] as string;
    const user = await User.findOne({ where: { apiKey } });

    if (!user) {
      res.status(404).json({ error: 'User not found' });
      return;
    }

    // Instead of creating a session, redirect to the shared portal URL
    const portalUrl = process.env.STRIPE_PORTAL_URL;
    
    // Log the redirect for tracking
    log.info(`Redirecting user ${user.id} to customer portal`);
    
    res.json({ url: portalUrl });
  } catch (error) {
    log.error(`Error handling portal redirect: ${error}`);
    res.status(500).json({ error: 'Failed to access customer portal' });
  }
});

// Handle subscription success
router.get('/subscription-success', (req: Request, res: Response) => {
  res.sendFile(path.join(__dirname, '../../public/subscription-success.html'));
});

// Handle subscription cancel
router.get('/subscription-cancel', (req: Request, res: Response) => {
  res.sendFile(path.join(__dirname, '../../public/subscription-cancel.html'));
});

export default router;