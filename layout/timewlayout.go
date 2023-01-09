package layout

import (
	"encoding"
	"strings"
	"time"
	"track/types"
)

// Timewarrior format:
// time format: `20060102T150405Z`
// (`Start` and `End` are represented in the time format above)
// Interval (without tags): 		`inc Start - End`
// Interval (with tags): 			`inc Start - End # tag1 tag2 tag3`
// Open Interval (without tags): 	`inc Start`
// Open Interval (with tags): 		`inc Start # tag1 tag2 tag3`

const timeFormat = "20060102T150405Z"
const entryPrefix = "inc "
const timeSeparator = " - "
const tagPrefix = " # "
const tagSeparator = " "

type timewLayout struct {
}

var TimewLayout = timewLayout{}

func makeInterval(str *strings.Builder, interval types.Interval) {
	startTime := interval.Start.Format(timeFormat)
	endTime := interval.End.Format(timeFormat)
	str.WriteString(entryPrefix)
	str.WriteString(startTime)
	str.WriteString(timeSeparator)
	str.WriteString(endTime)
	if interval.Tags != nil && len(interval.Tags) > 0 {
		str.WriteString(tagPrefix)
		str.WriteString(strings.Join(types.TagsToStringSlice(interval.Tags...), tagSeparator))
	}
}

func makeOpenInterval(str *strings.Builder, o types.OpenInterval) {
	startTime := o.Start.Format(timeFormat)
	str.WriteString(entryPrefix)
	str.WriteString(startTime)

	if o.Tags != nil && len(o.Tags) > 0 {
		str.WriteString(tagPrefix)
		str.WriteString(strings.Join(types.TagsToStringSlice(o.Tags...), " "))
	}
}

func (t timewLayout) ToIntervals(m encoding.BinaryMarshaler) (types.Intervals, error) {
	data, err := m.MarshalBinary()
	if err != nil {
		return types.Intervals{}, err
	}
	intervals := types.MakeIntervals()
	for _, line := range strings.Split(string(data), "\n") {
		if !strings.HasPrefix(line, entryPrefix) {
			continue
		}
		line = line[len(entryPrefix):]
		start, err := time.Parse(timeFormat, line[:len(timeFormat)])
		if err != nil {
			return types.Intervals{}, err
		}
		line = line[len(timeFormat):]
		if strings.HasPrefix(line, timeSeparator) {
			line = line[len(timeSeparator):]
			end, err := time.Parse(timeFormat, line[:len(timeFormat)])
			if err != nil {
				return types.Intervals{}, err
			}
			intervals.Intervals = append(intervals.Intervals, types.MakeInterval(start, end, parseTags(line)...))
		} else {
			intervals.OpenIntervals = append(intervals.OpenIntervals, types.MakeOpenInterval(start, parseTags(line)...))
		}
	}
	return intervals, nil
}

func parseTags(str string) []types.Tag {
	if !strings.HasPrefix(str, tagPrefix) {
		return types.MakeTags()
	}
	str = str[len(tagPrefix):]
	return types.MakeTags(strings.Split(str, tagSeparator)...)
}

func (t timewLayout) FromIntervals(intervals types.Intervals) (encoding.BinaryMarshaler, error) {
	var res strings.Builder
	for _, i := range intervals.Intervals {
		makeInterval(&res, i)
		res.WriteRune('\n')
	}
	for _, o := range intervals.OpenIntervals {
		makeOpenInterval(&res, o)
		res.WriteRune('\n')
	}
	return types.Bytes(res.String()), nil
}
