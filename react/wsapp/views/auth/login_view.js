/**
 * Created by svkior on 11/08/15.
 *
 * Стандартная форма для логиненья
 */

import React from 'react'
import Reflux from 'reflux'

import AuthStore from "../../store/auth_store.js"
import AuthActions from "../../actions/auth_actions.js"

var LoginView = React.createClass({
    mixins: [
        Reflux.connect(AuthStore, 'auth')
    ],
    handleSubmit(e){
        e.preventDefault();
        var login = React.findDOMNode(this.refs.login);
        var passwd = React.findDOMNode(this.refs.password);
        AuthActions.doLogin(login.value, passwd.value);
    },
    render(){
        var status;
        if(this.state.auth && this.state.auth.authErr){
            status = this.state.auth.authErr;
        } else {
            status = "";
        }
        return(
            <div>
                <h1>Здесь будет логин</h1>
                <form onSubmit={this.handleSubmit}>
                    <input type="text" placeholder="Имя Пользователя" ref="login"/>
                    <input type="password" placeholder="Пароль" ref="password"/>
                    <input type="submit" name="Login"/>
                </form>
                <h2>{status}</h2>
            </div>
        );
    }
});

export default LoginView;
