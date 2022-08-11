<template>
  <img v-if="hasImage" :src="`${url}&s=${size}`" @error="setText" />
  <span v-else>{{ name[0] }}</span>
</template>

<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";

@Component
export default class UserAvatar extends Vue {
  @Prop() url!: string;
  @Prop() size!: number;
  @Prop() name!: string;
  private hasImage: boolean = true;

  setText() {
    this.hasImage = false;
  }
}
</script>

<style lang="scss" scoped></style>
