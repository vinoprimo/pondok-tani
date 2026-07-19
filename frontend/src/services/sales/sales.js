import { api } from "../api";

export const listSalesOrders = (params = {}) => {
  return api.get("/sales/orders", { params });
};

export const createSalesOrder = (payload) => {
  return api.post("/sales/orders", payload);
};

export const updateSalesOrder = (id, payload) => {
  return api.put(`/sales/orders/${id}`, payload);
};

export const deleteSalesOrder = (id) => {
  return api.delete(`/sales/orders/${id}`);
};
