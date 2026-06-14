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

export const getAdminUsers = (params) => {
  return api.get('/admin/users', { params });
};

export const getAdminUser = (userId) => {
  return api.get(`/admin/users/${userId}`);
};

export const createAdminUser = (payload) => {
  return api.post('/admin/users', payload);
};

export const updateAdminUser = (userId, payload) => {
  return api.put(`/admin/users/${userId}`, payload);
};

export const deleteAdminUser = (userId) => {
  return api.delete(`/admin/users/${userId}`);
};

export const getAdminUserPlantBatches = (userId) => {
  return api.get(`/admin/users/${userId}/plant-batches`);
};

export const activateUserPackage = (userId, payload) => {
  if (payload instanceof FormData) {
    return api.put(`/admin/users/${userId}/activate-package`, payload, {
      headers: {
        "Content-Type": "multipart/form-data",
      },
    });
  }

  return api.put(`/admin/users/${userId}/activate-package`, payload);
};

export const updateInvestmentStatus = (userId, status) => {
  return api.put(`/admin/users/${userId}/investment-status`, { status });
};

export const updateMyFCMToken = (payload) => {
  return api.put('/users/me/fcm-token', payload);
};