<script setup lang="ts">
import logo from './logo.png';
</script>

<template>
    <div class="logo-wrapper" key="scale-down">
        <el-image style="width: 40%" :src="logo" fit="scale-down" />
    </div>
</template>

<style scoped>
.logo-wrapper {
    display: flex;
    padding-top: 20px;
    padding-left: 30px;
}
</style>
