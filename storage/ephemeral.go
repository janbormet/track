package storage

import (
	"track/types"
)

type ephemeral struct {
	data types.Intervals
}

func copyTags(tags []types.Tag) []types.Tag {
	newTags := make([]types.Tag, len(tags))
	copy(newTags, tags)
	return newTags
}

func copyInterval(i types.Interval) types.Interval {
	return types.Interval{
		Start: i.Start,
		End:   i.End,
		Tags:  copyTags(i.Tags),
	}
}

func copyOpenInterval(o types.OpenInterval) types.OpenInterval {
	return types.OpenInterval{
		Start: o.Start,
		Tags:  copyTags(o.Tags),
	}
}

func copyIntervalSlice(intervals []types.Interval) []types.Interval {
	newIntervals := make([]types.Interval, len(intervals))
	for i, interval := range intervals {
		newIntervals[i] = copyInterval(interval)
	}
	return newIntervals
}
func copyOpenIntervalSlice(intervals []types.OpenInterval) []types.OpenInterval {
	newOpenIntervals := make([]types.OpenInterval, len(intervals))
	for i, interval := range intervals {
		newOpenIntervals[i] = copyOpenInterval(interval)
	}
	return newOpenIntervals
}

func copyIntervals(intervals types.Intervals) types.Intervals {
	return types.Intervals{
		Intervals:     copyIntervalSlice(intervals.Intervals),
		OpenIntervals: copyOpenIntervalSlice(intervals.OpenIntervals),
	}
}

func (e *ephemeral) Save(intervals types.Intervals) error {
	e.data = copyIntervals(intervals)
	return nil
}

func (e *ephemeral) Load() (types.Intervals, error) {
	return copyIntervals(e.data), nil
}

func NewEphemeralStorage() types.Storage {
	return &ephemeral{
		data: types.Intervals{},
	}
}

type ephemeralWithLayout struct {
	layout types.Layout
	data   types.Bytes
}

func NewEphemeralStorageWithLayout(layout types.Layout) types.Storage {
	return &ephemeralWithLayout{
		layout: layout,
		data:   nil,
	}
}

func (e *ephemeralWithLayout) Save(intervals types.Intervals) error {
	data, err := e.layout.FromIntervals(intervals)
	if err != nil {
		return err
	}
	bytes, err := data.MarshalBinary()
	if err != nil {
		return err
	}
	e.data = bytes
	return nil
}

func (e *ephemeralWithLayout) Load() (types.Intervals, error) {
	if e.data == nil {
		return types.Intervals{}, nil
	}
	return e.layout.ToIntervals(e.data)
}
