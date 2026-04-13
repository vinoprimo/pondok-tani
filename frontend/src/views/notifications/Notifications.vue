<script setup lang="ts">
import { computed, ref } from 'vue';
import { Bell, CheckCircle2, AlertCircle, Info, TrendingUp, X } from 'lucide-vue-next';

const filter = ref('all');

const notifications = [
  {
    id: 1,
    type: 'success',
    title: 'Imbal bulanan telah dicairkan',
    message: 'Imbal Juni sebesar $918 telah masuk ke rekening Anda.',
    time: '2 jam lalu',
    read: false,
    category: 'payment',
  },
  {
    id: 2,
    type: 'info',
    title: 'Inspeksi tanaman dijadwalkan',
    message: 'Inspeksi kesehatan rutin untuk Tanaman #001 dijadwalkan 25 Januari 2026.',
    time: '5 jam lalu',
    read: false,
    category: 'plant',
  },
  {
    id: 3,
    type: 'success',
    title: 'Tahap pertumbuhan baru tercapai',
    message: 'Tanaman #001 memasuki fase berbunga. Perkiraan pembentukan polong dalam 3 bulan.',
    time: '1 hari lalu',
    read: true,
    category: 'plant',
  },
  {
    id: 4,
    type: 'alert',
    title: 'Peringatan cuaca',
    message: 'Hujan lebat diperkirakan di area budidaya. Memantau sistem drainase.',
    time: '1 hari lalu',
    read: true,
    category: 'alert',
  },
  {
    id: 5,
    type: 'info',
    title: 'Laporan kuartal tersedia',
    message: 'Laporan investasi Q4 2025 Anda siap diunduh.',
    time: '2 hari lalu',
    read: true,
    category: 'report',
  },
  {
    id: 6,
    type: 'success',
    title: 'Panen selesai',
    message: 'Panen batch #28 selesai. Hasil: 142 kg biji vanili premium.',
    time: '3 hari lalu',
    read: true,
    category: 'plant',
  },
  {
    id: 7,
    type: 'info',
    title: 'Pembaruan ROI',
    message: 'ROI portofolio Anda naik menjadi 24,5% bulan ini (sebelumnya 23,2%).',
    time: '1 minggu lalu',
    read: true,
    category: 'investment',
  },
  {
    id: 8,
    type: 'alert',
    title: 'Perawatan dijadwalkan',
    message: 'Perawatan sistem irigasi di Sektor A dijadwalkan 1 Februari.',
    time: '1 minggu lalu',
    read: true,
    category: 'alert',
  },
];

const filterItems = [
  { id: 'all', label: 'Semua' },
  { id: 'unread', label: 'Belum dibaca' },
  { id: 'plant', label: 'Pembaruan tanaman' },
  { id: 'payment', label: 'Pembayaran' },
  { id: 'alert', label: 'Peringatan' },
  { id: 'report', label: 'Laporan' },
  { id: 'investment', label: 'Investasi' },
];

const preferenceRows = [
  {
    label: 'Notifikasi pembayaran',
    description: 'Diberitahu saat imbal hasil dicairkan',
    enabled: true,
  },
  {
    label: 'Pembaruan tanaman',
    description: 'Perubahan fase pertumbuhan dan status kesehatan',
    enabled: true,
  },
  {
    label: 'Laporan tersedia',
    description: 'Saat laporan baru siap diunduh',
    enabled: true,
  },
  {
    label: 'Peringatan cuaca',
    description: 'Informasi cuaca penting untuk tanaman Anda',
    enabled: true,
  },
  {
    label: 'Wawasan investasi',
    description: 'Pembaruan ROI dan tren pasar',
    enabled: false,
  },
  {
    label: 'Email pemasaran',
    description: 'Berita dan konten promosi',
    enabled: false,
  },
];

const filteredNotifications = computed(() => {
  if (filter.value === 'all') return notifications;
  if (filter.value === 'unread') return notifications.filter((n) => !n.read);
  return notifications.filter((n) => n.category === filter.value);
});

const unreadCount = computed(() => notifications.filter((n) => !n.read).length);

