<script setup lang="ts">
import { computed, ref } from 'vue';
import { Bell, CheckCircle2, AlertCircle, Info, TrendingUp, X } from 'lucide-vue-next';

const filter = ref('all');

const notifications = [
  {
    id: 1,
    type: 'success',
    title: 'Monthly return deposited',
    message: 'Your June return of $918 has been deposited to your account.',
    time: '2 hours ago',
    read: false,
    category: 'payment',
  },
  {
    id: 2,
    type: 'info',
    title: 'Plant inspection scheduled',
    message: 'Routine health inspection scheduled for Plant #001 on January 25, 2026.',
    time: '5 hours ago',
    read: false,
    category: 'plant',
  },
  {
    id: 3,
    type: 'success',
    title: 'New growth stage reached',
    message: 'Plant #001 has entered the flowering stage. Expected pod development in 3 months.',
    time: '1 day ago',
    read: true,
    category: 'plant',
  },
  {
    id: 4,
    type: 'alert',
    title: 'Weather advisory',
    message: 'Heavy rainfall expected in cultivation area. Monitoring drainage systems.',
    time: '1 day ago',
    read: true,
    category: 'alert',
  },
  {
    id: 5,
    type: 'info',
    title: 'Quarterly report available',
    message: 'Your Q4 2025 investment report is now available for download.',
    time: '2 days ago',
    read: true,
    category: 'report',
  },
  {
    id: 6,
    type: 'success',
    title: 'Harvest completed',
    message: 'Batch #28 harvest completed successfully. Yield: 142kg of premium vanilla beans.',
    time: '3 days ago',
    read: true,
    category: 'plant',
  },
  {
    id: 7,
    type: 'info',
    title: 'ROI update',
    message: 'Your portfolio ROI increased to 24.5% this month, up from 23.2%.',
    time: '1 week ago',
    read: true,
    category: 'investment',
  },
  {
    id: 8,
    type: 'alert',
    title: 'Maintenance scheduled',
    message: 'Irrigation system maintenance scheduled for Sector A on February 1.',
    time: '1 week ago',
    read: true,
    category: 'alert',
  },
];

const filterItems = [
  { id: 'all', label: 'All' },
  { id: 'unread', label: 'Unread' },
  { id: 'plant', label: 'Plant Updates' },
  { id: 'payment', label: 'Payments' },
  { id: 'alert', label: 'Alerts' },
  { id: 'report', label: 'Reports' },
  { id: 'investment', label: 'Investment' },
];

const preferenceRows = [
  {
    label: 'Payment notifications',
    description: 'Get notified when returns are deposited',
    enabled: true,
  },
  {
    label: 'Plant updates',
    description: 'Growth stages and health status changes',
    enabled: true,
  },
  {
    label: 'Reports available',
    description: 'When new reports are ready for download',
    enabled: true,
  },
  {
    label: 'Weather alerts',
    description: 'Important weather advisories for your plants',
    enabled: true,
  },
  {
    label: 'Investment insights',
    description: 'ROI updates and market trends',
    enabled: false,
  },
  {
    label: 'Marketing emails',
    description: 'News and promotional content',
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
        <h2 class="text-2xl font-semibold text-gray-900">Notifications</h2>
        <p class="text-gray-600 mt-1">Stay updated with your investment activities</p>
      </div>
      <div class="flex items-center gap-3">
        <button type="button" class="px-4 py-2 text-gray-600 hover:text-gray-900 font-medium text-sm">
          Mark all as read
        </button>
        <button
          type="button"
          class="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors"
        >
          Settings
        </button>
      </div>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
      <div class="bg-white rounded-xl border border-gray-200 p-4">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-gray-600">Unread</p>
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
            <p class="text-sm text-gray-600">Plant Updates</p>
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
            <p class="text-sm text-gray-600">Alerts</p>
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
            <p class="text-sm text-gray-600">Financial</p>
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
                Mark as read
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="bg-white rounded-xl border border-gray-200 p-6">
      <h3 class="font-semibold text-gray-900 mb-4">Notification Preferences</h3>
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
