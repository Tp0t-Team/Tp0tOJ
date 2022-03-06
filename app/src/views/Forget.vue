<template>
  <v-container fill-height align-center justify-center>
    <v-card width="500">
      <v-form v-model="valid" class="pa-6" ref="form">
        <v-text-field
          v-model="mail"
          label="E-mail"
          :rules="[rules.required]"
          :disabled="loading"
        ></v-text-field>
        <v-spacer class="ma-6"></v-spacer>
        <v-btn
          color="primary"
          absolute
          bottom
          right
          @click="forget"
          :disabled="loading"
          :loading="loading"
        >submit</v-btn>
      </v-form>
    </v-card>
    <v-snackbar v-model="hasInfo" right bottom :timeout="3000">
      {{ infoText }}
      <v-spacer></v-spacer>
      <v-btn icon>
        <v-icon @click="hasInfo = false">close</v-icon>
      </v-btn>
    </v-snackbar>
  </v-container>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import gql from "graphql-tag";
import { ForgetResult } from "../struct";

@Component
export default class Forget extends Vue {
  private valid: boolean = false;

  private mail: string = "";

  private loading: boolean = false;

  private infoText: string = "";
  private hasInfo: boolean = false;

  private rules = {
    required: (value: string) => !!value || "请填写"
  };

  async forget() {
    if (!this.valid) {
      (this.$refs.form as any).validate();
      return;
    }
    this.loading = true;
    try {
      let res = await this.$apollo.mutate<ForgetResult, { input: string }>({
        mutation: gql`
          mutation($input: String!) {
            forget(input: $input) {
              message
            }
          }
        `,
        variables: {
          input: this.mail
        }
      });
      if (res.errors) throw res.errors.map(v => v.message).join(",");
      if (res.data!.forget.message) throw res.data!.forget.message;
      this.loading = false;
      this.infoText = "已经发送邮件至邮箱";
      this.hasInfo = true;
    } catch (e) {
      this.loading = false;
      this.infoText = e.toString();
      this.hasInfo = true;
    }
  }
}
</script>

<style lang="scss" scoped>
</style>