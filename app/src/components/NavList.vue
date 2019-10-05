<template>
  <v-list dense>
    <v-list-item @click="$router.push('/')">
      <v-list-item-action>
        <v-icon>home</v-icon>
      </v-list-item-action>
      <v-list-item-content>
        <v-list-item-title>Home</v-list-item-title>
      </v-list-item-content>
    </v-list-item>
    <v-list-item @click="$router.push({path:'/bulletin',query:{time:Date.now().toLocaleString()}})">
      <v-list-item-action>
        <v-icon>notifications</v-icon>
      </v-list-item-action>
      <v-list-item-content>
        <v-list-item-title>Announcement</v-list-item-title>
      </v-list-item-content>
    </v-list-item>
    <v-list-item
      v-if="isLogin"
      @click="$router.push({path:'/challenge',query:{time:Date.now().toLocaleString()}})"
    >
      <v-list-item-action>
        <v-icon>list_alt</v-icon>
      </v-list-item-action>
      <v-list-item-content>
        <v-list-item-title>Challenges</v-list-item-title>
      </v-list-item-content>
    </v-list-item>
    <v-list-item @click="$router.push({path:'/rank/1',query:{time:Date.now().toLocaleString()}})">
      <v-list-item-action>
        <v-icon>assessment</v-icon>
      </v-list-item-action>
      <v-list-item-content>
        <v-list-item-title>Rank</v-list-item-title>
      </v-list-item-content>
    </v-list-item>
    <v-list-item v-if="isLogin" @click="$router.push(`/profile/${$store.state.global.userId}`)">
      <v-list-item-action>
        <v-icon>person</v-icon>
      </v-list-item-action>
      <v-list-item-content>
        <v-list-item-title>Profile</v-list-item-title>
      </v-list-item-content>
    </v-list-item>
    <v-list-item v-if="isLogin" @click="logout">
      <v-list-item-action>
        <v-icon>exit_to_app</v-icon>
      </v-list-item-action>
      <v-list-item-content>
        <v-list-item-title>Logout</v-list-item-title>
      </v-list-item-content>
    </v-list-item>
    <v-list-item v-if="!isLogin" @click="$router.push(`/login`)">
      <v-list-item-action>
        <v-icon>person</v-icon>
      </v-list-item-action>
      <v-list-item-content>
        <v-list-item-title>Login | Register</v-list-item-title>
      </v-list-item-content>
    </v-list-item>
    <div v-if="$store.state.global.role=='admin' || $store.state.global.role=='team'">
      <v-divider></v-divider>
      <v-list-group prepend-icon="build" no-action>
        <template v-slot:activator>
          <v-list-item-title>Admin</v-list-item-title>
        </template>
        <v-list-item v-if="$store.state.global.role=='admin'" @click="$router.push('/admin/user')">
          <v-list-item-content>
            <v-list-item-title>Users</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-item @click="$router.push('/admin/challenge')">
          <v-list-item-content>
            <v-list-item-title>Challenges</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list-group>
      <v-divider></v-divider>
    </div>
    <v-list-item-group multiple v-model="settings">
      <v-list-item value="dark">
        <template v-slot:default="{ active, toggle }">
          <v-list-item-action>
            <v-checkbox v-model="active" @click="toggle"></v-checkbox>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title>Dark theme</v-list-item-title>
          </v-list-item-content>
        </template>
      </v-list-item>
    </v-list-item-group>
  </v-list>
</template>

<script lang="ts">
import { Component, Vue, Watch, Prop } from "vue-property-decorator";
import gql from "graphql-tag";
import { LogoutResult } from "@/struct";

@Component
export default class NavList extends Vue {
  @Prop() isLogin!: boolean;
  private settings: string[] = [];

  created() {
    if (this.$vuetify.theme.dark) {
      this.settings = ["dark"];
    }
  }

  @Watch("settings")
  settingsChanged() {
    this.$vuetify.theme.dark =
      this.settings.findIndex(val => val == "dark") >= 0;
  }

  async logout() {
    try {
      let res = await this.$apollo.mutate<LogoutResult, {}>({
        mutation: gql`
          mutation {
            logout {
              message
            }
          }
        `
      });
      if (res.errors) throw res.errors.map(v => v.message).join(",");
      if (res.data!.logout.message) throw res.data!.logout.message;
      this.$store.commit("global/resetUserIdAndRole");
      sessionStorage.removeItem("user_id");
      sessionStorage.removeItem("role");
      this.$router.push("/");
    } catch (e) {
      console.log(e.toString());
    }
  }
}
</script>

<style lang="scss" scoped>
</style>