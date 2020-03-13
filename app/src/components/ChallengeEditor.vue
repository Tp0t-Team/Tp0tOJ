<template>
  <v-form>
    <v-row>
      <v-col cols="6">
        <v-row>
          <v-spacer></v-spacer>
          <v-switch
            v-model="state"
            label="enabled"
            :disabled="loading || disabled"
            @change="Changed"
          ></v-switch>
          <v-spacer></v-spacer>
        </v-row>
      </v-col>
      <v-col cols="6">
        <v-file-input
          v-model="configFile"
          accept=".yaml, .yml"
          outlined
          label="load config file"
          :disabled="loading || disabled"
          @change="Changed"
        ></v-file-input>
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="6">
        <v-text-field
          v-model="name"
          outlined
          label="name"
          :disabled="loading || disabled"
          @change="Changed"
        ></v-text-field>
      </v-col>
      <v-col cols="6">
        <v-text-field
          v-model="type"
          outlined
          label="type"
          readonly
          value="web"
          :disabled="loading"
        ></v-text-field>
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="6">
        <v-text-field
          v-model="score"
          outlined
          label="score"
          type="number"
          :disabled="loading || disabled"
          @change="Changed"
        ></v-text-field>
      </v-col>
      <v-col cols="6">
        <v-text-field
          v-model="flag"
          outlined
          label="flag"
          :disabled="loading || disabled"
          @change="Changed"
        ></v-text-field>
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="6">
        <v-row>
          <v-spacer></v-spacer>
          <v-switch
            v-model="dynamicScore"
            label="DynamicScore"
            :disabled="loading || disabled"
            @change="Changed"
          ></v-switch>
          <v-spacer></v-spacer>
        </v-row>
      </v-col>
      <v-col cols="6">
        <v-row>
          <v-spacer></v-spacer>
          <v-switch
            v-model="proxiedFlag"
            label="ProxiedFlag"
            :disabled="loading || disabled"
            @change="Changed"
          ></v-switch>
          <v-spacer></v-spacer>
        </v-row>
      </v-col>
    </v-row>
    <v-row v-if="proxiedFlag">
      <v-col cols="6">
        <v-text-field
          v-model="portFrom"
          :value="portFrom"
          outlined
          label="port from"
          type="number"
          :disabled="loading || disabled"
          @change="Changed"
        ></v-text-field>
      </v-col>
      <v-col cols="6">
        <v-text-field
          v-model="portTo"
          :value="portTo"
          outlined
          label="portTo"
          type="number"
          :disabled="loading || disabled"
          @change="Changed"
        ></v-text-field>
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="12">
        <v-textarea
          v-model="description"
          outlined
          label="description"
          :disabled="loading || disabled"
          @change="Changed"
        ></v-textarea>
      </v-col>
      <!--  -->
      <v-col cols="12">
        <v-card>
          <v-list dense>
            <v-list-item>
              <v-text-field
                v-model="link"
                label="external link"
                append-icon="add"
                @click:append="addLink"
                :disabled="loading || disabled"
              ></v-text-field>
            </v-list-item>
            <v-list-item v-for="(l, index) in links" :key="l" @click="">
              <v-list-item-content>{{ l }}</v-list-item-content>
              <v-list-item-icon>
                <v-btn icon :disabled="loading" @click="removeLink(index)">
                  <v-icon>close</v-icon>
                </v-btn>
              </v-list-item-icon>
            </v-list-item>
          </v-list>
        </v-card>
      </v-col>
      <!--  -->
      <v-col cols="12">
        <v-card>
          <v-list dense>
            <v-list-item>
              <v-text-field
                v-model="hint"
                label="hint"
                append-icon="add"
                @click:append="addHint"
                :disabled="loading || disabled"
              ></v-text-field>
            </v-list-item>
            <v-list-item v-for="(l, index) in hints" :key="l" @click="">
              <v-list-item-content>{{ l }}</v-list-item-content>
              <v-list-item-icon>
                <v-btn
                  icon
                  :disabled="loading || disabled"
                  @click="removeHint(index)"
                >
                  <v-icon>close</v-icon>
                </v-btn>
              </v-list-item-icon>
            </v-list-item>
          </v-list>
        </v-card>
      </v-col>
    </v-row>
    <v-btn
      fab
      absolute
      right
      bottom
      color="primary"
      :loading="loading"
      :disabled="loading || !changed || disabled"
      @click="EmitSubmit"
    >
      <v-icon>done</v-icon>
    </v-btn>
  </v-form>
