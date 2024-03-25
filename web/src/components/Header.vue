<script setup lang="ts">
import SecretMessageLogo from './icons/secretMessage/SecretMessageLogo.vue';
import { useRouter } from 'vue-router';

const router = useRouter();

const onClickLogo = () => {
    router.push('/');
};
</script>

<template>
    <div class="header">
        <SecretMessageLogo class="logo" @click="onClickLogo" />
    </div>
</template>

<style scoped>
.header {
    display: flex;
    justify-content: flex-start;
    height: 100%;
}
.logo {
    cursor: pointer;
}
</style>
