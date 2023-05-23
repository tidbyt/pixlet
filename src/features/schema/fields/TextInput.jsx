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
            if (config[field.id].value != value) {
                setValue(config[field.id].value);
            }
        } else if (field.default) {
            if (field.default != value) {
                setValue(field.default);
            }
            dispatch(set({
                id: field.id,
                value: field.default,
            }));
        }
    }, [config])

    const debounce = (callback, wait) => {
        let timeoutId = null;
        return (...args) => {
            window.clearTimeout(timeoutId);
            timeoutId = window.setTimeout(() => {
                callback.apply(null, args);
            }, wait);
        };
    };

    const onChange = (event) => {
        setValue(event.target.value);
        if (PIXLET_WASM) {
            dispatch(set({
                id: field.id,
                value: event.target.value,
            }));
        } else {
            debounceConfig(event);
        }
    }

    const debounceConfig = useCallback(
        debounce((event) => {
            dispatch(set({
                id: field.id,
                value: event.target.value,
            }));
        }, 10),
        [value]
    );

    return (
        <TextField fullWidth value={value} label={field.name} variant="outlined" onChange={onChange} />
    )
}