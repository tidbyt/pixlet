import React, { useEffect, useState } from 'react';
import { useSelector, useDispatch } from 'react-redux';

import { callGeneratedHandler } from '../../handlers/actions';
import { set as setError } from '../../errors/errorSlice';


export default function Generated({ field }) {
    const [source, setSource] = useState(null);
    const config = useSelector(state => state.config);
    const schema = useSelector(state => state.schema);
    const dispatch = useDispatch();

    useEffect(() => {
        onChange(source);
    }, [config])

    useEffect(() => {
        setSource(getSourceField());
    }, [schema])

    const onChange = (source_field) => {
        if (source_field && source_field.id in config) {
            callGeneratedHandler(field.id, field.handler, config[source_field.id].value);
        }
    }

    const getSourceField = () => {
        if (schema.value.schema.length == 0) {
            return null;
        }

        for (let i = 0; i < schema.value.schema.length; i++) {
            if (schema.value.schema[i].id === field.source) {
                return schema.value.schema[i];
            }
        }

        let msg = `schema.Generated references source that does not exist: ${field.source}`;
        dispatch(setError({ id: msg, message: msg }));
        return null;
    }

    return null;
}