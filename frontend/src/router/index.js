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
        meta: { roles: ["investor", "admin"] },
      },
      {
        path: "portfolio",
        name: "dashboard-portfolio",
        component: () => import("../views/investment/InvestmentPortfolio.vue"),
        meta: { roles: ["investor"] },
      },
      {
        path: "plants",
        name: "dashboard-plants",
        component: () => import("../views/plant-monitoring/PlantMonitoring.vue"),
        meta: { roles: ["investor", "admin"] },
      },
      {
        path: "maintenance",
        name: "dashboard-maintenance",
        component: () => import("../views/plant-monitoring/MaintenanceActivities.vue"),
        meta: { roles: ["investor"] },
      },
      {
        path: "reminder-settings",
        name: "dashboard-reminder-settings",
        component: () => import("../views/notifications/ReminderSettings.vue"),
        meta: { roles: ["investor"] },
      },
      {
        path: "financials",
        name: "dashboard-financials",
        component: () => import("../views/investment/FinancialProjections.vue"),
        meta: { roles: ["investor", "admin"] },
      },
      {
        path: "reports",
        name: "dashboard-reports",
        component: () => import("../views/reporting/Reports.vue"),
        meta: { roles: ["investor", "admin"] },
      },
      {
        path: "notifications",
        name: "dashboard-notifications",
        component: () => import("../views/notifications/Notifications.vue"),
        meta: { roles: ["investor", "admin"] },
      },
      {
        path: "investors",
        name: "dashboard-investors",
        component: () => import("../views/investment/InvestorManagement.vue"),
        meta: { roles: ["admin"] },
      },
      {
        path: "maintenance-validation",
        name: "dashboard-maintenance-validation",
        component: () => import("../views/plant-monitoring/MaintenanceValidation.vue"),
        meta: { roles: ["admin"] },
      },
      {
        path: "harvest-sales",
        name: "dashboard-harvest-sales",
        component: () => import("../views/harvest/HarvestSales.vue"),
        meta: { roles: ["admin"] },
      },
      {
        path: "warehouse",
        name: "dashboard-warehouse",
        component: () => import("../views/warehouse/WarehouseStock.vue"),
        meta: { roles: ["admin"] },
      },
    ],
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem("token");
  const role = localStorage.getItem("userRole");
  const isAuth = Boolean(token);

  if (to.matched.some((record) => record.meta.requiresAuth) && !isAuth) {
    return next({ name: "login" });
  }

  if (to.name === "login" && isAuth) {
    return next({ name: "dashboard" });
  }

  const roleRule = to.matched.find((record) => Array.isArray(record.meta?.roles));
  if (roleRule) {
    const allowedRoles = roleRule.meta.roles;
    if (!role || !allowedRoles.includes(role)) {
      return next({ name: "dashboard" });
    }
  }

  next();
});

export default router;
