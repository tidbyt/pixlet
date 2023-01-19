import React, { useState, useEffect } from 'react';
import { useSelector, useDispatch } from 'react-redux';

import InputLabel from '@mui/material/InputLabel';
import MenuItem from '@mui/material/MenuItem';
import FormControl from '@mui/material/FormControl';
import Select from '@mui/material/Select';
import TextField from '@mui/material/TextField';
import Typography from '@mui/material/Typography';

import InputSlider from './InputSlider';
import { set } from '../../../config/configSlice';

export default function LocationForm({ field }) {
    const [value, setValue] = useState(field.default);
    const config = useSelector(state => state.config);

    const dispatch = useDispatch();

    useEffect(() => {
        // if (field.id in config) {
        //     setValue(config[field.id].value);
        // } else if (field.default) {
        //     dispatch(set({
        //         id: field.id,
        //         value: field.default,
        //     }));
        // }
    }, [])

    const onChange = (event) => {
        console.log(field);
        // setValue(event.target.value);
        // dispatch(set({
        //     id: field.id,
        //     value: event.target.value,
        // }));
    }

    return (
        <FormControl fullWidth>
            <Typography>Latitude</Typography>
            <InputSlider
            	min={-90}
            	max={90}
            	step={0.1}
            >
            </InputSlider>
            <Typography>Longitude</Typography>
            <InputSlider
            	min={-180}
            	max={180}
            	step={0.1}
            >
            </InputSlider>
            <Typography>Locality</Typography>
            <TextField
            	fullWidth
            	defaultValue="Somewhere"
            	variant="outlined"
            	// onChange={onChange}
            	style={{ marginBottom: '0.5rem' }} 
            />
            <Typography>Timezone</Typography>
            <Select
                // onChange={onChange}
            >
                {Intl.supportedValuesOf('timeZone').map((zone) => {
                    return <MenuItem value={zone}>{zone}</MenuItem>
                })}
            </Select>
        </FormControl>
    );
}