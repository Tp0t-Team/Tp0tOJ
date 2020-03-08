import Vue from "vue";
import Vuex from "vuex";
import GlobalStateStore from "@/stores/GlobalState";

Vue.use(Vuex);

export default new Vuex.Store({
  state: {},
  mutations: {},
  actions: {},
  modules: {
    global: GlobalStateStore
  }
});
