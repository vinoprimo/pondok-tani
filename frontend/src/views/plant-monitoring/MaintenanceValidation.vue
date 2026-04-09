<script setup lang="ts">
import { ref, computed } from 'vue';
import { CheckCircle, XCircle, Eye, User, Calendar } from 'lucide-vue-next';

const selectedActivity = ref<string | null>(null);
const showReviewModal = ref(false);
const rejectionReason = ref('');

const submittedActivities = [
  {
    id: 'MA-003',
    userName: 'Sarah Johnson',
    userId: 'INV-001',
    batchId: 'BATCH-045',
    activityType: 'Pruning',
    submissionDate: '2024-01-23',
    activityDate: '2024-01-23',
    description:
      'Performed regular pruning of vanilla vines, removed dead leaves and excess shoots to promote better flowering.',
    quantity: '15 vines pruned',
    notes: 'All tools sanitized before use',
    photoUrl: 'https://images.unsplash.com/photo-1464226184884-fa280b87c399?w=800',
    status: 'Pending',
  },
  {
    id: 'MA-006',
    userName: 'Michael Chen',
    userId: 'INV-002',
    batchId: 'BATCH-038',
    activityType: 'Fertilizing',
    submissionDate: '2024-01-23',
    activityDate: '2024-01-22',
    description: 'Applied organic fertilizer to all vanilla plants. Used NPK 15-15-15 formula.',
    quantity: '25 kg fertilizer',
    notes: 'Weather was optimal, no rain expected for 24 hours',
    photoUrl: 'https://images.unsplash.com/photo-1416879595882-3373a0480b5b?w=800',
    status: 'Pending',
  },
  {
    id: 'MA-007',
    userName: 'Emma Davis',
    userId: 'INV-003',
    batchId: 'BATCH-052',
    activityType: 'Watering',
    submissionDate: '2024-01-22',
    activityDate: '2024-01-22',
    description: 'Deep watering session for all plants in the morning.',
    quantity: '500 liters',
    notes: 'Soil moisture checked before watering',
    photoUrl: 'https://images.unsplash.com/photo-1523348837708-15d4a09cfac2?w=800',
    status: 'Pending',
  },
];

const selectedActivityData = computed(() =>
  submittedActivities.find((a) => a.id === selectedActivity.value)
);

function handleReview(activityId: string) {
  selectedActivity.value = activityId;
  showReviewModal.value = true;
}

function handleApprove() {
  alert(`Activity ${selectedActivity.value} approved successfully!`);
  showReviewModal.value = false;
  selectedActivity.value = null;
}

function handleReject() {
  if (!rejectionReason.value.trim()) {
    alert('Please provide a rejection reason');
    return;
  }
  alert(`Activity ${selectedActivity.value} rejected. Reason: ${rejectionReason.value}`);
  showReviewModal.value = false;
  selectedActivity.value = null;
  rejectionReason.value = '';
}

function closeReviewModal() {
  showReviewModal.value = false;
  rejectionReason.value = '';
}
</script>

