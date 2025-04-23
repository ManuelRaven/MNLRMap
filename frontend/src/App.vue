<script setup lang="ts">
import { onMounted, watch } from "vue";
import { RouterView } from "vue-router";
import { usePb } from "./composeables/usePb";
import { useQuasar } from "quasar";

const pb = usePb();

onMounted(async () => {
  let health = await pb.health.check();
  console.log(health);
});

const $q = useQuasar();

// Restore dark mode from local storage
const darkMode = localStorage.getItem("darkMode");
if (darkMode === "true") {
  $q.dark.set(true);
} else {
  $q.dark.set(false);
}

watch(
  () => $q.dark.isActive,
  (val) => {
    localStorage.setItem("darkMode", val.toString());
  },
  { immediate: true }
);
</script>

<template>
  <RouterView />
</template>

<style scoped></style>
