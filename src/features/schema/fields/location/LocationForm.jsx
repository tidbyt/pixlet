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
    const [value, setValue] = useState({
        // Default to Brooklyn, because that's where tidbyt folks
        // are and  we can only dispatch a location object which
        // has all fields set.
        'lat': 40.678,
        'lng': -73.944,
        'locality': 'Brooklyn, New York',
        'timezone': 'America/New_York',
        // But overwrite with app-specific defaults set in config.
        ...field.default
    });

    const config = useSelector(state => state.config);

    const dispatch = useDispatch();

    useEffect(() => {
        if (field.id in config) {
            setValue(JSON.parse(config[field.id].value));
        } else if (field.default) {
            dispatch(set({
                id: field.id,
                value: field.default,
            }));
        }
    }, [config])

    const setPart = (partName, partValue) => {
        let newValue = { ...value };
        newValue[partName] = partValue;
        setValue(newValue);
        dispatch(set({
            id: field.id,
            value: JSON.stringify(newValue),
        }));
    }

    const truncateLatLng = (value) => {
        return String(Number(value).toFixed(3));
    }

    const onChangeLatitude = (event) => {
        setPart('lat', truncateLatLng(event.target.value));
    }

    const onChangeLongitude = (event) => {
        setPart('lng', truncateLatLng(event.target.value));
    }

    const onChangeLocality = (event) => {
        setPart('locality', event.target.value);
    }

    const onChangeTimezone = (event) => {
        setPart('timezone', event.target.value);
    }

    return (
        <FormControl fullWidth>
            <Typography>Latitude</Typography>
            <InputSlider
                min={-90}
                max={90}
                step={0.1}
                onChange={onChangeLatitude}
                defaultValue={value['lat']}
            >
            </InputSlider>
            <Typography>Longitude</Typography>
            <InputSlider
                min={-180}
                max={180}
                step={0.1}
                onChange={onChangeLongitude}
                defaultValue={value['lng']}
            >
            </InputSlider>
            <Typography>Locality</Typography>
            <TextField
                fullWidth
                variant="outlined"
                onChange={onChangeLocality}
                style={{ marginBottom: '0.5rem' }}
                defaultValue={value['locality']}
            />
            <Typography>Timezone</Typography>
            <Select
                onChange={onChangeTimezone}
                defaultValue={value['timezone']}
            >
                {Intl.supportedValuesOf('timeZone').map((zone) => {
                    return <MenuItem value={zone}>{zone}</MenuItem>
                })}
            </Select>
        </FormControl>
    );
}