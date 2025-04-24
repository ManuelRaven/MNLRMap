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
      <q-btn color="primary" label="Extract" @click="openExtractDialog" />
    </div>

    <!-- Extract Dialog -->
    <q-dialog v-model="extractDialog">
      <q-card style="min-width: 400px">
        <q-card-section>
          <div class="text-h6">Create Map Extract</div>
        </q-card-section>

        <q-card-section class="q-pt-none">
          <div v-if="rectangleLayer">
            <p>Bounding Box Coordinates:</p>
            <pre>{{ formatBoundingBox() }}</pre>
          </div>
          <div v-else class="text-negative">
            Please draw a rectangle on the map first.
          </div>

          <q-input
            v-model="extractName"
            label="Extract Name"
            :rules="[(val) => !!val || 'Name is required']"
            autofocus
          />
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="Cancel" color="negative" v-close-popup />
          <q-btn
            flat
            label="Create"
            color="primary"
            @click="executeExtract"
            :disable="!rectangleLayer || !extractName"
            v-close-popup
          />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script lang="ts" setup>
import L from "leaflet";
import "leaflet/dist/leaflet.css";
import "leaflet-scribe/dist/leaflet.draw.css";
import { LMap, LTileLayer } from "@vue-leaflet/vue-leaflet";
import { ref, onMounted } from "vue";

const pb = usePb(); // Assuming you have a usePb() function to get your PocketBase instance

// Import Leaflet Draw plugin
import "leaflet-scribe";
import { usePb } from "@/composeables/usePb";
import { useRouter } from "vue-router";

const map = ref<any>(null);
const zoom = ref(13);
let drawControl: any = null;
let rectangleLayer: any = null;
let drawnItems: L.FeatureGroup;

// Dialog state variables
const extractDialog = ref(false);
const extractName = ref("");

const router = useRouter();

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

// Format the bounding box for display in the dialog
const formatBoundingBox = () => {
  if (!rectangleLayer) return "";

  const bounds = rectangleLayer.getBounds();
  return JSON.stringify(
    {
      southwest: [
        bounds.getSouthWest().lat.toFixed(6),
        bounds.getSouthWest().lng.toFixed(6),
      ],
      northeast: [
        bounds.getNorthEast().lat.toFixed(6),
        bounds.getNorthEast().lng.toFixed(6),
      ],
    },
    null,
    2
  );
};

// Open the extract dialog
const openExtractDialog = () => {
  if (!rectangleLayer) {
    // If no rectangle is drawn, show a notification
    alert("Please draw a rectangle on the map first");
    return;
  }

  extractDialog.value = true;
};

// Execute the extract with the provided name
const executeExtract = async () => {
  if (!rectangleLayer || !extractName.value) return;

  await createMapExtract();
  // Reset the extract name after successful creation
  extractName.value = "";
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

const createMapExtract = async () => {
  if (!map.value || !rectangleLayer) return;

  // Get the bounding box coordinates
  const bounds = rectangleLayer.getBounds();
  const bbox = [
    bounds.getSouthWest().lng,
    bounds.getSouthWest().lat,
    bounds.getNorthEast().lng,
    bounds.getNorthEast().lat,
  ];

  // Call the API to create the map extract
  try {
    let response: Response = await pb.send("/maps/create", {
      body: {
        bbox: bbox.join(","),
        name: extractName.value, // Use the name from the dialog
      },

      method: "POST",
    });
    router.push({
      name: "localmapslist",
    });
  } catch (error) {
    console.error("Error creating map extract:", error);
    alert("Error creating map extract: " + error);
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
