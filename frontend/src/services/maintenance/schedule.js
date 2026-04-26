import { api } from "../api";

export const getAdminMaintenanceSchedules = (params = {}) => {
  return api.get("/admin/maintenance-schedules", { params });
};

export const getAdminMaintenanceScheduleSummary = () => {
  return api.get("/admin/maintenance-schedules/summary");
};

export const createAdminMaintenanceSchedule = (payload) => {
  return api.post("/admin/maintenance-schedules", payload);
};
