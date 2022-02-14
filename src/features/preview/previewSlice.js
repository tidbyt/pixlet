import { createSlice } from '@reduxjs/toolkit';

export const previewSlice = createSlice({
    name: 'preview',
    initialState: {
        loading: false,
        value: {
            webp: '',
            title: 'Pixlet',
        }
    },
    reducers: {
        update: (state = initialState, action) => {
            let up = state;

            if ('webp' in action.payload) {
                up.value.webp = action.payload.webp;
            }

            if ('title' in action.payload) {
                up.value.title = action.payload.title;
            }

            return up;
        },
        loading: (state = initialState, action) => {
            return { ...state, loading: action.payload }
        },
    },
});

export const { update, loading } = previewSlice.actions;
export default previewSlice.reducer;