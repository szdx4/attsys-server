package common

import "time"

// ParseTime 解析时间字符串
func ParseTime(timeStr string) (time.Time, error) {
	location, _ := time.LoadLocation("Asia/Shanghai")
	t, err := time.ParseInLocation(time.RFC3339, timeStr, location)
	if err != nil {
		return time.Now(), err
	}
	return t, nil
}
