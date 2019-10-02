<template>
  <v-container fill-width>
    <v-row>
      <v-col v-for="(t,index) in topUser" :key="t.userId" cols="4">
        <v-hover v-slot:default="{ hover }">
          <v-card
            :elevation="hover ? 12 : 2"
            max-width="300"
            class="mx-auto d-flex flex-row mb-6 px-4"
            @click="$router.push(`/profile/${t.userId}`)"
          >
            <div class="pa-2 align-self-center">
              <v-avatar size="64" color="blue">
                <span class="headline">{{ t.name[0] }}</span>
              </v-avatar>
            </div>
            <div class="pa-2">
              <v-card-title>
                <v-chip>
                  <v-avatar large left :class="rankColor[index]+' white--text'">{{ index + 1 }}</v-avatar>
                  {{ t.name }}
                </v-chip>
              </v-card-title>
              <v-card-text>{{ t.score }}pt</v-card-text>
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
          v-for="r in ranks"
          :key="r.rank"
          @click="$router.push(`/profile/${r.userId}`)"
        >
          <td>{{ r.rank }}</td>
          <td>{{ r.name }}</td>
          <td>{{ r.score }}</td>
        </tr>
      </tbody>
    </v-simple-table>
    <v-row>
      <v-pagination v-model="page" :page="page" :length="pageCount"></v-pagination>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { RankDesc } from "@/struct";

@Component
export default class Rank extends Vue {
  private rankColor = ["amber", "light-blue", "green"];
  private page: number = 1;

  private topUser: RankDesc[] = [];
  private ranks: RankDesc[] = [];
  private pageCount: number = 2;

  mounted() {
    this.topUser = [
      { rank: 1, userId: "1", name: "Zenis", score: 1000 },
      { rank: 2, userId: "2", name: "Mio", score: 800 },
      { rank: 3, userId: "3", name: "DRSN", score: 600 }
    ];
    this.ranks = [
      { rank: 4, userId: "1", name: "Zenis", score: 1000 },
      { rank: 5, userId: "2", name: "Mio", score: 800 },
      { rank: 6, userId: "3", name: "DRSN", score: 600 },
      { rank: 7, userId: "1", name: "Zenis", score: 1000 },
      { rank: 8, userId: "2", name: "Mio", score: 800 },
      { rank: 9, userId: "3", name: "DRSN", score: 600 },
      { rank: 10, userId: "1", name: "Zenis", score: 1000 },
      { rank: 11, userId: "2", name: "Mio", score: 800 },
      { rank: 12, userId: "3", name: "DRSN", score: 600 }
    ];
    this.page = parseInt(this.$route.params.page);
  }
}
</script>

<style lang="scss" scoped>
.table-item {
  cursor: pointer;
}
</style>