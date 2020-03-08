<template>
  <div class="content" v-if="allow">
    <v-form ref="wrtieupForm">
      <v-file-input
        name="writeup"
        hide-details
        label="upload writeup"
        accept="application/pdf,.doc,.docx,application/msword,application/vnd.openxmlformats-officedocument.wordprocessingml.document,.md,text/plain"
        @change="fileUpload"
        :loading="loading"
        :disabled="loading"
      ></v-file-input>
    </v-form>
    <v-snackbar v-model="hasInfo" right top :timeout="3000">
      {{ infoText }}
      <v-spacer></v-spacer>
      <v-btn icon>
        <v-icon @click="hasInfo = false">close</v-icon>
      </v-btn>
    </v-snackbar>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";

@Component
export default class UserEditor extends Vue {
  private infoText: string = "";
  private hasInfo: boolean = false;

  private loading: boolean = false;
  private allow: boolean = false;

  mounted() {
    let allowTime: Date = new Date();
    allowTime.setHours(1);
    if (
      this.$store.state.competition.endtime != null &&
      this.$store.state.competition.competition === true
    ) {
      let deltaTime =
        Date.now() - this.$store.state.competition.endtime.getTime();
      if (deltaTime > 0 && deltaTime <= allowTime.getTime()) {
        this.allow = true;
      }
    }
  }

  async fileUpload(event: File) {
    if (event != null) {
      this.loading = true;
      try {
        let formData = new FormData();
        formData.append("wrtieup", event);
        let res = await fetch("/writeup", {
          method: "POST",
          headers: {
            "content-type": "multipart/form-data"
          },
          body: formData,
          cache: "no-cache"
        });
        let result = res.status;
        if (result != 200) throw `${res.status} ${res.statusText}`;
        this.infoText = "success";
        this.hasInfo = true;
        this.loading = false;
      } catch (e) {
        this.infoText = e.toString();
        this.hasInfo = true;
        this.loading = false;
      }
      (this.$refs.wrtieupForm as any).reset();
    }
  }
}
</script>

<style lang="scss" scoped>
.content {
  width: 180px;
}
</style>
