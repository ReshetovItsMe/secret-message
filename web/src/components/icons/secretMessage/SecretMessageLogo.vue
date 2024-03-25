<script setup lang="ts">
import logo from './logo.png';
</script>

<template>
    <div class="logo-wrapper" key="scale-down">
        <el-image style="height: 100%" :src="logo" fit="scale-down" />
    </div>
</template>

<style scoped>
.logo-wrapper {
    display: flex;
    height: 90%;
}
</style>
