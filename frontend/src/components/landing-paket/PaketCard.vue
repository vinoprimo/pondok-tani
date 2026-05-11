<script setup lang="ts">
import { CheckCircle2, Clock, TrendingUp, BarChart3 } from 'lucide-vue-next';

const props = defineProps<{
  name: string;
  range: string;
  description?: string;
  roi: string;
  duration: string;
  benefits?: string[];
  minQuantity?: number | string;
  highlight?: boolean;
  popular?: boolean;
  selected?: boolean;
  disabled?: boolean;
  buttonLabel?: string;
  infoMode?: 'full' | 'modal-description';
}>();

const emit = defineEmits<{
  choose: [];
}>();

function handleChoose() {
  if (!props.disabled) {
    emit('choose');
  }
}
</script>

<template>
  <article
    class="rounded-2xl p-8 flex flex-col h-full transform transition-all duration-300 cursor-pointer relative overflow-hidden"
    :class="[
      selected
        ? 'bg-gradient-to-br from-green-600 to-green-600 text-white relative overflow-hidden md:scale-105 shadow-2xl z-10 border-4 border-green-400/50 hover:scale-[1.07] hover:shadow-2xl'
        : 'bg-white border-2 border-gray-200 hover:border-green-500 hover:shadow-xl hover:-translate-y-2',
      highlight && !selected ? 'shadow-xl' : '',
    ]"
  >
    <div v-if="popular" class="absolute top-0 right-0 z-20 rounded-bl-lg bg-yellow-400 px-4 py-1.5 text-sm font-bold text-gray-900 shadow-sm">
      PALING POPULER
    </div>

    <div class="flex items-center justify-between mb-6" :class="{ 'mt-2': popular }">
      <h3 class="text-2xl font-bold" :class="selected ? 'text-white' : 'text-gray-900'">{{ name }}</h3>
      <div class="w-12 h-12 rounded-full flex items-center justify-center bg-green-100">
        <CheckCircle2 class="w-6 h-6 text-green-600" />
      </div>
    </div>

    <div class="mb-6">
      <div class="text-sm mb-2" :class="selected ? 'text-white/80' : 'text-gray-600'">Mulai dari</div>
      <div class="text-4xl font-bold mb-1" :class="selected ? 'text-white' : 'text-gray-900'">{{ range }}</div>
      <div class="text-sm" :class="selected ? 'text-white/80' : 'text-gray-600'">per 1 unit</div>
    </div>

    <div :class="selected ? 'bg-white/10 backdrop-blur-sm border border-white/20' : 'bg-green-50'" class="p-4 rounded-lg mb-6 space-y-3">
      <div class="flex items-center gap-2" :class="selected ? 'text-white' : 'text-green-600'">
        <Clock class="w-5 h-5" :class="selected ? 'text-white' : 'text-green-600'" />
        <span class="font-semibold">Durasi: {{ duration }}</span>
      </div>
      <div class="flex items-center gap-2" :class="selected ? 'text-white' : 'text-green-600'">
        <TrendingUp class="w-5 h-5" :class="selected ? 'text-white' : 'text-green-600'" />
        <span class="font-semibold">Estimasi ROI: {{ roi }}</span>
      </div>
      <div class="flex items-center gap-2" :class="selected ? 'text-white' : 'text-green-600'">
        <BarChart3 class="w-5 h-5" :class="selected ? 'text-white' : 'text-green-600'" />
        <span class="font-semibold">Min. Investasi: {{ props.minQuantity ?? '-' }}</span>
      </div>
    </div>

    <div class="border-t pt-6 mb-6 flex-grow" :class="selected ? 'border-white/20' : 'border-gray-200'">
      <div v-if="props.benefits?.length" class="space-y-3">
        <h5 class="text-sm font-semibold" :class="selected ? 'text-white' : 'text-gray-900'">Keuntungan Paket:</h5>
        <ul class="space-y-2 text-sm">
          <li
            v-for="benefit in props.benefits"
            :key="benefit"
            class="flex items-start gap-3"
          >
            <CheckCircle2
              class="mt-0.5 h-5 w-5 flex-shrink-0"
              :class="selected ? 'text-white' : 'text-green-600'"
            />
            <span :class="selected ? 'text-white/90' : 'text-gray-700'">{{ benefit }}</span>
          </li>
        </ul>
      </div>
    </div>

    <button
      type="button"
      class="w-full py-3.5 rounded-lg font-bold transition-colors mt-auto disabled:cursor-not-allowed disabled:opacity-70"
      :class="selected ? 'bg-white text-green-600 hover:bg-green-50 shadow-lg' : 'bg-green-600 text-white hover:bg-green-700'"
      :disabled="props.disabled"
      @click="handleChoose"
    >
      {{ props.buttonLabel || 'Pilih Paket' }}
    </button>
  </article>
</template>
