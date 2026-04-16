<script setup lang="ts">
defineProps<{
  items: string[];
}>();
</script>

<template>
  <section class="max-w-7xl mx-auto px-6 pb-14 lg:pb-20">
    <div class="bg-white border border-gray-200 rounded-2xl p-6 lg:p-8">
      <h2 class="text-2xl font-semibold text-gray-900 mb-6">Yang Anda Dapatkan</h2>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div
          v-for="(item, index) in items"
          :key="index"
          class="rounded-xl bg-gray-50 border border-gray-100 p-4 text-gray-700"
        >
          {{ item }}
        </div>
      </div>
    </div>
  </section>
</template>
