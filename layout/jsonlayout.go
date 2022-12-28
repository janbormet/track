package layout

import (
	"encoding"
	"encoding/json"
	"time"
	"track/types"
)

type (
	Tags          []string
	jsonLayout    struct{}
	JSONIntervals struct {
		Intervals     []JSONInterval     `json:"intervals"`
		OpenIntervals []JSONOpenInterval `json:"openIntervals"`
	}
	JSONInterval struct {
		Start time.Time `json:"start"`
		End   time.Time `json:"end"`
		Tags  Tags      `json:"tags"`
	}
	JSONOpenInterval struct {
		Start time.Time `json:"start"`
		Tags  Tags      `json:"tags"`
	}
)

var JSONLayout = jsonLayout{}

func (J jsonLayout) ToIntervals(m encoding.BinaryMarshaler) (types.Intervals, error) {
	data, err := m.MarshalBinary()
	if err != nil {
		return types.Intervals{}, err
	}
	var jsonIntervals JSONIntervals
	err = json.Unmarshal(data, &jsonIntervals)
	if err != nil {
		return types.Intervals{}, err
	}
	return jsonIntervals.decode(), nil
}

func (J jsonLayout) FromIntervals(intervals types.Intervals) (encoding.BinaryMarshaler, error) {
	jsonIntervals := makeJSONIntervals(intervals)
	data, err := json.Marshal(jsonIntervals)
	if err != nil {
		return nil, err
	}
	return types.Bytes(data), nil
}

func makeTags(tags []types.Tag) Tags {
	t := make(Tags, len(tags))
	for i, tag := range tags {
		t[i] = tag.String()
	}
	return t
}

func makeJSONInterval(interval types.Interval) JSONInterval {
	return JSONInterval{
		Start: interval.Start,
		End:   interval.End,
		Tags:  makeTags(interval.Tags),
	}
}

func makeJSONOpenInterval(i types.OpenInterval) JSONOpenInterval {
	return JSONOpenInterval{
		Start: i.Start,
		Tags:  makeTags(i.Tags),
	}
}

func makeJSONIntervals(intervals types.Intervals) JSONIntervals {
	jsonIntervalList := make([]JSONInterval, 0, len(intervals.Intervals))
	for _, interval := range intervals.Intervals {
		jsonIntervalList = append(jsonIntervalList, makeJSONInterval(interval))
	}
	jsonOpenIntervalList := make([]JSONOpenInterval, 0, len(intervals.OpenIntervals))
	for _, openInterval := range intervals.OpenIntervals {
		jsonOpenIntervalList = append(jsonOpenIntervalList, makeJSONOpenInterval(openInterval))
	}

	return JSONIntervals{
		Intervals:     jsonIntervalList,
		OpenIntervals: jsonOpenIntervalList,
	}
}

func (i JSONIntervals) decode() types.Intervals {
	intervals := make([]types.Interval, 0, len(i.Intervals))
	for _, interval := range i.Intervals {
		intervals = append(intervals, interval.decode())
	}
	openIntervals := make([]types.OpenInterval, 0, len(i.OpenIntervals))
	for _, openInterval := range i.OpenIntervals {
		openIntervals = append(openIntervals, openInterval.decode())
	}
	return types.Intervals{
		Intervals:     intervals,
		OpenIntervals: openIntervals,
	}
}

func (i JSONInterval) decode() types.Interval {
	return types.Interval{
		Start: i.Start,
		End:   i.End,
		Tags:  i.Tags.decode(),
	}
}

func (i JSONOpenInterval) decode() types.OpenInterval {
	return types.OpenInterval{
		Start: i.Start,
		Tags:  i.Tags.decode(),
	}
}

func (t Tags) decode() []types.Tag {
	tags := make([]types.Tag, len(t))
	for i, tag := range t {
		tags[i] = types.MakeTag(tag)
	}
	return tags
}

var _ types.Layout = jsonLayout{}
