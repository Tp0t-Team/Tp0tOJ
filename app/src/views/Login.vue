<template>
  <v-container fill-height align-center justify-center>
    <v-card width="500">
      <v-tabs v-model="tab" centered grow color="primary">
        <v-tab href="#tab-login" :disabled="loading">login</v-tab>
        <v-tab href="#tab-register" :disabled="loading">register</v-tab>
        <v-tab-item value="tab-login">
          <v-form v-model="loginValid" class="pa-6" ref="loginForm">
            <v-text-field
              v-model="mail"
              label="Mail"
              :rules="[rules.required, rules.email, rules.mailLen]"
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
              absolute
              bottom
              right
              text
              @click="$router.push('/forget')"
              class="forget"
              >forget</v-btn
            >
            <v-btn
              color="primary"
              absolute
              bottom
              right
              @click="login"
              :disabled="loading"
              :loading="loading"
              >login</v-btn
            >
          </v-form>
        </v-tab-item>
        <v-tab-item value="tab-register">
          <v-form v-model="regValid" class="pa-6" ref="registerForm">
            <v-layout row>
              <v-flex sm6 pl-3 pr-3>
                <!-- Zero space is used to disable autofill. -->
                <v-text-field
                  v-model="name"
                  :label="'\u200bName'"
                  :rules="[rules.required]"
                  :disabled="loading"
                ></v-text-field>
              </v-flex>
              <v-flex sm6 pl-3 pr-3>
                <v-text-field
                  v-model="mail"
                  label="Mail"
                  :rules="[rules.required, rules.email, rules.mailLen]"
                  :disabled="loading"
                ></v-text-field>
              </v-flex>
            </v-layout>
            <v-layout row>
              <v-flex sm6 pl-3 pr-3>
                <v-text-field
                  v-model="regPassword"
                  label="Password"
                  :append-icon="showPassword ? 'visibility' : 'visibility_off'"
                  :type="showPassword ? 'text' : 'password'"
                  @click:append="showPassword = !showPassword"
                  :rules="[
                    rules.required,
                    rules.passLen(8, 18),
                    rules.password
                  ]"
                  :disabled="loading"
                ></v-text-field>
              </v-flex>
              <v-flex sm6 pl-3 pr-3>
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
            </v-layout>
            <v-layout row>
              <v-spacer class="ma-3"></v-spacer>
              <v-btn
                color="primary"
                absolute
                bottom
                right
                @click="register"
                :disabled="loading"
                :loading="loading"
                >register</v-btn
              >
            </v-layout>
          </v-form>
        </v-tab-item>
      </v-tabs>
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
import {
  LoginInput,
  LoginResult,
  RegisterInput,
  RegisterResult
} from "@/struct";
import { commonChecker } from "@/utils";

@Component
export default class Login extends Vue {
  private tab: string = "tab-login";

  private username: string = "";
  private password: string = "";

  private name: string = "";
  private mail: string = "";
  private regPassword: string = "";
  private repeat: string = "";

  private loginValid: boolean = false;
  private regValid: boolean = false;
  // private departmentItems = constValue.departmentItems;
  private gradeItems: string[] = [];

  private showPassword: boolean = false;
  private againError: string = "";
  private rules = {
    required: (value: string) => !!value || "请填写",
    email: (value: string) =>
      !!(value || "").match(
        /^[_A-Za-z0-9-+]+(.[_A-Za-z0-9-]+)*@[A-Za-z0-9-]+(.[A-Za-z0-9]+)*(.[A-Za-z]{2,})$/
      ) || "非法的邮箱地址",
    passLen: commonChecker.passLen,
    mailLen: commonChecker.mailLen,
    password: commonChecker.password
  };

  private loading: boolean = false;

  private infoText: string = "";
  private hasInfo: boolean = false;

  created() {
    let year = new Date().getFullYear() - 9;
    let items = [];
    for (let i = 0; i < 10; i++) {
      items.push((year + i).toString());
    }
    this.gradeItems = items;
  }

  check() {
    if (this.regPassword != this.repeat) this.againError = "密码不一致";
  }

  async login() {
    if (!this.loginValid) {
      (this.$refs.loginForm as any).validate();
      return;
    }
    this.loading = true;
    try {
      let res = await this.$apollo.mutate<LoginResult, LoginInput>({
        mutation: gql`
          mutation($input: LoginInput!) {
            login(input: $input) {
              message
              userId
              role
            }
          }
        `,
        variables: {
          input: {
            mail: this.mail,
            password: this.password
          }
        }
      });
      if (res.errors) throw res.errors.map(v => v.message).join(",");
      if (res.data!.login.message) throw res.data!.login.message;
      this.loading = false;
      this.$store.commit("global/setUserIdAndRole", {
        userId: res.data!.login.userId,
        role: res.data!.login.role
      });
      sessionStorage.setItem("user_id", res.data!.login.userId);
      sessionStorage.setItem("role", res.data!.login.role);
      this.$router.replace("/");
    } catch (e) {
      this.loading = false;
      this.infoText = e.toString();
      this.hasInfo = true;
    }
  }

  async register() {
    if (!this.regValid) {
      (this.$refs.registerForm as any).validate();
      return;
    }
    this.loading = true;
    try {
      let res = await this.$apollo.mutate<RegisterResult, RegisterInput>({
        mutation: gql`
          mutation($input: RegisterInput!) {
            register(input: $input) {
              message
            }
          }
        `,
        variables: {
          input: {
            name: this.name,
            password: this.regPassword,
            mail: this.mail
          }
        }
      });
      if (res.errors) throw res.errors.map(v => v.message).join(",");
      if (res.data!.register.message) throw res.data!.register.message;
      this.loading = false;
      this.tab = "tab-login";
    } catch (e) {
      this.loading = false;
      this.infoText = e.toString();
      this.hasInfo = true;
    }
  }
}
</script>

<style lang="scss" scoped>
.forget {
  margin-right: 100px;
}
</style>
