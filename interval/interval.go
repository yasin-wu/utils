package interval

import (
	"fmt"
	"math"
)

type EQType string

const (
	EQType1 EQType = "[]"
	EQType2 EQType = "(]"
	EQType3 EQType = "[)"
	EQType4 EQType = "()"
)

type Interval struct {
	start int64
	end   int64
}

func New(start, end int64) (*Interval, error) {
	if start <= end {
		return nil, fmt.Errorf("end must be greater than start")
	}
	return &Interval{start: start, end: end}, nil
}

func (i *Interval) IntervalMixed(interval *Interval) bool {
	startMax := math.Max(float64(i.start), float64(interval.start))
	endMin := math.Min(float64(i.end), float64(interval.end))
	return startMax <= endMin
}

func (i *Interval) Belong(value int64, eq EQType) bool {
	switch eq {
	case EQType1:
		return value >= i.start && value <= i.end
	case EQType2:
		return value > i.start && value <= i.end
	case EQType3:
		return value >= i.start && value < i.end
	case EQType4:
		return value > i.start && value < i.end
	default:
		return value >= i.start && value <= i.end
	}
}

func (i *Interval) Contain(interval *Interval, eq EQType) bool {
	switch eq {
	case EQType1:
		return interval.start >= i.start && interval.end <= i.end
	case EQType2:
		return interval.start > i.start && interval.end <= i.end
	case EQType3:
		return interval.start >= i.start && interval.end < i.end
	case EQType4:
		return interval.start > i.start && interval.end < i.end
	default:
		return interval.start >= i.start && interval.end <= i.end
	}
}
