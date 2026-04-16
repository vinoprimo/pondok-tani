<script setup lang="ts">
import LandingLayout from '../../layouts/LandingLayout.vue';
import KemitraanHero from '../../components/landing-kemitraan/KemitraanHero.vue';
import KemitraanBenefitCard from '../../components/landing-kemitraan/KemitraanBenefitCard.vue';
import KemitraanStepList from '../../components/landing-kemitraan/KemitraanStepList.vue';

const benefits = [
  {
    title: 'Akses Pasar Terarah',
    description: 'Mitra mendapatkan dukungan pemasaran dan channel distribusi yang lebih stabil untuk hasil panen.',
  },
  {
    title: 'Standar Operasional',
    description: 'Setiap mitra dibekali SOP budidaya, pengeringan, dan grading agar kualitas produk konsisten.',
  },
  {
    title: 'Pendampingan Berkelanjutan',
    description: 'Tim lapangan melakukan kunjungan rutin untuk membantu evaluasi dan optimasi aktivitas kebun.',
  },
  {
    title: 'Transparansi Data',
    description: 'Aktivitas kemitraan terdokumentasi dalam sistem sehingga progres bisa dipantau bersama.',
  },
];

const steps = [
  'Isi formulir ketertarikan dan unggah profil kebun singkat untuk proses screening awal.',
  'Ikuti sesi wawancara dan asesmen lapangan bersama tim Omah Vanili.',
  'Aktifkan kemitraan, dapatkan onboarding operasional, lalu mulai fase eksekusi.',
];
</script>

<template>
  <LandingLayout>

    <KemitraanHero
      title="Bangun Kemitraan Vanili yang Tumbuh Bersama"
      subtitle="Halaman ini adalah placeholder informasi kemitraan. Nantinya akan berisi skema kerja sama resmi, persyaratan lengkap, dan benefit bisnis yang tervalidasi."
    />

    <section class="max-w-7xl mx-auto px-6 pb-10 lg:pb-14">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <KemitraanBenefitCard
          v-for="item in benefits"
          :key="item.title"
          :title="item.title"
          :description="item.description"
        />
      </div>
    </section>

    <KemitraanStepList :steps="steps" />
  </LandingLayout>
</template>
