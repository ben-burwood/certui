<template>
    <div class="min-h-screen w-full flex flex-col">
        <div class="navbar bg-base-100 shadow-sm">
            <div class="flex items-center gap-4 lg:ml-20">
                <img src="/logo.svg" alt="CertUI Logo" class="h-15" />
                <span class="text-2xl font-bold">CertUI</span>
            </div>
        </div>

        <div class="p-5 flex-1 bg-gray-50">
            <div class="max-w-6xl mx-auto">
                <div
                    class="my-4"
                    v-if="endpointsData.length"
                    v-for="(endpointData, index) in endpointsData"
                    :key="index"
                >
                    <div
                        v-if="endpointData.ssl"
                        class="collapse collapse-arrow bg-base-100 border border-base-300"
                    >
                        <input type="radio" name="certificate-accordion" />
                        <div class="collapse-title font-semibold">
                            {{ endpointData.endpoint }}
                            <ExpiryStatus
                                class="ml-4"
                                :ssl="endpointData.ssl"
                                :daysRemainingLimit="14"
                            />
                        </div>
                        <div class="collapse-content">
                            <EndpointCard
                                class="p-5"
                                :endpoint="endpointData.endpoint"
                                :ssl="endpointData.ssl"
                            />
                        </div>
                    </div>
                </div>
                <div v-else class="text-center text-gray-500 italic">
                    No endpoints.
                </div>
            </div>
        </div>

        <footer class="p-4 bg-gray-200 flex items-center justify-between">
            <aside class="text-left">
                Powered by
                <a
                    href="https://github.com/ben-burwood/certui"
                    target="_blank"
                    class="font-medium text-green-800 hover:text-green-600"
                    >CertUI</a
                >
            </aside>
            <ThemeSwitcher />
            <nav class="text-right">
                <a
                    href="https://github.com/ben-burwood/certui"
                    target="_blank"
                    title="CertUI on GitHub"
                >
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        width="32"
                        height="32"
                        viewBox="0 0 16 16"
                        class="hover:scale-110"
                    >
                        <path
                            fill="gray"
                            d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z"
                        />
                    </svg>
                </a>
            </nav>
        </footer>
    </div>
</template>

<script setup lang="ts">
import ThemeSwitcher from "./components/ThemeSwitcher.vue";
import EndpointCard from "@/components/EndpointCard.vue";
import { SERVER_URL } from "@/main";
import type { SSLDetails } from "@/types/certificate";
import { onMounted, ref } from "vue";
import ExpiryStatus from "./components/ExpiryStatus.vue";

const endpointsData = ref<{ endpoint: string; ssl: SSLDetails | null }[]>([]);

onMounted(async () => {
    try {
        const res = await fetch(`${SERVER_URL}/endpoints`);
        const data: Record<string, SSLDetails | null> = await res.json();
        endpointsData.value = Object.entries(data).map(([endpoint, ssl]) => ({
            endpoint,
            ssl,
        }));
    } catch (e) {
        console.error("Failed to fetch endpoints:", e);
        endpointsData.value = [];
    }
});
</script>
