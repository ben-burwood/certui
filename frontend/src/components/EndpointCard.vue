<template>
  <div>
    <div v-if="ssl">
      <div class="flex gap-4">
        <p><strong>TLS Version:</strong> {{ ssl.Version }}</p>
        <p><strong>Handshake Complete:</strong> {{ ssl.HandshakeComplete ? 'Yes' : 'No' }}</p>
        <p><strong>Did Resume:</strong> {{ ssl.DidResume ? 'Yes' : 'No' }}</p>
        <p><strong>Cipher Suite:</strong> {{ ssl.CipherSuite }}</p>
      </div>

      <ul v-if="ssl.PeerCertificates && ssl.PeerCertificates.length" class="list mt-2">
        <li class="p-4 pb-2 text-md opacity-80 tracking-wide">Peer Certificates</li>
        <li v-for="(cert, idx) in ssl.PeerCertificates" :key="idx" class="list-row">
          <EndpointCertificateCard :certificate="cert" />
        </li>
      </ul>
    </div>
    <div v-else>
      <em>No SSL data available.</em>
    </div>
  </div>
</template>

<script setup lang="ts">
import EndpointCertificateCard from '@/components/EndpointCertificateCard.vue';
import type { SSLDetails } from '@/types/certificate';

defineProps<{
  endpoint: string;
  ssl?: SSLDetails | null
}>();
</script>
