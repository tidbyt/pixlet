import React, { useState, useCallback, useEffect } from 'react';
import Cropper from 'react-easy-crop';
import { useSelector, useDispatch } from 'react-redux';

import Button from '@mui/material/Button';
import Modal from '@mui/material/Modal';
import PhotoCamera from '@mui/icons-material/PhotoCamera';
import Slider from '@mui/material/Slider';
import Stack from '@mui/material/Stack';
import DeleteIcon from '@mui/icons-material/Delete';
import Box from '@mui/material/Box';

import { set, remove } from '../../../config/configSlice';
import getCroppedImg from './cropImage';
import styles from './styles.css';


export default function RawPhotoSelect({ field }) {
    const config = useSelector(state => state.config);
    const dispatch = useDispatch();
    const [image, setImage] = useState("");

    useEffect(() => {
        if (field.id in config) {
            setImage(config[field.id].value);
        } else if (field.default) {
            setImage(field.default);
            dispatch(set({
                id: field.id,
                value: field.default,
            }));
        }
    }, [config])

    const handleCapture = ({ target }) => {
        const fileReader = new FileReader();
        fileReader.readAsDataURL(target.files[0]);
        fileReader.onload = (e) => {
            let base64String = e.target.result.split(",")[1];
            setImage(base64String);
            dispatch(set({
                id: field.id,
                value: base64String,
            }));

        };
    }

    const handleClear = () => {
        setImage("");
        dispatch(remove(field.id));
    };

    let buttons;

    if (image) {
        buttons = <Stack spacing={2} direction="row">
            <Button
                variant="contained"
                component="label"
                startIcon={<PhotoCamera htmlColor='white' />}
            >
                Upload Image
                <input
                    accept="image/*"
                    type="file"
                    hidden
                    onChange={handleCapture}
                />
            </Button >
            <Button
                variant="contained"
                onClick={handleClear}
                startIcon={<DeleteIcon htmlColor='white' />}
            >
                Clear Image
            </Button >
        </Stack>
    } else {
        buttons = <Button
            variant="contained"
            component="label"
            startIcon={<PhotoCamera htmlColor='white' />}
        >
            Upload Image
            <input
                accept="image/*"
                type="file"
                hidden
                onChange={handleCapture}
            />
        </Button >
    }

    return (
        <React.Fragment>
            {buttons}
        </React.Fragment>
    );
}