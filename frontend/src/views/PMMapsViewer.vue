<template>
  <div class="map-container">
    <div class="map-selector">
      <select v-model="selectedMap" :disabled="isLoading">
        <option disabled value="">Select a map</option>
        <option
          v-for="mapInDB in availableMaps"
          :key="mapInDB.id"
          :value="mapInDB"
        >
          {{ mapInDB.name }}
        </option>
      </select>
      <div v-if="isLoading" class="loading-indicator">Loading...</div>
    </div>
    <div id="map" ref="mapContainer"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from "vue";
import maplibregl from "maplibre-gl";
import "maplibre-gl/dist/maplibre-gl.css";
import { Protocol } from "pmtiles";
import { layers, namedFlavor } from "@protomaps/basemaps";
import { usePb } from "@/composeables/usePb";
import type { MapsRecord } from "@/types/pocketbase-types";
import type { MapInfo } from "@/types/custom-types";

const mapContainer = ref<HTMLElement | null>(null);
const availableMaps = ref<MapsRecord[]>([]);
const selectedMap = ref<MapsRecord | "">("");
const isLoading = ref(false);
let map: maplibregl.Map | null = null;
const pb = usePb();

// Fetch available maps from PocketBase
const fetchMaps = async () => {
  isLoading.value = true;
  try {
    const records = await pb.collection("maps").getList(1, 50, {
      sort: "name",
      filter: 'status = "completed"',
    });

    availableMaps.value = records.items;

    // Select the first map by default if available
    if (records.items.length > 0) {
      selectedMap.value = records.items[0];
    }
  } catch (error) {
    console.error("Error fetching maps:", error);
  } finally {
    isLoading.value = false;
  }
};

watch(
  selectedMap,
  (newMap) => {
    if (newMap) {
      DestroyCurrentMap();
      CreateMap(newMap.name);
    }
  },
  { immediate: true }
);

const DestroyCurrentMap = () => {
  if (map) {
    map.remove();
    map = null;
  }
};

const CreateMap = async (name: string) => {
  if (mapContainer.value) {
    // Fetch Map Info to get center and zoom level
    const mapInfo = await fetchMapInfo(name);

    let protocol = new Protocol();
    maplibregl.addProtocol("pmtiles", protocol.tile);
    map = new maplibregl.Map({
      container: mapContainer.value,
      style: {
        version: 8,
        glyphs:
          window.location.origin + "/basemap/fonts/{fontstack}/{range}.pbf",
        sprite: window.location.origin + "/basemap/sprites/v4/light",
        sources: {
          protomaps: {
            type: "vector",
            url: "pmtiles://maps/serve/" + name + ".pmtiles",
          },
        },
        layers: layers("protomaps", namedFlavor("light"), { lang: "de" }),
      },
      center: [mapInfo.center[0], mapInfo.center[1]], // Default center (longitude, latitude)
      zoom: 12, // Default zoom level
    });

    map.addControl(new maplibregl.NavigationControl(), "top-right");
    map.addControl(new maplibregl.ScaleControl(), "bottom-left");

    // Fetch available maps after map initialization
  }
};

onMounted(() => {
  fetchMaps();
});

onUnmounted(() => {
  if (map) {
    map.remove();
  }
});

const fetchMapInfo = async (mapName: string): Promise<MapInfo> => {
  try {
    const mapInfoResponse: MapInfo = await pb.send("/maps/info/" + mapName, {
      method: "GET",
    });

    console.log("Map info:", mapInfoResponse);
    return mapInfoResponse;
  } catch (error) {
    console.error("Error fetching map info:", error);
    throw error;
  }
};
</script>

<style scoped>
.map-container {
  position: relative;
  width: 100%;
  height: calc(100vh - 51px); /* Adjust height as needed */
}

.map-selector {
  position: absolute;
  top: 10px;
  left: 10px;
  z-index: 1000;
  background-color: white;
  padding: 10px;
  border-radius: 5px;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.3);
}

#map {
  position: absolute;
  top: 0;
  bottom: 0;
  width: 100%;
  height: 100%;
}
</style>
