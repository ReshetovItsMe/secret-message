<script setup lang="ts">
import Footer from './components/Footer.vue';
import Header from './components/Header.vue';
import SecretInput from './components/SecretInput.vue';
</script>

<template>
    <div class="common-layout">
        <el-container class="common-layout">
            <el-header>
                <Header />
            </el-header>
            <el-main>
                <SecretInput class="secret-input" />
            </el-main>
            <el-footer>
                <Footer />
            </el-footer>
        </el-container>
    </div>
</template>

<style scoped>
.common-layout {
    min-height: 100vh;

    .el-main {
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .secret-input {
        width: 50%;
    }
}
</style>
