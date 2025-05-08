package daytime

import (
	"database/sql/driver"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Parallel()

	type args struct {
		hour   int
		minute int
		second int
	}
	type expectedResult struct {
		daytime DayTime
		err     error
	}
	tests := []struct {
		name           string
		args           args
		expectedResult expectedResult
	}{
		{
			name: "Checking standard work",
			args: args{
				hour:   1,
				minute: 2,
				second: 3,
			},
			expectedResult: expectedResult{
				daytime: DayTime{
					hour:   1,
					minute: 2,
					second: 3,
				},
				err: nil,
			},
		},
		{
			name: "Checking the processing of an invalid hour value",
			args: args{
				hour:   24,
				minute: 2,
				second: 3,
			},
			expectedResult: expectedResult{
				daytime: DayTime{},
				err:     ErrInvalid,
			},
		},
		{
			name: "Checking the processing of an invalid minute value",
			args: args{
				hour:   1,
				minute: 60,
				second: 3,
			},
			expectedResult: expectedResult{
				daytime: DayTime{},
				err:     ErrInvalid,
			},
		},
		{
			name: "Checking the processing of an invalid second value",
			args: args{
				hour:   1,
				minute: 2,
				second: 60,
			},
			expectedResult: expectedResult{
				daytime: DayTime{},
				err:     ErrInvalid,
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()

			daytime, err := New(test.args.hour, test.args.minute, test.args.second)
			assert.EqualValues(tt, test.expectedResult.daytime, daytime)
			assert.ErrorIs(tt, err, test.expectedResult.err)
		})
	}
}

func TestParse(t *testing.T) {
	t.Parallel()

	type args struct {
		value string
	}
	type expectedResult struct {
		daytime DayTime
		err     error
	}
	tests := []struct {
		name           string
		args           args
		expectedResult expectedResult
	}{
		{
			name: "Checking standard work",
			args: args{
				value: "00:01:02",
			},
			expectedResult: expectedResult{
				daytime: DayTime{
					hour:   0,
					minute: 1,
					second: 2,
				},
				err: nil,
			},
		},
		{
			name: "Checking the processing of an invalid value",
			args: args{
				value: "",
			},
			expectedResult: expectedResult{
				daytime: DayTime{},
				err:     ErrInvalid,
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()

			daytime, err := Parse(test.args.value)
			assert.EqualValues(tt, test.expectedResult.daytime, daytime)
			assert.ErrorIs(tt, err, test.expectedResult.err)
		})
	}
}

func TestString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		daytime        *DayTime
		expectedResult string
	}{
		{
			name: "Checking to get the full value",
			daytime: &DayTime{
				hour:   1,
				minute: 2,
				second: 3,
			},
			expectedResult: "01:02:03",
		},
		{
			name:           "Checking to get the default value",
			daytime:        nil,
			expectedResult: "00:00",
		},
		{
			name: "Checking to get the value without second",
			daytime: &DayTime{
				hour:   1,
				minute: 2,
				second: 0,
			},
			expectedResult: "01:02",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()

			value := test.daytime.String()
			assert.EqualValues(tt, test.expectedResult, value)
		})
	}
}

func TestTime(t *testing.T) {
	t.Parallel()

	now := time.Now()
	year := now.Year()
	month := now.Month()
	day := now.Day()
	location := now.Location()

	tests := []struct {
		name           string
		daytime        *DayTime
		expectedResult time.Time
	}{
		{
			name: "Checking standard work",
			daytime: &DayTime{
				hour:   1,
				minute: 2,
				second: 3,
			},
			expectedResult: time.Date(
				year,
				month,
				day,
				1,
				2,
				3,
				0,
				location,
			),
		},
		{
			name:    "Checking the nil value",
			daytime: nil,
			expectedResult: time.Date(
				year,
				month,
				day,
				0,
				0,
				0,
				0,
				location,
			),
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()

			value := test.daytime.Time()
			assert.EqualValues(tt, test.expectedResult, value)
		})
	}
}

