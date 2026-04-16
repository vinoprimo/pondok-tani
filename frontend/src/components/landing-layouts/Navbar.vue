<script setup lang="ts">
import { computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { Sprout } from 'lucide-vue-next';

const router = useRouter();
const route = useRoute();

const navItems = [
  { label: 'Beranda', path: '/' },
  { label: 'Paket Investasi', path: '/paket-investasi' },
  { label: 'Kemitraan', path: '/kemitraan' },
];

const isActive = (path: string) => {
  if (path === '/') return route.path === '/';
  return route.path.startsWith(path);
};

const mobileTitle = computed(() => {
  if (route.path === '/paket-investasi') return 'Paket Investasi';
  if (route.path === '/kemitraan') return 'Kemitraan';
  return 'Beranda';
});
</script>

<template>
  <header class="border-b border-gray-200 bg-white/80 backdrop-blur-sm sticky top-0 z-50">
    <div class="max-w-7xl mx-auto px-6 py-4 flex items-center justify-between gap-4">
      <button type="button" class="flex items-center gap-3" @click="router.push('/')">
        <div class="w-10 h-10 bg-green-600 rounded-lg flex items-center justify-center">
          <Sprout class="w-6 h-6 text-white" />
        </div>
        <div class="text-left">
          <h1 class="font-semibold text-gray-900">Omah Vanili</h1>
          <p class="text-xs text-gray-500">Sistem Investasi Vanili</p>
        </div>
      </button>

      <nav class="hidden md:flex items-center gap-1">
        <button
          v-for="item in navItems"
          :key="item.path"
          type="button"
          class="px-4 py-2 rounded-lg font-medium transition-colors"
          :class="isActive(item.path) ? 'bg-green-100 text-green-700' : 'text-gray-700 hover:bg-gray-100'"
          @click="router.push(item.path)"
        >
          {{ item.label }}
        </button>
      </nav>

      <div class="hidden md:flex items-center gap-3">
        <button
          type="button"
          class="px-4 py-2 text-green-700 hover:bg-green-50 rounded-lg transition-colors font-medium"
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
      </div>
    </div>

    <div class="md:hidden px-6 pb-3 flex items-center justify-between">
      <p class="text-sm text-gray-600 font-medium">{{ mobileTitle }}</p>
      <button
        type="button"
        class="text-sm text-green-700 font-medium"
        @click="router.push({ path: '/login', query: { mode: 'register' } })"
      >
        Daftar
      </button>
    </div>
  </header>
</template>
