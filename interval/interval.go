package interval

import "math"

type EQType string

const (
	EQType1 EQType = "[]"
	EQType2 EQType = "(]"
	EQType3 EQType = "[)"
	EQType4 EQType = "()"
)

type Interval struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}

// IntervalMixed
/**
 * @author: yasin
 * @date: 2021/4/2 15:36
 * @description：判断两个区间是否存在交集,true is mixed
 */
func (this *Interval) IntervalMixed(interval *Interval) bool {
	startMax := math.Max(float64(this.Start), float64(interval.Start))
	endMin := math.Min(float64(this.End), float64(interval.End))
	return startMax <= endMin
}

// Belong
/**
 * @author: yasin
 * @date: 2021/4/22 13:14
 * @description：i 是否属于interval,取等
 */
func (this *Interval) Belong(i int64, eq EQType) bool {
	switch eq {
	case EQType1:
		return i >= this.Start && i <= this.End
	case EQType2:
		return i > this.Start && i <= this.End
	case EQType3:
		return i >= this.Start && i < this.End
	case EQType4:
		return i > this.Start && i < this.End
	default:
		return i >= this.Start && i <= this.End
	}
}

// Contain
/**
 * @author: yasin
 * @date: 2021/4/22 14:58
 * @description：interval是否this的子区间
 */
func (this *Interval) Contain(interval *Interval, eq EQType) bool {
	switch eq {
	case EQType1:
		return interval.Start >= this.Start && interval.End <= this.End
	case EQType2:
		return interval.Start > this.Start && interval.End <= this.End
	case EQType3:
		return interval.Start >= this.Start && interval.End < this.End
	case EQType4:
		return interval.Start > this.Start && interval.End < this.End
	default:
		return interval.Start >= this.Start && interval.End <= this.End
	}
}
