<template>
  <v-container fluid>
    <v-row>
      <v-col cols="6" class="content-col">
        <v-data-table :headers="headers"></v-data-table>
      </v-col>
      <v-col cols="6" class="content-col">
        <v-divider class="avatar-divider" color="primary"></v-divider>
        <v-row>
          <v-spacer></v-spacer>
          <v-avatar color="blue" size="64">
            <span>吴</span>
          </v-avatar>
          <v-spacer></v-spacer>
        </v-row>
        <v-form>
          <v-row>
            <v-col cols="6">
              <v-text-field outlined label="name"></v-text-field>
            </v-col>
            <v-col cols="6">
              <v-text-field outlined disabled label="student number"></v-text-field>
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="6">
              <v-select outlined label="role"></v-select>
            </v-col>
            <v-col cols="6">
              <v-select outlined label="state"></v-select>
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="6">
              <v-text-field outlined disabled type="number" label="score"></v-text-field>
            </v-col>
            <v-col cols="6">
              <v-menu ref="dateMenu" v-model="dateMenu" :close-on-content-click="false">
                <template v-slot:activator="{ on }">
                  <v-text-field outlined type="date" label="protected time" v-on="on"></v-text-field>
                </template>
                <v-date-picker v-model="menuDate">
                  <v-spacer></v-spacer>
                  <v-btn text @click="dateMenu=false">Cancel</v-btn>
                  <v-btn text color="primary" @click="$refs.dateMenu.save(menuDate)">OK</v-btn>
                </v-date-picker>
              </v-menu>
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="6">
              <v-select outlined label="department"></v-select>
            </v-col>
            <v-col cols="6">
              <v-select outlined label="grade"></v-select>
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="6">
              <v-text-field outlined label="QQ"></v-text-field>
            </v-col>
            <v-col cols="6">
              <v-text-field outlined label="mail"></v-text-field>
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="6">
              <v-text-field outlined disabled label="top rank"></v-text-field>
            </v-col>
            <v-col>
              <v-text-field outlined disabled type="date" label="join time"></v-text-field>
            </v-col>
          </v-row>
        </v-form>
        <v-toolbar dense>
          解题列表
          <v-spacer></v-spacer>
          <v-btn icon>
            <v-icon>refresh</v-icon>
          </v-btn>
        </v-toolbar>
        <v-simple-table dense>
          <thead>
            <tr>
              <th>submit time</th>
              <th>challenge name</th>
              <th>mark</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td>2019-10-01</td>
              <td>are you ok</td>
              <td>
                <v-badge class="mark-badge" color="yellow">
                  <template v-slot:badge>1</template>
                </v-badge>
              </td>
            </tr>
            <tr>
              <td>2019-10-01</td>
              <td>are you ok</td>
              <td>
                <v-badge class="mark-badge" color="blue">
                  <template v-slot:badge>2</template>
                </v-badge>
              </td>
            </tr>
            <tr>
              <td>2019-10-01</td>
              <td>are you ok</td>
              <td>
                <v-badge class="mark-badge" color="green">
                  <template v-slot:badge>3</template>
                </v-badge>
              </td>
            </tr>
          </tbody>
        </v-simple-table>
        <v-btn fab color="primary" absolute right bottom>
          <v-icon>done</v-icon>
        </v-btn>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";

@Component
export default class User extends Vue {
  private headers = [
    { text: "name", value: "name" },
    { text: "role", value: "role" },
    { text: "state", value: "state" },
    { text: "grade", value: "grade" },
    { text: "protected_time", value: "protected_time" },
    { text: "score", value: "score" }
  ];

  private dateMenu: boolean = false;
  private menuDate: Date = new Date();
}
</script>

<style lang="scss" scoped>
.content-col {
  height: calc(100vh - 120px);
  overflow-y: auto;
}

.avatar-divider {
  position: relative;
  bottom: -32px;
}

.mark-badge {
  position: relative;
  top: -4px;
}
</style>