<template>
  <v-container fluid class="fill-height">
    <v-row class="fill-height">
      <v-spacer></v-spacer>
      <v-col cols="6" class="content-col">
        <!-- show-expand item-key="challengeId" single-expand -->
        <v-data-table
          v-model="selected"
          :headers="headers"
          :items="events"
          @click:row="select"
          show-select
          :loading="loading"
          item-key="challengeId"
        >
          <template v-slot:body.append="{ headers }">
            <td :colspan="headers.length">
              <div class="action-group">
                <div class="action-item">
                  <v-btn
                    text
                    tile
                    block
                    color="primary"
                    @click="newEvent"
                    :disabled="loading"
                    >new</v-btn
                  >
                </div>
                <div class="action-item">
                  <v-btn
                    text
                    tile
                    block
                    color="accent"
                    :disabled="selected.length == 0 || loading"
                    @click="deleteEvent"
                    >delete</v-btn
                  >
                </div>
              </div>
            </td>
          </template>
          <template v-slot:item.time="{ item }">
            {{ formatDate(item.time) }}
          </template>
          <template v-slot:item.action="{ item }">
            {{ actionItems[item.action] }}
          </template>
        </v-data-table>
      </v-col>
      <v-dialog v-model="showDiscardDialog" width="300px">
        <v-card>
          <v-card-title>Are you sure to discard changes?</v-card-title>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn text @click="showDiscardDialog = false">cancel</v-btn>
            <v-btn text color="primary" @click="continueChange">sure</v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
      <v-col cols="6" class="content-col" v-if="!withoutValue">
        <v-form>
          <v-row>
            <v-col cols="6">
              <v-select
                v-model="action"
                :items="actionItems"
                outlined
                hide-details
                label="action"
                :disabled="loading"
                @change="Changed"
                ></v-select>
            </v-col>
            <v-col cols="6">
              <v-menu top :close-on-content-click="false">
                <template v-slot:activator="{ on, attrs }">
                  <v-text-field :value="time.toLocaleString()" outlined readonly label="time" v-bind="attrs" v-on="on"></v-text-field>
                </template>
                <v-card width="300px">
                  <v-tabs v-model="tab">
                    <v-tabs-slider color="yellow"></v-tabs-slider>
                    <v-tab key="date">date</v-tab>
                    <v-tab key="time">time</v-tab>
                  </v-tabs>
                  <v-tabs-items v-model="tab">
                    <v-tab-item key="date">
                      <v-date-picker full-width v-model="timeDate"></v-date-picker>
                    </v-tab-item>
                    <v-tab-item key="time">
                      <v-time-picker full-width format="ampm" v-model="timeTime"></v-time-picker>
                    </v-tab-item>
                  </v-tabs-items>
                </v-card>
              </v-menu>
            </v-col>
          </v-row>
          <v-btn
            fab
            absolute
            right
            bottom
            color="primary"
            :loading="loading"
            :disabled="loading || !changed"
            @click="submit"
          >
            <v-icon>done</v-icon>
          </v-btn>
        </v-form>
      </v-col>
      <v-spacer></v-spacer>
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
</template>

<script lang="ts">
import dayjs from "dayjs";
import { Component, Vue, Watch } from "vue-property-decorator";
import gql from "graphql-tag";
import {
    GameEventWithId,
    AddEventInput,
    AddEventResult,
    UpdateEventInput,
    UpdateEventResult,
    DeleteEventInput,
    DeleteEventResult,
    AllEventResult
} from "@/struct";
import constValue from "@/constValue";

@Component({})
export default class Challenge extends Vue {
  private headers = [
    { text: "time", value: "time" },
    { text: "action", value: "action" },
  ];

  private selected: GameEventWithId[] = [];

  private showDiscardDialog: boolean = false;
  private withoutValue: boolean = true;
  private loading: boolean = false;
  private changed: boolean = false;

  private events: GameEventWithId[] = [];
  private currentEvent: GameEventWithId | null = null;
  private tempEvent: GameEventWithId | null = null;

  private infoText: string = "";
  private hasInfo: boolean = false;

  private action: string = "";
  private actionItems = constValue.eventType;
  private time: Date = new Date();
  private timeDate: string = "";
  private timeTime: string = "";
  private tab: string = "date";

  formatDate(time: Date) {
    return dayjs(time).format("YYYY-MM-DD HH:mm (UTCZ)");
  }

  @Watch("currentEvent")
  editLoaded() {
    if (this.currentEvent !== null) {
      this.action = this.actionItems[this.currentEvent.action];
      this.time = new Date(this.currentEvent.time);
      let d = dayjs(this.time);
      this.timeDate = d.format("YYYY-MM-DD");
      this.timeTime = d.format("HH:mm");
    }
  }

  @Watch("timeDate")
  updateDate() {
    this.changed = true;
    this.time = dayjs(`${this.timeDate} ${this.timeTime}`).toDate();
  }

