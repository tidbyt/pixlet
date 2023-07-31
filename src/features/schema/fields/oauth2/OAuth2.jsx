import React, { useState, useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import OAuth2Login from 'react-simple-oauth2-login';

import Button from '@mui/material/Button';

import { callHandlerSetValue } from '../../../handlers/actions';
import { set as setError } from '../../../errors/errorSlice';
import { set, remove } from '../../../config/configSlice';


export default function OAuth2({ field }) {
    const [loggedIn, setLoggedIn] = useState("");
    const dispatch = useDispatch();
    const config = useSelector(state => state.config);
    const redirectUri = document.location.protocol + "//" + document.location.host + "/oauth-callback"

    useEffect(() => {
        if (field.id in config) {
            setLoggedIn(config[field.id].value);
        }
    }, [config])

    const onSuccess = (response) => {
        if (!response.code) {
            return onFailure("access was not granted");
        }

        callHandlerSetValue(field.id, field.handler, {
            code: response.code,
            client_id: field.client_id,
            redirect_uri: redirectUri,
            grant_type: "authorization_code",
        }, (value) => {
            setLoggedIn(value);
            dispatch(set({
                id: field.id,
                value: value,
            }));
        });
    }

    const logout = () => {
        setLoggedIn("");
        dispatch(remove(field.id));
    }

    const onFailure = (response) => {
        let msg = `failed login: ${response}`;
        dispatch(setError({ id: msg, message: msg }));
        console.error(response);
    }

    const renderButton = (params) => {
        return (
            <Button
                variant="contained"
                onClick={params.onClick}
            >
                Login
            </Button >
        )
    }


    if (loggedIn) {
        return (
            <Button
                variant="contained"
                onClick={logout}
            >
                Logout
            </Button>
        )
    }

    let scope = field.scopes.join(" ");
    return (
        <OAuth2Login
            isCrossOrigin={true}
            authorizationUrl={field.authorization_endpoint}
            responseType="code"
            scope={scope}
            state="abc123"
            clientId={field.client_id}
            redirectUri={redirectUri}
            render={renderButton}
            onSuccess={onSuccess}
            onFailure={onFailure} />
    )
}