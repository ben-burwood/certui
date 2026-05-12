<template>
    <div>
        <div v-if="ssl">
            <div class="flex flex-wrap gap-4">
                <p>
                    <strong>Handshake Complete:</strong>
                    {{ ssl.HandshakeComplete ? "Yes" : "No" }}
                </p>
                <p>
                    <strong>Did Resume:</strong>
                    {{ ssl.DidResume ? "Yes" : "No" }}
                </p>
                <p><strong>Cipher Suite:</strong> {{ ssl.CipherSuite }}</p>
            </div>

            <div v-if="tlsProtocols.length" class="mt-2">
                <p class="flex items-center gap-2">
                    <strong>TLS Protocols:</strong>
                    <span class="flex flex-wrap gap-2">
                        <span v-for="p in tlsProtocols" :key="p.Protocol" class="badge badge-info">
                            {{ p.Protocol }}
                        </span>
                    </span>
                </p>
            </div>

            <div v-if="httpVersions.length" class="mt-2">
                <p class="flex items-center gap-2">
                    <strong>HTTP Versions:</strong>
                    <span class="flex flex-wrap gap-2">
                        <span v-for="p in httpVersions" :key="p.Version" class="badge badge-info">
                            {{ p.Version }}
                        </span>
                    </span>
                </p>
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
import { computed } from "vue";
import EndpointCertificateCard from "@/components/EndpointCertificateCard.vue";
import type { SSLDetails } from "@/types/certificate";

const props = defineProps<{
    endpoint: string;
    ssl?: SSLDetails | null;
}>();

const tlsProtocols = computed(() => props.ssl?.TLSProtocols?.filter((p) => p.Supported) ?? []);
const httpVersions = computed(() => props.ssl?.HTTPVersions?.filter((p) => p.Supported) ?? []);
</script>
