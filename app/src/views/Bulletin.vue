<template>
  <div class="content-col">
    <v-container fluid>
      <div class="content">
        <v-card
          :class="`ma-4 bulletin bulletin-${item.style}`"
          v-for="item in bulletins"
          :key="item.time"
        >
          <svg
            stroke-width="1.5"
            viewBox="0 0 24 24"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
            height="1em"
            width="1em"
          >
            <path
              d="M14.272 10.445 18 2m-8.684 8.632L5 2m7.761 8.048L8.835 2m5.525 0-1.04 2.5M6 16a6 6 0 1 0 12 0 6 6 0 0 0-12 0Z"
              stroke="currentColor"
              stroke-linecap="round"
              stroke-linejoin="round"
            />
            <text
              x="9"
              y="19.5"
              fill="white"
              style="font-size: 10px; font-weight: bold;"
            >
              {{
                item.style == "first" ? "1" : item.style == "second" ? "2" : "3"
              }}
            </text>
          </svg>
          <span class="bulletin-time ma-4">{{
            new Date(item.publishTime).toLocaleString()
          }}</span>
          <v-card-title>{{ item.title }}</v-card-title>
          <v-card-text
            ><pre class="pl-4">{{ item.content }}</pre></v-card-text
          >
        </v-card>
      </div>

      <v-btn
        v-if="$store.state.global.role == 'admin'"
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
            <v-text-field
              :disabled="loading"
              v-model="title"
              label="Title"
              :rules="[rules.required]"
            ></v-text-field>
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
                >publish</v-btn
              >
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
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import gql from "graphql-tag";
import {
  BulletinPubInput,
  BulletinPubResult,
  AllBulletinResult,
  BulletinItem
  // BulletinSubResult
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
                style
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
      res.data!.allBulletin.bulletins.sort(
        (a, b) =>
          Number(a.publishTime < b.publishTime) -
          Number(a.publishTime > b.publishTime)
      );
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
      if (e === "unauthorized") {
        this.$store.commit("global/resetUserIdAndRole");
        this.$router.push("/login?unauthorized");
        return;
      }
      this.infoText = e.toString();
      this.hasInfo = true;
    }
  }
}
</script>

<style lang="scss" scoped>
.content-col {
  height: calc(100vh - 96px);
  overflow-y: auto;
}

.content {
  margin: 0 auto;
  max-width: 80%;
}

.bulletin-time {
  position: absolute;
  top: 0;
  right: 0;
  font-size: 12px;
  opacity: 0.4;
}

.bulletin > svg {
  display: none;
  border: 1px solid transparent;
}

.bulletin-first > svg,
.bulletin-second > svg,
.bulletin-third > svg {
  display: flex;
  position: absolute;
  left: 8px;
  width: 32px;
  height: 100%;
  justify-content: center;
  align-items: center;
}

.bulletin-first,
.bulletin-second,
.bulletin-third {
  padding-left: 32px;
}

.bulletin-first {
  border: 1px solid rgba(255, 193, 7, 0.5);
}

.bulletin-second {
  border: 1px solid rgba(3, 169, 244, 0.5);
}

.bulletin-third {
  border: 1px solid rgba(76, 175, 80, 0.5);
}

.bulletin-first > svg {
  color: #ffc107;
}

.bulletin-first > svg > text {
  fill: #ffc107;
}

.bulletin-second > svg {
  color: #03a9f4;
}

.bulletin-second > svg > text {
  fill: #03a9f4;
}

.bulletin-third > svg {
  color: #4caf50;
}

.bulletin-third > svg > text {
  fill: #4caf50;
}
</style>
