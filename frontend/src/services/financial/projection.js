// Estimasi panen per pohon per tahun (kg) berdasarkan umur pohon
// Dari data terlampir
const HARVEST_ESTIMATION_BY_YEAR = {
  1: 0,      // Pertama (tahun pertama)
  2: 0,      // Kedua (tahun kedua)
  3: 0.25,   // Ketiga
  4: 0.5,    // Keempat
  5: 1,      // Kelima
  6: 2,      // Keenam
};

/**
 * Hitung proyeksi pendapatan 6 tahun dengan range min-max
 * @param {Object} params
 * @param {number} params.plantCount - Jumlah pohon
 * @param {Object} params.gradePercentages - { gradeName: percent, ... }
 * @param {Array} params.pricesByGrade - Array of { gradeName, priceMin, priceMax }
 * @returns {Array} Array of { year, harvestKg, minRevenue, maxRevenue, minPerTree, maxPerTree }
 */
export function calculateRevenueProjection(params) {
  const { plantCount = 100, gradePercentages = {}, pricesByGrade = [] } = params;

  // Validasi
  if (!plantCount || plantCount <= 0) {
    throw new Error("Jumlah pohon harus lebih dari 0");
  }

  // Buat map harga dari array
  const priceMap = new Map();
  pricesByGrade.forEach((item) => {
    priceMap.set(item.gradeName, {
      priceMin: item.priceMin || 0,
      priceMax: item.priceMax || 0,
    });
  });

  // Hitung total persen (should be 100)
  const totalPercent = Object.values(gradePercentages).reduce((a, b) => a + b, 0);
  if (Math.abs(totalPercent - 100) > 0.01) {
    console.warn(`Total grade percentage is ${totalPercent}%, expected 100%`);
  }

  // Proyeksi per tahun
  const projections = [];

  for (let year = 1; year <= 6; year++) {
    const estimatePerTree = HARVEST_ESTIMATION_BY_YEAR[year] || 0;
    const totalHarvestKg = estimatePerTree * plantCount;

    let totalMinRevenue = 0;
    let totalMaxRevenue = 0;

    // Hitung revenue per grade
    Object.entries(gradePercentages).forEach(([gradeName, percentage]) => {
      const gradePercent = percentage / 100;
      const gradeHarvestKg = totalHarvestKg * gradePercent;

      const prices = priceMap.get(gradeName) || { priceMin: 0, priceMax: 0 };
      totalMinRevenue += gradeHarvestKg * prices.priceMin;
      totalMaxRevenue += gradeHarvestKg * prices.priceMax;
    });

    projections.push({
      year,
      harvestKg: totalHarvestKg,
      minRevenue: Math.round(totalMinRevenue),
      maxRevenue: Math.round(totalMaxRevenue),
      minPerTree: Math.round(totalMinRevenue / plantCount),
      maxPerTree: Math.round(totalMaxRevenue / plantCount),
    });
  }

  return projections;
}

/**
 * Hitung statistik agregat proyeksi
 */
export function getProjectionStats(projections) {
  let totalHarvestKg = 0;
  let totalMinRevenue = 0;
  let totalMaxRevenue = 0;

  projections.forEach((p) => {
    totalHarvestKg += p.harvestKg;
    totalMinRevenue += p.minRevenue;
    totalMaxRevenue += p.maxRevenue;
  });

  const avgYearlyMin = totalMinRevenue / projections.length;
  const avgYearlyMax = totalMaxRevenue / projections.length;

  return {
    totalHarvestKg,
    totalMinRevenue,
    totalMaxRevenue,
    avgYearlyMin,
    avgYearlyMax,
  };
}
