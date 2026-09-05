import React, {useEffect, useRef, useState} from 'react';
import {
  ActivityIndicator,
  StyleSheet,
  Text,
  TouchableOpacity,
  View,
} from 'react-native';

import {startTracking, stopTracking} from '../services/location.service';
import {enqueuePoint, startSyncWorker, stopSyncWorker} from '../services/telemetry.service';
import {useTelemetryStore} from '../store/telemetry.store';
import type {Coordinates} from '../types';

const VEHICLE_ID = 'demo-vehicle-001';

export function HomeScreen() {
  const [coords, setCoords] = useState<Coordinates | null>(null);
  const [tracking, setTracking] = useState(false);
  const [connected, setConnected] = useState<boolean | null>(null);
  const pendingCount = useTelemetryStore(s => s.getPending().length);
  const trackingRef = useRef(tracking);
  trackingRef.current = tracking;

  useEffect(() => {
    return () => {
      stopTracking();
      stopSyncWorker();
    };
  }, []);

  async function handleStartTracking() {
    const ok = await startTracking(position => {
      setCoords(position);
      enqueuePoint({
        vehicle_id: VEHICLE_ID,
        latitude: position.latitude,
        longitude: position.longitude,
        device_timestamp: new Date().toISOString(),
      });
    });
    if (ok) {
      startSyncWorker();
      setTracking(true);
      setConnected(true);
    } else {
      setConnected(false);
    }
  }

  function handleStopTracking() {
    stopTracking();
    stopSyncWorker();
    setTracking(false);
  }

  function handlePanic() {
    if (coords) {
      enqueuePoint({
        vehicle_id: VEHICLE_ID,
        latitude: coords.latitude,
        longitude: coords.longitude,
        device_timestamp: new Date().toISOString(),
      });
    }
  }

  return (
    <View style={styles.container}>
      <View style={styles.statusRow}>
        <View
          style={[
            styles.dot,
            connected === true
              ? styles.dotGreen
              : connected === false
              ? styles.dotRed
              : styles.dotGray,
          ]}
        />
        <Text style={styles.statusText}>
          {connected === null
            ? 'Not started'
            : connected
            ? 'Connected'
            : 'Permission denied'}
        </Text>
      </View>

      {coords ? (
        <View style={styles.coordBox}>
          <Text style={styles.coordLabel}>Latitude</Text>
          <Text style={styles.coordValue}>{coords.latitude.toFixed(6)}</Text>
          <Text style={styles.coordLabel}>Longitude</Text>
          <Text style={styles.coordValue}>{coords.longitude.toFixed(6)}</Text>
          {coords.speed != null && (
            <>
              <Text style={styles.coordLabel}>Speed</Text>
              <Text style={styles.coordValue}>
                {(coords.speed * 3.6).toFixed(1)} km/h
              </Text>
            </>
          )}
        </View>
      ) : tracking ? (
        <ActivityIndicator style={styles.spinner} size="large" color="#3b82f6" />
      ) : null}

      <Text style={styles.queueText}>Pending: {pendingCount} points</Text>

      {!tracking ? (
        <TouchableOpacity style={styles.btnStart} onPress={handleStartTracking}>
          <Text style={styles.btnText}>Start Tracking</Text>
        </TouchableOpacity>
      ) : (
        <TouchableOpacity style={styles.btnStop} onPress={handleStopTracking}>
          <Text style={styles.btnText}>Stop Tracking</Text>
        </TouchableOpacity>
      )}

      <TouchableOpacity
        style={[styles.btnPanic, !tracking && styles.btnDisabled]}
        onPress={handlePanic}
        disabled={!tracking}>
        <Text style={styles.btnText}>! PANIC</Text>
      </TouchableOpacity>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {flex: 1, padding: 24, backgroundColor: '#f9fafb'},
  statusRow: {flexDirection: 'row', alignItems: 'center', marginBottom: 24},
  dot: {width: 12, height: 12, borderRadius: 6, marginRight: 8},
  dotGreen: {backgroundColor: '#22c55e'},
  dotRed: {backgroundColor: '#ef4444'},
  dotGray: {backgroundColor: '#9ca3af'},
  statusText: {fontSize: 14, color: '#374151'},
  coordBox: {
    backgroundColor: '#fff',
    borderRadius: 12,
    padding: 16,
    marginBottom: 16,
    borderWidth: 1,
    borderColor: '#e5e7eb',
  },
  coordLabel: {fontSize: 11, color: '#9ca3af', marginTop: 8},
  coordValue: {fontSize: 18, fontWeight: '600', color: '#111827', fontVariant: ['tabular-nums']},
  spinner: {marginVertical: 32},
  queueText: {fontSize: 12, color: '#6b7280', marginBottom: 24},
  btnStart: {backgroundColor: '#3b82f6', borderRadius: 10, padding: 14, alignItems: 'center', marginBottom: 12},
  btnStop: {backgroundColor: '#6b7280', borderRadius: 10, padding: 14, alignItems: 'center', marginBottom: 12},
  btnPanic: {backgroundColor: '#ef4444', borderRadius: 10, padding: 14, alignItems: 'center'},
  btnDisabled: {opacity: 0.4},
  btnText: {color: '#fff', fontWeight: '700', fontSize: 16},
});
