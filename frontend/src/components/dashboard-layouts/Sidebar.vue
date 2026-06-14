<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { useRouter } from 'vue-router';
import {
  ArrowUpRight,
  CalendarDays,
  ChevronDown,
  LayoutDashboard,
  ListChecks,
  PlusCircle,
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

const router = useRouter();

const emit = defineEmits<{
  setActiveView: [view: string];
  logout: [];
}>();

const maintenanceOpen = ref(props.activeView.startsWith('maintenance'));
const plantsOpen = ref(props.activeView.startsWith('plants'));

watch(
  () => props.activeView,
  (value) => {
    if (value.startsWith('maintenance')) {
      maintenanceOpen.value = true;
    }
    if (value.startsWith('plants')) {
      plantsOpen.value = true;
    }
  }
);

const investorMenuItems = [
  { id: 'dashboard', label: 'Beranda', icon: LayoutDashboard },
  { id: 'portfolio', label: 'Portofolio Saya', icon: TrendingUp },
  { id: 'harvest', label: 'Panen', icon: Sprout },
  { id: 'my-warehouse', label: 'Stok Gudang Saya', icon: Package },
  {
    id: 'plants-group',
    label: 'Monitoring Tanaman',
    icon: Sprout,
    children: [
      { id: 'plants-input', label: 'Input Monitoring', icon: PlusCircle },
      { id: 'plants-history', label: 'Riwayat Monitoring', icon: ListChecks },
    ],
  },
  {
    id: 'maintenance-group',
    label: 'Aktivitas Perawatan',
    icon: ClipboardCheck,
    children: [
      { id: 'maintenance-schedule', label: 'Jadwal Perawatan', icon: CalendarDays },
      { id: 'maintenance-history', label: 'Riwayat Aktivitas', icon: ListChecks },
    ],
  },
  { id: 'reminder-settings', label: 'Pengaturan Pengingat', icon: Bell },
  { id: 'financials', label: 'Keuangan', icon: DollarSign },
  { id: 'reports', label: 'Laporan', icon: FileText },
  { id: 'notifications', label: 'Notifikasi', icon: Bell },
];

const mitraMenuItems = investorMenuItems.filter((item) => item.id !== 'maintenance-group');

// Admin memiliki akses ke semua menu investor plus menu tambahan untuk manajemen
const adminMenuItems = [
  { id: 'dashboard', label: 'Beranda', icon: LayoutDashboard },
  { id: 'investors', label: 'Investor', icon: Users },
  { id: 'maintenance-validation', label: 'Validasi Perawatan', icon: ClipboardCheck },
  { id: 'monitoring-schedule', label: 'Jadwal Monitoring', icon: CalendarDays },
  { id: 'harvest', label: 'Panen', icon: Sprout },
  { id: 'harvest-sales', label: 'Penjualan Panen', icon: ShoppingCart },
  { id: 'warehouse', label: 'Stok Gudang', icon: Package },
  { id: 'financials', label: 'Keuangan', icon: DollarSign },
  { id: 'reports', label: 'Laporan', icon: FileText },
  { id: 'vanili-management', label: 'Harga & Grade Vanili', icon: DollarSign },
  { id: 'notifications', label: 'Notifikasi', icon: Bell },
];

const menuItems = computed(() =>
  props.userRole === 'admin' ? adminMenuItems : props.userRole === 'mitra' ? mitraMenuItems : investorMenuItems
);

const roleLabel = computed(() => {
  if (props.userRole === 'admin') return 'Panel Admin';
  if (props.userRole === 'mitra') return 'Portal Mitra';
  return 'Portal Investor';
});
</script>

<template>
  <div class="w-64 h-dvh bg-white border-r border-gray-200 flex min-h-0 flex-col overflow-hidden">
    <div class="shrink-0 px-5 py-4 border-b border-gray-200">
      <div class="flex items-center gap-3">
        <img src="../../assets/logo-pondok-tani.png" alt="Pondok Tani Land" class="w-9 h-9 object-contain rounded-lg bg-white p-1" />
        <div>
          <h1 class="font-bold leading-tight text-green-700">Pondok Tani Land</h1>
          <p class="text-xs text-gray-500">{{ roleLabel }}</p>
        </div>
      </div>
    </div>

    <nav class="flex-1 min-h-0 overflow-y-auto px-3 py-4 space-y-1">
      <div v-for="item in menuItems" :key="item.id">
        <button
          v-if="!item.children"
          type="button"
          class="w-full flex items-center gap-3 px-3 py-2.5 rounded-lg transition-colors"
          :class="
            activeView === item.id
              ? 'bg-green-50 text-green-700'
              : 'text-gray-600 hover:bg-gray-50'
          "
          @click="emit('setActiveView', item.id)"
        >
          <component :is="item.icon" class="w-5 h-5 flex-shrink-0" />
          <span class="text-sm font-medium text-left">{{ item.label }}</span>
        </button>

        <div v-else>
          <button
            type="button"
            class="w-full flex items-center gap-3 px-3 py-2.5 rounded-lg transition-colors"
            :class="
              item.id === 'maintenance-group'
                ? (activeView === 'maintenance-schedule' || activeView === 'maintenance-history')
                  ? 'bg-green-50 text-green-700'
                  : 'text-gray-600 hover:bg-gray-50'
                : (activeView === 'plants-input' || activeView === 'plants-history')
                  ? 'bg-green-50 text-green-700'
                  : 'text-gray-600 hover:bg-gray-50'
            "
            @click="item.id === 'maintenance-group' ? maintenanceOpen = !maintenanceOpen : plantsOpen = !plantsOpen"
          >
            <component :is="item.icon" class="w-5 h-5 flex-shrink-0" />
            <span class="text-sm font-medium text-left">{{ item.label }}</span>
            <ChevronDown
              class="w-4 h-4 ml-auto transition-transform"
              :class="item.id === 'maintenance-group' ? (maintenanceOpen ? 'rotate-180' : '') : (plantsOpen ? 'rotate-180' : '')"
            />
          </button>
          <div v-show="item.id === 'maintenance-group' ? maintenanceOpen : plantsOpen" class="space-y-1 pl-9 mt-1">
            <button
              v-for="child in item.children"
              :key="child.id"
              type="button"
              class="w-full flex items-center gap-2 justify-start rounded-lg px-3 py-2 text-left transition-colors"
              :class="
                activeView === child.id
                  ? 'bg-green-50 text-green-700'
                  : 'text-gray-600 hover:bg-gray-50'
              "
              @click="emit('setActiveView', child.id)"
            >
              <component v-if="child.icon" :is="child.icon" class="w-4 h-4 flex-shrink-0" />
              <span class="text-sm font-medium text-left">{{ child.label }}</span>
            </button>
          </div>
        </div>
      </div>
    </nav>

    <div class="shrink-0 px-3 py-3 border-t border-gray-200 space-y-1">
      <button
        type="button"
        class="w-full inline-flex items-center justify-between gap-2 px-3 py-2 rounded-md text-xs font-medium text-gray-600 hover:bg-gray-100 transition-colors"
        @click="router.push('/')"
      >
        <span>Kembali ke landing</span>
        <ArrowUpRight class="w-4 h-4" />
      </button>
      <button
        type="button"
        class="w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-gray-600 hover:bg-gray-50 transition-colors"
      >
        <Settings class="w-5 h-5" />
        <span class="text-sm font-medium">Pengaturan</span>
      </button>
      <button
        type="button"
        class="w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-red-600 hover:bg-red-50 transition-colors"
        @click="emit('logout')"
      >
        <LogOut class="w-5 h-5" />
        <span class="text-sm font-medium">Keluar</span>
      </button>
    </div>
  </div>
</template>
