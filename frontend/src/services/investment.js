import { api } from "./api";

export const getInvestments = () => {
  return api.get("/investments");
};