<template>
  <v-container fluid>
    <v-row>
      <v-col cols="6">
        <v-card>
          <v-toolbar dense>web</v-toolbar>
          <v-list dense>
            <v-list-item @click=";">
              <v-list-item-content>123</v-list-item-content>
            </v-list-item>
            <v-list-item @click=";">
              <v-list-item-content>123</v-list-item-content>
            </v-list-item>
            <v-list-item @click=";">
              <v-list-item-content>123</v-list-item-content>
            </v-list-item>
            <v-list-item @click=";">
              <v-list-item-content>123</v-list-item-content>
            </v-list-item>
          </v-list>
        </v-card>
      </v-col>
      <v-col cols="6">
        <v-form>
          <v-text-field disabled value="web"></v-text-field>
        </v-form>
      </v-col>
    </v-row>
  </v-container>
</template>

<style lang="scss" scoped>
</style>