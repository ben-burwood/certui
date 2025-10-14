<template>
    <div v-if="ssl.IsExpired" class="badge badge-error">Expired</div>
    <div v-else-if="daysRemaining >= daysRemainingLimit" class="badge badge-success">{{ daysRemaining }} Days Remaining</div>
    <div v-else-if="daysRemaining !== null && daysRemaining < daysRemainingLimit" class="badge badge-warning">{{ daysRemaining }} Days Remaining</div>
    <div v-else class="badge badge-info">No Certificates</div>
</template>

<script setup lang="ts">
import type { SSLDetails } from '@/types/certificate';
import { daysUntil } from '@/utils/certificate';
import { computed } from 'vue';

const props = defineProps<{
    ssl: SSLDetails;
    daysRemainingLimit?: number
}>();

const daysRemainingLimit = props.daysRemainingLimit ?? 30;

const daysRemaining = computed(() => {
    const notAfter = props.ssl?.PeerCertificates?.[0]?.NotAfter ?? '';
    return daysUntil(notAfter, new Date());
});
</script>
