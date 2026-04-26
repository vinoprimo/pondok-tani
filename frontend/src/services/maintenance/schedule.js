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

export const getAdminMaintenanceActivitiesSummary = () => {
  return api.get("/admin/maintenance-activities/summary");
};

export const getAdminMaintenanceActivities = (params = {}) => {
  return api.get("/admin/maintenance-activities", { params });
};

export const reviewAdminMaintenanceActivity = (id, payload) => {
  return api.put(`/admin/maintenance-activities/${id}/review`, payload);
};

export const getMyMaintenanceSchedules = () => {
  return api.get("/users/me/maintenance-schedules");
};

export const getMyMaintenanceActivities = () => {
  return api.get("/users/me/maintenance-activities");
};

export const submitMyMaintenanceActivity = (payload) => {
  return api.post("/users/me/maintenance-activities", payload);
};
