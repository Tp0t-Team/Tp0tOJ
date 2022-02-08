<template>
  <v-app id="app">
    <v-navigation-drawer v-model="drawer" app clipped>
      <nav-list :isLogin="!!$store.state.global.userId"></nav-list>
    </v-navigation-drawer>

    <v-app-bar app clipped-left class="higher pr-2">
      <v-app-bar-nav-icon @click.stop="drawer = !drawer"></v-app-bar-nav-icon>
      <v-toolbar-title>Tp0t OJ</v-toolbar-title>
      <v-spacer></v-spacer>
      <v-btn
        icon
        @click="WarmUp"
        v-if="$store.state.global.role==='admin'"
      >
        <v-icon color="primary">whatshot</v-icon>
      </v-btn>
    </v-app-bar>

    <v-content>
      <router-view :key="$route.query.time" />
    </v-content>

    <v-footer app padless class="higher">
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
import gql from 'graphql-tag';
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

  async WarmUp() {
    try{
      let res = await this.$apollo.mutate<Boolean>({
        mutation: gql`
          mutation{
            warmUp
          }`
      });
      if (res.errors) throw res.errors.map(v => v.message).join(",");
      if (!res.data!) throw "error";
    }catch(e) {
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
</style>
