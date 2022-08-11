<template>
  <v-form>
    <v-row>
      <v-col cols="6">
        <v-text-field
          v-model="name"
          :disabled="loading || disabled"
          outlined
          label="name"
        ></v-text-field>
      </v-col>
      <v-col cols="6">
        <v-select
          v-model="role"
          :disabled="loading || disabled"
          outlined
          label="role"
          :items="roleType"
        ></v-select>
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="6">
        <v-select
          v-model="state"
          :disabled="loading || disabled"
          outlined
          label="state"
          :items="stateType"
        ></v-select>
      </v-col>
      <v-col cols="6">
        <v-text-field
          v-model="mail"
          :disabled="loading || disabled"
          outlined
          label="mail"
        ></v-text-field>
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <v-text-field
          v-model="joinTime"
          outlined
          disabled
          label="join time"
        ></v-text-field>
      </v-col>
    </v-row>
    <v-btn
      fab
      color="primary"
      :loading="loading"
      :disabled="disabled"
      absolute
      right
      bottom
      @click="EmitSubmit"
    >
      <v-icon>done</v-icon>
    </v-btn>
  </v-form>
</template>

<script lang="ts">
import { Component, Vue, Prop, Emit } from "vue-property-decorator";
import constValue from "@/constValue";
import { UserInfo } from "@/struct";

@Component
export default class UserEditor extends Vue {
  private roleType = constValue.roleType;
  private stateType = constValue.stateType;
  // private departmentItems = constValue.departmentItems;
  private gradeItems: string[] = [];

  @Prop() user!: UserInfo | null;
  @Prop() disabled!: boolean;
  @Prop() loading!: boolean;

  private name: string = "";
  private role: string = "";
  private state: string = "";
  private score: number = 0;
  private mail: string = "";
  private joinTime: string = "";

  created() {
    let year = new Date().getFullYear() - 9;
    let items = [];
    for (let i = 0; i < 10; i++) {
      items.push((year + i).toString());
    }
    this.gradeItems = items;
  }

  mounted() {
    if (!this.user) return;
    this.name = this.user.name;
    this.role = this.user.role;
    this.state = this.user.state;
    this.score = parseInt(this.user.score);
    this.mail = this.user.mail;
    this.joinTime = new Date(
      this.user.joinTime //.replace(/\//g, "-") + "+00:00"
    ).toLocaleString();
  }

  @Emit("submit")
  EmitSubmit() {
    return {
      userId: this.user!.userId,
      name: this.name,
      role: this.role,
      mail: this.mail,
      joinTime: this.user!.joinTime,
      score: this.user!.score.toString(),
      state: this.state
    } as UserInfo;
  }
}
</script>
