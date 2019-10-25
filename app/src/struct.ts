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
    score: string
    state: string
    rank: number
}

export interface ChallengeDesc {
    challengeId: string
    type: string
    name: string
    score: number
    description: string
    external_link: string[]
    hint: string[]
    blood: {
        userId: string
        name: string
    }[]
    done: boolean
}

export interface RankDesc {
    userId: string
    name: string
    score: string
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
        base_score: string | number
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

export interface SubmitInfo {
    submitTime: String
    challengeName: String
    mark: number
}

// query userInfo
export type UserInfoResult = { userInfo: { userInfo: UserInfo } & Result }

export type AllUserInfoResult = {
    allUserInfos: { allUserInfos: UserInfo[] } & Result
}

// query challenges
export type ChallengeResult = {
    challengeInfos: { challengeInfos: ChallengeDesc[] } & Result
}

// query challengeConfig
export type ChallengeConfigResult = {
    challengeConfigs: {
        challengeConfigs: ChallengeConfigWithId[]
    } & Result
}

// query rank
export type RankResult = { rank: { rankResultDescs: RankDesc[] } & Result }

// query allBulletin
export type AllBulletinResult = {
    allBulletin: { bulletins: BulletinItem[] } & Result
}

// query submitHistory
export type SubmitHistoryResult = {
    submitHistory: { submitInfos: SubmitInfo[] } & Result
}

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
    input: {
        challengeId: string
        flag: string
    }
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

// mutation challengeMutate
export type ChallengeMutateResult = { challengeMutate: Result }

// subscription bulletin
export type BulletinSubResult = { bulletin: BulletinItem }
