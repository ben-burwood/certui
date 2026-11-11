<template>
    <div class="min-h-screen w-full flex flex-col">
        <!-- HEADER -->
        <div class="navbar bg-base-100 shadow-sm">
            <div class="flex items-center gap-4 lg:ml-20">
                <img src="/logo.svg" alt="CertUI Logo" class="h-15" />
                <span class="text-2xl font-bold">CertUI</span>
            </div>
            <div class="flex-1"></div>
            <button
                class="btn btn-ghost btn-circle lg:mr-20"
                :disabled="loading"
                title="Force refresh (bypass cache)"
                aria-label="Force refresh"
                @click="loadEndpoints(true)"
            >
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="24"
                    height="24"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    :class="{ 'animate-spin': loading }"
                >
                    <path d="M3 12a9 9 0 0 1 9-9 9.75 9.75 0 0 1 6.74 2.74L21 8" />
                    <path d="M21 3v5h-5" />
                    <path d="M21 12a9 9 0 0 1-9 9 9.75 9.75 0 0 1-6.74-2.74L3 16" />
                    <path d="M8 16H3v5" />
                </svg>
            </button>
        </div>

        <!-- BODY -->
        <div class="p-5 flex-1 bg-base-200">
            <div class="max-w-6xl mx-auto" v-if="endpointsData.length">
                <details
                    class="my-4 collapse collapse-arrow bg-base-100"
                    name="certificate-accordion"
                    v-for="(endpointData, index) in endpointsData"
                    :key="index"
                >
                    <summary class="collapse-title flex flex-wrap gap-2 sm:gap-4">
                        <span class="font-semibold">{{ endpointData.endpoint }}</span>
                        <span v-if="endpointData.details?.Domain && endpointData.details.Domain.Address">
                            ({{ endpointData.details.Domain.Address }})
                        </span>
                        <div class="font-semibold flex gap-2">
                            <div v-if="endpointData.details?.Domain && endpointData.details.Domain.Resolves === false" class="badge badge-error">
                                Domain Can't Resolve
                            </div>
                            <ExpiryStatus
                                v-if="endpointData.details?.SSL?.PeerCertificates?.[0]?.NotAfter"
                                :notAfter="endpointData.details.SSL?.PeerCertificates?.[0]?.NotAfter"
                                :daysRemainingLimit="14"
                            />
                        </div>
                    </summary>
                    <div class="collapse-content">
                        <div v-if="endpointData.details?.Whois">
                            <div class="divider">WHOIS</div>
                            <WhoIs :whois="endpointData.details.Whois" />
                        </div>
                        <div v-if="endpointData.details?.SSL">
                            <div class="divider">SSL</div>
                            <EndpointCard class="p-5" :endpoint="endpointData.endpoint" :ssl="endpointData.details?.SSL" />
                        </div>
                    </div>
                </details>
            </div>
            <div v-else-if="loading" class="text-center mt-10">
                <span class="loading loading-dots loading-xl"></span>
            </div>
            <div v-else class="text-center italic">No endpoints.</div>
        </div>

        <!-- FOOTER -->
        <footer class="border-t border-base-300 bg-base-100">
            <div class="px-4 py-6 text-sm text-base-content/70 text-center">
                Powered by
                <a href="https://github.com/ben-burwood/certui" target="_blank" class="font-medium text-emerald-800 hover:text-emerald-600">CertUI</a>
            </div>
        </footer>

        <ThemeSwitcher />
        <Social />
    </div>
</template>

<script setup lang="ts">
import ThemeSwitcher from "./components/ThemeSwitcher.vue";
import Social from "./components/Social.vue";
import EndpointCard from "@/components/EndpointCard.vue";
import { SERVER_URL } from "@/main";
import type { EndpointDetails } from "@/types/endpoint";
import { ref, onMounted, onUnmounted } from "vue";
import ExpiryStatus from "./components/ExpiryStatus.vue";
import WhoIs from "./components/WhoIs.vue";

const endpointsData = ref<{ endpoint: string; details: EndpointDetails | null }[]>([]);

const loading = ref(true);

// Basic HTTP Inplementation
// onMounted(async () => {
//     loading.value = true;
//     try {
//         const res = await fetch(`${SERVER_URL}/endpoints`);
//         const data: Record<string, EndpointDetails | null> = await res.json();
//         endpointsData.value = Object.entries(data).map(
//             ([endpoint, details]) => ({
//                 endpoint,
//                 details,
//             }),
//         );
//     } catch (e) {
//         console.error("Failed to fetch endpoints:", e);
//         endpointsData.value = [];
//     } finally {
//         loading.value = false;
//     }
// });

// Server Sent Events Implementation
let eventSource: EventSource | null = null;

const loadEndpoints = (force = false) => {
    if (eventSource) eventSource.close();
    endpointsData.value = [];
    loading.value = true;

    eventSource = new EventSource(`${SERVER_URL}/endpoints-sse${force ? "?refresh=true" : ""}`);
    eventSource.onmessage = (event) => {
        const data = JSON.parse(event.data);
        endpointsData.value.push(data);
    };
    eventSource.addEventListener("done", () => {
        console.log("EventSource done");
        if (eventSource) eventSource.close();
        loading.value = false;
    });
    eventSource.onerror = (e) => {
        console.error("EventSource error:", e);
        loading.value = false;
        if (eventSource) eventSource.close();
    };
    eventSource.onopen = () => {
        console.log("EventSource opened");
    };
};

onMounted(() => loadEndpoints());

onUnmounted(() => {
    if (eventSource) eventSource.close();
});
</script>
