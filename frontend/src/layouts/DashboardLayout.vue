<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { Bell } from "lucide-vue-next";
import Sidebar from "../components/dashboard-layouts/Sidebar.vue";
import { clearAuthSession } from "../utils/session";
import {
  fetchSystemNotifications,
  subscribeSystemNotifications,
} from "../services/firebase/systemNotification";

const router = useRouter();
const route = useRoute();

const userRole = ref<"investor" | "mitra" | "admin">(
  (localStorage.getItem("userRole") as "investor" | "mitra" | "admin") || "investor"
);
const unreadNotificationCount = ref(0);
const recentNotifications = ref<any[]>([]);
const showNotifications = ref(false);
let unsubscribeNotifications: (() => void) | null = null;

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

async function loadUnreadNotificationCount() {
  const userId = getCurrentUserId();
  if (!userId) {
    unreadNotificationCount.value = 0;
    return;
  }

  try {
    const items = await fetchSystemNotifications(userId);
    unreadNotificationCount.value = items.filter((item) => !item.read).length;
  } catch (error) {
    console.error('Failed to load unread notifications:', error);
    unreadNotificationCount.value = 0;
  }
}

function startUnreadNotificationListener() {
  const userId = getCurrentUserId();
  if (!userId) return;

  unsubscribeNotifications = subscribeSystemNotifications(userId, (items) => {
    unreadNotificationCount.value = items.filter((item) => !item.read).length;
    recentNotifications.value = items.slice(0, 5); // Tampilkan 5 notifikasi terbaru
  });
}

function toggleNotifications() {
  showNotifications.value = !showNotifications.value;
}

function goToNotifications() {
  showNotifications.value = false;
  router.push('/dashboard/notifications');
}

const roleAllowedViews: Record<"investor" | "mitra" | "admin", string[]> = {
  investor: [
    "dashboard",
    "portfolio",
    "harvest",
    "plants-input",
    "plants-history",
    "maintenance-schedule",
    "maintenance-history",
    "reminder-settings",
    "financials",
    "reports",
    "notifications",
  ],
  mitra: [
    "dashboard",
    "portfolio",
    "harvest",
    "plants-input",
    "plants-history",
    "maintenance-schedule",
    "maintenance-history",
    "reminder-settings",
    "financials",
    "reports",
    "notifications",
  ],
  admin: [
    "dashboard",
    "investors",
    "vanili-management",
    "maintenance-validation",
    "monitoring-schedule",
    "harvest",
    "harvest-sales",
    "warehouse",
    "financials",
    "reports",
    "notifications",
  ],
};

const activeView = computed(() => {
  if (route.path === "/dashboard") return "dashboard";
  if (route.path.startsWith("/dashboard/monitoring-schedule")) return "monitoring-schedule";
  if (route.path.startsWith("/dashboard/plants/history")) return "plants-history";
  if (route.path.startsWith("/dashboard/plants")) return "plants-input";
  if (route.path.startsWith("/dashboard/maintenance-validation")) return "maintenance-validation";
  if (route.path.startsWith("/dashboard/maintenance/history")) return "maintenance-history";
  if (route.path.startsWith("/dashboard/maintenance")) return "maintenance-schedule";

  const matchMap: Record<string, string> = {
    "/dashboard/portfolio": "portfolio",
    "/dashboard/vanili-management": "vanili-management",
    "/dashboard/reminder-settings": "reminder-settings",
    "/dashboard/financials": "financials",
    "/dashboard/reports": "reports",
    "/dashboard/notifications": "notifications",
    "/dashboard/investors": "investors",
    "/dashboard/harvest": "harvest",
    "/dashboard/harvest-sales": "harvest-sales",
    "/dashboard/warehouse": "warehouse",
  };

  return matchMap[route.path] || "dashboard";
});

const pageTitle = computed(() =>
  userRole.value === "admin"
    ? "Dashboard Admin"
    : userRole.value === "mitra"
      ? "Portal Mitra"
      : "Portal Investor"
);

