<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue';
import { Bell, CheckCircle2, AlertCircle, Info, TrendingUp, X } from 'lucide-vue-next';
import {
  fetchSystemNotifications,
  subscribeSystemNotifications,
  markSystemNotificationRead,
} from '../../services/firebase/systemNotification';

const filter = ref('all');
const notifications = ref([]);

function getCurrentUserId() {
  const storedUserId = localStorage.getItem('userId');
  if (storedUserId) return storedUserId;

  const token = localStorage.getItem('token');
  if (!token) return null;

  try {
    const payload = JSON.parse(atob(token.split('.')[1]));
    return String(payload?.id || payload?.user_id || payload?.sub || '');
  } catch {
    return null;
  }
}

const filterItems = [
  { id: 'all', label: 'Semua' },
  { id: 'unread', label: 'Belum dibaca' },
];

const filteredNotifications = computed(() => {
  if (filter.value === 'all') return notifications.value;
  if (filter.value === 'unread') return notifications.value.filter((n) => !n.read);
  return notifications.value.filter((n) => n.category === filter.value);
});

const unreadCount = computed(() => notifications.value.filter((n) => !n.read).length);
const plantCount = computed(() => notifications.value.filter((n) => n.category === 'plant').length);
const alertCount = computed(() => notifications.value.filter((n) => n.category === 'alert').length);
const paymentCount = computed(() => notifications.value.filter((n) => n.category === 'payment').length);

let unsubscribeNotifications = null;

async function markAllAsRead() {
  const unreadItems = notifications.value.filter((n) => !n.read);
  if (!unreadItems.length) return;

  await Promise.all(
    unreadItems.map((notification) =>
      markSystemNotificationRead(notification.id).catch((error) => {
        console.error('Failed to mark notification read:', error);
      })
    )
  );

  notifications.value = notifications.value.map((notification) => ({
    ...notification,
    read: true,
  }));
}

async function markAsRead(notification) {
  if (notification.read || !notification.id) return;
  try {
    await markSystemNotificationRead(notification.id);
  } catch (error) {
    console.error('Failed to mark notification read:', error);
  }
}

function normalizeNotification(doc) {
  const createdAt = doc.createdAt && typeof doc.createdAt.toDate === 'function'
    ? doc.createdAt.toDate().toLocaleString('id-ID', { day: '2-digit', month: 'long', year: 'numeric' })
    : doc.time || 'Baru saja';

  return {
    ...doc,
    time: createdAt,
  };
}

async function loadNotifications() {
  const userId = getCurrentUserId();
  if (!userId) return;

  try {
    const items = await fetchSystemNotifications(userId);
    notifications.value = items.map(normalizeNotification);
  } catch (error) {
    console.error('Failed to load system notifications:', error);
  }
}

function startNotificationListener() {
  const userId = getCurrentUserId();
  if (!userId) return;

  unsubscribeNotifications = subscribeSystemNotifications(userId, (items) => {
    notifications.value = items.map(normalizeNotification);
  });
}

onMounted(() => {
  loadNotifications();
  startNotificationListener();
});

onUnmounted(() => {
  if (unsubscribeNotifications) unsubscribeNotifications();
});

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
        <button
          type="button"
          class="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          :disabled="unreadCount === 0"
          @click="markAllAsRead"
        >
          Tandai semua dibaca
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
            <p class="text-2xl font-semibold text-gray-900 mt-1">{{ plantCount }}</p>
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
            <p class="text-2xl font-semibold text-gray-900 mt-1">{{ alertCount }}</p>
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
            <p class="text-2xl font-semibold text-gray-900 mt-1">{{ paymentCount }}</p>
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
                @click.prevent="markAsRead(notification)"
              >
                Tandai dibaca
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

  </div>
</template>
