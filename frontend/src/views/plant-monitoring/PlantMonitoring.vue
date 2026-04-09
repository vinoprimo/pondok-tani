<script setup lang="ts">
import { computed, ref } from 'vue'
import {
  Sprout,
  Calendar,
  Camera,
  TrendingUp,
  Droplets,
  Sun,
  ThermometerSun,
  CheckCircle2,
} from 'lucide-vue-next'

const ERROR_IMG_SRC =
  'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iODgiIGhlaWdodD0iODgiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyIgc3Ryb2tlPSIjMDAwIiBzdHJva2UtbGluZWpvaW49InJvdW5kIiBvcGFjaXR5PSIuMyIgZmlsbD0ibm9uZSIgc3Ryb2tlLXdpZHRoPSIzLjciPjxyZWN0IHg9IjE2IiB5PSIxNiIgd2lkdGg9IjU2IiBoZWlnaHQ9IjU2IiByeD0iNiIvPjxwYXRoIGQ9Im0xNiA1OCAxNi0xOCAzMiAzMiIvPjxjaXJjbGUgY3g9IjUzIiBjeT0iMzUiIHI9IjciLz48L3N2Zz4KCg=='

function onImgError(e: Event) {
  const el = e.target as HTMLImageElement
  if (el) el.src = ERROR_IMG_SRC
}

const selectedPlant = ref('P-001')

const plants = [
  {
    id: 'P-001',
    name: 'Vanilla Plant #001',
    batch: 'Premium Batch A',
    age: '18 months',
    status: 'Flowering',
    health: 95,
    lastInspection: '2026-01-18',
    nextInspection: '2026-01-25',
    location: 'Sector A-12',
  },
  {
    id: 'P-002',
    name: 'Vanilla Plant #002',
    batch: 'Organic Batch B',
    age: '14 months',
    status: 'Growing',
    health: 92,
    lastInspection: '2026-01-19',
    nextInspection: '2026-01-26',
    location: 'Sector B-08',
  },
  {
    id: 'P-003',
    name: 'Vanilla Plant #003',
    batch: 'Premium Batch A',
    age: '22 months',
    status: 'Harvesting',
    health: 98,
    lastInspection: '2026-01-20',
    nextInspection: '2026-01-27',
    location: 'Sector A-15',
  },
]

const timeline = [
  {
    stage: 'Planting',
    date: '2024-07-15',
    status: 'completed' as const,
    description: 'Initial planting and setup',
    images: 1,
  },
  {
    stage: 'Early Growth',
    date: '2024-10-20',
    status: 'completed' as const,
    description: 'First 3 months of development',
    images: 3,
  },
  {
    stage: 'Vegetative Stage',
    date: '2025-02-10',
    status: 'completed' as const,
    description: 'Strong vine development',
    images: 4,
  },
  {
    stage: 'Pre-Flowering',
    date: '2025-08-05',
    status: 'completed' as const,
    description: 'Flower buds forming',
    images: 5,
  },
  {
    stage: 'Flowering',
    date: '2025-11-15',
    status: 'current' as const,
    description: 'Active flowering phase',
    images: 6,
  },
  {
    stage: 'Pod Development',
    date: 'Expected: 2026-03-01',
    status: 'upcoming' as const,
    description: 'Bean pod formation',
    images: 0,
  },
  {
    stage: 'Harvest',
    date: 'Expected: 2026-08-15',
    status: 'upcoming' as const,
    description: 'Ready for harvesting',
    images: 0,
  },
]

const environmentalData = [
  { label: 'Temperature', value: '26°C', status: 'optimal', icon: ThermometerSun },
  { label: 'Humidity', value: '75%', status: 'optimal', icon: Droplets },
  { label: 'Sunlight', value: '6.5 hrs', status: 'good', icon: Sun },
  { label: 'Soil Health', value: '8.2/10', status: 'excellent', icon: Sprout },
]

const currentPlant = computed(() => plants.find((p) => p.id === selectedPlant.value) ?? plants[0])