function handleSetActiveView(view: string) {
  const routeMap: Record<string, string> = {
    dashboard: "/dashboard",
    portfolio: "/dashboard/portfolio",
    harvest: "/dashboard/harvest",
    "plants-input": "/dashboard/plants/input",
    "plants-history": "/dashboard/plants/history",
    "vanili-management": "/dashboard/vanili-management",
    maintenance: "/dashboard/maintenance/schedule",
    "maintenance-schedule": "/dashboard/maintenance/schedule",
    "maintenance-history": "/dashboard/maintenance/history",
    "reminder-settings": "/dashboard/reminder-settings",
    financials: "/dashboard/financials",
    reports: "/dashboard/reports",
    notifications: "/dashboard/notifications",
    investors: "/dashboard/investors",
    "maintenance-validation": "/dashboard/maintenance-validation",
    "monitoring-schedule": "/dashboard/monitoring-schedule",
    "harvest-sales": "/dashboard/harvest-sales",
    warehouse: "/dashboard/warehouse",
  };

  const isAllowed = roleAllowedViews[userRole.value].includes(view);
  if (!isAllowed) {
    router.push("/dashboard");
    return;
  }

  router.push(routeMap[view] || "/dashboard");
}

function handleLogout() {
  clearAuthSession();
  router.push("/login");
}

function ensureAllowedCurrentRoute() {
  const currentView = activeView.value;
  const allowed = roleAllowedViews[userRole.value];
  if (!allowed.includes(currentView)) {
    router.replace("/dashboard");
  }
}

onMounted(() => {
  const storedRole = localStorage.getItem("userRole") as "investor" | "mitra" | "admin" | null;
  if (storedRole === "investor" || storedRole === "mitra" || storedRole === "admin") {
    userRole.value = storedRole;
  }
  ensureAllowedCurrentRoute();
  loadUnreadNotificationCount();
  startUnreadNotificationListener();
});

onUnmounted(() => {
  if (unsubscribeNotifications) {
    unsubscribeNotifications();
  }
});

watch(
  () => route.path,
  () => {
    ensureAllowedCurrentRoute();
  }
);
</script>

<template>
  <div class="flex h-dvh overflow-hidden bg-gray-50">
    <Sidebar
      :active-view="activeView"
      :user-role="userRole"
      @set-active-view="handleSetActiveView"
      @logout="handleLogout"
    />

    <div class="flex-1 flex min-h-0 flex-col overflow-hidden">
      <header class="shrink-0 bg-white border-b border-gray-200 px-8 py-4">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-xl font-semibold text-gray-900">{{ pageTitle }}</h1>
            <p class="text-sm text-gray-600 mt-0.5">
              Selamat datang kembali! Berikut perkembangan terkini perkebunan vanili Anda.
            </p>
          </div>
          <div class="relative">
            <button
              type="button"
              class="relative p-2 hover:bg-gray-100 rounded-lg transition-colors"
              @click="toggleNotifications"
            >
              <Bell class="w-6 h-6 text-gray-600" />
              <span
                v-if="unreadNotificationCount > 0"
                class="absolute top-1 right-1 min-w-[10px] h-5 px-1.5 bg-red-500 text-white text-[10px] font-semibold rounded-full flex items-center justify-center"
              >
                {{ unreadNotificationCount }}
              </span>
            </button>

            <!-- Dropdown Notifikasi -->
            <div
              v-if="showNotifications"
              class="absolute right-0 mt-2 w-80 bg-white rounded-xl shadow-lg border border-gray-100 z-50 overflow-hidden"
            >
              <div class="p-4 border-b border-gray-100 flex items-center justify-between">
                <h3 class="font-semibold text-gray-900">Notifikasi</h3>
                <button @click="goToNotifications" class="text-xs text-green-600 hover:text-green-700 font-medium">Lihat Semua</button>
              </div>
              <div class="max-h-96 overflow-y-auto">
                <div v-if="recentNotifications.length === 0" class="p-6 text-center text-gray-500 text-sm">
                  Tidak ada notifikasi baru
                </div>
                <div
                  v-for="notif in recentNotifications"
                  :key="notif.id"
                  class="p-4 border-b border-gray-50 hover:bg-gray-50 transition-colors cursor-pointer"
                  :class="!notif.read ? 'bg-green-50/30' : ''"
                  @click="goToNotifications"
                >
                  <h4 class="text-sm font-medium text-gray-900 mb-1">{{ notif.title || 'Pemberitahuan' }}</h4>
                  <p class="text-xs text-gray-600 line-clamp-2">{{ notif.message }}</p>
                  <p class="text-[10px] text-gray-400 mt-2">{{ notif.createdAt?.toDate ? notif.createdAt.toDate().toLocaleString('id-ID') : 'Baru saja' }}</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </header>

      <main class="flex-1 overflow-auto bg-gray-50">
        <div class="mx-auto w-full max-w-[1400px] px-4 py-5 md:px-6 md:py-6 lg:px-8 lg:py-8">
          <router-view />
        </div>
      </main>
    </div>
  </div>
</template>
