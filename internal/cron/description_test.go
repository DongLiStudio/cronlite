package cron

import (
	"strings"
	"testing"
)

func TestDescribe(t *testing.T) {
	for _, tc := range []struct{ expression, want string }{
		{"* * * * *", "每分钟"},
		{"*/5 * * * *", "每 5 分钟"},
		{"0/5 * * * *", "每 5 分钟"},
		{"*/15 * * * *", "每 15 分钟"},
		{"0 * * * *", "每小时"},
		{"15 * * * *", "每小时第 15 分钟"},
		{"0 */2 * * *", "每 2 小时"},
		{"30 */6 * * *", "每天 00:30、06:30、12:30、18:30"},
		{"0 9 * * *", "每天 09:00"},
		{"30 2 * * *", "每天 02:30"},
		{"0 9 * * 1-5", "每天 09:00 每周一–五"},
		{"0 9 * * MON-FRI", "每天 09:00 每周一–五"},
		{"0 12 15 * *", "每天 12:00 每月第 15 天"},
		{"0 0 1 JAN *", "每天 00:00 每月第 1 天，月份：1月"},
		{"0 9 1 * MON", "每天 09:00 每月第 1 天或每周一"},
		{"0 9 1-31 * MON", "每天 09:00"},
		{"0 9 1 * 0-6", "每天 09:00"},
		{"0 9 */2 * MON", "每天 09:00 每月第 1, 3, 5, 7, 9, 11, 13, 15, 17, 19, 21, 23, 25, 27, 29, 31 天或每周一"},
		{"*/7 * * * *", "每小时的第 0, 7, 14, 21, 28, 35, 42, 49, 56 分钟"},
		{"5/15 * * * *", "每小时的第 5, 20, 35, 50 分钟"},
		{"0,15 9-11 * MAR-MAY MON,WED,FRI", "在 09–11 时的第 0, 15 分钟 每周一, 三, 五，月份：3月–5月"},
		{"CRON_TZ=America/New_York 0 9 * * *", "每天 09:00 (America/New_York)"},
		{"0 0 31 2 *", "每天 00:00 每月第 31 天，月份：2月"},
		{"bad", "无效计划"},
		{"@daily", "无效计划"},
		{"TZ=", "无效计划"},
		{"CRON_TZ=UTC", "无效计划"},
	} {
		t.Run(tc.expression, func(t *testing.T) {
			if got := Describe(tc.expression); got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestParserRejectsIncompleteTimezonePrefix(t *testing.T) {
	for _, expression := range []string{"TZ=", "TZ=UTC", "CRON_TZ=", "CRON_TZ=UTC", "TZ= * * * * *"} {
		if _, err := NewParser().Parse(expression, "UTC"); err == nil {
			t.Fatalf("accepted malformed cron %q", expression)
		}
	}
}

func FuzzDescribe(f *testing.F) {
	for _, expression := range []string{"*/5 * * * *", "0 9 1 * MON", "", "CRON_TZ=UTC * * * * *", "* * * * *"} {
		f.Add(expression)
	}
	f.Fuzz(func(t *testing.T, expression string) {
		text := Describe(expression)
		if text == "" || strings.Contains(text, "<") {
			t.Fatalf("invalid description %q", text)
		}
	})
}
