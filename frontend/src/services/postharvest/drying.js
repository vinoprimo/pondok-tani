import { api } from "../api";

export const listDryingProcesses = (params = {}) => {
  return api.get("/drying-processes", { params });
};

export const completeDryingProcess = (id, payload) => {
  return api.patch(`/drying-processes/${id}/complete`, payload);
};
