import {
  collection,
  query,
  where,
  orderBy,
  onSnapshot,
  getDocs,
  addDoc,
  updateDoc,
  doc,
  serverTimestamp,
} from "firebase/firestore";
import { firestore } from "../../firebase";

const notificationsCollection = collection(firestore, "notifications");

export const fetchSystemNotifications = async (userId) => {
  const q = userId
    ? query(notificationsCollection, where("recipientId", "==", userId), orderBy("createdAt", "desc"))
    : query(notificationsCollection, orderBy("createdAt", "desc"));

  const snapshot = await getDocs(q);
  return snapshot.docs.map((doc) => ({ id: doc.id, ...doc.data() }));
};

export const subscribeSystemNotifications = (userId, callback) => {
  const q = userId
    ? query(notificationsCollection, where("recipientId", "==", userId), orderBy("createdAt", "desc"))
    : query(notificationsCollection, orderBy("createdAt", "desc"));

  return onSnapshot(q, (snapshot) => {
    const notifications = snapshot.docs.map((doc) => ({ id: doc.id, ...doc.data() }));
    callback(notifications);
  });
};

export const sendSystemNotification = async ({ recipientId, title, message, category = "general", type = "info" }) => {
  const docRef = await addDoc(notificationsCollection, {
    recipientId: recipientId || null,
    title,
    message,
    category,
    type,
    read: false,
    createdAt: serverTimestamp(),
  });

  return docRef.id;
};

export const markSystemNotificationRead = async (notificationId) => {
  const notificationDoc = doc(firestore, "notifications", notificationId);
  await updateDoc(notificationDoc, { read: true, readAt: serverTimestamp() });
};
