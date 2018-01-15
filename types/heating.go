package types

import (
    "time"
)

type HeatingTime struct {
  Time       TimeString `yaml:"time"`
  Degrees    int64      `yaml:"degrees"`
}

type TimeString time.Time

func (ts *TimeString) UnmarshalYAML(unmarshal func(interface{}) error) error {
    var s string
    err := unmarshal(&s)
    if err != nil {
        return err
    }
    
    //keep switch dates as number of seconds from unix stimestamp start
    t, err := time.Parse("15:04 2006-02-01", s + " 1970-01-01")
    if err != nil {
        return err
    }
    *ts = TimeString(t)

    return nil
}
