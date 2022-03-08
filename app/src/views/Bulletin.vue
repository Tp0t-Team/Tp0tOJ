<template>
  <v-container class="bulletin-list">
    <v-card v-for="item in bulletins" :key="item.time">
      <span
        class="bulletin-time ma-4"
      >{{new Date(item.publishTime).toLocaleString()}}</span>
      <v-card-title>{{item.title}}</v-card-title>
      <v-card-text>{{item.content}}</v-card-text>
    </v-card>
    <v-btn
      v-if="$store.state.global.role=='admin'||$store.state.global.role=='team'"
      fab
      absolute
      right
      bottom
      color="light-blue"
      @click="enterEdit"
    >
      <v-icon>add</v-icon>
    </v-btn>
    <v-dialog :persistent="loading" v-model="edit" width="400px">
      <v-card width="400px" class="pa-4">
        <v-form v-model="valid" ref="edit">
          <v-text-field :disabled="loading" v-model="title" label="Title" :rules="[rules.required]"></v-text-field>
          <v-textarea
            :disabled="loading"
            v-model="description"
            outlined
            label="Description"
            :rules="[rules.required]"
          ></v-textarea>
          <v-row>
            <v-spacer></v-spacer>
            <v-btn
              :disabled="loading"
              :loading="loading"
              color="primary"
              text
              @click="publish"
            >publish</v-btn>
          </v-row>
        </v-form>
      </v-card>
    </v-dialog>
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
import { Component, Vue } from "vue-property-decorator";
import gql from "graphql-tag";
import {
  BulletinPubInput,
  BulletinPubResult,
  AllBulletinResult,
  BulletinItem,
  BulletinSubResult
} from "@/struct";

@Component
export default class Bulletin extends Vue {
  private edit: boolean = false;
  private title: string = "";
  private description: string = "";

  private bulletins: BulletinItem[] = [];

  private valid: boolean = false;
  private loading: boolean = false;

  private rules = {
    required: (value: string) => !!value || "请填写"
  };

  private infoText: string = "";
  private hasInfo: boolean = false;

  async mounted() {
    await this.loadAll();
  }

  async loadAll() {
    try {
      let res = await this.$apollo.query<AllBulletinResult, {}>({
        query: gql`
          query {
            allBulletin {
              message
              bulletins {
                title
                content
                publishTime
              }
            }
          }
        `,
        fetchPolicy: "no-cache"
      });
      if (res.errors) throw res.errors.map(v => v.message).join(",");
      if (res.data!.allBulletin.message) throw res.data!.allBulletin.message;
      this.bulletins = res.data!.allBulletin.bulletins;
      // const observer = this.$apollo.subscribe<BulletinSubResult, {}>({
      //   query: gql`
      //     subscription {
      //       bulletin {
      //         title
      //         description
      //         time
      //       }
      //     }
      //   `
      // });
      // observer.subscribe({
      //   next: data => {
      //     this.bulletins = [data.data!.bulletin, ...this.bulletins];
      //   }
      // });
    } catch (e) {
      this.infoText = e.toString();
      this.hasInfo = true;
    }
  }

  enterEdit() {
    this.title = "";
    this.description = "";
    this.edit = true;
  }

  async publish() {
    if (!this.valid) {
      (this.$refs.edit as any).validate();
      return;
    }
    this.loading = true;
    try {
      let res = await this.$apollo.mutate<BulletinPubResult, BulletinPubInput>({
        mutation: gql`
          mutation($input: BulletinPubInput!) {
            bulletinPub(input: $input) {
              message
            }
          }
        `,
        variables: {
          input: {
            title: this.title,
            content: this.description,
            topping: false
          }
        }
      });
      if (res.errors) throw res.errors.map(v => v.message).join(",");
      if (res.data!.bulletinPub.message) throw res.data!.bulletinPub.message;
      this.loading = false;
      this.edit = false;
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
.bulletin-list {
  max-width: 80%;
}

.bulletin-time {
  position: absolute;
  top: 0;
  right: 0;
  font-size: 12px;
  opacity: 0.4;
}
</style>