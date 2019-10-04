<template>
  <v-container fill-width>
    <v-row>
      <v-col v-for="(r,index) in topRank" :key="r.userId" cols="4">
        <v-hover v-slot:default="{ hover }">
          <v-card
            :elevation="hover ? 12 : 2"
            max-width="300"
            class="mx-auto d-flex flex-row mb-6 px-4"
            @click="if(!!$store.state.global.userId)$router.push(`/profile/${r.userId}`)"
          >
            <div class="pa-2 align-self-center">
              <v-avatar size="64" color="blue">
                <span class="headline">{{ r.name[0] }}</span>
              </v-avatar>
            </div>
            <div class="pa-2">
              <v-card-title>
                <v-chip>
                  <v-avatar large left :class="rankColor[index]+' white--text'">{{ index + 1 }}</v-avatar>
                  {{ r.name }}
                </v-chip>
              </v-card-title>
              <v-card-text>{{ r.score }}pt</v-card-text>
            </div>
          </v-card>
        </v-hover>
      </v-col>
    </v-row>
    <v-simple-table class="ma-4">
      <thead>
        <tr>
          <th class="text-left">Rank</th>
          <th class="text-left">Name</th>
          <th class="text-left">Score</th>
        </tr>
      </thead>
      <tbody>
        <tr
          class="table-item"
          v-for="(r,index) in pageRank"
          :key="r.rank"
          @click="$router.push(`/profile/${r.userId}`)"
        >
          <td>{{ pageBase + index + 1 }}</td>
          <td>{{ r.name }}</td>
          <td>{{ r.score }}</td>
        </tr>
      </tbody>
    </v-simple-table>
    <v-row>
      <v-pagination v-model="page" :page="page" :length="pageCount"></v-pagination>
    </v-row>
    <v-snackbar v-model="hasInfo" right bottom :timeout="3000">
      {{ infoText }}
      <v-spacer></v-spacer>
      <v-btn icon>
        <v-icon @click="hasInfo = false">close</v-icon>
      </v-btn>
    </v-snackbar>
  </v-container>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import gql from "graphql-tag";
import { RankDesc, RankResult } from "@/struct";

const UserPerPage = 10;

@Component
export default class Rank extends Vue {
  private rankColor = ["amber", "light-blue", "green"];
  private page: number = 1;

  private ranks: RankDesc[] = [];
  private pageCount: number = 1;

  private infoText: string = "";
  private hasInfo: boolean = false;

  private get topRank() {
    return this.ranks.slice(0, 3);
  }
  private get pageBase() {
    return (this.page - 1) * 10 + 3;
  }
  private get pageRank() {
    return this.ranks.slice(this.pageBase, this.pageBase + UserPerPage);
  }

  async mounted() {
    // example data
    this.ranks = [
      { userId: "1", name: "Zenis", score: "1000" },
      { userId: "2", name: "Mio", score: "800" },
      { userId: "3", name: "DRSN", score: "600" },
      { userId: "1", name: "Zenis", score: "1000" },
      { userId: "2", name: "Mio", score: "800" },
      { userId: "3", name: "DRSN", score: "600" },
      { userId: "1", name: "Zenis", score: "1000" },
      { userId: "2", name: "Mio", score: "800" },
      { userId: "3", name: "DRSN", score: "600" },
      { userId: "1", name: "Zenis", score: "1000" },
      { userId: "2", name: "Mio", score: "800" },
      { userId: "3", name: "DRSN", score: "600" },
      { userId: "1", name: "Zenis", score: "1000" },
      { userId: "2", name: "Mio", score: "800" },
      { userId: "3", name: "DRSN", score: "600" },
      { userId: "1", name: "Zenis", score: "1000" },
      { userId: "2", name: "Mio", score: "800" },
      { userId: "3", name: "DRSN", score: "600" },
      { userId: "1", name: "Zenis", score: "1000" },
      { userId: "2", name: "Mio", score: "800" },
      { userId: "3", name: "DRSN", score: "600" },
      { userId: "1", name: "Zenis", score: "1000" },
      { userId: "2", name: "Mio", score: "800" },
      { userId: "3", name: "DRSN", score: "600" }
    ];
    //
    this.page = parseInt(this.$route.params.page);
    try {
      let res = await this.$apollo.query<RankResult>({
        query: gql`
          query {
            rank {
              message
              userInfos {
                userId
                name
                score
              }
            }
          }
        `
      });
      if (res.errors) throw res.errors.map(v => v.message).join(",");
      if (res.data!.rank.message) throw res.data!.rank.message;
      this.ranks = res.data!.rank.rankResultDescs;
      this.pageCount = Math.floor(
        (this.ranks.length + UserPerPage - 1) / UserPerPage
      );
    } catch (e) {
      this.infoText = e.toString();
      this.hasInfo = true;
    }
  }
}
</script>

<style lang="scss" scoped>
.table-item {
  cursor: pointer;
}
</style>