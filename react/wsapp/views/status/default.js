/**
 * Created by s.kior on 23.07.2015.
 */


import {Line} from 'react-chartjs'

import React from 'react';
import Reflux from 'reflux'
import WSActions from "../../actions/wsactions.js"
import WSStore from "../../store/wsstore.js"


//import ArtGateStatusLog from "./log.js"

import ArtGateChatRoom from "../chatroom/chatroom_view.js"


var ArtGateDefault = React.createClass({
    mixins: [
        Reflux.connect(WSStore, 'send')
    ],
    click(e){
        e.preventDefault();
        WSActions.sendMessage("Test")
    },
    pageReload(e){
        e.preventDefault();
        location.reload();
    },
    render(){
        var messages;
        if(this.state.send){
            messages = this.state.send.map(function(message, key){
                return <li key={key}>{message}</li>;
            });
        } else {
            messages = "";
        }
        var labels= [1,2,3,4,5,6,7,8,9,10];
        var datas = [1,2,3,4,5,6,7,8,9,10];

        if(this.state.send) {
            for(var i=0; i < this.state.send.length; i++){
                var mmm = JSON.parse(this.state.send[i]);
                if(mmm.payload && mmm.payload.Alloc){
                    datas[i] = mmm.payload.Alloc;
                }
            }
        }

        var chartData = {
            labels: labels,
            datasets: [
                {
                    label: "My First dataset",
                    fillColor: "rgba(220,220,220,0.2)",
                    strokeColor: "rgba(220,220,220,1)",
                    pointColor: "rgba(220,220,220,1)",
                    pointStrokeColor: "#fff",
                    pointHighlightFill: "#fff",
                    pointHighlightStroke: "rgba(220,220,220,1)",
                    data: datas
                }
            ]
        };


        return (
            <div>
                <h3>Окно статуса</h3>
                <ArtGateChatRoom/>
                <hr/>
                <ul>
                    <li><a href onClick={this.click}>Тестовое сообщение</a></li>
                    <li><a href onClick={this.pageReload}>Перегрузить</a></li>
                </ul>
                <hr/>
            </div>
        )
    }
});

/*
 <Line data={chartData} options={{bezierCurve: false}} width="200" height="200"/>
 <hr/>
 <ol>
 {messages}
 </ol>
 */

export default ArtGateDefault