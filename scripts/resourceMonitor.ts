import fs from 'fs';
import { execSync } from 'child_process';

const BAR_WIDTH = 24;
const DIVIDER = '═'.repeat(64);

// ─── Types ─────────────────────────────────────────────────────────────────

interface CpuTick {
  idle: number;
  total: number;
}

interface GpuSample {
  name: string;
  utilization: number;
  memUsedMb: number;
  memTotalMb: number;
}

interface Sample {
  cpuPercent: number | null;
  ramUsedBytes: number;
  ramTotalBytes: number;
  gpus: GpuSample[];
}

export interface ResourceReport {
  startTime: number;
  endTime: number;
  intervalMs: number;
  samples: Sample[];
}

export interface ResourceMonitor {
  stop(): ResourceReport;
}

export interface TestStats {
  suitesTotal: number;
  suitesFail: number;
  testsTotal: number;
  testsFail: number;
}

// ─── System readers ────────────────────────────────────────────────────────

function readCpuTick(): CpuTick {
  const stat = fs.readFileSync('/proc/stat', 'utf8');
  const parts = stat.split('\n')[0].trim().split(/\s+/).slice(1).map(Number);
  const idle = parts[3] + (parts[4] ?? 0); // idle + iowait
  const total = parts.reduce((a, b) => a + b, 0);
  return { idle, total };
}

function readRam(): { used: number; total: number } {
  const meminfo = fs.readFileSync('/proc/meminfo', 'utf8');
  const kbOf = (key: string): number => {
    const m = meminfo.match(new RegExp(`^${key}:\\s+(\\d+)`, 'm'));
    return m ? parseInt(m[1]) * 1024 : 0;
  };
  const total = kbOf('MemTotal');
  const available = kbOf('MemAvailable');
  return { used: total - available, total };
}

function readGpu(): GpuSample[] {
  try {
    const out = execSync(
      'nvidia-smi --query-gpu=name,utilization.gpu,memory.used,memory.total --format=csv,noheader,nounits',
      { encoding: 'utf8', timeout: 2000, stdio: 'pipe' }
    ).trim();
    return out.split('\n').map(line => {
      const p = line.split(',').map(s => s.trim());
      return {
        name: p[0],
        utilization: parseInt(p[1]) || 0,
        memUsedMb: parseInt(p[2]) || 0,
        memTotalMb: parseInt(p[3]) || 0,
      };
    }).filter(g => g.memTotalMb > 0);
  } catch {
    return [];
  }
}

// ─── Monitor ───────────────────────────────────────────────────────────────

export function startResourceMonitor(intervalMs = 3000): ResourceMonitor {
  const startTime = Date.now();
  const samples: Sample[] = [];
  let prevCpu: CpuTick | null = null;

  try {
    prevCpu = readCpuTick();
  } catch {
    // /proc/stat not available (non-Linux)
  }

  const timer = setInterval(() => {
    try {
      let cpuPercent: number | null = null;
      const curCpu = readCpuTick();
      if (prevCpu !== null) {
        const idleDelta = curCpu.idle - prevCpu.idle;
        const totalDelta = curCpu.total - prevCpu.total;
        if (totalDelta > 0) cpuPercent = (1 - idleDelta / totalDelta) * 100;
      }
      prevCpu = curCpu;

      const ram = readRam();
      samples.push({
        cpuPercent,
        ramUsedBytes: ram.used,
        ramTotalBytes: ram.total,
        gpus: readGpu(),
      });
    } catch {
      // Ignore sampling errors
    }
  }, intervalMs);

  return {
    stop() {
      clearInterval(timer);
      return { startTime, endTime: Date.now(), intervalMs, samples };
    },
  };
}

// ─── Test stats from junit.xml ─────────────────────────────────────────────

export function readTestStats(junitPath = './test-results/junit.xml'): TestStats | null {
  try {
    const xml = fs.readFileSync(junitPath, 'utf8');

    let suitesTotal = 0;
    let suitesFail = 0;
    const suiteRe = /<testsuite [^>]*>/g;
    let m: RegExpExecArray | null;
    while ((m = suiteRe.exec(xml)) !== null) {
      suitesTotal++;
      const f = parseInt(m[0].match(/failures="(\d+)"/)?.[1] ?? '0');
      const e = parseInt(m[0].match(/errors="(\d+)"/)?.[1] ?? '0');
      if (f > 0 || e > 0) suitesFail++;
    }

    const root = xml.match(/<testsuites[^>]*>/)?.[0] ?? '';
    const testsTotal = parseInt(root.match(/tests="(\d+)"/)?.[1] ?? '0');
    const failures = parseInt(root.match(/failures="(\d+)"/)?.[1] ?? '0');
    const errors = parseInt(root.match(/errors="(\d+)"/)?.[1] ?? '0');

    return { suitesTotal, suitesFail, testsTotal, testsFail: failures + errors };
  } catch {
    return null;
  }
}

// ─── Formatting ────────────────────────────────────────────────────────────

function bar(pct: number): string {
  const filled = Math.round((Math.min(100, Math.max(0, pct)) / 100) * BAR_WIDTH);
  return '█'.repeat(filled) + '░'.repeat(BAR_WIDTH - filled);
}

function fmtPct(n: number): string {
  return `${n.toFixed(1)}%`;
}

function fmtGb(bytes: number): string {
  return `${(bytes / 1024 ** 3).toFixed(1)} GB`;
}

function fmtMbAsGb(mb: number): string {
  return `${(mb / 1024).toFixed(1)} GB`;
}

