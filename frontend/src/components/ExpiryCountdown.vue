<template>
  <div v-if="expired">
    <p class="font-mono text-2xl font-bold">EXPIRED</p>
  </div>
  <div v-else class="flex gap-5">
    <div>
      <span
        class="countdown font-mono text-2xl"
        :style="{ '--value': displayValue(timeLeft.days), '--digits': digits(timeLeft.days) } as any">
        {{ displayValue(timeLeft.days) }}
      </span>
      days
    </div>
    <div>
      <span
        class="countdown font-mono text-2xl"
        :style="{ '--value': displayValue(timeLeft.hours), '--digits': digits(timeLeft.hours) } as any">
        {{ displayValue(timeLeft.hours) }}
      </span>
      hours
    </div>
    <div>
      <span
        class="countdown font-mono text-2xl"
        :style="{ '--value': displayValue(timeLeft.minutes), '--digits': digits(timeLeft.minutes) } as any">
        {{ displayValue(timeLeft.minutes) }}
      </span>
      min
    </div>
  </div>
</template>

<script setup lang="ts">
import { expiryUntilCountdown } from '@/utils/certificate';
import { computed, onMounted, ref } from 'vue';

const props = defineProps<{
  validTo: string
}>();

const timeLeft = computed(() => {
  return expiryUntilCountdown(props.validTo, now.value);
});

const expired = computed(() => {
  return (
    timeLeft.value.days <= 0 &&
    timeLeft.value.hours <= 0 &&
    timeLeft.value.minutes <= 0
  );
});

const displayValue = (v: number) => {
  return Math.max(0, Math.min(999, Math.round(v))).toString();
};
const digits = (v: number) => {
  return v.toString().length;
};

const now = ref(new Date());
onMounted(() => {
  setInterval(() => {
    now.value = new Date();
  }, 30000);
});
</script>
