package types

import (
	"time"
)

type (
	Tag       string
	Intervals struct {
		Intervals     []Interval
		OpenIntervals []OpenInterval
	}
	Interval struct {
		Start time.Time
		End   time.Time
		Tags  []Tag
	}
	OpenInterval struct {
		Start time.Time
		Tags  []Tag
	}
)

func (o OpenInterval) Close(time time.Time) Interval {
	return Interval{
		Start: o.Start,
		End:   time,
		Tags:  o.Tags,
	}
}

func (t Tag) String() string {
	return string(t)
}

func MakeTag(s string) Tag {
	return Tag(s)
}

func MakeTags(t ...string) []Tag {
	tags := make([]Tag, len(t))
	for i, tag := range t {
		tags[i] = MakeTag(tag)
	}
	return tags
}
