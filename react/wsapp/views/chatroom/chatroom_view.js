/**
 * Created by svkior on 10/08/15.
 */
import React from 'react'

import Reflux from 'reflux'

import ChatRoomStore from "../../store/chatroom_store.js"
import ChatRoomActions from "../../actions/chatroom_actions.js"



var ChatMessages = React.createClass({
    componentWillUpdate(){
        var node = this.getDOMNode();
        // Позволяем проматывать сообщения дальше, если мы уже внизу
        this.shouldScrollBottom = node.scrollTop + node.offsetHeight === node.scrollHeight;
        console.log("1");
    },
    componentDidUpdate(){
        console.log(this.shouldScrollBottom);
        if(this.shouldScrollBottom){
            var node = this.getDOMNode();
            node.scrollTop = node.scrollHeight;
        }
    },
    render(){
       return (
           <div style={{overflowY:"scroll", height:"200px"}}>
           <ul>
               {this.props.messages}
           </ul>
       </div>
       );
    }
});

var ArtGateChatRoom = React.createClass({
    mixins: [
        Reflux.connect(ChatRoomStore, 'chat')
    ],
    handleSubmit(e){
        e.preventDefault();
        var comment = React.findDOMNode(this.refs.comment);

        //console.log("Сообщение",comment.value);
        ChatRoomActions.sendMessage(comment.value);
        comment.value = "";
    },
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
                <ChatMessages messages={messages} />
                <form onSubmit={this.handleSubmit}>
                    <input type="text" placeholder="Ваш комментарий" ref="comment"/>
                </form>
            </div>
        );
    }
});

export default ArtGateChatRoom