import { createSlice } from '@reduxjs/toolkit';

export const previewSlice = createSlice({
    name: 'preview',
    initialState: {
        loading: false,
        value: {
            img: '',
            img_type: '',
            title: 'Pixlet',
        }
    },
    reducers: {
        update: (state = initialState, action) => {
            let up = state;

            if ('img' in action.payload) {
                up.value.img = action.payload.img;
            }

            if ('img_type' in action.payload) {
                up.value.img_type = action.payload.img_type;
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