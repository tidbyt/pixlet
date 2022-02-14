import { useEffect } from 'react';
import { useDispatch } from 'react-redux';
import { set } from './configSlice';
import { loading } from './paramSlice';


export default function ParamSetter() {
    const params = new URLSearchParams(document.location.search);
    const dispatch = useDispatch();

    useEffect(() => {
        params.forEach((value, key) => {
            dispatch(set({
                id: key,
                value: value,
            }));
        });
        dispatch(loading(false));
    }, []);

    return null;
};