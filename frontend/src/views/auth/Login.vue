<script setup lang="ts">
import { onMounted, reactive, ref, watch } from 'vue';
import { Sprout, Mail, Lock, Eye, EyeOff, ArrowLeft } from 'lucide-vue-next';
import { useRoute, useRouter } from "vue-router";
import { login, register } from "../../services/auth/auth";
import { getCurrentUser } from '../../services/user/user';
import { touchSession } from '../../utils/session';
import vaniliIllustration from '../../assets/vanili-ilustration.png';

const router = useRouter();
const route = useRoute();

const props = withDefaults(
  defineProps<{
    initialRole?: 'investor' | 'mitra' | 'admin';
  }>(),
  { initialRole: 'investor' }
);

const emit = defineEmits<{
  login: [role: 'investor' | 'mitra' | 'admin', credentials: { email: string; password: string }];
  back: [];
}>();

const role = ref<'investor' | 'mitra' | 'admin'>(props.initialRole);
const showPassword = ref(false);
const isRegistering = ref(false);
const formData = reactive({
  email: '',
  password: '',
  confirmPassword: '',
  fullName: '',
});
const errors = reactive({
  email: '',
  password: '',
  confirmPassword: '',
  fullName: '',
});

function extractRoleFromToken(token: string): 'investor' | 'mitra' | 'admin' {
  try {
    const payload = JSON.parse(atob(token.split('.')[1]));
    if (payload?.role === 'admin') return 'admin';
    if (payload?.role === 'mitra') return 'mitra';
    return 'investor';
  } catch {
    return 'investor';
  }
}

function validateEmail(email: string) {
  const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return re.test(email);
}

function resetErrors() {
  errors.email = '';
  errors.password = '';
  errors.confirmPassword = '';
  errors.fullName = '';
}

function syncRegisterModeFromQuery() {
  const mode = String(route.query.mode || '').toLowerCase();
  const shouldRegister = mode === 'register' || mode === 'daftar' || mode === '1';
  if (shouldRegister) {
    isRegistering.value = true;
  }
}

async function handleSubmit(e: Event) {
  e.preventDefault();
  resetErrors();

  let hasError = false;
  const newErrors = {
    email: '',
    password: '',
    confirmPassword: '',
    fullName: '',
  };

  if (!formData.email) {
    newErrors.email = 'Email harus diisi';
    hasError = true;
  } else if (!validateEmail(formData.email)) {
    newErrors.email = 'Format email tidak valid';
    hasError = true;
  }

  if (!formData.password) {
    newErrors.password = 'Password harus diisi';
    hasError = true;
  } else if (formData.password.length < 6) {
    newErrors.password = 'Password minimal 6 karakter';
    hasError = true;
  }

  if (isRegistering.value) {
    if (!formData.fullName) {
      newErrors.fullName = 'Nama lengkap harus diisi';
      hasError = true;
    }

    if (!formData.confirmPassword) {
      newErrors.confirmPassword = 'Konfirmasi password harus diisi';
      hasError = true;
    } else if (formData.password !== formData.confirmPassword) {
      newErrors.confirmPassword = 'Password tidak cocok';
      hasError = true;
    }
  }

  if (hasError) {
    Object.assign(errors, newErrors);
    return;
  }

  try {
    if (isRegistering.value) {
      await handleRegister();
    } else {
      await handleLogin();
    }
  } catch (err) {
    console.error("AUTH ERROR:", err);
    if (isRegistering.value) {
      errors.email = "Email sudah terdaftar atau terjadi kesalahan";
    } else {
      errors.email = "Email atau password salah";
    }
  }
}

async function handleLogin() {
  console.log("LOGIN REQUEST...");

  const res = await login({
    email: formData.email,
    password: formData.password,
  });

  console.log("LOGIN SUCCESS:", res);

  const token = res.data.token;
  const detectedRole = extractRoleFromToken(token);

  localStorage.setItem("token", token);
  localStorage.setItem("userRole", detectedRole);
  touchSession();
  role.value = detectedRole;

  const selectedPackageId = String(route.query.package_id || '').trim();
  if (selectedPackageId) {
    router.push({
      path: '/pilih-paket',
      query: { package_id: selectedPackageId },
    });
    return;
  }

  if (detectedRole === 'investor' || detectedRole === 'mitra') {
    try {
      const userRes = await getCurrentUser();
      const hasSelectedPackage = Boolean(userRes?.data?.investments?.length || userRes?.data?.selected_package_id);
      if (!hasSelectedPackage) {
        router.push('/pilih-paket');
        return;
      }
    } catch (err) {
      console.error('FAILED TO VERIFY PACKAGE STATUS:', err);
      router.push('/pilih-paket');
      return;
    }
  }

  router.push("/dashboard");
}

