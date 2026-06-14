<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { Camera, Loader2, NotebookText, PlusCircle, Send, Sprout } from 'lucide-vue-next'
import { createHarvestRequest } from '../../services/harvest/harvest'
import { createPlantMonitoring, listPlantMonitorings } from '../../services/plant-monitoring/monitoring'

const monitorings = ref<any[]>([])
const loading = ref(false)
const saving = ref(false)
const requestingHarvest = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

const selectedBatchId = ref<number | null>(null)
const photoFile = ref<File | null>(null)
const photoPreview = ref('')

const form = ref({
  phase: 'penanaman',
  healthStatus: 'sehat',
  disease: '',
  diseaseNote: '',
  affectedCount: 0,
  totalPlants: 0,
  note: '',
})

const phaseOptions = [
  { value: 'penanaman', label: 'Penanaman' },
  { value: 'pertumbuhan_awal', label: 'Pertumbuhan Awal' },
  { value: 'vegetatif', label: 'Vegetatif' },
  { value: 'pra-berbunga', label: 'Pra-berbunga' },
  { value: 'berbunga', label: 'Berbunga' },
  { value: 'panen', label: 'Panen' },
]

const phaseLabelMap: Record<string, string> = phaseOptions.reduce((acc, item) => {
  acc[item.value] = item.label
  return acc
}, {} as Record<string, string>)

const healthStatusOptions = [
  { value: 'sehat', label: 'Sehat' },
  { value: 'sebagian_terdampak', label: 'Sebagian terdampak' },
  { value: 'mati', label: 'Mati' },
]

const diseaseOptions = [
  { value: '', label: 'Tidak ada' },
  { value: 'batang_busuk', label: 'Batang busuk' },
  { value: 'lainnya', label: 'Lainnya' },
]

const groupedBatches = computed(() => {
  const map = new Map<number, any>()
  monitorings.value.forEach((item) => {
    if (!map.has(item.plant_batch_id)) {
      map.set(item.plant_batch_id, {
        id: item.plant_batch_id,
        batchCode: item.batch_code,
        packageName: item.package_name,
      })
    }
  })
  return Array.from(map.values())
})

const selectedBatch = computed(() => groupedBatches.value.find((item) => item.id === selectedBatchId.value) ?? null)

const selectedBatchMonitorings = computed(() => {
  if (!selectedBatchId.value) return []
  return monitorings.value
    .filter((item) => item.plant_batch_id === selectedBatchId.value)
    .sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
})

const latestMonitoring = computed(() => selectedBatchMonitorings.value[0] ?? null)

const route = useRoute()
const activeSection = computed(() => {
  if (route.name === 'dashboard-plants-history') return 'history'
  return 'input'
})

function formatDate(value: string | Date) {
  const date = value instanceof Date ? value : new Date(value)
  if (Number.isNaN(date.getTime())) return '-'
  return date.toLocaleDateString('id-ID', {
    day: '2-digit',
    month: 'long',
    year: 'numeric',
  })
}

function healthBadgeClass(status: string) {
  if (status === 'sehat') return 'bg-green-100 text-green-700'
  if (status === 'sebagian_terdampak') return 'bg-amber-100 text-amber-700'
  return 'bg-red-100 text-red-700'
}

function healthLabel(status: string) {
  const match = healthStatusOptions.find((item) => item.value === status)
  return match ? match.label : status
}

function diseaseLabel(value: string) {
  const match = diseaseOptions.find((item) => item.value === value)
  return match ? match.label : value || 'Tidak ada'
}

function resetFormDefaults() {
  form.value.phase = 'penanaman'
  form.value.healthStatus = 'sehat'
  form.value.disease = ''
  form.value.diseaseNote = ''
  form.value.affectedCount = 0
  form.value.totalPlants = Number(latestMonitoring.value?.total_plants || 0)
  form.value.note = ''
}

