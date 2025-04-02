package ktime

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertToZoneTimeStr(t *testing.T) {
	// 测试用例的时间戳：2023-05-15 12:40:45 UTC
	// 2023-05-15 20:40:45 UTC+8
	timestamp := int64(1684154445)

	// 测试用例
	tests := []struct {
		name        string
		timestamp   int64
		format      string
		zone        []string
		expected    string
		expectError bool
	}{
		{
			name:      "默认时区UTC+8",
			timestamp: timestamp,
			format:    "2006-01-02 15:04:05",
			zone:      []string{},
			expected:  "2023-05-15 20:40:45",
		},
		{
			name:      "使用数字时区+8",
			timestamp: timestamp,
			format:    "2006-01-02 15:04:05",
			zone:      []string{"+8"},
			expected:  "2023-05-15 20:40:45",
		},
		{
			name:      "使用数字时区-5",
			timestamp: timestamp,
			format:    "2006-01-02 15:04:05",
			zone:      []string{"-5"},
			expected:  "2023-05-15 07:40:45",
		},
		{
			name:      "使用标准时区名称",
			timestamp: timestamp,
			format:    "2006-01-02 15:04:05",
			zone:      []string{"UTC"},
			expected:  "2023-05-15 12:40:45",
		},
		{
			name:      "使用不同格式",
			timestamp: timestamp,
			format:    "2006/01/02 15:04:05 -07:00",
			zone:      []string{"+8"},
			expected:  "2023/05/15 20:40:45 +08:00",
		},
		{
			name:      "使用简短格式",
			timestamp: timestamp,
			format:    "15:04:05",
			zone:      []string{"+8"},
			expected:  "20:40:45",
		},
		{
			name:        "无效的时区格式",
			timestamp:   timestamp,
			format:      "2006-01-02 15:04:05",
			zone:        []string{"invalid_zone"},
			expectError: true,
		},
		{
			name:      "零时间戳",
			timestamp: 0,
			format:    "2006-01-02 15:04:05",
			zone:      []string{"+8"},
			expected:  "1970-01-01 08:00:00",
		},
		{
			name:      "负时间戳",
			timestamp: -86400, // 减去一天
			format:    "2006-01-02 15:04:05",
			zone:      []string{"+8"},
			expected:  "1969-12-31 08:00:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ConvertToZoneTimeStr(tt.timestamp, tt.format, tt.zone...)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
