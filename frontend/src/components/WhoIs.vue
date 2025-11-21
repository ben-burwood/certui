<template>
    <div class="flex flex-col gap-4">
        <div><strong>Registrar:</strong> {{ whois.Registrar }}</div>
        <div class="flex gap-2">
            <strong>Name Servers:</strong>
            <ul class="flex gap-2">
                <li v-for="ns in whois.NameServers" :key="ns">{{ ns }}</li>
            </ul>
        </div>
        <div class="flex gap-10 items-center">
            <ExpiryCountdown :validTo="whois.ExpirationDate" />
            <div>
                <strong>Expiration Date:</strong> {{ whois.ExpirationDate }}
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import ExpiryCountdown from "@/components/ExpiryCountdown.vue";
import type { WhoisDetails } from "@/types/whois";

defineProps<{
    whois: WhoisDetails;
}>();
</script>
