<template>
  <v-app id="app">
    <v-navigation-drawer v-model="drawer" app clipped>
      <nav-list :isLogin="!!$store.state.global.userId"></nav-list>
    </v-navigation-drawer>

    <v-app-bar app clipped-left class="higher pr-2">
      <v-app-bar-nav-icon @click.stop="drawer = !drawer"></v-app-bar-nav-icon>
      <v-toolbar-title>{{ title }}</v-toolbar-title>
      <v-spacer></v-spacer>
      <!-- <v-btn icon @click="WarmUp" v-if="$store.state.global.role === 'admin'">
        <v-icon color="primary">whatshot</v-icon>
      </v-btn> -->
    </v-app-bar>

    <v-content>
      <router-view :key="$route.query.time" />
    </v-content>
    <v-snackbar
      v-model="hasSSEInfo"
      right
      top
      multi-line
      color="success"
      :timeout="-1"
      class="sse-info"
    >
      <strong>{{ sseInfoTitle }}</strong
      ><br />
      <span>{{ sseInfoText }}</span>
      <!-- <v-spacer></v-spacer> -->
      <template v-slot:action="{ attrs }">
        <v-btn icon>
          <v-icon v-bind="attrs" @click="hasSSEInfo = false">close</v-icon>
        </v-btn>
      </template>
    </v-snackbar>

    <v-footer app padless class="higher">
      <v-col class="text-center pa-1">
        <span>
          Powered by
          <a href="https://github.com/Tp0t-Team/Tp0tOJ" target="_blank"
            ><strong>Tp0t OJ</strong></a
          >, under AGPL license
        </span>
      </v-col>
    </v-footer>
  </v-app>
</template>

<script lang="ts">
import { Component, Vue, Watch } from "vue-property-decorator";
import gql from "graphql-tag";
import NavList from "@/components/NavList.vue";
import constValue from "./constValue";

@Component({
  components: {
    NavList
  }
})
export default class App extends Vue {
  private drawer: boolean | null = null;
  private title = constValue.navTitle;

  private sseInfoText: string = "";
  private sseInfoTitle: string = "";
  private hasSSEInfo: boolean = false;
  private sseSource: EventSource = new EventSource("/sse?stream=message");
  private timer: number | null = null;

  @Watch("$route.path")
  async isMonitor() {
    if (this.$route.path.split("/")[1] == "monitor") {
      this.drawer = false;
    }
  }

  mounted() {
    document.title = this.title;
    let userId = localStorage.getItem("user_id") || null;
    let role = localStorage.getItem("role") || null;
    if (!userId || !role) {
      userId = null;
      role = null;
    }
    this.$store.commit("global/setUserIdAndRole", { userId, role });
    this.sseSource.addEventListener("message", e => {
      let parsed = JSON.parse(e.data);
      this.sseInfoTitle = parsed.title;
      this.sseInfoText = parsed.info;
      this.hasSSEInfo = true;
      if (this.timer !== null) {
        clearTimeout(this.timer);
        this.timer = null;
      }
      this.timer = setTimeout(() => {
        this.hasSSEInfo = false;
      }, 10 * 1000) as any;
    });
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
      if (e === "unauthorized") {
        this.$store.commit("global/resetUserIdAndRole");
        this.$router.push("/login?");
      }
      console.log(e);
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

.sse-info * {
  text-overflow: ellipsis;
  overflow: hidden;
}
</style>