func TestInTheNearFuture(t *testing.T) {
	t.Parallel()

	now := time.Now()
	hour := now.Hour()
	minute := now.Minute()
	second := now.Second()
	nanosecond := now.Nanosecond()
	location := now.Location()
	nextHour := now.Add(time.Hour)
	nextDay := now.Add(Day)
	time.Sleep(time.Millisecond)

	tests := []struct {
		name           string
		daytime        *DayTime
		expectedResult time.Time
	}{
		{
			name: "Checking to get the current day",
			daytime: &DayTime{
				hour:   nextHour.Hour(),
				minute: minute,
				second: second,
			},
			expectedResult: time.Date(
				now.Year(),
				now.Month(),
				now.Day(),
				nextHour.Hour(),
				minute,
				second,
				nanosecond,
				location,
			),
		},
		{
			name: "Checking to get the next day",
			daytime: &DayTime{
				hour:   hour,
				minute: minute,
				second: second,
			},
			expectedResult: time.Date(
				nextDay.Year(),
				nextDay.Month(),
				nextDay.Day(),
				hour,
				minute,
				second,
				nanosecond,
				location,
			),
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()

			value := test.daytime.InTheNearFuture()
			assert.EqualValues(tt, test.expectedResult.Year(), value.Year())
			assert.EqualValues(tt, test.expectedResult.Month(), value.Month())
			assert.EqualValues(tt, test.expectedResult.Day(), value.Day())
			assert.EqualValues(tt, test.expectedResult.Hour(), value.Hour())
			assert.EqualValues(tt, test.expectedResult.Minute(), value.Minute())
			assert.EqualValues(tt, test.expectedResult.Second(), value.Second())
		})
	}
}

func TestInTheRecentPast(t *testing.T) {
	t.Parallel()

	now := time.Now()
	hour := now.Hour()
	minute := now.Minute()
	second := now.Second()
	nanosecond := now.Nanosecond()
	nextHour := now.Add(time.Hour)
	previousDay := now.Add(-Day)
	location := now.Location()
	time.Sleep(time.Millisecond)

	tests := []struct {
		name           string
		daytime        *DayTime
		expectedResult time.Time
	}{
		{
			name: "Checking to get the current day",
			daytime: &DayTime{
				hour:   hour,
				minute: minute,
				second: second,
			},
			expectedResult: time.Date(
				now.Year(),
				now.Month(),
				now.Day(),
				hour,
				minute,
				second,
				nanosecond,
				location,
			),
		},
		{
			name: "Checking to get the previous day",
			daytime: &DayTime{
				hour:   nextHour.Hour(),
				minute: minute,
				second: second,
			},
			expectedResult: time.Date(
				now.Year(),
				now.Month(),
				previousDay.Day(),
				nextHour.Hour(),
				minute,
				second,
				nanosecond,
				location,
			),
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()

			value := test.daytime.InTheRecentPast()
			assert.EqualValues(tt, test.expectedResult.Year(), value.Year())
			assert.EqualValues(tt, test.expectedResult.Month(), value.Month())
			assert.EqualValues(tt, test.expectedResult.Day(), value.Day())
			assert.EqualValues(tt, test.expectedResult.Hour(), value.Hour())
			assert.EqualValues(tt, test.expectedResult.Minute(), value.Minute())
			assert.EqualValues(tt, test.expectedResult.Second(), value.Second())
		})
	}
}

func TestMarshalBinary(t *testing.T) {
	t.Parallel()

	type expectedResult struct {
		value []byte
		err   error
	}
	tests := []struct {
		name           string
		daytime        *DayTime
		expectedResult expectedResult
	}{
		{
			name: "Checking standard work",
			daytime: &DayTime{
				hour:   1,
				minute: 2,
				second: 3,
			},
			expectedResult: expectedResult{
				value: []byte("01:02:03"),
				err:   nil,
			},
		},
		{
			name:    "Checking to process nil",
			daytime: nil,
			expectedResult: expectedResult{
				value: []byte("00:00"),
				err:   nil,
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()

			value, err := test.daytime.MarshalBinary()
			assert.EqualValues(tt, test.expectedResult.value, value)
			assert.EqualValues(tt, test.expectedResult.err, err)
		})
	}
}

func TestUnmarshalBinary(t *testing.T) {
	t.Parallel()

	type args struct {
		data []byte
	}
	type expectedResult struct {
		daytime *DayTime
		err     error
	}
	tests := []struct {
		name           string
		daytime        *DayTime
		args           args
		expectedResult expectedResult
	}{
		{
			name:    "Checking standard work",
			daytime: &DayTime{},
			args: args{
				data: []byte("01:02:03"),
			},
			expectedResult: expectedResult{
				daytime: &DayTime{
					hour:   1,
					minute: 2,
					second: 3,
				},
				err: nil,
			},
		},
		{
			name:    "Checking to process parse error",
			daytime: &DayTime{},
			args: args{
				data: []byte("24:02:03"),
			},
			expectedResult: expectedResult{
				daytime: &DayTime{},
				err:     ErrInvalid,
			},
		},
		{
			name:    "Checking to process nil",
			daytime: nil,
			args: args{
				data: []byte("01:02:03"),
			},
			expectedResult: expectedResult{
				daytime: nil,
				err:     ErrObjIsNil,
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()

			err := test.daytime.UnmarshalBinary(test.args.data)
			assert.EqualValues(tt, test.expectedResult.daytime, test.daytime)
			assert.ErrorIs(tt, err, test.expectedResult.err)
		})
	}
}

