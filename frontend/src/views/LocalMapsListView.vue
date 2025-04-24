<template>
  <div class="maps-container q-pa-md">
    <h1 class="text-h4 q-mb-md">Local Maps</h1>

    <q-card flat bordered>
      <q-card-section class="row items-center">
        <div class="text-h6">Maps</div>
        <q-space />
        <q-input
          v-model="searchQuery"
          dense
          outlined
          placeholder="Search maps..."
          class="q-ml-md"
          clearable
          style="width: 200px"
        >
          <template v-slot:append>
            <q-icon name="search" />
          </template>
        </q-input>
        <q-btn
          flat
          round
          color="primary"
          icon="refresh"
          @click="fetchMaps"
          class="q-ml-sm"
          :loading="isLoading"
        />
      </q-card-section>

      <q-separator />

      <q-card-section v-if="isLoading">
        <div class="text-center q-pa-md">
          <q-spinner color="primary" size="3em" />
          <div class="q-mt-sm">Loading maps...</div>
        </div>
      </q-card-section>

      <q-card-section v-else-if="!isLoading && maps.length === 0">
        <div class="text-center text-grey q-pa-md">No maps found</div>
      </q-card-section>

      <q-table
        v-else
        :rows="filteredMaps"
        :columns="tableColumns"
        row-key="id"
        :pagination="pagination"
        :rows-per-page-options="[10, 20, 50, 0]"
        dense
        flat
      >
        <!-- Status column with badges -->
        <template v-slot:body-cell-status="props">
          <q-td :props="props">
            <q-badge
              :color="getStatusColor(props.value)"
              text-color="white"
              :label="props.value"
            />
          </q-td>
        </template>
        <!-- File size column formatted -->
        <template v-slot:body-cell-size="props">
          <q-td :props="props">
            <span v-if="mapSizes[props.value]">{{
              formatFileSize(mapSizes[props.value])
            }}</span>
            <q-spinner v-else size="xs" color="grey" />
          </q-td>
        </template>

        <!-- Created date column formatted -->
        <template v-slot:body-cell-created="props">
          <q-td :props="props">
            {{ formatDate(props.value) }}
          </q-td>
        </template>
        <template v-slot:body-cell-updated="props">
          <q-td :props="props">
            {{ formatDate(props.value) }}
          </q-td>
        </template>

        <!-- Actions column -->
        <template v-slot:body-cell-actions="props">
          <q-td :props="props" class="q-gutter-sm">
            <q-btn
              size="sm"
              color="info"
              icon="info"
              round
              dense
              @click="showMapInfo(props.row)"
              :disable="props.row.status !== 'completed'"
            >
              <q-tooltip>Map Info</q-tooltip>
            </q-btn>

            <q-btn
              size="sm"
              color="primary"
              icon="refresh"
              round
              dense
              @click="confirmRecreateMap(props.row)"
              :disable="props.row.status === 'processing'"
              :loading="processingMaps[props.row.id]"
            >
              <q-tooltip>Recreate Map</q-tooltip>
            </q-btn>

            <q-btn
              size="sm"
              color="negative"
              icon="delete"
              round
              dense
              @click="confirmDeleteMap(props.row)"
              :disable="props.row.status === 'processing'"
            >
              <q-tooltip>Delete Map</q-tooltip>
            </q-btn>
          </q-td>
        </template>
      </q-table>
    </q-card>

    <!-- Confirm Delete Dialog -->
    <q-dialog v-model="confirmDeleteDialog" persistent>
      <q-card>
        <q-card-section class="row items-center">
          <q-avatar icon="delete" color="negative" text-color="white" />
          <span class="q-ml-sm">Are you sure you want to delete this map?</span>
        </q-card-section>

        <q-card-section v-if="selectedMap">
          <div><strong>Name:</strong> {{ selectedMap.name }}</div>
          <div><strong>Status:</strong> {{ selectedMap.status }}</div>
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="Cancel" color="primary" v-close-popup />
          <q-btn
            flat
            label="Delete"
            color="negative"
            @click="deleteMap"
            :loading="isDeleting"
            v-close-popup
          />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Confirm Recreate Dialog -->
    <q-dialog v-model="confirmRecreateDialog" persistent>
      <q-card>
        <q-card-section class="row items-center">
          <q-avatar icon="refresh" color="primary" text-color="white" />
          <span class="q-ml-sm"
            >Are you sure you want to recreate this map?</span
          >
        </q-card-section>

        <q-card-section v-if="selectedMap">
          <div><strong>Name:</strong> {{ selectedMap.name }}</div>
          <div><strong>Status:</strong> {{ selectedMap.status }}</div>
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="Cancel" color="primary" v-close-popup />
          <q-btn
            flat
            label="Recreate"
            color="primary"
            @click="recreateMapAction"
            :loading="isRecreating"
            v-close-popup
          />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Map Info Dialog -->
    <q-dialog v-model="mapInfoDialog">
      <q-card style="width: 700px; max-width: 80vw">
        <q-card-section class="row items-center">
          <q-avatar icon="map" color="info" text-color="white" />
          <span class="text-h6 q-ml-sm">Map Information</span>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section v-if="selectedMap">
          <div class="q-mb-md">
            <div class="text-subtitle1"><strong>Map Details</strong></div>
            <div class="q-pl-md">
              <div><strong>Name:</strong> {{ selectedMap.name }}</div>
              <div><strong>Status:</strong> {{ selectedMap.status }}</div>
              <div v-if="mapSizes[selectedMap.id]">
                <strong>Size:</strong>
                {{ formatFileSize(mapSizes[selectedMap.id]) }}
              </div>
            </div>
          </div>

          <div class="q-mb-md">
            <div class="text-subtitle1"><strong>PMTiles URL</strong></div>
            <div class="q-pa-sm rounded-borders">
              <code>{{ pmtilesUrl }}</code>
              <q-btn
                size="xs"
                flat
                color="primary"
                icon="content_copy"
                class="q-ml-sm"
                @click="copyToClipboard(pmtilesUrl)"
              >
                <q-tooltip>Copy to clipboard</q-tooltip>
              </q-btn>
            </div>
          </div>

          <div class="q-mb-md">
            <div class="text-subtitle1">
              <strong>Using with Protomaps</strong>
            </div>
            <div class="q-pa-sm rounded-borders text-caption">
              <pre class="code-snippet">
