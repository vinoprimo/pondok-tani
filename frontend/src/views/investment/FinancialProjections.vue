<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from "vue";
import {
  BarChart3,
  BadgeDollarSign,
  Info,
  Pencil,
  Plus,
  RefreshCcw,
  Search,
  Trash2,
  Wallet,
  X,
  TrendingUp,
} from "lucide-vue-next";
import { listPlantMonitorings } from "../../services/plant-monitoring/monitoring";
import {
  createOperationalCost,
  updateOperationalCost,
  deleteOperationalCost,
  listOperationalCosts,
} from "../../services/financial/operational-cost";
import { listActualRevenues } from "../../services/financial/revenue";
import { getCurrentUser } from "../../services/user/user";
import { getGrades, getNationalPrices } from "../../services/price/vanili";
import { calculateRevenueProjection, getProjectionStats } from "../../services/financial/projection";
import Pagination from "../../components/common/Pagination.vue";

type MonitoringRow = {
  plant_batch_id: number;
  batch_code: string;
  package_name: string;
};

type BatchOption = {
  id: number;
  batchCode: string;
  packageName: string;
};

type OperationalCostEntry = {
  id: number;
  plantBatchId: number;
  batchCode: string;
  packageName: string;
  category: string;
  amount: number;
  date: string;
  note: string;
  createdAt: string;
};

type OperationalCostApiRow = {
  id: number;
  plant_batch_id: number;
  batch_code: string;
  package_name: string;
  category: string;
  amount: number;
  cost_date: string;
  note: string;
  created_at: string;
};

type RevenueEntry = {
  id: number;
  plantBatchId: number;
  batchCode: string;
  packageName: string;
  source: string;
  amount: number;
  date: string;
  note: string;
  status: string;
  createdAt: string;
};

type RevenueApiRow = {
  id: number;
  plant_batch_id: number;
  batch_code: string;
  package_name: string;
  source: string;
  amount: number;
  revenue_date: string;
  note: string;
  status: string;
  created_at: string;
};

type DeleteTarget = {
  id: number | string;
  type: "expense";
  label: string;
};

type GradeOption = {
  id: number;
  grade_name: string;
};

type NationalPriceRow = {
  national_price_id: number;
  grade_id: number;
  grade_name?: string;
  province: string;
  price_min_per_kg: number;
  price_max_per_kg: number;
  price_per_kg: number;
  harvest_type: string;
};

type ProjectionInput = {
  plantCount: number;
  selectedProvince: string;
  gradePercentages: Record<string, number>;
};

type ProjectionRow = {
  year: number;
  harvestKg: number;
  minRevenue: number;
  maxRevenue: number;
  minPerTree: number;
  maxPerTree: number;
};

type VisualPeriodKey = "ytd" | "6m" | "12m" | "24m" | "5y";
type VisualSeriesKey = "expense" | "revenue" | "net";

type VisualBucket = {
  key: string;
  label: string;
  expense: number;
  revenue: number;
  net: number;
};

const COST_PAGE_SIZE = 8;
const REVENUE_PAGE_SIZE = 8;

const today = new Date().toISOString().slice(0, 10);
const tableView = ref<"expenses" | "revenues">("expenses");

const batches = ref<BatchOption[]>([]);
const loadingBatches = ref(false);
const batchError = ref("");
const userInvestmentError = ref("");

const expenses = ref<OperationalCostEntry[]>([]);
const revenues = ref<RevenueEntry[]>([]);
const userInvestmentTotal = ref(0);
const userInvestmentCount = ref(0);

const infoMessage = ref("");
const errorMessage = ref("");

const expenseModalOpen = ref(false);
const savingExpense = ref(false);

const editExpenseModalOpen = ref(false);
const editingExpense = ref<OperationalCostEntry | null>(null);
const savingEditExpense = ref(false);

const editExpenseForm = reactive({
  plantBatchId: "",
  category: "Pupuk",
  amount: 0,
  date: today,
  note: "",
});

const deleteModalOpen = ref(false);
const deleteSubmitting = ref(false);
const deleteTarget = ref<DeleteTarget | null>(null);

const expenseSearchQuery = ref("");
const revenueSearchQuery = ref("");
const expensePage = ref(1);
const revenuePage = ref(1);

const visualPeriod = ref<VisualPeriodKey>("ytd");
const visualSeriesVisibility = reactive<Record<VisualSeriesKey, boolean>>({
  expense: true,
  revenue: true,
  net: true,
});

const visualPeriodOptions: Array<{ value: VisualPeriodKey; label: string }> = [
  { value: "ytd", label: "Year to date" },
  { value: "6m", label: "6 bulan terakhir" },
  { value: "12m", label: "12 bulan terakhir" },
  { value: "24m", label: "24 bulan terakhir" },
  { value: "5y", label: "5 tahun terakhir" },
];

const visualSeriesOptions: Array<{
  key: VisualSeriesKey;
  label: string;
  color: string;
  activeClass: string;
  inactiveClass: string;
}> = [
  {
    key: "expense",
    label: "Biaya operasional",
    color: "#ef4444",
    activeClass: "border-red-200 bg-red-50 text-red-700",
    inactiveClass: "border-gray-200 bg-white text-gray-500",
  },
  {
    key: "revenue",
    label: "Pendapatan",
    color: "#10b981",
    activeClass: "border-emerald-200 bg-emerald-50 text-emerald-700",
    inactiveClass: "border-gray-200 bg-white text-gray-500",
  },
  {
    key: "net",
    label: "Net profit",
    color: "#3b82f6",
    activeClass: "border-blue-200 bg-blue-50 text-blue-700",
    inactiveClass: "border-gray-200 bg-white text-gray-500",
  },
];

const expenseForm = reactive({
  plantBatchId: "",
  category: "Pupuk",
  amount: 0,
  date: today,
  note: "",
});


const costCategoryOptions = [
  "Pupuk",
  "Pestisida",
  "Tenaga kerja",
  "Transportasi",
  "Peralatan",
  "Lainnya",
];

// Projection modal state
const projectionModalOpen = ref(false);
const projectionFormulaModalOpen = ref(false);
const loadingGradesAndPrices = ref(false);
const projectingRevenue = ref(false);
const autoFilledProjectionGradeName = ref("");

const grades = ref<GradeOption[]>([]);
const nationalPrices = ref<NationalPriceRow[]>([]);

const projectionInput = reactive<ProjectionInput>({
  plantCount: 100,
  selectedProvince: "",
  gradePercentages: {},
});

const projectionResults = ref<ProjectionRow[]>([]);
const projectionStats = ref({ totalHarvestKg: 0, totalMinRevenue: 0, totalMaxRevenue: 0, avgYearlyMin: 0, avgYearlyMax: 0 });
const projectionChartDimensions = {
  width: 800,
  height: 320,
  marginTop: 24,
  marginRight: 24,
  marginBottom: 48,
  marginLeft: 80,
};

const filteredExpenses = computed(() => {
  const query = expenseSearchQuery.value.trim().toLowerCase();
  const sorted = [...expenses.value].sort((a, b) => {
    const byDate = b.date.localeCompare(a.date);
    if (byDate !== 0) return byDate;
    return b.createdAt.localeCompare(a.createdAt);
  });

  if (!query) return sorted;

  return sorted.filter((item) => {
    return (
      item.batchCode.toLowerCase().includes(query) ||
      item.packageName.toLowerCase().includes(query) ||
      item.category.toLowerCase().includes(query) ||
      item.note.toLowerCase().includes(query)
    );
  });
});

const filteredRevenues = computed(() => {
  const query = revenueSearchQuery.value.trim().toLowerCase();
  const sorted = [...revenues.value].sort((a, b) => {
    const byDate = b.date.localeCompare(a.date);
    if (byDate !== 0) return byDate;
    return b.createdAt.localeCompare(a.createdAt);
  });

  if (!query) return sorted;

  return sorted.filter((item) => {
    return (
      item.batchCode.toLowerCase().includes(query) ||
      item.packageName.toLowerCase().includes(query) ||
      item.source.toLowerCase().includes(query) ||
      item.note.toLowerCase().includes(query)
    );
  });
});

const expenseTotalPages = computed(() => Math.max(1, Math.ceil(filteredExpenses.value.length / COST_PAGE_SIZE)));
const revenueTotalPages = computed(() => Math.max(1, Math.ceil(filteredRevenues.value.length / REVENUE_PAGE_SIZE)));

const paginatedExpenses = computed(() => {
  const start = (expensePage.value - 1) * COST_PAGE_SIZE;
  return filteredExpenses.value.slice(start, start + COST_PAGE_SIZE);
});

const paginatedRevenues = computed(() => {
  const start = (revenuePage.value - 1) * REVENUE_PAGE_SIZE;
  return filteredRevenues.value.slice(start, start + REVENUE_PAGE_SIZE);
});

const totalExpense = computed(() => expenses.value.reduce((acc, item) => acc + item.amount, 0));
const totalRevenue = computed(() => revenues.value.reduce((acc, item) => acc + item.amount, 0));
const netBalance = computed(() => totalRevenue.value - totalExpense.value);
const totalInvestment = computed(() => userInvestmentTotal.value);
const roiNetGain = computed(() => netBalance.value - totalInvestment.value);
const roiPercent = computed(() => {
  if (totalInvestment.value <= 0) return 0;
  return (roiNetGain.value / totalInvestment.value) * 100;
});

const revenueVsCostPercent = computed(() => {
  if (totalExpense.value <= 0) return 0;
  return (totalRevenue.value / totalExpense.value) * 100;
});

const netStatusLabel = computed(() => {
  if (netBalance.value <= 0) {
    return "Arus kas masih negatif";
  }

  if (totalInvestment.value <= 0 || roiNetGain.value >= 0) {
    return "Break-even tercapai";
  }

  return "Arus kas positif, belum break-even investasi";
});

const netStatusClass = computed(() => {
  if (netBalance.value <= 0) return "text-rose-600";
  if (totalInvestment.value <= 0 || roiNetGain.value >= 0) return "text-emerald-600";
  return "text-amber-600";
});

// Projection computed properties
const availableProvinces = computed(() => {
  const provinces = new Set<string>();
  nationalPrices.value.forEach((price) => {
    provinces.add(price.province);
  });
  return Array.from(provinces).sort();
});

const filteredPricesByProvince = computed(() => {
  return nationalPrices.value.filter((price) => price.province === projectionInput.selectedProvince);
});

