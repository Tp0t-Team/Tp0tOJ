<template>
  <div class="nav-wrapper">
    <v-list dense>
      <v-list-item @click="$router.push('/')">
        <v-list-item-action>
          <v-icon>home</v-icon>
        </v-list-item-action>
        <v-list-item-content>
          <v-list-item-title>Home</v-list-item-title>
        </v-list-item-content>
      </v-list-item>
      <v-list-item
        @click="
          $router.push({
            path: '/bulletin',
            query: { time: Date.now().toLocaleString() }
          })
        "
      >
        <v-list-item-action>
          <v-icon>notifications</v-icon>
        </v-list-item-action>
        <v-list-item-content>
          <v-list-item-title>Announcement</v-list-item-title>
        </v-list-item-content>
      </v-list-item>
      <v-list-item
        v-if="isLogin"
        @click="
          $router.push({
            path: '/challenge',
            query: { time: Date.now().toLocaleString() }
          })
        "
      >
        <v-list-item-action>
          <v-icon>list_alt</v-icon>
        </v-list-item-action>
        <v-list-item-content>
          <v-list-item-title>Challenges</v-list-item-title>
        </v-list-item-content>
      </v-list-item>
      <v-list-item
        @click="
          $router.push({
            path: '/rank/1',
            query: { time: Date.now().toLocaleString() }
          })
        "
      >
        <v-list-item-action>
          <v-icon>assessment</v-icon>
        </v-list-item-action>
        <v-list-item-content>
          <v-list-item-title>Rank</v-list-item-title>
        </v-list-item-content>
      </v-list-item>
      <v-list-item
        v-if="isLogin"
        @click="$router.push(`/profile/${$store.state.global.userId}`)"
      >
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
      <div v-if="$store.state.global.role == 'admin'">
        <v-divider></v-divider>
        <v-list-group prepend-icon="build" no-action>
          <template v-slot:activator>
            <v-list-item-title>Admin</v-list-item-title>
          </template>
          <v-list-item @click="$router.push('/admin/user')">
            <v-list-item-content>
              <v-list-item-title>Users</v-list-item-title>
            </v-list-item-content>
          </v-list-item>
          <v-list-item
            @click="
              $router.push({
                path: '/admin/challenge',
                query: { time: Date.now().toLocaleString() }
              })
            "
          >
            <v-list-item-content>
              <v-list-item-title>Challenges</v-list-item-title>
            </v-list-item-content>
          </v-list-item>
          <v-list-item
            @click="
              $router.push({
                path: '/admin/images',
                query: { time: Date.now().toLocaleString() }
              })
            "
          >
            <v-list-item-content>
              <v-list-item-title>Images</v-list-item-title>
            </v-list-item-content>
          </v-list-item>
          <v-list-item
            @click="
              $router.push({
                path: '/admin/cluster',
                query: { time: Date.now().toLocaleString() }
              })
            "
          >
            <v-list-item-content>
              <v-list-item-title>Cluster</v-list-item-title>
            </v-list-item-content>
          </v-list-item>
          <v-list-item
            @click="
              $router.push({
                path: '/admin/writeup',
                query: { time: Date.now().toLocaleString() }
              })
            "
          >
            <v-list-item-content>
              <v-list-item-title>Writeup</v-list-item-title>
            </v-list-item-content>
          </v-list-item>
          <v-list-item
            @click="
              $router.push({
                path: '/admin/event',
                query: { time: Date.now().toLocaleString() }
              })
            "
          >
            <v-list-item-content>
              <v-list-item-title>Event</v-list-item-title>
            </v-list-item-content>
          </v-list-item>
          <v-list-item
            @click="
              $router.push({
                path: '/admin/analyse',
                query: { time: Date.now().toLocaleString() }
              })
            "
          >
            <v-list-item-content>
              <v-list-item-title>Analyse</v-list-item-title>
            </v-list-item-content>
          </v-list-item>
        </v-list-group>
        <v-divider></v-divider>
      </div>
      <v-list-item-group multiple v-model="settings">
        <v-list-item value="dark">
          <template v-slot:default="{ active }">
            <v-list-item-action>
              <v-checkbox :value="active"></v-checkbox>
            </v-list-item-action>
            <v-list-item-content>
              <v-list-item-title>Dark theme</v-list-item-title>
            </v-list-item-content>
          </template>
        </v-list-item>
      </v-list-item-group>
    </v-list>
    <v-tooltip top v-if="$store.state.global.role == 'member'">
      <template v-slot:activator="{ on, attrs }">
        <v-btn
          tile
          v-bind="attrs"
          v-on="on"
          @click="writeupClick"
          color="primary"
        >
          <v-icon>upload_file</v-icon>
        </v-btn>
      </template>
      <span>Upload Writeup</span>
      <input
        type="file"
        class="file-input"
        ref="writeup"
        @change="uploadWriteup"
      />
    </v-tooltip>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Watch, Prop } from "vue-property-decorator";
import gql from "graphql-tag";
import { LogoutResult } from "@/struct";

@Component
export default class NavList extends Vue {
  @Prop() isLogin!: boolean;
  private settings: string[] = [];

  mounted() {
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
      this.$store.commit("global/resetUserIdAndRole");
      this.$router.push("/");
      // must logout success
      if (res.errors) throw res.errors.map(v => v.message).join(",");
      if (res.data!.logout.message) throw res.data!.logout.message;
    } catch (e) {
      console.log(e.toString());
    }
  }

  writeupClick() {
    (this.$refs.writeup! as any).click();
  }

  async uploadWriteup(e: any) {
    console.log(e);
    let file: File = e.target.files[0];
    let formData = new FormData();
    formData.append("writeup", file, file.name);
    try {
      let res = await fetch("/writeup", {
        method: "POST",
        headers: {
          "X-CSRF-Token": (globalThis as any).CsrfToken as string
        },
        body: formData
      });
      if (res.status == 401) {
        this.$store.commit("global/resetUserIdAndRole");
        this.$router.push("/login?unauthorized");
        return;
      }
      if (res.status != 200) {
        throw res.statusText;
      }
    } catch (err) {
      alert(err); // TODO:
      e.target.value = "";
      return;
    }
    alert("upload writeup success");
    e.target.value = "";
  }
}
</script>

<style lang="scss" scoped>
.nav-wrapper {
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.file-input {
  display: none;
}
</style>
