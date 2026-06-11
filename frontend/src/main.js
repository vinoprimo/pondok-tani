import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import router from "./router";
import { startSessionTimeout } from "./utils/session";
import "./firebase";
import { registerFCMServiceWorker, registerFCMToken } from "./services/firebase/fcm";

createApp(App).use(router).mount("#app");

registerFCMServiceWorker()
  .then(() => registerFCMToken())
  .catch((err) => console.error("FCM initialization failed:", err));

startSessionTimeout(() => {
	if (router.currentRoute.value.name !== "login") {
		router.push({ name: "login", query: { expired: "1" } });
	}
});