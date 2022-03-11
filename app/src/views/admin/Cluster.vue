<template>
  <div>
    <v-container fill-width>
      <v-row>
        <v-col cols="6">
          <v-data-table :headers="nodeHeaders" :items="nodes"></v-data-table>
        </v-col>
        <v-col cols="6">
          <v-data-table :headers="replicaHeaders" :items="replicas"></v-data-table>
        </v-col>
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
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import gql from "graphql-tag";
import { ImageInfo, ImageInfoResult, Result } from "@/struct";

@Component
export default class Cluster extends Vue {
  private nodeHeaders = [
    { text: "name", value: "name" },
    { text: "os type", value: "osType" },
    { text: "distribution", value: "distribution" },
    { text: "kernel version", value: "kernel" },
    { text: "arch", value: "arch" },
  ];

  private replicaHeaders = [
    { text: "name", value: "name" },
    { text: "node", value: "node" },
    { text: "status", value: "status" },
  ];

//   private nodes: ImageInfo[] = [];

  private infoText: string = "";
  private hasInfo: boolean = false;

  async mounted() {
  }
}
</script>

<style lang="scss" scoped>
</style>