import { Protocol } from "pmtiles";
import maplibregl from "maplibre-gl";
import { layers, namedFlavor } from "@protomaps/basemaps";

// Register pmtiles protocol
const protocol = new Protocol();
maplibregl.addProtocol("pmtiles", protocol.tile);

// Create map with pmtiles source
const map = new maplibregl.Map({
  container: "map",
  style: {
    version: 8,
    glyphs: "https://your-domain.com/basemap/fonts/{fontstack}/{range}.pbf",
    sprite: "https://your-domain.com/basemap/sprites/v4/light",
    sources: {
      protomaps: {
        type: "vector",
        url: "{{ pmtilesUrl }}",
      },
    },
    layers: layers("protomaps", namedFlavor("light")),
  },
  center: [11.0, 49.5], // Default center
  zoom: 12, // Default zoom level
});</pre
              >
              <q-btn
                size="xs"
                flat
                color="primary"
                icon="content_copy"
                class="q-mt-sm"
                @click="copyCodeSnippet"
              >
                Copy code snippet
              </q-btn>
            </div>
          </div>
        </q-card-section>
      </q-card>
    </q-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from "vue";
import { usePb } from "@/composeables/usePb";
import { type MapsRecord, MapsStatusOptions } from "@/types/pocketbase-types";
import { useQuasar } from "quasar";
import type { MapSize } from "@/types/custom-types";

const $q = useQuasar();
const pb = usePb();

// Table state
const maps = ref<MapsRecord[]>([]);
const isLoading = ref(false);
const searchQuery = ref("");
const pagination = ref({
  rowsPerPage: 10,
});

// Action state
const processingMaps = ref<Record<string, boolean>>({});
const confirmDeleteDialog = ref(false);
const confirmRecreateDialog = ref(false);
const mapInfoDialog = ref(false);
const selectedMap = ref<MapsRecord | null>(null);
const isDeleting = ref(false);
const isRecreating = ref(false);
const pmtilesUrl = computed(() => {
  if (!selectedMap.value || !selectedMap.value.name) return "";
  return `${window.location.origin}/maps/serve/${selectedMap.value.name}.pmtiles`;
});

// Table columns
const tableColumns = [
  {
    name: "name",
    label: "Name",
    field: "name",
    sortable: true,
    align: "left" as const,
  },
  {
    name: "status",
    label: "Status",
    field: "status",
    sortable: true,
    align: "left" as const,
  },
  {
    name: "size",
    label: "File Size",
    field: "id",
    sortable: true,
    align: "right" as const,
  },
  {
    name: "created",
    label: "Created",
    field: "created",
    sortable: true,
    align: "left" as const,
  },
  {
    name: "updated",
    label: "Updated",
    field: "updated",
    sortable: true,
    align: "left" as const,
  },
  {
    name: "actions",
    label: "Actions",
    field: "actions",
    align: "center" as const,
    sortable: false,
  },
];

// Filter maps based on search query
const filteredMaps = computed(() => {
  if (!searchQuery.value) return maps.value;

  const query = searchQuery.value.toLowerCase();
  return maps.value.filter((map) => {
    return (
      map.name.toLowerCase().includes(query) ||
      (map.status && map.status.toLowerCase().includes(query))
    );
  });
});

// Fetch all maps
const fetchMaps = async () => {
  isLoading.value = true;
  try {
    const result = await pb.collection("maps").getFullList({
      sort: "-created",
    });
    maps.value = result;

    // Fetch map sizes for all maps
    for (const map of result) {
      if (map.name) {
        mapSizes.value[map.id] = await getMapSize(map.name);
      }
    }
  } catch (error) {
    console.error("Error fetching maps:", error);
    $q.notify({
      color: "negative",
      message: "Failed to load maps",
      icon: "error",
    });
  } finally {
    isLoading.value = false;
  }
};

