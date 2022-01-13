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
 * @author: yasin
 * @date: 2022/1/13 14:45
 * @description: Interval
 */
type Interval struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}

/**
 * @author: yasin
 * @date: 2022/1/13 14:45
 * @params: interval *Interval
 * @return: bool
 * @description: 判断interval是否与该区间有交集
 */
func (this *Interval) IntervalMixed(interval *Interval) bool {
	startMax := math.Max(float64(this.Start), float64(interval.Start))
	endMin := math.Min(float64(this.End), float64(interval.End))
	return startMax <= endMin
}

/**
 * @author: yasin
 * @date: 2022/1/13 14:46
 * @params: i int64, eq EQType
 * @return: bool
 * @description: 判断i是否属于该区间
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

/**
 * @author: yasin
 * @date: 2022/1/13 14:46
 * @params: interval *Interval, eq EQType
 * @return: bool
 * @description: 判断interval是否是该区间子集
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
