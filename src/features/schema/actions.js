import axios from 'axios';
import axiosRetry from 'axios-retry';

import { update, loading, error } from './schemaSlice';
import store from "../../store";


export default function refreshSchema() {
    store.dispatch(loading(true));

    const client = axios.create();
    axiosRetry(client, {
        retries: 5,
        retryDelay: () => 1000,
        retryCondition: (err) => {
            return err.response.status === 404;
        },
    });

    client.get(`${PIXLET_API_BASE}/api/v1/schema`)
        .then(res => {
            store.dispatch(update(res.data));
        })
        .catch(err => {
            if (err.response.status == 404) {
                store.dispatch(error({ id: err, message: "error with pixlet, please refresh page" }));
                return;
            }
            store.dispatch(error(err));
        })
        .then(() => {
            store.dispatch(loading(false));
        });
}