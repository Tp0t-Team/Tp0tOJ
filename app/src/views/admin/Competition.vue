<template>
  <v-container fluid class="fill-height">
    <v-row class="fill-height">
      <v-col cols="12" class="content-col">
        <v-row>
          <v-col cols="2"></v-col>
          <v-col cols="8">
            <v-card class="pa-4">
              <v-form>
                <v-row>
                  <v-col cols="6">
                    <v-row>
                      <v-spacer></v-spacer>
                      <v-switch
                        v-model="competition"
                        label="competition mode"
                        :disabled="loading"
                      ></v-switch>
                      <v-spacer></v-spacer>
                    </v-row>
                  </v-col>
                  <v-col cols="6">
                    <v-row>
                      <v-spacer></v-spacer>
                      <v-switch
                        v-model="registerAllow"
                        label="allow register"
                        :disabled="loading"
                      ></v-switch>
                      <v-spacer></v-spacer>
                    </v-row>
                  </v-col>
                </v-row>
                <v-row>
                  <v-col>
                    <date-time-picker
                      v-model="beginTime"
                      label="begin time"
                      :disabled="loading"
                    ></date-time-picker>
                  </v-col>
                  <v-col>
                    <date-time-picker
                      v-model="endTime"
                      label="end time"
                      :disabled="loading"
                    ></date-time-picker>
                  </v-col>
                </v-row>
                <v-spacer class="pt-4"></v-spacer>
                <v-btn fab absolute right bottom color="primary" @click="done">
                  <v-icon>done</v-icon>
                </v-btn>
              </v-form>
            </v-card>
          </v-col>
        </v-row>
      </v-col>
    </v-row>
    <v-snackbar v-model="hasInfo" right bottom :timeout="3000">
      {{ infoText }}
      <v-spacer></v-spacer>
      <v-btn icon>
        <v-icon @click="hasInfo = false">close</v-icon>
      </v-btn>
    </v-snackbar>
  </v-container>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import DateTimePicker from "@/components/DateTimePicker.vue";
import {
  CompetitionResult,
  CompetitionMutationResult,
  CompetitionMutationInput
} from "../../struct";
import gql from "graphql-tag";

@Component({
  components: {
    DateTimePicker
  }
})
export default class Competition extends Vue {
  private competition: boolean = false;
  private registerAllow: boolean = false;
  private beginTime: Date | null = null;
  private endTime: Date | null = null;

  private loading: boolean = false;

  private infoText: string = "";
  private hasInfo: boolean = false;

  async mounted() {
    this.loading = true;
    try {
      let res = await this.$apollo.query<CompetitionResult>({
        query: gql`
          query {
            competition {
              message
              competition
              registerAllow
              beginTime
              endTime
            }
          }
        `
      });
      if (res.errors) throw res.errors.map(v => v.message).join(",");
      if (res.data!.competition.message) throw res.data!.competition.message;
      this.competition = res.data!.competition.competition;
      this.registerAllow = res.data!.competition.registerAllow;
      this.beginTime = new Date(res.data!.competition.beginTime);
      this.endTime = new Date(res.data!.competition.endTime);
      this.loading = false;
    } catch (e) {
      this.loading = false;
      this.infoText = e.toString();
      this.hasInfo = true;
    }
  }

  async done() {
    this.loading = true;
    try {
      let res = await this.$apollo.mutate<
        CompetitionMutationResult,
        CompetitionMutationInput
      >({
        mutation: gql`
          mutation($input: CompetitionMutationIEnput!) {
            competition(input: $input) {
              message
            }
          }
        `,
        variables: {
          input: {
            competition: this.competition,
            registerAllow: this.registerAllow,
            beginTime: this.beginTime!.toUTCString(),
            endTime: this.endTime!.toUTCString()
          }
        }
      });
      if (res.errors) throw res.errors.map(v => v.message).join(",");
      if (res.data!.competition.message) throw res.data!.competition.message;
      this.loading = false;
      this.infoText = "success";
      this.hasInfo = true;
    } catch (e) {
      this.loading = false;
      this.infoText = e.toString();
      this.hasInfo = true;
    }
  }
}
</script>

<style lang="scss" scoped>
.content-col {
  height: calc(100vh - 120px);
  overflow-y: auto;
}
</style>
