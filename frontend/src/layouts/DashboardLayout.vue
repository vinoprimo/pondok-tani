<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { Bell } from "lucide-vue-next";
import Sidebar from "../components/dashboard-layouts/Sidebar.vue";
import { clearAuthSession } from "../utils/session";

const router = useRouter();
const route = useRoute();

const userRole = ref<"investor" | "mitra" | "admin">(
  (localStorage.getItem("userRole") as "investor" | "mitra" | "admin") || "investor"
);

const roleAllowedViews: Record<"investor" | "mitra" | "admin", string[]> = {
  investor: [
    "dashboard",
    "portfolio",
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
});

watch(
  () => route.path,
  () => {
    ensureAllowedCurrentRoute();
  }
);
</script>

<template>
  <div class="flex h-screen bg-gray-50">
    <Sidebar
      :active-view="activeView"
      :user-role="userRole"
      @set-active-view="handleSetActiveView"
      @logout="handleLogout"
    />

    <div class="flex-1 flex flex-col overflow-hidden">
      <header class="bg-white border-b border-gray-200 px-8 py-4">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-xl font-semibold text-gray-900">{{ pageTitle }}</h1>
            <p class="text-sm text-gray-600 mt-0.5">
              Selamat datang kembali! Berikut perkembangan terkini perkebunan vanili Anda.
            </p>
          </div>
          <button type="button" class="relative p-2 hover:bg-gray-100 rounded-lg transition-colors">
            <Bell class="w-6 h-6 text-gray-600" />
            <span class="absolute top-1 right-1 w-2 h-2 bg-red-500 rounded-full" />
          </button>
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