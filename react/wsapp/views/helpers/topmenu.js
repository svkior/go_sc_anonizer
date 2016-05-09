/**
 * Created by svkior on 18/08/15.
 */
import React from 'react';
import Reflux from 'reflux';

import AuthActions from "../../actions/auth_actions.js"

var ArtGateMenu = React.createClass({
    click(e){
        e.preventDefault();
        AuthActions.doLogout()
    },
    render(){
        return (
            <div>
                <a onClick={this.click} href="/">Logout</a>
            </div>
        );
    }
});

export default ArtGateMenu;