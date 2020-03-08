<template>
  <v-container>
    <!-- {{$route.params.user_id}} -->
    <v-row class="avatar">
      <v-spacer></v-spacer>
      <v-tooltip right>
        <template v-slot:activator="{ on }">
          <v-avatar
            class="hoverable"
            :color="loading ? 'grey' : 'blue'"
            size="64"
            @click="editAvatar"
            v-on="on"
          >
            <!-- <span class="white--text">{{userInfo.name[0]}}</span> -->
            <user-avatar
              class="white--text"
              :url="userInfo.avatar"
              :size="64"
              :name="userInfo.name"
              :key="userInfo.avatar"
            ></user-avatar>
          </v-avatar>
        </template>
        <span>avatar service based on Gravatar</span>
      </v-tooltip>
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
                      <strong v-if="$store.state.competition.competition">{{ userInfo.name }}</strong>
                      <br v-if="$store.state.competition.competition" />
                      <strong>
                        Rank:
                        <span v-if="userInfo.rank != 0">{{
                          userInfo.rank
                        }}</span>
                        <span v-else>∞</span>
                      </strong>
                      <br />
                      {{ userInfo.score }}pt
                    </span>
                  </v-subheader>
                </div>
              </v-row>
              <v-row v-if="!$store.state.competition.competition">
                <v-col cols="6">
                  <v-text-field
                    :loading="loading"
                    readonly
                    label="name"
                    :value="userInfo.name"
                  ></v-text-field>
                </v-col>
                <v-col cols="6">
                  <v-text-field
                    :loading="loading"
                    readonly
                    label="role"
                    :value="userInfo.role"
                  ></v-text-field>
                </v-col>
              </v-row>
              <v-row
                v-if="
                  $store.state.global.userId == $route.params.user_id &&
                    !$store.state.competition.competition
                "
              >
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
              <v-row v-if="!$store.state.competition.competition">
                <v-col cols="6">
                  <v-text-field
                    :loading="loading"
                    readonly
                    label="state"
                    :value="userInfo.state"
                  ></v-text-field>
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
              <v-row v-if="!$store.state.competition.competition">
                <v-col cols="6">
                  <v-text-field
                    :loading="loading"
                    readonly
                    label="grade"
                    :value="userInfo.grade"
                  ></v-text-field>
                </v-col>
                <v-col cols="6">
                  <v-text-field
                    :loading="loading"
                    readonly
                    label="join time"
                    :value="showJoinTime"
                  ></v-text-field>
                </v-col>
              </v-row>
              <v-row
                v-if="
                  $store.state.global.userId == $route.params.user_id &&
                    !$store.state.competition.competition
                "
              >
                <v-col cols="6">
                  <v-text-field
                    :loading="loading"
                    readonly
                    label="QQ"
                    :value="userInfo.qq"
                  ></v-text-field>
                </v-col>
                <v-col cols="6">
                  <v-text-field
                    :loading="loading"
                    readonly
                    label="mail"
                    :value="userInfo.mail"
                  ></v-text-field>
                </v-col>
              </v-row>
            </v-form>
            <v-row v-if="$store.state.competition.competition">
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
                    <td>{{ item.submitTime }}</td>
                    <td>{{ item.challengeName }}</td>
                    <td>
                      <v-badge
                        v-if="!!item.mark"
                        class="mark-badge"
                        :color="rankColor[item.mark - 1]"
                      >
                        <template v-slot:badge>{{ item.mark }}</template>
                      </v-badge>
                    </td>
                  </tr>
                </tbody>
              </v-simple-table>
            </v-row>
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
import { Component, Vue, Watch } from "vue-property-decorator";
import { Route } from "vue-router";
import UserAvatar from "@/components/UserAvatar.vue";
import {
  UserInfo,
  UserInfoResult,
  SubmitInfo,
  SubmitHistoryResult
} from "@/struct";
import gql from "graphql-tag";

@Component({
  components: {
    UserAvatar
  }
})
export default class Profile extends Vue {
  private loading: boolean = false;
  private userInfo: UserInfo = {
    userId: "",
    name: "",
    avatar: "",
    role: "",
    stuNumber: "",
    department: "",
    grade: "",
    protectedTime: "",
    qq: "",
    mail: "",
    topRank: "",
    joinTime: "",
    score: "0",
    state: "",
    rank: 0
  };

  private resolves: SubmitInfo[] = [];

  private infoText: string = "";
  private hasInfo: boolean = false;

  private get showJoinTime() {
    return new Date(
      this.userInfo.joinTime.toString().replace(/\//g, "-") + "+00:00"
    ).toLocaleString();
  }

  async mounted() {
    await this.load(this.$route.params.user_id);
    if (this.$store.state.competition.competition == true) {
      await this.loadResolve(this.$route.params.user_id);
    }
  }

  @Watch("$route")
  async beforeRouteUpdate(to: Route) {
    await this.load(to.params.user_id);
  }

  editAvatar() {
    window.open("https://www.gravatar.com/");
  }

  async load(userId: string) {
    try {
      let res = await this.$apollo.query<UserInfoResult, { userId: string }>({
        query: gql`
          query($userId: String!) {
            userInfo(userId: $userId) {
              message
              userInfo {
                userId
                name
                avatar
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
          userId: userId
        },
        fetchPolicy: "no-cache"
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

  async loadResolve(userId: string) {
    try {
      let res = await this.$apollo.query<
        SubmitHistoryResult,
        { input: string }
      >({
        query: gql`
          query($input: String!) {
            submitHistory(userId: $input) {
              message
              submitInfos {
                submitTime
                challengeName
                mark
              }
            }
          }
        `,
        variables: {
          input: userId
        },
        fetchPolicy: "no-cache"
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
}
</script>

<style lang="scss" scoped>
.avatar {
  position: absolute;
  left: 12px;
  width: 100%;
  z-index: 1;
}

.hoverable:hover {
  cursor: pointer;
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
