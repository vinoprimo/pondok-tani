<script setup lang="ts">
defineProps<{
  title: string;
  description: string;
}>();
</script>

<template>
  <section class="max-w-7xl mx-auto px-6 py-16 lg:py-20">
    <div class="bg-gradient-to-r from-green-600 to-green-700 rounded-3xl p-8 lg:p-12 text-white shadow-xl">
      <p class="inline-block px-3 py-1 rounded-full bg-white/20 text-sm font-medium mb-4">
        Paket Investasi
      </p>
      <h1 class="text-3xl lg:text-5xl font-bold leading-tight mb-4">{{ title }}</h1>
      <p class="text-green-50 text-lg max-w-3xl leading-relaxed">{{ description }}</p>
    </div>
  </section>
</template>
