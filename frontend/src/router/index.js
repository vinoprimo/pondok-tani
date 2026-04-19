import { createRouter, createWebHistory } from "vue-router";
import { getCurrentUser } from "../services/user/user";

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
    path: "/paket-investasi",
    name: "landing-paket",
    component: () => import("../views/landing/LandingPaket.vue"),
  },
  {
    path: "/kemitraan",
    name: "landing-kemitraan",
    component: () => import("../views/landing/LandingKemitraan.vue"),
  },
  {
    path: "/login",
    name: "login",
    component: () => import("../views/auth/Login.vue"),
  },
  {
    path: "/pilih-paket",
    name: "package-selection",
    component: () => import("../views/investment/PackageSelection.vue"),
    meta: { requiresAuth: true, roles: ["investor", "mitra"] },
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
        meta: { roles: ["investor", "mitra", "admin"] },
      },
      {
        path: "portfolio",
        name: "dashboard-portfolio",
        component: () => import("../views/investment/InvestmentPortfolio.vue"),
        meta: { roles: ["investor", "mitra"] },
      },
      {
        path: "plants",
        name: "dashboard-plants",
        component: () => import("../views/plant-monitoring/PlantMonitoring.vue"),
        meta: { roles: ["investor", "mitra", "admin"] },
      },
      {
        path: "maintenance",
        name: "dashboard-maintenance",
        component: () => import("../views/plant-monitoring/MaintenanceActivities.vue"),
        meta: { roles: ["investor", "mitra"] },
      },
      {
        path: "reminder-settings",
        name: "dashboard-reminder-settings",
        component: () => import("../views/notifications/ReminderSettings.vue"),
        meta: { roles: ["investor", "mitra"] },
      },
      {
        path: "financials",
        name: "dashboard-financials",
        component: () => import("../views/investment/FinancialProjections.vue"),
        meta: { roles: ["investor", "mitra", "admin"] },
      },
      {
        path: "reports",
        name: "dashboard-reports",
        component: () => import("../views/reporting/Reports.vue"),
        meta: { roles: ["investor", "mitra", "admin"] },
      },
      {
        path: "notifications",
        name: "dashboard-notifications",
        component: () => import("../views/notifications/Notifications.vue"),
        meta: { roles: ["investor", "mitra", "admin"] },
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

router.beforeEach(async (to, from, next) => {
  const token = localStorage.getItem("token");
  const role = localStorage.getItem("userRole");
  const isAuth = Boolean(token);
  const normalizedRole = role === "admin" ? "admin" : role === "mitra" ? "mitra" : "investor";
  const requiresAuthRoute = to.matched.some((record) => record.meta.requiresAuth);

  if (requiresAuthRoute && !isAuth) {
    return next({ name: "login" });
  }

  if (
    requiresAuthRoute &&
    isAuth &&
    (normalizedRole === "investor" || normalizedRole === "mitra") &&
    to.name !== "package-selection"
  ) {
    try {
      const userRes = await getCurrentUser();
      const hasSelectedPackage = Boolean(userRes?.data?.investments?.length || userRes?.data?.selected_package_id);
      if (!hasSelectedPackage) {
        return next({ name: "package-selection" });
      }
    } catch (err) {
      localStorage.removeItem("token");
      localStorage.removeItem("userRole");
      return next({ name: "login" });
    }
  }

  if (to.name === "login" && isAuth) {
    return next({ name: "dashboard" });
  }

  const roleRule = to.matched.find((record) => Array.isArray(record.meta?.roles));
  if (roleRule) {
    const allowedRoles = roleRule.meta.roles;
    if (!allowedRoles.includes(normalizedRole)) {
      return next({ name: "dashboard" });
    }
  }

  next();
});

export default router;
