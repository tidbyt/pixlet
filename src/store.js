import { configureStore } from '@reduxjs/toolkit'

import configSlice from './features/config/configSlice';
import errorSlice from './features/errors/errorSlice';
import handlerSlice from './features/handlers/handlerSlice';
import paramSlice from './features/config/paramSlice';
import previewSlice from './features/preview/previewSlice';
import schemaSlice from './features/schema/schemaSlice';

export default configureStore({
    reducer: {
        config: configSlice,
        errors: errorSlice,
        handlers: handlerSlice,
        param: paramSlice,
        preview: previewSlice,
        schema: schemaSlice,
    },
});