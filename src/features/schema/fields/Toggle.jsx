import React, { useState, useEffect } from 'react';

import { useSelector, useDispatch } from 'react-redux';
import Switch from '@mui/material/Switch';
import { set } from '../../config/configSlice';


export default function Toggle({ field }) {
    const [toggle, setToggle] = useState(JSON.parse(field.default));
    const config = useSelector(state => state.config);
    const dispatch = useDispatch();

    useEffect(() => {
        if (field.id in config) {
            setToggle(JSON.parse(config[field.id].value));
        } else if (JSON.parse(field.default)) {
            dispatch(set({
                id: field.id,
                value: field.default,
            }));
        }
    }, [])

    const onChange = (event) => {
        setToggle(event.target.checked);
        dispatch(set({
            id: field.id,
            value: JSON.stringify(event.target.checked),
        }))
    }

    return (
        <Switch checked={toggle} onChange={onChange} />
    )
}