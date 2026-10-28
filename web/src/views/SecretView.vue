<script setup lang="ts">
import axios from 'axios';
import { ElMessage } from 'element-plus';
import { computed, onMounted, ref } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute();
const secretId = computed(() => route.params.id ?? null);

const secret = ref<string | null>(null);
const loading = ref(true);

onMounted(async () => {
    try {
        const { data } = await axios.get<{ message: string }>('/message', {
            params: { messageId: secretId.value },
        });
        secret.value = data.message;
    } catch (error) {
        if (axios.isAxiosError(error) && error.response?.status === 404) {
            ElMessage.error(
                'This secret has already been opened or does not exist.',
            );
        } else {
            ElMessage.error(
                'Something went wrong while loading the secret. Please try again.',
            );
        }
    } finally {
        loading.value = false;
    }
});
</script>

<template>
    <div class="secret-container">
        <p v-if="loading">Loading…</p>
        <p v-else>{{ secret }}</p>
    </div>
</template>

<style scoped>
.secret-container {
    background: var(--el-color-info-light-9);
    border-radius: 10px 10px 10px 10px;
    overflow-y: scroll;
    overflow-wrap: anywhere;
    padding: 15px;
    text-align: center;
    font-size: var(--el-font-size-medium);
    max-height: 80%;
    max-width: 70%;
}

@media (max-width: 767px) {
    .secret-container {
        max-height: 85%;
        max-width: 100%;
        padding: 12px;
    }
}
</style>
