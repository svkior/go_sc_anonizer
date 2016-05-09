/**
 * Created by svkior on 26/08/15.
 * Events для работы со страницами
 *
 */

import Reflux from 'reflux';


var PageActions = Reflux.createActions([
    // Диспечер сообщений от WS
    "gotMessage",
    "loggedIn",
    // Для редактируемого документа
    "editPage"
]);

export default PageActions;