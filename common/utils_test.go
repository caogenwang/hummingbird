package hummingbird

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestParseRange(t *testing.T) {
    //Setting up individual test data
    tests := []struct {
        rangeHeader string
        exResBegin  int64
        exResEnd    int64
        exError     string
    }{
        {" ", 0, 0, ""},
        {"bytes=", 0, 0, "invalid range format"},
        {"bytes=-", 0, 0, "invalid range format"},
        {"bytes=-cv", 0, 0, "invalid end with no begin"},
        {"bytes=cv-", 0, 0, "invalid begin with no end"},
        {"bytes=-0", 0, 0, "zero end with no begin"},
        {"bytes=-12346", 0, 12345, ""},
        {"bytes=-12344", 1, 12345, ""},
        {"bytes=12344-", 12344, 12345, ""},
        {"bytes=12345-cv", 0, 0, "invalid end"},
        {"bytes=cv-12345", 0, 0, "invalid begin"},
        {"bytes=12346-12", 0, 0, "end before begin"},
        {"bytes=12346-123457", 0, 0, "Begin bigger than file"},
        {"bytes=12342-12343", 12342, 12344, ""},
        {"bytes=12342-12344", 12342, 12345, ""},
    }

    //Run tests with data from above
    for _, test := range tests {
        result, err := ParseRange(test.rangeHeader, 12345)
        if test.rangeHeader == " " {
            assert.Nil(t, result)
            assert.Nil(t, err)
            continue
        }
        if test.exError == "" {
            httpResult := httpRange{test.exResBegin, test.exResEnd}
            assert.Nil(t, err)
            assert.Contains(t, result, httpResult)
        } else {
            assert.Equal(t, err.Error(), test.exError)
        }
    }
}

func TestParseRange_NoEnd_BeginLargerThanFilesize(t *testing.T) {
    result, err := ParseRange("bytes=12346-", 12345)
    assert.Nil(t, err)
    assert.Empty(t, result)
}

func TestParseDate(t *testing.T) {
    //Setup tests with individual data
    tests := []string{
        "Mon, 02 Jan 2006 15:04:05 MST",
        "Mon, 02 Jan 2006 15:04:05 -0700",
        "Mon Jan 02 15:04:05 2006",
        "Monday, 02-Jan-06 15:04:05 MST",
        "1136214245",
        "2006-01-02 15:04:05",
    }

    //Run Tests from above
    for _, timestamp := range tests {
        timeResult, err := ParseDate(timestamp)
        if err == nil {
            assert.Equal(t, timeResult.Day(), 02)
            assert.Equal(t, timeResult.Month(), 01)
            assert.Equal(t, timeResult.Year(), 2006)
            assert.Equal(t, timeResult.Hour(), 15)
            assert.Equal(t, timeResult.Minute(), 04)
            assert.Equal(t, timeResult.Second(), 05)
        } else {
            assert.Equal(t, err.Error(), "invalid time")
        }
    }

}

func TestParseTimestamp(t *testing.T) {
    tests := []string{
        "2006-01-02 15:04:05",
        "Mon, 02 Jan 2006 15:04:05 MST",
    }

    for _, timestamp := range tests {
        timeResult, err := FormatTimestamp(timestamp)
        if err != nil {
            assert.Equal(t, err.Error(), "invalid time")
            assert.Empty(t, timeResult)
        }else{
            assert.Equal(t, "2006-01-02T15:04:05", timeResult)
        }
    }
}

func TestLooksTrue(t *testing.T) {
    tests := []string{
        "true ",
        "true",
        "t",
        "yes",
        "y",
        "1",
        "on",
    }

    for _, test := range tests {
        isTrue := LooksTrue(test)
        assert.True(t, isTrue)
    }
}
