import { api } from "../api";

export const listOperationalCosts = (params = {}) => {
  return api.get("/financial/operational-costs", { params });
};

export const createOperationalCost = (payload) => {
  return api.post("/financial/operational-costs", payload);
};

export const updateOperationalCost = (id, payload) => {
  return api.put(`/financial/operational-costs/${id}`, payload);
};

export const deleteOperationalCost = (id) => {
  return api.delete(`/financial/operational-costs/${id}`);
};
