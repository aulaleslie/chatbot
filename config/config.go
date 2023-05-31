package config

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

var runtimeEnvironment = "local"
var cfg *confData

func GetRuntimeEnvironment() string {
	runtimeEnvironment = strings.ToLower(os.Getenv(runtimeEnvironmentEnvKey))
	if len(runtimeEnvironment) == 0 {
		runtimeEnvironment = "local"
		log.Warn().Msg("no runtime environment set. defaulting to \"local\" runtime environment.")
	}
	return runtimeEnvironment
}

func getConfFilepath() string {
	env := GetRuntimeEnvironment()
	return confPathPrefix + env + confPathSuffix
}

func newConfData(path string) (*confData, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		log.Error().Msgf("unable to read conf file at path: %s. %v\n", path, err)
		return nil, err
	}

	data := &confData{}
	dec := json.NewDecoder(bytes.NewBuffer(buf))
	dec.DisallowUnknownFields()
	if err = dec.Decode(data); err != nil {
		log.Error().Err(err).Msgf("unable to deserialize conf data from %s", path)
		return nil, err
	}

	return data, nil
}

func Reload(override *string) error {
	confPath := getConfFilepath()

	log.Info().Msgf("Starting server with runtime environment: %s", runtimeEnvironment)
	log.Info().Msgf("Loading config file: %s", confPath)

	data, err := newConfData("C:\\Workspace\\chatbot\\configs\\local_config.json")
	if err != nil {
		log.Error().Err(err).Msg("error while reading config file")
		return err
	}

	cfg = data

	return nil
}

func NewFacebookConfig() *FacebookConf {
	return &FacebookConf{
		RequestUrl:   cfg.Facebook.RequestUrl,
		RequestToken: cfg.Facebook.RequestToken,
		AudioId:      cfg.Facebook.AudioId,
		VerifyToken:  cfg.Facebook.VerifyToken,
	}
}

func GetListenPort() string {
	return cfg.Server.ListenPort
}
