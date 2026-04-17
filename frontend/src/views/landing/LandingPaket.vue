<script setup lang="ts">
import { onMounted, ref } from 'vue';
import LandingLayout from '../../layouts/LandingLayout.vue';
import PaketHero from '../../components/landing-paket/PaketHero.vue';
import PaketCard from '../../components/landing-paket/PaketCard.vue';
import PaketBenefitList from '../../components/landing-paket/PaketBenefitList.vue';
import { getInvestmentPackages } from '../../services/investment/package';

type PaketItem = {
  id: number;
  name: string;
  range: string;
  roi: string;
  duration: string;
  highlight?: boolean;
};

const paketItems = ref<PaketItem[]>([]);
const loadingPackages = ref(false);

const formatRupiah = (value: number) =>
  new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    maximumFractionDigits: 0,
  }).format(value);

async function loadPackages() {
  loadingPackages.value = true;
  try {
    const res = await getInvestmentPackages({ status: 'active' });
    const list = Array.isArray(res.data) ? res.data : [];
    paketItems.value = list.map((item: any, idx: number) => ({
      id: item.id,
      name: item.package_name,
      range: `${formatRupiah(item.price)} (min ${item.min_quantity} unit)`,
      roi: 'Estimasi placeholder 15% - 25% / tahun',
      duration: 'Placeholder 12 - 24 bulan',
      highlight: idx === 1,
    }));
  } catch (err) {
    console.error('Failed to load investment packages:', err);
    paketItems.value = [];
  } finally {
    loadingPackages.value = false;
  }
}

const benefits = [
  'Akses dashboard perkembangan kebun secara real-time.',
  'Laporan bulanan performa kebun dan estimasi hasil panen.',
  'Notifikasi milestone penting pada fase budidaya.',
  'Pendampingan tim agronomi selama periode investasi.',
  'Dokumentasi aktivitas lapangan dengan foto dan catatan.',
  'Ringkasan finansial dan proyeksi ROI berkala.',
];

onMounted(() => {
  loadPackages();
});
</script>

<template>
  <LandingLayout>

    <PaketHero
      title="Temukan Paket Investasi Sesuai Target Anda"
      description="Berikut simulasi paket investasi placeholder untuk memudahkan eksplorasi skema pendanaan di Omah Vanili. Konten ini nantinya akan diganti dengan data resmi."
    />

    <section class="max-w-7xl mx-auto px-6 pb-10 lg:pb-14">
      <p v-if="loadingPackages" class="mb-4 text-sm text-gray-500">Memuat paket investasi terbaru...</p>
      <p v-else-if="paketItems.length === 0" class="mb-4 text-sm text-gray-500">
        Belum ada paket aktif. Konten akan diperbarui setelah data tersedia.
      </p>
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <PaketCard
          v-for="item in paketItems"
          :key="item.id"
          :name="item.name"
          :range="item.range"
          :roi="item.roi"
          :duration="item.duration"
          :highlight="item.highlight"
        />
      </div>
    </section>

    <PaketBenefitList :items="benefits" />
  </LandingLayout>
</template>
