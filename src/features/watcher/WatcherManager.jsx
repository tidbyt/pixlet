import { useEffect } from 'react';

import Watcher from './watcher';


export default function WatcherManager() {
    if (PIXLET_WASM) {
        // watcher isn't supported in the browser
        return null;
    }

    useEffect(() => {
        new Watcher();
    }, []);

    return null;
}