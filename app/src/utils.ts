export const commonChecker = {
  passLen: (minlen: number, maxlen: number) => (v: string) =>
    ((v || "").length >= minlen && (v || "").length <= maxlen) ||
    `非法的密码长度，需要 ${minlen}~${maxlen} 位`,
  mailLen: (maxlen: number) => (v: string) =>
    (v || "").length <= maxlen || `非法的邮箱长度，最多 ${maxlen} 位`,
  password: (value: string) =>
    (!!(value || "").match(/[A-Z]/) &&
      !!(value || "").match(/[a-z]/) &&
      !!(value || "").match(/\d/)) ||
    "密码必须由大小写字母数字和特殊符号组成" //TODO: 正则好像不对
};
