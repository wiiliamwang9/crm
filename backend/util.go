package main

import (
	"fmt"
	"time"
)

// joinStrings 用指定分隔符连接字符串数组
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}

	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}

// parseCommaSeparatedInt64 解析逗号分隔的int64数组
func parseCommaSeparatedInt64(str string) []int64 {
	if str == "" {
		return nil
	}

	var result []int64
	parts := splitString(str, ",")
	for _, part := range parts {
		part = trimSpace(part)
		if part != "" {
			if id, err := parseInt64(part); err == nil {
				result = append(result, id)
			}
		}
	}
	return result
}

// splitString 分割字符串
func splitString(s, sep string) []string {
	if s == "" {
		return nil
	}
	// 简单实现字符串分割
	var result []string
	start := 0
	for i := 0; i < len(s); i++ {
		if i+len(sep)-1 < len(s) && s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
			i += len(sep) - 1
		}
	}
	result = append(result, s[start:])
	return result
}

// trimSpace 去除字符串前后空格
func trimSpace(s string) string {
	start := 0
	end := len(s)

	// 去除前面的空格
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}

	// 去除后面的空格
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}

	return s[start:end]
}

// parseInt64 解析字符串为int64
func parseInt64(s string) (int64, error) {
	if s == "" {
		return 0, fmt.Errorf("empty string")
	}

	var result int64
	var negative bool
	start := 0

	if s[0] == '-' {
		negative = true
		start = 1
	} else if s[0] == '+' {
		start = 1
	}

	for i := start; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return 0, fmt.Errorf("invalid character")
		}
		result = result*10 + int64(s[i]-'0')
	}

	if negative {
		result = -result
	}

	return result, nil
}

// generatePreferenceID 生成偏好项唯一ID
func generatePreferenceID() string {
	return fmt.Sprintf("pref_%d", getCurrentTimestamp())
}

// getCurrentTimestamp 获取当前时间戳（毫秒）
func getCurrentTimestamp() int64 {
	return getCurrentTime().UnixNano() / 1000000
}

// getCurrentTime 获取当前时间
func getCurrentTime() time.Time {
	return time.Now()
}

// getStringFromMap 从map中安全获取字符串值
func getStringFromMap(m map[string]interface{}, key string) string {
	if value, ok := m[key]; ok {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return ""
}

// getTimeFromMap 从map中安全获取时间值
func getTimeFromMap(m map[string]interface{}, key string) time.Time {
	if value, ok := m[key]; ok {
		if timeStr, ok := value.(string); ok {
			if t, err := time.Parse(time.RFC3339, timeStr); err == nil {
				return t
			}
		}
		if t, ok := value.(time.Time); ok {
			return t
		}
	}
	return time.Time{}
}
