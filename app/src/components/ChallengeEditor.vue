<template>
  <v-form>
    <v-row>
      <!-- <v-col cols="4">
        <v-row>
          <v-spacer></v-spacer>
          <v-switch
            v-model="state"
            label="enabled"
            :disabled="loading || disabled || config.challengeId[0] == '-'"
            @change="Changed"
          ></v-switch>
          <v-spacer></v-spacer>
        </v-row>
      </v-col> -->
      <v-col cols="6">
        <v-row>
          <v-spacer></v-spacer>
          <v-switch
            v-model="dynamicScore"
            hide-details
            label="DynamicScore"
            :disabled="loading || disabled || config.challengeId[0] != '-'"
            @change="Changed"
          ></v-switch>
          <v-spacer></v-spacer>
        </v-row>
      </v-col>
      <v-col cols="6">
        <v-row>
          <v-spacer></v-spacer>
          <v-switch
            v-model="singleton"
            hide-details
            label="Singleton"
            :disabled="loading || disabled || config.challengeId[0] != '-'"
            @change="Changed"
          ></v-switch>
          <v-spacer></v-spacer>
        </v-row>
      </v-col>
      <v-col cols="12">
        <v-file-input
          v-model="configFile"
          accept=".yaml, .yml"
          outlined
          hide-details
          prepend-icon=""
          prepend-inner-icon="attach_file"
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
          hide-details
          label="name"
          :disabled="loading || disabled"
          @change="Changed"
        ></v-text-field>
      </v-col>
      <v-col cols="6">
        <!-- <v-text-field v-model="type" outlined label="type" readonly value="web" :disabled="loading"></v-text-field> -->
        <v-select
          v-model="type"
          :items="typeItems"
          outlined
          hide-details
          label="type"
          :disabled="loading || disabled"
          @change="Changed"
        ></v-select>
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="6">
        <v-text-field
          v-model="score"
          outlined
          hide-details
          label="score"
          type="number"
          :disabled="loading || disabled"
          @change="Changed"
        ></v-text-field>
      </v-col>
      <v-col cols="6">
        <!-- <v-text-field v-model="type" outlined label="type" readonly value="web" :disabled="loading"></v-text-field> -->
        <v-select
          v-model="flagType"
          :items="flagTypeItems"
          outlined
          hide-details
          label="flag type"
          :disabled="loading || disabled"
          @change="Changed"
        ></v-select>
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="12" v-if="flagType == 'Multiple'">
        <v-textarea
          v-model="flag"
          outlined
          hide-details
          label="flag"
          :disabled="loading || disabled"
          @change="Changed"
        ></v-textarea>
      </v-col>
      <v-col cols="12" v-else>
        <v-text-field
          v-model="flag"
          outlined
          hide-details
          label="flag"
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
          hide-details
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
            <v-list-item v-for="(l,index) in links" :key="l" @click=";">
              <v-list-item-content>{{l}}</v-list-item-content>
              <v-list-item-icon>
                <v-btn icon :disabled="loading" @click="removeLink(index)">
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
import { ChallengeConfig, ChallengeConfigWithId, NodeConfig } from "@/struct";
import constValue from "@/constValue";

@Component
export default class ChallengeEditor extends Vue {
  private challengeType = constValue.challengeType;

  private typeItems = constValue.challengeType;
  private flagTypeItems = constValue.flagType;

  @Prop() config!: ChallengeConfigWithId | null;
  @Prop() disabled!: boolean;
  @Prop() loading!: boolean;
  @Prop() changed!: boolean;

  private configFile: File | null = null;
  private reader: FileReader = new FileReader();

  private state: boolean = false;
  private dynamicScore: boolean = false;
  // private dynamicFlag: boolean = false;
  private name: string = "";
  private type: string = "";
  private score: number = 0;
  private flagType: string = "Dynamic";
  private flag: string = "";
  private description: string = "";
  private singleton: boolean = true;
  private nodeConfigs: NodeConfig[]|undefined = undefined;
  private links: string[] = [];

  private link: string = "";

  private setValue: boolean = false;

  mounted() {
    this.reader.addEventListener("load", e => {
      try {
        let config: ChallengeConfig &{name: string} = loadYaml(e.target!.result as string);
        if (this.challengeType.findIndex(v => v == config.category) < 0)
          throw "类型错误";
        // if (config.category != this.type) throw "不可修改类型";
        this.setValue = true;
        this.name = config.name || this.name;
        this.type = config.category || this.type;
        this.score = (config.score.baseScore as number) || this.score;
        this.flag = config.flag.value || this.flag;
        this.description = config.description || this.description;
        this.links = config.externalLink || this.links;
        this.setValue = false;
        this.state = false;
        this.dynamicScore = config.score.dynamic || this.dynamicScore;
        this.flagType = config.flag.type != undefined ? this.flagTypeItems[config.flag.type]: null ?? this.flagType;
        // this.dynamicFlag = config.flag.dynamic || this.dynamicFlag;
        this.singleton = config.singleton;
        this.nodeConfigs = config.nodeConfig;
      } catch (e) {
        this.EmitError(e.toString());
      }
    });
    if (!this.config) return;
    this.name = this.config.name;
    this.type = this.config.config.category;
    this.score = parseInt(this.config.config.score.baseScore as string);
    this.flag = this.config.config.flag.value;
    this.description = this.config.config.description;
    this.links = this.config.config.externalLink;
    this.state = this.config.state == "enabled";
    this.dynamicScore = this.config.config.score.dynamic;
    // this.dynamicFlag = this.config.config.flag.dynamic;
    this.flagType = this.flagTypeItems[this.config.config.flag.type];
    this.singleton = this.config.config.singleton;
    this.nodeConfigs = this.config.config.nodeConfig;
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
      state: this.state ? "enabled" : "disabled",
      config: {
        category: this.type,
        score: { dynamic: this.dynamicScore, baseScore: this.score.toString() },
        flag: { type: this.flagTypeItems.indexOf(this.flagType), value: this.flag },
        description: this.description,
        externalLink: this.links,
        singleton: this.singleton,
        nodeConfig: this.nodeConfigs
      }
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

  Changed() {
    if (this.setValue) return;
    this.EmitChange();
  }
}
</script>