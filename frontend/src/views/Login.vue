<script setup>
import { ref } from "vue";
import { login } from "../services/auth";
import { useRouter } from "vue-router";

const router = useRouter();

const email = ref("");
const password = ref("");
const error = ref("");

const handleLogin = async () => {
  try {
    const res = await login({
      email: email.value,
      password: password.value,
    });

    // simpan token
    localStorage.setItem("token", res.data.token);

    router.push("/dashboard");
  } catch (err) {
    error.value = "Login gagal";
    console.error(err);
  }
};
</script>

<template>
  <div>
    <h2>Login</h2>

    <input v-model="email" placeholder="Email" />
    <input v-model="password" type="password" placeholder="Password" />

    <button @click="handleLogin">Login</button>

    <p v-if="error">{{ error }}</p>
  </div>
</template>