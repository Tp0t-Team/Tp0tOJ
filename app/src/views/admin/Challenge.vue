<template>
  <v-container fluid class="fill-height">
    <v-row class="fill-height">
      <v-col cols="6" class="content-col">
        <v-card class="ma-4" v-for="type in challengeType" :key="type">
          <v-toolbar dense>{{type}}</v-toolbar>
          <v-list dense>
            <v-list-item
              v-for="conf in challengeConfigs.filter((c)=>c.type==type)"
              :key="conf.challengeId"
              @click="editChallenge(conf.challengeId)"
              :disabled="loading"
            >
              <v-list-item-content>{{conf.name}}</v-list-item-content>
              <v-list-item-icon>
                <v-btn icon :disabled="loading" @click.stop="tryDelete(conf.challengeId)">
                  <v-icon>close</v-icon>
                </v-btn>
              </v-list-item-icon>
            </v-list-item>
            <v-list-item @click="newChallenge(type)" :disabled="loading">
              <v-layout row>
                <v-spacer></v-spacer>
                <v-icon>add</v-icon>
                <v-spacer></v-spacer>
              </v-layout>
            </v-list-item>
          </v-list>
        </v-card>
      </v-col>
      <v-dialog v-model="showDiscardDialog" width="300px">
        <v-card>
          <v-card-title>Are you sure to discard changes?</v-card-title>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn text @click="showDiscardDialog=false">cancel</v-btn>
            <v-btn text color="primary" @click="continueChange">sure</v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
      <v-dialog v-model="showDeleteDialog" width="300px">
        <v-card>
          <v-card-title>Are you sure to delete this challenge?</v-card-title>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn text @click="showDeleteDialog=false">cancel</v-btn>
            <v-btn text color="primary" @click="deleteConfig">accept</v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
      <v-col cols="6" class="content-col">
        <challenge-editor
          :config="currentConfig"
          :disabled="withoutInit"
          :loading="loading"
          :changed="changed"
          @error="error"
          @change="Changed"
          @submit="submit"
          :key="currentConfig && currentConfig.challengeId || ''"
        ></challenge-editor>
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
import { Component, Vue, Watch } from "vue-property-decorator";
import gql from "graphql-tag";
import {
  ChallengeConfig,
  ChallengeConfigWithId,
  ChallengeConfigResult,
  ChallengeMutateResult
} from "@/struct";
import constValue from "@/constValue";
import ChallengeEditor from "@/components/ChallengeEditor.vue";

@Component({
  components: {
    ChallengeEditor
  }
})
export default class Challenge extends Vue {
  private challengeType = constValue.challengeType;

  private showDiscardDialog: boolean = false;
  private showDeleteDialog: boolean = false;
  private withoutInit: boolean = true;
  private loading: boolean = false;
  private changed: boolean = false;

  private challengeConfigs: ChallengeConfigWithId[] = [];
  private currentConfig: ChallengeConfigWithId | null = null;
  private tempConfig: ChallengeConfigWithId | null = null;
  private tempChallengeId: string = "";

  private infoText: string = "";
  private hasInfo: boolean = false;

  async mounted() {
    // example data
    this.challengeConfigs = [
      {
        challengeId: "1",
        type: "WEB",
        name: "easy_php",
        score: { dynamic: false, base_score: 1000 },
        flag: { dynamic: false, value: "qweertytryrty" },
        description: "123456",
        external_link: ["http://www.google.com"],
        hint: ["123456"],
        state: "enabled"
      }
    ];
    //
    try {
      let res = await this.$apollo.query<ChallengeConfigResult>({
        query: gql`
          query {
            challenges {
              message
              challengeInfos {
                challengeId
                type
                name
                score
                flag
                description
                externalLink
                hint
                state
              }
            }
          }
        `
      });
      if (res.errors) throw res.errors.map(v => v.message).join(",");
      if (res.data!.challengeConfigs.message)
        throw res.data!.challengeConfigs.message;
      this.challengeConfigs = res.data!.challengeConfigs.challengeConfigs;
    } catch (e) {
      this.infoText = e.toString();
      this.hasInfo = true;
    }
  }

  error(error: string) {
    this.infoText = error;
    this.hasInfo = true;
  }

  Changed() {
    this.changed = true;
  }

  async submit(config: ChallengeConfigWithId) {
    this.loading = true;
    let tempConfig = JSON.parse(JSON.stringify(config));
    tempConfig.challengeId =
      tempConfig.challengeId[0] == "-" ? "" : tempConfig.challengeId;
    try {
      let res = await this.$apollo.mutate<
        ChallengeMutateResult,
        ChallengeConfigWithId
      >({
        mutation: gql`
          mutation($input: ChallengeMutateInput) {
            challengeMutate(input: $input) {
              message
            }
          }
        `,
        variables: config
      });
      this.loading = false;
    } catch (e) {
      this.loading = false;
      this.infoText = e.toString();
      this.hasInfo = true;
    }
  }

  editChallenge(id: string) {
    let config = this.challengeConfigs.find(v => v.challengeId == id);
    if (!config) return;
    if (this.changed) {
      this.tempConfig = config;
      this.showDiscardDialog = true;
    } else {
      this.changed = false;
      this.withoutInit = false;
      this.currentConfig = config;
    }
  }

  newChallenge(type: string) {
    let config: ChallengeConfigWithId = {
      challengeId: "-" + Date.now().toLocaleString(),
      type: type,
      name: "",
      score: { dynamic: false, base_score: 0 },
      flag: { dynamic: false, value: "" },
      description: "",
      external_link: [],
      hint: [],
      state: "disabled"
    };
    if (this.changed) {
      this.tempConfig = config;
      this.showDiscardDialog = true;
    } else {
      this.changed = false;
      this.withoutInit = false;
      this.currentConfig = config;
    }
  }

  continueChange() {
    this.showDiscardDialog = false;
    this.changed = false;
    this.withoutInit = false;
    this.currentConfig = this.tempConfig;
  }

  tryDelete(id: string) {
    this.tempChallengeId = id;
    this.showDeleteDialog = true;
  }

  async deleteConfig() {
    console.log(this.tempChallengeId);
  }
}
</script>

<style lang="scss" scoped>
.content-col {
  height: calc(100vh - 120px);
  overflow-y: auto;
}
</style>