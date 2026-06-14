import { api } from "../api";

export const listGradingBatches = (params = {}) => {
  return api.get("/grading-batches", { params });
};

export const createGradingBatch = (payload) => {
  return api.post("/grading-batches", payload);
};
