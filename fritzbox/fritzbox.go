package fritzbox

import (
  "net/url"

  "github.com/coolduke/doedel/config"

  "github.com/op/go-logging"
  "github.com/bpicode/fritzctl/fritz"
  "github.com/bpicode/fritzctl/logger"
)

var log = logging.MustGetLogger("doedel")

type FritzBox struct {
  HomeAuto fritz.HomeAuto
}

func NewFritzBox() (*FritzBox, error) {
  fritzboxUrl, err := url.Parse(config.Conf.FritzBox.Url)
  if err != nil {
    return nil, err
  }

  log.Debugf("Trying %s", config.Conf.FritzBox.Url)
  homeAuto := fritz.NewHomeAuto(
    fritz.SkipTLSVerify(),
    fritz.URL(fritzboxUrl),
    fritz.Credentials(config.Conf.FritzBox.Username, config.Conf.FritzBox.Password),
  )

  l := &logger.Level{}
  l.Set("warn")

  err = homeAuto.Login()
  if err != nil {
    log.Errorf("Unable to login: %s", err.Error())
    return nil, err
  }

  return &FritzBox{HomeAuto: homeAuto}, nil
}

func (fb *FritzBox) LogCurrentTemperatures() error {
  devices, err := fb.HomeAuto.List()
  if err != nil {
    log.Errorf("Unable to log current temperatures: %s", err.Error())
    return err
  }
  
  for _, device := range devices.Thermostats() {
    log.Infof("Current temperature for %s: %s°C", device.Name, device.Thermostat.FmtMeasuredTemperature())
  }
  
  return nil
}

func (fb *FritzBox) SetTemperature(thermostat string, value float64) (error) {
  err := fb.HomeAuto.Temp(value, thermostat)
  if err != nil {
    log.Errorf("Unable to set temperature for %s: %s", thermostat, err.Error())
    return err
  }
  
  return nil
}
