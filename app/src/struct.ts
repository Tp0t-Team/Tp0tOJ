export interface Result {
    message: string
}

export type LoginResult = {
    userId: string
    role: string
} & Result
