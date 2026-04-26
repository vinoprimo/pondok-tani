import { api } from "../api";

export const listPlantMonitorings = (params = {}) => {
  return api.get("/plant-monitorings", { params });
};

export const createPlantMonitoring = (payload) => {
  if (payload instanceof FormData) {
    return api.post("/plant-monitorings", payload, {
      headers: {
        "Content-Type": "multipart/form-data",
      },
    });
  }

  return api.post("/plant-monitorings", payload);
};
