import Vue from "vue";
import SocketIO from "socket.io-client";

declare module "vue/types/vue" {
  interface Vue {
    $socket: SocketIOClient.Socket;
  }
}
Vue.prototype.$socket = SocketIO.connect("/socketIO");
