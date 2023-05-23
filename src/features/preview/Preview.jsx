import React from 'react';
import { useSelector } from 'react-redux';

import { Paper } from '@mui/material';

import styles from './styles.css';

export default function Preview() {
    const preview = useSelector(state => state.preview);

    let displayType = 'data:image/webp;base64,';
    if (PIXLET_WASM) {
        displayType = 'data:image/gif;base64,';
    }

    let webp = 'UklGRhoAAABXRUJQVlA4TA4AAAAvP8AHAAcQEf0PRET/Aw==';
    if (preview.value.webp) {
        webp = preview.value.webp;
    }

    return (
        <Paper sx={{ bgcolor: "black" }}>
            <img src={displayType + webp} className={styles.image} />
        </Paper>
    );
}