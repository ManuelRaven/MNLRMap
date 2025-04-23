<template>
  <div
    class="login-container row justify-center items-center"
    style="height: 100vh"
  >
    <div class="login-bg"></div>
    <q-card class="login-card q-pa-lg">
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
      <q-card-section class="text-center">
        <h4 class="q-my-md">{{ GlobalAppName }}</h4>
        <h5 class="q-my-md">Login</h5>
      </q-card-section>
      <q-card-section>
        <q-form @submit="onSubmit" class="q-gutter-md">
          <!-- Login method selector -->
          <div
            v-if="availableMethods.length > 1"
            class="row q-col-gutter-sm q-mb-md"
          >
            <div class="col-12">
              <q-btn-toggle
                v-model="loginMode"
                toggle-color="primary"
                color="white"
                text-color="primary"
                class="full-width"
                :options="availableMethods"
              />
            </div>
          </div>

          <!-- Password login form -->
          <template v-if="loginMode === 'password' && showPasswordLogin">
            <q-input
              v-model="email"
              label="Email"
              type="email"
              :rules="[(val) => !!val || 'Email is required', isValidEmail]"
              outlined
            >
              <template v-slot:prepend>
                <q-icon name="email" />
              </template>
            </q-input>

            <q-input
              v-model="password"
              label="Password"
              :type="isPwd ? 'password' : 'text'"
              :rules="[(val) => !!val || 'Password is required']"
              outlined
            >
              <template v-slot:prepend>
                <q-icon name="lock" />
              </template>
              <template v-slot:append>
                <q-icon
                  :name="isPwd ? 'visibility_off' : 'visibility'"
                  class="cursor-pointer"
                  @click="isPwd = !isPwd"
                />
              </template>
            </q-input>
          </template>

          <!-- OTP login form -->
          <template v-else-if="loginMode === 'otp' && showOtpOption">
            <q-input
              v-model="email"
              label="Email"
              type="email"
              :rules="[(val) => !!val || 'Email is required', isValidEmail]"
              outlined
            >
              <template v-slot:prepend>
                <q-icon name="email" />
              </template>
              <template v-slot:append>
                <q-btn
                  flat
                  dense
                  label="Send Code"
                  color="primary"
                  @click="requestOTP"
                  :disable="!isValidEmail(email) || loading"
                />
              </template>
            </q-input>

            <q-input
              v-model="otpCode"
              label="One-Time Code"
              type="number"
              :rules="[(val) => !!val || 'Code is required']"
              outlined
            >
              <template v-slot:prepend>
                <q-icon name="key" />
              </template>
            </q-input>
          </template>

          <!-- MFA login form -->
          <template v-else-if="loginMode === 'mfa' && showMfaOption">
            <q-input
              v-model="email"
              label="Email"
              type="email"
              :rules="[(val) => !!val || 'Email is required', isValidEmail]"
              outlined
            >
              <template v-slot:prepend>
                <q-icon name="email" />
              </template>
            </q-input>

            <q-input
              v-model="password"
              label="Password"
              :type="isPwd ? 'password' : 'text'"
              :rules="[(val) => !!val || 'Password is required']"
              outlined
            >
              <template v-slot:prepend>
                <q-icon name="lock" />
              </template>
              <template v-slot:append>
                <q-icon
                  :name="isPwd ? 'visibility_off' : 'visibility'"
                  class="cursor-pointer"
                  @click="isPwd = !isPwd"
                />
              </template>
            </q-input>

            <q-input
              v-model="mfaCode"
              label="MFA Code"
              type="number"
              :rules="[(val) => !!val || 'MFA code is required']"
              outlined
            >
              <template v-slot:prepend>
                <q-icon name="security" />
              </template>
            </q-input>
          </template>

          <div class="q-mt-md">
            <q-btn
              :label="loginButtonLabel"
              type="submit"
              color="primary"
              class="full-width"
              :loading="loading"
            />
          </div>

          <!-- OAuth providers section -->
          <div v-if="authProviders.length > 0" class="q-mt-md">
            <q-separator class="q-my-md">
              <div
                class="absolute-center bg-white dark:bg-dark q-px-sm text-grey-7"
              >
                or login with
              </div>
            </q-separator>

            <div class="row q-col-gutter-sm q-mt-md">
              <div
                v-for="provider in authProviders"
                :key="provider.name"
                class="col-12 col-sm-6"
              >
                <q-btn
                  :label="provider.displayName"
                  :color="getProviderColor(provider.name)"
                  class="full-width"
                  :icon="getProviderIcon(provider.name)"
                  @click="loginWithOAuth(provider)"
                  :disable="loading"
                />
              </div>
            </div>
          </div>
        </q-form>
      </q-card-section>

      <q-card-section v-if="errorMsg" class="text-negative text-center">
        {{ errorMsg }}
      </q-card-section>

      <q-card-section class="text-center q-pt-none">
        <p class="text-grey-7">
          Don't have an account? <a href="#" class="text-primary">Register</a>
        </p>
      </q-card-section>
    </q-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from "vue";
