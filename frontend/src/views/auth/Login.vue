<script setup lang="ts">
import { onMounted, reactive, ref, watch } from 'vue';
import { Sprout, Mail, Lock, Eye, EyeOff, ArrowLeft, ArrowRight, Upload } from 'lucide-vue-next';
import { useRoute, useRouter } from "vue-router";
import { login, register, registerMitra } from "../../services/auth/auth";
import { getCurrentUser } from '../../services/user/user';
import { registerFCMToken } from '../../services/firebase/fcm';
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
const registerStep = ref<'account' | 'mitra-detail'>('account');
const mitraPhotoPreview = ref('');
const formData = reactive({
  email: '',
  password: '',
  confirmPassword: '',
  fullName: '',
});
const mitraForm = reactive({
  batchCode: '',
  plantingDate: new Date().toISOString().slice(0, 10),
  location: '',
  seedCount: '',
  landArea: '',
  healthStatus: 'sehat',
  affectedCount: '0',
  disease: '',
  summary: '',
  photo: null as File | null,
});
const errors = reactive({
  email: '',
  password: '',
  confirmPassword: '',
  fullName: '',
  mitra: '',
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
  errors.mitra = '';
}

function syncRegisterModeFromQuery() {
  const mode = String(route.query.mode || '').toLowerCase();
  const shouldRegister = mode === 'register' || mode === 'daftar' || mode === '1';
  if (shouldRegister) {
    isRegistering.value = true;
  }
  if (String(route.query.role || '').toLowerCase() === 'mitra') {
    role.value = 'mitra';
  }
}

function validateAccountFields() {
  resetErrors();

  let hasError = false;
  if (!formData.email) {
    errors.email = 'Email harus diisi';
    hasError = true;
  } else if (!validateEmail(formData.email)) {
    errors.email = 'Format email tidak valid';
    hasError = true;
  }

  if (!formData.password) {
    errors.password = 'Password harus diisi';
    hasError = true;
  } else if (formData.password.length < 6) {
    errors.password = 'Password minimal 6 karakter';
    hasError = true;
  }

  if (!formData.fullName) {
    errors.fullName = 'Nama lengkap harus diisi';
    hasError = true;
  }

  if (!formData.confirmPassword) {
    errors.confirmPassword = 'Konfirmasi password harus diisi';
    hasError = true;
  } else if (formData.password !== formData.confirmPassword) {
    errors.confirmPassword = 'Password tidak cocok';
    hasError = true;
  }

  return !hasError;
}

function validateMitraFields() {
  const seedCount = Number(mitraForm.seedCount);
  const landArea = Number(mitraForm.landArea);
  const affectedCount = Number(mitraForm.affectedCount || 0);

  if (!mitraForm.batchCode.trim()) {
    errors.mitra = 'Batch code wajib diisi';
    return false;
  }
  if (!mitraForm.plantingDate) {
    errors.mitra = 'Tanggal tanam wajib diisi';
    return false;
  }
  if (!mitraForm.location.trim()) {
    errors.mitra = 'Lokasi penanaman wajib diisi';
    return false;
  }
  if (!Number.isFinite(seedCount) || seedCount <= 0) {
    errors.mitra = 'Jumlah bibit harus lebih dari 0';
    return false;
  }
  if (!Number.isFinite(landArea) || landArea <= 0) {
    errors.mitra = 'Luas lahan harus lebih dari 0';
    return false;
  }
  if (!Number.isFinite(affectedCount) || affectedCount < 0 || affectedCount > seedCount) {
    errors.mitra = 'Tanaman terdampak tidak boleh melebihi jumlah bibit';
    return false;
  }
  if (!mitraForm.photo) {
    errors.mitra = 'Foto batch wajib diunggah';
    return false;
  }

  errors.mitra = '';
  return true;
}

function goToMitraDetail() {
  if (!validateAccountFields()) return;
  registerStep.value = 'mitra-detail';
}

function handleMitraPhotoChange(event: Event) {
  const target = event.target as HTMLInputElement;
  const file = target.files?.[0] ?? null;
  if (mitraPhotoPreview.value) {
    URL.revokeObjectURL(mitraPhotoPreview.value);
  }
  mitraForm.photo = file;
  mitraPhotoPreview.value = file ? URL.createObjectURL(file) : '';
}

async function handleSubmit(e: Event) {
  e.preventDefault();
  resetErrors();

  if (isRegistering.value && role.value === 'mitra' && registerStep.value === 'account') {
    goToMitraDetail();
    return;
  }

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

  if (isRegistering.value && role.value === 'mitra' && !validateMitraFields()) {
    return;
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
  await registerFCMToken();

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
      if (detectedRole === 'mitra' && userRes?.data?.package_status === 'pending_validation') {
        router.push('/dashboard');
        return;
      }
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

  const res = role.value === 'mitra'
    ? await registerMitra(buildMitraRegistrationPayload())
    : await register({
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
  await registerFCMToken();

  const selectedPackageId = String(route.query.package_id || '').trim();
  if (selectedPackageId) {
    router.push({
      path: '/pilih-paket',
      query: { package_id: selectedPackageId },
    });
    return;
  }

  if (detectedRole === 'mitra') {
    router.push('/dashboard');
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

function buildMitraRegistrationPayload() {
  const payload = new FormData();
  payload.append('name', formData.fullName);
  payload.append('email', formData.email);
  payload.append('password', formData.password);
  payload.append('batch_code', mitraForm.batchCode.trim());
  payload.append('planting_date', mitraForm.plantingDate);
  payload.append('location', mitraForm.location.trim());
  payload.append('seed_count', String(Math.floor(Number(mitraForm.seedCount))));
  payload.append('land_area', String(Number(mitraForm.landArea)));
  payload.append('health_status', mitraForm.healthStatus);
  payload.append('affected_count', String(Math.floor(Number(mitraForm.affectedCount || 0))));
  payload.append('disease', mitraForm.disease);
  payload.append('summary', mitraForm.summary.trim());
  if (mitraForm.photo) {
    payload.append('photo', mitraForm.photo);
  }
  return payload;
}

function toggleRegistering() {
  isRegistering.value = !isRegistering.value;
  resetErrors();
  role.value = 'investor';
  registerStep.value = 'account';
  // Reset form data saat toggle
  formData.email = '';
  formData.password = '';
  formData.confirmPassword = '';
  formData.fullName = '';
  mitraForm.batchCode = '';
  mitraForm.plantingDate = new Date().toISOString().slice(0, 10);
  mitraForm.location = '';
  mitraForm.seedCount = '';
  mitraForm.landArea = '';
  mitraForm.healthStatus = 'sehat';
  mitraForm.affectedCount = '0';
  mitraForm.disease = '';
  mitraForm.summary = '';
  mitraForm.photo = null;
  if (mitraPhotoPreview.value) {
    URL.revokeObjectURL(mitraPhotoPreview.value);
    mitraPhotoPreview.value = '';
  }
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
                {{ isRegistering && role === 'mitra' && registerStep === 'mitra-detail' ? 'Data Awal Mitra' : isRegistering ? 'Daftar Akun Baru' : 'Masuk ke Akun' }}
              </h3>
              <p class="text-slate-600 mt-1">
                {{
                  isRegistering && role === 'mitra' && registerStep === 'mitra-detail'
                    ? 'Lengkapi data penanaman dan kondisi tanaman awal'
                    : isRegistering
                    ? 'Buat akun untuk memulai investasi vanili'
                    : 'Selamat datang kembali! Silakan login untuk melanjutkan'
                }}
              </p>
            </div>

            <div class="p-8 flex flex-col h-full">
              <div v-if="isRegistering && registerStep === 'account'" class="flex gap-2 mb-6 bg-gray-100 p-1 rounded-lg">
                <button
                  type="button"
                  class="flex-1 py-2.5 rounded-lg font-medium transition-all"
                  :class="
                    role === 'investor'
                      ? 'bg-white text-green-700 shadow-sm'
                      : 'text-gray-600 hover:text-gray-900'
                  "
                  @click="role = 'investor'; registerStep = 'account'"
                >
                  Investor
                </button>
                <button
                  type="button"
                  class="flex-1 py-2.5 rounded-lg font-medium transition-all"
                  :class="
                    role === 'mitra' ? 'bg-white text-green-700 shadow-sm' : 'text-gray-600 hover:text-gray-900'
                  "
                  @click="role = 'mitra'; registerStep = 'account'"
                >
                  Mitra
                </button>
              </div>

              <button
                v-if="isRegistering && role === 'mitra' && registerStep === 'mitra-detail'"
                type="button"
                class="mb-4 inline-flex items-center gap-2 text-sm font-medium text-green-700 hover:text-green-800"
                @click="registerStep = 'account'"
              >
                <ArrowLeft class="h-4 w-4" />
                Kembali ke data akun
              </button>

              <p v-else class="mb-6 text-sm text-gray-500">
              </p>

              <form class="flex flex-col flex-1 gap-4" @submit="handleSubmit">
                <div class="space-y-4">
                  <template v-if="!isRegistering || role !== 'mitra' || registerStep === 'account'">
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
                  </template>

                  <template v-if="isRegistering && role === 'mitra' && registerStep === 'mitra-detail'">
                    <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
                      <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">Batch Code <span class="text-red-500">*</span></label>
                        <input
                          v-model="mitraForm.batchCode"
                          type="text"
                          class="w-full rounded-lg border border-gray-300 px-4 py-3 focus:border-transparent focus:ring-2 focus:ring-green-500"
                          placeholder="Contoh: MITRA-001"
                        />
                      </div>

                      <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">Tanggal Tanam <span class="text-red-500">*</span></label>
                        <input
                          v-model="mitraForm.plantingDate"
                          type="date"
                          class="w-full rounded-lg border border-gray-300 px-4 py-3 focus:border-transparent focus:ring-2 focus:ring-green-500"
                        />
                      </div>

                      <div class="md:col-span-2">
                        <label class="block text-sm font-medium text-gray-700 mb-2">Lokasi Penanaman <span class="text-red-500">*</span></label>
                        <input
                          v-model="mitraForm.location"
                          type="text"
                          class="w-full rounded-lg border border-gray-300 px-4 py-3 focus:border-transparent focus:ring-2 focus:ring-green-500"
                          placeholder="Contoh: Blok A - Kebun Timur"
                        />
                      </div>

                      <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">Jumlah Bibit <span class="text-red-500">*</span></label>
                        <input
                          v-model="mitraForm.seedCount"
                          type="number"
                          min="1"
                          step="1"
                          class="w-full rounded-lg border border-gray-300 px-4 py-3 focus:border-transparent focus:ring-2 focus:ring-green-500"
                          placeholder="Contoh: 120"
                        />
                      </div>

                      <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">Luas Lahan (m2) <span class="text-red-500">*</span></label>
                        <input
                          v-model="mitraForm.landArea"
                          type="number"
                          min="0"
                          step="0.01"
                          class="w-full rounded-lg border border-gray-300 px-4 py-3 focus:border-transparent focus:ring-2 focus:ring-green-500"
                          placeholder="Contoh: 250"
                        />
                      </div>

                      <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">Status Kesehatan <span class="text-red-500">*</span></label>
                        <select
                          v-model="mitraForm.healthStatus"
                          class="w-full rounded-lg border border-gray-300 px-4 py-3 focus:border-transparent focus:ring-2 focus:ring-green-500"
                        >
                          <option value="sehat">Sehat</option>
                          <option value="sebagian_terdampak">Sebagian terdampak</option>
                          <option value="mati">Mati</option>
                        </select>
                      </div>

                      <div>
                        <label class="block text-sm font-medium text-gray-700 mb-2">Tanaman Terdampak</label>
                        <input
                          v-model="mitraForm.affectedCount"
                          type="number"
                          min="0"
                          step="1"
                          class="w-full rounded-lg border border-gray-300 px-4 py-3 focus:border-transparent focus:ring-2 focus:ring-green-500"
                          placeholder="0"
                        />
                      </div>

                      <div class="md:col-span-2">
                        <label class="block text-sm font-medium text-gray-700 mb-2">Disease</label>
                        <select
                          v-model="mitraForm.disease"
                          class="w-full rounded-lg border border-gray-300 px-4 py-3 focus:border-transparent focus:ring-2 focus:ring-green-500"
                        >
                          <option value="">Tidak ada</option>
                          <option value="batang_busuk">Batang busuk</option>
                          <option value="lainnya">Lainnya</option>
                        </select>
                      </div>

                      <div class="md:col-span-2">
                        <label class="block text-sm font-medium text-gray-700 mb-2">Catatan</label>
                        <textarea
                          v-model="mitraForm.summary"
                          rows="3"
                          class="w-full rounded-lg border border-gray-300 px-4 py-3 focus:border-transparent focus:ring-2 focus:ring-green-500"
                          placeholder="Tambahkan ringkasan kondisi tanaman atau catatan panen awal..."
                        />
                      </div>

                      <div class="md:col-span-2">
                        <label class="block text-sm font-medium text-gray-700 mb-2">Foto Batch <span class="text-red-500">*</span></label>
                        <label class="flex cursor-pointer items-center justify-center gap-2 rounded-lg border border-dashed border-gray-300 px-4 py-4 text-sm font-medium text-gray-600 hover:border-green-400 hover:bg-green-50">
                          <Upload class="h-5 w-5" />
                          <span>{{ mitraForm.photo ? mitraForm.photo.name : 'Unggah foto batch' }}</span>
                          <input type="file" accept="image/png,image/jpeg,image/jpg,image/webp" class="hidden" @change="handleMitraPhotoChange" />
                        </label>
                        <div v-if="mitraPhotoPreview" class="mt-3 overflow-hidden rounded-lg border border-gray-200">
                          <img :src="mitraPhotoPreview" alt="Preview foto batch" class="h-40 w-full object-cover" />
                        </div>
                      </div>
                    </div>

                    <p v-if="errors.mitra" class="text-sm text-red-600">{{ errors.mitra }}</p>
                  </template>
                </div>

                <div class="flex flex-col gap-4 mt-auto">
                  <button
                    type="submit"
                    class="w-full py-3 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors font-semibold text-lg shadow-lg hover:shadow-xl"
                  >
                    <span v-if="isRegistering && role === 'mitra' && registerStep === 'account'" class="inline-flex items-center justify-center gap-2">
                      Lanjutkan
                      <ArrowRight class="h-5 w-5" />
                    </span>
                    <span v-else>{{ isRegistering ? 'Daftar Sekarang' : 'Masuk' }}</span>
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
