<template>
  <v-container fluid>
    <v-row>
      <v-col cols="6" class="content-col">
        <v-data-table :headers="headers" :items="users" @click:row="select"></v-data-table>
      </v-col>
      <v-col cols="6" class="content-col">
        <v-divider class="avatar-divider" color="primary"></v-divider>
        <v-row>
          <v-spacer></v-spacer>
          <v-avatar color="blue" size="64">
            <span>{{currentUser && currentUser.name[0] || "?"}}</span>
          </v-avatar>
          <v-spacer></v-spacer>
        </v-row>
        <user-editor
          :disabled="withoutInit"
          :loading="loading"
          :key="currentUser && currentUser.stuNumber || ''"
          :user="currentUser"
          @submit="submit"
        ></user-editor>
        <v-toolbar dense>
          解题列表
          <v-spacer></v-spacer>
          <v-btn icon @click="refresh" :disabled="withoutInit">
            <v-icon>refresh</v-icon>
          </v-btn>
        </v-toolbar>
        <v-simple-table dense>
          <thead>
            <tr>
              <th>submit time</th>
              <th>challenge name</th>
              <th>mark</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in resolves" :key="item.submitTime">
              <td>{{item.submitTime}}</td>
              <td>{{item.challengeName}}</td>
              <td>
                <v-badge v-if="!!item.mark" class="mark-badge" :color="rankColor[item.mark - 1]">
                  <template v-slot:badge>{{item.mark}}</template>
                </v-badge>
              </td>
            </tr>
          </tbody>
        </v-simple-table>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import UserEditor from "@/components/UserEditor.vue";
import { UserInfo, ResolveInfo } from "@/struct";

@Component({
  components: {
    UserEditor
  }
})
export default class User extends Vue {
  private headers = [
    { text: "name", value: "name" },
    { text: "role", value: "role" },
    { text: "state", value: "state" },
    { text: "grade", value: "grade" },
    { text: "protected_time", value: "protectedTime" },
    { text: "score", value: "score" }
  ];
  private rankColor = ["amber", "light-blue", "green"];

  private withoutInit: boolean = true;
  private loading: boolean = false;

  private users: UserInfo[] = [];
  private currentUser: UserInfo | null = null;

  private resolves: ResolveInfo[] = [];

  async mounted() {
    this.users = [
      {
        name: "DRSN",
        role: "team",
        stuNumber: "9161",
        department: "计算机学院",
        grade: "2016",
        protectedTime: "2019-01-01",
        qq: "784227594",
        mail: "test@test.com",
        topRank: "1",
        joinTime: "2018-01-01",
        score: 1000,
        state: "normal",
        rank: 1
      }
    ];
  }

  select(user: UserInfo) {
    this.withoutInit = false;
    this.currentUser = user;
    this.resolves = [];
  }

  refresh() {
    console.log(this.currentUser!.stuNumber);
  }

  submit(info: UserInfo) {
    this.loading = true;
    console.log(info);
  }
}
</script>

<style lang="scss" scoped>
.content-col {
  height: calc(100vh - 120px);
  overflow-y: auto;
}

.avatar-divider {
  position: relative;
  bottom: -32px;
}

.mark-badge {
  position: relative;
  top: -4px;
}
</style>