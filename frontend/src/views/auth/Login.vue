<script setup lang="ts">
import { onMounted, reactive, ref, watch } from 'vue';
import { Sprout, Mail, Lock, Eye, EyeOff, ArrowLeft } from 'lucide-vue-next';
import { useRoute, useRouter } from "vue-router";
import { login, register } from "../../services/auth/auth";
import { getCurrentUser } from '../../services/user/user';

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
  <div class="min-h-screen bg-gradient-to-br from-green-50 via-white to-green-50 flex items-center justify-center p-4">
    <div class="w-full max-w-6xl grid grid-cols-1 lg:grid-cols-2 gap-8 items-center">
      <!-- Left Side - Branding -->
      <div class="hidden lg:block">
        <div class="space-y-6">
          <div class="flex items-center gap-3">
            <div class="w-16 h-16 bg-green-600 rounded-2xl flex items-center justify-center">
              <Sprout class="w-10 h-10 text-white" />
            </div>
            <div>
              <h1 class="text-3xl font-bold text-gray-900">Omah Vanili</h1>
              <p class="text-gray-600">Sistem Investasi Vanili</p>
            </div>
          </div>

          <div class="space-y-4">
            <h2 class="text-4xl font-bold text-gray-900 leading-tight">
              Investasi Cerdas,<br />
              Keuntungan Nyata
            </h2>
            <p class="text-lg text-gray-600">
              Platform investasi perkebunan vanili terpercaya dengan transparansi penuh dan monitoring real-time.
            </p>
          </div>

          <div class="grid grid-cols-2 gap-4 pt-6">
            <div class="bg-white rounded-xl p-6 border border-gray-200 shadow-sm">
              <p class="text-3xl font-bold text-green-600">1,200+</p>
              <p class="text-sm text-gray-600 mt-1">Investor Aktif</p>
            </div>
            <div class="bg-white rounded-xl p-6 border border-gray-200 shadow-sm">
              <p class="text-3xl font-bold text-green-600">25%+</p>
              <p class="text-sm text-gray-600 mt-1">ROI Rata-rata</p>
            </div>
          </div>

          <div class="bg-gradient-to-r from-green-600 to-green-700 rounded-2xl p-6 text-white">
            <p class="text-sm opacity-90 mb-2">Testimoni Investor</p>
            <p class="text-lg italic mb-3">
              "Platform yang sangat transparan dan mudah digunakan. ROI saya konsisten di atas 20% setiap tahun!"
            </p>
            <p class="font-medium">— Budi Santoso, Investor sejak 2022</p>
          </div>
        </div>
      </div>

      <!-- Right Side - Login Form -->
      <div class="bg-white rounded-2xl shadow-xl border border-gray-200 overflow-hidden">
        <div class="bg-gradient-to-r from-green-600 to-green-700 px-8 py-6">
          <button
            type="button"
            class="flex items-center gap-2 text-green-100 hover:text-white transition-colors mb-4"
            @click="router.push('/')"
          >
            <ArrowLeft class="w-4 h-4" />
            <span class="text-sm">Kembali</span>
          </button>
          <h3 class="text-2xl font-bold text-white">
            {{ isRegistering ? 'Daftar Akun Baru' : 'Masuk ke Akun' }}
          </h3>
          <p class="text-green-100 mt-1">
            {{
              isRegistering
                ? 'Buat akun untuk memulai investasi vanili'
                : 'Selamat datang kembali! Silakan login untuk melanjutkan'
            }}
          </p>
        </div>

        <div class="p-8">
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
            Role akun akan terdeteksi otomatis setelah login.
          </p>

          <form class="space-y-4" @submit="handleSubmit">
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

            <div v-if="!isRegistering" class="flex items-center justify-between">
              <label class="flex items-center gap-2 cursor-pointer">
                <input type="checkbox" class="w-4 h-4 text-green-600 border-gray-300 rounded focus:ring-green-500" />
                <span class="text-sm text-gray-600">Ingat saya</span>
              </label>
              <button type="button" class="text-sm text-green-600 hover:text-green-700 font-medium">
                Lupa password?
              </button>
            </div>

            <button
              type="submit"
              class="w-full py-3 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors font-semibold text-lg shadow-lg hover:shadow-xl"
            >
              {{ isRegistering ? 'Daftar Sekarang' : 'Masuk' }}
            </button>

            <div class="text-center pt-4">
              <p class="text-sm text-gray-600">
                {{ isRegistering ? 'Sudah punya akun?' : 'Belum punya akun?' }}
                {{ ' ' }}
                <button type="button" class="text-green-600 hover:text-green-700 font-semibold" @click="toggleRegistering">
                  {{ isRegistering ? 'Masuk di sini' : 'Daftar di sini' }}
                </button>
              </p>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>
