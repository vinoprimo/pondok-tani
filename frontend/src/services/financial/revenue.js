import { api } from "../api";

export const listActualRevenues = (params = {}) => {
  return api.get("/financial/revenues", { params });
};
