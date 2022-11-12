package types

import (
	"golang.org/x/text/unicode/norm"
	"log"
	"regexp"
	"server/entity"
	"strconv"
	"strings"
)

var blankRegexp *regexp.Regexp

func init() {
	blankRegexp, _ = regexp.Compile("\\s")
}

func lengthLimit(data string, min int, max int) bool {
	return len([]rune(data)) >= min && len([]rune(data)) <= max
}

type RegisterInput struct {
	Name     string
	Password string
	Mail     string
}

func (input *RegisterInput) CheckPass() bool {
	input.Name = blankRegexp.ReplaceAllString(input.Name, "")
	input.Mail = blankRegexp.ReplaceAllString(input.Mail, "")
	input.Name = norm.NFC.String(input.Name)
	return lengthLimit(input.Name, 1, 20) && lengthLimit(input.Password, 8, 18) && lengthLimit(input.Mail, 1, 50)
}

type RegisterResult struct {
	// success, xxx already in use, invalid xxx, failed, already login, not empty error
	Message string
}

type LoginInput struct {
	Mail     string
	Password string
}

func (input *LoginInput) CheckPass() bool {
	input.Mail = blankRegexp.ReplaceAllString(input.Mail, "")
	return lengthLimit(input.Password, 8, 18) && lengthLimit(input.Mail, 1, 50)
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

func (input *ResetInput) CheckPass() bool {
	return lengthLimit(input.Password, 8, 18) && lengthLimit(input.Token, 1, 60)
}

type ResetResult struct {
	Message string
}

type SubmitInput struct {
	ChallengeId string
	Flag        string
}

func (input *SubmitInput) CheckPass() bool {
	input.Flag = blankRegexp.ReplaceAllString(input.Flag, "")
	return input.Flag != ""
}

type SubmitResult struct {
	Message string
}

type BulletinPubInput struct {
	Title   string
	Content string
	Topping bool
}

func (input *BulletinPubInput) CheckPass() bool {
	input.Title = strings.TrimSpace(input.Title)
	input.Content = strings.TrimSpace(input.Content)
	return input.Title != ""
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

func (input *UserInfoUpdateInput) CheckPass() bool {
	input.Name = blankRegexp.ReplaceAllString(input.Name, "")
	input.Mail = blankRegexp.ReplaceAllString(input.Mail, "")
	return lengthLimit(input.Name, 1, 20) && lengthLimit(input.Mail, 1, 50) && checkRole(input.Role) && checkUserState(input.State)
}

func checkRole(str string) bool {
	return str == "member" || str == "team" || str == "admin"
}

func checkUserState(str string) bool {
	return str == "normal" || str == "disabled"
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
	State        string
	Singleton    bool
	NodeConfig   *[]NodeConfigInput
}

func (input *ChallengeMutateInput) CheckPass() bool {
	input.Name = strings.TrimSpace(input.Name)
	lines := strings.Split(input.Description, "\n")
	for index, line := range lines {
		lines[index] = strings.Join(strings.Fields(line), " ")
	}
	input.Description = strings.Join(lines, "\n")
	//log.Println(input)
	if input.NodeConfig == nil {
		if input.Flag.Type == Dynamic {
			return false
		}
		input.NodeConfig = &[]NodeConfigInput{}
	}
	if len(*input.NodeConfig) == 0 {
		input.Singleton = true
	}
	nodeNameSet := map[string]struct{}{}
	for _, node := range *input.NodeConfig {
		if !node.CheckPass() {
			return false
		}
		nodeNameSet[node.Name] = struct{}{}
	}
	if len(nodeNameSet) != len(*input.NodeConfig) {
		return false
	}
	//log.Println(input)
	return input.Name != "" && checkChallengeCategory(input.Category) && input.Score.CheckPass() && input.Flag.CheckPass() && checkChallengeState(input.State) && input.Score.CheckPass() && input.Flag.CheckPass()
}

func checkChallengeCategory(str string) bool {
	return str == "WEB" || str == "RE" || str == "PWN" || str == "MISC" || str == "CRYPTO" || str == "HARDWARE" || str == "RW"
}

func checkChallengeState(str string) bool {
	return str == "enabled" || str == "disabled"
}

type ScoreTypeInput struct {
	Dynamic   bool
	BaseScore string
}

func (input *ScoreTypeInput) CheckPass() bool {
	parsedScore, err := strconv.Atoi(input.BaseScore)
	if err != nil {
		return false
	}
	return parsedScore >= 0
}

type FlagTypeInput struct {
	Type  int32
	Value string
}

func (input *FlagTypeInput) CheckPass() bool {
	if input.Type != Multiple {
		input.Value = blankRegexp.ReplaceAllString(input.Value, "")
		if strings.Contains(input.Value, "\n") {
			return false
		}
	} else {
		flags := []string{}
		parts := strings.Split(input.Value, "\n")
		for _, line := range parts {
			s := blankRegexp.ReplaceAllString(line, "")
			if len(s) == 0 {
				continue
			}
			flags = append(flags, s)
		}
		input.Value = strings.Join(flags, "\n")
	}
	if input.Type == Regexp {
		_, err := regexp.Compile(input.Value)
		if err != nil {
			return false
		}
	}
	return input.Value != "" && 0 <= input.Type && input.Type < Max
}

type NodeConfigInput struct {
	Name         string
	Image        string
	ServicePorts []ServicePortInput
}

func (input *NodeConfigInput) CheckPass() bool {
	input.Name = strings.ToLower(strings.TrimSpace(input.Name))
	input.Image = strings.ToLower(strings.TrimSpace(input.Image))
	matched, err := regexp.MatchString("[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*", input.Image)
	if !matched || err != nil {
		log.Println("NodeConfig Check fault", err)
		return false
	}
	portNameSet := map[string]struct{}{}
	for _, port := range input.ServicePorts {
		if !port.CheckPass() {
			return false
		}
		portNameSet[port.Name] = struct{}{}
	}
	if len(portNameSet) != len(input.ServicePorts) {
		return false
	}
	return input.Name != "" && input.Image != ""
}

func (input *NodeConfigInput) ToNodeConfig() NodeConfig {
	servicePorts := []ServicePort{}
	for _, port := range input.ServicePorts {
		servicePorts = append(servicePorts, port.ToServicePort())
	}
	return NodeConfig{
		Name:  input.Name,
		Image: input.Name,
		//Ports:        ports,
		ServicePorts: servicePorts,
	}
}

type ServicePortInput struct {
	Name     string
	Protocol string
	External int32
	Internal int32
	Pod      int32
}

func (input *ServicePortInput) ToServicePort() ServicePort {
	return ServicePort{
		Name:     input.Name,
		Protocol: input.Protocol,
		External: input.External,
		Internal: input.Internal,
		Pod:      input.Pod,
	}
}

func (input *ServicePortInput) CheckPass() bool {
	input.Name = strings.ToLower(strings.TrimSpace(input.Name))
	return input.Name != "" &&
		(input.Protocol == "TCP" || input.Protocol == "UDP") &&
		input.External > 0 && input.External < 65535 &&
		input.Internal > 0 && input.Internal < 65535 &&
		input.Pod >= 0 && input.Pod < 65535
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
	Score    int32
	Role     string
	State    string
	//Rank     int
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
	Score  int32
}

type ChallengeInfo struct {
	ChallengeId string
	Category    string
	Name        string
	Score       int32
	SolvedNum   int32
	Blood       []BloodInfo
	Done        bool
}

type ChallengeDesc struct {
	ChallengeId  string
	Description  string
	ExternalLink []string
	Manual       bool
	Allocated    int32
}

const AllocatedUndone int32 = 0
const AllocatedDoing int32 = 1
const AllocatedDone int32 = 2

type BloodInfo struct {
	UserId string
	Name   string
	Avatar string
}

type ChallengeConfigsResult struct {
	Message          string
	ChallengeConfigs []ChallengeConfigAndState
}

type ChallengeConfigAndState struct {
	ChallengeId string
	Name        string
	Config      ChallengeConfig
	State       string
}

type ChallengeConfig struct {
	Category     string
	Score        ScoreType
	Flag         FlagType
	Description  string
	ExternalLink []string
	Singleton    bool
	NodeConfig   []NodeConfig
}

type NodeConfig struct {
	Name         string
	Image        string
	ServicePorts []ServicePort
}

type ServicePort struct {
	Name     string
	Protocol string
	External int32
	Internal int32
	Pod      int32
}

type ScoreType struct {
	Dynamic   bool
	BaseScore string
}

type FlagType struct {
	Type  int32
	Value string
}

const (
	Single = iota
	Multiple
	Regexp
	Dynamic
	Max
)

type BulletinResult struct {
	Message   string
	Bulletins []Bulletin
}

const BulletinStyleCommon = "common"
const BulletinStyleBlood1st = "first"
const BulletinStyleBlood2nd = "second"
const BulletinStyleBlood3rd = "third"

type Bulletin struct {
	Style       string
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
}

type WriteUpInfoResult struct {
	Message string
	Infos   []WriteUpInfo
}

type WriteUpInfo struct {
	UserId string
	Name   string
	Mail   string
	Solved int32
}

type StartReplicaResult struct {
	Message string
}

type ImageInfoResult struct {
	Message string
	Infos   []ImageInfo
}

type ImageInfo struct {
	Name     string
	Platform string
	Size     string
	Digest   string
}

type DeleteImageResult struct {
	Message string
}

type DeleteReplicaResult struct {
	Message string
}

type ClusterNodeInfo struct {
	Name         string
	OsType       string
	Distribution string
	Kernel       string
	Arch         string
}

type ClusterReplicaInfo struct {
	Name   string
	Node   string
	Status string
}

type ClusterInfoResult struct {
	Message  string
	Nodes    []ClusterNodeInfo
	Replicas []ClusterReplicaInfo
}

type ChallengeActionInput struct {
	Action       string
	ChallengeIds []string
}

func (input *ChallengeActionInput) CheckPass() bool {
	if input.Action != "enable" && input.Action != "disable" && input.Action != "delete" {
		return false
	}
	return true
}

type ChallengeActionResult struct {
	Message    string
	Successful []string
}

type WatchDescriptionResult struct {
	Message     string
	Description ChallengeDesc
}

type AllocStatusResult struct {
	Message   string
	Allocated int32
}

type GameEvent struct {
	EventId string
	Time    string
	Action  int32
}

type AddEventInput struct {
	Action int32
	Time   string
}

func (input *AddEventInput) CheckPass() bool {
	return input.Action < entity.EventMax
}

type AddEventResult struct {
	Message string
}

type UpdateEventInput struct {
	EventId string
	Time    string
}

func (input *UpdateEventInput) CheckPass() bool {
	//TODO
	return true
}

type UpdateEventResult struct {
	Message string
}

type DeleteEventInput struct {
	EventIds []string
}

func (input *DeleteEventInput) CheckPass() bool {
	//TODO
	return true
}

type DeleteEventResult struct {
	Message string
}

type AllEventResult struct {
	Message   string
	AllEvents []GameEvent
}
