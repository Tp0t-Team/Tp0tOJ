<template>
  <v-container fluid class="fill-height">
    <iframe ref="iframe" width="100%" height="100%" src="/home.html"></iframe>
  </v-container>
</template>

<script lang="ts">
import { Component, Vue, Prop, Watch } from "vue-property-decorator";

@Component
export default class HomeFrame extends Vue {
  @Prop() isDark!: boolean;

  mounted() {
    (this.$refs.iframe as HTMLIFrameElement).addEventListener("load", () => {
      this.setDark();
    });
  }

  @Watch("isDark", { immediate: true })
  setDark() {
    if (this.$refs.iframe == null) return;
    (this.$refs
      .iframe as HTMLIFrameElement).contentDocument!.body.style.color = this
      .isDark
      ? "white"
      : "black";
  }
}
</script>

<style lang="scss" scoped>
iframe {
  border: 0;
}
</style>