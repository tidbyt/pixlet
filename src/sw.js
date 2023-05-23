import cmsStarlark from './apps/cms.star';

importScripts('https://cdn.jsdelivr.net/gh/golang/go@go1.18.4/misc/wasm/wasm_exec.js');
importScripts('https://cdn.jsdelivr.net/gh/nlepage/go-wasm-http-server@v1.1.0/sw.js');

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