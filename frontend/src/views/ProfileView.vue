<template>
  <div class="profile-container q-pa-md">
    <q-card class="profile-card q-pa-md">
      <q-card-section>
        <div class="text-h5 q-mb-md">Profile Settings</div>

        <!-- Avatar upload/display section -->
        <div class="text-center q-mb-lg">
          <q-avatar size="100px" class="q-mb-sm">
            <img :src="userAvatar" alt="User avatar" />
          </q-avatar>
          <div>
            <q-file
              v-model="avatarFile"
              filled
              accept="image/*"
              label="Update profile picture"
              class="q-mt-sm"
              style="max-width: 300px; margin: 0 auto"
            >
              <template v-slot:prepend>
                <q-icon name="photo" />
              </template>
            </q-file>
          </div>
        </div>

        <!-- User information form -->
        <q-form @submit="updateProfile" class="q-gutter-md">
          <!-- Display name -->
          <q-input
            v-model="form.name"
            label="Full Name"
            filled
            :rules="[
              (val) =>
                val.length <= 100 || 'Name must be 100 characters or less',
            ]"
          />

          <!-- Email -->
          <q-input
            v-model="form.email"
            label="Email"
            filled
            type="email"
            :rules="[
              (val) => !!val || 'Email is required',
              (val) =>
                /^\w+([.-]?\w+)*@\w+([.-]?\w+)*(\.\w{2,})+$/.test(val) ||
                'Invalid email format',
            ]"
          />

          <q-separator class="q-my-lg" />

          <!-- Password change section -->
          <div class="text-subtitle1 q-mb-md">Change Password</div>

          <q-input
            v-model="form.oldPassword"
            label="Current Password"
            filled
            type="password"
          />

          <q-input
            v-model="form.password"
            label="New Password"
            filled
            type="password"
            :rules="[
              (val) =>
                !val ||
                val.length >= 8 ||
                'Password must be at least 8 characters',
            ]"
          />

          <q-input
            v-model="form.passwordConfirm"
            label="Confirm New Password"
            filled
            type="password"
            :rules="[
              (val) =>
                !form.password ||
                val === form.password ||
                'Passwords do not match',
            ]"
          />

          <!-- Submit and reset buttons -->
          <div class="row justify-center q-mt-lg">
            <q-btn
              label="Update Profile"
              type="submit"
              color="primary"
              :loading="loading"
            />
            <q-btn
              label="Reset"
              type="reset"
              flat
              class="q-ml-sm"
              @click="resetForm"
            />
          </div>
        </q-form>
      </q-card-section>
    </q-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { useAuth } from "@/composeables/useAuth";
import { usePb } from "@/composeables/usePb";
import { useQuasar } from "quasar";
import type { AuthStoreUserRecord } from "@/types/custom-types";

const pb = usePb();
const { authModel } = useAuth();
const $q = useQuasar();
const loading = ref(false);

// Form data
const form = ref({
  name: "",
  email: "",
  oldPassword: "",
  password: "",
  passwordConfirm: "",
});

// Avatar handling
const avatarFile = ref<File | null>(null);
const userAvatar = computed(() => {
  if (avatarFile.value) {
    return URL.createObjectURL(avatarFile.value);
  } else if (authModel.value?.avatar) {
    // Get avatar from PocketBase
    if (authModel.value.avatar.includes("http")) {
      // If it's already a full URL
      return authModel.value.avatar;
    } else {
      return pb.files.getURL(authModel.value, authModel.value.avatar);
    }
  }
  // Default avatar
  return "https://cdn.quasar.dev/img/boy-avatar.png";
});

// Initialize form data from current user
onMounted(() => {
  if (authModel.value) {
    form.value.name = authModel.value.name || "";
    form.value.email = authModel.value.email;
  }
});

// Form submission
const updateProfile = async () => {
  try {
    loading.value = true;

    if (!authModel.value) {
      throw new Error("Not logged in");
    }

    // Create form data for update
    const formData = new FormData();

    // Only add fields that have changed
    if (form.value.name !== authModel.value.name) {
      formData.append("name", form.value.name);
    }

    if (form.value.email !== authModel.value.email) {
      formData.append("email", form.value.email);
    }

    // Handle avatar upload if a file was selected
    if (avatarFile.value) {
      formData.append("avatar", avatarFile.value);
    }

    // Handle password change if provided
    if (form.value.password) {
      if (!form.value.oldPassword) {
        throw new Error("Current password is required to change password");
      }

      // Verify old password first
      try {
        await pb
          .collection("users")
          .authWithPassword(authModel.value.email, form.value.oldPassword);
      } catch (error) {
        throw new Error("Current password is incorrect");
      }

      formData.append("password", form.value.password);
      formData.append("passwordConfirm", form.value.passwordConfirm);
    }

    // Only proceed with update if we have changes
    if (formData.entries().next().done) {
      $q.notify({
        message: "No changes to update",
        color: "info",
      });
      loading.value = false;
      return;
    }

    // Make the update request
    await pb.collection("users").update(authModel.value.id, formData);

    // Refresh auth data
    await pb.collection("users").authRefresh();

    // Success notification
    $q.notify({
      message: "Profile updated successfully",
      color: "positive",
    });

    // Reset password fields
    form.value.oldPassword = "";
    form.value.password = "";
    form.value.passwordConfirm = "";
    avatarFile.value = null;
  } catch (error: any) {
    console.error("Failed to update profile:", error);
    $q.notify({
      message: `Failed to update profile: ${error.message || "Unknown error"}`,
      color: "negative",
    });
  } finally {
    loading.value = false;
  }
};

// Reset form to current values
const resetForm = () => {
  if (authModel.value) {
    form.value.name = authModel.value.name || "";
    form.value.email = authModel.value.email;
    form.value.oldPassword = "";
    form.value.password = "";
    form.value.passwordConfirm = "";
    avatarFile.value = null;
  }
};
</script>

<style scoped>
.profile-container {
  max-width: 800px;
  margin: 0 auto;
}
.profile-card {
  width: 100%;
}
</style>
