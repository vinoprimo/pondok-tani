import { api } from "../api";

export const listSalesOrders = (params = {}) => {
  return api.get("/sales/orders", { params });
};

export const createSalesOrder = (payload) => {
  return api.post("/sales/orders", payload);
};
