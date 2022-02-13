import React from 'react';
import { useSelector } from 'react-redux';

import { Paper } from '@mui/material';

import styles from './styles.css';

export default function Preview() {
    const preview = useSelector(state => state.preview);

    return (
        <Paper sx={{ bgcolor: "black" }}>
            <img src={"data:image/webp;base64," + preview.value.webp} className={styles.image} />
        </Paper>
    );
}