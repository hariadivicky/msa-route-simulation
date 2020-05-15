package main

import (
	"strings"
	"testing"
)

var jsonTestCases = []struct {
	Case   string
	Expect []Point
}{
	{
		Case: `[{"lng": 112.6371, "lat": -7.95993},
		{"lng": 112.63715, "lat": -7.95982},
		{"lng": 112.63748, "lat": -7.95881},
		{"lng": 112.63764, "lat": -7.95825},
		{"lng": 112.63801, "lat": -7.95673}]`,
		Expect: []Point{
			{Longitude: 112.6371, Latitude: -7.95993},
			{Longitude: 112.63715, Latitude: -7.95982},
			{Longitude: 112.63748, Latitude: -7.95881},
			{Longitude: 112.63764, Latitude: -7.95825},
			{Longitude: 112.63801, Latitude: -7.95673},
		},
	},
	{
		Case: `[{"lng": 112.63947, "lat": -7.95014},
		{"lng": 112.63962, "lat": -7.94944},
		{"lng": 112.63966, "lat": -7.94929}]`,
		Expect: []Point{
			{Longitude: 112.63947, Latitude: -7.95014},
			{Longitude: 112.63962, Latitude: -7.94944},
			{Longitude: 112.63966, Latitude: -7.94929},
		},
	},
}

func TestReadRouteJSON(t *testing.T) {
	r := createJSONReader()

	for _, tc := range jsonTestCases {
		result, err := r.Read(strings.NewReader(tc.Case))
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		// compare with expected
		if result == nil {
			t.Error("result is null")
			t.FailNow()
		}

		if len(result) != len(tc.Expect) {
			t.Error("result length not same with expected length")
			t.FailNow()
		}

		for i := 0; i < len(tc.Expect); i++ {
			if result[i].Longitude != tc.Expect[i].Longitude {
				t.Errorf("expect %f but result is %f", tc.Expect[i].Longitude, result[i].Longitude)
				t.FailNow()
			}

			if result[i].Latitude != tc.Expect[i].Latitude {
				t.Errorf("expect %f but result is %f", tc.Expect[i].Latitude, result[i].Latitude)
				t.FailNow()
			}
		}
	}
}
