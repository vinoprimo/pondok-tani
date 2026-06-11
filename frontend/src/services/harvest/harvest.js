import { api } from "../api";

export const listHarvests = (params = {}) => {
  return api.get("/harvest", { params });
};

export const getHarvestSummary = (params = {}) => {
  return api.get("/harvest/summary", { params });
};

export const listHarvestRequests = (params = {}) => {
  return api.get("/harvest-requests", { params });
};

export const createHarvestRequest = (payload) => {
  return api.post("/harvest-requests", payload);
};

export const validateHarvestRequest = (id, payload) => {
  return api.patch(`/harvest-requests/${id}/validate`, payload);
};
