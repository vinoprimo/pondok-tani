import { api } from "../api";

export const listRevenueSimulations = (params = {}) => {
  return api.get("/financial/revenue-simulations", { params });
};

export const createRevenueSimulation = (payload) => {
  return api.post("/financial/revenue-simulations", payload);
};

export const updateRevenueSimulation = (id, payload) => {
  return api.put(`/financial/revenue-simulations/${id}`, payload);
};

export const deleteRevenueSimulation = (id) => {
  return api.delete(`/financial/revenue-simulations/${id}`);
};
