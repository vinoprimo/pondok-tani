import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import router from "./router";
import { startSessionTimeout } from "./utils/session";

createApp(App).use(router).mount("#app");

startSessionTimeout(() => {
	if (router.currentRoute.value.name !== "login") {
		router.push({ name: "login", query: { expired: "1" } });
	}
});