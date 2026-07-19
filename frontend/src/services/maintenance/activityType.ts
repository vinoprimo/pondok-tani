import { api } from "../api";

export type ActivityType = {
  id: number;
  name: string;
  description: string;
  created_at: string;
  updated_at: string;
};

export async function getActivityTypes(params?: { page?: number; limit?: number; search?: string }) {
  return await api.get('/admin/activity-types', { params });
}

export async function createActivityType(data: { name: string; description: string }) {
  return await api.post('/admin/activity-types', data);
}

export async function updateActivityType(id: number, data: { name: string; description: string }) {
  return await api.put(`/admin/activity-types/${id}`, data);
}

export async function deleteActivityType(id: number) {
  return await api.delete(`/admin/activity-types/${id}`);
}
