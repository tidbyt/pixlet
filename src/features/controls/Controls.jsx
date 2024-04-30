import React from 'react';
import { useSelector, useDispatch } from 'react-redux';

import { Button, Stack } from '@mui/material';
import { resetConfig, setConfig } from '../config/actions';
import { set } from '../config/configSlice';

export default function Controls() {
    const preview = useSelector(state => state.preview);
    const config = useSelector(state => state.config);
    const schema = useSelector(state => state.schema);
    const dispatch = useDispatch();

    let imageType = 'webp';
    if (PIXLET_WASM || preview.value.img_type === "gif") {
        imageType = 'gif';
    }

    function downloadPreview() {
        const date = new Date().getTime();
        const element = document.createElement("a");

        // convert base64 to raw binary data held in a string
        let byteCharacters = atob(preview.value.img);

        // create an ArrayBuffer with a size in bytes
        let arrayBuffer = new ArrayBuffer(byteCharacters.length);

        // create a new Uint8Array view
        let uint8Array = new Uint8Array(arrayBuffer);

        // assign the values
        for (let i = 0; i < byteCharacters.length; i++) {
            uint8Array[i] = byteCharacters.charCodeAt(i);
        }

        const file = new Blob([uint8Array], { type: 'image/' + imageType });
        element.href = URL.createObjectURL(file);
        element.download = `tidbyt-preview-${date}.${imageType}`;
        document.body.appendChild(element); // Required for this to work in FireFox
        element.click();
    }

    function downloadConfig() {
        const date = new Date().getTime();
        const element = document.createElement("a");
        const jsonData = config;

        // Use Blob object for JSON
        const file = new Blob([JSON.stringify(jsonData)], { type: 'application/json' });
        element.href = URL.createObjectURL(file);
        element.download = `config-${date}.json`;
        document.body.appendChild(element); // Required for this to work in FireFox
        element.click();
    }

    function selectConfig() {
        const input = document.createElement('input');
        input.type = 'file';
        input.accept = 'application/json';

        input.onchange = function (event) {
            const file = event.target.files[0];
            if (file.type !== "application/json") {
                return;
            }

            const reader = new FileReader();

            reader.onload = function () {
                let contents = reader.result;
                let json = JSON.parse(contents);
                setConfig(json);
            };

            reader.onerror = function () {
                console.log(reader.error);
            };

            reader.readAsText(file);
        };

        input.click();
    }


    function resetSchema() {
        history.replaceState(null, '', location.pathname);
        resetConfig();
        schema.value.schema.forEach((field) => {
            if (field.default) {
                dispatch(set({
                    id: field.id,
                    value: field.default,
                }));
            };
        });
    };

    return (
        <Stack sx={{ marginTop: '32px' }} spacing={2} direction="row">
            <Button variant="outlined" onClick={() => selectConfig()}>Open Config</Button>
            <Button variant="outlined" onClick={() => downloadConfig()}>Save Config</Button>
            <Button variant="outlined" onClick={() => resetSchema()}>Reset</Button>
            <Button variant="contained" onClick={() => downloadPreview()}>Export Image</Button>
        </Stack>
    );
}