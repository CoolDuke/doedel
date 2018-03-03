package types

import "time"

type HeatingTime struct {
  Time       TimeString `yaml:"time"`
  Degrees    int64      `yaml:"degrees"`
}

type WorktimeEntry struct {
  From time.Time
  To time.Time
}

type TimetableOffsetStart struct {
  At         TimeString `yaml:"at"`
  Offset     float64    `yaml:"offset"`
  LeastAt    TimeString `yaml:"leastAt"`
}

type TimetableOffsetEnd struct {
  At         TimeString `yaml:"at"`
  Offset     float64    `yaml:"offset"`
  EarliestAt TimeString `yaml:"earliestAt"`
}

type TimetableOffsetDaytime struct {
  Start      TimetableOffsetStart `yaml:"start"`
  End        TimetableOffsetEnd   `yaml:"end"`
}

type TimetableOffsets struct {
  Morning    TimetableOffsetDaytime `yaml:"morning"`
  Evening    TimetableOffsetDaytime `yaml:"evening"`
}

type TimeString time.Time

func (ts *TimeString) UnmarshalYAML(unmarshal func(interface{}) error) error {
    var s string
    err := unmarshal(&s)
    if err != nil {
        return err
    }
    
    //keep switch dates as number of seconds from unix timestamp start
    t, err := time.Parse("15:04 2006-02-01", s + " 1970-01-01")
    if err != nil {
        return err
    }
    *ts = TimeString(t)

    return nil
}
