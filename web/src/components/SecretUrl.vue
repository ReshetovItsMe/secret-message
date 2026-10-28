<script lang="ts" setup>
import { computed, ref } from 'vue';
import { CopyDocument } from '@element-plus/icons-vue';
import { useClipboard } from '@vueuse/core';

const props = defineProps<{
    secretId: string;
}>();

const emit = defineEmits<{
    (e: 'makeNewSecret'): void;
}>();

const fullUrl = computed(() => {
    const domain = window.location.origin;
    return `${domain}/secret/${props.secretId}`;
});
const textToCopy = ref(fullUrl.value);

const makeNewSecret = () => {
    try {
        emit('makeNewSecret');
    } catch (error) {
        console.error(error);
    }
};

const { copy } = useClipboard({ source: textToCopy });
</script>

<template>
    <div class="url-block">
        <div class="url-block-link">
            <el-text size="large">{{ fullUrl }}</el-text>

            <el-button
                text
                size="large"
                :icon="CopyDocument"
                @click="copy()"
            />
        </div>
        <el-button
            class="procceed-button"
            type="primary"
            @click="makeNewSecret()"
            >New Secret</el-button
        >
    </div>
</template>

<style scoped>
.url-block {
    display: flex;
    flex-direction: column;
    max-width: 100%;
}
.url-block-link {
    display: flex;
    flex-direction: row;
    align-items: center;
    overflow-wrap: anywhere;
    min-width: 0;
}
.procceed-button {
    margin-top: 5pt;
    align-self: end;
}

@media (max-width: 767px) {
    .url-block-link {
        flex-direction: column;
        align-items: stretch;
        gap: 8px;
    }

    .procceed-button {
        align-self: stretch;
    }
}
</style>
