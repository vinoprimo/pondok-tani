<script setup lang="ts">
defineProps<{
  steps: string[];
}>();
</script>

<template>
  <section class="max-w-7xl mx-auto px-6 pb-14 lg:pb-20">
    <h2 class="text-2xl font-semibold text-gray-900 mb-6">Alur Bergabung Kemitraan</h2>
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <div
        v-for="(step, index) in steps"
        :key="index"
        class="rounded-2xl bg-white border border-gray-200 p-6"
      >
        <div class="w-10 h-10 rounded-full bg-green-600 text-white flex items-center justify-center font-bold mb-4">
          {{ index + 1 }}
        </div>
        <p class="text-gray-700 leading-relaxed">{{ step }}</p>
      </div>
    </div>
  </section>
</template>
