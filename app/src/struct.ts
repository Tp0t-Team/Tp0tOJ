export interface Result {
    message: string
}

export interface UserInfo {
    userId: string
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

export interface RankDesc {
    userId: string
    name: string
    score: number
}

export interface BulletinItem {
    title: string
    content: string
    publishTime: string
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

export type ChallengeConfigWithId = {
    challengeId: string
    state: string
} & ChallengeConfig

// query userInfo
export type UserInfoResult = { userInfo: { userInfo: UserInfo } & Result }

// query challenges
export type ChallengeResult = {
    challenges: { challengeInfos: ChallengeDesc[] } & Result
}

// query rank
export type RankResult = { rank: { userInfos: RankDesc[] } & Result }

// mutation login
export interface LoginInput {
    input: {
        stuNumber: string
        password: string
    }
}
export type LoginResult = {
    login: {
        userId: string
        role: string
    } & Result
}

// query allBulletin
export type AllBulletinResult = {
    allBulletin: { bulletins: BulletinItem[] } & Result
}

// mutation register
export interface RegisterInput {
    input: {
        name: string
        stuNumber: string
        password: string
        department: string
        grade: string
        qq: string
        mail: string
    }
}
export type RegisterResult = { register: Result }

// mutation logout
export type LogoutResult = { logout: Result }

// mutation sumbit
export interface SubmitInput {
    challengeId: string
    flag: string
}
export type SubmitResult = {
    submit: Result
}

// mutation bulletinPub
export interface BulletinPubInput {
    input: {
        title: string
        content: string
        topping: boolean
    }
}
export type BulletinPubResult = { bulletinPub: Result }

// subscription bulletin
export type BulletinSubResult = { bulletin: BulletinItem }

// mutation userInfoUpdate
export interface UserInfoUpdateInput {
    input: {
        userId: string
        name: string
        role: string
        department: string
        grade: string
        protectedTime: string
        qq: string
        mail: string
        state: string
    }
}
export type UserInfoUpdateResult = { userInfoUpdate: Result }

// TODO:

// query challengeConfig
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
