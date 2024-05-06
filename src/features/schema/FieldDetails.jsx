import React from 'react';

import OAuth2 from './fields/oauth2/OAuth2';
import PhotoSelect from './fields/photoselect/PhotoSelect';
import RawPhotoSelect from './fields/photoselect/RawPhotoSelect';
import Toggle from './fields/Toggle';
import Color from './fields/Color';
import DateTime from './fields/DateTime';
import Dropdown from './fields/Dropdown';
import LocationBased from './fields/location/LocationBased';
import LocationForm from './fields/location/LocationForm';
import TextInput from './fields/TextInput';
import Typeahead from './fields/Typeahead';
import Typography from '@mui/material/Typography';


export default function FieldDetails({ field }) {
    switch (field.type) {
        case 'datetime':
            return <DateTime field={field} />
        case 'dropdown':
            return <Dropdown field={field} />
        case 'location':
            return <LocationForm field={field} />
        case 'locationbased':
            return <LocationBased field={field} />
        case 'oauth2':
            return <OAuth2 field={field} />
        case 'png':
            return <PhotoSelect field={field} />
        case 'text':
            return <TextInput field={field} />
        case 'onoff':
            return <Toggle field={field} />
        case 'typeahead':
            return <Typeahead field={field} />
        case 'color':
            return <Color field={field} />
        default:
            return <Typography>Unsupported type: {field.type}</Typography>
    }
}