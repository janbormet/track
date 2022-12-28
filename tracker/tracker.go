package tracker

import (
	"time"
	"track/types"
)

type Context struct {
	Storage types.Storage
	Time    time.Time
}

func Open(time time.Time, tags []types.Tag) types.OpenInterval {
	return types.OpenInterval{
		Start: time,
		Tags:  tags,
	}
}

func Start(ctx Context, tags []types.Tag) error {
	intervals, err := ctx.Storage.Load()
	if err != nil {
		return err
	}
	intervals.OpenIntervals = append(intervals.OpenIntervals, Open(ctx.Time, tags))

	return ctx.Storage.Save(intervals)
}

func Stop(ctx Context) error {
	intervals, err := ctx.Storage.Load()
	if err != nil {
		return err
	}
	for _, openInterval := range intervals.OpenIntervals {
		intervals.Intervals = append(intervals.Intervals, openInterval.Close(ctx.Time))
	}
	intervals.OpenIntervals = []types.OpenInterval{}
	return ctx.Storage.Save(intervals)
}
