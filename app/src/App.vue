<template>
  <v-app id="inspire">
    <v-navigation-drawer v-model="drawer" app clipped>
      <nav-list></nav-list>
    </v-navigation-drawer>

    <v-app-bar app clipped-left>
      <v-app-bar-nav-icon @click.stop="drawer = !drawer"></v-app-bar-nav-icon>
      <v-toolbar-title>Tp0t OJ</v-toolbar-title>
    </v-app-bar>

    <v-content>
      <router-view />
    </v-content>

    <v-footer app padless>
      <v-col class="text-center pa-1">
        <span>
          <strong>Tp0t</strong> &copy; 2019
        </span>
      </v-col>
    </v-footer>
  </v-app>
  <!-- <div id="app">
    <div id="nav">
      <router-link to="/">Home</router-link> |
      <router-link to="/about">About</router-link>
    </div>
    <router-view />
  </div>-->
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
}
</script>

<style lang="scss">
// #app {
//   font-family: "Avenir", Helvetica, Arial, sans-serif;
//   -webkit-font-smoothing: antialiased;
//   -moz-osx-font-smoothing: grayscale;
//   text-align: center;
//   color: #2c3e50;
// }
// #nav {
//   padding: 30px;
//   a {
//     font-weight: bold;
//     color: #2c3e50;
//     &.router-link-exact-active {
//       color: #42b983;
//     }
//   }
// }
</style>
