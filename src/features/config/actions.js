import { update, clear } from './configSlice';
import store from '../../store';

export function setConfig(data) {
    store.dispatch(update(data));
}

export function resetConfig() {
    store.dispatch(clear());
}