import { fileURLToPath, URL } from 'node:url';

import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import { VitePWA } from 'vite-plugin-pwa';

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [
        vue(),
        VitePWA({
            // The web app manifest is a static file (`public/manifest.webmanifest`)
            // wired into `index.html`; the plugin only generates the service worker.
            manifest: false,
            registerType: 'autoUpdate',
            injectRegister: false,
            workbox: {
                globPatterns: ['**/*.{js,css,html,ico,png,svg,webmanifest}'],
                navigateFallback: '/index.html',
                // One-time secret guarantee: `/api/*` must never be served from
                // cache. NetworkOnly means no cache write and no cache fallback.
                runtimeCaching: [
                    {
                        urlPattern: /\/api\//,
                        handler: 'NetworkOnly',
                    },
                ],
            },
        }),
    ],
    resolve: {
        alias: {
            '@': fileURLToPath(new URL('./src', import.meta.url)),
        },
    },
});
