<template>
  <v-container fill-height>
      <v-container>
          <v-row>
          </v-row>
            <v-row>
                <v-col v-for="t in head" :key="t.rank" cols="12" sm="4">
                    <v-card max-width="300" class="mx-auto d-flex flex-row mb-6 px-4">
                        <div class="pa-2 align-self-center">
                            <h1>{{ t.rank }}</h1>
                        </div>
                        <div class="pa-2">
                            <v-card-title>{{ t.name }}</v-card-title>
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

class Team {
    public rank: number;
    public name: string;
    public score: number;
    constructor(rank: number, name: string, score: number) {this.name = name, this.score = score, this.rank = rank}
}

@Component
export default class Rank extends Vue {
    page: number = 0;
    ordinal: Array<Team> = Array.from({length: 26}, (v,k) => new Team(k + 1, String.fromCharCode(k + 97), 26 - k));
    pageCount: number = Math.ceil((this.ordinal.length - 3) / 8);
    head: Array<Team> = this.ordinal.length > 3 ? this.ordinal.slice(0, 3) : this.ordinal;
    others: Array<Team> = this.ordinal.length > 3 ? this.ordinal.slice(this.page * 8 + 3, (this.page + 1) * 8 + 3) : [];
    
    mounted() {
        console.log(this)
    }
}
</script>
