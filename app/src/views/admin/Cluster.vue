<template>
  <div>
    <v-container fill-width>
      <v-row>
        <v-col cols="7">
          <v-data-table :headers="nodeHeaders" :items="nodes"></v-data-table>
        </v-col>
        <v-col cols="5">
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
import { ClusterInfoResult, ClusterNodeInfo, ClusterReplicaInfo, ImageInfo, ImageInfoResult, Result } from "@/struct";

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

  private nodes: ClusterNodeInfo[] = [];
  private replicas: ClusterReplicaInfo[] = [];

  private infoText: string = "";
  private hasInfo: boolean = false;

  async mounted() {
    await this.loadData();
  }

  async loadData() {
    try {
      let res = await this.$apollo.query<ClusterInfoResult, {}>({
        query: gql`
          query {
            clusterInfo {
              message
              nodes {
                name
                osType
                distribution
                kernel
                arch
              }
              replicas {
                name
                node
                status
              }
            }
          }
        `,
        fetchPolicy: "no-cache",
      });
      if (res.errors) throw res.errors.map((v) => v.message).join(",");
      if (res.data!.clusterInfo.message) throw res.data!.clusterInfo.message;
      this.nodes = res.data!.clusterInfo.nodes;
      this.replicas = res.data!.clusterInfo.replicas;
    } catch (e) {
      this.infoText = e.toString();
      this.hasInfo = true;
    }
  }
}
</script>

<style lang="scss" scoped>
</style>
