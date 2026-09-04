import type {TelemetryPayload} from '../types';
import {useTelemetryStore} from '../store/telemetry.store';

const TELEMETRY_SERVICE_URL = 'http://10.0.2.2:8082/api/telemetry';
const BATCH_SIZE = 10;
const RETRY_DELAY_MS = 2000;

let syncInterval: ReturnType<typeof setInterval> | null = null;

export async function sendBatch(payloads: TelemetryPayload[]): Promise<boolean> {
  for (const payload of payloads) {
    try {
      const res = await fetch(TELEMETRY_SERVICE_URL, {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify(payload),
      });
      if (!res.ok && res.status < 400) {
        return false;
      }
    } catch {
      return false;
    }
  }
  return true;
}

export async function syncQueue(): Promise<void> {
  const store = useTelemetryStore.getState();
  const pending = store.getPending().slice(0, BATCH_SIZE);
  if (pending.length === 0) {
    return;
  }

  const ids = pending.map(p => p.id);
  const payloads: TelemetryPayload[] = pending.map(p => ({
    vehicle_id: p.vehicle_id,
    latitude: p.latitude,
    longitude: p.longitude,
    device_timestamp: p.device_timestamp,
  }));

  const ok = await sendBatch(payloads);
  if (ok) {
    store.markSynced(ids);
    store.clearSynced();
  }
}

export function startSyncWorker(): void {
  if (syncInterval !== null) {
    return;
  }
  syncInterval = setInterval(() => {
    syncQueue().catch(() => {});
  }, RETRY_DELAY_MS);
}

export function stopSyncWorker(): void {
  if (syncInterval !== null) {
    clearInterval(syncInterval);
    syncInterval = null;
  }
}

export function enqueuePoint(payload: TelemetryPayload): void {
  useTelemetryStore.getState().enqueue(payload);
}
