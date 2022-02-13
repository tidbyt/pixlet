import React from 'react';

import { useEffect } from 'react';
import CircularProgress from '@mui/material/CircularProgress';
import Grid from '@mui/material/Grid';


export default function OAuth2Handler() {
    const params = new URLSearchParams(window.location.search);

    useEffect(() => {
        window.addEventListener("message", function (event) {
            if (event.data.message === "requestResult") {
                event.source.postMessage({ "message": "deliverResult", result: { code: params.get("code") } });
            }
        });
    }, []);

    return (
        <Grid
            container
            spacing={0}
            direction="column"
            alignItems="center"
            justifyContent="center"
            style={{ minHeight: '100vh' }}
        >
            <Grid item xs={3}>
                <CircularProgress />
            </Grid>
        </Grid>
    )
}