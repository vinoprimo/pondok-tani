import { createRouter, createWebHistory } from "vue-router";
import Landing from "../views/landing/Landing.vue";
import Login from "../views/auth/Login.vue";
import DashboardLayout from "../layouts/DashboardLayout.vue";
import Dashboard from "../views/dashboard/Dashboard.vue";
import InvestmentPortfolio from "../views/investment/InvestmentPortfolio.vue";
import PlantMonitoring from "../views/plant-monitoring/PlantMonitoring.vue";
import MaintenanceActivities from "../views/plant-monitoring/MaintenanceActivities.vue";
import ReminderSettings from "../views/notifications/ReminderSettings.vue";
import FinancialProjections from "../views/investment/FinancialProjections.vue";
import Reports from "../views/reporting/Reports.vue";
import Notifications from "../views/notifications/Notifications.vue";
import InvestorManagement from "../views/investment/InvestorManagement.vue";
import MaintenanceValidation from "../views/plant-monitoring/MaintenanceValidation.vue";
import HarvestSales from "../views/harvest/HarvestSales.vue";
import WarehouseStock from "../views/warehouse/WarehouseStock.vue";

const routes = [
  { path: "/", component: Landing },
  { path: "/login", component: Login },
  {
    path: "/dashboard",
    component: DashboardLayout,
    meta: { requiresAuth: true },
    children: [
      { path: "", component: Dashboard },
      { path: "portfolio", component: InvestmentPortfolio },
      { path: "plants", component: PlantMonitoring },
      { path: "maintenance", component: MaintenanceActivities },
      { path: "reminder-settings", component: ReminderSettings },
      { path: "financials", component: FinancialProjections },
      { path: "reports", component: Reports },
      { path: "notifications", component: Notifications },
      { path: "investors", component: InvestorManagement },
      { path: "maintenance-validation", component: MaintenanceValidation },
      { path: "harvest-sales", component: HarvestSales },
      { path: "warehouse", component: WarehouseStock },
    ],
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach((to, from, next) => {
  const isAuth = localStorage.getItem("token");

  if (to.matched.some((record) => record.meta.requiresAuth) && !isAuth) {
    return next("/login");
  }

  next();
});

export default router;
