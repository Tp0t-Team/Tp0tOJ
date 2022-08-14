package utils

type Config struct {
	Server Server `yaml:"server"`
	Email  Email  `yaml:"email"`
	//Redis      Redis      `yaml:"redis"`
	Challenge  Challenge  `yaml:"challenge"`
	Kubernetes Kubernetes `yaml:"kubernetes"`
	Database   Database `yaml:"database"`
}

type Server struct {
	Host        string `yaml:"host"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	Mail        string `yaml:"mail"`
	Port        int    `yaml:"port"`
	Salt        string `yaml:"salt"`
	BehaviorLog bool   `yaml:"behaviorLog"`
}

type Email struct {
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

//type Redis struct {
//	Password  string `yaml:"password"`
//	MaxActive int    `yaml:"maxActive"`
//	MaxIdle   int    `yaml:"maxIdle"`
//	Database  int    `yaml:"database"`
//	Host      string `yaml:"host"`
//	Port      int    `yaml:"port"`
//	MaxWait   string `yaml:"maxWait"`
//	MinIdle   int    `yaml:"minIdle"`
//	Timeout   int    `yaml:"timeout"`
//}

type Database struct {
	Dsn string `yaml:"dsn"`
}

type Challenge struct {
	SecondBloodReward float64 `yaml:"secondBloodReward"`
	ThirdBloodReward  float64 `yaml:"thirdBloodReward"`
	HalfLife          int     `yaml:"halfLife"`
	FirstBloodReward  float64 `yaml:"firstBloodReward"`
}

type Kubernetes struct {
	PortAllocBegin int32  `yaml:"portAllocBegin"`
	PortAllocEnd   int32  `yaml:"portAllocEnd"`
	Username       string `yaml:"username"`
	Password       string `yaml:"password"`
	RegistryHost   string `yaml:"registryHost"`
}
