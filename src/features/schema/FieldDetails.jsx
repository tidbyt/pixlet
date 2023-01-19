import React from 'react';

import OAuth2 from './fields/oauth2/OAuth2';
import PhotoSelect from './fields/photoselect/PhotoSelect';
import Toggle from './fields/Toggle';
import DateTime from './fields/DateTime';
import Dropdown from './fields/Dropdown';
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
            return <Typography>schema.LocationBased() is not yet supported in pixlet, but is supported in the community repo. Be on the lookout for this field to be available in a future release.</Typography>
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
        default:
            return <Typography>Unsupported type: {field.type}</Typography>
    }
}