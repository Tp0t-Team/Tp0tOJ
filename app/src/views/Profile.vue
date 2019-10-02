<template>
  <v-container>
    <!-- {{$route.params.user_id}} -->
    <v-row class="avatar">
      <v-spacer></v-spacer>
      <v-avatar color="blue" size="64">
        <span class="white--text">吴</span>
      </v-avatar>
      <v-spacer></v-spacer>
    </v-row>
    <v-row class="mt-5">
      <v-spacer></v-spacer>
      <v-col cols="6">
        <v-card class="outer">
          <v-card class="inner pa-4">
            <v-form>
              <v-row>
                <v-col cols="6">
                  <v-text-field readonly label="name"></v-text-field>
                </v-col>
                <v-col cols="6">
                  <v-text-field readonly label="role"></v-text-field>
                </v-col>
              </v-row>
              <v-row>
                <v-col cols="6">
                  <v-text-field readonly label="student number"></v-text-field>
                </v-col>
                <v-col cols="6">
                  <v-text-field readonly label="department"></v-text-field>
                </v-col>
              </v-row>
              <v-row>
                <v-col cols="6">
                  <v-text-field readonly label="score"></v-text-field>
                </v-col>
                <v-col cols="6">
                  <v-text-field readonly label="state"></v-text-field>
                </v-col>
              </v-row>
              <v-row>
                <v-col cols="6">
                  <v-text-field readonly label="grade"></v-text-field>
                </v-col>
                <v-col cols="6">
                  <v-text-field readonly label="protected time"></v-text-field>
                </v-col>
              </v-row>
              <v-row>
                <v-col cols="6">
                  <v-text-field readonly label="QQ"></v-text-field>
                </v-col>
                <v-col cols="6">
                  <v-text-field readonly label="mail"></v-text-field>
                </v-col>
              </v-row>
              <v-row>
                <v-col cols="6">
                  <v-text-field readonly label="top rank"></v-text-field>
                </v-col>
                <v-col cols="6">
                  <v-text-field readonly label="join time"></v-text-field>
                </v-col>
              </v-row>
            </v-form>
          </v-card>
        </v-card>
      </v-col>
      <v-spacer></v-spacer>
    </v-row>
  </v-container>
</template>

<style lang="scss" scoped>
.avatar {
  position: absolute;
  width: 100%;
  z-index: 1;
}

.outer {
  background: transparent;
  box-shadow: 0 0 4px rgba(0, 0, 0, 0.4);
}

.inner {
  background: transparent;
  box-shadow: 0 0 4px rgba(0, 0, 0, 0.4) inset;
}
</style>