<template>
  <v-container fluid>
    <div>
      <div class="type-title display-1 ml-4">web</div>
      <v-divider color="primary"></v-divider>
      <v-row>
        <v-col sm="3">
          <v-layout justify-center>
            <v-hover disabled v-slot:default="{ hover }">
              <v-card
                :elevation="hover ? 12 : 2"
                width="20vw"
                height="220px"
                class="ma-4"
                @click="showDialog=true"
              >
                <v-badge left overlap>
                  <template v-slot:badge>1000</template>
                  <v-card-title class="subtitle-1 text-truncate">do you konw J language</v-card-title>
                </v-badge>
                <v-divider color="primary" class="divider"></v-divider>
                <v-card-text class="description">description</v-card-text>
                <v-overlay absolute :value="true" color="green" z-index="0">
                  <v-icon>done</v-icon>
                </v-overlay>
                <v-card-actions>
                  <v-row>
                    <v-col cols="4">
                      <v-layout row>
                        <v-spacer></v-spacer>
                        <v-btn fab @click.stop=";">
                          <v-icon>whatshot</v-icon>
                        </v-btn>
                        <v-spacer></v-spacer>
                      </v-layout>
                    </v-col>
                    <v-col cols="4">
                      <v-layout row>
                        <v-spacer></v-spacer>
                        <v-btn fab @click.stop=";">
                          <v-icon>whatshot</v-icon>
                        </v-btn>
                        <v-spacer></v-spacer>
                      </v-layout>
                    </v-col>
                    <v-col cols="4">
                      <v-layout row>
                        <v-spacer></v-spacer>
                        <v-btn fab @click.stop=";">
                          <v-icon>whatshot</v-icon>
                        </v-btn>
                        <v-spacer></v-spacer>
                      </v-layout>
                    </v-col>
                  </v-row>
                </v-card-actions>
              </v-card>
            </v-hover>
          </v-layout>
        </v-col>
      </v-row>
    </div>
    <v-dialog v-model="showDialog" :persistent="loading" width="800px">
      <v-card width="800px" height="600px">
        <v-toolbar dense>
          <v-spacer></v-spacer>
          <v-toolbar-title>
            do you konw J language
            <v-chip class="ml-2">1000pt</v-chip>
          </v-toolbar-title>
          <v-spacer></v-spacer>
        </v-toolbar>
        <v-text-field
          outlined
          class="ma-4 dialog-flag"
          label="flag"
          append-icon="send"
          :disabled="loading"
          :loading="loading"
          :rules="[rules.required]"
          @click:append="loading=true"
        ></v-text-field>
        <div class="dialog-discription pa-6">
          qqqqqqqqqq
          <br />qqqqqqqqqq
          <br />qqqqqqqqqq
          <br />qqqqqqqqqq
          <br />qqqqqqqqqq
          <br />qqqqqqqqqq
          <br />qqqqqqqqqq
          <br />qqqqqqqqqq
          <br />qqqqqqqqqq
          <br />qqqqqqqqqq
          <br />qqqqqqqqqq
          <br />qqqqqqqqqq
          <br />qqqqqqqqqq
          <br />qqqqqqqqqq
          <br />qqqqqqqqqq
          <br />qqqqqqqqqq
          <br />qqqqqqqqqq
          <br />qqqqqqqqqq
          <br />qqqqqqqqqq
          <br />qqqqqqqqqq
          <br />qqqqqqqqqq
          <br />
        </div>
        <div class="url-list">
          <v-chip color="primary" label outlined class="ma-4">
            http://www.google.com/
            <v-btn class="ml-4" icon color="primary" @click="openUrl('http://www.google.com/')">
              <v-icon>navigate_next</v-icon>
            </v-btn>
          </v-chip>
        </div>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";

@Component
export default class Challenge extends Vue {
  private showDialog: boolean = false;
  private loading: boolean = false;

  private rules = {
    required: (value: string) => !!value || "请填写"
  };

  openUrl(url: string) {
    window.open(url);
  }
}
</script>

<style lang="scss" scoped>
.type-title {
  width: 100%;
}

.divider {
  position: relative;
  left: 5%;
  width: 90%;
}

.description {
  height: 120px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 5;
  -webkit-box-orient: vertical;
  display: -moz-box;
  -moz-line-clamp: 5;
  -moz-box-orient: vertical;
}

.dialog-discription {
  height: 300px;
  overflow-y: auto;
}

.url-list {
  height: 128px;
  overflow-y: auto;
}
</style>