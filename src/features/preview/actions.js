import axios from 'axios';
import { update, loading } from './previewSlice';
import { set as setError, clear as clearErrors } from '../errors/errorSlice';
import store from '../../store';

let timeout = null;

export default function fetchPreview(formData) {
    if (PIXLET_WASM) {
        store.dispatch(loading(true));
        clearTimeout(timeout);
        timeout = setTimeout(function () {
            axios.post(`${PIXLET_API_BASE}/api/v1/preview`, formData)
                .then(res => {
                    document.title = res.data.title;
                    store.dispatch(update(res.data));
                    if ('error' in res.data) {
                        store.dispatch(setError({ id: res.data.error, message: res.data.error }));
                    } else {
                        store.dispatch(clearErrors());
                    }
                })
                .catch(err => {
                    // TODO: fix this.
                    store.dispatch(setError({ id: err, message: err }));
                })
                .then(() => {
                    store.dispatch(loading(false));
                })
        }, 300);

    } else {
        axios.post(`${PIXLET_API_BASE}/api/v1/preview`, formData)
            .then(res => {
                document.title = res.data.title;
                store.dispatch(update(res.data));
                if ('error' in res.data) {
                    store.dispatch(setError({ id: res.data.error, message: res.data.error }));
                } else {
                    store.dispatch(clearErrors());
                }
            })
            .catch(err => {
                // TODO: fix this.
                store.dispatch(setError({ id: err, message: err }));
            })
            .then(() => {
                store.dispatch(loading(false));
            })

    }
}