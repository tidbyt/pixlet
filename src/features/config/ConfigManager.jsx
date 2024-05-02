import { useEffect } from 'react';
import { useSelector } from 'react-redux';
import { useNavigate } from 'react-router-dom';

import fetchPreview from '../preview/actions';


export default function ConfigManager() {
    const config = useSelector(state => state.config);
    const loading = useSelector(state => state.param.loading);
    const preview = useSelector(state => state.preview);
    const navigate = useNavigate();

    const updatePreviews = (formData, params) => {
        navigate({ search: params.toString() });
        fetchPreview(formData);
    }

    useEffect(() => {
        const formData = new FormData();
        const params = new URLSearchParams();

        Object.entries(config).forEach((entry) => {
            const [id, item] = entry;

            // Not all config values fit inside a query parameter, most notably
            // images. If they don't fit, simply leave them out of the query
            // string. The downside is a refresh will lose that state.
            if (item.value.length < 1024) {
                params.set(id, item.value)
            }

            formData.set(id, item.value);
        });

        if (!loading || !('img' in preview)) {
            updatePreviews(formData, params);
        }
    }, [config]);

    return null;
}