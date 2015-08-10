/**
 * Created by s.kior on 23.07.2015.
 */
import React from 'react'
import Reflux from 'reflux'

import StatusStore from "../../store/status_store.js"

var ArtGateStatusLog = React.createClass({
    mixins: [
        Reflux.connect(StatusStore, 'status')
    ],
    render(){
        var messages;
        if(this.state.status){
            messages = this.state.status.map(function(message, key){
                return <li key={key}>{message}</li>;
            });
        } else {
            messages = "";
        }
        return (
            <ol>
                {messages}
            </ol>
        );
    }
});

export default ArtGateStatusLog
