<template>
  <v-container fluid class="fill-height">
    <v-row class="fill-height">
      <v-spacer></v-spacer>
      <v-col cols="6" class="content-col">
        <!-- show-expand item-key="challengeId" single-expand -->
        <v-data-table
          v-model="selected"
          :headers="headers"
          :items="challengeConfigs"
          @click:row="select"
          show-select
          :loading="loading"
          item-key="challengeId"
        >
          <template v-slot:body.append="{ headers }">
            <td :colspan="headers.length">
              <div class="action-group">
                <div class="action-item">
                  <v-btn
                    text
                    tile
                    block
                    color="primary"
                    @click="newChallenge"
                    :disabled="loading"
                    >new</v-btn
                  >
                </div>
                <div class="action-item">
                  <v-btn
                    text
                    tile
                    block
                    :disabled="!enableAble || loading"
                    @click="challengeAction('enable')"
                    >enable</v-btn
                  >
                </div>
                <div class="action-item">
                  <v-btn
                    text
                    tile
                    block
                    :disabled="!disableAble || loading"
                    @click="challengeAction('disable')"
                    >disable</v-btn
                  >
                </div>
                <div class="action-item">
                  <v-btn
                    text
                    tile
                    block
                    color="accent"
                    :disabled="selected.length == 0 || loading"
                    @click="challengeAction('delete')"
                    >delete</v-btn
                  >
                </div>
              </div>
            </td>
          </template>
          <template v-slot:item.category="{ item }">
            {{ item.config.category }}
          </template>
          <template v-slot:item.baseScore="{ item }">
            {{ item.config.score.baseScore }}
          </template>
        </v-data-table>
        <!-- <v-card class="ma-4" v-for="item in challengeConfigFiltered" :key="item.category">
          <v-toolbar dense>{{item.category}}</v-toolbar>
          <v-list dense>
            <v-list-item
              v-for="conf in item.items"
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
            <v-list-item @click="newChallenge(item.category)" :disabled="loading">
              <v-layout row>
                <v-spacer></v-spacer>
                <v-icon>add</v-icon>
                <v-spacer></v-spacer>
              </v-layout>
            </v-list-item>
          </v-list>
        </v-card> -->
      </v-col>
      <v-dialog v-model="showDiscardDialog" width="300px">
        <v-card>
          <v-card-title>Are you sure to discard changes?</v-card-title>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn text @click="showDiscardDialog = false">cancel</v-btn>
            <v-btn text color="primary" @click="continueChange">sure</v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
      <!-- <v-dialog v-model="showDeleteDialog" width="300px">
        <v-card>
          <v-card-title>Are you sure to delete this challenge?</v-card-title>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn text @click="showDeleteDialog = false">cancel</v-btn>
            <v-btn text color="primary" @click="deleteConfig">accept</v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog> -->
      <v-col cols="6" class="content-col">
        <challenge-editor
          :config="currentConfig"
          :disabled="withoutInit"
          :loading="loading"
          :changed="changed"
          @error="error"
          @change="Changed"
          @submit="submit"
          :key="(currentConfig && currentConfig.challengeId) || ''"
        ></challenge-editor>
      </v-col>
      <v-spacer></v-spacer>
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
import { Component, Vue, Watch } from "vue-property-decorator";
import gql from "graphql-tag";
import {
  // ChallengeConfig,
  ChallengeConfigWithId,
  ChallengeConfigResult,
  ChallengeMutateResult,
  ChallengeMutateInput,
  ChallengeActionResult,
  ChallengeActionInput
} from "@/struct";
import constValue from "@/constValue";
import ChallengeEditor from "@/components/ChallengeEditor.vue";

@Component({
  components: {
    ChallengeEditor
  }
})
export default class Challenge extends Vue {
  private headers = [
    { text: "name", value: "name" },
    { text: "category", value: "category" },
    { text: "base score", value: "baseScore" },
    { text: "state", value: "state" }
    // { text: '', value: 'data-table-expand' },
  ];

  private selected: ChallengeConfigWithId[] = [];

  private challengeType = constValue.challengeType;

  private showDiscardDialog: boolean = false;
  // private showDeleteDialog: boolean = false;
  private withoutInit: boolean = true;
  private loading: boolean = false;
  private changed: boolean = false;

  private challengeConfigs: ChallengeConfigWithId[] = [];
  private currentConfig: ChallengeConfigWithId | null = null;
  private tempConfig: ChallengeConfigWithId | null = null;
  private tempChallengeId: string = "";

