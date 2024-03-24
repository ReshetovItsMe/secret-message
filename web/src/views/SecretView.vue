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
        <p>Id of secret:</p>
        <span>{{ secretId }}</span>
    </div>
</template>

<style scoped>
.secret-container {
    align-items: center;
    display: flex;
    flex-direction: column;
    flex: 1;
    width: 100%;
}
</style>
