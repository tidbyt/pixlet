import React, { useEffect } from 'react';
import { useSelector } from 'react-redux';

import refreshSchema from './actions';
import Field from './Field';


export default function Schema() {
    const schema = useSelector(state => state.schema);

    useEffect(() => {
        refreshSchema();
    }, []);

    return (
        <div>
            {schema.value.schema.map((field) => {
                return <Field key={field.id} field={field} />
            })}
        </div>
    );
}