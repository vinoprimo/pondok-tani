import { api } from "../api";

export const getAdminDashboard = async () => {
  const response = await api.get("/admin/dashboard");
  return response.data;
};

export const getUserDashboard = async () => {
  const response = await api.get("/users/me/dashboard");
  return response.data;
};
