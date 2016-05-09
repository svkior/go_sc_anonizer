/**
 * Created by svkior on 26/08/15.
 */

import React from 'react';
import Reflux from 'reflux';

import PageStore from '../store/pages.js';
import PageActions from '../actions/page_actions.js';

import ReactQuill from 'react-quill';

var PageEdit = React.createClass({
    render(){
        return (
            <div className="form-group">
                <label>Содержание страницы</label>
                <ReactQuill
                    theme="snow"
                    value={this.props.content}
                    onChange={this.props.handler}
                    />
            </div>
        );
    }
});

var PageItem = React.createClass({
    handleClick(e){
        e.preventDefault();
        PageActions.editPage(this.props.page);
    },
    render(){
        return (
            <a href="#" onClick={this.handleClick} key={this.props.key} className="btn btn-default">{this.props.page.DocTitle}</a>
        );

    }
});

function htmlDecode(input){
    var e = document.createElement('div');
    e.innerHTML = input;
    return e.childNodes[0].nodeValue;
}

function stripScripts(s) {
    var div = document.createElement('div');
    div.innerHTML = s;
    var scripts = div.getElementsByTagName('script');
    var i = scripts.length;
    while (i--) {
        scripts[i].parentNode.removeChild(scripts[i]);
    }
    return div.innerHTML;
}

var PageList = React.createClass({
    mixins:[
        Reflux.connect(PageStore, 'pages')
    ],
    getInitialStates(){
        return {
            pages: {
                pages: [],
                page: null
            },
            content: []
        };
    },
    handleChangeContent(content){
        //console.log("Content changed: ", content);
        this.state.conent = content;
    },
    render(){
        //console.log(this.state.pages);
        var content;

        if(this.state.pages.page){
            content = stripScripts(htmlDecode(this.state.pages.page.DocContent));
        } else {
            content = "<h1>Нажмите на заголовок статьи для редактирования</h1>";
        }
        var pages = this.state.pages.pages.map(function(page, key){
            return <PageItem page={page} key={key}/>
        });
        return (
            <div>
                <div className="btn-group">{pages}</div>
                <PageEdit content={content} handler={this.handleChangeContent}/>
            </div>
        );
    }
});

export default PageList;
