import { log } from '../utils/logger';
import { checkRedisHealth, isLocalMode } from '../config/redis';
import os from 'os';

// System health check
export interface HealthStatus {
  healthy: boolean;
  services: {
    redis: { healthy: boolean; message?: string; skipped?: boolean };
    system: { 
      healthy: boolean;
      freeMem: number;
      totalMem: number;
      memUsedPercent: number;
      uptime: number;
      cpuLoad: number[];
    };
  };
  timestamp: number;
  instanceId: string;
}

export function getSystemHealth(): HealthStatus {
  // Get Redis health
  const redisHealth = checkRedisHealth();
  
  // Get system metrics
  const freeMem = os.freemem();
  const totalMem = os.totalmem();
  const memUsedPercent = ((totalMem - freeMem) / totalMem) * 100;
  const uptime = os.uptime();
  const cpuLoad = os.loadavg();
  
  // System is healthy if memory usage is under 90%
  const systemHealthy = memUsedPercent < 90;
  
  // Overall health status - skip Redis check in local mode
  const healthy = (isLocalMode || redisHealth.healthy) && systemHealthy;
  
  return {
    healthy,
    services: {
      redis: redisHealth,
      system: {
        healthy: systemHealthy,
        freeMem,
        totalMem,
        memUsedPercent,
        uptime,
        cpuLoad
      }
    },
    timestamp: Date.now(),
    instanceId: process.env.FLY_ALLOC_ID || 'local'
  };
}

export function logSystemHealth() {
  const health = getSystemHealth();
  
  log.info(`System health: ${health.healthy ? 'HEALTHY' : 'UNHEALTHY'}`);
  log.info(`  Memory: ${Math.round(health.services.system.memUsedPercent)}% used (${Math.round(health.services.system.freeMem/1024/1024)}MB free)`);
  log.info(`  CPU Load: ${health.services.system.cpuLoad.map(v => v.toFixed(2)).join(', ')}`);
  if (!isLocalMode) {
    log.info(`  Redis: ${health.services.redis.healthy ? 'CONNECTED' : 'DISCONNECTED'} ${health.services.redis.message || ''}`);
  }
}

export function startHealthMonitoring(intervalMs = 300000) { // Default: 5 minutes
  return setInterval(logSystemHealth, intervalMs);
}