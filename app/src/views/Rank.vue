<template>
  <div class="content-col">
    <v-container>
      <v-row>
        <v-col v-for="(r, index) in topRank" :key="r.desc.userId" cols="4">
          <v-hover v-slot:default="{ hover }">
            <v-card
              :elevation="hover ? 12 : 2"
              max-width="300"
              class="mx-auto d-flex flex-row mb-6 px-4"
              @click="
                if (!!$store.state.global.userId)
                  $router.push(`/profile/${r.desc.userId}`);
              "
            >
              <div class="pa-2 align-self-center">
                <v-avatar size="64" color="blue">
                  <!-- <span class="headline">{{ r.name[0] }}</span> -->
                  <user-avatar
                    class="headline white--text"
                    :url="r.desc.avatar"
                    :size="64"
                    :name="r.desc.name"
                  ></user-avatar>
                </v-avatar>
              </div>
              <div class="pa-2">
                <v-card-title>
                  <v-chip>
                    <v-avatar
                      large
                      left
                      :class="rankColor[index] + ' white--text'"
                      >{{ index + 1 }}</v-avatar
                    >
                    {{ r.desc.name.toUpperCase() }}
                  </v-chip>
                </v-card-title>
                <v-card-text>{{ r.desc.score }}pt</v-card-text>
              </div>
            </v-card>
          </v-hover>
        </v-col>
      </v-row>
      <rank-table :value="commonRank" :userPerPage="10" :pageInit="page" />
      <v-snackbar v-model="hasInfo" right bottom :timeout="3000">
        {{ infoText }}
        <!-- <v-spacer></v-spacer> -->
        <template v-slot:action="{ attrs }">
          <v-btn icon>
            <v-icon v-bind="attrs" @click="hasInfo = false">close</v-icon>
          </v-btn>
        </template>
      </v-snackbar>
    </v-container>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import gql from "graphql-tag";
import UserAvatar from "@/components/UserAvatar.vue";
import RankTable from "@/components/RankTable.vue";
import { RankWithIndex, RankResult } from "@/struct";

@Component({
  components: {
    UserAvatar,
    RankTable
  }
})
export default class Rank extends Vue {
  private rankColor = ["amber", "light-blue", "green"];
  private page: number = 1; // only for init

  private ranks: RankWithIndex[] = [];

  private infoText: string = "";
  private hasInfo: boolean = false;

  private get topRank() {
    return this.ranks.slice(0, 3);
  }
  private get commonRank() {
    return this.ranks.slice(3);
  }

  async mounted() {
    this.page = Math.max(parseInt(this.$route.params.page), 1);
    await this.loadData();
  }

  async loadData() {
    try {
      let res = await this.$apollo.query<RankResult>({
        query: gql`
          query {
            rank {
              message
              rankResultDescs {
                userId
                name
                avatar
                score
              }
            }
          }
        `,
        fetchPolicy: "no-cache"
      });
      if (res.errors) throw res.errors.map(v => v.message).join(",");
      if (res.data!.rank.message) throw res.data!.rank.message;
      this.ranks = res
        .data!.rank.rankResultDescs.sort(
          (a, b) => parseInt(b.score) - parseInt(a.score)
        )
        .map((it, index) => ({ index, desc: it } as RankWithIndex));
    } catch (e) {
      this.infoText = e.toString();
      this.hasInfo = true;
    }
  }
}
</script>

<style lang="scss" scoped>
.content-col {
  height: calc(100vh - 96px);
  overflow-y: auto;
}

.table-item {
  cursor: pointer;
}

.progress-bar {
  position: absolute;
  bottom: 0;
  left: 0;
  height: 4px;
  width: 50%;
  background-color: rgb(245, 124, 0);
  transition: 0.5s width linear;
}
</style>
