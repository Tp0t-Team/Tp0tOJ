<template>
  <v-form>
    <v-row>
      <v-col cols="6">
        <v-text-field v-model="name" :disabled="loading || disabled" outlined label="name"></v-text-field>
      </v-col>
      <v-col cols="6">
        <v-text-field v-model="studentNumber" outlined disabled label="student number"></v-text-field>
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="6">
        <v-select
          v-model="role"
          :disabled="loading || disabled"
          outlined
          label="role"
          :items="roleType"
        ></v-select>
      </v-col>
      <v-col cols="6">
        <v-select
          v-model="state"
          :disabled="loading || disabled"
          outlined
          label="state"
          :items="stateType"
        ></v-select>
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="6">
        <v-text-field v-model="score" outlined disabled label="score"></v-text-field>
      </v-col>
      <v-col cols="6">
        <v-menu
          ref="dateMenu"
          v-model="dateMenu"
          :disabled="loading || disabled"
          :close-on-content-click="false"
        >
          <template v-slot:activator="{ on }">
            <v-text-field
              v-model="protectedTime"
              :disabled="loading || disabled"
              outlined
              label="protected time"
              v-on="on"
            ></v-text-field>
          </template>
          <v-date-picker v-model="protectedTime">
            <v-spacer></v-spacer>
            <v-btn text @click="dateMenu=false">Cancel</v-btn>
            <v-btn text color="primary" @click="$refs.dateMenu.save(menuDate)">OK</v-btn>
          </v-date-picker>
        </v-menu>
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="6">
        <v-select
          v-model="department"
          :disabled="loading || disabled"
          outlined
          label="department"
          :items="departmentItems"
        ></v-select>
      </v-col>
      <v-col cols="6">
        <v-select
          v-model="grade"
          :disabled="loading || disabled"
          outlined
          label="grade"
          :items="gradeItems"
        ></v-select>
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="6">
        <v-text-field v-model="qq" :disabled="loading || disabled" outlined label="QQ"></v-text-field>
      </v-col>
      <v-col cols="6">
        <v-text-field v-model="mail" :disabled="loading || disabled" outlined label="mail"></v-text-field>
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="6">
        <v-text-field v-model="topRank" outlined disabled label="top rank"></v-text-field>
      </v-col>
      <v-col>
        <v-text-field v-model="joinTime" outlined disabled label="join time"></v-text-field>
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
  private departmentItems = constValue.departmentItems;
  private gradeItems: string[] = [];

  @Prop() user!: UserInfo | null;
  @Prop() disabled!: boolean;
  @Prop() loading!: boolean;

  private name: string = "";
  private studentNumber: string = "";
  private role: string = "";
  private state: string = "";
  private score: number = 0;
  private protectedTime: string = new Date()
    .toLocaleDateString()
    .replace(/\//g, "-");
  private department: string = "";
  private grade: string = "";
  private qq: string = "";
  private mail: string = "";
  private topRank: string = "";
  private joinTime: string = "";

  private dateMenu: boolean = false;

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
    this.studentNumber = this.user.stuNumber;
    this.role = this.user.role;
    this.state = this.user.state;
    this.score = parseInt(this.user.score);
    this.protectedTime = this.user.protectedTime;
    this.department = this.user.department;
    this.grade = this.user.grade;
    this.qq = this.user.qq;
    this.mail = this.user.mail;
    this.topRank = this.user.topRank;
    this.joinTime = this.user.joinTime;
  }

  @Emit("submit")
  EmitSubmit() {
    return {
      userId: this.user!.userId,
      name: this.name,
      role: this.role,
      stuNumber: this.user!.stuNumber,
      department: this.department,
      grade: this.user!.grade,
      protectedTime: this.protectedTime,
      qq: this.qq,
      mail: this.mail,
      topRank: this.user!.topRank,
      joinTime: this.user!.joinTime,
      score: this.user!.score.toString(),
      state: this.state,
      rank: this.user!.rank
    } as UserInfo;
  }
}
</script>