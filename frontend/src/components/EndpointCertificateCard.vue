<template>
    <div>
        <div class="flex gap-10 mb-2 items-center">
            <ExpiryCountdown :validTo="certificate.NotAfter" />

            <div>
                <p><strong>Valid From:</strong> {{ certificate.NotBefore }}</p>
                <p><strong>Valid To:</strong> {{ certificate.NotAfter }}</p>
            </div>
        </div>

        <p><strong>Subject:</strong> {{ certificate.Subject }}</p>
        <p><strong>Issuer:</strong> {{ certificate.Issuer }}</p>

        <div class="mt-2 flex gap-5">
            <p><strong>Signature Algorithm:</strong> {{ certificate.SignatureAlgorithm }}</p>
            <p><strong>Public Key Algorithm:</strong> {{ certificate.PublicKeyAlgorithm }}</p>
        </div>
    </div>
</template>

<script setup lang="ts">
import ExpiryCountdown from '@/components/ExpiryCountdown.vue';
import type { CertificateDetails } from '@/types/certificate';

defineProps<{
    certificate: CertificateDetails
}>();
</script>
