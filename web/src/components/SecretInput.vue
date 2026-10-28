<script lang="ts" setup>
import { ref } from 'vue';
import axios from 'axios';
import { ElMessage } from 'element-plus';

const emit = defineEmits<{
    (e: 'returnSecretUrl', url: string): void;
}>();

const message = ref('');
const placeholder: string =
    '- Write your secret and hit "Share" button. 🔘 \n- The link will be available 24 or until opened. ⏱️ \n- All data is encrypted. 🤫 \n- Tell your friends about us! 🕺';

const postSecret = async () => {
    try {
        const { data } = await axios.post<{ messageId: string }>('/message', {
            message: message.value,
        });

        emit('returnSecretUrl', data.messageId);
    } catch (error) {
        ElMessage.error(`Cannot post secret: ${error}`);
    }
};
</script>

<template>
    <div class="input-block">
        <textarea
            v-model="message"
            class="form-control"
            :placeholder="placeholder"
        ></textarea>
        <el-button
            class="procceed-button"
            type="primary"
            @click="postSecret()"
            :disabled="message.length < 1"
            >Share</el-button
        >
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
    min-height: 40vh;
    padding: 12px;
    font-size: 16px;
    flex: 1;
}

.procceed-button {
    margin-top: 5pt;
    align-self: end;
}

@media (max-width: 767px) {
    .procceed-button {
        align-self: stretch;
    }
}
</style>
