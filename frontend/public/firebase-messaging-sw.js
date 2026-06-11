importScripts('https://www.gstatic.com/firebasejs/12.13.0/firebase-app-compat.js');
importScripts('https://www.gstatic.com/firebasejs/12.13.0/firebase-messaging-compat.js');

const firebaseConfig = {
  apiKey: 'AIzaSyChmGelULiHf-kqVVlUZDk30h0Qz2VCpxg',
  authDomain: 'pondok-tani.firebaseapp.com',
  projectId: 'pondok-tani',
  storageBucket: 'pondok-tani.firebasestorage.app',
  messagingSenderId: '157374030333',
  appId: '1:157374030333:web:148767688eba6309c1f994',
};

firebase.initializeApp(firebaseConfig);
const messaging = firebase.messaging();

messaging.onBackgroundMessage(function (payload) {
  const notificationTitle = payload.notification?.title || 'Pondok Tani';
  const notificationOptions = {
    body: payload.notification?.body || '',
    icon: '/favicon.ico',
  };

  self.registration.showNotification(notificationTitle, notificationOptions);
});
