package utils

type Config struct {
	Server       Server     `yaml:"server"`
	Email        Email      `yaml:"email"`
	Challenge    Challenge  `yaml:"challenge"`
	Kubernetes   Kubernetes `yaml:"kubernetes"`
	Database     Database   `yaml:"database"`
	TimelineFile string     `yaml:"timelineFile"`
}

type Server struct {
	Host        string `yaml:"host"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	Mail        string `yaml:"mail"`
	Port        int    `yaml:"port"`
	Salt        string `yaml:"salt"`
	BehaviorLog bool   `yaml:"behaviorLog"`
	Debug       Debug  `yaml:"debug"`
}

type Email struct {
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

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

type Debug struct {
	DockerOpDetail bool `yaml:"dockerOpDetail"`
	NoOriginCheck  bool `yaml:"noOriginCheck"`
	DBOpDetail     bool `yaml:"dbOpDetail"`
}
