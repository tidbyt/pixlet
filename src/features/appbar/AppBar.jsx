import React from 'react';

import AppBar from '@mui/material/AppBar';
import Toolbar from '@mui/material/Toolbar';

import Logo from './logo.svg';
import styles from './styles.css';
import { solarized } from '../theme/colors';


export default function NavBar() {
    return (
        <AppBar sx={{ backgroundColor: solarized.base02 }} position="static">
            <Toolbar>
                <div className={styles.title}>
                    <Logo className={styles.logo} />
                </div>
            </Toolbar>
        </AppBar>
    )
}