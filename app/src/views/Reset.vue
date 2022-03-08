<template>
  <v-container fill-height align-center justify-center>
    <v-card width="500">
      <v-form v-model="valid" class="pa-6" ref="form">
        <v-text-field
          v-model="password"
          label="Password"
          :append-icon="showPassword ? 'visibility' : 'visibility_off'"
          :type="showPassword ? 'text' : 'password'"
          @click:append="showPassword = !showPassword"
          :rules="[rules.required]"
          :disabled="loading"
        ></v-text-field>
        <v-text-field
          v-model="repeat"
          label="Repeat"
          :append-icon="showPassword ? 'visibility' : 'visibility_off'"
          :type="showPassword ? 'text' : 'password'"
          @click:append="showPassword = !showPassword"
          :rules="[rules.required]"
          :disabled="loading"
          :error-messages="againError"
          @focus="againError = ''"
          @blur="check"
        ></v-text-field>
        <v-spacer class="ma-6"></v-spacer>
        <v-btn
          color="primary"
          absolute
          bottom
          right
          @click="reset"
          :disabled="loading"
          :loading="loading"
        >reset</v-btn>
      </v-form>
    </v-card>
    <v-snackbar v-model="hasInfo" right bottom :timeout="3000">
      {{ infoText }}
      <!-- <v-spacer></v-spacer> -->
      <template v-slot:action="{ attrs }">
        <v-btn icon>
          <v-icon v-bind="attrs" @click="hasInfo = false">close</v-icon>
        </v-btn>
      </template>
    </v-snackbar>
  </v-container>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import gql from "graphql-tag";
import { ResetResult, ResetInput } from "../struct";

@Component
export default class Reset extends Vue {
  private showPassword: boolean = false;
  private againError: string = "";

  private valid: boolean = false;

  private password: string = "";
  private repeat: string = "";

  private loading: boolean = false;

  private infoText: string = "";
  private hasInfo: boolean = false;

  private rules = {
    required: (value: string) => !!value || "请填写"
  };

  check() {
    if (this.password != this.repeat) this.againError = "密码不一致";
  }

  mounted() {
    if (!this.$route.query.token) {
      this.$router.replace("/");
    }
  }

  async reset() {
    if (!this.valid) {
      (this.$refs.form as any).validate();
      return;
    }
    this.loading = true;
    try {
      let res = await this.$apollo.mutate<ResetResult, ResetInput>({
        mutation: gql`
          mutation($input: ResetInput!) {
            reset(input: $input) {
              message
            }
          }
        `,
        variables: {
          input: {
            password: this.password,
            token: this.$route.query.token.toString()
          }
        }
      });
      if (res.errors) throw res.errors.map(v => v.message).join(",");
      if (res.data!.reset.message) throw res.data!.reset.message;
      this.loading = false;
      this.$router.replace("/login");
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