// Format date for display
const formatDate = (dateString: string) => {
  if (!dateString) return "";
  const date = new Date(dateString);
  return date.toLocaleString();
};

// Format file size for display
const formatFileSize = (bytes: number) => {
  if (bytes === 0) return "0 Bytes";

  const units = ["Bytes", "KB", "MB", "GB", "TB"];
  const i = Math.floor(Math.log(bytes) / Math.log(1024));

  return parseFloat((bytes / Math.pow(1024, i)).toFixed(2)) + " " + units[i];
};

// Get color based on status
const getStatusColor = (status: MapsStatusOptions) => {
  switch (status) {
    case "completed":
      return "positive";
    case "processing":
      return "warning";
    case "pending":
      return "info";
    case "failed":
      return "negative";
    default:
      return "grey";
  }
};

// Delete map confirmation
const confirmDeleteMap = (map: MapsRecord) => {
  selectedMap.value = map;
  confirmDeleteDialog.value = true;
};

// Delete map action
const deleteMap = async () => {
  if (!selectedMap.value) return;

  isDeleting.value = true;
  try {
    await pb.collection("maps").delete(selectedMap.value.id);

    $q.notify({
      color: "positive",
      message: `Map "${selectedMap.value.name}" deleted successfully`,
      icon: "check",
    });

    // Refresh the map list
    await fetchMaps();
  } catch (error) {
    console.error("Error deleting map:", error);
    $q.notify({
      color: "negative",
      message: "Failed to delete map",
      icon: "error",
    });
  } finally {
    isDeleting.value = false;
    selectedMap.value = null;
  }
};

// Recreate map confirmation
const confirmRecreateMap = (map: MapsRecord) => {
  selectedMap.value = map;
  confirmRecreateDialog.value = true;
};
// Get map file size
const mapSizes = ref<Record<string, number>>({});

const getMapSize = async (name: string): Promise<number> => {
  try {
    const size: MapSize = await pb.send("/maps/size/" + name, {
      method: "GET",
    });
    return size.sizeBytes;
  } catch (error) {
    console.error("Error fetching map size:", error);
    return 0;
  }
};
// Recreate map action
const recreateMapAction = async () => {
  if (!selectedMap.value) return;

  isRecreating.value = true;
  processingMaps.value[selectedMap.value.id] = true;

  try {
    await pb.send("/maps/recreate/" + selectedMap.value.id, {
      method: "POST",
    });

    $q.notify({
      spinner: true,
      message: `Map "${selectedMap.value.name}" wird nun neu erstellt, je nach Größe kann das einige Minuten dauern`,
      icon: "refresh",
    });
  } catch (error) {
    console.error("Error recreating map:", error);
    $q.notify({
      color: "negative",
      message: "Failed to recreate map",
      icon: "error",
    });
  } finally {
    isRecreating.value = false;
    processingMaps.value[selectedMap.value?.id || ""] = false;
    selectedMap.value = null;
  }
};

// Show map info dialog
const showMapInfo = (map: MapsRecord) => {
  selectedMap.value = map;
  mapInfoDialog.value = true;
};

// Copy to clipboard function
const copyToClipboard = (text: string) => {
  navigator.clipboard.writeText(text).then(
    () => {
      $q.notify({
        color: "positive",
        message: "Copied to clipboard",
        icon: "content_copy",
        timeout: 2000,
      });
    },
    (err) => {
      console.error("Could not copy text: ", err);
      $q.notify({
        color: "negative",
        message: "Failed to copy to clipboard",
        icon: "error",
      });
    }
  );
};

// Copy code snippet
const copyCodeSnippet = () => {
  if (!selectedMap.value) return;

  const codeTemplate = `import { Protocol } from "pmtiles";
import maplibregl from "maplibre-gl";
import { layers, namedFlavor } from "@protomaps/basemaps";

// Register pmtiles protocol
const protocol = new Protocol();
maplibregl.addProtocol("pmtiles", protocol.tile);

// Create map with pmtiles source
const map = new maplibregl.Map({
  container: "map",
  style: {
    version: 8,
    glyphs: "${window.location.origin}/basemap/fonts/{fontstack}/{range}.pbf",
    sprite: "${window.location.origin}/basemap/sprites/v4/light",
    sources: {
      protomaps: {
        type: "vector",
        url: "${pmtilesUrl.value}",
      },
    },
    layers: layers("protomaps", namedFlavor("light")),
  },
  center: [11.0, 49.5], // Default center
  zoom: 12, // Default zoom level
});`;

  copyToClipboard(codeTemplate);
};

// Initialize the component
onMounted(() => {
  fetchMaps();
  pb.collection("maps").subscribe("*", (e) => {
    if (e.action === "update" || e.action === "delete") {
      fetchMaps();
    }
  });
});

onUnmounted(() => {
  pb.collection("maps").unsubscribe("*");
});
</script>

<style scoped>
.maps-container {
  max-width: 1200px;
  margin: 0 auto;
}

.code-snippet {
  white-space: pre;
  overflow-x: auto;
  margin: 0;
  font-family: monospace;
  font-size: 12px;
  line-height: 1.4;
}
</style>