function fmtDuration(ms: number): string {
  const s = Math.round(ms / 1000);
  if (s < 60) return `${s}s`;
  return `${Math.floor(s / 60)}m ${s % 60}s`;
}

function calcStats(vals: number[]): { peak: number; avg: number; min: number } {
  if (vals.length === 0) return { peak: 0, avg: 0, min: 0 };
  return {
    peak: Math.max(...vals),
    min: Math.min(...vals),
    avg: vals.reduce((a, b) => a + b, 0) / vals.length,
  };
}

function cpuRow(label: string, pct: number, last: boolean): string {
  const prefix = last ? '└─' : '├─';
  return `  ${prefix} ${label.padEnd(8)} ${fmtPct(pct).padStart(6)}  ${bar(pct)}`;
}

function ramRow(label: string, used: number, total: number, last: boolean): string {
  const prefix = last ? '└─' : '├─';
  const pct = (used / total) * 100;
  return `  ${prefix} ${label.padEnd(8)} ${fmtGb(used).padStart(9)}  ${fmtPct(pct).padStart(6)}  ${bar(pct)}`;
}

function vramRow(label: string, usedMb: number, totalMb: number, last: boolean): string {
  const prefix = last ? '└─' : '├─';
  const pct = (usedMb / totalMb) * 100;
  return `  ${prefix} ${label.padEnd(8)} ${fmtMbAsGb(usedMb).padStart(9)}  ${fmtPct(pct).padStart(6)}  ${bar(pct)}`;
}

// ─── Report printer ────────────────────────────────────────────────────────

export function printResourceReport(
  report: ResourceReport,
  testStats?: TestStats | null
): void {
  const { samples, startTime, endTime, intervalMs } = report;

  console.log('\n' + DIVIDER);
  console.log('  Test Run Summary');
  console.log(
    `  Duration: ${fmtDuration(endTime - startTime)}` +
    `  |  Samples: ${samples.length}` +
    `  |  Interval: ${intervalMs / 1000}s`
  );

  if (testStats) {
    const suitesPass = testStats.suitesTotal - testStats.suitesFail;
    const testsPass = testStats.testsTotal - testStats.testsFail;
    const suiteStr = testStats.suitesFail === 0
      ? `${suitesPass} passed, ${testStats.suitesTotal} total`
      : `${suitesPass} passed, ${testStats.suitesFail} failed, ${testStats.suitesTotal} total`;
    const testStr = testStats.testsFail === 0
      ? `${testsPass} passed, ${testStats.testsTotal} total`
      : `${testsPass} passed, ${testStats.testsFail} failed, ${testStats.testsTotal} total`;
    console.log(`  Test Suites: ${suiteStr}`);
    console.log(`  Tests:       ${testStr}`);
  }

  console.log(DIVIDER);

  if (samples.length === 0) {
    console.log('  (No samples collected)');
    console.log('\n' + DIVIDER + '\n');
    return;
  }

  // CPU
  const cpuVals = samples.map(s => s.cpuPercent).filter((v): v is number => v !== null);
  if (cpuVals.length > 0) {
    const { peak, avg, min } = calcStats(cpuVals);
    console.log('\n  CPU');
    console.log(cpuRow('Peak:', peak, false));
    console.log(cpuRow('Average:', avg, false));
    console.log(cpuRow('Min:', min, true));
  }

  // RAM
  const ramTotal = samples[0]?.ramTotalBytes ?? 0;
  const ramVals = samples.map(s => s.ramUsedBytes);
  if (ramVals.length > 0 && ramTotal > 0) {
    const { peak, avg, min } = calcStats(ramVals);
    console.log(`\n  RAM  (Total: ${fmtGb(ramTotal)})`);
    console.log(ramRow('Peak:', peak, ramTotal, false));
    console.log(ramRow('Average:', avg, ramTotal, false));
    console.log(ramRow('Min:', min, ramTotal, true));
  }

  // GPU / VRAM (one section per device)
  const gpuCount = (samples.find(s => s.gpus.length > 0)?.gpus ?? []).length;
  for (let i = 0; i < gpuCount; i++) {
    const gpuSamples = samples.filter(s => s.gpus[i] !== undefined);
    if (gpuSamples.length === 0) continue;

    const name = gpuSamples[0].gpus[i].name;
    const utilVals = gpuSamples.map(s => s.gpus[i].utilization);
    const vramVals = gpuSamples.map(s => s.gpus[i].memUsedMb);
    const vramTotal = gpuSamples[0].gpus[i].memTotalMb;

    if (utilVals.length > 0) {
      const { peak, avg, min } = calcStats(utilVals);
      const label = gpuCount === 1 ? 'GPU' : `GPU ${i}`;
      console.log(`\n  ${label} — ${name}`);
      console.log(cpuRow('Peak:', peak, false));
      console.log(cpuRow('Average:', avg, false));
      console.log(cpuRow('Min:', min, true));
    }

    if (vramVals.length > 0 && vramTotal > 0) {
      const { peak, avg, min } = calcStats(vramVals);
      const label = gpuCount === 1 ? 'VRAM' : `VRAM ${i}`;
      console.log(`\n  ${label}  (Total: ${fmtMbAsGb(vramTotal)})`);
      console.log(vramRow('Peak:', peak, vramTotal, false));
      console.log(vramRow('Average:', avg, vramTotal, false));
      console.log(vramRow('Min:', min, vramTotal, true));
    }
  }

  console.log('\n' + DIVIDER + '\n');
}
