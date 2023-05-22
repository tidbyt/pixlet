import axios from 'axios';

import { update, loading, error } from './schemaSlice';
import store from "../../store";


export default function refreshSchema() {
    store.dispatch(loading(true));

    axios.get(`${PIXLET_API_BASE}/api/v1/schema`)
        .then(res => {
            store.dispatch(update(res.data));
        })
        .catch(err => {
            // TODO: fix this.
            store.dispatch(error(err));
        })
        .then(() => {
            store.dispatch(loading(false));
        });
}