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
      <rank-table :value="ranks" :userPerPage="10" :pageInit="page" />
      <v-btn fab absolute right bottom color="primary" @click="pause = !pause">
        <v-icon>{{ pause ? "play_arrow" : "pause" }}</v-icon>
      </v-btn>
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
import RankTable from "@/components/RankTable.vue";
import { RankWithIndex, RankResult, ChartData } from "@/struct";
import constValue from "../constValue";

const userLocale =
  navigator.languages && navigator.languages.length
    ? navigator.languages[0]
    : navigator.language;

@Component({
  components: {
    UserAvatar,
    RankTable
  }
})
export default class Monitor extends Vue {
  private monitorMode = false;
  private renewCounter = 0;
  private CountMax = constValue.CountMax;
  private sseSource: EventSource | undefined;

  private page: number = 1; // only for init

  private ranks: RankWithIndex[] = [];

  private pause: boolean = false;

  private infoText: string = "";
  private hasInfo: boolean = false;

  private series: { name: string; data: number[][] }[] = [];

  get chartOptions() {
    let isDark = this.$vuetify.theme.dark;
    return {
      chart: {
        type: "line",
        zoom: {
          enabled: false
        }
      },
      colors: [
        "#1f77b4", // muted blue
        "#ff7f0e", // safety orange
        "#2ca02c", // cooked asparagus green
        "#d62728", // brick red
        "#9467bd", // muted purple
        "#8c564b", // chestnut brown
        "#e377c2", // raspberry yogurt pink
        "#7f7f7f", // middle gray
        "#bcbd22", // curry yellow-green
        "#17becf" // blue-teal
      ],
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
      tooltip: {
        x: {
          formatter: function(val: number) {
            return Intl.DateTimeFormat(userLocale, {
              dateStyle: "short",
              timeStyle: "medium"
            } as Intl.DateTimeFormatOptions).format(val);
          }
        }
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
      },
      xaxis: {
        type: "datetime",
        labels: {
          datetimeUTC: false
        }
      }
    };
  }

  async mounted() {
    this.monitorMode = this.$route.path.split("/")[1] == "monitor";
    this.page = Math.max(parseInt(this.$route.params.page), 1);
    await this.loadData();
    if (this.monitorMode) {
      this.sseSource = new EventSource("/sse?stream=message");
      this.sseSource.addEventListener("message", async () => {
        if (this.pause) {
          this.renewCounter = this.CountMax - 1;
          return;
        }
        await this.loadData();
      });
      setInterval(this.renew, 500);
    }
  }

  async renew() {
    if (this.pause) {
      return;
    }
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
      this.ranks = res
        .data!.rank.rankResultDescs.sort(
          (a, b) => parseInt(b.score) - parseInt(a.score)
        )
        .map((it, index) => ({ index, desc: it } as RankWithIndex));

      let chartRes = await fetch(
        `/chart/?num=${this.$route.query["num"] ?? "10"}`,
        {
          cache: "no-cache",
          headers: {
            "X-CSRF-Token": (globalThis as any).CsrfToken as string
          }
        }
      );
      let originChartData: ChartData = await chartRes.json();
      let seriesList = [];
      let now = Math.max(
        Date.now(),
        originChartData.x.length > 0
          ? originChartData.x[originChartData.x.length - 1]
          : 0
      );
      for (let series of originChartData.y) {
        let seriesItem = {
          name: series.name,
          data: [] as number[][]
        };
        for (let index = 0; index < series.score.length; index++) {
          seriesItem.data.push([originChartData.x[index], series.score[index]]);
        }
        if (seriesItem.data.length > 0) {
          seriesItem.data.push([
            now,
            seriesItem.data[seriesItem.data.length - 1][1]
          ]);
        } else {
          seriesItem.data.push([now, 0]);
        }
        seriesList.push(seriesItem);
      }
      this.series = seriesList;
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
  margin: 16px 24px;
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
