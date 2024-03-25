<script setup lang="ts">
import Footer from './components/Footer.vue';
import Header from './components/Header.vue';
</script>

<template>
    <div class="common-layout">
        <el-container class="common-layout">
            <el-header>
                <Header />
            </el-header>
            <el-main>
                <Suspense>
                    <RouterView />
                </Suspense>
            </el-main>
            <el-footer>
                <Footer />
            </el-footer>
        </el-container>
    </div>
</template>

<style scoped>
.common-layout {
    height: 100vh;

    .el-main {
        display: flex;
        align-items: center;
        justify-content: center;
        flex-direction: column;
    }

    .el-container {
        padding: 15px;
    }
}
</style>
