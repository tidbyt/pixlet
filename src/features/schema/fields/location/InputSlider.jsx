// Largely based on https://mui.com/material-ui/react-slider/#InputSlider.js
import React, { useState } from 'react';

import { styled } from '@mui/material/styles';
import Box from '@mui/material/Box';
import Grid from '@mui/material/Grid';
import Typography from '@mui/material/Typography';
import Slider from '@mui/material/Slider';
import MuiInput from '@mui/material/Input';

const Input = styled(MuiInput)`
  width: 80px;
`;

export default function InputSlider({ min, max, step, defaultValue, onChange}) {
  const [value, setValue] = useState(defaultValue);

  const handleSliderChange = (event, newValue) => {
    setValue(newValue);
    onChange(event);
  };

  const handleInputChange = (event) => {
    if (event.target.value === '') {
      setValue('');
      onChange(event);
      return;
    }
    const value = Number(event.target.value);
    if (value < min) {
      setValue(min);
    } else if (value > max) {
      setValue(max);
    } else {
      setValue(value);
    }
    onChange(event);
  };

  const handleBlur = () => {
    if (value < min) {
      setValue(min);
    } else if (value > max) {
      setValue(max);
    }
  };

  return (
    <Box sx={{ width: 250 }}>
      <Grid container spacing={2} alignItems="center">
        <Grid item xs>
          <Slider
            value={value}
            min={min}
            max={max}
            step={step}
            onChange={handleSliderChange}
            aria-labelledby="input-slider"
          />
        </Grid>
        <Grid item>
          <Input
            value={value}
            size="small"
            onChange={handleInputChange}
            onBlur={handleBlur}
            inputProps={{
              step: {step},
              min: {min},
              max: {max},
              type: 'number',
              'aria-labelledby': 'input-slider',
            }}
          />
        </Grid>
      </Grid>
    </Box>
  );
}