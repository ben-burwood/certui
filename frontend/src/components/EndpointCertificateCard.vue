<template>
    <div>
        <div class="flex flex-col sm:flex-row gap-4 mb-2 items-start sm:items-center">
            <ExpiryCountdown :validTo="certificate.NotAfter" />

            <div>
                <p><strong>Valid From:</strong> {{ certificate.NotBefore }}</p>
                <p><strong>Valid To:</strong> {{ certificate.NotAfter }}</p>
            </div>
        </div>

        <p class="flex items-center gap-2">
            <span><strong>Subject:</strong> {{ certificate.Subject }}</span>
            <span v-if="certificate.IsWildcard" class="badge badge-secondary">Wildcard</span>
        </p>
        <p><strong>Issuer:</strong> {{ certificate.Issuer }}</p>

        <div v-if="certificate.DNSNames?.length" class="mt-2">
            <p class="flex items-center gap-2">
                <strong>Subject Alternative Names:</strong>
                <span class="flex flex-wrap gap-2">
                    <span v-for="name in certificate.DNSNames" :key="name" class="badge badge-outline">
                        {{ name }}
                    </span>
                </span>
            </p>
        </div>

        <p class="mt-2 break-all"><strong>SHA-256 Fingerprint:</strong> <code>{{ certificate.Fingerprint }}</code></p>

        <div class="mt-2 flex flex-wrap gap-4">
            <p><strong>Signature Algorithm:</strong> {{ certificate.SignatureAlgorithm }}</p>
            <p><strong>Public Key Algorithm:</strong> {{ certificate.PublicKeyAlgorithm }}</p>
        </div>
    </div>
</template>

<script setup lang="ts">
import ExpiryCountdown from "@/components/ExpiryCountdown.vue";
import type { CertificateDetails } from "@/types/certificate";

defineProps<{
    certificate: CertificateDetails;
}>();
</script>
