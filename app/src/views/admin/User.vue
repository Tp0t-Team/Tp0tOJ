<template>
  <v-container fluid>
    <v-row>
      <v-col cols="6" class="content-col">
        <v-data-table
          :headers="headers"
          :items="users"
          @click:row="select"
        ></v-data-table>
      </v-col>
      <v-col cols="6" class="content-col">
        <v-divider class="avatar-divider" color="primary"></v-divider>
        <v-row class="avatar">
          <v-spacer></v-spacer>
          <v-btn fab @click="seeProfile()" :disabled="!currentUser">
            <v-avatar color="blue" size="64">
              <!-- <span>{{currentUser && currentUser.name[0] || "?"}}</span> -->
              <user-avatar
                class="white--text"
                v-if="!!currentUser"
                :url="currentUser.avatar"
                :size="64"
                :name="currentUser.name"
                :key="currentUser.avatar"
              ></user-avatar>
            </v-avatar>
          </v-btn>
          <v-spacer></v-spacer>
        </v-row>
        <user-editor
          :disabled="withoutInit"
          :loading="loading"
          :key="(currentUser && currentUser.mail) || ''"
          :user="currentUser"
          @submit="submit"
        ></user-editor>
        <!-- <v-toolbar dense>
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
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in resolves" :key="item.submitTime">
              <td>{{item.submitTime}}</td>
              <td>{{item.challengeName}}</td>
            </tr>
          </tbody>
        </v-simple-table> -->
      </v-col>
    </v-row>
    <v-snackbar v-model="hasInfo" right bottom :timeout="3000">
      {{ infoText }}
      <!-- <v-spacer></v-spacer> -->
      <template v-slot:action="{ attrs }">
        <v-btn icon>
          <v-icon v-bind="attrs" @click="hasInfo = false">close</v-icon>
        </v-btn>
      </template>
    </v-snackbar>
  </v-container>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import gql from "graphql-tag";
import UserEditor from "@/components/UserEditor.vue";
import UserAvatar from "@/components/UserAvatar.vue";
import {
  UserInfo,
  AllUserInfoResult,
  SubmitInfo,
  UserInfoUpdateInput,
  UserInfoUpdateResult
} from "@/struct";

@Component({
  components: {
    UserEditor,
    UserAvatar
  }
})
export default class User extends Vue {
  private headers = [
    { text: "name", value: "name" },
    { text: "role", value: "role" },
    { text: "state", value: "state" },
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
    await this.loadAll();
  }

  async loadAll() {
    try {
      let res = await this.$apollo.query<AllUserInfoResult, {}>({
        query: gql`
          query {
            allUserInfos {
              message
              allUserInfos {
                userId
                name
                avatar
                role
                mail
                joinTime
                score
                state
              }
            }
          }
        `,
        fetchPolicy: "no-cache"
      });
      if (res.errors) throw res.errors.map(v => v.message).join(",");
      if (res.data!.allUserInfos.message) throw res.data!.allUserInfos.message;
      this.users = res.data!.allUserInfos.allUserInfos;
    } catch (e) {
      this.infoText = e.toString();
      this.hasInfo = true;
    }
  }

  select(user: UserInfo) {
    this.withoutInit = false;
    this.currentUser = user;
    this.resolves = [];
  }

  seeProfile() {
    if (this.currentUser)
      this.$router.push(`/profile/${this.currentUser.userId}`);
  }

  // async refresh() {
  //   try {
  //     let res = await this.$apollo.query<
  //       SubmitHistoryResult,
  //       { input: string }
  //     >({
  //       query: gql`
  //         query($input: String!) {
  //           submitHistory(userId: $input) {
  //             message
  //             submitInfos {
  //               submitTime
  //               challengeName
  //             }
  //           }
  //         }
  //       `,
  //       variables: {
  //         input: this.currentUser!.userId
  //       },
  //       fetchPolicy: "no-cache"
  //     });
  //     if (res.errors) throw res.errors.map(v => v.message).join(",");
  //     if (res.data!.submitHistory.message)
  //       throw res.data!.submitHistory.message;
  //     this.resolves = res.data!.submitHistory.submitInfos;
  //   } catch (e) {
  //     this.infoText = e.toString();
  //     this.hasInfo = true;
  //   }
  // }

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
            mail: info.mail,
            state: info.state
          }
        }
      });
      if (res.errors) throw res.errors.map(v => v.message).join(",");
      if (res.data!.userInfoUpdate.message)
        throw res.data!.userInfoUpdate.message;
      this.loading = false;
      await this.loadAll();
      this.infoText = "change user info success";
      this.hasInfo = true;
    } catch (e) {
      this.loading = false;
      this.infoText = e.toString();
      this.hasInfo = true;
    }
  }
}
</script>

<style lang="scss" scoped>
.avatar {
  padding: 12px;
}

.content-col {
  height: calc(100vh - 120px);
  overflow-y: auto;
}

.avatar-divider {
  position: relative;
  bottom: -32px;
}

// .mark-badge {
//   position: relative;
//   top: -4px;
// }
</style>
