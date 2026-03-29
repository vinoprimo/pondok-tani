<script setup>
import { ref, onMounted } from "vue";
import { getInvestments } from "../services/investment";

const investments = ref([]);

onMounted(async () => {
  try {
    const res = await getInvestments();
    investments.value = res.data;
  } catch (err) {
    console.error(err);
  }
});
</script>

<template>
  <div>
    <h1>Dashboard</h1>

    <div v-for="inv in investments" :key="inv.id">
      <p><b>{{ inv.package_name }}</b></p>
      <p>Rp {{ inv.amount }}</p>
      <p>Status: {{ inv.status }}</p>
      <hr />
    </div>
  </div>
</template>