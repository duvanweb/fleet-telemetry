import {create} from 'zustand';
import type {AlertRecord} from '../types';

interface AlertsState {
  alerts: AlertRecord[];
  addAlert: (alert: AlertRecord) => void;
  clearAlerts: () => void;
}

export const useAlertsStore = create<AlertsState>(set => ({
  alerts: [],

  addAlert: alert => {
    set(state => ({
      alerts: [alert, ...state.alerts].slice(0, 200),
    }));
  },

  clearAlerts: () => {
    set({alerts: []});
  },
}));