function buildMonitoringPayload() {
  if (!selectedBatchId.value) {
    errorMessage.value = 'Pilih batch terlebih dahulu.'
    return null
  }

  const affectedCount = Number(form.value.affectedCount)
  const totalPlants = Number(form.value.totalPlants)
  if (!Number.isFinite(totalPlants) || totalPlants <= 0) {
    errorMessage.value = 'Total tanaman harus lebih dari 0.'
    return null
  }
  if (!Number.isFinite(affectedCount) || affectedCount < 0) {
    errorMessage.value = 'Jumlah tanaman terdampak tidak valid.'
    return null
  }
  if (affectedCount > totalPlants) {
    errorMessage.value = 'Jumlah tanaman terdampak tidak boleh melebihi total tanaman.'
    return null
  }

  if (!photoFile.value) {
    errorMessage.value = 'Foto monitoring wajib diunggah.'
    return null
  }

  const payload = new FormData()
  payload.append('plant_batch_id', String(selectedBatchId.value))
  payload.append('phase', form.value.phase)
  payload.append('health_status', form.value.healthStatus)
  payload.append('disease', form.value.disease)
  payload.append('disease_note', form.value.disease === 'lainnya' ? form.value.diseaseNote : '')
  payload.append('affected_count', String(Math.floor(affectedCount)))
  payload.append('total_plants', String(Math.floor(totalPlants)))
  payload.append('note', form.value.note.trim())
  payload.append('photo', photoFile.value)

  return payload
}

function buildHarvestMessage(monitoring: any) {
  const batchLabel = selectedBatch.value?.batchCode || '-'
  const packageLabel = selectedBatch.value?.packageName || '-'
  const monitoringDate = formatDate(monitoring?.created_at || monitoring?.monitoring_date || new Date())
  const requestDate = formatDate(new Date())

  const lines = [
    'Halo Admin Pondok Tani,',
    '',
    'Saya ingin mengajukan panen dengan detail monitoring berikut:',
    `Tanggal pengajuan: ${requestDate}`,
    `Batch: ${batchLabel}`,
    `Paket: ${packageLabel}`,
    `Fase: ${phaseLabelMap[monitoring?.phase] || monitoring?.phase || '-'}`,
    `Status kesehatan: ${healthLabel(monitoring?.health_status || monitoring?.healthStatus || '-')}`,
    `Disease: ${diseaseLabel(monitoring?.disease || '')}`,
    `Tanaman terdampak: ${monitoring?.affected_count ?? monitoring?.affectedCount ?? 0} / ${monitoring?.total_plants ?? monitoring?.totalPlants ?? 0}`,
  ]

  if (monitoring?.disease_note || monitoring?.diseaseNote) {
    lines.push(`Disease note: ${monitoring?.disease_note || monitoring?.diseaseNote}`)
  }

  lines.push(`Catatan: ${monitoring?.note || '-'}`)
  lines.push(`Tanggal monitoring: ${monitoringDate}`)
  lines.push('', 'Terima kasih.')

  return lines.join('\n')
}

function openWhatsApp(message: string) {
  const phoneNumber = '6281328164003'
  const url = `https://wa.me/${phoneNumber}?text=${encodeURIComponent(message)}`
  window.open(url, '_blank')
}

function handlePhotoChange(event: Event) {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0] ?? null

  if (photoPreview.value) {
    URL.revokeObjectURL(photoPreview.value)
  }

  photoFile.value = file
  photoPreview.value = file ? URL.createObjectURL(file) : ''
}

async function fetchMonitorings() {
  loading.value = true
  errorMessage.value = ''
  try {
    const response = await listPlantMonitorings()
    monitorings.value = Array.isArray(response.data) ? response.data : []

    if (!selectedBatchId.value && groupedBatches.value.length) {
      selectedBatchId.value = groupedBatches.value[0].id
    }

    if (!selectedBatchId.value) {
      resetFormDefaults()
      return
    }

    const selectedStillExists = groupedBatches.value.some((item) => item.id === selectedBatchId.value)
    if (!selectedStillExists) {
      selectedBatchId.value = groupedBatches.value[0]?.id ?? null
    }

    resetFormDefaults()
  } catch (error) {
    console.error('Failed to load plant monitorings:', error)
    errorMessage.value = 'Gagal memuat data monitoring tanaman.'
  } finally {
    loading.value = false
  }
}

