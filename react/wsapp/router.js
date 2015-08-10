import React from 'react';

import Router, {Route, DefaultRoute, Link, RouteHandler} from 'react-router';

import App from "./views/app.js"
import ArtGateDefault from "./views/status/default.js"

var routes = (
    <Route name="app" path="/" handler={App}>
        <DefaultRoute handler={ArtGateDefault}/>
    </Route>
);

function analytics(state){
    if(typeof ga !== 'undefined'){
        window.ga('send', 'pageview',{'page': state.path});
    }
}

export default function(agClient){
    agClient.router = Router.run(routes, Router.HistoryLocation, function(Handler, state){
        React.render(<Handler/>, document.body);
        analytics(state);
    })
}