func TestMarshalText(t *testing.T) {
	t.Parallel()

	type expectedResult struct {
		value []byte
		err   error
	}
	tests := []struct {
		name           string
		daytime        *DayTime
		expectedResult expectedResult
	}{
		{
			name: "Checking standard work",
			daytime: &DayTime{
				hour:   1,
				minute: 2,
				second: 3,
			},
			expectedResult: expectedResult{
				value: []byte("01:02:03"),
				err:   nil,
			},
		},
		{
			name:    "Checking to process nil",
			daytime: nil,
			expectedResult: expectedResult{
				value: []byte("00:00"),
				err:   nil,
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()

			value, err := test.daytime.MarshalText()
			assert.EqualValues(tt, test.expectedResult.value, value)
			assert.EqualValues(tt, test.expectedResult.err, err)
		})
	}
}

func TestUnmarshalText(t *testing.T) {
	t.Parallel()

	type args struct {
		data []byte
	}
	type expectedResult struct {
		daytime *DayTime
		err     error
	}
	tests := []struct {
		name           string
		daytime        *DayTime
		args           args
		expectedResult expectedResult
	}{
		{
			name:    "Checking standard work",
			daytime: &DayTime{},
			args: args{
				data: []byte("01:02:03"),
			},
			expectedResult: expectedResult{
				daytime: &DayTime{
					hour:   1,
					minute: 2,
					second: 3,
				},
				err: nil,
			},
		},
		{
			name:    "Checking to process parse error",
			daytime: &DayTime{},
			args: args{
				data: []byte("24:02:03"),
			},
			expectedResult: expectedResult{
				daytime: &DayTime{},
				err:     ErrInvalid,
			},
		},
		{
			name:    "Checking to process nil",
			daytime: nil,
			args: args{
				data: []byte("01:02:03"),
			},
			expectedResult: expectedResult{
				daytime: nil,
				err:     ErrObjIsNil,
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()

			err := test.daytime.UnmarshalText(test.args.data)
			assert.EqualValues(tt, test.expectedResult.daytime, test.daytime)
			assert.ErrorIs(tt, err, test.expectedResult.err)
		})
	}
}

func TestMarshalJSON(t *testing.T) {
	t.Parallel()

	type expectedResult struct {
		value []byte
		err   error
	}
	tests := []struct {
		name           string
		daytime        *DayTime
		expectedResult expectedResult
	}{
		{
			name: "Checking standard work",
			daytime: &DayTime{
				hour:   1,
				minute: 2,
				second: 3,
			},
			expectedResult: expectedResult{
				value: []byte("01:02:03"),
				err:   nil,
			},
		},
		{
			name:    "Checking to process nil",
			daytime: nil,
			expectedResult: expectedResult{
				value: []byte("00:00"),
				err:   nil,
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()

			value, err := test.daytime.MarshalJSON()
			assert.EqualValues(tt, test.expectedResult.value, value)
			assert.EqualValues(tt, test.expectedResult.err, err)
		})
	}
}

func TestUnmarshalJSON(t *testing.T) {
	t.Parallel()

	type args struct {
		data []byte
	}
	type expectedResult struct {
		daytime *DayTime
		err     error
	}
	tests := []struct {
		name           string
		daytime        *DayTime
		args           args
		expectedResult expectedResult
	}{
		{
			name:    "Checking standard work",
			daytime: &DayTime{},
			args: args{
				data: []byte("01:02:03"),
			},
			expectedResult: expectedResult{
				daytime: &DayTime{
					hour:   1,
					minute: 2,
					second: 3,
				},
				err: nil,
			},
		},
		{
			name:    "Checking to process parse error",
			daytime: &DayTime{},
			args: args{
				data: []byte("24:02:03"),
			},
			expectedResult: expectedResult{
				daytime: &DayTime{},
				err:     ErrInvalid,
			},
		},
		{
			name:    "Checking to process nil",
			daytime: nil,
			args: args{
				data: []byte("01:02:03"),
			},
			expectedResult: expectedResult{
				daytime: nil,
				err:     ErrObjIsNil,
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()

			err := test.daytime.UnmarshalJSON(test.args.data)
			assert.EqualValues(tt, test.expectedResult.daytime, test.daytime)
			assert.ErrorIs(tt, err, test.expectedResult.err)
		})
	}
}

