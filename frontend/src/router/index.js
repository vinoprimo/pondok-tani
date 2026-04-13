import { createRouter, createWebHistory } from "vue-router";
import Landing from "../views/Landing.vue";
import Login from "../views/Login.vue";
import Dashboard from "../views/Dashboard.vue";
import DashboardLayout from "../layouts/DashboardLayout.vue";
import InvestmentPortfolio from "../views/InvestmentPortfolio.vue";
import PlantMonitoring from "../views/PlantMonitoring.vue";
import MaintenanceActivities from "../views/MaintenanceActivities.vue";
import ReminderSettings from "../views/ReminderSettings.vue";
import FinancialProjections from "../views/FinancialProjections.vue";
import Reports from "../views/Reports.vue";
import Notifications from "../views/Notifications.vue";
import InvestorManagement from "../views/InvestorManagement.vue";
import MaintenanceValidation from "../views/MaintenanceValidation.vue";
import HarvestSales from "../views/HarvestSales.vue";
import WarehouseStock from "../views/WarehouseStock.vue";

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