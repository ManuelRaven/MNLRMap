<template>
  <div class="collections-container">
    <h1 class="text-h4 q-mb-md">Collections Explorer</h1>

    <div class="row q-col-gutter-md">
      <!-- Collections sidebar -->
      <div class="col-12 col-sm-4 col-md-3">
        <q-card flat bordered>
          <q-card-section>
            <div class="text-h6">Collections</div>
          </q-card-section>

          <q-separator />

          <q-list padding>
            <q-item
              v-for="collection in collections"
              :key="collection.name"
              clickable
              v-ripple
              :active="selectedCollection === collection.name"
              @click="selectCollection(collection.name)"
              active-class="bg-primary text-white"
            >
              <q-item-section>
                {{ collection.name }}
              </q-item-section>
              <q-item-section side>
                <q-badge :label="collectionCounts[collection.name] || 0" />
              </q-item-section>
            </q-item>
          </q-list>
        </q-card>
      </div>

      <!-- Records display -->
      <div class="col-12 col-sm-8 col-md-9">
        <q-card flat bordered>
          <q-card-section class="row items-center">
            <div class="text-h6" v-if="selectedCollection">
              {{ selectedCollection }} Records
            </div>
            <div class="text-h6" v-else>Select a collection</div>
            <q-space />
            <q-input
              v-if="selectedCollection"
              v-model="searchQuery"
              dense
              outlined
              placeholder="Search records..."
              class="q-ml-md"
              clearable
              style="width: 200px"
            >
              <template v-slot:append>
                <q-icon name="search" />
              </template>
            </q-input>
            <q-btn
              v-if="selectedCollection"
              flat
              round
              color="primary"
              icon="refresh"
              @click="fetchRecords"
              class="q-ml-sm"
              :loading="isLoading"
            />
          </q-card-section>

          <q-separator />

          <q-card-section
            v-if="selectedCollection && !isLoading && records.length === 0"
          >
            <div class="text-center text-grey q-pa-md">No records found</div>
          </q-card-section>

          <q-card-section v-if="isLoading">
            <div class="text-center q-pa-md">
              <q-spinner color="primary" size="3em" />
              <div class="q-mt-sm">Loading records...</div>
            </div>
          </q-card-section>

          <template
            v-if="selectedCollection && !isLoading && records.length > 0"
          >
            <q-table
              :rows="filteredRecords"
              :columns="tableColumns"
              row-key="id"
              :pagination="pagination"
              :rows-per-page-options="[10, 20, 50, 0]"
              dense
              flat
            >
              <template v-slot:body="props">
                <q-tr :props="props">
                  <q-td
                    v-for="col in props.cols"
                    :key="col.name"
                    :props="props"
                  >
                    <template v-if="col.name === 'actions'">
                      <q-btn
                        flat
                        round
                        size="sm"
                        color="primary"
                        icon="visibility"
                        @click="viewRecord(props.row)"
                      />
                    </template>
                    <template v-else-if="isJsonField(col.value)">
                      <pre class="json-field">{{ formatJson(col.value) }}</pre>
                    </template>
                    <template v-else>
                      {{ formatField(col.value) }}
                    </template>
                  </q-td>
                </q-tr>
              </template>
            </q-table>
          </template>
        </q-card>
      </div>
    </div>

    <!-- Record details dialog -->
    <q-dialog v-model="recordDialog" maximized>
      <q-card>
        <q-card-section class="row items-center">
          <div class="text-h6">Record Details</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-separator />

        <q-card-section class="scroll" style="max-height: 70vh">
          <pre>{{ JSON.stringify(selectedRecord, null, 2) }}</pre>
        </q-card-section>
      </q-card>
    </q-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from "vue";
import { usePb } from "@/composeables/usePb";
import { Collections } from "@/types/pocketbase-types";

// State setup
const pb = usePb();
const records = ref<any[]>([]);
const selectedCollection = ref<string | null>(null);
const collectionCounts = ref<Record<string, number>>({});
const isLoading = ref(false);
const searchQuery = ref("");
const selectedRecord = ref<any>(null);
const recordDialog = ref(false);
const pagination = ref({
  rowsPerPage: 20,
});

// Table columns will be dynamically generated
const tableColumns = ref<any[]>([]);

