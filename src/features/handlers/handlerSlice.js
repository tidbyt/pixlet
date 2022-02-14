import { createSlice } from '@reduxjs/toolkit';

export const handlerSlice = createSlice({
    name: 'handler',
    initialState: {
        loading: false,
        values: {}
    },
    reducers: {
        update: (state = initialState, action) => {
            let up = state;
            up.values[action.payload.id] = action.payload.value
            return up;
        },
        loading: (state = initialState, action) => {
            return { ...state, loading: action.payload }
        },
    },
});

export const { update, loading } = handlerSlice.actions;
export default handlerSlice.reducer;