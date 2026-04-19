import { api } from "../api";

export const getCurrentUser = () => {
  return api.get("/auth/me");
};

export const saveSelectedPackage = (packageIds) => {
  const normalizedPackageIds = Array.isArray(packageIds)
    ? packageIds.filter(Boolean)
    : [packageIds].filter(Boolean);

  return api.put('/users/me/package', {
    package_id: normalizedPackageIds[0],
    package_ids: normalizedPackageIds,
  });
};

export const getAdminUsers = () => {
  return api.get('/admin/users');
};

export const activateUserPackage = (userId) => {
  return api.put(`/admin/users/${userId}/activate-package`);
};