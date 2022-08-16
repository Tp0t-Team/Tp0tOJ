<template>
  <v-container fluid fill-height style="padding: 0;">
    <iframe class="sanddance" src="/sanddance.html" ref="sanddance"></iframe>
    <v-snackbar v-model="hasInfo" right bottom :timeout="3000">
      {{ infoText }}
      <!-- <v-spacer></v-spacer> -->
      <template v-slot:action="{ attrs }">
        <v-btn icon>
          <v-icon v-bind="attrs" @click="hasInfo = false">close</v-icon>
        </v-btn>
      </template>
    </v-snackbar>
    <v-btn
      fab
      absolute
      right
      bottom
      color="primary"
      :loading="loading"
      :disable="loading"
      @click="loadData"
    >
      <v-icon>refresh</v-icon>
    </v-btn>
  </v-container>
</template>

<script lang="ts">
import { Component, Vue, Watch } from "vue-property-decorator";

@Component({})
export default class Analyse extends Vue {
  private infoText: string = "";
  private hasInfo: boolean = false;
  private loading: boolean = false;
  private data: any = {};

  async mounted() {
    (this.$refs.sanddance as HTMLIFrameElement).addEventListener("load", () => {
      this.update();
    });
    await this.loadData();
  }

  async loadData() {
    this.loading = true;
    try {
      let res = await fetch("/data");
      let data = await res.json();
      this.data = data;
      this.update();
    } catch (e) {
      this.infoText = e.toString();
      this.hasInfo = true;
    }
    this.loading = false;
  }

  @Watch("$vuetify.theme.dark")
  update() {
    let theme = this.$vuetify.theme.dark ? "dark-theme" : "";
    (this.$refs.sanddance as HTMLIFrameElement).contentWindow?.postMessage(
      { theme: theme, data: this.data },
      "*"
    );
  }
}
</script>

<style lang="scss" scoped>
.sanddance {
  width: 100%;
  height: 100%;
  border: 0;
}
</style>
