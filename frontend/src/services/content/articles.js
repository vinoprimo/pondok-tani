import { api } from '../api';

function toAbsoluteAssetUrl(value) {
  if (!value) return value;
  if (/^https?:\/\//i.test(value)) return value;

  const baseURL = (api.defaults.baseURL || '').replace(/\/$/, '');
  return `${baseURL}${value.startsWith('/') ? value : `/${value}`}`;
}

function normalizeArticle(item) {
  return {
    ...item,
    image_url: toAbsoluteAssetUrl(item.image_url),
  };
}

export const getPublishedArticles = async (params = {}) => {
  const response = await api.get('/articles', { params });
  return {
    ...response,
    data: Array.isArray(response.data)
      ? response.data.map(normalizeArticle)
      : normalizeArticle(response.data),
  };
};

export const getArticleById = async (id) => {
  const response = await api.get(`/articles/${id}`);
  return {
    ...response,
    data: normalizeArticle(response.data),
  };
};
export const createArticle = (payload) => api.post('/admin/articles', payload, {
  headers: { 'Content-Type': 'multipart/form-data' },
});
export const updateArticle = (id, payload) => api.put(`/admin/articles/${id}`, payload, {
  headers: { 'Content-Type': 'multipart/form-data' },
});
export const deleteArticle = (id) => api.delete(`/admin/articles/${id}`);
