package heating

import (
    "time"

    "github.com/coolduke/doedel/config"
//    "github.com/coolduke/doedel/types"

    "github.com/op/go-logging"
)

type Timetable struct {
  Log *logging.Logger
  Config *config.ConfigHeating
  Entries []TimetableEntry
}

type TimetableEntry struct {
  SwitchAt time.Time
  Degrees int64
}

func NewHeating(log *logging.Logger, conf config.ConfigHeating) (*Timetable, error) {
  now := time.Now()

  var entries []TimetableEntry

  d := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
  for i := 0; i < 92; i = i + 1 {
    weekday := d.Weekday().String()
    times, ok := conf.Defaults[weekday]
    if ok == false {
      log.Warningf("Missing default definition for %s, using 16 degrees as default %s", weekday)
      entries = append(entries, TimetableEntry{d, 16})
    } else {
      for _, heatingTime := range times {
        entries = append(entries, TimetableEntry{d.Add(time.Duration(time.Time(heatingTime.Time).Unix()) * 1000000000), heatingTime.Degrees})
      }
    }
    d = d.AddDate(0, 0, 1)
  }

  if log.IsEnabledFor(logging.DEBUG) {
    for _, entry := range entries {
      log.Debugf("At %s switch to %d", entry.SwitchAt.String(), entry.Degrees)
    }
  }
  
  return &Timetable{Log: log, Config: &conf, Entries: entries}, nil
}
