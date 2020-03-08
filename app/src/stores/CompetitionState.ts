import { Module } from "vuex";
import { CompetitionInfo } from "@/struct";

interface CompetitionState {
  competition: boolean;
  registerAllow: boolean;
  beginTime: number;
  endTime: number;
}

const CompetitionStore: Module<CompetitionState, {}> = {
  namespaced: true,
  state: {
    competition: false,
    registerAllow: true,
    beginTime: Date.now(),
    endTime: Date.now()
  },
  mutations: {
    setCompetitionInfo(state, competition: CompetitionInfo) {
      state.competition = competition.competition;
      state.registerAllow = competition.registerAllow;
      state.beginTime = new Date(competition.beginTime).getTime();
      state.endTime = new Date(competition.endTime).getTime();
    }
  }
};

export default CompetitionStore;
