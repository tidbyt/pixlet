import React, { useState, useCallback, useEffect } from 'react';
import Cropper from 'react-easy-crop';
import { useDispatch } from 'react-redux';

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


export default function PhotoSelect({ field }) {
    const [crop, setCrop] = useState({ x: 0, y: 0 });
    const dispatch = useDispatch();
    const [zoom, setZoom] = useState(1);
    const [open, setOpen] = useState(false);
    const [image, setImage] = useState("");
    const [croppedImage, setCroppedImage] = useState("");
    const [croppedAreaPixels, setCroppedAreaPixels] = useState(null);

    useEffect(() => {
        if (croppedImage) {
            dispatch(set({
                id: field.id,
                value: croppedImage,
            }))
        } else if (field.default) {
            dispatch(set({
                id: field.id,
                value: field.default,
            }));
        } else {
            dispatch(remove(field.id));
        }
    }, [croppedImage]);

    const onCropComplete = useCallback((croppedArea, croppedAreaPixels) => {
        setCroppedAreaPixels(croppedAreaPixels)
    }, []);

    const handleOpen = () => setOpen(true);

    const handleClose = () => setOpen(false);

    const handleCapture = ({ target }) => {
        const fileReader = new FileReader();
        fileReader.readAsDataURL(target.files[0]);
        fileReader.onload = (e) => {
            setImage(e.target.result);
            handleOpen();
        };
    }

    const handleSave = useCallback(() => {
        const croppedImage = getCroppedImg(image, croppedAreaPixels);
        setCroppedImage(croppedImage);
        handleClose();
    }, [croppedAreaPixels]);

    const handleClear = useCallback(() => {
        setImage("");
        setCroppedImage("");
        setCroppedAreaPixels(null);
    }, []);

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
            <Modal
                open={open}
                onClose={handleClose}
                aria-labelledby="modal-modal-title"
                aria-describedby="modal-modal-description"
            >
                <div className={styles.imageCropper}>
                    <div className={styles.cropContainer}>
                        <Cropper
                            image={image}
                            crop={crop}
                            zoom={zoom}
                            aspect={2 / 1}
                            onCropChange={setCrop}
                            onCropComplete={onCropComplete}
                            onZoomChange={setZoom}
                            objectFit="horizontal-cover"
                        />
                    </div>
                    <div className={styles.controls}>
                        <Stack spacing={10} direction="row">
                            <Box sx={{ width: 300 }}>
                                <Slider
                                    value={zoom}
                                    aria-label="Zoom"
                                    valueLabelDisplay="auto"
                                    min={1}
                                    max={3}
                                    step={0.1}
                                    onChange={(e) => {
                                        setZoom(e.target.value)
                                    }}
                                />
                            </Box>
                            <Button
                                variant="contained"
                                onClick={handleSave}
                            >
                                Save
                            </Button>

                        </Stack>
                    </div>
                </div>
            </Modal>
        </React.Fragment>
    );
}