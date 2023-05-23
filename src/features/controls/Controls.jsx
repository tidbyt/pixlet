import React from 'react';
import { useSelector } from 'react-redux';

import { Button, Stack, Grid } from '@mui/material';

export default function Controls() {
    const preview = useSelector(state => state.preview);

    let imageType = 'webp';
    if (PIXLET_WASM) {
        imageType = 'gif';
    }

    function downloadPreview() {
        const date = new Date().getTime();
        const element = document.createElement("a");

        // convert base64 to raw binary data held in a string
        let byteCharacters = atob(preview.value.webp);

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

    function resetSchema() {
        history.replaceState(null, '', location.pathname);
        window.location.reload();
    }

    return (
        <Stack sx={{ marginTop: '32px' }} spacing={2} direction="row">
            <Button variant="contained" onClick={() => downloadPreview()}>Save</Button>
            <Button variant="contained" onClick={() => resetSchema()}>Reset</Button>
        </Stack>
    );
}