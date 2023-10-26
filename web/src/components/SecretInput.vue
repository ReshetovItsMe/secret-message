<script lang="ts" setup>
import { ref } from 'vue';

const textarea = ref('');
const placeholder: string =
    '- Write your secret and hit "Share" button. 🔘 \n- The link will be available 24 or until opened. ⏱️ \n- All data is encrypted. 🤫 \n- Tell your friends about us! 🕺';
</script>

<template>
    <div class="input-block">
        <textarea
            v-model="textarea"
            class="form-control"
            :placeholder="placeholder"
        ></textarea>
        <el-button class="procceed-button" type="primary">Share</el-button>
    </div>
</template>

<style scoped>
.input-block {
    display: flex;
    flex-direction: column;
}

.form-control {
    resize: none;
    height: 100%;
    width: 100%;
    flex: 1;
}

.procceed-button {
    margin-top: 5pt;
    align-self: end;
}
</style>
