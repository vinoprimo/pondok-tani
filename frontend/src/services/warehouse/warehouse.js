import { api } from "../api";

export const getWarehouseSummary = (params = {}) => {
  return api.get("/warehouse/summary", { params });
};

export const listWarehouseStocks = (params = {}) => {
  return api.get("/warehouse/stocks", { params });
};

export const listStockMovements = (params = {}) => {
  return api.get("/warehouse/movements", { params });
};