async function handleRegister() {
  console.log("REGISTER REQUEST...");

  const res = await register({
    name: formData.fullName,
    email: formData.email,
    password: formData.password,
    role: role.value,
  });

  console.log("REGISTER SUCCESS:", res);

  // Auto login setelah register
  const loginRes = await login({
    email: formData.email,
    password: formData.password,
  });

  console.log("AUTO LOGIN SUCCESS:", loginRes);

  const token = loginRes.data.token;
  const detectedRole = extractRoleFromToken(token);

  localStorage.setItem("token", token);
  localStorage.setItem("userRole", detectedRole);
  touchSession();
  role.value = detectedRole;

  const selectedPackageId = String(route.query.package_id || '').trim();
  if (selectedPackageId) {
    router.push({
      path: '/pilih-paket',
      query: { package_id: selectedPackageId },
    });
    return;
  }

  if (detectedRole === 'investor' || detectedRole === 'mitra') {
    try {
      const userRes = await getCurrentUser();
      const hasSelectedPackage = Boolean(userRes?.data?.investments?.length || userRes?.data?.selected_package_id);
      if (!hasSelectedPackage) {
        router.push('/pilih-paket');
        return;
      }
    } catch (err) {
      console.error('FAILED TO VERIFY PACKAGE STATUS:', err);
      router.push('/pilih-paket');
      return;
    }
  }

  router.push("/dashboard");
}

function toggleRegistering() {
  isRegistering.value = !isRegistering.value;
  resetErrors();
  role.value = 'investor';
  // Reset form data saat toggle
  formData.email = '';
  formData.password = '';
  formData.confirmPassword = '';
  formData.fullName = '';
}

onMounted(() => {
  syncRegisterModeFromQuery();
});

watch(
  () => route.query.mode,
  () => {
    syncRegisterModeFromQuery();
  }
);
</script>

