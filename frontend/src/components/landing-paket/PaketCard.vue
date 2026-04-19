<script setup lang="ts">
import { CheckCircle2, Sprout } from 'lucide-vue-next';

const props = defineProps<{
  name: string;
  range: string;
  roi: string;
  duration: string;
  highlight?: boolean;
  selected?: boolean;
  buttonLabel?: string;
}>();

const emit = defineEmits<{
  choose: [];
}>();
</script>

<template>
  <article
    class="rounded-2xl border p-6 transition-all relative overflow-hidden"
    :class="[
      highlight ? 'border-green-500 bg-green-50 shadow-lg' : 'border-gray-200 bg-white hover:shadow-md',
      selected ? 'ring-2 ring-green-600 ring-offset-2' : '',
    ]"
  >
    <div v-if="selected" class="absolute right-4 top-4 flex items-center gap-1 rounded-full bg-green-700 px-3 py-1 text-xs font-semibold text-white">
      <CheckCircle2 class="h-3.5 w-3.5" />
      Terpilih
    </div>

    <div class="flex items-center justify-between mb-4">
      <h3 class="text-xl font-semibold text-gray-900 pr-16">{{ name }}</h3>
      <Sprout class="w-6 h-6" :class="highlight ? 'text-green-700' : 'text-green-600'" />
    </div>

    <ul class="space-y-2 text-sm text-gray-700 mb-6">
      <li><span class="font-medium">Modal:</span> {{ range }}</li>
      <li><span class="font-medium">Estimasi ROI:</span> {{ roi }}</li>
      <li><span class="font-medium">Durasi:</span> {{ duration }}</li>
    </ul>

    <button
      type="button"
      class="w-full py-2.5 rounded-lg font-medium transition-colors"
      :class="highlight || selected ? 'bg-green-700 text-white hover:bg-green-800' : 'bg-green-600 text-white hover:bg-green-700'"
      @click="emit('choose')"
    >
      {{ props.buttonLabel || 'Pilih Paket' }}
    </button>
  </article>
</template>
