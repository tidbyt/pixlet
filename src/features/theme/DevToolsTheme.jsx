import React from 'react';
import { library } from '@fortawesome/fontawesome-svg-core';
import { fab } from '@fortawesome/free-brands-svg-icons';
import { fas } from '@fortawesome/free-solid-svg-icons';

import ThemeProvider from '@mui/system/ThemeProvider';

import { theme } from './theme';
import './styles.css';


export default function DevToolsTheme(props) {
    library.add(fas, fab);

    return (
        <ThemeProvider theme={theme}>
            {props.children}
        </ThemeProvider>
    );
}