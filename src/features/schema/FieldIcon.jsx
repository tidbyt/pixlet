import React from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { findIconDefinition } from '@fortawesome/fontawesome-svg-core';

import Icon from '@mui/material/Icon';

import styles from './styles.css';


export default function FieldIcon(props) {
    const iconName = props.icon;
    const faIconName = iconName.replace(/[A-Z]/g, m => "-" + m.toLowerCase());

    let icoDef = findIconDefinition({ prefix: 'fas', iconName: faIconName });
    if (icoDef) {
        return <FontAwesomeIcon icon={icoDef} />;
    }

    icoDef = findIconDefinition({ prefix: 'fab', iconName: faIconName });
    if (icoDef) {
        return <FontAwesomeIcon icon={icoDef} />;
    }

    return (
        <Icon className={styles.materialIcons}>{iconName}</Icon>
    )
}