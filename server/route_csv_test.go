package main

import (
	"strings"
	"testing"
)

var csvTestCases = []struct {
	Case   string
	Expect []Point
}{
	{
		Case: `112.6371,-7.95993
112.63715,-7.95982
112.63748,-7.95881
112.63764,-7.95825
112.63801,-7.95673`,
		Expect: []Point{
			{Longitude: 112.6371, Latitude: -7.95993},
			{Longitude: 112.63715, Latitude: -7.95982},
			{Longitude: 112.63748, Latitude: -7.95881},
			{Longitude: 112.63764, Latitude: -7.95825},
			{Longitude: 112.63801, Latitude: -7.95673},
		},
	},
	{
		Case: `112.63947,-7.95014
112.63962,-7.94944
112.63966,-7.94929`,
		Expect: []Point{
			{Longitude: 112.63947, Latitude: -7.95014},
			{Longitude: 112.63962, Latitude: -7.94944},
			{Longitude: 112.63966, Latitude: -7.94929},
		},
	},
}

func TestReadRouteCSV(t *testing.T) {
	r := createCSVReader()

	for _, tc := range csvTestCases {
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
