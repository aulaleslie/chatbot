package config

const (
	confPathPrefix           = "./configs/"
	confPathSuffix           = "_config.json"
	runtimeEnvironmentEnvKey = "APP_ENV"
)

type confData struct {
	Server   serverConf   `json:"server,omitempty"`
	Facebook FacebookConf `json:"facebook,omitempty"`
}

type FacebookConf struct {
	RequestUrl   string `json:"requestUrl"`
	RequestToken string `json:"requestToken"`
	AudioId      string `json:"audioId"`
	VerifyToken  string `json:"verifyToken"`
}

type serverConf struct {
	ListenPort string `json:"listenPort,omitempty"`
}
