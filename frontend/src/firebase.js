// Import the functions you need from the Firebase SDKs you need
import { initializeApp } from "firebase/app";
import { getAnalytics } from "firebase/analytics";
import { getFirestore } from "firebase/firestore";

// Your web app's Firebase configuration
const firebaseConfig = {
  apiKey: "AIzaSyChmGelULiHf-kqVVlUZDk30h0Qz2VCpxg",
  authDomain: "pondok-tani.firebaseapp.com",
  projectId: "pondok-tani",
  storageBucket: "pondok-tani.firebasestorage.app",
  messagingSenderId: "157374030333",
  appId: "1:157374030333:web:148767688eba6309c1f994",
  measurementId: "G-PSP2ZHYLLS",
};

// Initialize Firebase
export const firebaseApp = initializeApp(firebaseConfig);
export const firestore = getFirestore(firebaseApp);
export const analytics = typeof window !== "undefined" ? getAnalytics(firebaseApp) : null;
