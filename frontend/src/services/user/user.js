import { api } from "../api";

export const getCurrentUser = () => {
  return api.get("/auth/me");
};