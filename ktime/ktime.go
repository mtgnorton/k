package ktime

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

var (
	ErrInvalidZoneFormat = errors.New("invalid zone format")
)

const (
	DefaultZone       = "UTC+8"
	DefaultZoneOffset = 8 * 60 * 60
)

// ConvertToZoneTimeStr 将时间戳转换为指定时区和格式的时间字符串
// 参数:
//   - timestamp: Unix时间戳（秒）
//   - format: 时间格式化模板，如 "2006-01-02 15:04:05"
//   - zone: 可选的时区参数，支持标准时区名称（如"UTC"）或小时偏移（如"+8"、"-5"）
//
// 返回值:
//   - string: 格式化后的时间字符串
//   - error: 如果提供了无效的时区格式，返回 ErrInvalidZoneFormat 错误
//
// 注意事项:
//   - 如果不提供时区参数，默认使用 UTC+8
//   - 时区参数可以是标准时区名称或小时偏移格式
//
// 示例:
//
//	str, err := ConvertToZoneTimeStr(1684154445, "2006-01-02 15:04:05", "+8")
//	// 返回: "2023-05-15 20:40:45", nil
func ConvertToZoneTimeStr(timestamp int64, format string, zone ...string) (string, error) {
	// 将时间戳转换为 time.Time 对象
	t := time.Unix(timestamp, 0)

	// 解析时区
	var location *time.Location
	if len(zone) == 0 {
		// 默认使用 UTC+8
		location = time.FixedZone(DefaultZone, DefaultZoneOffset)
	} else {
		zoneStr := zone[0]
		// 尝试解析标准时区
		loc, err := time.LoadLocation(zoneStr)
		if err != nil {
			// 如果不是标准时区，尝试解析为小时偏移
			// 例如 "+8" 或 "-5"
			var offset int
			_, err := fmt.Sscanf(zoneStr, "%d", &offset)
			if err == nil {
				location = time.FixedZone(fmt.Sprintf("UTC%+d", offset), offset*60*60)
			} else {
				// 解析失败，返回错误
				return "", ErrInvalidZoneFormat
			}
		} else {
			location = loc
		}
	}

	// 将时间转换为指定时区
	tInZone := t.In(location)

	// 格式化时间字符串并返回
	return tInZone.Format(format), nil
}

// MustConvertToZoneTimeStr 将时间戳转换为指定时区和格式的时间字符串，如果转换失败，会 panic,
// 参考 ConvertToZoneTimeStr
func MustConvertToZoneTimeStr(timestamp int64, format string, zone ...string) string {
	str, err := ConvertToZoneTimeStr(timestamp, format, zone...)
	if err != nil {
		panic(err)
	}
	return str
}
