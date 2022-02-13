import React from 'react';

import Accordion from '@mui/material/Accordion';
import AccordionSummary from '@mui/material/AccordionSummary';
import AccordionDetails from '@mui/material/AccordionDetails';
import Typography from '@mui/material/Typography';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';

import FieldDetails from './FieldDetails';
import FieldIcon from './FieldIcon';

export default function Field(props) {
    const field = props.field;

    const [expanded, setExpanded] = React.useState(false);

    const handleChange = (panel) => (event, isExpanded) => {
        setExpanded(isExpanded ? panel : false);
    };

    return (
        <Accordion expanded={expanded === 'panel1'} onChange={handleChange('panel1')}>
            <AccordionSummary
                expandIcon={<ExpandMoreIcon />}
                aria-controls="panel1bh-content"
                id="panel1bh-header"
            >
                <Typography sx={{ width: '10%', flexShrink: 0 }}>
                    <FieldIcon icon={field.icon} />
                </Typography>
                <Typography sx={{ width: '33%', flexShrink: 0 }}>
                    {field.name}
                </Typography>
                <Typography sx={{ color: 'text.secondary' }}>{field.description}</Typography>
            </AccordionSummary>
            <AccordionDetails>
                <FieldDetails field={field} />
            </AccordionDetails>
        </Accordion>
    );
}