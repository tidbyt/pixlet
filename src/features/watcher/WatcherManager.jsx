import { useEffect } from 'react';

import Watcher from './watcher';


export default function WatcherManager() {
    useEffect(() => {
        new Watcher();
    }, []);

    return null;
}