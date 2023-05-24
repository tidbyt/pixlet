import React from 'react';
import { useSelector } from 'react-redux';

import Container from '@mui/material/Container';
import Box from '@mui/material/Box';
import Grid from '@mui/material/Grid';

import AppBar from './features/appbar/AppBar';
import ConfigManager from './features/config/ConfigManager';
import ErrorManager from './features/errors/ErrorManager';
import ErrorSnackbar from './features/errors/ErrorSnackbar';
import ParamSetter from './features/config/ParamSetter';
import Preview from './features/preview/Preview';
import Schema from './features/schema/Schema';
import WatcherManager from './features/watcher/WatcherManager';
import Controls from './features/controls/Controls';
import { Typography } from '@mui/material';


export default function Main() {
    const schema = useSelector(state => state.schema);

    let size = 12;
    if (schema.value.schema.length > 0) {
        size = 8;
    }

    function iOS() {
        return [
            'iPad Simulator',
            'iPhone Simulator',
            'iPod Simulator',
            'iPad',
            'iPhone',
            'iPod'
        ].includes(navigator.platform)
    };

    if (PIXLET_WASM && iOS()) {
        return (
            <ErrorSnackbar >
                <AppBar />
                <Container maxWidth='xl' sx={{ marginTop: '32px' }}>
                    <Box sx={{ flexGrow: 1 }}>
                        <Grid container spacing={4}>
                            <Grid item xs={12} lg={12}>
                                <Typography variant='h4' sx={{ textAlign: 'center' }} color='text.secondary'>
                                    Sorry, iOS is not supported.
                                </Typography>
                            </Grid>
                            <Grid item xs={12} lg={12}>
                                <Typography sx={{ textAlign: 'center' }} color='text.secondary'>
                                    Please try again on a desktop browser.
                                </Typography>
                            </Grid>
                        </Grid>
                    </Box>
                </Container>
            </ErrorSnackbar>
        );
    }


    return (
        <ErrorSnackbar >
            <ParamSetter />
            <ConfigManager />
            <ErrorManager />
            <WatcherManager />

            <AppBar />
            <Container maxWidth='xl' sx={{ marginTop: '32px' }}>
                <Box sx={{ flexGrow: 1 }}>
                    <Grid container spacing={4}>
                        <Grid item xs={12} lg={size}>
                            <Preview scale={10} />
                            <Controls />
                        </Grid>
                        <Grid item xs={12} lg={4}>
                            <Schema />
                        </Grid>
                    </Grid>
                </Box>
            </Container>
        </ErrorSnackbar>
    )
}