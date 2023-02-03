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
import { callHandler } from '../../../handlers/actions';

export default function LocationBased({ field }) {
    const [locationValue, setLocationValue] = useState({
    	// Default to Brooklyn, because that's where tidbyt folks
    	// are and  we can only dispatch a location object which
    	// has all fields set.
    	'lat': 40.6782,
    	'lng': -73.9442,
    	'locality': 'Brooklyn, New York',
    	'timezone': 'America/New_York'
   	});

    const [value, setValue] = useState(field.default);

    const config = useSelector(state => state.config);
    const dispatch = useDispatch();
    const handlerResults = useSelector(state => state.handlers)

    useEffect(() => {
        if (field.id in config) {
            setValue(config[field.id].value);
        } else if (field.default) {
            dispatch(set({
                id: field.id,
                value: field.default,
            }));
        }
        callHandler(field.id, field.handler, JSON.stringify(locationValue));
    }, [])

    const setPart = (partName, partValue) => {
    	let newLocationValue = {...locationValue};
    	newLocationValue[partName] = partValue;
    	setLocationValue(newLocationValue);
        callHandler(field.id, field.handler, JSON.stringify(newLocationValue));
    }

    const onChangeLatitude = (event) => {
        setPart('lat', event.target.value);
    }

    const onChangeLongitude = (event) => {
    	setPart('lng', event.target.value);
    }

    const onChangeLocality = (event) => {
    	setPart('locality', event.target.value);
    }

    const onChangeTimezone = (event) => {
    	setPart('timezone', event.target.value);
    }

    const onChangeOption = (event) => {
        setValue(event.target.value);
        dispatch(set({
            id: field.id,
            value: JSON.stringify({'value': event.target.value})
        }));
    }

    let options = [];
    if (field.id in handlerResults.values) {
        options = handlerResults.values[field.id];
    }

    return (
        <FormControl fullWidth>
            <Typography>Latitude</Typography>
            <InputSlider
            	min={-90}
            	max={90}
            	step={0.1}
            	onChange={onChangeLatitude}
            	defaultValue={locationValue['lat']}
            >
            </InputSlider>
            <Typography>Longitude</Typography>
            <InputSlider
            	min={-180}
            	max={180}
            	step={0.1}
            	onChange={onChangeLongitude}
            	defaultValue={locationValue['lng']}
            >
            </InputSlider>
            <Typography>Locality</Typography>
            <TextField
            	fullWidth
            	variant="outlined"
            	onChange={onChangeLocality}
            	style={{ marginBottom: '0.5rem' }} 
            	defaultValue={locationValue['locality']}
            />
            <Typography>Timezone</Typography>
            <Select
                onChange={onChangeTimezone}
                style={{ marginBottom: '0.5rem' }} 
                defaultValue={locationValue['timezone']}
            >
                {Intl.supportedValuesOf('timeZone').map((zone) => {
                    return <MenuItem value={zone}>{zone}</MenuItem>
                })}
            </Select>
            <Typography>Options for chosen location</Typography>
            <Select
                onChange={onChangeOption}
            >
                {options.map((option) => {
                    return <MenuItem key={option.value} value={option.value}>{option.display}</MenuItem>
                })}
            </Select>
        </FormControl>
    );
}