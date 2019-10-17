package heating

import (
  "time"
  "sort"

  "github.com/coolduke/doedel/config"
  "github.com/coolduke/doedel/types"

  "github.com/op/go-logging"

_  "github.com/kr/pretty"
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
  for i := 0; i < 92; i = i + 1 {
    weekday := d.Weekday().String()
    times, ok := config.Conf.Heating.Defaults[weekday]
    if ok == false {
      log.Warningf("Missing default definition for %s, using 16 degrees as default %s", weekday)
      entries = append(entries, TimetableEntry{d, 16})
    } else {
      for _, heatingTime := range times {
        entries = append(entries, TimetableEntry{d.Add(time.Duration(time.Time(heatingTime.TimeString.Time).Unix()) * 1000000000), heatingTime.Degrees})
      }
    }
    d = d.AddDate(0, 0, 1)
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
    year, month, day := worktime.From.Date()
    var morningStart, morningEnd, eveningStart, eveningEnd time.Time

    //TODO: Move to type.TimeString
    morningStart = time.Date(year, month, day,
                              config.Conf.Heating.TimeTableOffsets.Morning.Start.LeastAt.Time.Hour(), 
                              config.Conf.Heating.TimeTableOffsets.Morning.Start.LeastAt.Time.Minute(),
                              0, 0, time.Local,
                             )

    morningEnd = time.Date(year, month, day,
                            config.Conf.Heating.TimeTableOffsets.Morning.End.EarliestAt.Time.Hour(), 
                            config.Conf.Heating.TimeTableOffsets.Morning.End.EarliestAt.Time.Minute(),
                            0, 0, time.Local,
                           )

    eveningStart = time.Date(year, month, day,
                              config.Conf.Heating.TimeTableOffsets.Evening.Start.LeastAt.Time.Hour(), 
                              config.Conf.Heating.TimeTableOffsets.Evening.Start.LeastAt.Time.Minute(),
                              0, 0, time.Local,
                             )

    //TODO: also implement if-clause for other times / use At time if defined
    var nilTimeString types.TimeString
    if config.Conf.Heating.TimeTableOffsets.Evening.End.At != nilTimeString {
      eveningEnd = time.Date(year, month, day,
                             config.Conf.Heating.TimeTableOffsets.Evening.End.At.Time.Hour(), 
                             config.Conf.Heating.TimeTableOffsets.Evening.End.At.Time.Minute(),
                             0, 0, time.Local,
                            )
    } else {
      eveningEnd = time.Date(year, month, day,
                             config.Conf.Heating.TimeTableOffsets.Evening.End.EarliestAt.Time.Hour(), 
                             config.Conf.Heating.TimeTableOffsets.Evening.End.EarliestAt.Time.Minute(),
                             0, 0, time.Local,
                            )
    }

    worktimeFromWithOffset := worktime.From.Add(time.Minute * time.Duration(config.Conf.Heating.TimeTableOffsets.Morning.Start.Offset))
    if worktimeFromWithOffset.Before(morningStart) {
      morningStart = worktimeFromWithOffset
    }
    worktimeFromWithOffset = worktime.From.Add(time.Minute * time.Duration(config.Conf.Heating.TimeTableOffsets.Morning.End.Offset))
    if worktimeFromWithOffset.After(morningEnd) {
      morningEnd = worktimeFromWithOffset
    }

    worktimeToWithOffset := worktime.To.Add(time.Minute * time.Duration(config.Conf.Heating.TimeTableOffsets.Evening.Start.Offset))
    if worktimeToWithOffset.Before(eveningStart) {
      eveningStart = worktimeToWithOffset
    }
    worktimeToWithOffset = worktime.To.Add(time.Minute * time.Duration(config.Conf.Heating.TimeTableOffsets.Evening.End.Offset))
    if worktimeToWithOffset.After(eveningEnd) {
      eveningEnd = worktimeToWithOffset
    }

    newEntries = append(newEntries, TimetableEntry{morningStart, 22},
                                    TimetableEntry{morningEnd, 17},
                                    TimetableEntry{eveningStart, 22},
                                    TimetableEntry{eveningEnd, 17})
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
