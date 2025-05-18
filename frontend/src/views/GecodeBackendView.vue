<template>
  <q-page padding>
    <q-card dark bordered class="bg-grey-9 my-card">
      <q-card-section>
        <div class="text-h6">Geocode Backend</div>
      </q-card-section>

      <q-separator dark inset />

      <q-card-section v-if="backendInfo">
        <div class="text-subtitle2">Backend Information:</div>
        <div class="q-my-md">
          <div><strong>Name:</strong> {{ backendInfo.name }}</div>
          <div>
            <strong>Status: </strong>
            <span
              :class="backendInfo.reachable ? 'text-positive' : 'text-negative'"
            >
              {{ backendInfo.reachable ? "OK" : "Not Reachable" }}
            </span>
          </div>
        </div>
      </q-card-section>

      <q-card-section v-else>
        <q-spinner color="primary" size="3em" />
        <div class="q-mt-sm">Lade Backend-Informationen...</div>
      </q-card-section>

      <q-card-actions align="right">
        <q-btn
          flat
          color="primary"
          label="Aktualisieren"
          @click="getBackendInfo"
        />
      </q-card-actions>
    </q-card>
    <!-- Documentation section -->
    <q-card dark bordered class="bg-grey-9 my-card q-mt-md">
      <q-card-section>
        <div class="text-h6">Documentation</div>
      </q-card-section>

      <q-separator dark inset />

      <q-card-section>
        <div class="text-subtitle2">Geocode Backend Configuration:</div>

        <div class="q-my-md">
          <p>
            The geocoding functionality of this application can operate with two
            different backends:
          </p>

          <ul class="q-mb-md">
            <li>
              <strong>OpenStreetMap (OSM):</strong> Default backend when no
              special configuration is provided.
            </li>
            <li>
              <strong>MNLR Address Server:</strong> Can be used for custom
              address databases.
            </li>
          </ul>

          <div class="text-subtitle2 q-mt-lg">Changing the Configuration:</div>
          <p>
            The geocode backend configuration is managed through an environment
            variable in the backend:
          </p>

          <q-card flat bordered class="bg-grey-8 q-pa-sm q-my-md">
            <code>MNLRADDRESSSERVER=https://my-address-data-url.com</code>
          </q-card>

          <p>When setting this environment variable:</p>
          <ul>
            <li>If empty or not set: OpenStreetMap will be used (default)</li>
            <li>
              If set to a URL: The MNLR Address Server at that URL will be used
            </li>
          </ul>

          <p class="text-caption q-mt-md">
            <strong>Note:</strong> After changing the environment variable, the
            backend server must be restarted.
          </p>
        </div>
      </q-card-section>
    </q-card>
  </q-page>
</template>

<script setup lang="ts">
import { usePb } from "@/composeables/usePb";
import type { BackendInfo } from "@/types/custom-types";
import { onMounted, ref } from "vue";

const backendInfo = ref<BackendInfo | null>(null);

const pb = usePb();

onMounted(() => {
  getBackendInfo();
});

const getBackendInfo = async () => {
  try {
    const response = await pb.send<BackendInfo>("/geoapi/backend", {
      method: "GET",
    });
    backendInfo.value = response;
  } catch (error) {
    console.error("Error fetching backend info:", error);
  }
};
</script>

<style scoped></style>
