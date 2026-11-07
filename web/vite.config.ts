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
    server: {
        proxy: {
            // Mirror nginx's rewrite: strip the `/api` prefix and forward to
            // the local gateway, so `pnpm dev` can reach the API without the
            // nginx container. The gateway serves routes like `/message`.
            '/api': {
                target: 'http://localhost:3000',
                changeOrigin: true,
                rewrite: (path) => path.replace(/^\/api/, ''),
            },
        },
    },
});
