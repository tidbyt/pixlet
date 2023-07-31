import { createSlice } from '@reduxjs/toolkit';


export const configSlice = createSlice({
    name: 'config',
    initialState: {
    },
    reducers: {
        set: (state = initialState, action) => {
            let config = state;
            config[action.payload.id] = action.payload
            return state;
        },
        update: (state = initialState, action) => {
            state = action.payload;
            return state;
        },
        clear: (state = initialState, action) => {
            state = {};
            return state;
        },
        remove: (state = initialState, action) => {
            let config = state;
            if (action.payload in config) {
                delete config[action.payload];
            }
            return state;
        },
    },
});

export const { set, remove, update, clear } = configSlice.actions;
export default configSlice.reducer;