const projectionGradeTotalPercent = computed(() => {
  return Object.values(projectionInput.gradePercentages).reduce((total: number, value) => {
    return total + (Number(value) || 0);
  }, 0);
});

const isProjectionGradeTotalValid = computed(() => Math.abs(projectionGradeTotalPercent.value - 100) < 0.01);

const canSubmitProjection = computed(() => {
  return Boolean(projectionInput.selectedProvince) && projectionInput.plantCount > 0 && isProjectionGradeTotalValid.value;
});

const projectionChartInnerWidth = computed(
  () => projectionChartDimensions.width - projectionChartDimensions.marginLeft - projectionChartDimensions.marginRight
);

const projectionChartInnerHeight = computed(
  () => projectionChartDimensions.height - projectionChartDimensions.marginTop - projectionChartDimensions.marginBottom
);

const projectionChartVisibleValues = computed(() => {
  if (!projectionResults.value.length) return [0];
  const values: number[] = [];
  projectionResults.value.forEach((row) => {
    values.push(row.minRevenue, row.maxRevenue);
  });
  return values.length ? values : [0];
});

const projectionChartMinValue = computed(() => Math.min(0, ...projectionChartVisibleValues.value));
const projectionChartMaxValue = computed(() => Math.max(0, ...projectionChartVisibleValues.value));
const projectionChartRange = computed(() => {
  const range = projectionChartMaxValue.value - projectionChartMinValue.value;
  return range === 0 ? 1 : range;
});

const getProjectionChartX = (index: number, pointCount: number) => {
  if (pointCount <= 1) {
    return projectionChartDimensions.marginLeft + projectionChartInnerWidth.value / 2;
  }
  const ratio = index / (pointCount - 1);
  return projectionChartDimensions.marginLeft + ratio * projectionChartInnerWidth.value;
};

const getProjectionChartY = (value: number) => {
  const ratio = (projectionChartMaxValue.value - value) / projectionChartRange.value;
  return projectionChartDimensions.marginTop + ratio * projectionChartInnerHeight.value;
};

const projectionChartYTicks = computed(() => {
  const tickCount = 5;
  return Array.from({ length: tickCount }, (_, index) => {
    const ratio = index / (tickCount - 1);
    const value = projectionChartMaxValue.value - ratio * projectionChartRange.value;
    return {
      value,
      y: getProjectionChartY(value),
    };
  });
});

const projectionChartXTicks = computed(() => {
  const count = projectionResults.value.length;
  return projectionResults.value.map((row, index) => ({
    label: `Tahun ${row.year}`,
    x: getProjectionChartX(index, count),
  }));
});

const projectionChartZeroLineY = computed(() => getProjectionChartY(0));

const roiStatusLabel = computed(() => {
  if (totalInvestment.value <= 0) return "Belum ada nilai investasi";
  if (roiPercent.value > 0) return "Return on investment positif";
  if (roiPercent.value === 0) return "Return on investment netral";
  return "Return on investment negatif";
});

const chartDimensions = {
  width: 980,
  height: 360,
  marginTop: 24,
  marginRight: 28,
  marginBottom: 46,
  marginLeft: 76,
};

const parseSafeDate = (rawDate: string) => {
  if (!rawDate) return null;
  const [year, month, day] = rawDate.split("-").map(Number);
  if (year && month && day) {
    return new Date(year, month - 1, day);
  }
  const parsed = new Date(rawDate);
  if (Number.isNaN(parsed.getTime())) return null;
  return parsed;
};

const visualPeriodLabel = computed(() => {
  return visualPeriodOptions.find((item) => item.value === visualPeriod.value)?.label ?? "-";
});

const visualBuckets = computed<VisualBucket[]>(() => {
  const now = new Date();
  const monthAnchor = new Date(now.getFullYear(), now.getMonth(), 1);
  const bucketMap = new Map<string, { label: string; revenue: number; expense: number }>();
  const orderedKeys: string[] = [];

  const addBucket = (key: string, label: string) => {
    orderedKeys.push(key);
    bucketMap.set(key, { label, revenue: 0, expense: 0 });
  };

  if (visualPeriod.value === "5y") {
    for (let i = 4; i >= 0; i -= 1) {
      const year = monthAnchor.getFullYear() - i;
      addBucket(String(year), String(year));
    }
  } else {
    const monthCount =
      visualPeriod.value === "ytd"
        ? monthAnchor.getMonth() + 1
        : visualPeriod.value === "6m"
        ? 6
        : visualPeriod.value === "12m"
        ? 12
        : 24;
    for (let i = monthCount - 1; i >= 0; i -= 1) {
      const d = new Date(monthAnchor.getFullYear(), monthAnchor.getMonth() - i, 1);
      const key = `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, "0")}`;
      addBucket(
        key,
        d.toLocaleDateString("id-ID", {
          month: "short",
          year: "2-digit",
        })
      );
    }
  }

  const toBucketKey = (rawDate: string) => {
    const parsed = parseSafeDate(rawDate);
    if (!parsed) return "";
    if (visualPeriod.value === "5y") {
      return String(parsed.getFullYear());
    }
    return `${parsed.getFullYear()}-${String(parsed.getMonth() + 1).padStart(2, "0")}`;
  };

  expenses.value.forEach((item) => {
    const key = toBucketKey(item.date);
    const target = bucketMap.get(key);
    if (target) {
      target.expense += Number(item.amount) || 0;
    }
  });

  revenues.value.forEach((item) => {
    const key = toBucketKey(item.date);
    const target = bucketMap.get(key);
    if (target) {
      target.revenue += Number(item.amount) || 0;
    }
  });

  return orderedKeys.map((key) => {
    const row = bucketMap.get(key)!;
    return {
      key,
      label: row.label,
      expense: row.expense,
      revenue: row.revenue,
      net: row.revenue - row.expense,
    };
  });
});

const visibleVisualSeries = computed(() => {
  return visualSeriesOptions.filter((item) => visualSeriesVisibility[item.key]);
});

const chartInnerWidth = computed(
  () => chartDimensions.width - chartDimensions.marginLeft - chartDimensions.marginRight
);
const chartInnerHeight = computed(
  () => chartDimensions.height - chartDimensions.marginTop - chartDimensions.marginBottom
);

const chartVisibleValues = computed(() => {
  if (!visibleVisualSeries.value.length) return [0];

  const values: number[] = [];
  visualBuckets.value.forEach((bucket) => {
    visibleVisualSeries.value.forEach((series) => {
      values.push(bucket[series.key]);
    });
  });
  return values.length ? values : [0];
});

const chartMinValue = computed(() => Math.min(0, ...chartVisibleValues.value));
const chartMaxValue = computed(() => Math.max(0, ...chartVisibleValues.value));
const chartRange = computed(() => {
  const range = chartMaxValue.value - chartMinValue.value;
  return range === 0 ? 1 : range;
});

const getChartX = (index: number, pointCount: number) => {
  if (pointCount <= 1) {
    return chartDimensions.marginLeft + chartInnerWidth.value / 2;
  }
  const ratio = index / (pointCount - 1);
  return chartDimensions.marginLeft + ratio * chartInnerWidth.value;
};

const getChartY = (value: number) => {
  const ratio = (chartMaxValue.value - value) / chartRange.value;
  return chartDimensions.marginTop + ratio * chartInnerHeight.value;
};

const chartYTicks = computed(() => {
  const tickCount = 5;
  return Array.from({ length: tickCount }, (_, index) => {
    const ratio = index / (tickCount - 1);
    const value = chartMaxValue.value - ratio * chartRange.value;
    return {
      value,
      y: getChartY(value),
    };
  });
});

const chartXTicks = computed(() => {
  const count = visualBuckets.value.length;
  return visualBuckets.value.map((bucket, index) => ({
    label: bucket.label,
    x: getChartX(index, count),
  }));
});

const chartZeroLineY = computed(() => getChartY(0));

const chartSeriesLines = computed(() => {
  const count = visualBuckets.value.length;
  return visibleVisualSeries.value.map((series) => {
    const points = visualBuckets.value.map((bucket, index) => {
      const value = bucket[series.key];
      return {
        x: getChartX(index, count),
        y: getChartY(value),
        value,
      };
    });

    const path = points
      .map((point, index) => `${index === 0 ? "M" : "L"} ${point.x.toFixed(2)} ${point.y.toFixed(2)}`)
      .join(" ");

    return {
      ...series,
      path,
      points,
    };
  });
});

const visualPeriodTotals = computed(() => {
  return visualBuckets.value.reduce(
    (acc, bucket) => {
      acc.expense += bucket.expense;
      acc.revenue += bucket.revenue;
      acc.net += bucket.net;
      return acc;
    },
    { expense: 0, revenue: 0, net: 0 }
  );
});

const compactNumberFormatter = new Intl.NumberFormat("id-ID", {
  notation: "compact",
  maximumFractionDigits: 1,
});

const formatAxisRupiah = (value: number) => {
  const sign = value < 0 ? "-" : "";
  return `${sign}Rp${compactNumberFormatter.format(Math.abs(value))}`;
};

const chartHasVisibleSeries = computed(() => visibleVisualSeries.value.length > 0);

const expenseCategorySeries = computed(() => {
  const grouped = new Map<string, number>();
  expenses.value.forEach((item) => {
    const key = item.category || "Tanpa kategori";
    grouped.set(key, (grouped.get(key) || 0) + (Number(item.amount) || 0));
  });

  return Array.from(grouped.entries())
    .map(([category, amount]) => ({
      category,
      amount,
      percent: totalExpense.value > 0 ? (amount / totalExpense.value) * 100 : 0,
    }))
    .sort((a, b) => b.amount - a.amount)
    .slice(0, 6);
});

const maxExpenseCategoryValue = computed(() => {
  return Math.max(1, ...expenseCategorySeries.value.map((item) => item.amount));
});

const revenueBatchSeries = computed(() => {
  const grouped = new Map<string, number>();
  revenues.value.forEach((item) => {
    const key = item.batchCode || "Tanpa batch";
    grouped.set(key, (grouped.get(key) || 0) + (Number(item.amount) || 0));
  });

  return Array.from(grouped.entries())
    .map(([batchCode, amount]) => ({ batchCode, amount }))
    .sort((a, b) => b.amount - a.amount)
    .slice(0, 6);
});

const maxRevenueBatchValue = computed(() => {
  return Math.max(1, ...revenueBatchSeries.value.map((item) => item.amount));
});

