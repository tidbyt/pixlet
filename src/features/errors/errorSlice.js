import { createSlice } from '@reduxjs/toolkit';

export const errorSlice = createSlice({
    name: 'errors',
    initialState: {
        active: {},
        inactive: {},
    },
    reducers: {
        set: (state = initialState, action) => {
            // TODO: Fix this in pixlet.
            if (action.payload.message.includes("didn't export a main() function")) {
                return state;
            }

            // If the error already exists, make no changes.
            if (action.payload.id in state.active) {
                return state;
            }

            let active = {};
            active[action.payload.id] = action.payload;

            return {
                active: active,
                inactive: state.active,
            }
        },
        clear: (state = initialState) => {
            return {
                active: {},
                inactive: state.active,
            }
        },
    },
});

export const { set, remove, clear } = errorSlice.actions;
export default errorSlice.reducer;