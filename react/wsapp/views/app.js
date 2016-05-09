import React from 'react';
import Reflux from 'reflux';
import {RouteHandler} from 'react-router';



import LoginView from "./auth/login_view.js"
import ConnectingView from "./auth/connecting_view.js"

import AuthStore from "../store/auth_store.js"
import WSStore from "../store/wsstore.js"

import ArtGateMenu from "./helpers/topmenu.js"


var App = React.createClass({
    mixins: [
        Reflux.connect(AuthStore, 'auth'),
        Reflux.connect(WSStore, 'ws')
    ],
    render(){
        var content;
        if(this.state.auth){

            if(this.state.auth.canLogin){

                if(this.state.auth.loggedIn){
                    // Здесь пойдет правильный контент
                    content = (
                        <div>
                            <ArtGateMenu/>
                            <RouteHandler/>
                        </div>
                    );
                } else {
                    // Здесь нужно отображать вид на логинство
                    content = <LoginView/>;
                }
            } else {
                // Здесь нужно сформировать контент пока мы грузимся
                content = <ConnectingView/>;
            }
        } else {
            content = <ConnectingView/>;
        }

        return content;
    }
});

export default App
