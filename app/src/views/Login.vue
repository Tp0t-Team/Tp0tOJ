<template>
  <v-container fill-height align-center justify-center>
    <v-card width="500">
      <v-tabs centered grow color="primary">
        <v-tab href="#tab-login" :disabled="loading">login</v-tab>
        <v-tab href="#tab-register" :disabled="loading">register</v-tab>
        <v-tab-item value="tab-login">
          <v-form v-model="loginValid" class="pa-6" ref="loginForm">
            <v-text-field
              v-model="username"
              label="User"
              :rules="[rules.required]"
              :disabled="loading"
            ></v-text-field>
            <v-text-field
              v-model="password"
              label="Password"
              :append-icon="showPassword ? 'visibility' : 'visibility_off'"
              :type="showPassword ? 'text' : 'password'"
              @click:append="showPassword = !showPassword"
              :rules="[rules.required]"
              :disabled="loading"
            ></v-text-field>
            <v-spacer class="ma-6"></v-spacer>
            <v-btn
              color="primary"
              absolute
              bottom
              right
              @click="login"
              :disabled="loading"
              :loading="loading"
            >login</v-btn>
          </v-form>
        </v-tab-item>
        <v-tab-item value="tab-register">
          <v-form v-model="regValid" class="pa-6" ref="registerForm">
            <v-layout row>
              <v-flex sm6 pl-3 pr-3>
                <v-text-field
                  v-model="name"
                  label="Name"
                  :rules="[rules.required]"
                  :disabled="loading"
                ></v-text-field>
                <v-text-field
                  v-model="department"
                  label="Department"
                  :rules="[rules.required]"
                  :disabled="loading"
                ></v-text-field>
                <v-text-field v-model="qq" label="QQ" :rules="[rules.required]" :disabled="loading"></v-text-field>
                <v-text-field
                  v-model="regPassword"
                  label="Password"
                  :append-icon="showPassword ? 'visibility' : 'visibility_off'"
                  :type="showPassword ? 'text' : 'password'"
                  @click:append="showPassword = !showPassword"
                  :rules="[rules.required]"
                  :disabled="loading"
                ></v-text-field>
              </v-flex>
              <v-flex sm6 pl-3 pr-3>
                <v-text-field
                  v-model="studentNumber"
                  label="Student Number"
                  :rules="[rules.required]"
                  :disabled="loading"
                ></v-text-field>
                <v-text-field
                  v-model="grade"
                  label="Grade"
                  :rules="[rules.required]"
                  :disabled="loading"
                ></v-text-field>
                <v-text-field
                  v-model="mail"
                  label="Mail"
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
              </v-flex>
              <v-spacer class="ma-3"></v-spacer>
              <v-btn
                color="primary"
                absolute
                bottom
                right
                @click="register"
                :disabled="loading"
                :loading="loading"
              >register</v-btn>
            </v-layout>
          </v-form>
        </v-tab-item>
      </v-tabs>
    </v-card>
  </v-container>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";

@Component
export default class Login extends Vue {
  private username: string = "";
  private password: string = "";

  private name: string = "";
  private studentNumber: string = "";
  private department: string = "";
  private grade: string = "";
  private qq: string = "";
  private mail: string = "";
  private regPassword: string = "";
  private repeat: string = "";

  private loginValid: boolean = false;
  private regValid: boolean = false;

  private showPassword: boolean = false;
  private againError: string = "";
  private rules = {
    required: (value: string) => !!value || "请填写"
  };

  private loading: boolean = false;

  check() {
    if (this.regPassword != this.repeat) this.againError = "密码不一致";
  }

  async login() {
    if (!this.loginValid) {
      (this.$refs.loginForm as any).validate();
      return;
    }
    this.loading = true;
    // do the action
  }

  async register() {
    if (!this.regValid) {
      (this.$refs.registerForm as any).validate();
      return;
    }
    this.loading = true;
    // do the action
  }
}
</script>

<style lang="scss" scoped>
</style>