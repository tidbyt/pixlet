import { createSlice } from '@reduxjs/toolkit';


export const paramSlice = createSlice({
    name: 'param',
    initialState: {
        loading: true,
    },
    reducers: {
        loading: (state = initialState, action) => {
            return { loading: action.payload };
        },
    },
});

export const { loading } = paramSlice.actions;
export default paramSlice.reducer;