<script setup lang="ts">
import code from '../assets/code.png';
import github from '../assets/github.png';
import linkedin from '../assets/linkedin.png';
</script>

<template>
    <div class="footer">
        <div class="logo-wrapper" key="scale-down">
            <el-image style="width: 50%" :src="code" fit="scale-down" />
        </div>
        <div class="logo-wrapper" key="scale-down">
            <el-image style="width: 50%" :src="github" fit="scale-down" />
        </div>
        <div class="logo-wrapper" key="scale-down">
            <el-image style="width: 50%" :src="linkedin" fit="scale-down" />
        </div>
    </div>
</template>

<style scoped>
.logo-wrapper {
    display: flex;
    padding: 0 10px;
}
.footer {
    display: flex;
    justify-content: flex-end;
}
</style>