async function submitMonitoring() {
  successMessage.value = ''
  errorMessage.value = ''

  const payload = buildMonitoringPayload()
  if (!payload) return

  saving.value = true
  try {
    await createPlantMonitoring(payload)
    successMessage.value = 'Monitoring tanaman berhasil disimpan.'
    if (photoPreview.value) {
      URL.revokeObjectURL(photoPreview.value)
    }
    photoFile.value = null
    photoPreview.value = ''
    resetFormDefaults()
    await fetchMonitorings()
  } catch (error) {
    console.error('Failed to create monitoring:', error)
    errorMessage.value = 'Gagal menyimpan monitoring tanaman.'
  } finally {
    saving.value = false
  }
}

async function submitHarvestRequest() {
  successMessage.value = ''
  errorMessage.value = ''

  const payload = buildMonitoringPayload()
  if (!payload) return

  requestingHarvest.value = true
  try {
    const response = await createPlantMonitoring(payload)
    const monitoringData = response?.data?.data
    if (!monitoringData?.id) {
      throw new Error('Monitoring response missing')
    }

    await createHarvestRequest({ plant_monitoring_id: monitoringData.id })
    successMessage.value = 'Ajuan panen berhasil dibuat.'

    if (photoPreview.value) {
      URL.revokeObjectURL(photoPreview.value)
    }
    photoFile.value = null
    photoPreview.value = ''
    resetFormDefaults()
    await fetchMonitorings()

    openWhatsApp(buildHarvestMessage(monitoringData))
  } catch (error: any) {
    console.error('Failed to create harvest request:', error)
    errorMessage.value = error?.response?.data?.error || 'Gagal mengajukan panen.'
  } finally {
    requestingHarvest.value = false
  }
}

onMounted(() => {
  fetchMonitorings()
})
</script>

