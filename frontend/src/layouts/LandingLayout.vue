<script setup lang="ts">
import Navbar from '../components/landing-layouts/Navbar.vue';
import Footer from '../components/landing-layouts/Footer.vue';
</script>

<template>
  <div class="min-h-screen bg-gradient-to-br from-green-50 via-white to-green-50 flex flex-col">
    <Navbar />
    <main class="flex-1">
      <slot />
    </main>
    <Footer />
  </div>
</template>
