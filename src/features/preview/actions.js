import axios from 'axios';
import { update, loading } from './previewSlice';
import { set as setError, clear as clearErrors } from '../errors/errorSlice';
import store from '../../store';
import axiosRetry from 'axios-retry';

let timeout = null;

export default function fetchPreview(formData) {
    const client = axios.create();
    axiosRetry(client, {
        retries: 5,
        retryDelay: () => 1000,
        retryCondition: (err) => {
            return err.response.status === 404;
        },
    });

    client.post(`${PIXLET_API_BASE}/api/v1/preview`, formData)
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
            if (err.response.status == 404) {
                store.dispatch(setError({ id: err, message: "error with pixlet, please refresh page" }));
                return;
            }
            store.dispatch(setError({ id: err, message: err }));
        })
        .then(() => {
            store.dispatch(loading(false));
        })
}