<template>
  <v-container fluid class="scrollable">
    <div v-for="type in challengeType" :key="type">
      <div v-if="challenges.filter((v)=>v.type==type).length!=0">
        <div class="type-title display-1 ml-4">{{type}}</div>
        <v-divider color="primary"></v-divider>
      </div>
      <v-row>
        <v-col sm="3" v-for="item in challenges.filter((v)=>v.type==type)" :key="item.challengeId">
          <v-layout justify-center>
            <v-hover :disabled="item.done" v-slot:default="{ hover }">
              <v-card
                :elevation="hover ? 12 : 2"
                width="20vw"
                height="220px"
                class="ma-4"
                @click="openDetial(item.challengeId)"
              >
                <v-badge left overlap class="score" z-index="2">
                  <template v-slot:badge>{{item.score}}</template>
                  <v-card-title class="subtitle-1 text-truncate">{{item.name}}</v-card-title>
                </v-badge>
                <v-divider color="primary" class="divider"></v-divider>
                <v-card-text class="description">{{item.description}}</v-card-text>
                <v-overlay absolute :value="item.done" color="green" z-index="0">
                  <v-icon>done</v-icon>
                </v-overlay>
                <v-card-actions>
                  <v-row>
                    <v-col
                      cols="4"
                      v-for="blood in item.blood"
                      :key="item.challengeId+blood.userId"
                    >
                      <v-layout row>
                        <v-spacer></v-spacer>
                        <v-btn color="blue" fab @click.stop="seeBlood(blood.userId)">
                          <!-- <v-avatar size="56">{{blood.name[0]}}</v-avatar> -->
                          <v-avatar size="56">
                            <user-avatar class="white--text" :url="blood.avatar" :size="56" :name="blood.name"></user-avatar>
                          </v-avatar>
                        </v-btn>
                        <v-spacer></v-spacer>
                      </v-layout>
                    </v-col>
                  </v-row>
                </v-card-actions>
              </v-card>
            </v-hover>
          </v-layout>
        </v-col>
      </v-row>
    </div>
    <v-dialog
      v-model="showDialog"
      :persistent="loading"
      width="800px"
      v-if="currentChallenge!=null"
    >
      <v-card width="800px" height="600px">
        <v-toolbar dense>
          <v-spacer></v-spacer>
          <v-toolbar-title>
            {{currentChallenge.name}}
            <v-chip class="ml-2">{{currentChallenge.score}}pt</v-chip>
          </v-toolbar-title>
          <v-spacer></v-spacer>
        </v-toolbar>
        <v-text-field
          v-model="sumbitFlag"
          outlined
          class="ma-4 dialog-flag"
          label="flag"
          append-icon="send"
          :disabled="loading"
          :loading="loading"
          @click:append="submit"
          :error-messages="submitError"
          @focus="submitError = ''"
          @blur="check"
        ></v-text-field>
        <div class="dialog-discription pa-6">
          <pre>{{formatedHint + currentChallenge.description}}</pre>
        </div>
        <div class="url-list">
          <v-chip
            color="primary"
            label
            outlined
            class="ma-4"
            v-for="link in currentChallenge.externalLink"
            :key="link"
          >
            {{link}}
            <v-btn class="ml-4" icon color="primary" @click="openUrl(link)">
              <v-icon>navigate_next</v-icon>
            </v-btn>
          </v-chip>
        </div>
      </v-card>
    </v-dialog>
    <v-snackbar v-model="hasInfo" top :timeout="3000">
      {{ infoText }}
      <v-spacer></v-spacer>
      <v-btn icon>
        <v-icon @click="hasInfo = false">close</v-icon>
      </v-btn>
    </v-snackbar>
  </v-container>
</template>

<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
import gql from "graphql-tag";
import UserAvatar from "@/components/UserAvatar.vue";
import {
  ChallengeDesc,
  ChallengeResult,
  SubmitResult,
  SubmitInput
} from "@/struct";
import constValue from "@/constValue";

@Component({
  components: {
    UserAvatar
  }
})
export default class Challenge extends Vue {
  private challengeType = constValue.challengeType;

  private challenges: ChallengeDesc[] = [];

  private sumbitFlag: string = "";
  private submitError: string = "";
  private valid: boolean = false;

  private showDialog: boolean = false;
  private currentChallenge: ChallengeDesc | null = null;
  private loading: boolean = false;

  private infoText: string = "";
  private hasInfo: boolean = false;

  private get formatedHint() {
    if (!this.currentChallenge) return "";
    return this.currentChallenge.hint.reduce(
      (pre: string, cur: string) => `Hint: ${cur}\n${pre}`,
      ""
    );
  }

  async mounted() {
    try {
      let res = await this.$apollo.query<ChallengeResult>({
        query: gql`
          query {
            challengeInfos {
              message
              challengeInfos {
                challengeId
                name
                type
                description
                externalLink
                hint
                score
                blood {
                  userId
                  name
                  avatar
                }
                done
              }
            }
          }
        `,
        fetchPolicy: "no-cache"
      });
      if (res.errors) throw res.errors.map(v => v.message).join(",");
      if (res.data!.challengeInfos.message)
        throw res.data!.challengeInfos.message;
      this.challenges = res.data!.challengeInfos.challengeInfos;
    } catch (e) {
      this.infoText = e.toString();
      this.hasInfo = true;
    }
  }

  openDetial(id: string) {
    let c = this.challenges.find(v => v.challengeId == id);
    if (!c || c.done) return;
    this.currentChallenge = c;
    this.sumbitFlag = "";
    this.submitError = "";
    this.showDialog = true;
  }

  seeBlood(id: string) {
    this.$router.push(`/profile/${id}`);
  }

  check() {
    if (!this.sumbitFlag) this.submitError = "请填写";
    else this.submitError = "";
  }

  openUrl(url: string) {
    window.open(url);
  }

  async submit() {
    this.check();
    if (!!this.submitError) return;
    let coreFlag = this.sumbitFlag.trim();
    this.loading = true;
    try {
      if (coreFlag.substr(0, 5) != "flag{") throw "error format";
      if (coreFlag[coreFlag.length - 1] != "}") throw "error format";
      coreFlag = coreFlag.substring(5, coreFlag.length - 1);
      let res = await this.$apollo.mutate<SubmitResult, SubmitInput>({
        mutation: gql`
          mutation($input: SubmitInput!) {
            submit(input: $input) {
              message
            }
          }
        `,
        variables: {
          input: {
            challengeId: this.currentChallenge!.challengeId,
            flag: coreFlag
          }
        }
      });
      if (res.errors) throw res.errors.map(v => v.message).join(",");
      if (res.data!.submit.message) throw res.data!.submit.message;
      // throw "提交成功";
      this.$router.replace({
        path: "/challenge",
        query: { time: Date.now().toLocaleString() }
      });
    } catch (e) {
      this.loading = false;
      this.infoText = e.toString();
      this.hasInfo = true;
    }
  }
}
</script>

<style lang="scss" scoped>
.scrollable {
  height: calc(100vh - 96px);
  overflow-x: hidden;
  overflow-y: auto;
}

.type-title {
  width: 100%;
}

.divider {
  position: relative;
  left: 5%;
  width: 90%;
}

.score {
  z-index: 5;
}

.description {
  height: 120px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 5;
  -webkit-box-orient: vertical;
  display: -moz-box;
  -moz-line-clamp: 5;
  -moz-box-orient: vertical;
}

.dialog-discription {
  height: 300px;
  overflow-y: auto;
}

.url-list {
  height: 128px;
  overflow-y: auto;
}
</style>