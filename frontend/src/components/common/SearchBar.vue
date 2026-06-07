<template>
  <div class="relative w-full max-w-sm">
    <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
      <Search class="h-5 w-5 text-gray-400" />
    </div>
    <input
      type="text"
      :value="modelValue"
      @input="onInput"
      :placeholder="placeholder"
      class="block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md leading-5 bg-white placeholder-gray-500 focus:outline-none focus:placeholder-gray-400 focus:ring-1 focus:ring-green-500 focus:border-green-500 sm:text-sm"
    />
  </div>
</template>

<script setup>
import { Search } from 'lucide-vue-next'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  placeholder: {
    type: String,
    default: 'Search...'
  }
})

const emit = defineEmits(['update:modelValue', 'search'])

let timeout = null

const onInput = (event) => {
  const value = event.target.value
  emit('update:modelValue', value)
  
  if (timeout) clearTimeout(timeout)
  timeout = setTimeout(() => {
    emit('search', value)
  }, 500) // 500ms debounce
}
</script>
