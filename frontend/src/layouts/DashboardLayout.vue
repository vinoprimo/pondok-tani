<template>
    <div class="flex h-screen bg-gray-50">
      <Sidebar />
  
      <div class="flex-1 flex flex-col">
        <Topbar />
  
        <div class="p-6 overflow-y-auto">
          <slot />
        </div>
      </div>
    </div>
  </template>
  
  <script setup>
  import Sidebar from "../components/Sidebar.vue";
  import Topbar from "../components/Topbar.vue";
  </script>