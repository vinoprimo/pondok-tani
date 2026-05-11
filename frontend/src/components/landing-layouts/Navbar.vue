<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { ChevronDown, LogOut, Sprout, UserCircle2 } from 'lucide-vue-next';
import { getCurrentUser } from '../../services/user/user';
import { clearAuthSession } from '../../utils/session';

const router = useRouter();
const route = useRoute();
const isProfileOpen = ref(false);
const isLoadingProfile = ref(false);
const isLoggedIn = ref(false);
const profile = ref({
  name: 'Guest',
  email: '-',
  role: 'guest',
  packageStatus: 'none',
});

const canEnterDashboard = computed(() => {
  if (!isLoggedIn.value) return false;
  if (profile.value.role === 'admin') return true;
  if (profile.value.role === 'investor' || profile.value.role === 'mitra') {
    return profile.value.packageStatus === 'active';
  }
  return false;
});

const navItems = [
  { label: 'Beranda', path: '/' },
  { label: 'Paket Investasi', path: '/paket-investasi' },
  { label: 'Kemitraan', path: '/kemitraan' },
  { label: 'Kontak Kami', path: '/kontak-kami' },
];

const isActive = (path: string) => {
  if (path === '/') return route.path === '/';
  return route.path.startsWith(path);
};

const mobileTitle = computed(() => {
  if (route.path === '/paket-investasi') return 'Paket Investasi';
  if (route.path === '/kemitraan') return 'Kemitraan';
  if (route.path === '/kontak-kami') return 'Kontak Kami';
  return 'Beranda';
});

async function loadProfile() {
  isLoggedIn.value = Boolean(localStorage.getItem('token'));

  if (!isLoggedIn.value) {
    isProfileOpen.value = false;
    profile.value = {
      name: 'Guest',
      email: '-',
      role: 'guest',
      packageStatus: 'none',
    };
    return;
  }

  isLoadingProfile.value = true;
  try {
    const response = await getCurrentUser();
    profile.value = {
      name: response.data?.name || 'User',
      email: response.data?.email || '-',
      role: response.data?.role || 'investor',
      packageStatus: response.data?.package_status || 'none',
    };
  } catch (error) {
    isLoggedIn.value = false;
    profile.value = {
      name: 'Guest',
      email: '-',
      role: 'guest',
      packageStatus: 'none',
    };
  } finally {
    isLoadingProfile.value = false;
  }
}

function toggleProfile() {
  isProfileOpen.value = !isProfileOpen.value;
}

function handleLogout() {
  clearAuthSession();
  isLoggedIn.value = false;
  isProfileOpen.value = false;
  profile.value = {
    name: 'Guest',
    email: '-',
    role: 'guest',
    packageStatus: 'none',
  };
  router.push('/login');
}

onMounted(() => {
  loadProfile();
});

watch(
  () => route.fullPath,
  () => {
    loadProfile();
  }
);
</script>