function getIconBg(type: string) {
  switch (type) {
    case 'success':
      return 'bg-green-50';
    case 'alert':
      return 'bg-yellow-50';
    default:
      return 'bg-blue-50';
  }
}
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-semibold text-gray-900">Notifikasi</h2>
        <p class="text-gray-600 mt-1">Ikuti perkembangan aktivitas investasi Anda</p>
      </div>
      <div class="flex items-center gap-3">
        <button type="button" class="px-4 py-2 text-gray-600 hover:text-gray-900 font-medium text-sm">
          Tandai semua dibaca
        </button>
        <button
          type="button"
          class="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors"
        >
          Pengaturan
        </button>
      </div>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
      <div class="bg-white rounded-xl border border-gray-200 p-4">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-gray-600">Belum dibaca</p>
            <p class="text-2xl font-semibold text-gray-900 mt-1">{{ unreadCount }}</p>
          </div>
          <div class="w-12 h-12 bg-red-50 rounded-lg flex items-center justify-center">
            <Bell class="w-6 h-6 text-red-600" />
          </div>
        </div>
      </div>

      <div class="bg-white rounded-xl border border-gray-200 p-4">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-gray-600">Pembaruan tanaman</p>
            <p class="text-2xl font-semibold text-gray-900 mt-1">3</p>
          </div>
          <div class="w-12 h-12 bg-green-50 rounded-lg flex items-center justify-center">
            <CheckCircle2 class="w-6 h-6 text-green-600" />
          </div>
        </div>
      </div>

      <div class="bg-white rounded-xl border border-gray-200 p-4">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-gray-600">Peringatan</p>
            <p class="text-2xl font-semibold text-gray-900 mt-1">2</p>
          </div>
          <div class="w-12 h-12 bg-yellow-50 rounded-lg flex items-center justify-center">
            <AlertCircle class="w-6 h-6 text-yellow-600" />
          </div>
        </div>
      </div>

      <div class="bg-white rounded-xl border border-gray-200 p-4">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-gray-600">Keuangan</p>
            <p class="text-2xl font-semibold text-gray-900 mt-1">2</p>
          </div>
          <div class="w-12 h-12 bg-blue-50 rounded-lg flex items-center justify-center">
            <TrendingUp class="w-6 h-6 text-blue-600" />
          </div>
        </div>
      </div>
    </div>

    <div class="bg-white rounded-xl border border-gray-200 p-4">
      <div class="flex flex-wrap gap-2">
        <button
          v-for="item in filterItems"
          :key="item.id"
          type="button"
          :class="[
            'px-4 py-2 rounded-lg text-sm font-medium transition-colors',
            filter === item.id
              ? 'bg-green-600 text-white'
              : 'bg-gray-100 text-gray-700 hover:bg-gray-200',
          ]"
          @click="filter = item.id"
        >
          {{ item.label }}
        </button>
      </div>
    </div>

    <div class="space-y-3">
      <div
        v-for="notification in filteredNotifications"
        :key="notification.id"
        :class="[
          'bg-white rounded-xl border p-6 transition-all',
          notification.read ? 'border-gray-200' : 'border-green-200 bg-green-50/30',
        ]"
      >
        <div class="flex items-start gap-4">
          <div
            class="w-10 h-10 rounded-lg flex items-center justify-center flex-shrink-0"
            :class="getIconBg(notification.type)"
          >
            <CheckCircle2
              v-if="notification.type === 'success'"
              class="w-5 h-5 text-green-600"
            />
            <AlertCircle
              v-else-if="notification.type === 'alert'"
              class="w-5 h-5 text-yellow-600"
            />
            <Info v-else class="w-5 h-5 text-blue-600" />
          </div>

          <div class="flex-1 min-w-0">
            <div class="flex items-start justify-between gap-4 mb-2">
              <h3 class="font-semibold text-gray-900">{{ notification.title }}</h3>
              <div class="flex items-center gap-2 flex-shrink-0">
                <span
                  v-if="!notification.read"
                  class="w-2 h-2 bg-green-600 rounded-full"
                />
                <button type="button" class="text-gray-400 hover:text-gray-600">
                  <X class="w-4 h-4" />
                </button>
              </div>
            </div>
            <p class="text-gray-600 mb-2">{{ notification.message }}</p>
            <div class="flex items-center gap-4">
              <span class="text-sm text-gray-500">{{ notification.time }}</span>
              <button
                v-if="!notification.read"
                type="button"
                class="text-sm text-green-600 hover:text-green-700 font-medium"
              >
                Tandai dibaca
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="bg-white rounded-xl border border-gray-200 p-6">
      <h3 class="font-semibold text-gray-900 mb-4">Preferensi notifikasi</h3>
      <div class="space-y-4">
        <div
          v-for="(pref, index) in preferenceRows"
          :key="index"
          class="flex items-center justify-between py-3 border-b border-gray-100 last:border-0"
        >
          <div class="flex-1">
            <p class="font-medium text-gray-900">{{ pref.label }}</p>
            <p class="text-sm text-gray-600 mt-0.5">{{ pref.description }}</p>
          </div>
          <button
            type="button"
            :class="[
              'relative inline-flex h-6 w-11 items-center rounded-full transition-colors',
              pref.enabled ? 'bg-green-600' : 'bg-gray-200',
            ]"
          >
            <span
              :class="[
                'inline-block h-4 w-4 transform rounded-full bg-white transition-transform',
                pref.enabled ? 'translate-x-6' : 'translate-x-1',
              ]"
            />
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