watch(expenseSearchQuery, () => {
  expensePage.value = 1;
});

watch(revenueSearchQuery, () => {
  revenuePage.value = 1;
});

watch(filteredExpenses, () => {
  if (expensePage.value > expenseTotalPages.value) {
    expensePage.value = expenseTotalPages.value;
  }
});

watch(filteredRevenues, () => {
  if (revenuePage.value > revenueTotalPages.value) {
    revenuePage.value = revenueTotalPages.value;
  }
});

const formatRupiah = (value: number) =>
  new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    maximumFractionDigits: 0,
  }).format(value);

const formatPercent = (value: number) => {
  if (!Number.isFinite(value)) return "0.0%";
  const sign = value > 0 ? "+" : "";
  return `${sign}${value.toFixed(1)}%`;
};

const formatDate = (value: string) => {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return "-";
  return date.toLocaleDateString("id-ID", {
    day: "2-digit",
    month: "short",
    year: "numeric",
  });
};

function getBatchById(batchId: number) {
  return batches.value.find((item) => item.id === batchId) || null;
}

function ensureDefaultBatchInForm() {
  const firstBatch = batches.value[0];
  if (!firstBatch) return;

  if (!expenseForm.plantBatchId) {
    expenseForm.plantBatchId = String(firstBatch.id);
  }
}

async function fetchBatches() {
  loadingBatches.value = true;
  batchError.value = "";
  try {
    const response = await listPlantMonitorings();
    const rows = Array.isArray(response.data) ? (response.data as MonitoringRow[]) : [];

    const map = new Map<number, BatchOption>();
    rows.forEach((item) => {
      if (!item.plant_batch_id || !item.batch_code) return;
      if (!map.has(item.plant_batch_id)) {
        map.set(item.plant_batch_id, {
          id: item.plant_batch_id,
          batchCode: item.batch_code,
          packageName: item.package_name || "Paket tanpa nama",
        });
      }
    });

    batches.value = Array.from(map.values()).sort((a, b) => a.batchCode.localeCompare(b.batchCode));
    ensureDefaultBatchInForm();
  } catch (error) {
    console.error("Failed to load plant batches from monitoring data", error);
    batchError.value = "Gagal memuat batch tanaman. Pastikan data plant monitoring sudah tersedia.";
  } finally {
    loadingBatches.value = false;
  }
}

async function fetchOperationalCosts() {
  try {
    const response = await listOperationalCosts();
    const rows = Array.isArray(response.data) ? (response.data as OperationalCostApiRow[]) : [];

    expenses.value = rows.map((item) => ({
      id: item.id,
      plantBatchId: item.plant_batch_id,
      batchCode: item.batch_code,
      packageName: item.package_name || "Paket tanpa nama",
      category: item.category,
      amount: Number(item.amount) || 0,
      date: item.cost_date,
      note: item.note || "",
      createdAt: item.created_at,
    }));
  } catch (error) {
    console.error("Failed to load operational costs", error);
    errorMessage.value = "Gagal memuat data biaya operasional dari server.";
  }
}

async function fetchRevenues() {
  try {
    const response = await listActualRevenues();
    const rows = Array.isArray(response.data) ? (response.data as RevenueApiRow[]) : [];

    revenues.value = rows.map((item) => ({
      id: item.id,
      plantBatchId: item.plant_batch_id,
      batchCode: item.batch_code,
      packageName: item.package_name || "Paket tanpa nama",
      source: item.source,
      amount: Number(item.amount) || 0,
      date: item.revenue_date,
      note: item.note || "",
      status: item.status || "realisasi",
      createdAt: item.created_at,
    }));
  } catch (error) {
    console.error("Failed to load revenues", error);
    errorMessage.value = "Gagal memuat data pendapatan dari server.";
  }
}

async function fetchCurrentUserInvestment() {
  userInvestmentError.value = "";
  try {
    const response = await getCurrentUser();
    const payload = response?.data ?? {};
    const investments = Array.isArray(payload.investments) ? payload.investments : [];
    userInvestmentCount.value = investments.length;

    userInvestmentTotal.value = investments.reduce((acc: number, item: { amount?: number | string }) => {
      return acc + (Number(item?.amount) || 0);
    }, 0);
  } catch (error) {
    console.error("Failed to load current user investment data", error);
    userInvestmentTotal.value = 0;
    userInvestmentCount.value = 0;
    userInvestmentError.value = "Gagal memuat total investasi user. ROI mungkin belum akurat.";
  }
}

function resetExpenseForm() {
  expenseForm.category = "Pupuk";
  expenseForm.amount = 0;
  expenseForm.date = today;
  expenseForm.note = "";
  if (batches.value.length > 0) {
    expenseForm.plantBatchId = String(batches.value[0].id);
  }
}

function openExpenseModal() {
  errorMessage.value = "";
  infoMessage.value = "";
  resetExpenseForm();
  expenseModalOpen.value = true;
}

function closeExpenseModal() {
  expenseModalOpen.value = false;
}

async function submitExpense() {
  errorMessage.value = "";
  infoMessage.value = "";

  const batchId = Number(expenseForm.plantBatchId);
  if (!batchId) {
    errorMessage.value = "Batch tanaman wajib dipilih.";
    return;
  }

  const batch = getBatchById(batchId);
  if (!batch) {
    errorMessage.value = "Batch tanaman tidak valid.";
    return;
  }

  if (!expenseForm.category.trim()) {
    errorMessage.value = "Kategori biaya wajib diisi.";
    return;
  }

  if (!Number.isFinite(expenseForm.amount) || expenseForm.amount <= 0) {
    errorMessage.value = "Nominal biaya harus lebih dari 0.";
    return;
  }

  if (!expenseForm.date) {
    errorMessage.value = "Tanggal biaya wajib diisi.";
    return;
  }

  savingExpense.value = true;
  try {
    await createOperationalCost({
      plant_batch_id: batch.id,
      category: expenseForm.category.trim(),
      amount: Number(expenseForm.amount),
      cost_date: expenseForm.date,
      note: expenseForm.note.trim(),
    });

    await fetchOperationalCosts();

    infoMessage.value = "Biaya operasional berhasil ditambahkan.";
    closeExpenseModal();
  } catch (error) {
    console.error("Failed to save operational cost", error);
    errorMessage.value = "Gagal menyimpan biaya operasional.";
  } finally {
    savingExpense.value = false;
  }
}

function openEditExpenseModal(item: OperationalCostEntry) {
  editingExpense.value = item;
  editExpenseForm.plantBatchId = String(item.plantBatchId);
  editExpenseForm.category = item.category;
  editExpenseForm.amount = item.amount;
  editExpenseForm.date = item.date ? item.date.slice(0, 10) : today;
  editExpenseForm.note = item.note || "";
  errorMessage.value = "";
  editExpenseModalOpen.value = true;
}

function closeEditExpenseModal() {
  editExpenseModalOpen.value = false;
  editingExpense.value = null;
}

async function submitEditExpense() {
  if (!editingExpense.value) return;
  errorMessage.value = "";
  infoMessage.value = "";

  const batchId = Number(editExpenseForm.plantBatchId);
  if (!batchId) {
    errorMessage.value = "Batch tanaman wajib dipilih.";
    return;
  }

  const batch = getBatchById(batchId);
  if (!batch) {
    errorMessage.value = "Batch tanaman tidak valid.";
    return;
  }

  if (!editExpenseForm.category.trim()) {
    errorMessage.value = "Kategori biaya wajib diisi.";
    return;
  }

  if (!Number.isFinite(editExpenseForm.amount) || editExpenseForm.amount <= 0) {
    errorMessage.value = "Nominal biaya harus lebih dari 0.";
    return;
  }

  if (!editExpenseForm.date) {
    errorMessage.value = "Tanggal biaya wajib diisi.";
    return;
  }

  savingEditExpense.value = true;
  try {
    await updateOperationalCost(editingExpense.value.id, {
      plant_batch_id: batch.id,
      category: editExpenseForm.category.trim(),
      amount: Number(editExpenseForm.amount),
      cost_date: editExpenseForm.date,
      note: editExpenseForm.note.trim(),
    });

    await fetchOperationalCosts();
    infoMessage.value = "Biaya operasional berhasil diperbarui.";
    closeEditExpenseModal();
  } catch (error) {
    console.error("Failed to update operational cost", error);
    errorMessage.value = "Gagal memperbarui biaya operasional.";
  } finally {
    savingEditExpense.value = false;
  }
}

function askDeleteExpense(item: OperationalCostEntry) {
  deleteTarget.value = {
    id: item.id,
    type: "expense",
    label: `${item.category} - ${item.batchCode}`,
  };
  deleteModalOpen.value = true;
}

function closeDeleteModal() {
  if (deleteSubmitting.value) return;
  deleteModalOpen.value = false;
  deleteTarget.value = null;
}

async function confirmDelete() {
  if (!deleteTarget.value) return;

  deleteSubmitting.value = true;
  let success = false;
  try {
    await deleteOperationalCost(deleteTarget.value.id);
    await fetchOperationalCosts();
    infoMessage.value = "Data biaya operasional berhasil dihapus.";
    success = true;
  } catch (error) {
    console.error("Failed to delete tracker data", error);
    errorMessage.value = "Gagal menghapus data.";
  } finally {
    deleteSubmitting.value = false;
  }

  if (success) {
    closeDeleteModal();
  }
}

async function refreshTrackerData() {
  await Promise.all([fetchBatches(), fetchOperationalCosts(), fetchRevenues(), fetchCurrentUserInvestment()]);
  errorMessage.value = "";
  infoMessage.value = "Data tracker diperbarui.";
}

async function fetchGradesAndNationalPrices() {
  loadingGradesAndPrices.value = true;
  try {
    const [gradesRes, pricesRes] = await Promise.all([getGrades(), getNationalPrices()]);

    grades.value = Array.isArray(gradesRes.data)
      ? gradesRes.data.map((g: any) => ({
          id: g.id,
          grade_name: g.grade_name,
        }))
      : [];

    nationalPrices.value = Array.isArray(pricesRes.data)
      ? pricesRes.data.map((p: any) => ({
          national_price_id: p.national_price_id,
          grade_id: p.grade_id,
          grade_name: p.grade?.grade_name,
          province: p.province,
          price_min_per_kg: p.price_min_per_kg,
          price_max_per_kg: p.price_max_per_kg,
          price_per_kg: p.price_per_kg,
          harvest_type: p.harvest_type,
        }))
      : [];
  } catch (error) {
    console.error("Failed to fetch grades and national prices", error);
    errorMessage.value = "Gagal memuat data grade dan harga nasional.";
  } finally {
    loadingGradesAndPrices.value = false;
  }
}

