import store from '../../store';
import { update } from '../preview/previewSlice';
import { update as updateSchema } from '../schema/schemaSlice';
import { set as setError, clear as clearErrors } from '../errors/errorSlice';

export default class Watcher {
    constructor() {
        this.connect();
    }

    connect() {
        const proto =  document.location.protocol === "https:" ? "wss:" : "ws:";
        this.conn = new WebSocket(proto + '//' + document.location.host + '/api/v1/ws');
        this.conn.open = this.open.bind(this);
        this.conn.onmessage = this.process.bind(this);
        this.conn.onclose = this.close.bind(this);
        setTimeout(this.check.bind(this), 5000);
    }

    open(e) {
        console.log('[watcher] connection established');
        store.dispatch(clearErrors());
    }

    process(e) {
        console.log('[watcher] received new message');
        const data = JSON.parse(e.data);

        switch (data.type) {
            case 'webp':
                store.dispatch(update({
                    webp: data.message,
                }));
                store.dispatch(clearErrors());
                break;
            case 'schema':
                store.dispatch(updateSchema(JSON.parse(data.message)));
                break;
            case 'error':
                store.dispatch(setError({ id: data.message, message: data.message }));
                break;
            default:
                console.log(`[watcher] unknown type ${data.type}`);
        }
    }

    check() {
        if (this.conn.readyState === WebSocket.CONNECTING) {
            console.log('[watcher] connection timed out');
            this.reconnect();
        }
    }

    close(e) {
        let msg = `lost connection to pixlet, please refresh page: ${e.code}`;
        store.dispatch(setError({ id: msg, message: msg }));
        // TODO: we should in theory be able to reconnect here.
        // this.reconnect();
    }

    reconnect() {
        console.log('[watcher] reestablishing connection');
        this.connect();
    }
}