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
import { usePb } from "@/composeables/usePb";
import type { GeoResponse, MapInfo } from "@/types/custom-types";
import type { MapsRecord } from "@/types/pocketbase-types";
import MaplibreGeocoder, {
  type CarmenGeojsonFeature,
  type MaplibreGeocoderApi,
  type MaplibreGeocoderSuggestionResults,
} from "@maplibre/maplibre-gl-geocoder";
import "@maplibre/maplibre-gl-geocoder/dist/maplibre-gl-geocoder.css";
import { layers, namedFlavor } from "@protomaps/basemaps";
import {
  DPI,
  Format,
  MaplibreExportControl,
  PageOrientation,
  Size,
} from "@watergis/maplibre-gl-export";
import "@watergis/maplibre-gl-export/dist/maplibre-gl-export.css";
import maplibregl from "maplibre-gl";
import "maplibre-gl/dist/maplibre-gl.css";
import { Protocol } from "pmtiles";
import { onMounted, onUnmounted, ref, watch } from "vue";

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

    var Geo: MaplibreGeocoderApi = {
      // required
      forwardGeocode: async (config) => {
        try {
          if (!config.query) {
            return {
              type: "FeatureCollection",
              features: [],
            };
          }

          const response = await fetchGeocode(config.query.toString());
          const features: CarmenGeojsonFeature[] = [];

          for (const item of response) {
            features.push({
              type: "Feature",
              geometry: {
                type: "Point",
                coordinates: [item.longitude, item.latitude],
              },
              properties: {
                label: `${item.street} ${item.house_number}, ${item.city}`,
                place_name: `${item.street} ${item.house_number}, ${item.city}`,
                id: item.id,
                place_id: item.id,
              },
              id: "",
              text: `${item.street} ${item.house_number}, ${item.city}`,
              place_name: `${item.street} ${item.house_number}, ${item.city}`,
              place_type: [],
            });
          }

          return {
            type: "FeatureCollection",
            features: features,
          };
        } catch (error) {
          console.error("Geocoding error:", error);
          return {
            type: "FeatureCollection",
            features: [],
          };
        }
      },

      async getSuggestions(config) {
        try {
          if (!config.query) {
            return {
              suggestions: [],
            };
          }

          const response = await fetchGeocode(config.query.toString());
          const suggestionsResult: MaplibreGeocoderSuggestionResults = {
            suggestions: [],
          };

          for (const item of response) {
            suggestionsResult.suggestions.push({
              text: `${item.street} ${item.house_number}, ${item.city}`,
              placeId: `${item.street} ${item.house_number}, ${item.city}`,
            });
          }

          return suggestionsResult;
        } catch (error) {
          console.error("Geocoding error:", error);
          return {
            suggestions: [],
          };
        }
      },
    };

    // Pass in or define a geocoding API that matches the above
    const geocoder = new MaplibreGeocoder(Geo, { maplibregl: maplibregl });
    map.addControl(geocoder, "bottom-left");

    const exportControl = new MaplibreExportControl({
      PageSize: Size.A3,
      PageOrientation: PageOrientation.Portrait,
      Format: Format.PNG,
      DPI: DPI[96],
      Crosshair: true,
      PrintableArea: true,
      Local: "de",
    });
    map.addControl(exportControl, "top-right");

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

const fetchGeocode = async (search: string): Promise<GeoResponse[]> => {
  try {
    const mapInfoResponse: GeoResponse[] = await pb.send(
      "/geoapi/geocode?q=" + search,
      {
        method: "GET",
      }
    );

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
