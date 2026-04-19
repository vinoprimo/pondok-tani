import { api } from "../api";

export const getGrades = () => {
  return api.get("/admin/vanili/grades");
};

export const createGrade = (payload) => {
  return api.post("/admin/vanili/grades", payload);
};

export const updateGrade = (id, payload) => {
  return api.put(`/admin/vanili/grades/${id}`, payload);
};

export const deleteGrade = (id) => {
  return api.delete(`/admin/vanili/grades/${id}`);
};

export const getNationalPrices = () => {
  return api.get("/admin/vanili/prices");
};

export const createNationalPrice = (payload) => {
  return api.post("/admin/vanili/prices", payload);
};

export const updateNationalPrice = (id, payload) => {
  return api.put(`/admin/vanili/prices/${id}`, payload);
};

export const deleteNationalPrice = (id) => {
  return api.delete(`/admin/vanili/prices/${id}`);
};
