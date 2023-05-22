import axios from 'axios';
import { update, loading } from './previewSlice';
import { set as setError, clear as clearErrors } from '../errors/errorSlice';
import store from '../../store';


export default function fetchPreview(formData) {
    store.dispatch(loading(true));
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