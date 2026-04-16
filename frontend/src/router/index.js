import { createRouter, createWebHistory } from "vue-router";

/**
 * Struktur views (relatif ke src/):
 *   views/auth/Login.vue
 *   views/landing/Landing.vue
 *   views/dashboard/Dashboard.vue
 *   views/harvest/HarvestSales.vue
 *   views/investment/{FinancialProjections,InvestmentPortfolio,InvestorManagement}.vue
 *   views/notifications/{Notifications,ReminderSettings}.vue
 *   views/plant-monitoring/{MaintenanceActivities,MaintenanceValidation,PlantMonitoring}.vue
 *   views/reporting/Reports.vue
 *   views/warehouse/WarehouseStock.vue
 * Layout: layouts/DashboardLayout.vue
 */

const routes = [
  {
    path: "/",
    name: "landing",
    component: () => import("../views/landing/Landing.vue"),
  },
  {
    path: "/login",
    name: "login",
    component: () => import("../views/auth/Login.vue"),
  },
  {
    path: "/dashboard",
    component: () => import("../layouts/DashboardLayout.vue"),
    meta: { requiresAuth: true },
    children: [
      {
        path: "",
        name: "dashboard",
        component: () => import("../views/dashboard/Dashboard.vue"),
      },
      {
        path: "portfolio",
        name: "dashboard-portfolio",
        component: () => import("../views/investment/InvestmentPortfolio.vue"),
      },
      {
        path: "plants",
        name: "dashboard-plants",
        component: () => import("../views/plant-monitoring/PlantMonitoring.vue"),
      },
      {
        path: "maintenance",
        name: "dashboard-maintenance",
        component: () => import("../views/plant-monitoring/MaintenanceActivities.vue"),
      },
      {
        path: "reminder-settings",
        name: "dashboard-reminder-settings",
        component: () => import("../views/notifications/ReminderSettings.vue"),
      },
      {
        path: "financials",
        name: "dashboard-financials",
        component: () => import("../views/investment/FinancialProjections.vue"),
      },
      {
        path: "reports",
        name: "dashboard-reports",
        component: () => import("../views/reporting/Reports.vue"),
      },
      {
        path: "notifications",
        name: "dashboard-notifications",
        component: () => import("../views/notifications/Notifications.vue"),
      },
      {
        path: "investors",
        name: "dashboard-investors",
        component: () => import("../views/investment/InvestorManagement.vue"),
      },
      {
        path: "maintenance-validation",
        name: "dashboard-maintenance-validation",
        component: () => import("../views/plant-monitoring/MaintenanceValidation.vue"),
      },
      {
        path: "harvest-sales",
        name: "dashboard-harvest-sales",
        component: () => import("../views/harvest/HarvestSales.vue"),
      },
      {
        path: "warehouse",
        name: "dashboard-warehouse",
        component: () => import("../views/warehouse/WarehouseStock.vue"),
      },
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
    return next({ name: "login" });
  }

  next();
});

export default router;
