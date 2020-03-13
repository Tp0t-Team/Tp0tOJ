<template>
  <v-app id="app">
    <v-navigation-drawer v-model="drawer" app clipped>
      <nav-list :isLogin="!!$store.state.global.userId"></nav-list>
    </v-navigation-drawer>

    <v-app-bar app clipped-left class="higher pr-2">
      <v-app-bar-nav-icon @click.stop="drawer = !drawer"></v-app-bar-nav-icon>
      <v-toolbar-title>Tp0t OJ</v-toolbar-title>
      <v-spacer></v-spacer>
      <writeup-upload
        v-if="$store.state.global.userId != null"
      ></writeup-upload>
      <v-btn icon @click="WarmUp" v-if="$store.state.global.role === 'admin'">
        <v-icon color="primary">whatshot</v-icon>
      </v-btn>
    </v-app-bar>

    <v-content>
      <router-view :key="$route.query.time" />
      <v-snackbar v-model="hasInfo" top :timeout="3000">
        {{ infoTitle }}
        <br />
        {{ infoText }}
        <v-spacer></v-spacer>
        <v-btn icon>
          <v-icon @click="hasInfo = false">close</v-icon>
        </v-btn>
      </v-snackbar>
    </v-content>

    <v-footer app padless class="higher">
      <v-col class="text-center pa-1">
        <span> <strong>Tp0t</strong> &copy; 2019 </span>
      </v-col>
    </v-footer>
  </v-app>
</template>

<script lang="ts">
import { Component, Vue, Prop, Watch } from "vue-property-decorator";
import gql from "graphql-tag";
import NavList from "@/components/NavList.vue";
import WriteupUpload from "@/components/WriteupUpload.vue";
import { CompetitionResult, BulletinItem } from "./struct";

@Component({
  components: {
    NavList,
    WriteupUpload
  }
})
export default class App extends Vue {
  private drawer: boolean | null = null;
  private permissioned: boolean = false;

  private infoTitle: string = "";
  private infoText: string = "";
  private hasInfo: boolean = false;

  async mounted() {
    let userId = sessionStorage.getItem("user_id") || null;
    let role = sessionStorage.getItem("role") || null;
    if (!userId || !role) {
      userId = null;
      role = null;
    }
    this.$store.commit("global/setUserIdAndRole", { userId, role });
    try {
      let res = await this.$apollo.query<CompetitionResult>({
        query: gql`
          query {
            competition {
              message
              competition
              registerAllow
              beginTime
              endTime
            }
          }
        `
      });
      if (res.errors) throw res.errors.map(v => v.message).join(",");
      if (res.data!.competition.message) throw res.data!.competition.message;
      this.$store.commit(
        "competition/setCompetitionInfo",
        res.data!.competition
      );
    } catch (e) {
      console.log(e);
    }
    this.initSocket();
    let permission = await Notification.requestPermission();
    this.permissioned = permission === "granted";
  }

  beforeDestroy() {
    this.destroySocket();
  }

  initSocket() {
    this.$socket.on("bulletin", (items: BulletinItem[]) => {
      this.$store.commit("Bulletin/addBulletins", items);
    });
  }

  destroySocket() {
    this.$socket.off("bulletin");
    this.$socket.disconnect();
  }

  async WarmUp() {
    try {
      let res = await this.$apollo.mutate<Boolean>({
        mutation: gql`
          mutation {
            warmUp
          }
        `
      });
      if (res.errors) throw res.errors.map(v => v.message).join(",");
      if (!res.data!) throw "error";
    } catch (e) {
      console.log(e);
    }
  }

  @Watch("$store.state.bulletin.bulletins", { deep: true })
  showBulletin() {
    if ((this.$store.state.bulletins as BulletinItem[]).length == 0) {
      return;
    }
    if (this.permissioned) {
      for (let item of this.$store.state.bulletins as BulletinItem[]) {
        new Notification(item.title, { body: item.content });
      }
    } else {
      let item: BulletinItem = this.$store.state.bulletins[
        this.$store.state.bulletins.length - 1
      ];
      this.infoTitle = item.title;
      this.infoText = item.content;
      this.hasInfo = true;
    }
  }
}
</script>

<style lang="scss">
html {
  overflow: hidden;
}

.higher {
  z-index: 3;
}
</style>
