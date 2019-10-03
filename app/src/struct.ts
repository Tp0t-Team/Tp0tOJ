export interface Result {
    message: string
}

export type LoginResult = {
    login: {
        userId: string
        role: string
    } & Result
}

export type RegisterResult = { register: Result }

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

export type UserInfoResult = { userInfo: { userInfo: UserInfo } & Result }

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

export type ChallengeResult = {
    challenges: { challengeInfos: ChallengeDesc[] } & Result
}

export interface RankDesc {
    userId: string
    name: string
    score: number
}

export type RankResult = { rank: { userInfos: RankDesc[] } & Result }

export type BulletinPubResult = { bulletinPub: Result }

export type BulletinSubResult = { bulletin: BulletinItem }

export interface BulletinItem {
    title: string
    description: string
    time: string
}

export type AllBulletinResult = {
    allBulletin: { bulletin: BulletinItem[] } & Result
}

export interface ChallengeConfig {
    name: string
    type: string
    score: {
        dynamic: boolean
        base_score: number
    }
    flag: {
        dynamic: boolean
        value: string
    }
    description: string
    external_link: string[]
    hint: string[]
}

export type ChallengeConfigWithId = { challengeId: string } & ChallengeConfig

export type ChallengeConfigResult = {
    challenges: {
        challengeInfos: ChallengeConfigWithId[]
    } & Result
}

export interface ResolveInfo {
    submitTime: string
    challengeName: string
    mark: number
}
