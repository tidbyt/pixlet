import React from 'react';
import { useSelector } from 'react-redux';

import { Paper } from '@mui/material';

import styles from './styles.css';

export default function Preview() {
    const preview = useSelector(state => state.preview);

    let webp = 'UklGRhoAAABXRUJQVlA4TA4AAAAvP8AHAAcQEf0PRET/Aw==';
    if (preview.value.webp) {
        webp = preview.value.webp;
    }

    return (
        <Paper sx={{ bgcolor: "black" }}>
            <img src={"data:image/webp;base64," + webp} className={styles.image} />
        </Paper>
    );
}