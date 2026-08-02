import { StatusBar } from "expo-status-bar";
import { StyleSheet, View } from "react-native";
import MapView, { Marker } from "react-native-maps";

const INITIAL_REGION = {
  latitude: 37.78825,
  longitude: -122.4324,
  latitudeDelta: 0.0922,
  longitudeDelta: 0.0421,
};

const MARKER_COORDINATE = {
  latitude: 37.78825,
  longitude: -122.4324,
};

export default function HomeScreen() {
  return (
    <View style={styles.container}>
      <MapView style={styles.map} initialRegion={INITIAL_REGION}>
        <Marker
          coordinate={MARKER_COORDINATE}
          title="{{project_name}}"
          description="An interactive map powered by {{project_name}}"
        />
      </MapView>
      <StatusBar style="auto" />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  map: {
    ...StyleSheet.absoluteFillObject,
  },
});
