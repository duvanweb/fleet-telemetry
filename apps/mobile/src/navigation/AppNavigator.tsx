import React from 'react';
import {NavigationContainer} from '@react-navigation/native';
import {createBottomTabNavigator} from '@react-navigation/bottom-tabs';
import {Text} from 'react-native';

import {HomeScreen} from '../screens/HomeScreen';
import {AlertsScreen} from '../screens/AlertsScreen';
import {useAlertsStore} from '../store/alerts.store';

export type RootTabParamList = {
  Home: undefined;
  Alerts: undefined;
};

const Tab = createBottomTabNavigator<RootTabParamList>();

function TabIcon({name, focused}: {name: string; focused: boolean}) {
  const icons: Record<string, string> = {Home: '🗺', Alerts: '🔔'};
  return (
    <Text style={{fontSize: 18, opacity: focused ? 1 : 0.5}}>
      {icons[name] ?? '•'}
    </Text>
  );
}

export function AppNavigator() {
  const openAlerts = useAlertsStore(s => s.alerts.filter(a => a.status === 'OPEN').length);

  return (
    <NavigationContainer>
      <Tab.Navigator
        screenOptions={({route}) => ({
          tabBarIcon: ({focused}) => <TabIcon name={route.name} focused={focused} />,
          tabBarActiveTintColor: '#3b82f6',
          tabBarInactiveTintColor: '#9ca3af',
          headerShown: true,
        })}>
        <Tab.Screen name="Home" component={HomeScreen} options={{title: 'Fleet Tracker'}} />
        <Tab.Screen
          name="Alerts"
          component={AlertsScreen}
          options={{
            tabBarBadge: openAlerts > 0 ? openAlerts : undefined,
          }}
        />
      </Tab.Navigator>
    </NavigationContainer>
  );
}
