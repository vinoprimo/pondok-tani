import { api } from "../api";

export const getCurrentUser = () => {
  return api.get("/auth/me");
};

export const saveSelectedPackage = (packageId) => {
  return api.put('/users/me/package', { package_id: packageId });
};

export const getAdminUsers = () => {
  return api.get('/admin/users');
};

export const activateUserPackage = (userId) => {
  return api.put(`/admin/users/${userId}/activate-package`);
};