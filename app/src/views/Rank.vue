<template>
  <div class="content-col">
    <v-container>
      <v-row>
        <div class="chart-container">
          <apex-chart
            type="line"
            height="350"
            width="100%"
            :options="chartOptions"
            :series="series"
          ></apex-chart>
        </div>
      </v-row>
      <!-- <v-row>
        <v-col v-for="(r, index) in topRank" :key="r.userId" cols="4">
          <v-hover v-slot:default="{ hover }">
            <v-card
              :elevation="hover ? 12 : 2"
              max-width="300"
              class="mx-auto d-flex flex-row mb-6 px-4"
              @click="
                if (!!$store.state.global.userId)
                  $router.push(`/profile/${r.userId}`);
              "
            >
              <div class="pa-2 align-self-center">
                <v-avatar size="64" color="blue">
                  <user-avatar
                    class="headline white--text"
                    :url="r.avatar"
                    :size="64"
                    :name="r.name"
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
                    {{ r.name.toUpperCase() }}
                  </v-chip>
                </v-card-title>
                <v-card-text>{{ r.score }}pt</v-card-text>
              </div>
            </v-card>
          </v-hover>
        </v-col>
      </v-row> -->
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
            v-for="(r, index) in pageRank"
            :key="r.rank"
            @click="
              if ($store.state.global.role == 'admin')
                $router.push(`/profile/${r.userId}`);
            "
          >
            <td>{{ pageBase + index + 1 }}</td>
            <td>{{ r.name }}</td>
            <td>{{ r.score }}</td>
          </tr>
        </tbody>
      </v-simple-table>
      <v-row justify="center">
        <v-pagination
          v-model="page"
          :page="page"
          :length="pageCount"
        ></v-pagination>
      </v-row>
      <v-snackbar v-model="hasInfo" right bottom :timeout="3000">
        {{ infoText }}
        <!-- <v-spacer></v-spacer> -->
        <template v-slot:action="{ attrs }">
          <v-btn icon>
            <v-icon v-bind="attrs" @click="hasInfo = false">close</v-icon>
          </v-btn>
        </template>
      </v-snackbar>
      <div
        class="progress-bar"
        :style="`width:${(renewCounter * 100) / CountMax}%;`"
      ></div>
    </v-container>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import gql from "graphql-tag";
import UserAvatar from "@/components/UserAvatar.vue";
import { RankDesc, RankResult } from "@/struct";
import constValue from "../constValue";

const UserPerPage = 10;

@Component({
  components: {
    UserAvatar
  }
})
export default class Rank extends Vue {
  private monitorMode = false;
  private renewCounter = 0;
  private CountMax = constValue.CountMax;
  private sseSource: EventSource | undefined;

  // private rankColor = ["amber", "light-blue", "green"];
  private page: number = 1;

  private ranks: RankDesc[] = [];
  private pageCount: number = 1;

  private infoText: string = "";
  private hasInfo: boolean = false;

  private series = [
    {
      name: "High - 2013",
      data: [28, 29, 33, 36, 32, 32, 33]
    },
    {
      name: "Low - 2013",
      data: [12, 11, 14, 18, 17, 13, 13]
    }
  ];

  get chartOptions() {
    let isDark = this.$vuetify.theme.dark;
    return {
      chart: {
        type: "line",
        zoom: {
          enabled: false
        }
      },
      theme: {
        mode: isDark ? "dark" : "light"
      },
      dataLabels: {
        enabled: false
      },
      stroke: {
        width: 2,
        curve: "stepline",
        dashArray: 0
      },
      // legend: {
      //   tooltipHoverFormatter: function(val: any, opts: any) {
      //     return val + ' - ' + opts.w.globals.series[opts.seriesIndex][opts.dataPointIndex] + ''
      //   }
      // },
      markers: {
        size: 0,
        hover: {
          sizeOffset: 6
        }
      },
      grid: {
        borderColor: "#7f7f7f"
      }
    };
  }

  // private get topRank() {
  //   return this.ranks.slice(0, 3);
  // }
  private get pageBase() {
    return (this.page - 1) * 10; // + 3;
  }
  private get pageRank() {
    return this.ranks.slice(this.pageBase, this.pageBase + UserPerPage);
  }

  async mounted() {
    this.monitorMode = this.$route.path.split("/")[1] == "monitor";
    this.page = Math.max(parseInt(this.$route.params.page), 1);
    await this.loadData();
    if (this.monitorMode) {
      this.sseSource = new EventSource("/sse?stream=message");
      this.sseSource.addEventListener("message", async () => {
        await this.loadData();
      });
      setInterval(this.renew, 500);
    }
  }

  async renew() {
    this.renewCounter = (this.renewCounter % this.CountMax) + 1;
    if (this.renewCounter == this.CountMax) {
      await this.loadData();
    }
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
      this.ranks = res.data!.rank.rankResultDescs.sort(
        (a, b) => parseInt(b.score) - parseInt(a.score)
      );
      this.pageCount = Math.floor(
        (this.ranks.length /*- 3*/ + UserPerPage - 1) / UserPerPage
      );
    } catch (e) {
      this.infoText = e.toString();
      this.hasInfo = true;
    }
  }
}
</script>

<style lang="scss" scoped>
.chart-container {
  width: 100%;
  margin-top: 16px;
  margin-bottom: 16px;
  background-color: rgb(127 127 127 / 0.1);
}

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
