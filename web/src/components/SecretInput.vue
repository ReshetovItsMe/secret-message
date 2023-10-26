<script lang="ts" setup>
import { ref } from 'vue';
const textarea = ref('');
</script>

<template>
    <div class="input-block">
        <textarea
            v-model="textarea"
            class="form-control"
            placeholder="Hello, input here"
        ></textarea>
        <el-button class="procceed-button" type="primary">Primary</el-button>
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
