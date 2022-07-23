import React, { useEffect } from 'react';
import { useSelector } from 'react-redux';

import refreshSchema from './actions';
import Field from './Field';
import Generated from './fields/Generated';


export default function Schema() {
    const schema = useSelector(state => state.schema);

    useEffect(() => {
        refreshSchema();
    }, []);

    return (
        <div>
            {
                schema.value.schema.map((field) => {
                    if (field.type === "generated") {
                        return <Generated key={field.id} field={field} />
                    }

                    return <Field key={field.id} field={field} />
                })
            }
            {
                schema.generated.schema.map((field) => {
                    // A generated field cannot return a generated field, that
                    // would be chaos!
                    return <Field key={field.id} field={field} />
                })
            }
        </div>
    );
}