import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux';
import { BrowserRouter, Route, Routes } from "react-router-dom";

import Main from './Main';
import OAuth2Handler from './features/schema/fields/oauth2/OAuth2Handler';
import store from './store';
import DevToolsTheme from './features/theme/DevToolsTheme';


const App = () => {
    return (
        <Provider store={store}>
            <DevToolsTheme>
                <BrowserRouter>
                    <Routes>
                        <Route exact path="/" element={<Main />} />
                        <Route path="oauth-callback" element={<OAuth2Handler />} />
                    </Routes>
                </BrowserRouter>
            </DevToolsTheme>
        </Provider >
    )
}

ReactDOM.render(<App />, document.getElementById('app'));