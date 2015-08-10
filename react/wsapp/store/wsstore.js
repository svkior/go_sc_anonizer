/**
 * Created by s.kior on 23.07.2015.
 */
import Reflux from 'reflux'

import StatusActions from "../actions/status_actions.js"
import WSActions from "../actions/wsactions.js"

var WSStore = Reflux.createStore({
    listenables: [WSActions],
    onWsSubscribe: function(channel){
        var m ={
            type: "subscribe",
            name: channel
        };
        if(this.socket){
          this.socket.send(JSON.stringify(m));
        }

    },
    onSendMessage: function(msg){
        var m = {
            type: "status",
            payload: {
                Message: "тестовое сообщение"
            }
        };
        if(this.socket){
            this.socket.send(JSON.stringify(m));
        }
    },
    processMessage(msg){
        var un = JSON.parse(msg);

        if(un && un.type){
            // У нас есть тип сообщения
            if(un.type === "reload"){
                console.log("Change of js detected. Reloading after 3 sec...");
                setTimeout(function(){location.reload();}, 3000);
            }
        }

        this.messages.push(msg);
        if(this.messages.length > 15){
            this.messages.splice (0, 1);
        }
        this.trigger(this.messages);
    },
    connectToWs(){
        this.socket = new WebSocket("ws://" + window.location.host +"/device");
        this.socket.onopen = function(){
            this.messages = [];
            this.trigger(this.messages);
            StatusActions.wsconnected();
        }.bind(this);

        this.socket.onclose = function(){
            console.log("Connection has been closed.");
            WSActions.disconnected();
            setTimeout(this.connectToWs, 2000);
        }.bind(this);
        this.socket.onmessage = function(e){
            this.processMessage(e.data);
        }.bind(this);
    },
    init(){
        if(!window["WebSocket"]){
            alert("Error: Your browser does not support web sockets.");
            this.messages = [];
            this.trigger(this.messages);
        } else {
            this.connectToWs();
        }
    },
    getInitialState(){
        return this.messages;
    }
});

export default WSStore