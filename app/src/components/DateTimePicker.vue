<template>
  <div class="content">
    <h4 class="body-2">{{ label }}</h4>
    <v-menu v-model="menu" :close-on-content-click="false">
      <template v-slot:activator="{ on }">
        <v-btn-toggle class="content">
          <v-btn class="btn" v-on="on" @click="DateClick" :disabled="disabled"
            ><v-spacer>{{ date != "" ? date : "DATE" }}</v-spacer></v-btn
          >
          <v-btn class="btn" v-on="on" @click="TimeClick" :disabled="disabled"
            ><v-spacer>{{ time != "" ? time : "TIME" }}</v-spacer></v-btn
          >
        </v-btn-toggle>
      </template>
      <v-card>
        <div class="d-flex justify-center">
          <v-date-picker
            class="elevation-0"
            v-if="!dateOrTime"
            v-model="menuDate"
          ></v-date-picker>
          <v-time-picker
            class="elevation-0"
            v-else
            v-model="menuTime"
          ></v-time-picker>
        </div>
        <v-card-actions>
          <v-btn text @click="cancel">cancel</v-btn>
          <v-btn text color="primary" @click="sure">ok</v-btn>
        </v-card-actions>
      </v-card>
    </v-menu>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Prop, Emit, Watch } from "vue-property-decorator";

@Component
export default class Competition extends Vue {
  @Prop() label!: String;
  @Prop() value!: Date;
  @Prop() disabled!: boolean;

  private date: String = "";
  private time: String = "";
  private dateOrTime: boolean = false;
  private menu: boolean = false;

  private menuDate: String = "";
  private menuTime: String = "";

  mounted() {
    this.refresh();
  }

  @Watch("value")
  refresh() {
    if (this.value != null) {
      this.date = `${this.value.getFullYear()}-${this.value.getMonth() +
        1}-${this.value.getDate()}`;
      this.time = `${this.value.getHours()}:${this.value.getMinutes()}`;
    } else {
      this.date = "";
      this.time = "";
    }
  }

  @Emit("input")
  SetDate() {
    if (this.date == "" || this.time == "") {
      return new Date();
    }
    let datetime = new Date(this.date + " " + this.time);
    return datetime;
  }

  DateClick() {
    this.menuDate = this.date || null;
    this.dateOrTime = false;
  }

  TimeClick() {
    this.menuTime = this.time || null;
    this.dateOrTime = true;
  }

  sure() {
    if (!this.dateOrTime) {
      this.date = this.menuDate;
    } else {
      this.time = this.menuTime;
    }
    this.menu = false;
    this.SetDate();
  }

  cancel() {
    this.menu = false;
  }
}
</script>

<style lang="scss" scoped>
.content {
  width: 100%;
}

.btn {
  width: 50%;
}

.menu-content {
  margin: 0 auto;
}
</style>
