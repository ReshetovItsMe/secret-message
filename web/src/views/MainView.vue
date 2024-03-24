<script setup lang="ts">
import { reactive } from 'vue';
import SecretInput from '../components/SecretInput.vue';
import SecretUrl from '../components/SecretUrl.vue';

const store = {
    state: reactive<{ url?: string }>({
        url: undefined,
    }),

    setUrlAction(newValue: string) {
        this.state.url = newValue;
    },

    clearUrlAction() {
        this.state.url = undefined;
    },
};

const setUrlValue = (url: string) => {
    store.setUrlAction(url);
};
const makeNewSecret = () => {
    store.clearUrlAction();
};
</script>

<template>
    <div class="secret-container">
        <SecretUrl
            v-if="!!store.state.url"
            :secretId="store.state.url"
            @make-new-secret="makeNewSecret"
        />
        <SecretInput
            v-else
            class="secret-input"
            @return-secret-url="setUrlValue"
        />
    </div>
</template>

<style scoped>
.secret-container {
    align-items: center;
    display: flex;
    flex-direction: column;
    flex: 1;
    width: 100%;
}

.secret-input {
    flex: 1;
}

@media (min-width: 1100px) {
    .secret-input {
        width: 70%;
    }
}

@media (max-width: 1099px) {
    .secret-input {
        width: 90%;
    }
}
</style>