const mainPlantSrc =
  'https://images.unsplash.com/photo-1763050233345-7945d724d009?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHx2YW5pbGxhJTIwcGxhbnQlMjBmbG93ZXJ8ZW58MXx8fHwxNzY5MDI4MDU4fDA&ixlib=rb-4.1.0&q=80&w=1080&utm_source=figma&utm_medium=referral'

const gallerySrcs = [
  'https://images.unsplash.com/photo-1763050233345-7945d724d009?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHx2YW5pbGxhJTIwcGxhbnQlMjBmbG93ZXJ8ZW58MXx8fHwxNzY5MDI4MDU4fDA&ixlib=rb-4.1.0&q=80&w=400',
  'https://images.unsplash.com/photo-1637922808382-0e5930886159?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHx2YW5pbGxhJTIwcGxhbnRhdGlvbiUyMGZhcm18ZW58MXx8fHwxNzY5MDI4MDU4fDA&ixlib=rb-4.1.0&q=80&w=400',
  'https://images.unsplash.com/photo-1675501343762-fc726e82809c?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w3Nzg4Nzd8MHwxfHNlYXJjaHwxfHx2YW5pbGxhJTIwYmVhbnMlMjBwb2RzfGVufDF8fHx8MTc2OTAyODA1OXww&ixlib=rb-4.1.0&q=80&w=400',
]

function plantButtonClass(id: string) {
  return selectedPlant.value === id
    ? 'border-green-600 bg-green-50'
    : 'border-gray-200 bg-white hover:border-gray-300'
}

function plantTitleClass(id: string) {
  return selectedPlant.value === id ? 'text-green-700' : 'text-gray-900'
}

function statusBadgeClass(status: string) {
  if (status === 'Flowering') return 'bg-purple-100 text-purple-700'
  if (status === 'Growing') return 'bg-blue-100 text-blue-700'
  return 'bg-yellow-100 text-yellow-700'
}

function envStatusClass(status: string) {
  if (status === 'optimal') return 'text-green-600'
  if (status === 'excellent') return 'text-blue-600'
  return 'text-yellow-600'
}

function timelineDotClass(status: string) {
  if (status === 'completed') return 'bg-green-100'
  if (status === 'current') return 'bg-blue-100'
  return 'bg-gray-100'
}
</script>

