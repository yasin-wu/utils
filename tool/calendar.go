package tool

import (
	"errors"
	"strconv"
	"time"
)

type Calendar struct {
	Calendar time.Time `json:"calendar"`
	FromTime time.Time `json:"from_time"`
	ToTime   time.Time `json:"to_time"`
	Usage    float64   `json:"usage"`
}

/**
 * @author: yasinWu
 * @date: 2022/1/13 14:32
 * @params: startTime, endTime time.Time, timeFormatTpl string
 * @return: []*Calendar, error
 * @description: 切割时间区间为天
 */
func CalendarDays(startTime, endTime time.Time, timeFormatTpl string) ([]*Calendar, error) {
	if endTime.Before(startTime) || endTime.Equal(startTime) {
		return nil, errors.New("startTime <= endTime")
	}
	if timeFormatTpl == "" {
		timeFormatTpl = "2006-01-02"
	}
	days := GetBetweenDates(startTime, endTime, timeFormatTpl)
	var data []*Calendar
	for i, v := range days {
		vt, _ := time.ParseInLocation("2006-01-02", v, time.Now().Location())
		var fromTime time.Time
		var toTime time.Time
		switch i {
		case 0:
			fromTime = startTime
			toTime = time.Date(vt.Year(), vt.Month(), vt.Day(), 23, 59, 59, 0, time.Now().Location())
		case len(days) - 1:
			fromTime = time.Date(vt.Year(), vt.Month(), vt.Day(), 0, 0, 0, 0, time.Now().Location())
			toTime = endTime
		default:
			fromTime = time.Date(vt.Year(), vt.Month(), vt.Day(), 0, 0, 0, 0, time.Now().Location())
			toTime = time.Date(vt.Year(), vt.Month(), vt.Day(), 23, 59, 59, 0, time.Now().Location())
		}
		usageTime, _ := strconv.ParseFloat(strconv.FormatFloat(toTime.Sub(fromTime).Hours(), 'f', 1, 64), 64)
		calendar := &Calendar{
			Calendar: vt,
			FromTime: fromTime,
			ToTime:   toTime,
			Usage:    usageTime,
		}
		data = append(data, calendar)
	}
	return data, nil
}