// Hard-coded collection names
const collections = ref([{ name: "users" }]);

// Initialize collections and counts
const initializeCollections = async () => {
  try {
    // Initialize the counts
    await updateCollectionCounts();

    // If collections are available, select the first one
    if (collections.value.length > 0 && !selectedCollection.value) {
      selectCollection(collections.value[0].name);
    }
  } catch (error) {
    console.error("Error initializing collections:", error);
  }
};

// Update the count for all collections
const updateCollectionCounts = async () => {
  for (const collection of collections.value) {
    try {
      const count = await pb.collection(collection.name).getList(1, 1, {
        skipTotal: false,
        sort: "-created",
      });
      collectionCounts.value[collection.name] = count.totalItems;
    } catch (error) {
      console.error(
        `Error fetching count for collection ${collection.name}:`,
        error
      );
      collectionCounts.value[collection.name] = 0;
    }
  }
};

// Select a collection to view
const selectCollection = async (collectionName: string) => {
  selectedCollection.value = collectionName;
  searchQuery.value = "";
  await fetchRecords();
};

// Fetch records for the selected collection
const fetchRecords = async () => {
  if (!selectedCollection.value) return;

  isLoading.value = true;
  records.value = [];

  try {
    const result = await pb.collection(selectedCollection.value).getFullList({
      sort: "-created",
      expand: "creator",
    });
    records.value = result;

    // Generate table columns based on the first record
    if (records.value.length > 0) {
      const firstRecord = records.value[0];
      const columns = [];

      // Always include id first
      columns.push({
        name: "id",
        label: "ID",
        field: "id",
        sortable: true,
        align: "left",
      });

      // Add other fields, excluding expand data
      Object.keys(firstRecord).forEach((key) => {
        if (key !== "id" && key !== "expand" && !key.startsWith("expand.")) {
          columns.push({
            name: key,
            label: key.charAt(0).toUpperCase() + key.slice(1),
            field: key,
            sortable: true,
            align: "left",
          });
        }
      });

      // Add actions column
      columns.push({
        name: "actions",
        label: "Actions",
        field: "actions",
        align: "center",
      });

      tableColumns.value = columns;
    }
  } catch (error) {
    console.error(
      `Error fetching records for collection ${selectedCollection.value}:`,
      error
    );
  } finally {
    isLoading.value = false;
  }
};

// Filter records based on search query
const filteredRecords = computed(() => {
  if (!searchQuery.value) return records.value;

  const query = searchQuery.value.toLowerCase();
  return records.value.filter((record) => {
    return Object.entries(record).some(([key, value]) => {
      if (key === "expand") return false;
      if (value === null || value === undefined) return false;
      return String(value).toLowerCase().includes(query);
    });
  });
});

// View a specific record's details
const viewRecord = (record: any) => {
  selectedRecord.value = record;
  recordDialog.value = true;
};

// Check if a field value is a JSON object or array
const isJsonField = (value: any) => {
  return (
    value !== null &&
    typeof value === "object" &&
    !Array.isArray(value) &&
    Object.keys(value).length > 0
  );
};

// Format JSON fields for display
const formatJson = (value: any) => {
  if (!value) return "";
  try {
    return JSON.stringify(value, null, 2);
  } catch (e) {
    return String(value);
  }
};

// Format field values for display
const formatField = (value: any) => {
  if (value === null || value === undefined) return "";
  if (typeof value === "boolean") return value ? "Yes" : "No";
  if (Array.isArray(value)) return value.join(", ");
  if (typeof value === "object") return JSON.stringify(value);
  return String(value);
};

// Initialize collections on component mount
onMounted(() => {
  initializeCollections();
});

// Watch for collection changes
watch(selectedCollection, () => {
  if (selectedCollection.value) {
    fetchRecords();
  }
});
</script>

<style scoped>
.collections-container {
  padding: 20px;
}

.json-field {
  max-width: 300px;
  max-height: 100px;
  overflow: auto;
  margin: 0;
  font-size: 0.8em;
  background-color: rgba(0, 0, 0, 0.02);
  padding: 4px;
  border-radius: 4px;
}

pre {
  white-space: pre-wrap;
  word-wrap: break-word;
}
</style>
