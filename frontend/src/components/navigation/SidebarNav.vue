<template>
  <q-list padding>
    <template v-for="route in mainRoutes" :key="route.path">
      <q-item
        :to="route.path || '/'"
        clickable
        :exact="true"
        v-ripple
        active-class="active-nav-link"
      >
        <q-item-section avatar>
          <q-icon :name="getIconName(route)" />
        </q-item-section>
        <q-item-section>{{ route.meta?.title || route.name }}</q-item-section>
      </q-item>
    </template>

    <q-separator spaced />

    <q-item clickable v-ripple @click="$q.dark.toggle()">
      <q-item-section avatar>
        <q-icon :name="$q.dark.isActive ? 'dark_mode' : 'light_mode'" />
      </q-item-section>
      <q-item-section>{{
        $q.dark.isActive ? "Light Mode" : "Dark Mode"
      }}</q-item-section>
      <q-item-section avatar>
        <q-toggle
          :model-value="$q.dark.isActive"
          @update:model-value="(e) => $q.dark.set(e)"
          color="primary"
          keep-color
        />
      </q-item-section>
    </q-item>
  </q-list>
</template>

<script setup lang="ts">
import { useQuasar } from "quasar";
import { computed, ref, watch } from "vue";
import { useRouter } from "vue-router";
import type { RouteRecordRaw } from "vue-router";

const router = useRouter();
const $q = useQuasar();
// Helper function to safely get the icon name
const getIconName = (route: RouteRecordRaw): string => {
  return (route.meta?.icon as string) || "navigate_next";
};

// Get main routes from the router (only those in MainLayout)
const mainRoutes = computed(() => {
  // Find the MainLayout route
  const mainLayoutRoute = router.options.routes.find(
    (route) =>
      route.component && route.component.toString().includes("MainLayout")
  );

  // Return only the children routes that should be visible in navigation
  return (
    mainLayoutRoute?.children?.filter((route) => {
      // Skip routes without names or marked as hidden in meta
      return route.name && !route.meta?.hideInNav;
    }) || []
  );
});
</script>

<style scoped>
.active-nav-link {
  background-color: rgba(0, 0, 0, 0.1);
  color: var(--q-primary);
  font-weight: bold;
}
</style>
