<script setup lang="ts">
defineProps<{
  title: string;
  subtitle: string;
}>();
</script>

<template>
  <section class="max-w-7xl mx-auto px-6 py-16 lg:py-20">
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8 items-center">
      <div>
        <p class="inline-block px-3 py-1 rounded-full bg-green-100 text-green-700 text-sm font-medium mb-4">
          Program Kemitraan
        </p>
        <h1 class="text-3xl lg:text-5xl font-bold text-gray-900 leading-tight mb-4">{{ title }}</h1>
        <p class="text-lg text-gray-600 leading-relaxed">{{ subtitle }}</p>
      </div>
      <div class="bg-gradient-to-br from-green-100 to-white rounded-3xl p-8 border border-green-200">
        <div class="grid grid-cols-2 gap-4">
          <div class="rounded-xl bg-white p-4 border border-gray-100">
            <p class="text-sm text-gray-500">Mitra aktif</p>
            <p class="text-2xl font-bold text-green-700 mt-1">120+</p>
          </div>
          <div class="rounded-xl bg-white p-4 border border-gray-100">
            <p class="text-sm text-gray-500">Wilayah</p>
            <p class="text-2xl font-bold text-green-700 mt-1">18</p>
          </div>
          <div class="rounded-xl bg-white p-4 border border-gray-100">
            <p class="text-sm text-gray-500">Pelatihan</p>
            <p class="text-2xl font-bold text-green-700 mt-1">Bulanan</p>
          </div>
          <div class="rounded-xl bg-white p-4 border border-gray-100">
            <p class="text-sm text-gray-500">Pendampingan</p>
            <p class="text-2xl font-bold text-green-700 mt-1">1 on 1</p>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
