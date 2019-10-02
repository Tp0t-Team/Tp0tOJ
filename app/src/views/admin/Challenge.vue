<template>
  <v-container fluid class="fill-height">
    <v-row class="fill-height">
      <v-col cols="6" class="content-col">
        <v-card>
          <v-toolbar dense>web</v-toolbar>
          <v-list dense>
            <v-list-item @click=";" :disabled="loading">
              <v-list-item-content>123</v-list-item-content>
              <v-list-item-icon>
                <v-btn icon :disabled="loading" @click="showDialog=true">
                  <v-icon>close</v-icon>
                </v-btn>
              </v-list-item-icon>
            </v-list-item>
            <v-list-item @click=";">
              <v-layout row>
                <v-spacer></v-spacer>
                <v-icon>add</v-icon>
                <v-spacer></v-spacer>
              </v-layout>
            </v-list-item>
          </v-list>
        </v-card>
      </v-col>
      <v-dialog v-model="showDialog" width="300px">
        <v-card>
          <v-card-title>Are you sure to delete this challenge?</v-card-title>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn text>cancel</v-btn>
            <v-btn text color="primary">accept</v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
      <v-col cols="6" class="content-col">
        <v-form>
          <v-row>
            <v-col cols="12">
              <v-file-input accept=".yaml, .yml" outlined label="load config file"></v-file-input>
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="6">
              <v-text-field outlined label="name" :disabled="loading"></v-text-field>
            </v-col>
            <v-col cols="6">
              <v-text-field outlined label="type" readonly value="web" :disabled="loading"></v-text-field>
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="6">
              <v-text-field outlined label="score" type="number" :disabled="loading"></v-text-field>
            </v-col>
            <v-col cols="6">
              <v-text-field outlined label="flag" :disabled="loading"></v-text-field>
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12">
              <v-textarea outlined label="description" :disabled="loading"></v-textarea>
            </v-col>
            <!--  -->
            <v-col cols="12">
              <v-card>
                <v-list dense>
                  <v-list-item>
                    <v-text-field
                      label="external link"
                      append-icon="add"
                      @click:append=";"
                      :disabled="loading"
                    ></v-text-field>
                  </v-list-item>
                  <v-list-item @click=";">
                    <v-list-item-content>http://www.google.com/</v-list-item-content>
                    <v-list-item-icon>
                      <v-btn icon :disabled="loading">
                        <v-icon>close</v-icon>
                      </v-btn>
                    </v-list-item-icon>
                  </v-list-item>
                </v-list>
              </v-card>
            </v-col>
            <!--  -->
            <v-col cols="12">
              <v-card>
                <v-list dense>
                  <v-list-item>
                    <v-text-field
                      label="hint"
                      append-icon="add"
                      @click:append=";"
                      :disabled="loading"
                    ></v-text-field>
                  </v-list-item>
                  <v-list-item @click=";">
                    <v-list-item-content>http://www.google.com/</v-list-item-content>
                    <v-list-item-icon>
                      <v-btn icon :disabled="loading">
                        <v-icon>close</v-icon>
                      </v-btn>
                    </v-list-item-icon>
                  </v-list-item>
                </v-list>
              </v-card>
            </v-col>
          </v-row>
        </v-form>
        <v-btn
          fab
          absolute
          right
          bottom
          color="primary"
          :loading="loading"
          :disabled="loading"
          @click="loading=true"
        >
          <v-icon>done</v-icon>
        </v-btn>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import yaml from "js-yaml";

@Component
export default class Challenge extends Vue {
  private showDialog: boolean = false;
  private loading: boolean = false;
}
</script>

<style lang="scss" scoped>
.content-col {
  height: calc(100vh - 120px);
  overflow-y: auto;
}
</style>