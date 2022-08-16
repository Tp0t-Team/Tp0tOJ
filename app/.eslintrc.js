module.exports = {
  root: true,
  env: {
    node: true
  },
  extends: ["plugin:vue/essential", "@vue/prettier", "@vue/typescript"],
  rules: {
    "no-console": process.env.NODE_ENV === "production" ? "error" : "off",
    "no-debugger": process.env.NODE_ENV === "production" ? "error" : "off",
    "vue/valid-v-slot": [
      "error",
      {
        allowModifiers: true
      }
    ]
  },
  parserOptions: {
    parser: "@typescript-eslint/parser"
  },
  globals: {
    globalThis: "writable"
  }
};
