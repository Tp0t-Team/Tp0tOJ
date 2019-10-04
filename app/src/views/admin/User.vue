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
    <v-snackbar v-model="hasInfo" right bottom :timeout="3000">
      {{ infoText }}
      <v-spacer></v-spacer>
      <v-btn icon>
        <v-icon @click="hasInfo = false">close</v-icon>
      </v-btn>
    </v-snackbar>
  </v-container>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import gql from "graphql-tag";
import UserEditor from "@/components/UserEditor.vue";
import {
  UserInfo,
  SubmitInfo,
  UserInfoUpdateInput,
  UserInfoUpdateResult,
  SubmitHistoryResult
} from "@/struct";

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

  private resolves: SubmitInfo[] = [];

  private infoText: string = "";
  private hasInfo: boolean = false;

  async mounted() {
    this.users = [
      {
        userId: "1",
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
        score: "1000",
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

  async refresh() {
    try {
      let res = await this.$apollo.query<
        SubmitHistoryResult,
        { userId: string }
      >({
        query: gql`
          query($userId: string) {
            submitHistory(userId: $userId) {
              message
              submitInfos {
                submitTime
                challengeName
                mark
              }
            }
          }
        `
      });
      if (res.errors) throw res.errors.map(v => v.message).join(",");
      if (res.data!.submitHistory.message)
        throw res.data!.submitHistory.message;
      this.resolves = res.data!.submitHistory.submitInfos;
    } catch (e) {
      this.infoText = e.toString();
      this.hasInfo = true;
    }
  }

  async submit(info: UserInfo) {
    this.loading = true;
    try {
      let res = await this.$apollo.mutate<
        UserInfoUpdateResult,
        UserInfoUpdateInput
      >({
        mutation: gql`
          mutation($input: UserInfoUpdateInput!) {
            userInfoUpdate(input: $input) {
              message
            }
          }
        `,
        variables: {
          input: {
            userId: info.userId,
            name: info.name,
            role: info.role,
            department: info.department,
            grade: info.grade,
            protectedTime: info.protectedTime,
            qq: info.qq,
            mail: info.mail,
            state: info.state
          }
        }
      });
      if (res.errors) throw res.errors.map(v => v.message).join(",");
      if (res.data!.userInfoUpdate.message)
        throw res.data!.userInfoUpdate.message;
      this.loading = false;
    } catch (e) {
      this.loading = false;
      this.infoText = e.toString();
      this.hasInfo = true;
    }
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