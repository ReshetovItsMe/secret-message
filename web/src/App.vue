<script setup lang="ts">
import { reactive } from 'vue';
import Footer from './components/Footer.vue';
import Header from './components/Header.vue';
import SecretInput from './components/SecretInput.vue';
import SecretUrl from './components/SecretUrl.vue';

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
    <div class="common-layout">
        <el-container class="common-layout">
            <el-header>
                <Header />
            </el-header>
            <el-main>
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
        flex-direction: column;
    }

    .el-container {
        padding: 15px;
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
}
</style>