  private infoText: string = "";
  private hasInfo: boolean = false;

  private enableAble: boolean = false;
  private disableAble: boolean = false;

  @Watch("selected")
  selectedChange() {
    this.enableAble = false;
    this.disableAble = false;
    for (let item of this.selected) {
      for (let config of this.challengeConfigs) {
        if (config.challengeId == item.challengeId) {
          if (config.state == "disabled") {
            this.enableAble = true;
          }
          if (config.state == "enabled") {
            this.disableAble = true;
          }
          break;
        }
      }
    }
  }

  // private get challengeConfigFiltered() {
  //   return this.challengeType.map(v => ({
  //     category: v,
  //     items: this.challengeConfigs.filter(c => c.config.category == v)
  //   }));
  // }

  async mounted() {
    await this.loadAll();
  }

  async loadAll() {
    try {
      let res = await this.$apollo.query<ChallengeConfigResult, {}>({
        query: gql`
          query {
            challengeConfigs {
              message
              challengeConfigs {
                challengeId
                state
                name
                config {
                  category
                  score {
                    dynamic
                    baseScore
                  }
                  flag {
                    type
                    value
                  }
                  description
                  externalLink
                  singleton
                  nodeConfig {
                    name
                    image
                    servicePorts {
                      name
                      protocol
                      external
                      internal
                      pod
                    }
                  }
                }
              }
            }
          }
        `,
        fetchPolicy: "no-cache"
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
    let tempConfig: ChallengeMutateInput = JSON.parse(
      JSON.stringify(config.config)
    );
    tempConfig.challengeId =
      config.challengeId[0] == "-" ? "" : config.challengeId;
    tempConfig.name = config.name;
    tempConfig.state = config.state;
    try {
      let res = await this.$apollo.mutate<
        ChallengeMutateResult,
        { input: ChallengeMutateInput }
      >({
        mutation: gql`
          mutation($input: ChallengeMutateInput!) {
            challengeMutate(input: $input) {
              message
            }
          }
        `,
        variables: { input: tempConfig }
      });
      if (res.errors) throw res.errors.map(v => v.message).join(",");
      if (res.data!.challengeMutate.message)
        throw res.data!.challengeMutate.message;
      this.loading = false;
      this.changed = false;
      this.currentConfig = null;
      this.withoutInit = true;
      await this.loadAll();
      this.infoText = "add / update success";
      this.hasInfo = true;
    } catch (e) {
      this.loading = false;
      this.infoText = e.toString();
      this.hasInfo = true;
    }
  }

  select(challenge: ChallengeConfigWithId) {
    this.editChallenge(challenge.challengeId);
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

  // newChallenge(type: string) {
  newChallenge() {
    let config: ChallengeConfigWithId = {
      challengeId: "-" + Date.now().toLocaleString(),
      name: "",
      state: "disabled",
      config: {
        // category: type,
        category: "",
        score: { dynamic: false, baseScore: "0" },
        flag: { type: 0, value: "" },
        description: "",
        externalLink: [],
        singleton: true,
        nodeConfig: undefined
      }
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
    if (this.tempConfig == null) {
      this.withoutInit = true;
    }
    this.currentConfig = this.tempConfig;
  }

  // tryDelete(id: string) {
  //   this.tempChallengeId = id;
  //   this.showDeleteDialog = true;
  // }

  async deleteConfig() {
    // console.log(this.tempChallengeId);
    this.infoText = "功能暂未实现";
    this.hasInfo = true;
  }

  async challengeAction(action: string) {
    if (this.changed) {
      this.showDiscardDialog = true;
      this.tempConfig = null;
      return;
    }
    this.loading = true;
    try {
      let res = await this.$apollo.mutate<
        ChallengeActionResult,
        { input: ChallengeActionInput }
      >({
        mutation: gql`
          mutation($input: ChallengeActionInput!) {
            challengeAction(input: $input) {
              message
              successful
            }
          }
        `,
        variables: {
          input: {
            action: action,
            challengeIds: this.selected.map(it => it.challengeId)
          }
        }
      });
      if (res.errors) throw res.errors.map(v => v.message).join(",");
      if (res.data!.challengeAction.message)
        throw res.data!.challengeAction.message;
      this.selected = [];
      this.loading = false;
      this.changed = false;
      this.currentConfig = null;
      this.withoutInit = true;
      this.infoText = `${action} success`;
      this.hasInfo = true;
      await this.loadAll();
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

.action-group {
  display: flex;
  flex-direction: row;
}

.action-item {
  flex-basis: 25%;
}
</style>
