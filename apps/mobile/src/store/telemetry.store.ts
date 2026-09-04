import {create} from 'zustand';
import type {QueuedTelemetry} from '../types';

interface TelemetryState {
  queue: QueuedTelemetry[];
  enqueue: (item: Omit<QueuedTelemetry, 'id' | 'synced' | 'queued_at'>) => void;
  markSynced: (ids: string[]) => void;
  getPending: () => QueuedTelemetry[];
  clearSynced: () => void;
}

let _idCounter = 0;
function nextId(): string {
  _idCounter += 1;
  return `${Date.now()}-${_idCounter}`;
}

export const useTelemetryStore = create<TelemetryState>((set, get) => ({
  queue: [],

  enqueue: item => {
    const entry: QueuedTelemetry = {
      ...item,
      id: nextId(),
      synced: false,
      queued_at: new Date().toISOString(),
    };
    set(state => ({queue: [...state.queue.slice(-499), entry]}));
  },

  markSynced: ids => {
    const idSet = new Set(ids);
    set(state => ({
      queue: state.queue.map(q => (idSet.has(q.id) ? {...q, synced: true} : q)),
    }));
  },

  getPending: () => get().queue.filter(q => !q.synced),

  clearSynced: () => {
    set(state => ({queue: state.queue.filter(q => !q.synced)}));
  },
}));
