export interface Result {
    message: string
}

export type LoginResult = {
    userId: string
    role: string
} & Result

export interface UserInfo {
    name: string
    role: string
    stuNumber: string
    department: string
    grade: string
    protectedTime: string
    qq: string
    mail: string
    topRank: string
    joinTime: string
}

export type UserInfoResult = { userInfo: UserInfo } & Result
