import { createRouter, createWebHistory } from "vue-router";
import { getCurrentUser } from "../services/user/user";
import { clearAuthSession, isSessionExpired, touchSession } from "../utils/session";

/**
 * Struktur views (relatif ke src/):
 *   views/auth/Login.vue
 *   views/landing/Landing.vue
 *   views/dashboard/Dashboard.vue
 *   views/harvest/HarvestRequests.vue
 *   views/harvest/HarvestSales.vue
 *   views/admin/VaniliManagement.vue
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
    path: "/kontak-kami",
    name: "landing-contact",
    component: () => import("../views/landing/ContactUs.vue"),
  },
  {
    path: "/artikel/:id",
    name: "article-detail",
    component: () => import("../views/landing/ArticleDetail.vue"),
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
    path: "/konfirmasi-pesanan",
    name: "order-confirmation",
    component: () => import("../views/investment/OrderConfirmation.vue"),
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
        meta: { roles: ["investor", "mitra"] },
      },
      {
        path: "plants/input",
        name: "dashboard-plants-input",
        component: () => import("../views/plant-monitoring/PlantMonitoring.vue"),
        meta: { roles: ["investor", "mitra"] },
      },
      {
        path: "plants/history",
        name: "dashboard-plants-history",
        component: () => import("../views/plant-monitoring/PlantMonitoring.vue"),
        meta: { roles: ["investor", "mitra"] },
      },
      {
        path: "vanili-management",
        name: "dashboard-vanili-management",
        component: () => import("../views/admin/VaniliManagement.vue"),
        meta: { roles: ["admin"] },
      },
      {
        path: "article-management",
        name: "dashboard-article-management",
        component: () => import("../views/admin/ArticleManagement.vue"),
        meta: { roles: ["admin"] },
      },
      {
        path: "maintenance",
        name: "dashboard-maintenance",
        component: () => import("../views/plant-monitoring/MaintenanceActivities.vue"),
        meta: { roles: ["investor", "mitra"] },
      },
      {
        path: "maintenance/schedule",
        name: "dashboard-maintenance-schedule",
        component: () => import("../views/plant-monitoring/MaintenanceActivities.vue"),
        meta: { roles: ["investor", "mitra"] },
      },
      {
        path: "maintenance/history",
        name: "dashboard-maintenance-history",
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
        path: "monitoring-schedule",
        name: "dashboard-monitoring-schedule",
        component: () => import("../views/plant-monitoring/MonitoringScheduleAdmin.vue"),
        meta: { roles: ["admin"] },
      },
      {
        path: "monitoring-schedule/:userId",
        name: "dashboard-monitoring-schedule-detail",
        component: () => import("../views/plant-monitoring/MonitoringScheduleAdminDetail.vue"),
        meta: { roles: ["admin"] },
      },
      {
        path: "harvest-sales",
        name: "dashboard-harvest-sales",
        component: () => import("../views/harvest/HarvestSales.vue"),
        meta: { roles: ["admin"] },
      },
      {
        path: "harvest",
        name: "dashboard-harvest",
        component: () => import("../views/harvest/HarvestRequests.vue"),
        meta: { roles: ["investor", "mitra", "admin"] },
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
  if (token && isSessionExpired()) {
    clearAuthSession();
    if (to.name !== "login") {
      return next({ name: "login", query: { expired: "1" } });
    }
  }

  const activeToken = localStorage.getItem("token");
  const role = localStorage.getItem("userRole");
  const isAuth = Boolean(activeToken);
  const normalizedRole = role === "admin" ? "admin" : role === "mitra" ? "mitra" : "investor";
  const requiresAuthRoute = to.matched.some((record) => record.meta.requiresAuth);
  const packageFlowRouteNames = ["package-selection", "order-confirmation"];

  if (isAuth) {
    touchSession();
  }

  if (requiresAuthRoute && !isAuth) {
    return next({ name: "login" });
  }

  if (
    requiresAuthRoute &&
    isAuth &&
    (normalizedRole === "investor" || normalizedRole === "mitra") &&
    !packageFlowRouteNames.includes(String(to.name || ""))
  ) {
    try {
      const userRes = await getCurrentUser();
      const hasActivePackage = userRes?.data?.package_status === "active";
      if (!hasActivePackage) {
        return next({ name: "package-selection" });
      }
    } catch (err) {
      clearAuthSession();
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
