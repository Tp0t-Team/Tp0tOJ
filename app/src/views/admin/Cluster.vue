<template>
  <div>
    <v-container fill-width>
      <v-row>
        <v-spacer></v-spacer>
        <v-col cols="7">
          <v-tabs v-model="tab" centered>
            <v-tabs-slider></v-tabs-slider>
            <v-tab href="#tab-1"> Nodes </v-tab>
            <v-tab href="#tab-2"> Replicas </v-tab>
          </v-tabs>
          <v-tabs-items v-model="tab">
            <v-tab-item value="tab-1">
              <v-data-table
                :headers="nodeHeaders"
                :items="nodes"
              ></v-data-table>
            </v-tab-item>
            <v-tab-item value="tab-2">
              <v-data-table
                :headers="replicaHeaders"
                :items="replicas"
                :loading="loading"
              >
                <template v-slot:item.delete="{ item }">
                  <v-btn
                    text
                    color="primary"
                    @click="delReplica(item.name)"
                    :disabled="loading"
                  >
                    delete
                  </v-btn>
                </template>
              </v-data-table>
            </v-tab-item>
          </v-tabs-items>
        </v-col>
        <v-spacer></v-spacer>
      </v-row>
      <!-- <v-row>
        <v-col cols="7">
          <v-data-table :headers="nodeHeaders" :items="nodes"></v-data-table>
        </v-col>
        <v-col cols="5">
          <v-data-table :headers="replicaHeaders" :items="replicas"></v-data-table>
        </v-col>
      </v-row> -->
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
import {
  ClusterInfoResult,
  ClusterNodeInfo,
  ClusterReplicaInfo,
  Result
} from "@/struct";

@Component
export default class Cluster extends Vue {
  private nodeHeaders = [
    { text: "name", value: "name" },
    { text: "os type", value: "osType" },
    { text: "distribution", value: "distribution" },
    { text: "kernel version", value: "kernel" },
    { text: "arch", value: "arch" }
  ];

  private replicaHeaders = [
    { text: "name", value: "name" },
    { text: "node", value: "node" },
    { text: "status", value: "status" },
    { text: "", value: "delete" }
  ];

  private tab: any = null;

  private nodes: ClusterNodeInfo[] = [];
  private replicas: ClusterReplicaInfo[] = [];

  private infoText: string = "";
  private hasInfo: boolean = false;

  private loading: boolean = false;

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
        fetchPolicy: "no-cache"
      });
      if (res.errors) throw res.errors.map(v => v.message).join(",");
      if (res.data!.clusterInfo.message) throw res.data!.clusterInfo.message;
      this.nodes = res.data!.clusterInfo.nodes;
      this.replicas = res.data!.clusterInfo.replicas;
    } catch (e) {
      this.infoText = e.toString();
      this.hasInfo = true;
    }
  }

  async delReplica(name: string) {
    this.loading = true;
    try {
      let res = await this.$apollo.mutate<
        { deleteReplica: Result },
        { input: string }
      >({
        mutation: gql`
          mutation($input: String!) {
            deleteReplica(input: $input) {
              message
            }
          }
        `,
        variables: {
          input: name
        }
      });
      if (res.errors) throw res.errors.map(v => v.message).join(",");
      if (res.data!.deleteReplica.message)
        throw res.data!.deleteReplica.message;
      this.loading = false;
      this.infoText = "delete success";
      this.hasInfo = true;
      await this.loadData();
    } catch (e) {
      this.loading = false;
      this.infoText = e.toString();
      this.hasInfo = true;
    }
  }
}
</script>

<style lang="scss" scoped></style>
