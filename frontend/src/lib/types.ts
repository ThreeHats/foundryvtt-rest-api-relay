export interface UserData {
  id: number;
  email: string;
  // The dashboard authenticates via `sessionToken` (Authorization: Bearer).
  // The master API key is shown ONCE at registration / regeneration in a
  // one-time modal and never persisted to the user store. The frontend
  // never has access to the master key after that initial display.
  sessionToken: string;
  sessionExpiresAt?: string;
  emailVerified?: boolean;
  apiKeyRotationRequired?: boolean;
  role?: 'user' | 'admin';
  subscriptionStatus: 'free' | 'active' | 'past_due';
  subscriptionEndsAt?: string;
  requestsThisMonth: number;
  limits?: {
    monthlyLimit: number;
    unlimitedMonthly: boolean;
  };
}

// AuthResponse is the shape returned by /auth/register and /auth/login.
// It includes the session token (for ongoing dashboard auth) and — only on
// register/regenerate — the plaintext apiKey for the one-time display modal.
export interface AuthResponse extends UserData {
  apiKey?: string; // present only on register / regenerate-key responses
}

// --- Admin dashboard types ---

export interface AdminUserView {
  id: number;
  email: string;
  role: 'user' | 'admin';
  disabled: boolean;
  emailVerified: boolean;
  subscriptionStatus: string;
  requestsToday: number;
  requestsThisMonth: number;
  maxHeadlessSessions: number;
  apiKeyRotationRequired: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface AdminApiKeyView {
  id: number;
  userId: number;
  name: string;
  keyPrefix: string;
  enabled: boolean;
  isExpired: boolean;
  scopes: string[];
  scopedClientIds: string[];
  scopedUserId: string;
  dailyLimit: number;
  requestsToday: number;
  hasCredentials: boolean;
  expiresAt?: string;
  createdAt: string;
}

export interface AdminClientInfo {
  clientId: string;
  instanceId?: string;
  lastSeen: number;
  connectedSince: number;
  worldId?: string;
  worldTitle?: string;
  foundryVersion?: string;
  systemId?: string;
  systemTitle?: string;
  customName?: string;
  ipAddress?: string;
}

export interface AdminAuditLogEntry {
  id: number;
  adminUserId: number;
  action: string;
  targetType: string;
  targetId: { String: string; Valid: boolean } | string | null;
  details: { String: string; Valid: boolean } | string | null;
  ipAddress: { String: string; Valid: boolean } | string | null;
  createdAt: string;
}

export interface AdminSystemHealth {
  status: string;
  version: string;
  goroutines: number;
  memoryAllocBytes: number;
  memoryHeapBytes: number;
  memorySysBytes: number;
  gcCount: number;
  wsConnectionCount: number;
  pendingRequestCount?: number;
  headlessSessionCount?: number;
  redisStatus?: string;
}

export interface AdminMetricsOverview {
  requestsPerMinute: number;
  requestsPerHour: number;
  requestsPerDay: number;
  errorsTotal: number;
}

export interface AdminFeatureFlags {
  disable_registration: boolean;
  maintenance_mode: boolean;
}

export interface AdminSession {
  email: string;
  role: 'admin';
  csrfToken: string;
}

export interface ScopedKey {
  id: number;
  name: string;
  key: string;
  enabled: boolean;
  isExpired: boolean;
  scopedClientId?: string | null;
  scopedClientIds?: string | string[] | null;
  scopedUserId?: string | null;
  scopedUserIds?: Record<string, string> | null;
  scopes: string[] | null;
  dailyLimit: number | null;
  requestsToday: number;
  expiresAt: string | null;
  createdAt: string;
}

export interface KnownUser {
  id: number;
  knownClientId: number;
  userId: string;
  name: string;
  role: number;
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

export interface ConnectionToken {
  id: number;
  name: string;
  // clientId links this token to its KnownClient (world). Empty for tokens
  // created before this field was added.
  clientId: string;
  allowedIps: string;
  // source: "dashboard" for normal pair flow, "headless" for relay-managed
  // per-session tokens.
  source: 'dashboard' | 'headless';
  lastUsedAt: string | null;
  createdAt: string;
}

export interface KnownClient {
  id: number;
  clientId: string;
  worldId: string;
  worldTitle: string;
  systemId: string;
  systemTitle: string;
  systemVersion: string;
  foundryVersion: string;
  customName: string;
  lastSeenAt: string | null;
  isOnline: boolean;
  // autoStartOnRemoteRequest, when true, lets the relay spawn a headless
  // session for this clientId in response to an incoming remote-request from
  // a sibling client (when this client is currently offline).
  autoStartOnRemoteRequest: boolean;
  // credentialId is the explicit link to a stored Credential used by the
  // headless auto-start path. null means "use the user's first credential
  // if there's exactly one" (works for single-Foundry-server deployments).
  credentialId: number | null;
  // Cross-world tunneling settings (world-level; all browsers for this world share these).
  // allowedTargetClients lists clientIds this world may invoke remote-request against.
  allowedTargetClients: string[];
  // remoteScopes lists scope strings this world holds for cross-world operations.
  remoteScopes: string[];
  // remoteRequestsPerHour: per-world rate limit for cross-world operations. 0 = unlimited.
  remoteRequestsPerHour: number;
  createdAt: string;
  // tokens lists all ConnectionTokens (browser pairings) for this client.
  // activeTokenId is the database ID of the token whose WS connection is
  // currently live; 0 means no active connection.
  tokens: KnownClientToken[];
  activeTokenId: number;
}

// KnownClientToken is a lightweight view of ConnectionToken returned as part
// of the GET /auth/known-clients response (no sensitive fields).
export interface KnownClientToken {
  id: number;
  name: string;
  clientId: string;
  source: 'dashboard' | 'headless';
  allowedIps: string;
  lastUsedAt: string | null;
  createdAt: string;
}

export interface RemoteRequestLog {
  id: number;
  sourceClientId: string;
  sourceTokenId: number;
  targetClientId: string;
  action: string;
  success: boolean;
  errorMessage?: string;
  sourceIp?: string;
  createdAt: string;
}

export interface Credential {
  id: number;
  name: string;
  foundryUrl: string;
  foundryUsername: string;
  createdAt: string;
  updatedAt: string;
}

export interface PairRequestDetails {
  code: string;
  worldId: string;
  worldTitle: string;
  systemId: string;
  systemTitle: string;
  systemVersion: string;
  foundryVersion: string;
  requestedRemoteScopes: string[];
  requestedTargetClients: string[];
  upgradeOnly: boolean;
  status: string;
  expiresAt: string;
  knownClients: KnownClient[];
}

export interface KeyRequestDetails {
  code: string;
  appName: string;
  appDescription: string;
  appUrl: string;
  requestedScopes: string[];
  requestedClientIds: string[];
  callbackUrl: string;
  suggestedDailyLimit: number;
  suggestedExpiry: string;
  status: string;
  expiresAt: string;
}

export interface NotificationSettings {
  discordWebhookUrl: string;
  notifyEmail: string;
  notifyOnNewClientConnect: boolean;
  notifyOnConnect: boolean;
  notifyOnDisconnect: boolean;
  notifyOnMetadataMismatch: boolean;
  notifyOnSettingsChange: boolean;
  notifyOnExecuteJs: boolean;
  notifyOnMacroExecute: boolean;
  notificationDebounceWindowSecs: number;
  remoteRequestBatchWindowSecs: number;
  logCrossWorldRequests: boolean;
  smtpAvailable: boolean;
}

export interface ActivityEvent {
  id: number;
  type: 'connection' | 'remote_request' | 'module_event';
  eventSubtype?: 'connect' | 'disconnect';
  clientId: string;
  worldTitle?: string;
  targetClientId?: string;
  action?: string;
  success?: boolean;
  errorMessage?: string;
  tokenName?: string;
  ipAddress?: string;
  flagged?: boolean;
  actor?: string;
  description?: string;
  userId?: number;
  createdAt: string;
}

export interface ApiKeyNotificationSettings {
  apiKeyId: number;
  discordWebhookUrl: string;
  notifyEmail: string;
  notifyOnExecuteJs: boolean;
  notifyOnMacroExecute: boolean;
  notifyOnRateLimit: boolean;
  notifyOnError: boolean;
  smtpAvailable: boolean;
}

export interface KnownClientNotificationSettings {
  knownClientId: number;
  discordWebhookUrl: string;
  notifyEmail: string;
  notifyOnConnect: boolean;
  notifyOnDisconnect: boolean;
  notifyOnExecuteJs: boolean;
  notifyOnMacroExecute: boolean;
  smtpAvailable: boolean;
}

export interface ConnectionLog {
  id: number;
  clientId: string;
  tokenName?: string;
  ipAddress: string;
  userAgent?: string;
  worldId?: string;
  worldTitle?: string;
  systemId?: string;
  foundryVersion?: string;
  metadataMatch: boolean;
  flagged: boolean;
  flagReason?: string;
  createdAt: string;
}

export interface AlertConfig {
  discordWebhookUrl: string;
  emailDestination: string;
  subscriptions: Array<{ alertType: string; channel: 'discord' | 'email' }>;
}

export interface AlertEvent {
  type: string;
  severity: 'info' | 'warning' | 'critical';
  message: string;
  details?: Record<string, unknown>;
  timestamp: string;
}
