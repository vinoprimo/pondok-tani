<script setup lang="ts">
import LandingLayout from '../../layouts/LandingLayout.vue';
import PaketHero from '../../components/landing-paket/PaketHero.vue';
import PaketCard from '../../components/landing-paket/PaketCard.vue';
import PaketBenefitList from '../../components/landing-paket/PaketBenefitList.vue';

const paketItems = [
  { name: 'Starter Vanili', range: 'Rp 5 - 15 Juta', roi: '14% - 18% / tahun', duration: '12 bulan' },
  { name: 'Growth Vanili', range: 'Rp 20 - 50 Juta', roi: '18% - 24% / tahun', duration: '18 bulan', highlight: true },
  { name: 'Prime Vanili', range: 'Rp 60 Juta ke atas', roi: '22% - 30% / tahun', duration: '24 bulan' },
];

const benefits = [
  'Akses dashboard perkembangan kebun secara real-time.',
  'Laporan bulanan performa kebun dan estimasi hasil panen.',
  'Notifikasi milestone penting pada fase budidaya.',
  'Pendampingan tim agronomi selama periode investasi.',
  'Dokumentasi aktivitas lapangan dengan foto dan catatan.',
  'Ringkasan finansial dan proyeksi ROI berkala.',
];
</script>

<template>
  <LandingLayout>

    <PaketHero
      title="Temukan Paket Investasi Sesuai Target Anda"
      description="Berikut simulasi paket investasi placeholder untuk memudahkan eksplorasi skema pendanaan di Omah Vanili. Konten ini nantinya akan diganti dengan data resmi."
    />

    <section class="max-w-7xl mx-auto px-6 pb-10 lg:pb-14">
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <PaketCard
          v-for="item in paketItems"
          :key="item.name"
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
