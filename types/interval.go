package types

import (
	"regexp"
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
	illegal := regexp.MustCompile("[^A-Za-z0-9-]")
	return Tag(illegal.ReplaceAll([]byte(s), []byte("")))
}

func TagsToStringSlice(tags ...Tag) []string {
	res := make([]string, len(tags))
	for i, tag := range tags {
		res[i] = tag.String()
	}
	return res
}

func MakeTags(t ...string) []Tag {
	tags := make([]Tag, len(t))
	for i, tag := range t {
		tags[i] = MakeTag(tag)
	}
	return tags
}

func MakeIntervals() Intervals {
	return Intervals{
		Intervals:     []Interval{},
		OpenIntervals: []OpenInterval{},
	}
}

func MakeInterval(start time.Time, end time.Time, tags ...Tag) Interval {
	return Interval{
		Start: start,
		End:   end,
		Tags:  tags,
	}
}

func MakeOpenInterval(start time.Time, tags ...Tag) OpenInterval {
	return OpenInterval{
		Start: start,
		Tags:  tags,
	}
}
