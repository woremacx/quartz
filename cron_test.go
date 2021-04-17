package quartz_test

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/woremacx/quartz"
)

const (
	readDateLayout = "Mon Jan 2 15:04:05 2006"
)

// for debug
var verbose = false

// on original(reugn/go-quartz), it was 1555351200000000000
// $ date -d @1555351200
// Mon Apr 15 18:00:00 UTC 2019
var PREV_1 = int64TimeFromStringByLocal("Mon Apr 15 18:00:00 2019")

// on original(reugn/go-quartz), it was 1555524000000000000
// $ date -d @1555524000
// Wed Apr 17 18:00:00 UTC 2019
var PREV_2 = int64TimeFromStringByLocal("Wed Apr 17 18:00:00 2019")

func int64TimeFromStringByLocal(value string) int64 {
	localTimeFromString, err := time.ParseInLocation(readDateLayout, value, time.Local)
	if err != nil {
		panic(err)
	}
	return localTimeFromString.UnixNano()
}

func TestCronExpression1(t *testing.T) {
	prev := PREV_1
	if verbose {
		t.Logf("prev: %s", getTimeFromInt64(prev))
	}
	result := ""
	//                                         0     1  2  3    4 5 6
	//                                         s     m  h  d    m w y
	cronTrigger, err := quartz.NewCronTrigger("10/20 15 14 5-10 * ? *")
	cronTrigger.Description()
	if err != nil {
		t.Fatal(err)
	} else {
		result, _ = iterate(prev, cronTrigger, 1000)
	}
	assertEqual(t, result, "Fri Dec 8 14:15:10 2023")
}

func TestCronExpression2(t *testing.T) {
	prev := PREV_1
	if verbose {
		t.Logf("prev: %s", getTimeFromInt64(prev))
	}
	result := ""
	//                                         0 1     2     3 4 5 6
	//                                         s min   hour  d m w y
	cronTrigger, err := quartz.NewCronTrigger("* 5,7,9 14-16 * * ? *")
	if err != nil {
		t.Fatal(err)
	} else {
		result, _ = iterate(prev, cronTrigger, 1000)
	}
	assertEqual(t, result, "Mon Aug 5 14:05:00 2019")
}

func TestCronExpression3(t *testing.T) {
	prev := PREV_1
	if verbose {
		t.Logf("prev: %s", getTimeFromInt64(prev))
	}
	result := ""
	//       0 1     2    3 4 5       6
	//       s m     h    d m w       y
	expr := "* 5,7,9 14/2 * * Wed,Sat *"
	cronTrigger, err := quartz.NewCronTrigger(expr)
	if err != nil {
		t.Fatal(err)
	} else {
		result, _ = iterate(prev, cronTrigger, 1000)
	}
	assertEqual(t, result, "Sat Dec 7 14:05:00 2019")
}

func TestCronExpression4(t *testing.T) {
	expression := "0 5,7 14 1 * Sun *"
	_, err := quartz.NewCronTrigger(expression)
	if err == nil {
		t.Fatalf("%s should fail", expression)
	}
}

func TestCronExpression5(t *testing.T) {
	prev := PREV_1
	if verbose {
		t.Logf("prev: %s", getTimeFromInt64(prev))
	}
	result := ""
	cronTrigger, err := quartz.NewCronTrigger("* * * * * ? *")
	if err != nil {
		t.Fatal(err)
	} else {
		result, _ = iterate(prev, cronTrigger, 1000)
	}
	assertEqual(t, result, "Mon Apr 15 18:16:40 2019")
}

func TestCronExpression6(t *testing.T) {
	prev := PREV_1
	if verbose {
		t.Logf("prev: %s", getTimeFromInt64(prev))
	}
	result := ""
	cronTrigger, err := quartz.NewCronTrigger("* * 14/2 * * Mon/3 *")
	if err != nil {
		t.Fatal(err)
	} else {
		result, _ = iterate(prev, cronTrigger, 1000)
	}
	assertEqual(t, result, "Mon Mar 15 18:00:00 2021")
}

func TestCronExpression7(t *testing.T) {
	prev := PREV_1
	if verbose {
		t.Logf("prev: %s", getTimeFromInt64(prev))
	}
	result := ""
	cronTrigger, err := quartz.NewCronTrigger("* 5-9 14/2 * * 0-2 *")
	if err != nil {
		t.Fatal(err)
	} else {
		result, _ = iterate(prev, cronTrigger, 1000)
	}
	assertEqual(t, result, "Tue Jul 16 16:09:00 2019")
}

