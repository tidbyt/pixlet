import React, { useState, useEffect } from 'react';

import { useSelector, useDispatch } from 'react-redux';
import { ColorPicker, createColor } from "material-ui-color";
import { set } from '../../config/configSlice';


export default function Color({ field }) {
    const [color, setColor] = useState(createColor(field.default));
    const config = useSelector(state => state.config);
    const dispatch = useDispatch();

    useEffect(() => {
        if (field.id in config) {
            setColor(createColor(config[field.id].value));
        } else if (field.default) {
            dispatch(set({
                id: field.id,
                value: field.default,
            }));
        }
    }, [])

    const onChange = (value) => {
        setColor(value);

        // Skip updates that contain an error.
        if (value.hasOwnProperty("error")) {
            return;
        }

        dispatch(set({
            id: field.id,
            value: "#" + value.hex,
        }));
    }

    return (
        <ColorPicker value={color} hideTextfield disablePlainColor onChange={onChange} />
    )
}