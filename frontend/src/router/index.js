import { createRouter, createWebHistory } from "vue-router";
import Landing from "../views/landing/Landing.vue";
import Login from "../views/auth/Login.vue";
import Dashboard from "../views/dashboard/Dashboard.vue";

const routes = [
  { path: "/", component: Landing },
  { path: "/login", component: Login },
  { path: "/dashboard", component: Dashboard },
];

// ✅ BUAT ROUTER DULU
const router = createRouter({
  history: createWebHistory(),
  routes,
});

// ✅ BARU PAKAI beforeEach
router.beforeEach((to, from, next) => {
  const isAuth = localStorage.getItem("token");

  if (to.path === "/dashboard" && !isAuth) {
    return next("/login");
  }

  next();
});

// ✅ EXPORT
export default router;