package main

import (
    "fmt"
    "os"
    
    "github.com/coolduke/doedel/config"
    "github.com/coolduke/doedel/zeitkonto"
    "github.com/coolduke/doedel/fritzbox"

    "github.com/op/go-logging"
)

var log = logging.MustGetLogger("doedel")
var format = logging.MustStringFormatter(
    `%{color}%{time:15:04:05.000} %{level:-8s} ▶ %{shortpkg:-10s} ▶%{color:reset} %{message}`,
)

func main() {
    logBackend := logging.NewLogBackend(os.Stderr, "", 0)
    logBackendFormatter := logging.NewBackendFormatter(logBackend, format)
    logBackendLeveled := logging.AddModuleLevel(logBackendFormatter)
    logBackendLeveled.SetLevel(logging.DEBUG, "")
    logging.SetBackend(logBackendLeveled)

    if len(os.Args) != 2 {
        fmt.Fprintf(os.Stderr, "Usage: doedel <pdf file>\n")
        os.Exit(2)
    }
    filename := os.Args[1]

    //load configuration
    config, err := config.GetConfig(log, "doedel.yml")
    if err != nil {
      log.Error(err.Error())
      os.Exit(1)
    }
    
    extractor, err := zeitkonto.NewExtractor(log, filename)
    if err != nil {
      log.Error(err.Error())
      os.Exit(1)
    }

    fritzbox, err := fritzbox.NewFritzBox(log, *config.FritzBox)
    if err != nil {
      log.Error(err.Error())
      os.Exit(1)
    }

    timetable := extractor.GetTimetable()
    
    log.Infof("Got %d records", len(timetable))

    fritzbox.LogCurrentTemperatures()
}
