<script setup lang="ts">
import { computed } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute();

const secretId = computed(() => {
    return route.params.id ?? null;
});
</script>

<template>
    <div class="secret-container">
        <p>SUPER LONG TEXT SUPER LONG TEXT SUPER LONG TEXT SUPER LONG TEXT</p>
    </div>
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
