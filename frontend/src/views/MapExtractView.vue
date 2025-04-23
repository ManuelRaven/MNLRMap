<template>
  <q-page class="map">
    <l-map
      ref="map"
      v-model:zoom="zoom"
      :center="[47.41322, -1.219482]"
      @ready="onMapReady"
    >
      <l-tile-layer
        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
        layer-type="base"
        name="OpenStreetMap"
      ></l-tile-layer>
    </l-map>
    <div class="map-controls">
      <q-btn color="primary" label="Log Bounding Box" @click="logBoundingBox" />
    </div>
  </q-page>
</template>

<script lang="ts" setup>
import L from "leaflet";
import "leaflet/dist/leaflet.css";
import "leaflet-draw/dist/leaflet.draw.css";
import { LMap, LTileLayer } from "@vue-leaflet/vue-leaflet";
import { ref, onMounted } from "vue";

// Import Leaflet Draw plugin
import "leaflet-scribe";

const map = ref<any>(null);
const zoom = ref(13);
let drawControl: any = null;
let rectangleLayer: any = null;
let drawnItems: L.FeatureGroup;

const onMapReady = () => {
  if (!map.value) return;

  // Initialize the FeatureGroup to store drawn items
  drawnItems = new L.FeatureGroup();
  map.value.leafletObject.addLayer(drawnItems);

  // Initialize the draw control and pass it the FeatureGroup of editable layers
  drawControl = new L.Control.Draw({
    edit: {
      featureGroup: drawnItems,
      edit: false,
      remove: true,
    },
    draw: {
      polyline: false,
      polygon: false,
      circle: false,
      marker: false,
      circlemarker: false,
      rectangle: {
        shapeOptions: {
          color: "#ff0000",
          weight: 2,
        },
      },
    },
  });

  map.value.leafletObject.addControl(drawControl);

  // Handle the created event
  map.value.leafletObject.on(L.Draw.Event.CREATED, (event: any) => {
    const layer = event.layer;
    rectangleLayer = layer;
    drawnItems.addLayer(layer);
  });
};

const logBoundingBox = () => {
  if (rectangleLayer) {
    const bounds = rectangleLayer.getBounds();
    console.log("Rectangle Bounding Box:", {
      southwest: bounds.getSouthWest(),
      northeast: bounds.getNorthEast(),
      bbox: [
        bounds.getSouth(),
        bounds.getWest(),
        bounds.getNorth(),
        bounds.getEast(),
      ],
    });
  } else {
    console.log("No rectangle has been drawn yet");
  }
};
</script>

<style scoped>
.map {
  height: calc(100vh - 51px);
  width: 100%;
  position: relative;
}

.map-controls {
  position: absolute;
  bottom: 20px;
  right: 20px;
  z-index: 1000;
}
</style>