async function openProjectionModal() {
  errorMessage.value = "";
  infoMessage.value = "";
  projectionInput.plantCount = 100;
  projectionInput.selectedProvince = availableProvinces.value[0] || "";
  projectionInput.gradePercentages = {};
  projectionResults.value = [];
  autoFilledProjectionGradeName.value = "";

  if (!grades.value.length || !nationalPrices.value.length) {
    await fetchGradesAndNationalPrices();
    projectionInput.selectedProvince = availableProvinces.value[0] || "";
  }

  // Initialize grade percentages
  grades.value.forEach((grade) => {
    projectionInput.gradePercentages[grade.grade_name] = 0;
  });

  projectionModalOpen.value = true;
}

function closeProjectionModal() {
  projectionModalOpen.value = false;
  projectionFormulaModalOpen.value = false;
}

function autofillRemainingGradePercentage(changedGradeName: string) {
  if (grades.value.length !== 3) return;

  const gradeNames = grades.value.map((grade) => grade.grade_name);
  if (autoFilledProjectionGradeName.value === changedGradeName) {
    autoFilledProjectionGradeName.value = "";
    return;
  }

  if (autoFilledProjectionGradeName.value) {
    const manualTotal = gradeNames
      .filter((gradeName) => gradeName !== autoFilledProjectionGradeName.value)
      .reduce((total, gradeName) => {
        return total + (Number(projectionInput.gradePercentages[gradeName]) || 0);
      }, 0);

    projectionInput.gradePercentages[autoFilledProjectionGradeName.value] = Math.max(0, Number((100 - manualTotal).toFixed(1)));
    return;
  }

  const filledGradeNames = gradeNames.filter((gradeName) => {
    return (Number(projectionInput.gradePercentages[gradeName]) || 0) > 0;
  });
  const emptyGradeNames = gradeNames.filter((gradeName) => {
    return (Number(projectionInput.gradePercentages[gradeName]) || 0) === 0;
  });

  if (filledGradeNames.length !== 2 || emptyGradeNames.length !== 1) return;

  const filledTotal = filledGradeNames.reduce((total, gradeName) => {
    return total + (Number(projectionInput.gradePercentages[gradeName]) || 0);
  }, 0);
  const remainingPercentage = Math.max(0, Number((100 - filledTotal).toFixed(1)));

  projectionInput.gradePercentages[emptyGradeNames[0]] = remainingPercentage;
  autoFilledProjectionGradeName.value = emptyGradeNames[0];
}

function submitProjection() {
  errorMessage.value = "";
  infoMessage.value = "";

  if (!projectionInput.selectedProvince) {
    errorMessage.value = "Provinsi harus dipilih.";
    return;
  }

  if (!projectionInput.plantCount || projectionInput.plantCount <= 0) {
    errorMessage.value = "Jumlah pohon harus lebih dari 0.";
    return;
  }

  // Validasi persentase total = 100%
  if (!isProjectionGradeTotalValid.value) {
    errorMessage.value = `Total persentase grade harus 100% (saat ini ${projectionGradeTotalPercent.value.toFixed(1)}%).`;
    return;
  }

  // Prepare price data for selected province and grades
  const pricesByGrade = grades.value
    .map((grade) => {
      const priceRow = filteredPricesByProvince.value.find((p) => p.grade_id === grade.id);
      if (!priceRow) return null;

      return {
        gradeName: grade.grade_name,
        priceMin: priceRow.price_min_per_kg,
        priceMax: priceRow.price_max_per_kg,
      };
    })
    .filter((item) => item !== null);

  if (!pricesByGrade.length) {
    errorMessage.value = "Tidak ada data harga untuk provinsi dan grade yang dipilih.";
    return;
  }

  projectingRevenue.value = true;
  try {
    const results = calculateRevenueProjection({
      plantCount: projectionInput.plantCount,
      gradePercentages: projectionInput.gradePercentages,
      pricesByGrade,
    });

    projectionResults.value = results;
    projectionStats.value = getProjectionStats(results);

    infoMessage.value = "Proyeksi pendapatan berhasil dihitung.";
  } catch (error) {
    console.error("Failed to calculate projection", error);
    errorMessage.value = `Gagal menghitung proyeksi: ${(error as any).message}`;
    projectionResults.value = [];
  } finally {
    projectingRevenue.value = false;
  }
}

onMounted(async () => {
  await Promise.all([fetchBatches(), fetchOperationalCosts(), fetchRevenues(), fetchCurrentUserInvestment(), fetchGradesAndNationalPrices()]);
});
</script>

