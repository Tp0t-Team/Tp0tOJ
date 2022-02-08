package types

type RegisterInput struct {
	Name     string
	Password string
	Mail     string
}

type RegisterResult struct {
	// success, xxx already in use, invalid xxx, failed, already login, not empty error
	Message string
}

type LoginInput struct {
	Mail     string
	Password string
}

type LoginResult struct {
	Message string
	UserId  string
	Role    string
}

type LogoutResult struct {
	Message string
}

type ForgetResult struct {
	Message string
}

type ResetInput struct {
	Password string
	Token    string
}

type ResetResult struct {
	Message string
}

type SubmitInput struct {
	ChallengeId string
	Flag        string
}

type SubmitResult struct {
	Message string
}

type BulletinPubInput struct {
	Title   string
	Content string
	Topping bool
}

type BulletinPubResult struct {
	Message string
}

type UserInfoUpdateInput struct {
	UserId string
	Name   string
	Role   string
	Mail   string
	State  string
}

type UserInfoUpdateResult struct {
	Message string
}

type ChallengeMutateInput struct {
	ChallengeId  string
	Name         string
	Category     string
	Score        ScoreTypeInput
	Flag         FlagTypeInput
	Description  string
	ExternalLink []string
	Hint         []string
	State        string
}

type ScoreTypeInput struct {
	Dynamic   bool
	BaseScore string
}

type FlagTypeInput struct {
	Dynamic bool
	Value   string
}

type ChallengeMutateResult struct {
	Message string
}

type ChallengeRemoveResult struct {
	Message string
}

type UserInfoResult struct {
	Message  string
	UserInfo UserInfo
}

type UserInfo struct {
	UserId   string
	Name     string
	Avatar   string
	Mail     string
	JoinTime string
	Score    string
	Role     string
	State    string
	Rank     int
}

type AllUserInfoResult struct {
	Message      string
	AllUserInfos []UserInfo
}

type RankResult struct {
	Message         string
	RankResultDescs []RankResultDesc
}

type ChallengeInfosResult struct {
	Message        string
	ChallengeInfos []ChallengeInfo
}

type RankResultDesc struct {
	UserId string
	Name   string
	Avatar string
	Score  int
}

type ChallengeInfo struct {
	ChallengeId  string
	Category     string
	Name         string
	Score        int
	Description  string
	ExternalLink []string
	Hint         []string
	Blood        []BloodInfo
	Done         bool
}

type BloodInfo struct {
	UserId string
	Name   string
	Avatar string
}

type ChallengeConfigsResult struct {
	Message          string
	ChallengeConfigs []ChallengeConfig
}

type ChallengeConfig struct {
	ChallengeId  string
	Name         string
	Category     string
	Score        ScoreType
	Flag         FlagType
	Description  string
	ExternalLink []string
	Hint         []string
	State        string
}

type ScoreType struct {
	Dynamic   bool
	BaseScore string
}

type FlagType struct {
	Dynamic bool
	Value   string
}

type BulletinResult struct {
	Message   string
	Bulletins []Bulletin
}

type Bulletin struct {
	Title       string
	Content     string
	PublishTime string
}

type SubmitHistoryResult struct {
	Message     string
	SubmitInfos []SubmitInfo
}

type SubmitInfo struct {
	SubmitTime    string
	ChallengeName string
	Mark          int
}
