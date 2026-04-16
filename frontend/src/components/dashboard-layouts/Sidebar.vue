<script setup lang="ts">
import { computed } from 'vue';
import {
  LayoutDashboard,
  Sprout,
  TrendingUp,
  DollarSign,
  FileText,
  Bell,
  Settings,
  Users,
  LogOut,
  ShoppingCart,
  Package,
  ClipboardCheck,
} from 'lucide-vue-next';

const props = defineProps<{
  activeView: string;
  userRole: 'investor' | 'mitra' | 'admin';
}>();

const emit = defineEmits<{
  setActiveView: [view: string];
  logout: [];
}>();

const investorMenuItems = [
  { id: 'dashboard', label: 'Beranda', icon: LayoutDashboard },
  { id: 'portfolio', label: 'Portofolio Saya', icon: TrendingUp },
  { id: 'plants', label: 'Monitoring Tanaman', icon: Sprout },
  { id: 'maintenance', label: 'Aktivitas Perawatan', icon: ClipboardCheck },
  { id: 'reminder-settings', label: 'Pengaturan Pengingat', icon: Bell },
  { id: 'financials', label: 'Keuangan', icon: DollarSign },
  { id: 'reports', label: 'Laporan', icon: FileText },
  { id: 'notifications', label: 'Notifikasi', icon: Bell },
];

const adminMenuItems = [
  { id: 'dashboard', label: 'Beranda', icon: LayoutDashboard },
  { id: 'investors', label: 'Investor', icon: Users },
  { id: 'plants', label: 'Manajemen Tanaman', icon: Sprout },
  { id: 'maintenance-validation', label: 'Validasi Perawatan', icon: ClipboardCheck },
  { id: 'harvest-sales', label: 'Penjualan Panen', icon: ShoppingCart },
  { id: 'warehouse', label: 'Stok Gudang', icon: Package },
  { id: 'financials', label: 'Keuangan', icon: DollarSign },
  { id: 'reports', label: 'Laporan', icon: FileText },
  { id: 'notifications', label: 'Notifikasi', icon: Bell },
];

const menuItems = computed(() =>
  props.userRole === 'admin' ? adminMenuItems : investorMenuItems
);

const roleLabel = computed(() => {
  if (props.userRole === 'admin') return 'Panel Admin';
  if (props.userRole === 'mitra') return 'Portal Mitra';
  return 'Portal Investor';
});
</script>

<template>
  <div class="w-64 bg-white border-r border-gray-200 flex flex-col h-screen">
    <div class="p-6 border-b border-gray-200">
      <div class="flex items-center gap-3">
        <div class="w-10 h-10 bg-green-600 rounded-lg flex items-center justify-center">
          <Sprout class="w-6 h-6 text-white" />
        </div>
        <div>
          <h1 class="font-semibold text-gray-900">Omah Vanili</h1>
          <p class="text-xs text-gray-500">{{ roleLabel }}</p>
        </div>
      </div>
    </div>

    <nav class="flex-1 p-4 space-y-1">
      <button
        v-for="item in menuItems"
        :key="item.id"
        type="button"
        class="w-full flex items-center gap-3 px-4 py-3 rounded-lg transition-colors"
        :class="
          activeView === item.id
            ? 'bg-green-50 text-green-700'
            : 'text-gray-600 hover:bg-gray-50'
        "
        @click="emit('setActiveView', item.id)"
      >
        <component :is="item.icon" class="w-5 h-5 flex-shrink-0" />
        <span class="font-medium text-left">{{ item.label }}</span>
      </button>
    </nav>

    <div class="p-4 border-t border-gray-200 space-y-1">
      <button
        type="button"
        class="w-full flex items-center gap-3 px-4 py-3 rounded-lg text-gray-600 hover:bg-gray-50 transition-colors"
      >
        <Settings class="w-5 h-5" />
        <span class="font-medium">Pengaturan</span>
      </button>
      <button
        type="button"
        class="w-full flex items-center gap-3 px-4 py-3 rounded-lg text-red-600 hover:bg-red-50 transition-colors"
        @click="emit('logout')"
      >
        <LogOut class="w-5 h-5" />
        <span class="font-medium">Keluar</span>
      </button>
    </div>
  </div>
</template>
