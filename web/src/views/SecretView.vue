<script async setup lang="ts">
import axios from 'axios';
import { ElMessage } from 'element-plus';
import { computed, reactive } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute();
const secretId = computed(() => route.params.id ?? null);
const store = {
    state: reactive<{ secret?: string }>({
        secret: undefined,
    }),

    setSecret(newValue: string) {
        this.state.secret = newValue;
    },
};

try {
    const { data } = await axios.get<{ message: string }>('/message', {
        params: { messageId: secretId.value },
    });
    store.setSecret(data.message);
} catch (error) {
    ElMessage.error(`No such secret id ${secretId.value}`);
}
</script>

<template>
    <Suspense>
        <div class="secret-container">
            <p>{{ store.state.secret }}</p>
        </div>
    </Suspense>
</template>

<style scoped>
.secret-container {
    background: var(--el-color-info-light-9);
    border-radius: 10px 10px 10px 10px;
    overflow-y: scroll;
    padding: 15px;
    text-align: center;
    font-size: var(--el-font-size-medium);
    max-height: 80%;
    max-width: 70%;
}
</style>
