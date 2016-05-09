/**
 * Created by svkior on 11/08/15.
 *
 * Отображаем загрузку до авторизации. При помощи Spin.js
 */

import React from 'react';
import Reflux from 'reflux';


var ConnectingView = React.createClass({
    getInitialState() {
        return {
            windowWidth: window.innerWidth,
            windowHeight: window.innerHeight
        };
    },
    handleResize(){
       this.setState({
           windowWidth: window.innerWidth,
           windowHeight: window.innerHeight
       })
    },
    componentDidMount: function() {
        window.addEventListener('resize', this.handleResize);
    },

    componentWillUnmount: function() {
        window.removeEventListener('resize', this.handleResize);
    },
    render(){
        var sizeX = this.state.windowWidth;
        var sizeY = this.state.windowHeight;

        var viewBox = [
            0, 0, sizeX, sizeY
        ].join(' ');

        return(
            <svg viewBox={viewBox} width={sizeX} height={sizeY}>
                <text x={100} y={sizeY/2} fontFamily="Verdana" fontSize="43pt">
                    Загрузка приложения...
                </text>
            </svg>

        );
    }
});

export default ConnectingView