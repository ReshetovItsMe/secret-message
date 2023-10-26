<script setup lang="ts">
import code from '../assets/code.png';
import github from '../assets/github.png';
import linkedin from '../assets/linkedin.png';
</script>

<template>
    <div class="footer">
        <div class="logo-wrapper" key="scale-down">
            <el-image class="image" :src="code" fit="scale-down" />
        </div>
        <div class="logo-wrapper" key="scale-down">
            <el-image class="image" :src="github" fit="scale-down" />
        </div>
        <div class="logo-wrapper" key="scale-down">
            <el-image class="image" :src="linkedin" fit="scale-down" />
        </div>
    </div>
</template>

<style scoped>
.logo-wrapper {
    display: flex;
    padding: 0 5px;
}
.footer {
    display: flex;
    justify-content: flex-end;
    height: 100%;
}

.image {
    height: 50%;
    margin: auto;
}
</style>