<template>
  <div class="space-y-6">
    <div class="rounded-xl border border-gray-200 bg-white p-6">
      <div class="flex flex-wrap items-start justify-between gap-4">
        <div>
          <h2 class="text-2xl font-semibold text-gray-900">Financial Tracker</h2>
          <p class="mt-1 text-gray-600">
            Catat biaya operasional per batch dan pantau pendapatan agar pengelolaan keuangan lebih rapi.
          </p>
          <p class="mt-2 text-xs text-emerald-700">
            Pendapatan otomatis diambil dari transaksi penjualan panen.
          </p>
        </div>

        <div class="flex flex-wrap items-center gap-2">
          <button
            type="button"
            class="inline-flex items-center gap-2 rounded-lg border border-gray-300 px-3 py-2 text-sm font-medium text-gray-700 transition hover:bg-gray-50"
            @click="refreshTrackerData"
          >
            <RefreshCcw class="h-4 w-4" />
            Refresh data
          </button>
          <button
            type="button"
            class="inline-flex items-center gap-2 rounded-lg bg-violet-600 px-3 py-2 text-sm font-medium text-white transition hover:bg-violet-700 disabled:cursor-not-allowed disabled:opacity-60"
            :disabled="loadingGradesAndPrices"
            @click="openProjectionModal"
          >
            <TrendingUp class="h-4 w-4" />
            Proyeksikan Pendapatan
          </button>
          <button
            type="button"
            class="inline-flex items-center gap-2 rounded-lg bg-emerald-600 px-3 py-2 text-sm font-medium text-white transition hover:bg-emerald-700 disabled:cursor-not-allowed disabled:opacity-60"
            :disabled="!batches.length || loadingBatches"
            @click="openExpenseModal"
          >
            <Plus class="h-4 w-4" />
            Input biaya operasional
          </button>
        </div>
      </div>
    </div>

    <div class="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
      <article class="rounded-2xl border border-gray-200 bg-white p-5 shadow-sm transition hover:-translate-y-0.5 hover:shadow-md">
        <div class="flex items-center gap-3">
          <span class="inline-flex h-11 w-11 items-center justify-center rounded-xl bg-emerald-100 text-emerald-600">
            <BadgeDollarSign class="h-5 w-5" />
          </span>
          <p class="text-base font-medium text-gray-700">Total Pendapatan</p>
        </div>
        <p class="mt-5 text-4xl font-semibold tracking-tight text-gray-900">{{ formatRupiah(totalRevenue) }}</p>
        <p class="mt-2 text-sm text-emerald-600">
          {{ revenues.length }} transaksi • {{ formatPercent(revenueVsCostPercent) }} dari biaya
        </p>
      </article>

      <article class="rounded-2xl border border-gray-200 bg-white p-5 shadow-sm transition hover:-translate-y-0.5 hover:shadow-md">
        <div class="flex items-center gap-3">
          <span class="inline-flex h-11 w-11 items-center justify-center rounded-xl bg-blue-100 text-blue-600">
            <BarChart3 class="h-5 w-5" />
          </span>
          <p class="text-base font-medium text-gray-700">Saldo Bersih</p>
        </div>
        <p class="mt-5 text-4xl font-semibold tracking-tight" :class="netBalance >= 0 ? 'text-gray-900' : 'text-rose-700'">
          {{ formatRupiah(netBalance) }}
        </p>
        <p class="mt-2 text-sm" :class="netStatusClass">
          {{ netStatusLabel }}
        </p>
      </article>

      <article class="rounded-2xl border border-gray-200 bg-white p-5 shadow-sm transition hover:-translate-y-0.5 hover:shadow-md">
        <div class="flex items-center gap-3">
          <span class="inline-flex h-11 w-11 items-center justify-center rounded-xl bg-violet-100 text-violet-600">
            <Wallet class="h-5 w-5" />
          </span>
          <p class="text-base font-medium text-gray-700">Total Biaya Operasional</p>
        </div>
        <p class="mt-5 text-4xl font-semibold tracking-tight text-gray-900">{{ formatRupiah(totalExpense) }}</p>
        <p class="mt-2 text-sm text-gray-600">{{ expenses.length }} transaksi biaya</p>
      </article>

      <article class="rounded-2xl border border-gray-200 bg-white p-5 shadow-sm transition hover:-translate-y-0.5 hover:shadow-md">
        <div class="flex items-center gap-3">
          <span class="inline-flex h-11 w-11 items-center justify-center rounded-xl bg-amber-100 text-amber-600">
            <BadgeDollarSign class="h-5 w-5" />
          </span>
          <p class="text-base font-medium text-gray-700">ROI %</p>
        </div>
        <p class="mt-5 text-4xl font-semibold tracking-tight" :class="roiPercent >= 0 ? 'text-gray-900' : 'text-rose-700'">
          {{ totalInvestment > 0 ? formatPercent(roiPercent) : '-' }}
        </p>
        <p class="mt-2 text-sm" :class="roiPercent >= 0 ? 'text-emerald-600' : 'text-rose-600'">
          {{ totalInvestment > 0 ? roiStatusLabel : 'Belum ada investasi user' }}
        </p>
        <p class="mt-1 text-xs text-gray-500">
          Total biaya investasi user: {{ formatRupiah(totalInvestment) }} ({{ userInvestmentCount }} data investasi)
        </p>
        <p class="mt-1 text-xs text-gray-500">Rumus ROI: (Saldo Bersih - Total Investasi) / Total Investasi</p>
      </article>
    </div>

    <section class="rounded-xl border border-gray-200 bg-white p-5">
      <div class="flex flex-wrap items-start justify-between gap-4">
        <div>
          <div class="flex items-center gap-2">
            <BarChart3 class="h-5 w-5 text-emerald-600" />
            <h3 class="text-lg font-semibold text-gray-900">Visualisasi Kinerja Keuangan</h3>
          </div>
          <p class="mt-1 text-sm text-gray-600">
            Grafik garis biaya operasional, pendapatan, dan net profit berdasarkan periode yang dipilih.
          </p>
        </div>

        <label class="w-full max-w-xs space-y-1 text-sm text-gray-600">
          <span class="font-medium text-gray-700">Periode visual</span>
          <select
            v-model="visualPeriod"
            class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm outline-none transition focus:border-emerald-500 focus:ring-2 focus:ring-emerald-100"
          >
            <option v-for="item in visualPeriodOptions" :key="item.value" :value="item.value">
              {{ item.label }}
            </option>
          </select>
        </label>
      </div>

      <div class="mt-4 flex flex-wrap items-center gap-2">
        <span class="text-xs font-semibold uppercase tracking-wide text-gray-500">Tampilkan:</span>
        <button
          v-for="series in visualSeriesOptions"
          :key="series.key"
          type="button"
          class="inline-flex items-center gap-2 rounded-lg border px-3 py-1.5 text-sm font-medium transition"
          :class="visualSeriesVisibility[series.key] ? series.activeClass : series.inactiveClass"
          @click="visualSeriesVisibility[series.key] = !visualSeriesVisibility[series.key]"
        >
          <span class="h-2.5 w-2.5 rounded-full" :style="{ backgroundColor: series.color }" />
          {{ series.label }}
        </button>
      </div>

      <div class="mt-4 grid gap-3 sm:grid-cols-3">
        <div class="rounded-lg border border-red-100 bg-red-50 px-3 py-2.5">
          <p class="text-xs font-semibold uppercase tracking-wide text-red-600">Biaya ({{ visualPeriodLabel }})</p>
          <p class="mt-1 text-sm font-semibold text-red-700">{{ formatRupiah(visualPeriodTotals.expense) }}</p>
        </div>
        <div class="rounded-lg border border-emerald-100 bg-emerald-50 px-3 py-2.5">
          <p class="text-xs font-semibold uppercase tracking-wide text-emerald-600">Pendapatan ({{ visualPeriodLabel }})</p>
          <p class="mt-1 text-sm font-semibold text-emerald-700">{{ formatRupiah(visualPeriodTotals.revenue) }}</p>
        </div>
        <div class="rounded-lg border border-blue-100 bg-blue-50 px-3 py-2.5">
          <p class="text-xs font-semibold uppercase tracking-wide text-blue-600">Net Profit ({{ visualPeriodLabel }})</p>
          <p class="mt-1 text-sm font-semibold" :class="visualPeriodTotals.net >= 0 ? 'text-blue-700' : 'text-rose-700'">
            {{ formatRupiah(visualPeriodTotals.net) }}
          </p>
        </div>
      </div>

      <div v-if="!chartHasVisibleSeries" class="mt-5 rounded-lg border border-dashed border-gray-300 px-4 py-6 text-center text-sm text-gray-500">
        Pilih minimal satu data untuk ditampilkan di grafik.
      </div>

      <div v-else class="mt-5 overflow-x-auto">
        <div class="min-w-[760px]">
          <svg
            class="h-[360px] w-full"
            :viewBox="`0 0 ${chartDimensions.width} ${chartDimensions.height}`"
            role="img"
            aria-label="Grafik garis kinerja keuangan"
          >
            <g>
              <line
                v-for="tick in chartYTicks"
                :key="`y-grid-${tick.value}`"
                :x1="chartDimensions.marginLeft"
                :x2="chartDimensions.width - chartDimensions.marginRight"
                :y1="tick.y"
                :y2="tick.y"
                stroke="#e5e7eb"
                stroke-dasharray="4 4"
              />

              <line
                v-for="tick in chartXTicks"
                :key="`x-grid-${tick.label}`"
                :x1="tick.x"
                :x2="tick.x"
                :y1="chartDimensions.marginTop"
                :y2="chartDimensions.height - chartDimensions.marginBottom"
                stroke="#f1f5f9"
              />

              <line
                :x1="chartDimensions.marginLeft"
                :x2="chartDimensions.width - chartDimensions.marginRight"
                :y1="chartZeroLineY"
                :y2="chartZeroLineY"
                stroke="#94a3b8"
                stroke-width="1.5"
              />

              <text
                v-for="tick in chartYTicks"
                :key="`y-label-${tick.value}`"
                :x="chartDimensions.marginLeft - 10"
                :y="tick.y + 4"
                text-anchor="end"
                class="fill-slate-500 text-[11px]"
              >
                {{ formatAxisRupiah(tick.value) }}
              </text>

              <text
                v-for="tick in chartXTicks"
                :key="`x-label-${tick.label}`"
                :x="tick.x"
                :y="chartDimensions.height - 18"
                text-anchor="middle"
                class="fill-slate-600 text-[11px]"
              >
                {{ tick.label }}
              </text>

              <g v-for="series in chartSeriesLines" :key="series.key">
                <path :d="series.path" fill="none" :stroke="series.color" stroke-width="3" stroke-linecap="round" stroke-linejoin="round" />
                <circle
                  v-for="(point, pointIndex) in series.points"
                  :key="`${series.key}-${pointIndex}`"
                  :cx="point.x"
                  :cy="point.y"
                  r="4.5"
                  :fill="series.color"
                />
              </g>
            </g>
          </svg>
        </div>
      </div>
    </section>

    <div class="grid gap-6 xl:grid-cols-2">
      <section class="rounded-xl border border-gray-200 bg-white p-5">
        <h3 class="mb-4 text-lg font-semibold text-gray-900">Breakdown Biaya per Kategori</h3>
        <div v-if="!expenseCategorySeries.length" class="rounded-lg border border-dashed border-gray-300 px-4 py-6 text-center text-sm text-gray-500">
          Belum ada biaya operasional untuk divisualisasikan.
        </div>
        <div v-else class="space-y-3">
          <div v-for="item in expenseCategorySeries" :key="item.category" class="space-y-1">
            <div class="flex items-center justify-between text-sm">
              <span class="font-medium text-gray-800">{{ item.category }}</span>
              <span class="text-gray-600">{{ formatRupiah(item.amount) }} ({{ item.percent.toFixed(1) }}%)</span>
            </div>
            <div class="h-2.5 rounded-full bg-gray-100">
              <div
                class="h-2.5 rounded-full bg-red-500"
                :style="{ width: `${(item.amount / maxExpenseCategoryValue) * 100}%` }"
              />
            </div>
          </div>
        </div>
      </section>

      <section class="rounded-xl border border-gray-200 bg-white p-5">
        <h3 class="mb-4 text-lg font-semibold text-gray-900">Pendapatan per Batch</h3>
        <div v-if="!revenueBatchSeries.length" class="rounded-lg border border-dashed border-gray-300 px-4 py-6 text-center text-sm text-gray-500">
          Belum ada pendapatan untuk divisualisasikan.
        </div>
        <div v-else class="space-y-3">
          <div v-for="item in revenueBatchSeries" :key="item.batchCode" class="space-y-1">
            <div class="flex items-center justify-between text-sm">
              <span class="font-medium text-gray-800">{{ item.batchCode }}</span>
              <span class="text-gray-600">{{ formatRupiah(item.amount) }}</span>
            </div>
            <div class="h-2.5 rounded-full bg-gray-100">
              <div
                class="h-2.5 rounded-full bg-emerald-500"
                :style="{ width: `${(item.amount / maxRevenueBatchValue) * 100}%` }"
              />
            </div>
          </div>
        </div>
      </section>
    </div>

    <div v-if="batchError" class="rounded-lg border border-amber-200 bg-amber-50 px-4 py-3 text-sm text-amber-800">
      {{ batchError }}
    </div>
    <div v-if="userInvestmentError" class="rounded-lg border border-amber-200 bg-amber-50 px-4 py-3 text-sm text-amber-800">
      {{ userInvestmentError }}
    </div>
    <div v-if="errorMessage" class="rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
      {{ errorMessage }}
    </div>
    <div v-if="infoMessage" class="rounded-lg border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm text-emerald-700">
      {{ infoMessage }}
    </div>

    <section class="overflow-hidden rounded-xl border border-gray-200 bg-white">
      <div class="border-b border-gray-200 px-4 py-3">
        <div class="inline-flex rounded-lg border border-gray-300 bg-gray-50 p-1">
          <button
            type="button"
            class="rounded-md px-3 py-1.5 text-sm font-medium transition"
            :class="tableView === 'expenses' ? 'bg-white text-emerald-700 shadow-sm' : 'text-gray-600 hover:text-gray-900'"
            @click="tableView = 'expenses'"
          >
            Tabel Biaya Operasional
          </button>
          <button
            type="button"
            class="rounded-md px-3 py-1.5 text-sm font-medium transition"
            :class="tableView === 'revenues' ? 'bg-white text-amber-700 shadow-sm' : 'text-gray-600 hover:text-gray-900'"
            @click="tableView = 'revenues'"
          >
            Tabel Pendapatan
          </button>
        </div>
      </div>

      <div v-if="tableView === 'expenses'">
        <div class="flex flex-wrap items-center justify-between gap-3 border-b border-gray-200 px-4 py-3">
          <h3 class="text-lg font-semibold text-gray-900">Biaya Operasional per Batch</h3>
          <div class="relative w-full max-w-sm">
            <Search class="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
            <input
              v-model="expenseSearchQuery"
              type="text"
              class="w-full rounded-lg border border-gray-300 bg-white py-2.5 pl-10 pr-3 text-sm outline-none transition focus:border-emerald-500 focus:ring-2 focus:ring-emerald-100"
              placeholder="Cari batch, kategori, catatan"
            />
          </div>
        </div>

        <div class="max-h-[460px] overflow-auto">
          <table class="w-full min-w-[980px]">
            <thead class="bg-gray-50 text-left">
              <tr>
                <th class="sticky top-0 z-10 bg-gray-50 px-4 py-3 text-xs font-semibold uppercase tracking-wide text-gray-600">Tanggal</th>
                <th class="sticky top-0 z-10 bg-gray-50 px-4 py-3 text-xs font-semibold uppercase tracking-wide text-gray-600">Batch</th>
                <th class="sticky top-0 z-10 bg-gray-50 px-4 py-3 text-xs font-semibold uppercase tracking-wide text-gray-600">Kategori</th>
                <th class="sticky top-0 z-10 bg-gray-50 px-4 py-3 text-xs font-semibold uppercase tracking-wide text-gray-600">Catatan</th>
                <th class="sticky top-0 z-10 bg-gray-50 px-4 py-3 text-xs font-semibold uppercase tracking-wide text-gray-600 text-right">Nominal</th>
                <th class="sticky top-0 z-10 bg-gray-50 px-4 py-3 text-xs font-semibold uppercase tracking-wide text-gray-600 text-right">Aksi</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200">
              <tr v-if="!filteredExpenses.length">
                <td colspan="6" class="px-4 py-6 text-center text-sm text-gray-500">
                  Belum ada biaya operasional. Tambahkan data melalui tombol input.
                </td>
              </tr>
              <tr v-for="item in paginatedExpenses" :key="item.id" class="hover:bg-gray-50">
                <td class="px-4 py-3 text-sm text-gray-700">{{ formatDate(item.date) }}</td>
                <td class="px-4 py-3 text-sm">
                  <p class="font-medium text-gray-900">{{ item.batchCode }}</p>
                  <p class="text-xs text-gray-500">{{ item.packageName }}</p>
                </td>
                <td class="px-4 py-3 text-sm text-gray-700">{{ item.category }}</td>
                <td class="px-4 py-3 text-sm text-gray-600">{{ item.note || '-' }}</td>
                <td class="px-4 py-3 text-right text-sm font-semibold text-red-600">{{ formatRupiah(item.amount) }}</td>
                <td class="px-4 py-3">
                  <div class="flex justify-end gap-2">
                    <button
                      type="button"
                      class="rounded-lg border border-amber-200 p-2 text-amber-600 transition hover:bg-amber-50"
                      title="Edit biaya"
                      @click="openEditExpenseModal(item)"
                    >
                      <Pencil class="h-4 w-4" />
                    </button>
                    <button
                      type="button"
                      class="rounded-lg border border-red-200 p-2 text-red-600 transition hover:bg-red-50"
                      title="Hapus biaya"
                      @click="askDeleteExpense(item)"
                    >
                      <Trash2 class="h-4 w-4" />
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <Pagination
          v-if="filteredExpenses.length > 0"
          :current-page="expensePage"
          :total-pages="expenseTotalPages"
          :total-rows="filteredExpenses.length"
          :limit="COST_PAGE_SIZE"
          @update:page="expensePage = $event"
        />
      </div>

      <div v-else>
        <div class="flex flex-wrap items-center justify-between gap-3 border-b border-gray-200 px-4 py-3">
          <h3 class="text-lg font-semibold text-gray-900">Pendapatan per Batch</h3>
          <div class="relative w-full max-w-sm">
            <Search class="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
            <input
              v-model="revenueSearchQuery"
              type="text"
              class="w-full rounded-lg border border-gray-300 bg-white py-2.5 pl-10 pr-3 text-sm outline-none transition focus:border-emerald-500 focus:ring-2 focus:ring-emerald-100"
              placeholder="Cari batch, sumber, catatan"
            />
          </div>
        </div>

        <div class="max-h-[460px] overflow-auto">
          <table class="w-full min-w-[980px]">
            <thead class="bg-gray-50 text-left">
              <tr>
                <th class="sticky top-0 z-10 bg-gray-50 px-4 py-3 text-xs font-semibold uppercase tracking-wide text-gray-600">Tanggal</th>
                <th class="sticky top-0 z-10 bg-gray-50 px-4 py-3 text-xs font-semibold uppercase tracking-wide text-gray-600">Batch</th>
                <th class="sticky top-0 z-10 bg-gray-50 px-4 py-3 text-xs font-semibold uppercase tracking-wide text-gray-600">Sumber</th>
                <th class="sticky top-0 z-10 bg-gray-50 px-4 py-3 text-xs font-semibold uppercase tracking-wide text-gray-600">Status</th>
                <th class="sticky top-0 z-10 bg-gray-50 px-4 py-3 text-xs font-semibold uppercase tracking-wide text-gray-600">Catatan</th>
                <th class="sticky top-0 z-10 bg-gray-50 px-4 py-3 text-xs font-semibold uppercase tracking-wide text-gray-600 text-right">Nominal</th>
                <th class="sticky top-0 z-10 bg-gray-50 px-4 py-3 text-xs font-semibold uppercase tracking-wide text-gray-600 text-right">Aksi</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200">
              <tr v-if="!filteredRevenues.length">
                <td colspan="7" class="px-4 py-6 text-center text-sm text-gray-500">
                  Belum ada data pendapatan.
                </td>
              </tr>
              <tr v-for="item in paginatedRevenues" :key="item.id" class="hover:bg-gray-50">
                <td class="px-4 py-3 text-sm text-gray-700">{{ formatDate(item.date) }}</td>
                <td class="px-4 py-3 text-sm">
                  <p class="font-medium text-gray-900">{{ item.batchCode }}</p>
                  <p class="text-xs text-gray-500">{{ item.packageName }}</p>
                </td>
                <td class="px-4 py-3 text-sm text-gray-700">{{ item.source }}</td>
                <td class="px-4 py-3 text-sm">
                  <span class="inline-flex rounded-full bg-emerald-100 px-2.5 py-1 text-xs font-semibold text-emerald-700">
                    {{ item.status }}
                  </span>
                </td>
                <td class="px-4 py-3 text-sm text-gray-600">{{ item.note || '-' }}</td>
                <td class="px-4 py-3 text-right text-sm font-semibold text-emerald-600">{{ formatRupiah(item.amount) }}</td>
              </tr>
            </tbody>
          </table>
        </div>

        <Pagination
          v-if="filteredRevenues.length > 0"
          :current-page="revenuePage"
          :total-pages="revenueTotalPages"
          :total-rows="filteredRevenues.length"
          :limit="REVENUE_PAGE_SIZE"
          @update:page="revenuePage = $event"
        />
      </div>
    </section>

    <div v-if="expenseModalOpen" class="fixed inset-0 z-50 bg-gray-900/45 p-4">
      <div class="mx-auto mt-10 w-full max-w-2xl rounded-xl border border-gray-200 bg-white p-6 shadow-xl">
        <div class="flex items-start justify-between gap-3">
          <div>
            <h3 class="text-xl font-semibold text-gray-900">Input Biaya Operasional</h3>
            <p class="mt-1 text-sm text-gray-600">Catat biaya yang terjadi pada batch tanaman tertentu.</p>
          </div>
          <button
            type="button"
            class="rounded-lg border border-gray-200 p-2 text-gray-600 transition hover:bg-gray-100"
            @click="closeExpenseModal"
          >
            <X class="h-4 w-4" />
          </button>
        </div>

        <form class="mt-5 grid gap-4 md:grid-cols-2" @submit.prevent="submitExpense">
          <label class="space-y-2">
            <span class="text-sm font-medium text-gray-700">Batch tanaman</span>
            <select
              v-model="expenseForm.plantBatchId"
              class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2.5 text-sm outline-none transition focus:border-emerald-500 focus:ring-2 focus:ring-emerald-100"
            >
              <option v-for="batch in batches" :key="batch.id" :value="String(batch.id)">
                {{ batch.batchCode }} - {{ batch.packageName }}
              </option>
            </select>
          </label>

          <label class="space-y-2">
            <span class="text-sm font-medium text-gray-700">Kategori biaya</span>
            <select
              v-model="expenseForm.category"
              class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2.5 text-sm outline-none transition focus:border-emerald-500 focus:ring-2 focus:ring-emerald-100"
            >
              <option v-for="item in costCategoryOptions" :key="item" :value="item">
                {{ item }}
              </option>
            </select>
          </label>

          <label class="space-y-2">
            <span class="text-sm font-medium text-gray-700">Nominal biaya</span>
            <input
              v-model.number="expenseForm.amount"
              type="number"
              min="0"
              step="1000"
              class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2.5 text-sm outline-none transition focus:border-emerald-500 focus:ring-2 focus:ring-emerald-100"
              placeholder="Contoh: 250000"
            />
          </label>

          <label class="space-y-2">
            <span class="text-sm font-medium text-gray-700">Tanggal biaya</span>
            <input
              v-model="expenseForm.date"
              type="date"
              class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2.5 text-sm outline-none transition focus:border-emerald-500 focus:ring-2 focus:ring-emerald-100"
            />
          </label>

          <label class="space-y-2 md:col-span-2">
            <span class="text-sm font-medium text-gray-700">Catatan (opsional)</span>
            <textarea
              v-model="expenseForm.note"
              rows="3"
              class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2.5 text-sm outline-none transition focus:border-emerald-500 focus:ring-2 focus:ring-emerald-100"
              placeholder="Contoh: pembelian pupuk organik untuk minggu ke-2"
            />
          </label>

          <div class="md:col-span-2 flex justify-end gap-2">
            <button
              type="button"
              class="rounded-lg border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50"
              @click="closeExpenseModal"
            >
              Batal
            </button>
            <button
              type="submit"
              :disabled="savingExpense"
              class="inline-flex items-center gap-2 rounded-lg bg-emerald-600 px-4 py-2 text-sm font-medium text-white hover:bg-emerald-700 disabled:opacity-70"
            >
              <Wallet class="h-4 w-4" />
              {{ savingExpense ? "Menyimpan..." : "Simpan Biaya" }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Edit Expense Modal -->
    <div v-if="editExpenseModalOpen" class="fixed inset-0 z-50 bg-gray-900/45 p-4">
      <div class="mx-auto mt-10 w-full max-w-2xl rounded-xl border border-gray-200 bg-white p-6 shadow-xl">
        <div class="flex items-start justify-between gap-3">
          <div>
            <h3 class="text-xl font-semibold text-gray-900">Edit Biaya Operasional</h3>
            <p class="mt-1 text-sm text-gray-600">Perbarui data biaya operasional yang sudah dicatat.</p>
          </div>
          <button
            type="button"
            class="rounded-lg border border-gray-200 p-2 text-gray-600 transition hover:bg-gray-100"
            @click="closeEditExpenseModal"
          >
            <X class="h-4 w-4" />
          </button>
        </div>

        <form class="mt-5 grid gap-4 md:grid-cols-2" @submit.prevent="submitEditExpense">
          <label class="space-y-2">
            <span class="text-sm font-medium text-gray-700">Batch tanaman</span>
            <select
              v-model="editExpenseForm.plantBatchId"
              class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2.5 text-sm outline-none transition focus:border-amber-500 focus:ring-2 focus:ring-amber-100"
            >
              <option v-for="batch in batches" :key="batch.id" :value="String(batch.id)">
                {{ batch.batchCode }} - {{ batch.packageName }}
              </option>
            </select>
          </label>

          <label class="space-y-2">
            <span class="text-sm font-medium text-gray-700">Kategori biaya</span>
            <select
              v-model="editExpenseForm.category"
              class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2.5 text-sm outline-none transition focus:border-amber-500 focus:ring-2 focus:ring-amber-100"
            >
              <option v-for="item in costCategoryOptions" :key="item" :value="item">
                {{ item }}
              </option>
            </select>
          </label>

          <label class="space-y-2">
            <span class="text-sm font-medium text-gray-700">Nominal biaya</span>
            <input
              v-model.number="editExpenseForm.amount"
              type="number"
              min="0"
              step="1000"
              class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2.5 text-sm outline-none transition focus:border-amber-500 focus:ring-2 focus:ring-amber-100"
              placeholder="Contoh: 250000"
            />
          </label>

          <label class="space-y-2">
            <span class="text-sm font-medium text-gray-700">Tanggal biaya</span>
            <input
              v-model="editExpenseForm.date"
              type="date"
              class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2.5 text-sm outline-none transition focus:border-amber-500 focus:ring-2 focus:ring-amber-100"
            />
          </label>

          <label class="space-y-2 md:col-span-2">
            <span class="text-sm font-medium text-gray-700">Catatan (opsional)</span>
            <textarea
              v-model="editExpenseForm.note"
              rows="3"
              class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2.5 text-sm outline-none transition focus:border-amber-500 focus:ring-2 focus:ring-amber-100"
              placeholder="Contoh: pembelian pupuk organik untuk minggu ke-2"
            />
          </label>

          <p v-if="errorMessage" class="md:col-span-2 rounded-lg border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700">{{ errorMessage }}</p>

          <div class="md:col-span-2 flex justify-end gap-2">
            <button
              type="button"
              class="rounded-lg border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50"
              @click="closeEditExpenseModal"
            >
              Batal
            </button>
            <button
              type="submit"
              :disabled="savingEditExpense"
              class="inline-flex items-center gap-2 rounded-lg bg-amber-500 px-4 py-2 text-sm font-medium text-white hover:bg-amber-600 disabled:opacity-70"
            >
              <Pencil class="h-4 w-4" />
              {{ savingEditExpense ? "Menyimpan..." : "Perbarui Biaya" }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <div v-if="deleteModalOpen" class="fixed inset-0 z-50 bg-gray-900/45 p-4">
      <div class="mx-auto mt-20 w-full max-w-md rounded-xl border border-gray-200 bg-white p-6 shadow-xl">
        <h3 class="text-lg font-semibold text-gray-900">Konfirmasi Hapus</h3>
        <p class="mt-2 text-sm text-gray-600">
          Yakin ingin menghapus
          <span class="font-semibold text-gray-900">{{ deleteTarget?.label }}</span>
          ?
        </p>
        <p class="mt-1 text-xs text-gray-500">Tindakan ini tidak dapat dibatalkan.</p>

        <div class="mt-5 flex justify-end gap-2">
          <button
            type="button"
            class="rounded-lg border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50"
            :disabled="deleteSubmitting"
            @click="closeDeleteModal"
          >
            Batal
          </button>
          <button
            type="button"
            class="rounded-lg bg-red-600 px-4 py-2 text-sm font-medium text-white hover:bg-red-700 disabled:opacity-60"
            :disabled="deleteSubmitting"
            @click="confirmDelete"
          >
            {{ deleteSubmitting ? "Menghapus..." : "Ya, Hapus" }}
          </button>
        </div>
      </div>
    </div>

    <div v-if="projectionModalOpen" class="fixed inset-0 z-50 overflow-y-auto bg-gray-900/45 p-4">
      <div class="mx-auto mt-10 w-full max-w-4xl rounded-xl border border-gray-200 bg-white p-6 shadow-xl">
        <div class="flex items-start justify-between gap-3">
          <div>
            <h3 class="text-xl font-semibold text-gray-900">Proyeksi Pendapatan 6 Tahun</h3>
            <p class="mt-1 text-sm text-gray-600">
              Masukkan parameter estimasi panen berdasarkan tabel perhitungan yang tersedia (tahun 1-6).
            </p>
          </div>
          <div class="flex items-center gap-2">
            <button
              type="button"
              class="inline-flex items-center gap-2 rounded-lg border border-violet-200 bg-violet-50 px-3 py-2 text-sm font-medium text-violet-700 transition hover:bg-violet-100"
              title="Lihat formula perhitungan"
              @click="projectionFormulaModalOpen = true"
            >
              <Info class="h-4 w-4" />
              Formula
            </button>
            <button
              type="button"
              class="rounded-lg border border-gray-200 p-2 text-gray-600 transition hover:bg-gray-100"
              @click="closeProjectionModal"
            >
              <X class="h-4 w-4" />
            </button>
          </div>
        </div>

        <!-- Form Section -->
        <form class="mt-5 space-y-4" @submit.prevent="submitProjection">
          <!-- Jumlah Pohon -->
          <label class="space-y-2">
            <span class="text-sm font-medium text-gray-700">Jumlah Pohon</span>
            <input
              v-model.number="projectionInput.plantCount"
              type="number"
              min="1"
              step="1"
              class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2.5 text-sm outline-none transition focus:border-violet-500 focus:ring-2 focus:ring-violet-100"
              placeholder="Contoh: 100"
            />
          </label>

          <!-- Provinsi -->
          <label class="space-y-2">
            <span class="text-sm font-medium text-gray-700">Provinsi (untuk referensi harga)</span>
            <select
              v-model="projectionInput.selectedProvince"
              class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2.5 text-sm outline-none transition focus:border-violet-500 focus:ring-2 focus:ring-violet-100"
            >
              <option value="">-- Pilih Provinsi --</option>
              <option v-for="province in availableProvinces" :key="province" :value="province">
                {{ province }}
              </option>
            </select>
          </label>

          <!-- Grade Percentages -->
          <div class="rounded-lg border border-gray-200 bg-gray-50 p-4">
            <div class="mb-4 flex flex-wrap items-start justify-between gap-3">
              <div>
                <p class="text-sm font-medium text-gray-700">Distribusi Hasil Panen per Grade</p>
                <p class="mt-1 text-xs text-amber-700">
                  Catatan: rata-rata hasil Grade A biasanya berada di kisaran 20-30% dari total panen.
                </p>
              </div>
              <div class="rounded-lg border border-white bg-white px-4 py-2 text-right shadow-sm">
                <p class="text-xs font-medium uppercase tracking-wide text-gray-500">Total</p>
                <p
                  class="text-2xl font-semibold"
                  :class="{
                    'text-emerald-600': isProjectionGradeTotalValid,
                    'text-red-600': !isProjectionGradeTotalValid,
                  }"
                >
                  {{ projectionGradeTotalPercent.toFixed(1) }}%
                </p>
              </div>
            </div>
            <div v-if="!grades.length" class="text-xs text-gray-500">
              Tidak ada data grade yang tersedia. Silakan refresh data.
            </div>
            <div v-else class="grid gap-3 md:grid-cols-2">
              <label v-for="grade in grades" :key="grade.id" class="space-y-1">
                <span class="text-sm font-medium text-gray-600">{{ grade.grade_name }} (%)</span>
                <input
                  v-model.number="projectionInput.gradePercentages[grade.grade_name]"
                  type="number"
                  min="0"
                  max="100"
                  step="0.1"
                  class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2.5 text-lg font-semibold outline-none transition focus:border-violet-500 focus:ring-2 focus:ring-violet-100"
                  placeholder="0"
                  @input="autofillRemainingGradePercentage(grade.grade_name)"
                />
              </label>
            </div>
            <div class="mt-3 text-sm font-medium" :class="isProjectionGradeTotalValid ? 'text-emerald-600' : 'text-red-600'">
              {{ isProjectionGradeTotalValid ? "Distribusi sudah 100%." : "Total distribusi harus tepat 100% sebelum proyeksi dihitung." }}
            </div>
          </div>

          <!-- Submit Button -->
          <div class="flex justify-end">
            <button
              type="submit"
              :disabled="projectingRevenue || !canSubmitProjection"
              class="inline-flex items-center gap-2 rounded-lg bg-violet-600 px-4 py-2.5 text-sm font-medium text-white transition hover:bg-violet-700 disabled:cursor-not-allowed disabled:bg-gray-300 disabled:text-gray-500"
            >
              <TrendingUp class="h-4 w-4" />
              {{ projectingRevenue ? "Menghitung..." : "Proyeksikan" }}
            </button>
          </div>
        </form>

        <!-- Results Section -->
        <div v-if="projectionResults.length" class="mt-6 space-y-6 border-t border-gray-200 pt-6">
          <!-- Statistics -->
          <div class="grid gap-4 md:grid-cols-2">
            <div class="rounded-lg border border-violet-100 bg-violet-50 p-4">
              <p class="text-xs font-semibold uppercase tracking-wide text-violet-600">Total Panen (6 Tahun)</p>
              <p class="mt-2 text-lg font-semibold text-violet-900">{{ projectionStats.totalHarvestKg.toLocaleString("id-ID") }} kg</p>
            </div>
            <div class="rounded-lg border border-emerald-100 bg-emerald-50 p-4">
              <p class="text-xs font-semibold uppercase tracking-wide text-emerald-600">Rata-rata Pendapatan/Tahun</p>
              <p class="mt-2 text-xs text-emerald-700">Min: Rp{{ Math.round(projectionStats.avgYearlyMin).toLocaleString("id-ID") }}</p>
              <p class="text-xs text-emerald-700">Max: Rp{{ Math.round(projectionStats.avgYearlyMax).toLocaleString("id-ID") }}</p>
            </div>
            <div class="rounded-lg border border-blue-100 bg-blue-50 p-4">
              <p class="text-xs font-semibold uppercase tracking-wide text-blue-600">Total Min Revenue (6 Tahun)</p>
              <p class="mt-2 text-lg font-semibold text-blue-900">Rp{{ Math.round(projectionStats.totalMinRevenue).toLocaleString("id-ID") }}</p>
            </div>
            <div class="rounded-lg border border-amber-100 bg-amber-50 p-4">
              <p class="text-xs font-semibold uppercase tracking-wide text-amber-600">Total Max Revenue (6 Tahun)</p>
              <p class="mt-2 text-lg font-semibold text-amber-900">Rp{{ Math.round(projectionStats.totalMaxRevenue).toLocaleString("id-ID") }}</p>
            </div>
          </div>

          <!-- Chart -->
          <div class="rounded-lg border border-gray-200 bg-gray-50 p-4">
            <p class="mb-3 text-sm font-medium text-gray-700">Grafik Range Proyeksi Pendapatan</p>
            <div class="overflow-x-auto">
              <div class="min-w-[400px]">
                <svg
                  class="h-80 w-full"
                  :viewBox="`0 0 ${projectionChartDimensions.width} ${projectionChartDimensions.height}`"
                  role="img"
                  aria-label="Grafik range proyeksi pendapatan"
                >
                  <g>
                    <!-- Background grid -->
                    <line
                      v-for="tick in projectionChartYTicks"
                      :key="`y-grid-${tick.value}`"
                      :x1="projectionChartDimensions.marginLeft"
                      :x2="projectionChartDimensions.width - projectionChartDimensions.marginRight"
                      :y1="tick.y"
                      :y2="tick.y"
                      stroke="#e5e7eb"
                      stroke-dasharray="4 4"
                    />

                    <!-- Zero line -->
                    <line
                      :x1="projectionChartDimensions.marginLeft"
                      :x2="projectionChartDimensions.width - projectionChartDimensions.marginRight"
                      :y1="projectionChartZeroLineY"
                      :y2="projectionChartZeroLineY"
                      stroke="#94a3b8"
                      stroke-width="1.5"
                    />

                    <!-- Y-axis labels -->
                    <text
                      v-for="tick in projectionChartYTicks"
                      :key="`y-label-${tick.value}`"
                      :x="projectionChartDimensions.marginLeft - 10"
                      :y="tick.y + 4"
                      text-anchor="end"
                      class="fill-slate-500 text-[10px]"
                    >
                      Rp{{ Math.round(tick.value / 1000000) }}M
                    </text>

                    <!-- X-axis labels -->
                    <text
                      v-for="tick in projectionChartXTicks"
                      :key="`x-label-${tick.label}`"
                      :x="tick.x"
                      :y="projectionChartDimensions.height - 18"
                      text-anchor="middle"
                      class="fill-slate-600 text-[10px]"
                    >
                      {{ tick.label }}
                    </text>

                    <!-- Range areas and lines -->
                    <g v-for="(result, idx) in projectionResults" :key="`range-${idx}`">
                      <!-- Range polygon (gradient area) -->
                      <line
                        :x1="getProjectionChartX(idx, projectionResults.length)"
                        :x2="getProjectionChartX(idx, projectionResults.length)"
                        :y1="getProjectionChartY(result.minRevenue)"
                        :y2="getProjectionChartY(result.maxRevenue)"
                        stroke="#3b82f6"
                        stroke-width="2"
                        opacity="0.6"
                      />

                      <!-- Min point (circle) -->
                      <circle
                        :cx="getProjectionChartX(idx, projectionResults.length)"
                        :cy="getProjectionChartY(result.minRevenue)"
                        r="3"
                        fill="#ef4444"
                      />

                      <!-- Max point (circle) -->
                      <circle
                        :cx="getProjectionChartX(idx, projectionResults.length)"
                        :cy="getProjectionChartY(result.maxRevenue)"
                        r="3"
                        fill="#10b981"
                      />
                    </g>

                    <!-- Legend markers -->
                    <circle cx="20" cy="20" r="2.5" fill="#ef4444" />
                    <text x="28" y="24" class="fill-slate-600 text-[10px]">Min</text>

                    <circle cx="55" cy="20" r="2.5" fill="#10b981" />
                    <text x="63" y="24" class="fill-slate-600 text-[10px]">Max</text>
                  </g>
                </svg>
              </div>
            </div>
          </div>

          <!-- Results Table -->
          <div class="overflow-x-auto rounded-lg border border-gray-200">
            <table class="w-full min-w-[500px] text-left text-sm">
              <thead class="border-b border-gray-200 bg-gray-50">
                <tr>
                  <th class="px-4 py-2 font-semibold text-gray-700">Tahun</th>
                  <th class="px-4 py-2 text-right font-semibold text-gray-700">Panen (kg)</th>
                  <th class="px-4 py-2 text-right font-semibold text-gray-700">Min Revenue</th>
                  <th class="px-4 py-2 text-right font-semibold text-gray-700">Max Revenue</th>
                  <th class="px-4 py-2 text-right font-semibold text-gray-700">Range</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-200">
                <tr v-for="result in projectionResults" :key="result.year" class="hover:bg-gray-50">
                  <td class="px-4 py-2 font-medium text-gray-900">{{ result.year }}</td>
                  <td class="px-4 py-2 text-right text-gray-700">{{ result.harvestKg.toLocaleString("id-ID") }}</td>
                  <td class="px-4 py-2 text-right text-red-600">
                    Rp{{ Math.round(result.minRevenue).toLocaleString("id-ID") }}
                  </td>
                  <td class="px-4 py-2 text-right text-emerald-600">
                    Rp{{ Math.round(result.maxRevenue).toLocaleString("id-ID") }}
                  </td>
                  <td class="px-4 py-2 text-right text-blue-600">
                    Rp{{ Math.round(result.maxRevenue - result.minRevenue).toLocaleString("id-ID") }}
                  </td>
                </tr>
                <tr class="bg-violet-50 font-semibold text-violet-900">
                  <td class="px-4 py-2">Total</td>
                  <td class="px-4 py-2 text-right">{{ projectionStats.totalHarvestKg.toLocaleString("id-ID") }} kg</td>
                  <td class="px-4 py-2 text-right">Rp{{ Math.round(projectionStats.totalMinRevenue).toLocaleString("id-ID") }}</td>
                  <td class="px-4 py-2 text-right">Rp{{ Math.round(projectionStats.totalMaxRevenue).toLocaleString("id-ID") }}</td>
                  <td class="px-4 py-2 text-right">
                    Rp{{ Math.round(projectionStats.totalMaxRevenue - projectionStats.totalMinRevenue).toLocaleString("id-ID") }}
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Closing button -->
        <div class="mt-6 flex justify-end gap-2 border-t border-gray-200 pt-4">
          <button
            type="button"
            class="rounded-lg border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50"
            @click="closeProjectionModal"
          >
            Tutup
          </button>
        </div>
      </div>

      <div v-if="projectionFormulaModalOpen" class="fixed inset-0 z-[60] flex items-center justify-center bg-gray-900/50 p-4">
        <div class="w-full max-w-2xl rounded-xl border border-gray-200 bg-white p-6 shadow-xl">
          <div class="flex items-start justify-between gap-3">
            <div>
              <h4 class="text-lg font-semibold text-gray-900">Formula Proyeksi Pendapatan</h4>
              <p class="mt-1 text-sm text-gray-600">
                Perhitungan menggunakan estimasi panen per pohon, distribusi grade, dan range harga nasional pada provinsi terpilih.
              </p>
            </div>
            <button
              type="button"
              class="rounded-lg border border-gray-200 p-2 text-gray-600 transition hover:bg-gray-100"
              @click="projectionFormulaModalOpen = false"
            >
              <X class="h-4 w-4" />
            </button>
          </div>

          <div class="mt-5 space-y-4 text-sm text-gray-700">
            <div class="rounded-lg border border-gray-200 bg-gray-50 p-4">
              <p class="font-semibold text-gray-900">1. Estimasi panen per tahun</p>
              <p class="mt-2">Total panen tahun ke-n = jumlah pohon x estimasi kg per pohon.</p>
              <div class="mt-3 grid grid-cols-2 gap-2 text-xs sm:grid-cols-3">
                <div class="rounded-md bg-white px-3 py-2">Tahun 1: 0 kg/pohon</div>
                <div class="rounded-md bg-white px-3 py-2">Tahun 2: 0 kg/pohon</div>
                <div class="rounded-md bg-white px-3 py-2">Tahun 3: 0,25 kg/pohon</div>
                <div class="rounded-md bg-white px-3 py-2">Tahun 4: 0,5 kg/pohon</div>
                <div class="rounded-md bg-white px-3 py-2">Tahun 5: 1 kg/pohon</div>
                <div class="rounded-md bg-white px-3 py-2">Tahun 6: 2 kg/pohon</div>
              </div>
            </div>

            <div class="rounded-lg border border-gray-200 bg-gray-50 p-4">
              <p class="font-semibold text-gray-900">2. Pembagian hasil panen per grade</p>
              <p class="mt-2">Panen per grade = total panen x persentase grade.</p>
              <p class="mt-1 text-xs text-gray-500">Total persentase semua grade harus tepat 100% agar proyeksi dapat dihitung.</p>
            </div>

            <div class="rounded-lg border border-gray-200 bg-gray-50 p-4">
              <p class="font-semibold text-gray-900">3. Pendapatan minimum dan maksimum</p>
              <p class="mt-2">Pendapatan minimum = jumlah seluruh panen per grade x harga minimum per kg.</p>
              <p class="mt-1">Pendapatan maksimum = jumlah seluruh panen per grade x harga maksimum per kg.</p>
              <p class="mt-1 text-xs text-gray-500">Harga minimum dan maksimum diambil dari data harga nasional untuk provinsi yang dipilih.</p>
            </div>

            <div class="rounded-lg border border-violet-100 bg-violet-50 p-4 text-violet-900">
              <p class="font-semibold">Rumus ringkas</p>
              <p class="mt-2">
                Revenue per tahun = Σ((jumlah pohon x estimasi kg/pohon tahun ke-n) x persentase grade x harga per kg grade).
              </p>
            </div>
          </div>

          <div class="mt-6 flex justify-end border-t border-gray-200 pt-4">
            <button
              type="button"
              class="rounded-lg bg-violet-600 px-4 py-2 text-sm font-medium text-white transition hover:bg-violet-700"
              @click="projectionFormulaModalOpen = false"
            >
              Mengerti
            </button>
          </div>
        </div>
      </div>
    </div>

    <div v-if="!batches.length && !loadingBatches" class="rounded-xl border border-dashed border-gray-300 bg-gray-50 p-5 text-sm text-gray-600">
      Belum ada data batch tanaman yang bisa dipakai untuk tracker keuangan. Silakan pastikan batch sudah tervalidasi dan memiliki data monitoring.
    </div>
  </div>
</template>
