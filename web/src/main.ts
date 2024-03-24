import './index.css';
import { createApp } from 'vue';
import ElementPlus from 'element-plus';
import 'element-plus/dist/index.css';
import App from './App.vue';
import { createMemoryHistory, createRouter } from 'vue-router';
import MainView from './views/MainView.vue';

const routes = [{ path: '/', component: MainView }];

const router = createRouter({
    history: createMemoryHistory(),
    routes,
});

const app = createApp(App);

app.use(router);
app.use(ElementPlus);
app.mount('#app');