  @Watch("timeTime")
  updateTime() {
    this.changed = true;
    this.time = dayjs(`${this.timeDate} ${this.timeTime}`).toDate();
  }

  @Watch("action")
  updateAction() {
    this.changed = true;
  }

  async mounted() {
    await this.loadAll();
  }

  async loadAll() {
    try {
      let res = await this.$apollo.query<AllEventResult, {}>({
        query: gql`
          query {
            allEvents {
              message
              allEvents {
                eventId
                time
                action
              }
            }
          }
        `,
        fetchPolicy: "no-cache",
      });
      if (res.errors) throw res.errors.map((v) => v.message).join(",");
      if (res.data!.allEvents.message)
        throw res.data!.allEvents.message;
      this.events = res.data!.allEvents.allEvents.map((it)=> ({
        eventId: it.eventId,
        action: it.action,
        time: new Date(parseInt(it.time as string)*1000),//解决unix time精度问题
      }));
    } catch (e) {
      this.infoText = e.toString();
      this.hasInfo = true;
    }
  }

  error(error: string) {
    this.infoText = error;
    this.hasInfo = true;
  }

  Changed() {
    this.changed = true;
  }

  async submit() {
    this.loading = true;
    let aimEvent: GameEventWithId = {
      eventId: this.currentEvent!.eventId,
      action: this.actionItems.indexOf(this.action),
      time: dayjs(this.time).unix().toString()
    };
    if (aimEvent.eventId[0] == "-") {
      aimEvent.eventId = "";
    }
    try {
      if (aimEvent.eventId == "") {
        let res = await this.$apollo.mutate<AddEventResult, {input: AddEventInput}>({
          mutation: gql`
            mutation ($input: AddEventInput!) {
              addEventAction(input: $input) {
                message
              }
            }
          `,
          variables: {input: {
            time: aimEvent.time as string,
            action: aimEvent.action
          }},
        });
        if (res.errors) throw res.errors.map((v) => v.message).join(",");
        if (res.data!.addEventAction.message)
          throw res.data!.addEventAction.message;
      } else {
        let res = await this.$apollo.mutate<UpdateEventResult, {input: UpdateEventInput}>({
          mutation: gql`
            mutation ($input: UpdateEventInput!) {
              updateEvent(input: $input) {
                message
              }
            }
          `,
          variables: {input: {
            eventId: aimEvent.eventId,
            time: aimEvent.time as string
          }},
        });
        if (res.errors) throw res.errors.map((v) => v.message).join(",");
        if (res.data!.updateEvent.message)
          throw res.data!.updateEvent.message;
      }
      this.loading = false;
      this.changed = false;
      this.currentEvent = null;
      this.withoutValue = true;
      await this.loadAll();
      this.infoText = "add / update success";
      this.hasInfo = true;
    } catch (e) {
      this.loading = false;
      this.infoText = e.toString();
      this.hasInfo = true;
    }
  }

  select(event: GameEventWithId) {
    let aimEvent = JSON.parse(JSON.stringify(event));
    if (this.changed) {
        this.tempEvent = aimEvent;
        this.showDiscardDialog = true;
    } else {
        this.changed = false;
        this.withoutValue = false;
        this.currentEvent = aimEvent;
    }
  }

  newEvent() {
    let aimEvent: GameEventWithId = {
        eventId: "-" + Date.now().toLocaleString(),
        time: new Date(Date.now()),
        action: 0
    };

    if (this.changed) {
      this.tempEvent = aimEvent;
      this.showDiscardDialog = true;
    } else {
      this.changed = false;
      this.withoutValue = false;
      this.currentEvent = aimEvent;
    }
  }

  continueChange() {
    this.showDiscardDialog = false;
    this.changed = false;
    if (this.tempEvent == null) {
      this.withoutValue = true;
    }
    this.currentEvent = this.tempEvent;
  }

  async deleteEvent() {
    this.loading = true;
    try {
      let res = await this.$apollo.mutate<
        DeleteEventResult,
        { input: DeleteEventInput }
      >({
        mutation: gql`
          mutation ($input: DeleteEventInput!) {
            deleteEvent(input: $input) {
              message
            }
          }
        `,
        variables: {
          input: {
            eventIds: this.selected.map((it) => it.eventId),
          },
        },
      });
      if (res.errors) throw res.errors.map((v) => v.message).join(",");
      if (res.data!.deleteEvent.message)
        throw res.data!.deleteEvent.message;
      this.selected = [];
      this.loading = false;
      this.changed = false;
      this.currentEvent = null;
      this.withoutValue = true;
      this.infoText = 'delete success';
      this.hasInfo = true;
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
.content-col {
  height: calc(100vh - 120px);
  overflow-y: auto;
}

.action-group {
  display: flex;
  flex-direction: row;
}

.action-item {
  flex-basis: 50%;
}
</style>
