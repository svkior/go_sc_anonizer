/**
 * Created by s.kior on 23.07.2015.
 */
import Reflux from 'reflux';

var WSActions = Reflux.createActions([
    "sendMessage",  // Послать сообщение по каналу. Должно быть {}. Функция сама делает JSON.Stringify
    "postMsg",     // Посылает {} Такого вида, которого нужно
    "connected",    // Экшн вызывается, если соединение установлено
    "disconnected",  // Экшн вызывается, если соединение разорвано
    // wsSubscribe - подписывание на обновление коллекции
    "wsSubscribe",
    // wsRegisterFeed - Создание потока. Все сообщения с этим потоком перенаправляются
    // На указанный обработчик
    "wsRegisterFeed"
]);

export default WSActions;