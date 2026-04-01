import Stripe from 'stripe';
import { log } from '../utils/logger';

// Stripe is disabled only when no secret key is provided
const isStripeDisabled = !process.env.STRIPE_SECRET_KEY;

// Initialize Stripe conditionally
let stripe: any;
const SUBSCRIPTION_PRICES = {
  monthly: process.env.STRIPE_PRICE_ID || '' // Your Stripe price ID for monthly subscription
};

if (isStripeDisabled) {
  log.info('Stripe disabled — STRIPE_SECRET_KEY not provided');
  stripe = {
    disabled: true,
    customers: { create: async () => ({ id: 'disabled' }) },
    checkout: { sessions: { create: async () => ({ url: '#' }) } },
    webhooks: { constructEvent: () => ({ type: 'disabled', data: { object: {} } }) }
  };
} else {
  stripe = new Stripe(process.env.STRIPE_SECRET_KEY, {
    apiVersion: '2025-02-24.acacia'
  });
}

export { stripe, SUBSCRIPTION_PRICES, isStripeDisabled };