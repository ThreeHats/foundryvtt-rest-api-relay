/**
 * Admin API client.
 *
 * Wraps each /admin/api/* endpoint in a typed function. All requests use
 * `adminFetch` which handles cookie credentials + CSRF.
 */

import { adminGet, adminMutate } from './adminAuth';
import type {
  AdminUserView,
  AdminApiKeyView,
  AdminClientInfo,
  AdminAuditLogEntry,
  AdminSystemHealth,
  AdminMetricsOverview,
  AdminFeatureFlags,
  AlertConfig,
  AlertEvent,
  ActivityEvent,
} from './types';

// --- Users ---

export interface UsersListResponse {
  users: AdminUserView[];
  total: number;
  offset: number;
  limit: number;
}

export interface UsersListParams {
  offset?: number;
  limit?: number;
  search?: string;
  role?: string;
  disabled?: boolean | null;
  verified?: boolean | null;
  rotation?: boolean | null;
  subscription?: string;
  sort?: string;
  dir?: 'asc' | 'desc';
}

export const adminApi = {
  // Users
  listUsers: (params: UsersListParams = {}) => {
    const q = new URLSearchParams();
    q.set('offset', String(params.offset ?? 0));
    q.set('limit', String(params.limit ?? 50));
    if (params.search) q.set('search', params.search);
    if (params.role) q.set('role', params.role);
    if (params.disabled != null) q.set('disabled', String(params.disabled));
    if (params.verified != null) q.set('verified', String(params.verified));
    if (params.rotation != null) q.set('rotation', String(params.rotation));
    if (params.subscription) q.set('subscription', params.subscription);
    if (params.sort) q.set('sort', params.sort);
    if (params.dir) q.set('dir', params.dir);
    return adminGet<UsersListResponse>(`/admin/api/users?${q}`);
  },
  getUser: (id: number) => adminGet<AdminUserView>(`/admin/api/users/${id}`),
  updateUser: (id: number, body: Partial<AdminUserView>) =>
    adminMutate<AdminUserView>('PATCH', `/admin/api/users/${id}`, body),
  disableUser: (id: number) => adminMutate<{ message: string }>('POST', `/admin/api/users/${id}/disable`),
  enableUser: (id: number) => adminMutate<{ message: string }>('POST', `/admin/api/users/${id}/enable`),
  rotateUserKey: (id: number) =>
    adminMutate<{ message: string; keyPrefix: string }>('POST', `/admin/api/users/${id}/rotate-key`),
  flagUserRotation: (id: number) =>
    adminMutate<{ message: string }>('POST', `/admin/api/users/${id}/force-rotation`),
  sendPasswordReset: (id: number) =>
    adminMutate<{ message: string }>('POST', `/admin/api/users/${id}/send-password-reset`),
  deleteUser: (id: number) => adminMutate<{ message: string }>('DELETE', `/admin/api/users/${id}`),

  // API Keys
  listKeys: (offset = 0, limit = 50) =>
    adminGet<{ keys: AdminApiKeyView[]; total: number; offset: number; limit: number }>(
      `/admin/api/keys?offset=${offset}&limit=${limit}`
    ),
  getKey: (id: number) => adminGet<AdminApiKeyView>(`/admin/api/keys/${id}`),
  revokeKey: (id: number) => adminMutate<{ message: string }>('DELETE', `/admin/api/keys/${id}`),
  patchKey: (id: number, body: { enabled?: boolean; scopes?: string[] }) =>
    adminMutate<AdminApiKeyView>('PATCH', `/admin/api/keys/${id}`, body),

  // Connected clients
  listClients: () =>
    adminGet<{ total: number; clients: AdminClientInfo[] }>('/admin/api/clients'),
  disconnectClient: (id: string) =>
    adminMutate<{ message: string }>('POST', `/admin/api/clients/${id}/disconnect`),

  // Audit logs
  listAuditLogs: (params: { offset?: number; limit?: number; action?: string; targetType?: string; adminUserId?: number }) => {
    const q = new URLSearchParams();
    if (params.offset !== undefined) q.set('offset', String(params.offset));
    if (params.limit !== undefined) q.set('limit', String(params.limit));
    if (params.action) q.set('action', params.action);
    if (params.targetType) q.set('targetType', params.targetType);
    if (params.adminUserId) q.set('adminUserId', String(params.adminUserId));
    return adminGet<{ entries: AdminAuditLogEntry[]; total: number }>(`/admin/api/audit-logs?${q.toString()}`);
  },

  // Headless sessions
  listHeadlessSessions: () =>
    adminGet<{ total: number; sessions: any[] }>('/admin/api/headless-sessions'),
  killHeadlessSession: (id: string) =>
    adminMutate<{ message: string }>('DELETE', `/admin/api/headless-sessions/${id}`),

  // System health
  getHealth: () => adminGet<AdminSystemHealth>('/admin/api/system/health'),

  // Metrics
  getMetricsOverview: () => adminGet<AdminMetricsOverview>('/admin/api/metrics/overview'),
  getMetricsByEndpoint: () =>
    adminGet<{ endpoints: Array<{ path: string; count: number }> }>('/admin/api/metrics/by-endpoint'),
  getTopConsumers: (limit = 10) =>
    adminGet<{ users: Array<{ userId: number; count: number }> }>(`/admin/api/metrics/top-consumers?limit=${limit}`),

  // Operational tools
  getFeatureFlags: () => adminGet<AdminFeatureFlags>('/admin/api/ops/feature-flags'),
  setFeatureFlag: (flag: string, value: boolean) =>
    adminMutate<AdminFeatureFlags>('POST', '/admin/api/ops/feature-flags', { flag, value }),
  forceDisconnect: (id: string) =>
    adminMutate<{ message: string }>('POST', `/admin/api/ops/force-disconnect/${id}`),

  // Subscriptions
  getSubscriptions: () =>
    adminGet<{ statusCounts: Record<string, number>; subscribers: AdminUserView[] }>('/admin/api/subscriptions'),

  // Alert config
  getAlertConfig: () =>
    adminGet<AlertConfig>('/admin/api/alerts/config'),
  saveAlertConfig: (cfg: AlertConfig) =>
    adminMutate<AlertConfig>('PUT', '/admin/api/alerts/config', cfg),
  testAlert: (channel: 'discord' | 'email') =>
    adminMutate<{ message: string }>('POST', '/admin/api/alerts/test', { channel }),
  getRecentAlerts: () =>
    adminGet<{ events: AlertEvent[] }>('/admin/api/alerts/recent'),

  // Activity log
  getActivity: (params: {
    userId?: number;
    type?: string;
    world?: string;
    action?: string;
    success?: boolean | null;
    since?: string;
    until?: string;
    limit?: number;
    offset?: number;
  } = {}) => {
    const q = new URLSearchParams();
    if (params.userId) q.set('userId', String(params.userId));
    if (params.type) q.set('type', params.type);
    if (params.world) q.set('world', params.world);
    if (params.action) q.set('action', params.action);
    if (params.success === true) q.set('success', 'true');
    else if (params.success === false) q.set('success', 'false');
    if (params.since) q.set('since', params.since);
    if (params.until) q.set('until', params.until);
    if (params.limit != null) q.set('limit', String(params.limit));
    if (params.offset != null) q.set('offset', String(params.offset));
    return adminGet<{ events: ActivityEvent[]; total: number; offset: number; limit: number }>(
      `/admin/api/activity?${q}`
    );
  },
};
