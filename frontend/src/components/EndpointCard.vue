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

            <div v-if="ssl.Validation" class="mt-2">
                <p class="flex items-center gap-2">
                    <strong>Issued Via:</strong>
                    <span class="tooltip" :data-tip="ssl.Validation.Reason">
                        <span class="badge" :class="ssl.Validation.Method === 'DNS-01' ? 'badge-secondary' : 'badge-accent'">
                            {{ ssl.Validation.Method }}{{ ssl.Validation.Confirmed ? "" : " (likely)" }}
                        </span>
                    </span>
                </p>
            </div>

            <div v-for="(row, rowIdx) in caaRows" :key="row.label" class="mt-2">
                <p class="flex items-center gap-2">
                    <strong>{{ row.label }}:</strong>
                    <span class="flex flex-wrap gap-2">
                        <span v-for="(issuer, idx) in row.issuers" :key="idx" class="tooltip" :data-tip="caaTooltip(issuer)">
                            <span class="badge badge-neutral">{{ issuer.CA || "None" }}</span>
                        </span>
                    </span>
                    <span v-if="rowIdx === 0 && caaInheritedFrom" class="text-sm opacity-60">
                        (from {{ caaInheritedFrom }})
                    </span>
                </p>
            </div>

            <div v-if="ssl.CAAError" class="mt-2">
                <p class="flex items-center gap-2">
                    <strong>Allowed CAs:</strong>
                    <span class="tooltip" :data-tip="ssl.CAAError">
                        <span class="badge badge-warning">Lookup failed</span>
                    </span>
                </p>
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
import type { CAAIssuer, SSLDetails } from "@/types/certificate";

const props = defineProps<{
    endpoint: string;
    ssl?: SSLDetails | null;
}>();

const tlsProtocols = computed(() => props.ssl?.TLSProtocols?.filter((p) => p.Supported) ?? []);
const httpVersions = computed(() => props.ssl?.HTTPVersions?.filter((p) => p.Supported) ?? []);

const hostname = computed(() => {
    try {
        return new URL(props.endpoint).hostname.replace(/\.$/, "");
    } catch {
        return props.endpoint;
    }
});

const caaInheritedFrom = computed(() => {
    const domain = props.ssl?.CAA?.Domain;
    return domain && domain !== hostname.value ? domain : "";
});

const caaRows = computed(() => {
    const caa = props.ssl?.CAA;
    if (!caa) return [];
    const rows = [{ label: "Allowed CAs", issuers: caa.Issuers ?? [] }];
    if (caa.WildcardIssuers?.length) rows.push({ label: "Allowed CAs (Wildcard)", issuers: caa.WildcardIssuers });
    return rows.filter((r) => r.issuers.length);
});

function caaTooltip(issuer: CAAIssuer): string {
    if (!issuer.CA) return "No CA may issue";
    const parts = [`Methods: ${issuer.ValidationMethods?.join(", ") ?? "any"}`];
    if (issuer.AccountURI) parts.push(`Account: ${issuer.AccountURI}`);
    return parts.join(" · ");
}
</script>
