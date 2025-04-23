<template>
  <q-layout view="hHh lpR fFf">
    <q-header reveal bordered class="bg-primary text-white">
      <q-toolbar>
        <q-btn dense flat round icon="menu" @click="toggleLeftDrawer" />
        <q-toolbar-title>
          <q-avatar>
            <img src="windows11/StoreLogo.scale-200.png" />
          </q-avatar>
          {{ GlobalAppName }}
        </q-toolbar-title>
        <q-space />

        <div class="q-gutter-sm row items-center no-wrap">
          <q-btn-dropdown
            flat
            round
            :ripple="false"
            class="disable-hover-effect"
          >
            <template v-slot:label>
              <q-avatar size="26px" color="dark">
                <img :src="userAvatarUrl" alt="User Avatar" />
              </q-avatar>
            </template>
            <q-list>
              <q-item v-if="authModel" clickable>
                <q-item-section>
                  <q-item-label>{{ authModel.email }}</q-item-label>
                  <q-item-label caption>{{
                    authModel.name || "User"
                  }}</q-item-label>
                </q-item-section>
              </q-item>

              <q-item clickable @click="navigateToProfile">
                <q-item-section>Profile Settings</q-item-section>
                <q-item-section avatar>
                  <q-icon name="settings" />
                </q-item-section>
              </q-item>

              <q-separator />

              <q-item clickable @click="handleLogout">
                <q-item-section>Logout</q-item-section>
                <q-item-section avatar>
                  <q-icon name="logout" />
                </q-item-section>
              </q-item>
            </q-list>
          </q-btn-dropdown>
        </div>
      </q-toolbar>
    </q-header>
    <q-drawer show-if-above v-model="leftDrawerOpen" side="left" bordered>
      <SidebarNav />
    </q-drawer>
    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script setup lang="ts">
import { ref, computed } from "vue";
import SidebarNav from "@/components/navigation/SidebarNav.vue";
import { useAuth } from "@/composeables/useAuth";
import { useRouter } from "vue-router";
import { usePb } from "@/composeables/usePb";

const pb = usePb();

const leftDrawerOpen = ref(false);
const toggleLeftDrawer = () => {
  leftDrawerOpen.value = !leftDrawerOpen.value;
};

const { authModel, logout } = useAuth();
const router = useRouter();

const userAvatarUrl = computed(() => {
  if (authModel.value?.avatar) {
    return pb.files.getURL(authModel.value, authModel.value.avatar);
  }
  return "https://cdn.quasar.dev/img/boy-avatar.png";
});

const handleLogout = async () => {
  await logout();
  router.push("/login"); // Redirect to login page after logout
};

const navigateToProfile = () => {
  router.push("/profile");
};
</script>

<style>
.disable-hover-effect > .q-focus-helper {
  opacity: 0 !important;
}
</style>
