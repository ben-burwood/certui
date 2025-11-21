<template>
    <div v-if="daysRemaining >= daysRemainingLimit" class="badge badge-success">
        {{ daysRemaining }} Days Remaining
    </div>
    <div
        v-else-if="
            daysRemaining !== null &&
            daysRemaining < daysRemainingLimit &&
            daysRemaining > 0
        "
        class="badge badge-warning"
    >
        {{ daysRemaining }} Days Remaining
    </div>
    <div
        v-else-if="daysRemaining !== null && daysRemaining <= 0"
        class="badge badge-error"
    >
        Expired
    </div>
</template>

<script setup lang="ts">
import { daysUntil } from "@/utils/expiry";
import { computed } from "vue";

const props = defineProps<{
    notAfter: string;
    daysRemainingLimit?: number;
}>();

const daysRemainingLimit = props.daysRemainingLimit ?? 30;

const daysRemaining = computed(() => {
    return daysUntil(props.notAfter, new Date());
});
</script>
