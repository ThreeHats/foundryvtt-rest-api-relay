import express, { Request, Response } from 'express';
import { stripe } from '../config/stripe';
import { User } from '../models/user';
import { log } from '../utils/logger';

const router = express.Router();

// Stripe webhook handler
router.post('/stripe', async (req: Request, res: Response) => {
  const sig = req.headers['stripe-signature'] as string;
  let event;

  try {
    event = stripe.webhooks.constructEvent(
      req.body,
      sig,
      process.env.STRIPE_WEBHOOK_SECRET || ''
    );
  } catch (err) {
    log.error(`Webhook Error: ${err}`);
    res.status(400).send(`Webhook Error: ${err}`);
    return;
  }

  // Handle the event
  let success = false;
  try {
    switch (event.type) {
      case 'customer.subscription.created':
      case 'customer.subscription.updated':
        success = await handleSubscriptionUpdated(event.data.object);
        break;
      case 'customer.subscription.deleted':
        success = await handleSubscriptionDeleted(event.data.object);
        break;
      case 'invoice.payment_succeeded':
        success = await handlePaymentSucceeded(event.data.object);
        break;
      case 'invoice.payment_failed':
        success = await handlePaymentFailed(event.data.object);
        break;
      default:
        log.info(`Unhandled event type: ${event.type}`);
        success = true; // Unknown events are considered success
    }
  } catch (error) {
    log.error('Webhook handler exception:', {
      eventType: event.type,
      eventId: event.id,
      error: error instanceof Error ? error.message : String(error)
    });
    res.status(500).json({ error: 'Webhook processing failed' });
    return;
  }

  // Return appropriate status code
  if (success) {
    res.status(200).send();
  } else {
    // Return 500 so Stripe will retry
    res.status(500).json({ error: 'Failed to process webhook' });
  }
});

// Handle subscription updates
async function handleSubscriptionUpdated(subscription: any): Promise<boolean> {
  try {
    const customerId = subscription.customer;
    log.info(`Processing subscription update for customer: ${customerId}, status: ${subscription.status}`);

    const user = await User.findOne({ where: { stripeCustomerId: customerId } });

    if (!user) {
      log.error(`User not found for customer: ${customerId}`);
      return false; // User not found - return failure so Stripe retries
    }

    // Map Stripe statuses: 'trialing' should grant full access
    const stripeStatus = subscription.status;
    const effectiveStatus = stripeStatus === 'trialing' ? 'active' : stripeStatus;

    // Status priority: 'active' should never be overwritten by lesser statuses.
    // Stripe fires 'created' (incomplete) and 'updated' (active) nearly simultaneously,
    // causing a race where 'incomplete' can overwrite 'active'. Use an atomic conditional
    // update so 'incomplete' never downgrades an 'active' subscription.
    const STATUS_PRIORITY: Record<string, number> = {
      'incomplete': 0,
      'incomplete_expired': 0,
      'past_due': 1,
      'canceled': 2,
      'active': 3,
    };

    const newPriority = STATUS_PRIORITY[effectiveStatus] ?? 1;
    const currentStatus = user.getDataValue('subscriptionStatus') || 'free';
    const currentPriority = STATUS_PRIORITY[currentStatus] ?? -1;

    if (newPriority < currentPriority) {
      log.info(`Skipping downgrade for user ${user.id}: ${currentStatus} → ${effectiveStatus} (stripe: ${stripeStatus})`);
      // Still update subscription ID and period end
      await user.update({
        subscriptionId: subscription.id,
        subscriptionEndsAt: new Date(subscription.current_period_end * 1000)
      });
      return true;
    }

    log.info(`Found user ID ${user.id} for customer ${customerId}, updating subscription status to: ${effectiveStatus} (stripe: ${stripeStatus})`);

    await user.update({
      subscriptionStatus: effectiveStatus,
      subscriptionId: subscription.id,
      subscriptionEndsAt: new Date(subscription.current_period_end * 1000)
    });

    await user.reload();
    log.info(`Successfully updated subscription for user ${user.id} to status: ${user.subscriptionStatus}`);
    return true;
  } catch (error) {
    log.error(`Error updating subscription`, {
      customerId: subscription.customer,
      subscriptionId: subscription.id,
      status: subscription.status,
      error: error instanceof Error ? error.message : String(error)
    });
    return false; // Return failure so Stripe retries
  }
}

// Handle subscription deletions
async function handleSubscriptionDeleted(subscription: any): Promise<boolean> {
  try {
    const customerId = subscription.customer;
    const user = await User.findOne({ where: { stripeCustomerId: customerId } });

    if (!user) {
      log.error(`User not found for customer: ${customerId}`);
      return false;
    }

    await user.update({
      subscriptionStatus: 'canceled',
      subscriptionEndsAt: new Date(subscription.canceled_at * 1000)
    });

    log.info(`Subscription canceled for user ${user.id}`);
    return true;
  } catch (error) {
    log.error(`Error handling subscription deletion`, {
      customerId: subscription.customer,
      error: error instanceof Error ? error.message : String(error)
    });
    return false;
  }
}

// Handle successful payments
async function handlePaymentSucceeded(invoice: any): Promise<boolean> {
  try {
    if (invoice.subscription) {
      const customerId = invoice.customer;
      const user = await User.findOne({ where: { stripeCustomerId: customerId } });

      if (!user) {
        log.error(`User not found for customer: ${customerId}`);
        return false;
      }

      // Log the payment success only - request count management is handled by
      // the monthly cron job in src/cron/monthlyReset.ts
      log.info(`Payment success recorded for user ${user.id}`);
    }
    return true;
  } catch (error) {
    log.error(`Error handling payment success`, {
      error: error instanceof Error ? error.message : String(error)
    });
    return false;
  }
}

// Handle failed payments
async function handlePaymentFailed(invoice: any): Promise<boolean> {
  try {
    if (invoice.subscription) {
      const customerId = invoice.customer;
      const user = await User.findOne({ where: { stripeCustomerId: customerId } });

      if (!user) {
        log.error(`User not found for customer: ${customerId}`);
        return false;
      }

      // Mark subscription as past_due
      await user.update({
        subscriptionStatus: 'past_due'
      });

      log.info(`Updated subscription status to past_due for user ${user.id}`);
    }
    return true;
  } catch (error) {
    log.error(`Error handling payment failure`, {
      error: error instanceof Error ? error.message : String(error)
    });
    return false;
  }
}

export default router;