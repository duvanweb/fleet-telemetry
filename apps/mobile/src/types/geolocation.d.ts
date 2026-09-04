declare module '@react-native-community/geolocation' {
  interface GeoCoordinates {
    latitude: number;
    longitude: number;
    altitude: number | null;
    accuracy: number;
    altitudeAccuracy: number | null;
    heading: number | null;
    speed: number | null;
  }

  interface GeoPosition {
    coords: GeoCoordinates;
    timestamp: number;
  }

  interface GeoError {
    code: number;
    message: string;
  }

  interface WatchOptions {
    enableHighAccuracy?: boolean;
    distanceFilter?: number;
    interval?: number;
    fastInterval?: number;
    timeout?: number;
    maximumAge?: number;
  }

  const Geolocation: {
    watchPosition(
      success: (position: GeoPosition) => void,
      error?: (error: GeoError) => void,
      options?: WatchOptions,
    ): number;
    clearWatch(watchId: number): void;
    getCurrentPosition(
      success: (position: GeoPosition) => void,
      error?: (error: GeoError) => void,
      options?: WatchOptions,
    ): void;
  };

  export default Geolocation;
}
