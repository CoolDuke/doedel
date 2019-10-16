package config

import (
  "io/ioutil"

  "github.com/coolduke/doedel/types"

  "gopkg.in/yaml.v2"
  "github.com/op/go-logging"
)

var log = logging.MustGetLogger("doedel")
var Conf Config

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

func LoadConfig(filename string) (error) {
  log.Noticef("Reading configuration from: %s", filename)

  bytes, err := ioutil.ReadFile(filename)
  if err != nil {
    log.Errorf("Unable to load configuration from %s: %s", filename, err.Error())
    return err
  }

  var config Config

  err = yaml.Unmarshal(bytes, &config)
  if err != nil {
    log.Errorf("Unable to parse YAML from %s: %s", filename, err.Error())
    return err
  }

  //TODO: validate configuration
//  log.Infof("Heating defaults %s\n", config)

  Conf = config

  return nil
}
