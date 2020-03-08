import { Module } from "vuex";

interface GlobalState {
  userId: string | null;
  role: string | null;
}

const GlobalStateStore: Module<GlobalState, {}> = {
  namespaced: true,
  state: {
    userId: null,
    role: null
  },
  mutations: {
    setUserIdAndRole(state, { userId, role }) {
      state.userId = userId;
      state.role = role;
    },
    resetUserIdAndRole(state) {
      state.userId = null;
      state.role = null;
    }
  },
  actions: {}
};

export default GlobalStateStore;