import { useRouter } from "vue-router";
import { useAuth } from "@/composeables/useAuth";
import { useQuasar } from "quasar";
import { usePb } from "@/composeables/usePb";

// Auth methods interfaces
interface AuthMethodsState {
  mfa: {
    enabled: boolean;
    duration: number;
  };
  otp: {
    enabled: boolean;
    duration: number;
  };
  password: {
    enabled: boolean;
    identityFields: string[];
  };
  oauth2: {
    enabled: boolean;
    providers: AuthProvider[];
  };
}

interface AuthProvider {
  name: string;
  displayName: string;
  state: string;
  authURL: string;
  codeVerifier: string;
  codeChallenge: string;
  codeChallengeMethod: string;
}

const $q = useQuasar();
const router = useRouter();
const { login, listAuthOptions } = useAuth();

// Form state
const email = ref("");
const password = ref("");
const otpCode = ref("");
const otpId = ref("");
const mfaCode = ref("");
const isPwd = ref(true);
const loading = ref(false);
const errorMsg = ref("");
const authProviders = ref<AuthProvider[]>([]);
const authMethods = ref<AuthMethodsState | null>(null);

// Login mode
const loginMode = ref<"password" | "otp" | "mfa">("password");

// Computed properties for showing different login methods
const showPasswordLogin = computed(
  () => authMethods.value?.password?.enabled ?? true
);
const showOtpOption = computed(() => authMethods.value?.otp?.enabled ?? false);
const showMfaOption = computed(() => authMethods.value?.mfa?.enabled ?? false);

// Computed property for available methods (only enabled ones)
const availableMethods = computed(() => {
  const methods = [];

  if (showPasswordLogin.value) {
    methods.push({ label: "Password", value: "password" });
  }

  if (showOtpOption.value) {
    methods.push({ label: "One-Time Code", value: "otp" });
  }

  if (showMfaOption.value) {
    methods.push({ label: "MFA", value: "mfa" });
  }

  return methods;
});

// Computed property for login button label
const loginButtonLabel = computed(() => {
  switch (loginMode.value) {
    case "password":
      return "Login with Password";
    case "otp":
      return "Login with One-Time Code";
    case "mfa":
      return "Login with MFA";
    default:
      return "Login";
  }
});

onMounted(async () => {
  try {
    const methods = await listAuthOptions();
    console.log("Auth methods:", methods);
    authMethods.value = methods;

    // Store OAuth providers if enabled
    if (methods?.oauth2?.enabled && methods.oauth2.providers) {
      authProviders.value = methods.oauth2.providers;
    }

    // Set default login mode to first available method
    setTimeout(() => {
      if (availableMethods.value.length > 0) {
        // Check if current login mode is not available
        if (!availableMethods.value.some((m) => m.value === loginMode.value)) {
          loginMode.value = availableMethods.value[0].value as
            | "password"
            | "otp"
            | "mfa";
        }
      }
    }, 0);
  } catch (error) {
    console.error("Failed to fetch auth methods:", error);
  }
});

const isValidEmail = (val: string) => {
  const emailPattern = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
  return emailPattern.test(val) || "Invalid email format";
};

