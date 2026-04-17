import { api } from "../api";

export const getInvestmentPackages = (params) => {
  return api.get("/investment-packages", { params });
};

export const getInvestmentPackageById = (id) => {
  return api.get(`/investment-packages/${id}`);
};

export const createInvestmentPackage = (payload) => {
  return api.post("/investment-packages", payload);
};

export const updateInvestmentPackage = (id, payload) => {
  return api.put(`/investment-packages/${id}`, payload);
};

export const deleteInvestmentPackage = (id) => {
  return api.delete(`/investment-packages/${id}`);
};
