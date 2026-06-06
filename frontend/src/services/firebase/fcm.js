import { getMessaging, getToken, onMessage } from 'firebase/messaging';
import { firebaseApp } from '../../firebase';
import { updateMyFCMToken } from '../user/user';

const messaging = getMessaging(firebaseApp);
const vapidKey = import.meta.env.VITE_FIREBASE_VAPID_KEY || '';

export async function registerFCMServiceWorker() {
  if (typeof navigator === 'undefined' || !('serviceWorker' in navigator)) {
    return null;
  }

  try {
    const registration = await navigator.serviceWorker.register('/firebase-messaging-sw.js');
    return registration;
  } catch (error) {
    console.error('Failed to register FCM service worker:', error);
    return null;
  }
}

export async function requestNotificationPermission() {
  if (typeof Notification === 'undefined') {
    throw new Error('Browser tidak mendukung notifikasi');
  }

  const permission = await Notification.requestPermission();
  return permission === 'granted';
}

export async function getFCMToken() {
  if (!vapidKey) {
    throw new Error('VITE_FIREBASE_VAPID_KEY belum dikonfigurasi');
  }

  return await getToken(messaging, { vapidKey });
}

export async function registerFCMToken() {
  const hasToken = Boolean(localStorage.getItem('token'));
  if (!hasToken) return null;

  try {
    await registerFCMServiceWorker();

    const granted = await requestNotificationPermission();
    if (!granted) return null;

    const fcmToken = await getFCMToken();
    if (!fcmToken) return null;

    await updateMyFCMToken({ fcm_token: fcmToken });
    return fcmToken;
  } catch (error) {
    console.error('Failed to register FCM token:', error);
    return null;
  }
}

export function onFCMMessage(callback) {
  return onMessage(messaging, callback);
}
