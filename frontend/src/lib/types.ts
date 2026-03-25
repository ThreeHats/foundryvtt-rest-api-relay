export interface UserData {
  id: number;
  email: string;
  apiKey: string;
  subscriptionStatus: 'free' | 'active' | 'past_due';
  subscriptionEndsAt?: string;
  requestsToday: number;
  requestsThisMonth: number;
  limits?: {
    dailyLimit: number;
    monthlyLimit: number;
    unlimitedMonthly: boolean;
  };
}

export interface ScopedKey {
  id: number;
  name: string;
  key: string;
  enabled: boolean;
  isExpired: boolean;
  scopedClientId: string | null;
  scopedUserId: string | null;
  dailyLimit: number | null;
  requestsToday: number;
  expiresAt: string | null;
  hasFoundryCredentials: boolean;
  foundryUrl?: string;
  foundryUsername?: string;
  createdAt: string;
}

export interface ConnectedClient {
  id: string;
  clientId: string;
  customName?: string;
}

export interface Player {
  id: string;
  name: string;
  role: number;
  active: boolean;
}

export interface ApiError {
  error: string;
}