<template>
  <div class="space-y-6">
    <div>
      <h2 class="text-2xl font-semibold text-gray-900">Maintenance Validation</h2>
      <p class="text-gray-600 mt-1">
        Review and validate maintenance activities submitted by investors
      </p>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center justify-between mb-2">
          <p class="text-sm text-gray-600">Pending Review</p>
          <Calendar class="w-5 h-5 text-yellow-600" />
        </div>
        <p class="text-2xl font-semibold text-yellow-600">3</p>
        <p class="text-xs text-gray-500 mt-1">Requires attention</p>
      </div>

      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center justify-between mb-2">
          <p class="text-sm text-gray-600">Approved Today</p>
          <CheckCircle class="w-5 h-5 text-green-600" />
        </div>
        <p class="text-2xl font-semibold text-gray-900">12</p>
        <p class="text-xs text-green-600 mt-1">+8 from yesterday</p>
      </div>

      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center justify-between mb-2">
          <p class="text-sm text-gray-600">Rejected Today</p>
          <XCircle class="w-5 h-5 text-red-600" />
        </div>
        <p class="text-2xl font-semibold text-gray-900">2</p>
        <p class="text-xs text-red-600 mt-1">Quality issues</p>
      </div>

      <div class="bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center justify-between mb-2">
          <p class="text-sm text-gray-600">Total This Month</p>
          <Calendar class="w-5 h-5 text-blue-600" />
        </div>
        <p class="text-2xl font-semibold text-gray-900">87</p>
        <p class="text-xs text-gray-500 mt-1">Activities reviewed</p>
      </div>
    </div>

    <div class="bg-white rounded-xl border border-gray-200 overflow-hidden">
      <div class="px-6 py-4 border-b border-gray-200">
        <h3 class="text-lg font-semibold text-gray-900">Pending Activities</h3>
        <p class="text-sm text-gray-600 mt-1">Review submitted maintenance activities</p>
      </div>
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-gray-50 border-b border-gray-200">
            <tr>
              <th class="text-left px-6 py-3 text-sm font-medium text-gray-900">Activity ID</th>
              <th class="text-left px-6 py-3 text-sm font-medium text-gray-900">User</th>
              <th class="text-left px-6 py-3 text-sm font-medium text-gray-900">Batch</th>
              <th class="text-left px-6 py-3 text-sm font-medium text-gray-900">Activity Type</th>
              <th class="text-left px-6 py-3 text-sm font-medium text-gray-900">Activity Date</th>
              <th class="text-left px-6 py-3 text-sm font-medium text-gray-900">Submission Date</th>
              <th class="text-center px-6 py-3 text-sm font-medium text-gray-900">Photo</th>
              <th class="text-center px-6 py-3 text-sm font-medium text-gray-900">Status</th>
              <th class="text-center px-6 py-3 text-sm font-medium text-gray-900">Action</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-for="activity in submittedActivities" :key="activity.id" class="hover:bg-gray-50">
              <td class="px-6 py-4 font-medium text-gray-900">{{ activity.id }}</td>
              <td class="px-6 py-4">
                <div class="flex items-center gap-2">
                  <div class="w-8 h-8 bg-green-100 rounded-full flex items-center justify-center">
                    <User class="w-4 h-4 text-green-600" />
                  </div>
                  <div>
                    <p class="font-medium text-gray-900">{{ activity.userName }}</p>
                    <p class="text-xs text-gray-500">{{ activity.userId }}</p>
                  </div>
                </div>
              </td>
              <td class="px-6 py-4 text-gray-900">{{ activity.batchId }}</td>
              <td class="px-6 py-4">
                <span
                  class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800"
                >
                  {{ activity.activityType }}
                </span>
              </td>
              <td class="px-6 py-4 text-gray-900">{{ activity.activityDate }}</td>
              <td class="px-6 py-4 text-gray-900">{{ activity.submissionDate }}</td>
              <td class="px-6 py-4 text-center">
                <div class="flex justify-center">
                  <img
                    :src="activity.photoUrl"
                    alt="Activity preview"
                    class="w-12 h-12 rounded-lg object-cover cursor-pointer hover:opacity-75 transition-opacity"
                    @click="handleReview(activity.id)"
                  />
                </div>
              </td>
              <td class="px-6 py-4 text-center">
                <span
                  class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800"
                >
                  Pending
                </span>
              </td>
              <td class="px-6 py-4 text-center">
                <button
                  type="button"
                  class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors text-sm font-medium flex items-center gap-2 mx-auto"
                  @click="handleReview(activity.id)"
                >
                  <Eye class="w-4 h-4" />
                  Review
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div
      v-if="showReviewModal && selectedActivityData"
      class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
    >
      <div class="bg-white rounded-2xl max-w-4xl w-full max-h-[90vh] overflow-y-auto">
        <div class="sticky top-0 bg-white border-b border-gray-200 px-6 py-4 flex items-center justify-between">
          <div>
            <h3 class="text-xl font-semibold text-gray-900">
              Review Activity - {{ selectedActivityData.id }}
            </h3>
            <p class="text-sm text-gray-600 mt-1">
              Submitted by {{ selectedActivityData.userName }}
            </p>
          </div>
          <button
            type="button"
            class="w-8 h-8 rounded-lg hover:bg-gray-100 flex items-center justify-center transition-colors"
            @click="closeReviewModal"
          >
            <XCircle class="w-5 h-5 text-gray-500" />
          </button>
        </div>

        <div class="p-6 space-y-6">
          <div>
            <h4 class="text-sm font-medium text-gray-700 mb-3">Photo Evidence</h4>
            <div class="border border-gray-200 rounded-xl overflow-hidden">
              <img
                :src="selectedActivityData.photoUrl"
                alt="Activity evidence"
                class="w-full h-96 object-cover"
              />
            </div>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div class="bg-gray-50 rounded-lg p-4">
              <p class="text-sm text-gray-600 mb-1">Activity Type</p>
              <p class="font-semibold text-gray-900">{{ selectedActivityData.activityType }}</p>
            </div>
            <div class="bg-gray-50 rounded-lg p-4">
              <p class="text-sm text-gray-600 mb-1">Batch ID</p>
              <p class="font-semibold text-gray-900">{{ selectedActivityData.batchId }}</p>
            </div>
            <div class="bg-gray-50 rounded-lg p-4">
              <p class="text-sm text-gray-600 mb-1">Activity Date</p>
              <p class="font-semibold text-gray-900">{{ selectedActivityData.activityDate }}</p>
            </div>
            <div class="bg-gray-50 rounded-lg p-4">
              <p class="text-sm text-gray-600 mb-1">Submission Date</p>
              <p class="font-semibold text-gray-900">{{ selectedActivityData.submissionDate }}</p>
            </div>
          </div>

          <div>
            <h4 class="text-sm font-medium text-gray-700 mb-2">Description</h4>
            <div class="bg-gray-50 rounded-lg p-4">
              <p class="text-gray-900">{{ selectedActivityData.description }}</p>
            </div>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div>
              <h4 class="text-sm font-medium text-gray-700 mb-2">Quantity Used</h4>
              <div class="bg-gray-50 rounded-lg p-4">
                <p class="text-gray-900">{{ selectedActivityData.quantity }}</p>
              </div>
            </div>
            <div>
              <h4 class="text-sm font-medium text-gray-700 mb-2">Additional Notes</h4>
              <div class="bg-gray-50 rounded-lg p-4">
                <p class="text-gray-900">{{ selectedActivityData.notes }}</p>
              </div>
            </div>
          </div>

          <div>
            <h4 class="text-sm font-medium text-gray-700 mb-2">Rejection Reason (Optional)</h4>
            <textarea
              v-model="rejectionReason"
              rows="3"
              placeholder="If rejecting, please provide a clear reason..."
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>

          <div class="flex gap-3 pt-4">
            <button
              type="button"
              class="flex-1 px-6 py-3 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors font-medium flex items-center justify-center gap-2"
              @click="handleReject"
            >
              <XCircle class="w-5 h-5" />
              Reject Activity
            </button>
            <button
              type="button"
              class="flex-1 px-6 py-3 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors font-medium flex items-center justify-center gap-2"
              @click="handleApprove"
            >
              <CheckCircle class="w-5 h-5" />
              Approve Activity
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
