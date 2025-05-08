package daytime

import (
	"database/sql/driver"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	Day         = 24 * time.Hour
	DefaultTime = "00:00"
)

var (
	daytimeRegex = regexp.MustCompile(`^\d\d:\d\d(:\d\d){0,1}$`)

	ErrObjIsNil   = errors.New("object is nil")
	ErrInvalid    = errors.New("invalid")
	ErrUnexpected = errors.New("unexpected")
)

type DayTime struct {
	hour   int
	minute int
	second int
}

// New create a new daytime.
func New(hour int, minute int, second int) (DayTime, error) {
	if hour < 0 || hour > 23 {
		return DayTime{}, errors.Wrap(ErrInvalid, fmt.Sprintf("value of hour is %d", hour))
	}
	if minute < 0 || minute > 59 {
		return DayTime{}, errors.Wrap(ErrInvalid, fmt.Sprintf("value of minute is %d", minute))
	}
	if second < 0 || second > 59 {
		return DayTime{}, errors.Wrap(ErrInvalid, fmt.Sprintf("value of second is %d", second))
	}

	return DayTime{
		hour:   hour,
		minute: minute,
		second: second,
	}, nil
}

// Parse parse a daytime.
func Parse(value string) (DayTime, error) {
	value = strings.Trim(string(value), " \t")
	submatches := daytimeRegex.FindAllStringSubmatch(value, -1)

	if len(submatches) == 0 {
		return DayTime{}, errors.Wrap(ErrInvalid, fmt.Sprintf("value '%s'", value))
	}

	values := strings.Split(value, ":")
	hour, err := strconv.Atoi(values[0])
	if err != nil {
		return DayTime{}, errors.Wrap(err, "hour")
	}

	minute, err := strconv.Atoi(values[1])
	if err != nil {
		return DayTime{}, errors.Wrap(err, "minute")
	}

	second := 0
	if len(values) > 2 {
		second, err = strconv.Atoi(values[2])
	}
	if err != nil {
		return DayTime{}, errors.Wrap(err, "second")
	}

	return New(hour, minute, second)
}

// String convert to string.
func (t *DayTime) String() string {
	if t == nil {
		return DefaultTime
	}

	hour := strconv.Itoa(t.hour)
	if len(hour) == 1 {
		hour = "0" + hour
	}

	minute := strconv.Itoa(t.minute)
	if len(minute) == 1 {
		minute = "0" + minute
	}

	value := hour + ":" + minute
	if t.second == 0 {
		return value
	}

	second := strconv.Itoa(t.second)
	if len(second) == 1 {
		second = "0" + second
	}

	return value + ":" + second
}

// Time bringing to the current day's time.
func (t *DayTime) Time() time.Time {
	now := time.Now()
	year, month, day := now.Date()
	hour := 0
	minute := 0
	second := 0
	if t != nil {
		hour = t.hour
		minute = t.minute
		second = t.second
	}

	return time.Date(
		year,
		month,
		day,
		hour,
		minute,
		second,
		0,
		now.Location(),
	)
}

// InTheNearFuture bringing to the current day's time.
func (t *DayTime) InTheNearFuture() time.Time {
	now := time.Now()
	datetime := t.Time()

	if datetime.Before(now) {
		datetime = datetime.Add(Day)
	}

	return datetime
}

// InTheRecentPast bringing to the near future.
func (t *DayTime) InTheRecentPast() time.Time {
	now := time.Now()
	datetime := t.Time()

	if datetime.After(now) {
		datetime = datetime.Add(-Day)
	}

	return datetime
}

func (t *DayTime) MarshalBinary() ([]byte, error) {
	return []byte(t.String()), nil
}

func (t *DayTime) UnmarshalBinary(data []byte) error {
	if t == nil {
		return ErrObjIsNil
	}

	value, err := Parse(string(data))
	if err != nil {
		return errors.Wrap(err, "parse")
	}

	*t = value

	return nil
}

func (t *DayTime) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}

func (t *DayTime) UnmarshalText(data []byte) error {
	if t == nil {
		return ErrObjIsNil
	}

	value, err := Parse(string(data))
	if err != nil {
		return errors.Wrap(err, "parse")
	}

	*t = value

	return nil
}

func (t *DayTime) MarshalJSON() ([]byte, error) {
	return []byte(t.String()), nil
}

func (t *DayTime) UnmarshalJSON(data []byte) error {
	if t == nil {
		return ErrObjIsNil
	}

	value, err := Parse(string(data))
	if err != nil {
		return errors.Wrap(err, "parse")
	}

	*t = value

	return nil
}

func (t *DayTime) MarshalCSV() (string, error) {
	return t.String(), nil
}

func (t *DayTime) UnmarshalCSV(str string) error {
	if t == nil {
		return ErrObjIsNil
	}

	value, err := Parse(str)
	if err != nil {
		return errors.Wrap(err, "parse")
	}

	*t = value

	return nil
}

func (t *DayTime) Scan(src any) error {
	if t == nil {
		return ErrObjIsNil
	}

	str := ""
	switch src := src.(type) {
	case []byte:
		str = string(src)
		fmt.Printf("str(bytes) = '%s'\n", str)
	case string:
		str = src
	default:
		return errors.Wrap(ErrUnexpected, fmt.Sprintf("type of value '%T'", src))
	}

	value, err := Parse(str)
	if err != nil {
		return errors.Wrap(err, "parse")
	}

	*t = value

	return nil
}

func (t *DayTime) Value() (driver.Value, error) {
	return t.String(), nil
}
