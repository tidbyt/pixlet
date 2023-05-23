import React, { useState, useEffect } from 'react';
import { useSelector, useDispatch } from 'react-redux';

import Autocomplete from '@mui/material/Autocomplete';
import TextField from '@mui/material/TextField';

import { set, remove } from '../../config/configSlice';
import { callHandler } from '../../handlers/actions';


export default function Typeahead({ field }) {
    const [value, setValue] = useState(null);
    const config = useSelector(state => state.config);
    const dispatch = useDispatch();
    const handlerResults = useSelector(state => state.handlers)

    useEffect(() => {
        if (field.id in config) {
            setValue(JSON.parse(config[field.id].value));
        }
    }, [config])

    const onChange = (event, newValue) => {
        if (newValue) {
            setValue(newValue);
            dispatch(set({
                id: field.id,
                value: JSON.stringify(newValue),
            }))
        } else {
            setValue(null);
            dispatch(remove(field.id));
        }
    }

    let options = [];
    if (field.id in handlerResults.values) {
        options = handlerResults.values[field.id];
    }

    return (
        <Autocomplete
            fullWidth
            disablePortal
            value={value}
            onInputChange={(event, v) => {
                callHandler(field.id, field.handler, v);
            }}
            onChange={onChange}
            options={options}
            getOptionLabel={(option) => option.display}
            renderInput={(params) => <TextField fullWidth {...params} label={field.name} />}
        />
    )
}