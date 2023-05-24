import cmsStarlark from './apps/cms.star';

importScripts(new URL('./go/wasm_exec.js', import.meta.url));
importScripts(new URL('./go/sw.js', import.meta.url));

registerWasmHTTPListener(
    new URL('./pixlet.wasm', import.meta.url),
    {
        'base': 'pixlet',
        'args': ['serve', cmsStarlark],
    }
);

// Skip installed stage and jump to activating stage
addEventListener('install', (event) => {
    event.waitUntil(skipWaiting())
});

// Start controlling clients as soon as the SW is activated
addEventListener('activate', event => {
    event.waitUntil(clients.claim())
});