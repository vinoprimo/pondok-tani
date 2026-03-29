<script setup>
import HelloWorld from './components/HelloWorld.vue'
import { api } from "./services/api";

api.get("/").then(res => {
  console.log(res.data);
});
</script>

<template>
  <router-view />
</template>