<template>
  <div class="space-y-6">
    <div>
      <h2 class="text-2xl font-semibold text-gray-900">
        {{ activeSection === 'input' ? 'Input Monitoring' : 'Riwayat Monitoring' }}
      </h2>
      <p class="mt-1 text-gray-600">
        {{ activeSection === 'input'
          ? 'Isi data monitoring tanaman untuk batch yang dipilih.'
          : 'Lihat riwayat monitoring yang telah tercatat setiap batch.' }}
      </p>
    </div>

    <div v-if="errorMessage" class="rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
      {{ errorMessage }}
    </div>
    <div v-if="successMessage" class="rounded-xl border border-green-200 bg-green-50 px-4 py-3 text-sm text-green-700">
      {{ successMessage }}
    </div>

    <div class="rounded-xl border border-gray-200 bg-white p-4">
      <p class="mb-3 text-sm font-medium text-gray-700">Pilih batch</p>
      <div v-if="loading" class="flex items-center gap-2 text-sm text-gray-500">
        <Loader2 class="h-4 w-4 animate-spin" />
        Memuat batch monitoring...
      </div>
      <div v-else-if="!groupedBatches.length" class="rounded-lg border border-dashed border-gray-300 p-4 text-sm text-gray-500">
        Belum ada data monitoring. Data awal akan otomatis terbentuk setelah validasi penanaman awal.
      </div>
      <div v-else class="flex gap-3 overflow-x-auto pb-1">
        <button
          v-for="item in groupedBatches"
          :key="item.id"
          type="button"
          :class="[
            'min-w-[230px] rounded-lg border px-4 py-3 text-left transition-colors',
            selectedBatchId === item.id ? 'border-green-500 bg-green-50' : 'border-gray-200 bg-white hover:border-gray-300',
          ]"
          @click="selectedBatchId = item.id; resetFormDefaults()"
        >
          <p class="font-semibold text-gray-900">{{ item.batchCode }}</p>
          <p class="mt-1 text-sm text-gray-600">{{ item.packageName || '-' }}</p>
        </button>
      </div>
    </div>

    <div v-if="selectedBatch" class="grid grid-cols-1 gap-6">
      <section
        v-if="activeSection === 'input'"
        class="rounded-xl border border-gray-200 bg-white p-5"
      >
        <div class="mb-4 flex items-center gap-2">
          <PlusCircle class="h-4 w-4 text-green-600" />
          <h3 class="text-lg font-semibold text-gray-900">Input monitoring</h3>
        </div>

        <div class="space-y-4">
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700">Batch</label>
            <input
              type="text"
              :value="`${selectedBatch.batchCode} • ${selectedBatch.packageName || '-'}`"
              class="w-full rounded-lg border border-gray-200 bg-gray-50 px-3 py-2 text-sm text-gray-600"
              disabled
            />
          </div>

          <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
            <div>
              <label class="mb-1 block text-sm font-medium text-gray-700">Fase</label>
              <select v-model="form.phase" class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-green-500 focus:outline-none focus:ring-2 focus:ring-green-100">
                <option v-for="phase in phaseOptions" :key="phase.value" :value="phase.value">{{ phase.label }}</option>
              </select>
            </div>

            <div>
              <label class="mb-1 block text-sm font-medium text-gray-700">Status kesehatan</label>
              <select v-model="form.healthStatus" class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-green-500 focus:outline-none focus:ring-2 focus:ring-green-100">
                <option v-for="item in healthStatusOptions" :key="item.value" :value="item.value">{{ item.label }}</option>
              </select>
            </div>

            <div>
              <label class="mb-1 block text-sm font-medium text-gray-700">Disease</label>
              <select v-model="form.disease" class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-green-500 focus:outline-none focus:ring-2 focus:ring-green-100">
                <option v-for="item in diseaseOptions" :key="item.value || 'none'" :value="item.value">{{ item.label }}</option>
              </select>
            </div>

            <div>
              <label class="mb-1 block text-sm font-medium text-gray-700">Tanaman terdampak</label>
              <input
                v-model.number="form.affectedCount"
                type="number"
                min="0"
                class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-green-500 focus:outline-none focus:ring-2 focus:ring-green-100"
              />
            </div>

            <div class="md:col-span-2">
              <label class="mb-1 block text-sm font-medium text-gray-700">Total tanaman</label>
              <input
                v-model.number="form.totalPlants"
                type="number"
                min="1"
                class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-green-500 focus:outline-none focus:ring-2 focus:ring-green-100"
              />
            </div>

            <div v-if="form.disease === 'lainnya'" class="md:col-span-2">
              <label class="mb-1 block text-sm font-medium text-gray-700">Disease note</label>
              <textarea
                v-model="form.diseaseNote"
                rows="2"
                class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-green-500 focus:outline-none focus:ring-2 focus:ring-green-100"
                placeholder="Deskripsi penyakit lainnya"
              />
            </div>

            <div class="md:col-span-2">
              <label class="mb-1 block text-sm font-medium text-gray-700">Catatan monitoring</label>
              <textarea
                v-model="form.note"
                rows="3"
                class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-green-500 focus:outline-none focus:ring-2 focus:ring-green-100"
                placeholder="Catatan tambahan"
              />
            </div>

            <div class="md:col-span-2">
              <label class="mb-1 block text-sm font-medium text-gray-700">Foto monitoring</label>
              <input
                type="file"
                accept="image/png,image/jpeg,image/jpg,image/webp"
                class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm file:mr-3 file:rounded-md file:border-0 file:bg-green-50 file:px-3 file:py-1.5 file:text-green-700 hover:file:bg-green-100"
                @change="handlePhotoChange"
              />
              <div v-if="photoPreview" class="mt-3 overflow-hidden rounded-lg border border-gray-200">
                <img :src="photoPreview" alt="Preview monitoring" class="h-40 w-full object-cover" />
              </div>
            </div>
          </div>

          <div class="flex flex-wrap items-center gap-3">
            <button
              type="button"
              class="inline-flex items-center gap-2 rounded-lg bg-green-600 px-4 py-2 text-sm font-semibold text-white hover:bg-green-700 disabled:opacity-60"
              :disabled="saving || requestingHarvest"
              @click="submitMonitoring"
            >
              <Loader2 v-if="saving" class="h-4 w-4 animate-spin" />
              <Camera v-else class="h-4 w-4" />
              {{ saving ? 'Menyimpan...' : 'Simpan Monitoring' }}
            </button>

            <button
              v-if="form.phase === 'panen'"
              type="button"
              class="inline-flex items-center gap-2 rounded-lg border border-green-600 px-4 py-2 text-sm font-semibold text-green-700 hover:bg-green-50 disabled:opacity-60"
              :disabled="saving || requestingHarvest"
              @click="submitHarvestRequest"
            >
              <Loader2 v-if="requestingHarvest" class="h-4 w-4 animate-spin" />
              <Send v-else class="h-4 w-4" />
              {{ requestingHarvest ? 'Mengajukan...' : 'Ajukan Panen' }}
            </button>
          </div>
        </div>
      </section>

      <section
        v-if="activeSection === 'history'"
        class="rounded-xl border border-gray-200 bg-white p-5"
      >
        <div class="mb-4 flex items-center justify-between gap-4">
          <div>
            <h3 class="text-lg font-semibold text-gray-900">Riwayat monitoring</h3>
            <p class="mt-1 text-sm text-gray-600">Batch {{ selectedBatch.batchCode }} • Paket {{ selectedBatch.packageName || '-' }}</p>
          </div>

          <div v-if="latestMonitoring" class="rounded-lg bg-green-50 px-3 py-2 text-right">
            <p class="text-xs text-green-700">Progress terbaru</p>
            <p class="text-xl font-semibold text-green-700">{{ latestMonitoring.progress }}%</p>
          </div>
        </div>

        <div v-if="!selectedBatchMonitorings.length" class="rounded-lg border border-dashed border-gray-300 p-5 text-sm text-gray-500">
          Belum ada riwayat monitoring pada batch ini.
        </div>

        <div v-else class="space-y-4">
          <article
            v-for="item in selectedBatchMonitorings"
            :key="item.id"
            class="overflow-hidden rounded-xl border border-gray-200"
          >
            <div class="grid grid-cols-1 gap-0 lg:grid-cols-5">
              <div class="lg:col-span-2">
                <img
                  v-if="item.photo_url"
                  :src="`http://localhost:8000${item.photo_url}`"
                  alt="Foto monitoring"
                  class="h-full min-h-[180px] w-full object-cover"
                />
                <div v-else class="flex h-full min-h-[180px] items-center justify-center bg-gray-100 text-sm text-gray-500">
                  Foto tidak tersedia
                </div>
              </div>

              <div class="space-y-3 p-4 lg:col-span-3">
                <div class="flex flex-wrap items-center justify-between gap-2">
                  <div class="inline-flex items-center gap-2">
                    <Sprout class="h-4 w-4 text-green-600" />
                    <span class="font-semibold text-gray-900">{{ phaseLabelMap[item.phase] || item.phase }}</span>
                    <span class="rounded-full bg-blue-100 px-2 py-0.5 text-xs font-medium text-blue-700">{{ item.progress }}%</span>
                  </div>
                  <span :class="['rounded-full px-2.5 py-1 text-xs font-medium', healthBadgeClass(item.health_status)]">
                    {{ healthLabel(item.health_status) }}
                  </span>
                </div>

                <div class="grid grid-cols-1 gap-2 text-sm text-gray-700 md:grid-cols-2">
                  <!-- <p class="md:col-span-2">
                    ID Monitoring:
                    <span class="inline-flex items-center rounded-md bg-gray-100 px-2 py-0.5 font-mono text-xs font-semibold text-gray-600">
                      #{{ item.id }}
                    </span>
                  </p> -->
                  <p>Disease: <span class="font-medium">{{ diseaseLabel(item.disease) }}</span></p>
                  <p>Terdampak: <span class="font-medium">{{ item.affected_count }} / {{ item.total_plants }}</span></p>
                  <p class="md:col-span-2">Tanggal: <span class="font-medium">{{ formatDate(item.created_at) }}</span></p>
                  <p v-if="item.disease_note" class="md:col-span-2">Disease note: <span class="font-medium">{{ item.disease_note }}</span></p>
                  <p v-if="item.note" class="md:col-span-2">Catatan: <span class="font-medium">{{ item.note }}</span></p>
                </div>
              </div>
            </div>
          </article>
        </div>
      </section>
    </div>

    <div v-else-if="!loading" class="rounded-xl border border-dashed border-gray-300 bg-white p-8 text-center text-sm text-gray-500">
      <NotebookText class="mx-auto mb-2 h-5 w-5 text-gray-400" />
      Tidak ada batch monitoring yang bisa dipilih.
    </div>
  </div>
</template>
