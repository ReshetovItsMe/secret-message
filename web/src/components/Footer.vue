<script setup lang="ts">
import code from '../assets/code.png';
import github from '../assets/github.png';
import linkedin from '../assets/linkedin.png';
</script>

<template>
    <div class="footer">
        <a
            class="logo-wrapper"
            key="scale-down"
            href="https://github.com/ReshetovItsMe/secret-message"
            target="_blank"
        >
            <el-image class="image" :src="code" fit="scale-down" />
        </a>
        <a
            class="logo-wrapper"
            key="scale-down"
            href="https://github.com/ReshetovItsMe"
            target="_blank"
        >
            <el-image class="image" :src="github" fit="scale-down" />
        </a>
        <a
            class="logo-wrapper"
            key="scale-down"
            href="https://www.linkedin.com/in/vladislav-reshetov/"
            target="_blank"
        >
            <el-image class="image" :src="linkedin" fit="scale-down" />
        </a>
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