</template>

<script lang="ts">
import { Component, Vue, Watch, Prop, Emit } from "vue-property-decorator";
import { load as loadYaml } from "js-yaml";
import { ChallengeConfig, ChallengeConfigWithId } from "@/struct";
import constValue from "@/constValue";

@Component
export default class ChallengeEditor extends Vue {
  private challengeType = constValue.challengeType;

  @Prop() config!: ChallengeConfigWithId | null;
  @Prop() disabled!: boolean;
  @Prop() loading!: boolean;
  @Prop() changed!: boolean;

  private configFile: File | null = null;
  private reader: FileReader = new FileReader();

  private state: boolean = false;
  private dynamicScore: boolean = false;
  private proxiedFlag: boolean = false;
  private portFrom: number = 0;
  private portTo: number = 0;
  private name: string = "";
  private type: string = "";
  private score: number = 0;
  private flag: string = "";
  private description: string = "";
  private links: string[] = [];
  private hints: string[] = [];

  private link: string = "";
  private hint: string = "";

  private setValue: boolean = false;

  mounted() {
    this.reader.addEventListener("load", e => {
      try {
        let config: ChallengeConfig = loadYaml(e.target!.result as string);
        if (this.challengeType.findIndex(v => v == config.type) < 0)
          throw "类型错误";
        if (config.type != this.type) throw "不可修改类型";
        this.setValue = true;
        this.name = config.name || this.name;
        this.score = (config.score.base_score as number) || this.score;
        this.flag = config.flag.value || this.flag;
        this.description = config.description || this.description;
        this.links = config.external_link || this.links;
        this.hints = config.hint || this.hints;
        this.setValue = false;
        this.state = false;
        this.dynamicScore = config.score.dynamic || this.dynamicScore;
        this.proxiedFlag = config.flag.dynamic || this.proxiedFlag;
        this.portFrom = config.flag.portFrom || this.portFrom;
        this.portTo = config.flag.portTo || this.portTo;
      } catch (e) {
        this.EmitError(e.toString());
      }
    });
    if (!this.config) return;
    this.name = this.config.name;
    this.type = this.config.type;
    this.score = parseInt(this.config.score.base_score as string);
    this.flag = this.config.flag.value;
    this.description = this.config.description;
    this.links = this.config.external_link;
    this.hints = this.config.hint;
    this.state = this.config.state == "enabled";
    this.dynamicScore = this.config.score.dynamic;
    this.proxiedFlag = this.config.flag.dynamic;
    this.portFrom = this.config.flag.portFrom;
    this.portTo = this.config.flag.portTo;
  }

  @Emit("error")
  EmitError(error: string) {
    return error;
  }

  @Emit("change")
  EmitChange() {}

  @Emit("submit")
  EmitSubmit() {
    return {
      challengeId: (this.config && this.config.challengeId) || "",
      name: this.name,
      type: this.type,
      score: { dynamic: this.dynamicScore, base_score: this.score.toString() },
      flag: { dynamic: this.proxiedFlag, value: this.flag, portFrom: this.portFrom, portTo: this.portTo },
      description: this.description,
      external_link: this.links,
      hint: this.hints,
      state: this.state ? "enabled" : "disabled"
    } as ChallengeConfigWithId;
  }

  @Watch("configFile")
  loadConfigFile() {
    if (!this.configFile) return;
    this.reader.readAsText(this.configFile);
  }

  addLink() {
    this.links = [...this.links, this.link];
    this.link = "";
    this.Changed();
  }

  removeLink(index: number) {
    this.links.splice(index, 1);
    this.links = this.links;
    this.Changed();
  }

  addHint() {
    this.hints = [...this.hints, this.hint];
    this.hint = "";
    this.Changed();
  }

  removeHint(index: number) {
    this.hints.splice(index, 1);
    this.hints = this.hints;
    this.Changed();
  }

  Changed() {
    if (this.setValue) return;
    this.EmitChange();
  }
}
</script>
