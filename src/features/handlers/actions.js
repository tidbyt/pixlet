
import axios from 'axios';
import { update, loading } from './handlerSlice';
import { set as setError } from '../errors/errorSlice';
import store from "../../store";


export function callHandler(id, handler, param) {
    let data = {
        id: id,
        param: param
    }

    store.dispatch(loading(true));
    axios.post('/api/v1/handlers/' + handler, JSON.stringify(data))
        .then(res => {
            store.dispatch(update({ id: id, value: res.data }));
        })
        .catch(err => {
            // TODO: make sure this clears.
            let msg = `could not call handler ${handler} with param ${param}`
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
    axios.post('/api/v1/handlers/' + handler, JSON.stringify(data))
        .then(res => {
            store.dispatch(update({ id: id, value: res.data }));
            valueHandler(res.data);
        })
        .catch(err => {
            // TODO: make sure this clears.
            let msg = `could not call handler ${handler} with param ${param}`
            store.dispatch(setError({ id: msg, message: msg }));
            console.log(err);
        })
        .then(() => {
            store.dispatch(loading(false));
        })
}