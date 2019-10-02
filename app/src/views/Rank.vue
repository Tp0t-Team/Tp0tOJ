<template>
  <v-container fill-height>
      <v-container>
          <v-row>
          </v-row>
            <v-row>
                <v-col v-for="t in head" :key="t.rank" cols="12" sm="4">
                    <v-card max-width="300" class="mx-auto d-flex flex-row mb-6 px-4">
                        <div class="pa-2 align-self-center">
                            <v-avatar v-if="t.avatar" size="40px">
                                <img alt="Avatar" src="t.avatar">
                            </v-avatar>
                            <v-avatar v-else size="40px" color="orange">
                                <span class="white--text headline">{{ t.rank }}</span>
                            </v-avatar>
                        </div>
                        <div class="pa-2">
                            <v-card-title>{{ t.rank }} {{ t.name }}</v-card-title>
                            <v-card-text>{{ t.score }}</v-card-text>
                        </div>
                    </v-card>
                </v-col>
            </v-row>
            <v-simple-table class="ma-4">
                <template v-slot:default>
                <thead>
                    <tr>
                        <th class="text-left">Rank</th>
                        <th class="text-left">Team/Individual Name</th>
                        <th class="text-left">Score</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="t in others" :key="t.rank">
                        <td>{{ t.rank }}</td>
                        <td>{{ t.name }}</td>
                        <td>{{ t.score }}</td>
                    </tr>
                </tbody>
                </template>
            </v-simple-table>
            <v-pagination v-model="page" :page="page" :length="pageCount"></v-pagination>
      </v-container>
  </v-container>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import gql from "graphql-tag";
import { RankDesc, RankResult } from "@/struct";

@Component
export default class Rank extends Vue {
    ranks: RankDesc[] = [];
    page: number = 0;
    get pageCount(): number {
        return Math.ceil((this.ranks.length - 3) / 8);
    }
    get head(): RankDesc[] {
        return this.ranks.length > 3 ? this.ranks.slice(0, 3) : this.ranks;
    }
    get others(): RankDesc[] {
        return this.ranks.length > 3 ? this.ranks.slice((this.page - 1) * 8 + 3, (this.page) * 8 + 3) : [];
    }

    mounted() {
        this.page = parseInt(this.$route.params.page);
        this.ranks = Array.from({length: 26}, (v,k) => ({rank: k + 1, name: String.fromCharCode(k + 97), avatar: "", score: 26 - k}))
        
    }
}
</script>
