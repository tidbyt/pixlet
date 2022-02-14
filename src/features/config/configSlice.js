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
        remove: (state = initialState, action) => {
            let config = state;
            if (action.payload in config) {
                delete config[action.payload];
            }
            return state;
        },
    },
});

export const { set, remove } = configSlice.actions;
export default configSlice.reducer;