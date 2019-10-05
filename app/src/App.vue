<template>
  <v-app id="app">
    <v-navigation-drawer v-model="drawer" app clipped>
      <nav-list :isLogin="!!$store.state.global.userId"></nav-list>
    </v-navigation-drawer>

    <v-app-bar app clipped-left>
      <v-app-bar-nav-icon @click.stop="drawer = !drawer"></v-app-bar-nav-icon>
      <v-toolbar-title>Tp0t OJ</v-toolbar-title>
    </v-app-bar>

    <v-content>
      <router-view :key="$route.query.time" />
    </v-content>

    <v-footer app padless>
      <v-col class="text-center pa-1">
        <span>
          <strong>Tp0t</strong> &copy; 2019
        </span>
      </v-col>
    </v-footer>
  </v-app>
</template>

<script lang="ts">
import { Component, Vue, Prop, Watch } from "vue-property-decorator";
import NavList from "@/components/NavList.vue";

@Component({
  components: {
    NavList
  }
})
export default class App extends Vue {
  private drawer: boolean | null = null;

  mounted() {
    let userId = sessionStorage.getItem("user_id") || null;
    let role = sessionStorage.getItem("role") || null;
    if (!userId || !role) {
      userId = null;
      role = null;
    }
    this.$store.commit("global/setUserIdAndRole", { userId, role });
  }
}
</script>

<style lang="scss">
html {
  overflow: hidden;
}
</style>
