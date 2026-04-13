<script setup lang="ts">
import { computed, ref } from 'vue';
import { Bell } from 'lucide-vue-next';
import LandingPage from '../views/Landing.vue';
import LoginPage from '../views/Login.vue';
import Sidebar from '../components/layouts/Sidebar.vue';
import DashboardOverview from '../views/Dashboard.vue';
import InvestmentPortfolio from '../views/InvestmentPortfolio.vue';
import PlantMonitoring from '../views/PlantMonitoring.vue';
import FinancialProjections from '../views/FinancialProjections.vue';
import Reports from '../views/Reports.vue';
import Notifications from '../views/Notifications.vue';
import InvestorManagement from '../views/InvestorManagement.vue';
import HarvestSales from '../views/HarvestSales.vue';
import WarehouseStock from '../views/WarehouseStock.vue';
import MaintenanceActivities from '../views/MaintenanceActivities.vue';
import MaintenanceValidation from '../views/MaintenanceValidation.vue';
import ReminderSettings from '../views/ReminderSettings.vue';

type AppState = 'landing' | 'login' | 'app';
type UserRole = 'investor' | 'admin';

const appState = ref<AppState>('landing');
const isLoggedIn = ref(false);
const userRole = ref<UserRole>('investor');
const activeView = ref('dashboard');
const userInfo = ref({ email: '', name: '' });

function handleGetStarted(role: UserRole) {
  userRole.value = role;
  appState.value = 'login';
}

function handleLogin(role: UserRole, credentials: { email: string; password: string }) {
  // In a real app, you would validate credentials with backend.
  userRole.value = role;
  isLoggedIn.value = true;
  appState.value = 'app';
  activeView.value = 'dashboard';
  userInfo.value = {
    email: credentials.email,
    name: role === 'admin' ? 'Pengguna admin' : 'Sarah Johnson',
  };
}

function handleLogout() {
  isLoggedIn.value = false;
  appState.value = 'landing';
  activeView.value = 'dashboard';
  userInfo.value = { email: '', name: '' };
}

function handleBackToLanding() {
  appState.value = 'landing';
}

const currentViewComponent = computed(() => {
  if (userRole.value === 'admin') {
    const adminViews: Record<string, unknown> = {
      dashboard: DashboardOverview,
      investors: InvestorManagement,
      plants: PlantMonitoring,
      'maintenance-validation': MaintenanceValidation,
      'harvest-sales': HarvestSales,
      warehouse: WarehouseStock,
      financials: FinancialProjections,
      reports: Reports,
      notifications: Notifications,
    };

    return adminViews[activeView.value] ?? DashboardOverview;
  }

  const investorViews: Record<string, unknown> = {
    dashboard: DashboardOverview,
    portfolio: InvestmentPortfolio,
    plants: PlantMonitoring,
    maintenance: MaintenanceActivities,
    'reminder-settings': ReminderSettings,
    financials: FinancialProjections,
    reports: Reports,
    notifications: Notifications,
  };

  return investorViews[activeView.value] ?? DashboardOverview;
});

const currentViewNeedsUserRole = computed(() => currentViewComponent.value === DashboardOverview);
</script>

<template>
  <LandingPage
    v-if="appState === 'landing'"
    @get-started="handleGetStarted"
  />

  <LoginPage
    v-else-if="appState === 'login'"
    :initial-role="userRole"
    @login="handleLogin"
    @back="handleBackToLanding"
  />

  <div v-else class="flex h-screen bg-gray-50">
    <Sidebar
      :active-view="activeView"
      :user-role="userRole"
      @set-active-view="(view) => (activeView = view)"
      @logout="handleLogout"
    />

    <div class="flex-1 flex flex-col overflow-hidden">
      <header class="bg-white border-b border-gray-200 px-8 py-4">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-xl font-semibold text-gray-900">
              {{ userRole === 'admin' ? 'Dashboard Admin' : 'Portal Investor' }}
            </h1>
            <p class="text-sm text-gray-600 mt-0.5">
              Selamat datang kembali! Berikut perkembangan terkini perkebunan vanili Anda.
            </p>
          </div>

          <div class="flex items-center gap-4">
            <div class="flex items-center gap-2 px-3 py-1.5 bg-gray-100 rounded-lg">
              <span class="text-sm text-gray-600">Lihat sebagai:</span>
              <button
                type="button"
                class="px-3 py-1 bg-white rounded-md text-sm font-medium text-gray-900 hover:bg-gray-50 transition-colors border border-gray-200"
                @click="
                  userRole = userRole === 'investor' ? 'admin' : 'investor';
                  activeView = 'dashboard';
                "
              >
                {{ userRole === 'investor' ? 'Ganti ke Admin' : 'Ganti ke Investor' }}
              </button>
            </div>

            <button type="button" class="relative p-2 hover:bg-gray-100 rounded-lg transition-colors">
              <Bell class="w-6 h-6 text-gray-600" />
              <span class="absolute top-1 right-1 w-2 h-2 bg-red-500 rounded-full" />
            </button>

            <div class="flex items-center gap-3 pl-4 border-l border-gray-200">
              <div class="text-right">
                <p class="text-sm font-medium text-gray-900">
                  {{ userInfo.name || (userRole === 'admin' ? 'Pengguna admin' : 'Sarah Johnson') }}
                </p>
                <p class="text-xs text-gray-500">
                  {{ userRole === 'admin' ? 'Administrator Sistem' : 'Investor Premium' }}
                </p>
              </div>
              <div
                class="w-10 h-10 bg-gradient-to-br from-green-400 to-green-600 rounded-full flex items-center justify-center text-white font-semibold"
              >
                {{ userRole === 'admin' ? 'AD' : 'SJ' }}
              </div>
            </div>
          </div>
        </div>
      </header>

      <main class="flex-1 overflow-y-auto px-8 py-6">
        <component
          :is="currentViewComponent"
          v-bind="currentViewNeedsUserRole ? { userRole } : {}"
        />
      </main>
    </div>
  </div>
</template>
