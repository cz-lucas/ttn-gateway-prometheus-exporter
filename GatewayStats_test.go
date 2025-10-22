package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoundTripTimes_ConvertToSeconds(t *testing.T) {
	t.Run("With all valid values", func(t *testing.T) {

		rtt := RoundTripTimes{
			Min:    "0.2342343454s",
			Median: "0.4356764s",
			Max:    "0.87654435s",
			Count:  42,
		}

		min, median, max, err := rtt.ConvertToSeconds()
		assert.InDelta(t, min, 0.2342343454, 0.0000001)
		assert.InDelta(t, median, 0.4356764, 0.0000001)
		assert.InDelta(t, max, 0.87654435, 0.0000001)
		assert.Nil(t, err)
	})

	t.Run("All valid except min", func(t *testing.T) {

		rtt := RoundTripTimes{
			Min:    "i am a string",
			Median: "0.4356764s",
			Max:    "0.87654435s",
			Count:  42,
		}

		min, median, max, err := rtt.ConvertToSeconds()
		assert.Equal(t, min, -1.0)
		assert.Equal(t, median, -1.0)
		assert.Equal(t, max, -1.0)

		expected := "error parsing min duration: time: invalid duration \"i am a string\""
		assert.EqualError(t, err, expected)
	})

	t.Run("All valid except median", func(t *testing.T) {

		rtt := RoundTripTimes{
			Min:    "0.23457s",
			Median: "another string",
			Max:    "0.87654435s",
			Count:  42,
		}

		min, median, max, err := rtt.ConvertToSeconds()
		assert.Equal(t, min, -1.0)
		assert.Equal(t, median, -1.0)
		assert.Equal(t, max, -1.0)

		expected := "error parsing median duration: time: invalid duration \"another string\""
		assert.EqualError(t, err, expected)
	})

	t.Run("All valid except max", func(t *testing.T) {

		rtt := RoundTripTimes{
			Min:    "0.734565234s",
			Median: "0.4356764s",
			Max:    "and another one",
			Count:  42,
		}

		min, median, max, err := rtt.ConvertToSeconds()
		assert.Equal(t, min, -1.0)
		assert.Equal(t, median, -1.0)
		assert.Equal(t, max, -1.0)

		expected := "error parsing max duration: time: invalid duration \"and another one\""
		assert.EqualError(t, err, expected)
	})
}

func Test_convertDurationToSeconds(t *testing.T) {
	t.Run("Valid time 1", func(t *testing.T) {

		result, err := convertDurationToSeconds("0.2343455635s")
		assert.InDelta(t, result, 0.2343455635, 0.0000001)
		assert.Nil(t, err)
	})

	t.Run("Valid time 2", func(t *testing.T) {

		result, err := convertDurationToSeconds("5.9876543243456456s")
		assert.InDelta(t, result, 5.9876543243456456, 0.0000001)
		assert.Nil(t, err)
	})

	t.Run("Invalid time", func(t *testing.T) {

		result, err := convertDurationToSeconds("fiftyoneseconds")
		assert.Equal(t, result, -1.0)

		expected := "time: invalid duration \"fiftyoneseconds\""
		assert.EqualError(t, err, expected)
	})
}

func Test_stringsToFloat64(t *testing.T) {
	t.Run("Valid float 1", func(t *testing.T) {

		result, err := stringsToFloat64("012.2334534534545")
		assert.InDelta(t, result, 12.2334534534545, 0.0000001)
		assert.Nil(t, err)

	})

	t.Run("Valid float 2", func(t *testing.T) {

		result, err := stringsToFloat64("2345345.324353")
		assert.InDelta(t, result, 2345345.324353, 0.001)
		assert.Nil(t, err)

	})

	t.Run("Invalid float", func(t *testing.T) {

		result, err := stringsToFloat64("three.onefour")
		assert.Equal(t, result, -1.0)

		expected := "strconv.ParseFloat: parsing \"three.onefour\": invalid syntax"
		assert.EqualError(t, err, expected)
	})
}

func TestGatewayStats_GetUplinkCount(t *testing.T) {
	t.Run("Valid value", func(t *testing.T) {

		gwStats := GatewayStats{
			UplinkCount: "1337",
		}

		uplinkCount, err := gwStats.GetUplinkCount()
		assert.Equal(t, uplinkCount, 1337.0)
		assert.Nil(t, err)
	})

	t.Run("Invalid value", func(t *testing.T) {

		gwStats := GatewayStats{
			UplinkCount: "invalid",
		}

		uplinkCount, err := gwStats.GetUplinkCount()
		assert.Equal(t, uplinkCount, -1.0)

		expected := "error parsing uplinkCount: strconv.ParseFloat: parsing \"invalid\": invalid syntax"
		assert.EqualError(t, err, expected)
	})
}