<template>
  <div class="space-y-6">
    <div>
      <h2 class="text-2xl font-semibold text-gray-900">Plant Monitoring</h2>
      <p class="text-gray-600 mt-1">Track growth stages and health of your vanilla plants</p>
    </div>

    <div class="flex gap-3 overflow-x-auto pb-2">
      <button
        v-for="plant in plants"
        :key="plant.id"
        type="button"
        :class="[
          'flex-shrink-0 px-6 py-3 rounded-lg border-2 transition-all',
          plantButtonClass(plant.id),
        ]"
        @click="selectedPlant = plant.id"
      >
        <div class="text-left">
          <p :class="['font-medium', plantTitleClass(plant.id)]">{{ plant.name }}</p>
          <p class="text-sm text-gray-600 mt-1">{{ plant.batch }}</p>
        </div>
      </button>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <div class="lg:col-span-2 bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-start justify-between mb-6">
          <div>
            <h3 class="text-xl font-semibold text-gray-900">{{ currentPlant.name }}</h3>
            <p class="text-gray-600 mt-1">
              Age: {{ currentPlant.age }} • Location: {{ currentPlant.location }}
            </p>
          </div>
          <span :class="['px-3 py-1 rounded-full text-sm font-medium', statusBadgeClass(currentPlant.status)]">
            {{ currentPlant.status }}
          </span>
        </div>

        <div class="aspect-video bg-gray-100 rounded-lg overflow-hidden mb-6">
          <img
            :src="mainPlantSrc"
            alt="Vanilla Plant"
            class="w-full h-full object-cover"
            @error="onImgError"
          />
        </div>

        <div>
          <h4 class="font-semibold text-gray-900 mb-4">Environmental Conditions</h4>
          <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div v-for="(item, index) in environmentalData" :key="index" class="bg-gray-50 rounded-lg p-4">
              <div class="flex items-center gap-2 mb-2">
                <component :is="item.icon" class="w-4 h-4 text-gray-600" />
                <span class="text-sm text-gray-600">{{ item.label }}</span>
              </div>
              <p class="text-lg font-semibold text-gray-900">{{ item.value }}</p>
              <span :class="['text-xs', envStatusClass(item.status)]">{{ item.status }}</span>
            </div>
          </div>
        </div>
      </div>

      <div class="space-y-6">
        <div class="bg-white rounded-xl border border-gray-200 p-6">
          <h4 class="font-semibold text-gray-900 mb-4">Plant Health</h4>
          <div class="text-center mb-4">
            <div class="inline-flex items-center justify-center w-24 h-24 bg-green-50 rounded-full mb-3">
              <span class="text-3xl font-semibold text-green-600">{{ currentPlant.health }}%</span>
            </div>
            <p class="text-sm text-gray-600">Overall Health Score</p>
          </div>
          <div class="space-y-3">
            <div class="flex items-center justify-between text-sm">
              <span class="text-gray-600">Growth Rate</span>
              <span class="font-medium text-gray-900">Excellent</span>
            </div>
            <div class="flex items-center justify-between text-sm">
              <span class="text-gray-600">Disease Risk</span>
              <span class="font-medium text-green-600">Low</span>
            </div>
            <div class="flex items-center justify-between text-sm">
              <span class="text-gray-600">Pest Activity</span>
              <span class="font-medium text-green-600">None</span>
            </div>
          </div>
        </div>

        <div class="bg-white rounded-xl border border-gray-200 p-6">
          <h4 class="font-semibold text-gray-900 mb-4">Inspections</h4>
          <div class="space-y-3">
            <div>
              <p class="text-sm text-gray-600">Last Inspection</p>
              <p class="font-medium text-gray-900">{{ currentPlant.lastInspection }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-600">Next Scheduled</p>
              <p class="font-medium text-gray-900">{{ currentPlant.nextInspection }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="bg-white rounded-xl border border-gray-200 p-6">
      <div class="flex items-center justify-between mb-6">
        <h3 class="text-lg font-semibold text-gray-900">Growth Timeline</h3>
        <button type="button" class="text-sm text-green-600 hover:text-green-700 font-medium">
          View All Photos
        </button>
      </div>

      <div class="relative">
        <div class="absolute left-6 top-0 bottom-0 w-0.5 bg-gray-200" />

        <div class="space-y-6">
          <div v-for="(item, index) in timeline" :key="index" class="relative flex gap-6">
            <div
              :class="[
                'relative z-10 flex-shrink-0 w-12 h-12 rounded-full flex items-center justify-center',
                timelineDotClass(item.status),
              ]"
            >
              <CheckCircle2 v-if="item.status === 'completed'" class="w-6 h-6 text-green-600" />
              <TrendingUp v-else-if="item.status === 'current'" class="w-6 h-6 text-blue-600" />
              <Calendar v-else class="w-6 h-6 text-gray-400" />
            </div>

            <div class="flex-1 pb-6">
              <div class="bg-gray-50 rounded-lg p-4">
                <div class="flex items-start justify-between mb-2">
                  <div>
                    <h4 class="font-semibold text-gray-900">{{ item.stage }}</h4>
                    <p class="text-sm text-gray-600 mt-1">{{ item.date }}</p>
                  </div>
                  <div v-if="item.images > 0" class="flex items-center gap-1 text-sm text-gray-600">
                    <Camera class="w-4 h-4" />
                    <span>{{ item.images }} photos</span>
                  </div>
                </div>
                <p class="text-sm text-gray-600">{{ item.description }}</p>

                <div v-if="item.status === 'current'" class="grid grid-cols-3 gap-2 mt-4">
                  <div
                    v-for="(src, gi) in gallerySrcs"
                    :key="gi"
                    class="aspect-square bg-white rounded-lg overflow-hidden"
                  >
                    <img
                      :src="src"
                      alt="Plant photo"
                      class="w-full h-full object-cover"
                      @error="onImgError"
                    />
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
