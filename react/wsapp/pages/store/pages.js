/**
 * Created by svkior on 26/08/15.
 *
 *  Store для списка страниц. Его нужно получить с сервера.
 */

import Reflux from 'reflux'


import AuthActions from "../../actions/auth_actions.js"
import WSActions from "../../actions/wsactions.js"
import PageActions from "../actions/page_actions.js"

var PageStore = Reflux.createStore({
    listenables: [PageActions],
    onLoggedIn(){
        WSActions.wsRegisterFeed('pages', PageActions);
    },
    onGotMessage(msg){
        switch(msg.type){
            case "page":
                console.log("PAGE: ", msg.payload);
                this.pageList.pages.push(msg.payload);
                this.trigger(this.pageList);
                break
        }
    },
    onEditPage(page){
        console.log('Try to edit page', page);
        this.pageList.page = page;
        this.trigger(this.pageList);
    },
    init(){
        AuthActions.handleLoggedIn(PageActions);
        this.pageList = {
            pages: [],
            page: null
        };
        this.trigger(this.pageList);
    },
    getInitialState(){
        return this.pageList;
    }
});

export default PageStore
