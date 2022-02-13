import { useEffect } from 'react';
import { useSelector } from 'react-redux';
import { useSnackbar } from 'notistack';


export default function ErrorManager() {
    const { enqueueSnackbar, closeSnackbar } = useSnackbar();
    const errors = useSelector(state => state.errors);

    useEffect(() => {
        for (const id in errors.active) {
            enqueueSnackbar(errors.active[id].message, { key: id, variant: 'error', persist: true });
        }
        for (const id in errors.inactive) {
            closeSnackbar(id);
        }
    }, [errors]);

    return null;
}