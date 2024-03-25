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
            <el-text size="large" v-model="textToCopy">{{ fullUrl }}</el-text>

            <el-button
                text
                size="large"
                type="''"
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
}
.url-block-link {
    display: flex;
    flex-direction: row;
}
.procceed-button {
    margin-top: 5pt;
    align-self: end;
}
</style>
