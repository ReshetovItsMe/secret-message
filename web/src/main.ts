import './index.css';
import { createApp } from 'vue';
import ElementPlus from 'element-plus';
import 'element-plus/dist/index.css';
import App from './App.vue';
import { createWebHistory, createRouter } from 'vue-router';
import MainView from './views/MainView.vue';
import SecretView from './views/SecretView.vue';
import axios from 'axios';

axios.defaults.baseURL = 'http://gateway:3000';

const routes = [
    { path: '/secret/:id', component: SecretView },
    { path: '/', component: MainView },
    { path: '/:pathMatch(.*)*', redirect: '/' },
];

const router = createRouter({
    history: createWebHistory(),
    routes,
});

const app = createApp(App);

app.use(router);
app.use(ElementPlus);
app.mount('#app');
