package interval

import "math"

type EQType string

const (
	EQType1 EQType = "[]"
	EQType2 EQType = "(]"
	EQType3 EQType = "[)"
	EQType4 EQType = "()"
)

/**
 * @author: yasinWu
 * @date: 2022/1/13 14:45
 * @description: Interval
 */
type Interval struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}

/**
 * @author: yasinWu
 * @date: 2022/1/13 14:45
 * @params: interval *Interval
 * @return: bool
 * @description: 判断interval是否与该区间有交集
 */
func (i *Interval) IntervalMixed(interval *Interval) bool {
	startMax := math.Max(float64(i.Start), float64(interval.Start))
	endMin := math.Min(float64(i.End), float64(interval.End))
	return startMax <= endMin
}

/**
 * @author: yasinWu
 * @date: 2022/1/13 14:46
 * @params: i int64, eq EQType
 * @return: bool
 * @description: 判断i是否属于该区间
 */
func (i *Interval) Belong(value int64, eq EQType) bool {
	switch eq {
	case EQType1:
		return value >= i.Start && value <= i.End
	case EQType2:
		return value > i.Start && value <= i.End
	case EQType3:
		return value >= i.Start && value < i.End
	case EQType4:
		return value > i.Start && value < i.End
	default:
		return value >= i.Start && value <= i.End
	}
}

/**
 * @author: yasinWu
 * @date: 2022/1/13 14:46
 * @params: interval *Interval, eq EQType
 * @return: bool
 * @description: 判断interval是否是该区间子集
 */
func (i *Interval) Contain(interval *Interval, eq EQType) bool {
	switch eq {
	case EQType1:
		return interval.Start >= i.Start && interval.End <= i.End
	case EQType2:
		return interval.Start > i.Start && interval.End <= i.End
	case EQType3:
		return interval.Start >= i.Start && interval.End < i.End
	case EQType4:
		return interval.Start > i.Start && interval.End < i.End
	default:
		return interval.Start >= i.Start && interval.End <= i.End
	}
}
