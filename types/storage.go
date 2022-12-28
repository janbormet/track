package types

type Storage interface {
	Save(Intervals) error
	Load() (Intervals, error)
}
