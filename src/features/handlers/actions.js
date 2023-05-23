
import axios from 'axios';
import { update, loading } from './handlerSlice';
import { updateGenerated } from '../schema/schemaSlice';
import { set as setError } from '../errors/errorSlice';
import store from "../../store";

function parseErrorMessage(handler, error) {
    let msg = `could not call handler ${handler}: ${error.message}`;
    if ("response" in error && "data" in error.response) {
        msg = error.response.data;
    }
    return msg;
}

export function callHandler(id, handler, param) {
    let data = {
        id: id,
        param: param
    }

    store.dispatch(loading(true));
    axios.post(`${PIXLET_API_BASE}/api/v1/handlers/` + handler, JSON.stringify(data))
        .then(res => {
            store.dispatch(update({ id: id, value: res.data }));
        })
        .catch(err => {
            // TODO: make sure this clears.
            const msg = parseErrorMessage(handler, err);
            store.dispatch(setError({ id: msg, message: msg }));
            console.log(err);
        })
        .then(() => {
            store.dispatch(loading(false));
        })
}

export function callGeneratedHandler(id, handler, param) {
    let data = {
        id: id,
        param: param
    }

    store.dispatch(loading(true));
    axios.post(`${PIXLET_API_BASE}/api/v1/handlers/` + handler, JSON.stringify(data))
        .then(res => {
            store.dispatch(updateGenerated(res.data));
        })
        .catch(err => {
            // TODO: make sure this clears.
            const msg = parseErrorMessage(handler, err);
            store.dispatch(setError({ id: msg, message: msg }));
            console.log(err);
        })
        .then(() => {
            store.dispatch(loading(false));
        })
}

export function callHandlerSetValue(id, handler, param, valueHandler) {
    let data = {
        id: id,
        param: JSON.stringify(param)
    }

    store.dispatch(loading(true));
    axios.post(`${PIXLET_API_BASE}/api/v1/handlers/` + handler, JSON.stringify(data))
        .then(res => {
            store.dispatch(update({ id: id, value: res.data }));
            valueHandler(res.data);
        })
        .catch(err => {
            // TODO: make sure this clears.
            const msg = parseErrorMessage(handler, err);
            store.dispatch(setError({ id: msg, message: msg }));
        })
        .then(() => {
            store.dispatch(loading(false));
        })
}