<template>
  <div class="min-h-screen bg-gradient-to-br from-green-50 via-white to-green-100 flex items-center justify-center px-4 py-10">
    <div class="w-full max-w-6xl">
      <div class="rounded-[32px] border border-green-100 bg-white/90 shadow-2xl backdrop-blur overflow-hidden">
        <div class="grid grid-cols-1 lg:grid-cols-[1.05fr_1fr]">
          <!-- Left Panel - Illustration -->
          <div class="hidden lg:flex flex-col gap-6 p-8 bg-gradient-to-br from-emerald-100 via-white to-green-50 relative overflow-hidden">
            <div class="absolute left-4 top-8 h-24 w-24 rounded-full bg-green-200/60 blur-3xl"></div>
            <div class="absolute right-8 top-24 h-16 w-16 rounded-full bg-emerald-200/50 blur-2xl"></div>
            <div class="absolute -bottom-8 left-12 h-32 w-32 rounded-full bg-green-200/50 blur-3xl"></div>
            <div class="absolute inset-x-6 bottom-10 h-24 rounded-3xl bg-white/40 blur-2xl"></div>
            <div class="flex items-center gap-3 relative z-10">
              <img src="../../assets/logo-pondok-tani.png" alt="Pondok Tani Land" class="w-16 h-16 object-contain rounded-2xl shadow-sm bg-white p-1.5" />
              <div>
                <h1 class="text-3xl font-bold text-green-700">Pondok Tani Land</h1>
              </div>
            </div>

            <div class="relative z-10 flex-1 flex items-center justify-center">
              <img :src="vaniliIllustration" alt="Ilustrasi Vanili" class="max-h-full w-auto max-w-full object-contain" />
            </div>
          </div>

          <!-- Right Panel - Login Form -->
          <div class="bg-white flex flex-col">
            <div class="bg-white px-8 py-7">
              <button
                type="button"
                class="flex items-center gap-2 text-green-600 hover:text-green-700 transition-colors mb-4"
                @click="router.push('/')"
              >
                <ArrowLeft class="w-4 h-4" />
                <span class="text-sm">Kembali</span>
              </button>
              <h3 class="text-2xl font-bold text-slate-900">
                {{ isRegistering ? 'Daftar Akun Baru' : 'Masuk ke Akun' }}
              </h3>
              <p class="text-slate-600 mt-1">
                {{
                  isRegistering
                    ? 'Buat akun untuk memulai investasi vanili'
                    : 'Selamat datang kembali! Silakan login untuk melanjutkan'
                }}
              </p>
            </div>

            <div class="p-8 flex flex-col h-full">
              <div v-if="isRegistering" class="flex gap-2 mb-6 bg-gray-100 p-1 rounded-lg">
                <button
                  type="button"
                  class="flex-1 py-2.5 rounded-lg font-medium transition-all"
                  :class="
                    role === 'investor'
                      ? 'bg-white text-green-700 shadow-sm'
                      : 'text-gray-600 hover:text-gray-900'
                  "
                  @click="role = 'investor'"
                >
                  Investor
                </button>
                <button
                  type="button"
                  class="flex-1 py-2.5 rounded-lg font-medium transition-all"
                  :class="
                    role === 'mitra' ? 'bg-white text-green-700 shadow-sm' : 'text-gray-600 hover:text-gray-900'
                  "
                  @click="role = 'mitra'"
                >
                  Mitra
                </button>
              </div>

              <p v-else class="mb-6 text-sm text-gray-500">
              </p>

              <form class="flex flex-col flex-1 gap-4" @submit="handleSubmit">
                <div class="space-y-4">
                  <div v-if="isRegistering">
                    <label class="block text-sm font-medium text-gray-700 mb-2">
                      Nama Lengkap <span class="text-red-500">*</span>
                    </label>
                    <input
                      v-model="formData.fullName"
                      type="text"
                      class="w-full px-4 py-3 border rounded-lg focus:ring-2 focus:ring-green-500 focus:border-transparent transition-all"
                      :class="errors.fullName ? 'border-red-500' : 'border-gray-300'"
                      placeholder="Masukkan nama lengkap Anda"
                    />
                    <p v-if="errors.fullName" class="mt-1 text-sm text-red-600">{{ errors.fullName }}</p>
                  </div>

                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-2">
                      Email <span class="text-red-500">*</span>
                    </label>
                    <div class="relative">
                      <Mail class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
                      <input
                        v-model="formData.email"
                        type="email"
                        class="w-full pl-10 pr-4 py-3 border rounded-lg focus:ring-2 focus:ring-green-500 focus:border-transparent transition-all"
                        :class="errors.email ? 'border-red-500' : 'border-gray-300'"
                        placeholder="nama@email.com"
                      />
                    </div>
                    <p v-if="errors.email" class="mt-1 text-sm text-red-600">{{ errors.email }}</p>
                  </div>

                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-2">
                      Password <span class="text-red-500">*</span>
                    </label>
                    <div class="relative">
                      <Lock class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
                      <input
                        v-model="formData.password"
                        :type="showPassword ? 'text' : 'password'"
                        class="w-full pl-10 pr-12 py-3 border rounded-lg focus:ring-2 focus:ring-green-500 focus:border-transparent transition-all"
                        :class="errors.password ? 'border-red-500' : 'border-gray-300'"
                        placeholder="Masukkan password"
                      />
                      <button
                        type="button"
                        class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
                        @click="showPassword = !showPassword"
                      >
                        <EyeOff v-if="showPassword" class="w-5 h-5" />
                        <Eye v-else class="w-5 h-5" />
                      </button>
                    </div>
                    <p v-if="errors.password" class="mt-1 text-sm text-red-600">{{ errors.password }}</p>
                  </div>

                  <div v-if="isRegistering">
                    <label class="block text-sm font-medium text-gray-700 mb-2">
                      Konfirmasi Password <span class="text-red-500">*</span>
                    </label>
                    <div class="relative">
                      <Lock class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
                      <input
                        v-model="formData.confirmPassword"
                        :type="showPassword ? 'text' : 'password'"
                        class="w-full pl-10 pr-4 py-3 border rounded-lg focus:ring-2 focus:ring-green-500 focus:border-transparent transition-all"
                        :class="errors.confirmPassword ? 'border-red-500' : 'border-gray-300'"
                        placeholder="Masukkan ulang password"
                      />
                    </div>
                    <p v-if="errors.confirmPassword" class="mt-1 text-sm text-red-600">
                      {{ errors.confirmPassword }}
                    </p>
                  </div>
                </div>

                <div class="flex flex-col gap-4 mt-auto">
                  <button
                    type="submit"
                    class="w-full py-3 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors font-semibold text-lg shadow-lg hover:shadow-xl"
                  >
                    {{ isRegistering ? 'Daftar Sekarang' : 'Masuk' }}
                  </button>

                  <div v-if="!isRegistering" class="flex items-center justify-between rounded-2xl bg-slate-50 border border-slate-200 px-4 py-3">
                    <label class="flex items-center gap-2 cursor-pointer text-sm text-gray-700">
                      <input type="checkbox" class="w-4 h-4 text-green-600 border-gray-300 rounded focus:ring-green-500" />
                      <span>Ingat saya</span>
                    </label>
                    <button type="button" class="text-sm text-green-600 hover:text-green-700 font-medium">
                      Lupa password?
                    </button>
                  </div>

                  <div class="text-center">
                    <p class="text-sm text-gray-600">
                      {{ isRegistering ? 'Sudah punya akun?' : 'Belum punya akun?' }}
                      {{ ' ' }}
                      <button type="button" class="text-green-600 hover:text-green-700 font-semibold" @click="toggleRegistering">
                        {{ isRegistering ? 'Masuk di sini' : 'Daftar di sini' }}
                      </button>
                    </p>
                  </div>
                </div>
              </form>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
