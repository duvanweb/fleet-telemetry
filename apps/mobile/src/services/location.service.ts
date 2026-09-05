import {PermissionsAndroid, Platform} from 'react-native';
import Geolocation from '@react-native-community/geolocation';
import type {Coordinates} from '../types';

const DISTANCE_FILTER_M = 10;
const MIN_SPEED_KMH = 1;
let watchId: number | null = null;

function haversineMeters(lat1: number, lon1: number, lat2: number, lon2: number): number {
  const R = 6371000;
  const dLat = ((lat2 - lat1) * Math.PI) / 180;
  const dLon = ((lon2 - lon1) * Math.PI) / 180;
  const a =
    Math.sin(dLat / 2) ** 2 +
    Math.cos((lat1 * Math.PI) / 180) *
      Math.cos((lat2 * Math.PI) / 180) *
      Math.sin(dLon / 2) ** 2;
  return R * 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
}

async function requestPermission(): Promise<boolean> {
  if (Platform.OS === 'android') {
    const granted = await PermissionsAndroid.request(
      PermissionsAndroid.PERMISSIONS.ACCESS_FINE_LOCATION,
    );
    return granted === PermissionsAndroid.RESULTS.GRANTED;
  }
  return true;
}

export async function startTracking(
  onPosition: (coords: Coordinates) => void,
): Promise<boolean> {
  const allowed = await requestPermission();
  if (!allowed) {
    return false;
  }

  let lastLat: number | null = null;
  let lastLon: number | null = null;
  let highFrequency = true;

  function setInterval(coords: Coordinates): void {
    const speedKmh = (coords.speed ?? 0) * 3.6;
    const wantsHigh = speedKmh >= MIN_SPEED_KMH;
    if (wantsHigh !== highFrequency) {
      highFrequency = wantsHigh;
      stopTracking();
      startTracking(onPosition);
    }
  }

  watchId = Geolocation.watchPosition(
    pos => {
      const {latitude, longitude, speed, accuracy} = pos.coords;

      if (lastLat !== null && lastLon !== null) {
        const dist = haversineMeters(lastLat, lastLon, latitude, longitude);
        if (dist < DISTANCE_FILTER_M) {
          return;
        }
      }

      lastLat = latitude;
      lastLon = longitude;

      const coords: Coordinates = {latitude, longitude, speed, accuracy};
      setInterval(coords);
      onPosition(coords);
    },
    () => {},
    {
      enableHighAccuracy: true,
      distanceFilter: DISTANCE_FILTER_M,
      interval: highFrequency ? 5000 : 30000,
      fastInterval: 2000,
    },
  );

  return true;
}

export function stopTracking(): void {
  if (watchId !== null) {
    Geolocation.clearWatch(watchId);
    watchId = null;
  }
}
