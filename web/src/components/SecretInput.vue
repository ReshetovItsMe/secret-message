<script lang="ts" setup>
import { ref } from 'vue';
const textarea = ref('');
</script>

<template>
    <div class="input-block">
        <el-input
            v-model="textarea"
            :rows="15"
            type="textarea"
            resize="none"
            placeholder="Please input"
        />
        <el-button class="procceed-button" type="primary">Primary</el-button>
    </div>
</template>

<style scoped>
.input-block {
    display: flex;
    flex-direction: column;
}

.procceed-button {
    margin-top: 5pt;
    align-self: end;
}
</style>
