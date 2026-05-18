import { collection, addDoc, serverTimestamp } from "firebase/firestore";
import { firestore } from "../../firebase";

const emailCollection = collection(firestore, "email_notifications");

export const sendEmailNotification = async ({ to, subject, body, template = null, metadata = {} }) => {
  const docRef = await addDoc(emailCollection, {
    to,
    subject,
    body,
    template,
    metadata,
    status: "pending",
    createdAt: serverTimestamp(),
  });

  return docRef.id;
};

export const logEmailNotification = async ({ to, subject, body, status = "pending", metadata = {} }) => {
  const docRef = await addDoc(emailCollection, {
    to,
    subject,
    body,
    status,
    metadata,
    createdAt: serverTimestamp(),
  });

  return docRef.id;
};
