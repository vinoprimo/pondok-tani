import { api } from "../api";

export const login = (data) => {
  return api.post("/auth/login", data);
};

export const register = (data) => {
  return api.post("/auth/register", data);
};

export const registerMitra = (payload) => {
  return api.post("/auth/register/mitra", payload, {
    headers: {
      "Content-Type": "multipart/form-data",
    },
  });
};
