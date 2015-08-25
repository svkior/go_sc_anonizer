/**
 * Created by svkior on 10/08/15.
 */
import Reflux from 'reflux'

import WSActions from "../actions/wsactions.js"
import StatusActions from "../actions/status_actions.js"
import ChatRoomActions from "../actions/chatroom_actions.js"


var ChatRoomStore = Reflux.createStore({
    listenables: [StatusActions,ChatRoomActions],
    onSendMessage(msg_string){
        WSActions.sendMessage(msg_string);
    },
    onGotMessage(msg){
        //console.log("ChatRoom: Realy got message", msg);
        if(this.connected){
            this.log.push(msg.payload.Message);
            this.trigger(this.log);
        } else {
            this.log = [msg.payload.Message];
            this.connected = true;
            this.trigger(this.log);
        }
    },
    onWsconnected(){
        console.log("ChatRoom: Connect to the Chat");
        WSActions.wsSubscribe('respawn'); // FIXME: Нужно вызывать из правильного места
        WSActions.wsRegisterFeed('chatroom', ChatRoomActions);
    },
    init(){
        this.connected = false
        this.log = ["Нет соединения"];
        this.trigger(this.log);
    },
    getInitialState(){
        return this.log;
    }
});

export default ChatRoomStore