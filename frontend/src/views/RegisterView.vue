<template>
  <div class="row justify-center items-center" style="height: 100vh;">
    <q-card class="register-card q-pa-lg">
      <q-card-section class="text-center">
        <h4 class="q-my-md">Register</h4>
      </q-card-section>

      <q-card-section>
        <q-form @submit="onSubmit" class="q-gutter-md">
          <q-input 
            v-model="email" 
            label="Email" 
            type="email" 
            :rules="[val => !!val || 'Email is required', isValidEmail]"
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
            :rules="[val => !!val || 'Password is required', val => val.length >= 8 || 'Password must be at least 8 characters']"
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
            v-model="passwordConfirm" 
            label="Confirm Password" 
            :type="isPwdConfirm ? 'password' : 'text'"
            :rules="[
              val => !!val || 'Password confirmation is required',
              val => val === password || 'Passwords do not match'
            ]"
            outlined
          >
            <template v-slot:prepend>
              <q-icon name="lock" />
            </template>
            <template v-slot:append>
              <q-icon
                :name="isPwdConfirm ? 'visibility_off' : 'visibility'"
                class="cursor-pointer"
                @click="isPwdConfirm = !isPwdConfirm"
              />
            </template>
          </q-input>

          <div class="full-width q-mt-md">
            <q-btn
              label="Register"
              type="submit"
              color="primary"
              class="full-width"
              :loading="loading"
            />
          </div>
        </q-form>
      </q-card-section>

      <q-card-section v-if="errorMsg" class="text-negative text-center">
        {{ errorMsg }}
      </q-card-section>

      <q-card-section class="text-center q-pt-none">
        <p class="text-grey-7">Already have an account? <router-link to="/login" class="text-primary">Login</router-link></p>
      </q-card-section>
    </q-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuth } from '@/composeables/useAuth';
import { useQuasar } from 'quasar';

const $q = useQuasar();
const router = useRouter();
const { register } = useAuth();

const email = ref('');
const password = ref('');
const passwordConfirm = ref('');
const isPwd = ref(true);
const isPwdConfirm = ref(true);
const loading = ref(false);
const errorMsg = ref('');

const isValidEmail = (val: string) => {
  const emailPattern = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
  return emailPattern.test(val) || 'Invalid email format';
};

const onSubmit = async () => {
  try {
    loading.value = true;
    errorMsg.value = '';
    
    if (password.value !== passwordConfirm.value) {
      errorMsg.value = 'Passwords do not match';
      return;
    }
    
    const userData = await register(email.value, password.value, passwordConfirm.value);
    console.log('Registration successful:', userData);
    
    $q.notify({
      type: 'positive',
      message: 'Registration successful! You can now login.',
      position: 'top'
    });
    
    // Redirect to login page after successful registration
    router.push('/login');
  } catch (error: any) {
    console.error('Registration error:', error);
    errorMsg.value = error.message || 'Failed to register. Please try again.';
    
    $q.notify({
      type: 'negative',
      message: errorMsg.value,
      position: 'top'
    });
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.register-card {
  width: 100%;
  max-width: 400px;
  border-radius: 8px;
}

@media (max-width: 599px) {
  .register-card {
    width: 90%;
  }
}
</style>
