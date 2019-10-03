<template>
  <v-container>
    <!-- {{$route.params.user_id}} -->
    <v-row class="avatar">
      <v-spacer></v-spacer>
      <v-avatar :color="loading?'grey':'blue'" size="64">
        <span class="white--text">{{userInfo.name[0]}}</span>
      </v-avatar>
      <v-spacer></v-spacer>
    </v-row>
    <v-row class="mt-5">
      <v-spacer></v-spacer>
      <v-col cols="6">
        <v-card class="outer">
          <v-card class="inner pa-4">
            <v-form>
              <v-row justify="center" class="mt-4">
                <div>
                  <v-subheader>
                    <span class="text-center">
                      <strong>
                        Rank:
                        <span v-if="userInfo.rank!=0">{{userInfo.rank}}</span>
                        <span v-else>∞</span>
                      </strong>
                      <br />
                      {{userInfo.score}}pt
                    </span>
                  </v-subheader>
                </div>
              </v-row>
              <v-row>
                <v-col cols="6">
                  <v-text-field :loading="loading" readonly label="name" :value="userInfo.name"></v-text-field>
                </v-col>
                <v-col cols="6">
                  <v-text-field :loading="loading" readonly label="role" :value="userInfo.role"></v-text-field>
                </v-col>
              </v-row>
              <v-row v-if="$store.state.global.userId==$route.params.user_id">
                <v-col cols="6">
                  <v-text-field
                    :loading="loading"
                    readonly
                    label="student number"
                    :value="userInfo.stuNumber"
                  ></v-text-field>
                </v-col>
                <v-col cols="6">
                  <v-text-field
                    :loading="loading"
                    readonly
                    label="department"
                    :value="userInfo.department"
                  ></v-text-field>
                </v-col>
              </v-row>
              <v-row>
                <v-col cols="6">
                  <v-text-field :loading="loading" readonly label="state" :value="userInfo.state"></v-text-field>
                </v-col>
                <v-col cols="6">
                  <v-text-field
                    :loading="loading"
                    readonly
                    label="protected time"
                    :value="userInfo.protectedTime"
                  ></v-text-field>
                </v-col>
              </v-row>
              <v-row>
                <v-col cols="6">
                  <v-text-field :loading="loading" readonly label="grade" :value="userInfo.grade"></v-text-field>
                </v-col>
                <v-col cols="6">
                  <v-text-field
                    :loading="loading"
                    readonly
                    label="join time"
                    :value="userInfo.joinTime"
                  ></v-text-field>
                </v-col>
              </v-row>
              <v-row v-if="$store.state.global.userId==$route.params.user_id">
                <v-col cols="6">
                  <v-text-field :loading="loading" readonly label="QQ" :value="userInfo.qq"></v-text-field>
                </v-col>
                <v-col cols="6">
                  <v-text-field :loading="loading" readonly label="mail" :value="userInfo.mail"></v-text-field>
                </v-col>
              </v-row>
            </v-form>
          </v-card>
        </v-card>
      </v-col>
      <v-spacer></v-spacer>
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
import { UserInfo, UserInfoResult } from "@/struct";
import gql from "graphql-tag";

@Component
export default class Profile extends Vue {
  private loading: boolean = false;
  private userInfo: UserInfo = {
    userId: "",
    name: "",
    role: "",
    stuNumber: "",
    department: "",
    grade: "",
    protectedTime: "",
    qq: "",
    mail: "",
    topRank: "",
    joinTime: "",
    score: 0,
    state: "",
    rank: 0
  };

  private infoText: string = "";
  private hasInfo: boolean = false;

  async mounted() {
    // example data
    this.loading = true;
    this.userInfo = {
      userId: "1",
      name: "CXK",
      role: "member",
      stuNumber: "0001",
      department: "xxxx",
      grade: "2010",
      protectedTime: "2019-10-10",
      qq: "12345678",
      mail: "123@test.com",
      topRank: "∞",
      joinTime: "2019-01-01",
      score: 0,
      state: "protected",
      rank: 0
    };
    //
    try {
      let res = await this.$apollo.query<UserInfoResult>({
        query: gql`
          query($userId: String!) {
            userInfo(userId: $userId) {
              message
              userInfo {
                userId
                name
                role
                stuNumber
                department
                grade
                protectedTime
                qq
                mail
                topRank
                joinTime
                score
                state
                rank
              }
            }
          }
        `,
        variables: {
          userId: this.$route.params.user_id
        }
      });
      if (res.errors) throw res.errors.map(v => v.message).join(",");
      if (res.data!.userInfo.message) throw res.data!.userInfo.message;
      this.userInfo = res.data!.userInfo.userInfo;
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
.avatar {
  position: absolute;
  left: 12px;
  width: 100%;
  z-index: 1;
}

.outer {
  background: transparent;
  box-shadow: 0 0 4px rgba(0, 0, 0, 0.4);
}

.inner {
  background: transparent;
  box-shadow: 0 0 4px rgba(0, 0, 0, 0.4) inset;
}
</style>