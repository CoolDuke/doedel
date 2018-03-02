package config

import (
    "io/ioutil"

    "github.com/coolduke/doedel/types"

    "gopkg.in/yaml.v2"
    "github.com/op/go-logging"
)

type Config struct {
  FritzBox         *ConfigFritzBox         `yaml:"fritzbox"`
  Heating          *ConfigHeating          `yaml:"heating"`
}

type ConfigFritzBox struct {
  Url      string `yaml:"url"`
  Username string `yaml:"username"`
  Password string `yaml:"password"`
}

type ConfigHeating struct {
  Defaults map[string][]types.HeatingTime `yaml:"defaults"`
  TimeTableOffsets types.TimetableOffsets `yaml:"timetableOffsets"`
}

func GetConfig(log *logging.Logger, filename string) (Config, error) {
  log.Noticef("Reading configuration from: %s", filename)

  bytes, err := ioutil.ReadFile(filename)
  if err != nil {
    return Config{}, err
  }

  var config Config

  err = yaml.Unmarshal(bytes, &config)
  if err != nil {
    return Config{}, err
  }

  //TODO: validate configuration
  log.Infof("Heating defaults %s\n", config)

  return config, nil
}
