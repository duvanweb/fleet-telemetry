import React from 'react';
import {
  FlatList,
  StyleSheet,
  Text,
  TouchableOpacity,
  View,
} from 'react-native';

import {useAlertsStore} from '../store/alerts.store';
import type {AlertRecord} from '../types';

function AlertItem({item}: {item: AlertRecord}) {
  const isOpen = item.status === 'OPEN';
  return (
    <View style={[styles.card, isOpen ? styles.cardOpen : styles.cardResolved]}>
      <View style={styles.cardHeader}>
        <Text style={styles.alertType}>{item.type.replace('_', ' ')}</Text>
        <View style={[styles.badge, isOpen ? styles.badgeOpen : styles.badgeResolved]}>
          <Text style={styles.badgeText}>{item.status}</Text>
        </View>
      </View>
      <Text style={styles.vehicleText}>Vehicle: {item.vehicle_id.slice(0, 12)}</Text>
      <Text style={styles.timeText}>
        {new Date(item.received_at).toLocaleString()}
      </Text>
    </View>
  );
}

export function AlertsScreen() {
  const {alerts, clearAlerts} = useAlertsStore();

  return (
    <View style={styles.container}>
      <View style={styles.header}>
        <Text style={styles.title}>Alert History</Text>
        {alerts.length > 0 && (
          <TouchableOpacity onPress={clearAlerts}>
            <Text style={styles.clearBtn}>Clear</Text>
          </TouchableOpacity>
        )}
      </View>

      {alerts.length === 0 ? (
        <View style={styles.empty}>
          <Text style={styles.emptyText}>No alerts received yet.</Text>
        </View>
      ) : (
        <FlatList
          data={alerts}
          keyExtractor={item => item.id}
          renderItem={({item}) => <AlertItem item={item} />}
          contentContainerStyle={styles.list}
        />
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {flex: 1, backgroundColor: '#f9fafb'},
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: 16,
    borderBottomWidth: 1,
    borderBottomColor: '#e5e7eb',
    backgroundColor: '#fff',
  },
  title: {fontSize: 17, fontWeight: '600', color: '#111827'},
  clearBtn: {color: '#ef4444', fontSize: 14},
  list: {padding: 12},
  empty: {flex: 1, alignItems: 'center', justifyContent: 'center'},
  emptyText: {color: '#9ca3af', fontSize: 14},
  card: {
    borderRadius: 10,
    padding: 14,
    marginBottom: 10,
    borderWidth: 1,
  },
  cardOpen: {backgroundColor: '#fef2f2', borderColor: '#fecaca'},
  cardResolved: {backgroundColor: '#fff', borderColor: '#e5e7eb'},
  cardHeader: {flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center', marginBottom: 6},
  alertType: {fontSize: 14, fontWeight: '600', color: '#111827'},
  badge: {borderRadius: 12, paddingHorizontal: 8, paddingVertical: 2},
  badgeOpen: {backgroundColor: '#fee2e2'},
  badgeResolved: {backgroundColor: '#f0fdf4'},
  badgeText: {fontSize: 11, fontWeight: '600', color: '#374151'},
  vehicleText: {fontSize: 12, color: '#6b7280', marginBottom: 2},
  timeText: {fontSize: 11, color: '#9ca3af'},
});