// Request OTP function
const requestOTP = async () => {
  if (!isValidEmail(email.value) || loading.value) {
    return;
  }

  try {
    loading.value = true;
    errorMsg.value = "";
    const pb = usePb();
    const result = await pb.collection("users").requestOTP(email.value);
    otpId.value = result.otpId;

    $q.notify({
      type: "positive",
      message: "One-time code sent to your email",
      position: "top",
    });
  } catch (error: any) {
    console.error("OTP request error:", error);
    errorMsg.value = error.message || "Failed to send one-time code";

    $q.notify({
      type: "negative",
      message: errorMsg.value,
      position: "top",
    });
  } finally {
    loading.value = false;
  }
};

const getProviderColor = (provider: string): string => {
  const colors: Record<string, string> = {
    google: "red",
    facebook: "blue-10",
    github: "grey-9",
    twitter: "light-blue",
    microsoft: "blue",
    spotify: "green",
    discord: "purple",
    // Add more providers as needed
  };

  return colors[provider] || "grey";
};

const getProviderIcon = (provider: string): string => {
  const icons: Record<string, string> = {
    google: "fab fa-google",
    facebook: "fab fa-facebook-f",
    github: "fab fa-github",
    twitter: "fab fa-twitter",
    microsoft: "fab fa-microsoft",
    spotify: "fab fa-spotify",
    discord: "fab fa-discord",
    // Add more providers as needed
  };

  return icons[provider] || "link";
};

const loginWithOAuth = async (provider: any) => {
  try {
    loading.value = true;
    errorMsg.value = "";

    // Open provider authentication in a popup
    window.open(provider.authURL, "oauth", "width=600,height=700");

    // Listen for the auth completion via PocketBase's OAuth flow
    const pb = usePb();
    await pb.collection("users").authWithOAuth2({
      provider: provider.name,
    });

    // After successful OAuth login
    $q.notify({
      type: "positive",
      message: "Login successful",
      position: "top",
    });

    router.push("/");
  } catch (error: any) {
    console.error("OAuth login error:", error);
    errorMsg.value =
      error.message || `Failed to login with ${provider.displayName}`;

    $q.notify({
      type: "negative",
      message: errorMsg.value,
      position: "top",
    });
  } finally {
    loading.value = false;
  }
};

const onSubmit = async () => {
  try {
    loading.value = true;
    errorMsg.value = "";

    const pb = usePb();
    let authData;

    switch (loginMode.value) {
      case "password":
        authData = await pb
          .collection("users")
          .authWithPassword(email.value, password.value);
        break;

      case "otp":
        if (!otpId.value || !otpCode.value) {
          throw new Error("Please request an OTP code and enter it.");
        }
        authData = await pb
          .collection("users")
          .authWithOTP(otpId.value, otpCode.value);
        break;

      case "mfa":
        // First authenticate with password
        const tempAuth = await pb
          .collection("users")
          .authWithPassword(email.value, password.value);
        // Then validate MFA (this is a simplified approach - production should handle this differently)
        // Note: This is a placeholder as PocketBase's MFA typically requires special handling
        authData = tempAuth;
        break;
    }

    console.log("Login successful:", authData);

    $q.notify({
      type: "positive",
      message: "Login successful",
      position: "top",
    });

    // Redirect to home page or dashboard after successful login
    router.push("/");
  } catch (error: any) {
    console.error("Login error:", error);
    errorMsg.value =
      error.message || "Failed to login. Please check your credentials.";

    $q.notify({
      type: "negative",
      message: errorMsg.value,
      position: "top",
    });
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.login-container {
  position: relative;
  overflow: hidden;
}

.login-bg {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-image: url("/bg.webp");
  background-size: cover;
  background-position: center;
  filter: blur(8px);
  z-index: -1;
}

.login-card {
  width: 100%;
  max-width: 800px;
  border-radius: 8px;
  position: relative;
  z-index: 1;
  backdrop-filter: blur(4px);
}

@media (max-width: 599px) {
  .login-card {
    width: 90%;
  }
}
</style>