<template>
  <header class="border-b border-gray-200 bg-white/80 backdrop-blur-sm sticky top-0 z-50">
    <div class="max-w-7xl mx-auto px-6 py-4 flex items-center justify-between gap-4">
      <button type="button" class="flex items-center gap-3" @click="router.push('/')">
        <img src="../../assets/logo-pondok-tani.png" alt="Pondok Tani Land" class="w-10 h-10 object-contain rounded-lg" />
        <div class="text-left">
          <h1 class="font-bold text-green-600">Pondok Tani Land</h1>
        </div>
      </button>
      <nav class="hidden md:flex items-center gap-1">
        <button
          v-for="item in navItems"
          :key="item.path"
          type="button"
          class="px-4 py-2 rounded-lg font-medium transition-colors"
          :class="isActive(item.path) ? 'bg-green-100 text-green-600' : 'text-gray-700 hover:bg-gray-100'"
          @click="router.push(item.path)"
        >
          {{ item.label }}
        </button>
      </nav>

      <div class="hidden md:flex items-center gap-3 relative">
        <template v-if="isLoggedIn && canEnterDashboard">
          <button
            type="button"
            class="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors font-medium"
            @click="router.push('/dashboard')"
          >
            Masuk ke Dashboard
          </button>
        </template>
        <template v-else-if="isLoggedIn">
          <button
            type="button"
            class="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors font-medium"
            @click="router.push('/pilih-paket')"
          >
            Pilih Paket
          </button>
        </template>
        <template v-else>
          <button
            type="button"
            class="px-4 py-2 text-green-600 hover:bg-green-50 rounded-lg transition-colors font-medium"
            @click="router.push('/login')"
          >
            Masuk
          </button>
          <button
            type="button"
            class="px-5 py-2.5 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors font-medium"
            @click="router.push({ path: '/login', query: { mode: 'register' } })"
          >
            Daftar Sekarang
          </button>
        </template>
        <button
          v-if="isLoggedIn"
          type="button"
          class="inline-flex items-center gap-1 rounded-lg border border-gray-200 bg-white px-3 py-2 text-gray-700 hover:bg-gray-50"
          @click="toggleProfile"
        >
          <UserCircle2 class="h-5 w-5" />
          <ChevronDown class="h-4 w-4" />
        </button>

        <div
          v-if="isLoggedIn && isProfileOpen"
          class="absolute right-0 top-12 w-64 rounded-xl border border-gray-200 bg-white p-4 shadow-lg"
        >
          <p class="text-xs uppercase tracking-wide text-gray-500 mb-2">Profil</p>
          <p class="text-sm font-semibold text-gray-900">{{ isLoadingProfile ? 'Memuat...' : profile.name }}</p>
          <p class="text-sm text-gray-600 mt-1">{{ isLoadingProfile ? '-' : profile.email }}</p>
          <p class="mt-2 inline-flex rounded-full bg-green-50 px-2.5 py-1 text-xs font-medium text-green-600">
            {{ isLoadingProfile ? 'loading' : profile.role }}
          </p>
          <button
            v-if="isLoggedIn"
            type="button"
            class="mt-4 inline-flex w-full items-center justify-center gap-2 rounded-lg border border-red-200 px-3 py-2 text-sm font-medium text-red-600 hover:bg-red-50 transition-colors"
            @click="handleLogout"
          >
            <LogOut class="h-4 w-4" />
            Keluar
          </button>
        </div>
      </div>
    </div>

    <div class="md:hidden px-6 pb-3 flex items-center justify-between">
      <p class="text-sm text-gray-600 font-medium">{{ mobileTitle }}</p>
      <div class="flex items-center gap-2">
        <button
          v-if="!isLoggedIn"
          type="button"
          class="text-sm text-green-600 font-medium"
          @click="router.push({ path: '/login', query: { mode: 'register' } })"
        >
          Daftar
        </button>
        <button
          v-else-if="canEnterDashboard"
          type="button"
          class="text-sm text-green-600 font-medium"
          @click="router.push('/dashboard')"
        >
          Dashboard
        </button>
        <button
          v-else
          type="button"
          class="text-sm text-green-600 font-medium"
          @click="router.push('/pilih-paket')"
        >
          Pilih Paket
        </button>
        <button
          v-if="isLoggedIn"
          type="button"
          class="inline-flex items-center gap-1 rounded-lg border border-gray-200 bg-white px-2.5 py-1.5 text-gray-700"
          @click="toggleProfile"
        >
          <UserCircle2 class="h-4 w-4" />
          <ChevronDown class="h-3.5 w-3.5" />
        </button>
      </div>
    </div>

    <div
      v-if="isLoggedIn && isProfileOpen"
      class="md:hidden mx-6 mb-3 rounded-xl border border-gray-200 bg-white p-3"
    >
      <p class="text-xs uppercase tracking-wide text-gray-500 mb-2">Profil</p>
      <p class="text-sm font-semibold text-gray-900">{{ isLoadingProfile ? 'Memuat...' : profile.name }}</p>
      <p class="text-sm text-gray-600 mt-1">{{ isLoadingProfile ? '-' : profile.email }}</p>
      <p class="mt-2 inline-flex rounded-full bg-green-50 px-2.5 py-1 text-xs font-medium text-green-600">
        {{ isLoadingProfile ? 'loading' : profile.role }}
      </p>
      <button
        v-if="isLoggedIn"
        type="button"
        class="mt-4 inline-flex w-full items-center justify-center gap-2 rounded-lg border border-red-200 px-3 py-2 text-sm font-medium text-red-600 hover:bg-red-50 transition-colors"
        @click="handleLogout"
      >
        <LogOut class="h-4 w-4" />
        Keluar
      </button>
    </div>
  </header>
</template>
