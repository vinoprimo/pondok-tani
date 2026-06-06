import { api } from '../api';

export const fetchSystemNotifications = async (userId) => {
  try {
    const response = await api.get('/users/me/notifications');
    return response.data || [];
  } catch (error) {
    console.error('Error fetching notifications:', error);
    return [];
  }
};

export const subscribeSystemNotifications = (userId, callback) => {
  // Fallback to polling for REST API
  fetchSystemNotifications(userId).then(callback);
  const intervalId = setInterval(() => {
    fetchSystemNotifications(userId).then(callback);
  }, 30000); // poll every 30 seconds

  return () => clearInterval(intervalId);
};

export const sendSystemNotification = async ({ recipientId, title, message, category = "general", type = "info" }) => {
  // Notifications are typically created by the backend automatically, 
  // but if needed from frontend, it would go here via an API POST.
  return null; 
};

export const markSystemNotificationRead = async (notificationId) => {
  try {
    await api.put(`/users/me/notifications/${notificationId}/read`);
  } catch (error) {
    console.error('Error marking notification as read:', error);
  }
};

