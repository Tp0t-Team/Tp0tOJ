<template>
  <v-container class="bulletin-list">
    <v-card>
      <span class="bulletin-time ma-4">2019年9月30日20:21:12</span>
      <v-card-title>the title</v-card-title>
      <v-card-text>123456</v-card-text>
    </v-card>
  </v-container>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";

@Component
export default class Bulletin extends Vue {}
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