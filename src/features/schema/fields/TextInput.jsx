import React, { useState, useEffect, useCallback } from 'react';
import { useSelector, useDispatch } from 'react-redux';

import TextField from '@mui/material/TextField';

import { set } from '../../config/configSlice';


export default function TextInput({ field }) {
    const config = useSelector(state => state.config);
    const [value, setValue] = useState(() => {
        if (field.id in config) {
            return config[field.id].value;
        }

        return field.default;
    });
    const dispatch = useDispatch();

    useEffect(() => {
        if (field.id in config) {
            setValue(config[field.id].value);
        } else {
            setValue(field.default);
            dispatch(set({
                id: field.id,
                value: field.default,
            }));
        }
    }, [])

    const debounce = (callback, wait) => {
        let timeoutId = null;
        return (...args) => {
            window.clearTimeout(timeoutId);
            timeoutId = window.setTimeout(() => {
                callback.apply(null, args);
            }, wait);
        };
    };

    const onChange = useCallback(
        debounce((event) => {
            setValue(event.target.value);
            dispatch(set({
                id: field.id,
                value: event.target.value,
            }));
        }, 10),
        [value]
    );

    return (
        <TextField fullWidth defaultValue={value} label={field.name} variant="outlined" onChange={onChange} />
    )
}