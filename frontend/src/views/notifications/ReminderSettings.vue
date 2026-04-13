<script setup lang="ts">
import { ref, reactive } from 'vue';
import { Bell, Mail, Smartphone, Clock, CheckCircle } from 'lucide-vue-next';

type ReminderTiming = '3days' | '1day' | 'sameday';

const settings = reactive({
  emailReminders: true,
  pushNotifications: false,
  reminderTiming: '1day' as ReminderTiming,
});

const saved = ref(false);

let savedTimer: ReturnType<typeof setTimeout> | undefined;

function handleSave() {
  saved.value = true;
  if (savedTimer) clearTimeout(savedTimer);
  savedTimer = setTimeout(() => {
    saved.value = false;
  }, 3000);
}
</script>

<template>
  <div class="space-y-6">
    <div>
      <h2 class="text-2xl font-semibold text-gray-900">Reminder Settings</h2>
      <p class="text-gray-600 mt-1">Configure your maintenance activity reminders</p>
    </div>

    <div class="bg-white rounded-xl border border-gray-200 p-6 space-y-6">
      <div class="flex items-start justify-between pb-6 border-b border-gray-200">
        <div class="flex items-start gap-4">
          <div class="w-12 h-12 bg-blue-50 rounded-lg flex items-center justify-center flex-shrink-0">
            <Mail class="w-6 h-6 text-blue-600" />
          </div>
          <div>
            <h3 class="text-lg font-semibold text-gray-900">Email Reminders</h3>
            <p class="text-sm text-gray-600 mt-1">
              Receive email notifications for upcoming maintenance tasks
            </p>
          </div>
        </div>
        <label class="relative inline-flex items-center cursor-pointer">
          <input
            v-model="settings.emailReminders"
            type="checkbox"
            class="sr-only peer"
          />
          <div
            class="w-14 h-7 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-green-300 rounded-full peer peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-0.5 after:start-[4px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-6 after:w-6 after:transition-all peer-checked:bg-green-600"
          />
        </label>
      </div>

      <div class="flex items-start justify-between pb-6 border-b border-gray-200">
        <div class="flex items-start gap-4">
          <div class="w-12 h-12 bg-green-50 rounded-lg flex items-center justify-center flex-shrink-0">
            <Smartphone class="w-6 h-6 text-green-600" />
          </div>
          <div>
            <h3 class="text-lg font-semibold text-gray-900">Push Notifications</h3>
            <p class="text-sm text-gray-600 mt-1">
              Get instant push notifications on your device
            </p>
          </div>
        </div>
        <label class="relative inline-flex items-center cursor-pointer">
          <input
            v-model="settings.pushNotifications"
            type="checkbox"
            class="sr-only peer"
          />
          <div
            class="w-14 h-7 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-green-300 rounded-full peer peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-0.5 after:start-[4px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-6 after:w-6 after:transition-all peer-checked:bg-green-600"
          />
        </label>
      </div>

      <div class="flex items-start gap-4">
        <div class="w-12 h-12 bg-orange-50 rounded-lg flex items-center justify-center flex-shrink-0">
          <Clock class="w-6 h-6 text-orange-600" />
        </div>
        <div class="flex-1">
          <h3 class="text-lg font-semibold text-gray-900 mb-1">Reminder Timing</h3>
          <p class="text-sm text-gray-600 mb-4">
            Choose when to receive reminders before the due date
          </p>
          <div class="space-y-3">
            <label
              class="flex items-center gap-3 p-4 border border-gray-200 rounded-lg hover:bg-gray-50 cursor-pointer transition-colors"
            >
              <input
                v-model="settings.reminderTiming"
                type="radio"
                name="timing"
                value="3days"
                class="w-4 h-4 text-green-600 focus:ring-green-500"
              />
              <div>
                <p class="font-medium text-gray-900">3 days before</p>
                <p class="text-sm text-gray-500">Early warning for better planning</p>
              </div>
            </label>

            <label
              class="flex items-center gap-3 p-4 border border-gray-200 rounded-lg hover:bg-gray-50 cursor-pointer transition-colors"
            >
              <input
                v-model="settings.reminderTiming"
                type="radio"
                name="timing"
                value="1day"
                class="w-4 h-4 text-green-600 focus:ring-green-500"
              />
              <div>
                <p class="font-medium text-gray-900">1 day before</p>
                <p class="text-sm text-gray-500">Standard reminder timing</p>
              </div>
            </label>

            <label
              class="flex items-center gap-3 p-4 border border-gray-200 rounded-lg hover:bg-gray-50 cursor-pointer transition-colors"
            >
              <input
                v-model="settings.reminderTiming"
                type="radio"
                name="timing"
                value="sameday"
                class="w-4 h-4 text-green-600 focus:ring-green-500"
              />
              <div>
                <p class="font-medium text-gray-900">Same day (morning)</p>
                <p class="text-sm text-gray-500">Last minute reminder</p>
              </div>
            </label>
          </div>
        </div>
      </div>
    </div>

    <div class="flex justify-end">
      <button
        type="button"
        :disabled="saved"
        :class="[
          'px-6 py-3 rounded-lg font-medium transition-all flex items-center gap-2',
          saved
            ? 'bg-green-100 text-green-700 cursor-not-allowed'
            : 'bg-green-600 text-white hover:bg-green-700',
        ]"
        @click="handleSave"
      >
        <template v-if="saved">
          <CheckCircle class="w-5 h-5" />
          Settings Saved!
        </template>
        <template v-else>
          <Bell class="w-5 h-5" />
          Save Settings
        </template>
      </button>
    </div>

    <div class="bg-white rounded-xl border border-gray-200 p-6">
      <h3 class="text-lg font-semibold text-gray-900 mb-4">Notification Preview</h3>
      <div class="space-y-3">
        <div class="p-4 bg-blue-50 border border-blue-200 rounded-lg">
          <div class="flex items-start gap-3">
            <Bell class="w-5 h-5 text-blue-600 mt-0.5" />
            <div>
              <p class="font-medium text-blue-900">Maintenance task due tomorrow</p>
              <p class="text-sm text-blue-700 mt-1">
                Your watering activity for BATCH-045 is scheduled for tomorrow. Don't forget to
                submit your report!
              </p>
            </div>
          </div>
        </div>

        <div class="p-4 bg-green-50 border border-green-200 rounded-lg">
          <div class="flex items-start gap-3">
            <CheckCircle class="w-5 h-5 text-green-600 mt-0.5" />
            <div>
              <p class="font-medium text-green-900">Your activity report has been verified</p>
              <p class="text-sm text-green-700 mt-1">
                Admin User has approved your pruning activity (MA-003). Great work!
              </p>
            </div>
          </div>
        </div>

        <div class="p-4 bg-red-50 border border-red-200 rounded-lg">
          <div class="flex items-start gap-3">
            <Bell class="w-5 h-5 text-red-600 mt-0.5" />
            <div>
              <p class="font-medium text-red-900">Your activity report was rejected</p>
              <p class="text-sm text-red-700 mt-1">
                Reason: Photo quality insufficient - please resubmit with clearer images
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
