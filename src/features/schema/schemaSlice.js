import { createSlice } from '@reduxjs/toolkit';

export const schemaSlice = createSlice({
    name: 'schema',
    initialState: {
        loading: false,
        error: '',
        value: {
            version: '1',
            schema: []
        },
    },
    reducers: {
        update: (state = initialState, action) => {
            return { ...state, value: action.payload, loading: false, error: '' }
        },
        loading: (state = initialState, action) => {
            return { ...state, loading: action.payload }
        },
        error: (state = initialState, action) => {
            return { ...state, error: action.payload }
        },
    },
});

export const { update, loading, error } = schemaSlice.actions;
export default schemaSlice.reducer;