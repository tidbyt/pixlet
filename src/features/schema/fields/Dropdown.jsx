import React, { useState, useEffect } from 'react';
import { useSelector, useDispatch } from 'react-redux';

import InputLabel from '@mui/material/InputLabel';
import MenuItem from '@mui/material/MenuItem';
import FormControl from '@mui/material/FormControl';
import Select from '@mui/material/Select';

import { set } from '../../config/configSlice';


export default function Dropdown({ field }) {
    const [value, setValue] = useState(field.default);
    const config = useSelector(state => state.config);
    const dispatch = useDispatch();

    useEffect(() => {
        if (field.id in config) {
            setValue(config[field.id].value);
        } else if (field.default) {
            dispatch(set({
                id: field.id,
                value: field.default,
            }));
        }
    }, [])

    const onChange = (event) => {
        setValue(event.target.value);
        dispatch(set({
            id: field.id,
            value: event.target.value,
        }));
    }

    return (
        <FormControl fullWidth>
            <InputLabel>{field.name}</InputLabel>
            <Select
                value={value}
                label={field.name}
                onChange={onChange}
            >
                {field.options.map((option) => {
                    return <MenuItem key={option.value} value={option.value}>{option.display}</MenuItem>
                })}
            </Select>
        </FormControl>
    );
}