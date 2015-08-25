/**
 * Created by s.kior on 23.07.2015.
 */
import React from 'react'
import Reflux from 'reflux'

import ChatRoomStore from "../../store/chatroom_store.js"

var ArtGateStatusLog = React.createClass({
    mixins: [
        Reflux.connect(ChatRoomStore, 'chat')
    ],
    render(){
        var messages;
        if(this.state.chat){
            messages = this.state.chat.map(function(message, key){
                return <li key={key}>{message}</li>;
            });
        } else {
            messages = <li>Нет подключения к серверу!!!</li>;
        }
        return (
            <div>
                <h3> Сообщения от сервера </h3>
                <ul>
                    {messages}
                </ul>
            </div>
        );
    }
});

export default ArtGateStatusLog
