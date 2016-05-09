/**
 * Created by svkior on 11/08/15.
 */

import Reflux from 'reflux'

import WSActions from "../actions/wsactions.js"
import StatusActions from "../actions/status_actions.js"
import AuthActions from "../actions/auth_actions.js"

var AuthStore = Reflux.createStore({
    listenables: [StatusActions, AuthActions],
    onWsconnected(){
        console.log("Разрешаем логиниться");
        WSActions.wsRegisterFeed('auth', AuthActions);
        WSActions.wsSubscribe('respawn'); // FIXME: Нужно вызывать из правильного места
    },
    onHandleLoggedIn(action){
        console.log("Want to register", action);
        this.auth.handlers.push(action)
        this.trigger(this.auth);
    },
    onDoLogin(login, passwd){
        console.log("Try to login with username: ", login, " and password ", passwd);
        WSActions.postMsg(
            {
                type: "login",
                name: "auth",
                broadcast: false,
                payload:{
                    login: login,
                    passwd: passwd
                }

            }
        );
    },
    onDoLogout(){
        console.log("Do Logout");
        WSActions.postMsg(
            {
                type: "logout",
                name: "auth"
            }
        );
    },
    onGotMessage(msg){
        //console.log("Got Message Auth:");
        //console.log(msg);
        switch(msg.type){
            case "subscribed":
                console.log("AUTH IS GOT SUBSCRIBED!!!");
                var session = sessionStorage.getItem("sessionID");
                if(session){
                    console.log("We got session: ", session);
                    // Пытаемся восстановить сессию
                     WSActions.postMsg(
                     {
                        type:"restore",
                        name:"auth",
                        broadcast: false,
                        payload:{
                            session: session
                        }
                     }
                     );
                } else {
                    this.auth.canLogin = true;
                    this.trigger(this.auth)
                }
                break;
            case "authoff":
                //console.log("Auth OFF!!!");
                this.auth.authErr = "";
                this.auth.loggedIn = false;
                this.auth.canLogin = true;
                //FIXME: Раскомментировать
                //sessionStorage.removeItem("sessionID");
                this.trigger(this.auth);
                break;
            case "authok":
                console.log("Auth OK!!!!");
                //console.log(msg.payload.authentication);
                sessionStorage.setItem("sessionID", msg.payload.authentication);
                this.auth.authErr = "";
                this.auth.loggedIn = true;
                this.auth.canLogin = true;
                this.trigger(this.auth);
                // Сказали всем, что вошли
                for(var i=0; i < this.auth.handlers.length; i++){
                    console.log(this.auth.handlers[i]);
                    this.auth.handlers[i].loggedIn();
                }
                break;
            case "autherr":
                console.log("Auth ERR!!!!");
                this.auth.canLogin = true;
                this.auth.authErr = "Can not auth: wrong password or login";
                this.trigger(this.auth);
                break;
        }
    },
    init(){
        this.auth = {
            canLogin: false,
            loggedIn: false,
            authPass: "",
            authErr: "",
            handlers: []
        };
        this.trigger(this.auth);
    },
    getInitialState(){
        return this.auth;
    }
});

export default AuthStore;