func TestCronDaysOfWeek(t *testing.T) {
	daysOfWeek := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	expected := []string{
		"Sun Apr 21 00:00:00 2019",
		"Mon Apr 22 00:00:00 2019",
		"Tue Apr 23 00:00:00 2019",
		"Wed Apr 24 00:00:00 2019",
		"Thu Apr 18 00:00:00 2019",
		"Fri Apr 19 00:00:00 2019",
		"Sat Apr 20 00:00:00 2019",
	}

	for i := 0; i < len(daysOfWeek); i++ {
		cronDayOfWeek(t, daysOfWeek[i], expected[i])
		cronDayOfWeek(t, strconv.Itoa(i), expected[i])
	}
}

func cronDayOfWeek(t *testing.T, dayOfWeek, expected string) {
	prev := PREV_2
	if verbose {
		t.Logf("prev: %s", getTimeFromInt64(prev))
	}
	expression := fmt.Sprintf("0 0 0 * * %s", dayOfWeek)
	cronTrigger, err := quartz.NewCronTrigger(expression)
	if err != nil {
		t.Fatal(err)
	} else {
		nextFireTime, err := cronTrigger.NextFireTime(prev)
		if err != nil {
			t.Fatal(err)
		} else {
			assertEqual(t, getTimeFromInt64(nextFireTime), expected)
		}
	}
}

func TestCronYearly(t *testing.T) {
	prev := PREV_1
	if verbose {
		t.Logf("prev: %s", getTimeFromInt64(prev))
	}
	result := ""
	cronTrigger, err := quartz.NewCronTrigger("@yearly")
	if err != nil {
		t.Fatal(err)
	} else {
		result, _ = iterate(prev, cronTrigger, 100)
	}
	assertEqual(t, result, "Sun Jan 1 00:00:00 2119")
}

func TestCronMonthly(t *testing.T) {
	prev := PREV_1
	if verbose {
		t.Logf("prev: %s", getTimeFromInt64(prev))
	}
	result := ""
	cronTrigger, err := quartz.NewCronTrigger("@monthly")
	if err != nil {
		t.Fatal(err)
	} else {
		result, _ = iterate(prev, cronTrigger, 100)
	}
	assertEqual(t, result, "Sun Aug 1 00:00:00 2027")
}

func TestCronWeekly(t *testing.T) {
	prev := PREV_1
	if verbose {
		t.Logf("prev: %s", getTimeFromInt64(prev))
	}
	result := ""
	cronTrigger, err := quartz.NewCronTrigger("@weekly")
	if err != nil {
		t.Fatal(err)
	} else {
		result, _ = iterate(prev, cronTrigger, 100)
	}
	assertEqual(t, result, "Sun Mar 14 00:00:00 2021")
}

func TestCronDaily(t *testing.T) {
	prev := PREV_1
	if verbose {
		t.Logf("prev: %s", getTimeFromInt64(prev))
	}
	result := ""
	cronTrigger, err := quartz.NewCronTrigger("@daily")
	if err != nil {
		t.Fatal(err)
	} else {
		result, _ = iterate(prev, cronTrigger, 1000)
	}
	assertEqual(t, result, "Sun Jan 9 00:00:00 2022")
}

func TestCronHourly(t *testing.T) {
	prev := PREV_1
	if verbose {
		t.Logf("prev: %s", getTimeFromInt64(prev))
	}
	result := ""
	cronTrigger, err := quartz.NewCronTrigger("@hourly")
	if err != nil {
		t.Fatal(err)
	} else {
		result, _ = iterate(prev, cronTrigger, 1000)
	}
	assertEqual(t, result, "Wed May 29 06:00:00 2019")
}

func iterate(prev int64, cronTrigger *quartz.CronTrigger, iterations int) (string, error) {
	var err error
	for i := 0; i < iterations; i++ {
		prev, err = cronTrigger.NextFireTime(prev)
		if verbose {
			fmt.Println(time.Unix(prev/int64(time.Second), 0).Local().Format(readDateLayout))
		}
		if err != nil {
			return "", err
		}
	}
	return time.Unix(prev/int64(time.Second), 0).Local().Format(readDateLayout), nil
}

func getTimeFromInt64(value int64) string {
	return time.Unix(value/int64(time.Second), 0).Local().Format(readDateLayout)
}