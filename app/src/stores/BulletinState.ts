import { Module } from "vuex";
import { BulletinItem } from "@/struct";

interface BulletinState {
  bulletins: BulletinItem[];
}

const BulletinStore: Module<BulletinState, {}> = {
  namespaced: true,
  state: {
    bulletins: []
  },
  mutations: {
    addBulletins(state, items: BulletinItem[]) {
      state.bulletins.push(...items);
    },
    cleanBulletins(state) {
      state.bulletins = [];
    }
  }
};

export default BulletinStore;
