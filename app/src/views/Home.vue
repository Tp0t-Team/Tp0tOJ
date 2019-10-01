<template>
  <home-frame :isDark="$vuetify.theme.dark"></home-frame>
</template>

<script lang="ts">
import { Component, Vue, Prop, Watch } from "vue-property-decorator";
import HomeFrame from "@/components/HomeFrame.vue";

@Component({
  components: {
    HomeFrame
  }
})
export default class Home extends Vue {}
</script>

<style lang="scss" scoped>
</style>