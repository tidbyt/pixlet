import React, { useState, useEffect } from 'react';
import { useSelector, useDispatch } from 'react-redux';

import AdapterDateFns from '@mui/lab/AdapterDateFns';
import LocalizationProvider from '@mui/lab/LocalizationProvider';
import TextField from '@mui/material/TextField';
import DateTimePicker from '@mui/lab/DateTimePicker';

import { set, remove } from '../../config/configSlice'


export default function DateTime({ field }) {
    const [dateTime, setDateTime] = useState(new Date());
    const config = useSelector(state => state.config);
    const dispatch = useDispatch();

    useEffect(() => {
        if (field.id in config) {
            setDateTime(new Date(config[field.id].value));
        }
    }, []);

    const onChange = (timestamp) => {
        if (!timestamp) {
            setDateTime(new Date());
            dispatch(remove(field.id));
            return;
        }

        setDateTime(timestamp);
        dispatch(set({
            id: field.id,
            value: timestamp.toISOString(),
        }));
    }

    return (
        <LocalizationProvider dateAdapter={AdapterDateFns}>
            <DateTimePicker
                renderInput={(props) => <TextField {...props} />}
                label={field.name}
                value={dateTime}
                onChange={onChange}
                onError={console.log}
            />
        </LocalizationProvider>
    );
}