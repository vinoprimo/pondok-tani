<script setup>
import { onMounted } from "vue";
import { getCurrentUser } from "../services/user";
import { user } from "../store/user";

onMounted(async () => {
  try {
    const res = await getCurrentUser();
    user.value = res.data;
  } catch (err) {
    console.error(err);
  }
});
</script>

<template>
  <div class="bg-white border-b px-6 py-4 flex justify-between items-center">

    <input
      type="text"
      placeholder="Cari..."
      class="border px-4 py-2 rounded-lg w-80 outline-none"
    />

    <div class="flex items-center gap-4">
      <span class="text-sm text-gray-600">
        Hi, {{ user?.name || "User" }}
      </span>

      <div class="w-10 h-10 bg-gray-300 rounded-full"></div>
    </div>

  </div>
</template>