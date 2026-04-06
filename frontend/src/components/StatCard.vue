<template>
    <div class="bg-white p-5 rounded-xl shadow-sm">
      <p class="text-sm text-gray-500">{{ title }}</p>
      <h2 class="text-2xl font-bold mt-2">{{ value }}</h2>
      <p class="text-green-600 text-sm mt-1">{{ change }}</p>
    </div>
  </template>
  
  <script setup>
  defineProps({
    title: String,
    value: String,
    change: String,
  });
  </script>