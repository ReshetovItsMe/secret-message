import { registerSW } from 'virtual:pwa-register';

export function registerServiceWorker(): void {
    // Service workers require a secure context (HTTPS or localhost). On plain
    // HTTP we skip registration entirely and the app degrades to a normal SPA.
    if (!window.isSecureContext || !('serviceWorker' in navigator)) {
        return;
    }

    registerSW({ immediate: true });
}
