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
    score: number
    state: string
    rank: number
}

export type UserInfoResult = { userInfo: UserInfo } & Result

export interface ChallengeDesc {
    challengeId: string
    type: string
    name: string
    score: number
    description: string
    externalLink: string[]
    hint: string[]
    blood: string[]
    done: boolean
}

export type ChallengeResult = { challenges: ChallengeDesc[] } & Result

export interface RankDesc {
    rank: number
    name: string
    avatar: string
    score: number
}

export type RankResult = { ranks: RankDesc[] } & Result