func TestMarshalCSV(t *testing.T) {
	t.Parallel()

	type expectedResult struct {
		value string
		err   error
	}
	tests := []struct {
		name           string
		daytime        *DayTime
		expectedResult expectedResult
	}{
		{
			name: "Checking standard work",
			daytime: &DayTime{
				hour:   1,
				minute: 2,
				second: 3,
			},
			expectedResult: expectedResult{
				value: "01:02:03",
				err:   nil,
			},
		},
		{
			name:    "Checking to process nil",
			daytime: nil,
			expectedResult: expectedResult{
				value: "00:00",
				err:   nil,
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()

			value, err := test.daytime.MarshalCSV()
			assert.EqualValues(tt, test.expectedResult.value, value)
			assert.EqualValues(tt, test.expectedResult.err, err)
		})
	}
}

func TestUnmarshalCSV(t *testing.T) {
	t.Parallel()

	type args struct {
		str string
	}
	type expectedResult struct {
		daytime *DayTime
		err     error
	}
	tests := []struct {
		name           string
		daytime        *DayTime
		args           args
		expectedResult expectedResult
	}{
		{
			name:    "Checking standard work",
			daytime: &DayTime{},
			args: args{
				str: "01:02:03",
			},
			expectedResult: expectedResult{
				daytime: &DayTime{
					hour:   1,
					minute: 2,
					second: 3,
				},
				err: nil,
			},
		},
		{
			name:    "Checking to process parse error",
			daytime: &DayTime{},
			args: args{
				str: "24:02:03",
			},
			expectedResult: expectedResult{
				daytime: &DayTime{},
				err:     ErrInvalid,
			},
		},
		{
			name:    "Checking to process nil",
			daytime: nil,
			args: args{
				str: "01:02:03",
			},
			expectedResult: expectedResult{
				daytime: nil,
				err:     ErrObjIsNil,
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()

			err := test.daytime.UnmarshalCSV(test.args.str)
			assert.EqualValues(tt, test.expectedResult.daytime, test.daytime)
			assert.ErrorIs(tt, err, test.expectedResult.err)
		})
	}
}

func TestScan(t *testing.T) {
	t.Parallel()

	type args struct {
		src any
	}
	type expectedResult struct {
		daytime *DayTime
		err     error
	}
	tests := []struct {
		name           string
		daytime        *DayTime
		args           args
		expectedResult expectedResult
	}{
		{
			name:    "Checking to process string",
			daytime: &DayTime{},
			args: args{
				src: "01:02:03",
			},
			expectedResult: expectedResult{
				daytime: &DayTime{
					hour:   1,
					minute: 2,
					second: 3,
				},
				err: nil,
			},
		},
		{
			name:    "Checking to process bytes",
			daytime: &DayTime{},
			args: args{
				src: []byte("01:02:03"),
			},
			expectedResult: expectedResult{
				daytime: &DayTime{
					hour:   1,
					minute: 2,
					second: 3,
				},
				err: nil,
			},
		},
		{
			name:    "Checking to process nil",
			daytime: nil,
			args: args{
				src: "01:02:03",
			},
			expectedResult: expectedResult{
				daytime: nil,
				err:     ErrObjIsNil,
			},
		},
		{
			name:    "Checking to process unexpected type",
			daytime: &DayTime{},
			args: args{
				src: map[int]string{1: "01:02:03"},
			},
			expectedResult: expectedResult{
				daytime: &DayTime{},
				err:     ErrUnexpected,
			},
		},
		{
			name:    "Checking to process parse error",
			daytime: &DayTime{},
			args: args{
				src: "24:02:03",
			},
			expectedResult: expectedResult{
				daytime: &DayTime{},
				err:     ErrInvalid,
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()

			err := test.daytime.Scan(test.args.src)
			assert.EqualValues(tt, test.expectedResult.daytime, test.daytime)
			assert.ErrorIs(tt, err, test.expectedResult.err)
		})
	}
}

func TestValue(t *testing.T) {
	t.Parallel()

	type expectedResult struct {
		value driver.Value
		err   error
	}
	tests := []struct {
		name           string
		daytime        *DayTime
		expectedResult expectedResult
	}{
		{
			name: "Checking standard work",
			daytime: &DayTime{
				hour:   1,
				minute: 2,
				second: 3,
			},
			expectedResult: expectedResult{
				value: "01:02:03",
				err:   nil,
			},
		},
		{
			name:    "Checking to process nil",
			daytime: nil,
			expectedResult: expectedResult{
				value: "00:00",
				err:   nil,
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()

			value, err := test.daytime.Value()
			assert.EqualValues(tt, test.expectedResult.value, value)
			assert.EqualValues(tt, test.expectedResult.err, err)
		})
	}
}
