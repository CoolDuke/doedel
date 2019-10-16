package heating

import (
  "time"
  "sort"

  "github.com/coolduke/doedel/config"
  "github.com/coolduke/doedel/types"

  "github.com/op/go-logging"
)

var log = logging.MustGetLogger("doedel")

type Heating struct {
  Timetable []TimetableEntry
}

type TimetableEntry struct {
  SwitchAt time.Time
  Degrees int64
}

func NewHeating() (*Heating, error) {
  now := time.Now()

  var entries []TimetableEntry

  d := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
  for i := 0; i < 31; i = i + 1 { //92
    weekday := d.Weekday().String()
    times, ok := config.Conf.Heating.Defaults[weekday]
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
  
  return &Heating{Timetable: entries}, nil
}

func (h *Heating) ApplyWorktimes(worktimes []types.WorktimeEntry) (error) {
  var newEntries []TimetableEntry

  //remove all dates we will override
  for _, timetableEntry := range h.Timetable {
    found := false
    for _, worktime := range worktimes {
      if worktime.From.Month() == timetableEntry.SwitchAt.Month() && worktime.From.Day() == timetableEntry.SwitchAt.Day() {
        found = true
        break
      }
    }
    if !found {
      newEntries = append(newEntries, timetableEntry)
    }
  }

  //add new entries
  for _, worktime := range worktimes {
    //TODO: implement offsets
    newEntries = append(newEntries, TimetableEntry{worktime.From, 22}, TimetableEntry{worktime.To, 22})
  }

  //get dates back into order
  sort.Slice(newEntries, func(i, j int) bool { return newEntries[i].SwitchAt.Before(newEntries[j].SwitchAt)})

  if log.IsEnabledFor(logging.DEBUG) {
    log.Debugf("Modified timetable:")
    for _, entry := range newEntries {
      log.Debugf("At %s switch to %d", entry.SwitchAt.String(), entry.Degrees)
    }
  }
  
  h.Timetable = newEntries

  return nil
}
