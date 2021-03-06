/*
 * Copyright (C) 2018 The Seer Authors. All rights reserved.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package stream_test

import (
	"math"
	"testing"

	"github.com/cshenton/seer/stream"

	"github.com/cshenton/seer/dist/uv"
)

func TestToLogNormal(t *testing.T) {
	tt := []struct {
		name     string
		loc      float64
		scale    float64
		logLoc   float64
		logScale float64
	}{
		{"simple", 1, 1, -0.3465735903, 0.8325546112},
		{"less simple", 100, 10, 4.6001950206, 0.09975134512},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			n, err := uv.NewNormal(tc.loc, tc.scale)
			if err != nil {
				t.Fatal("unexpected error while creating Normal,", err)
			}
			ln, err := stream.ToLogNormal(n)
			if math.Abs(tc.logLoc-ln.Location) > 1e-8 {
				t.Errorf("expected new location %v, but got %v", tc.logLoc, ln.Location)
			}
			if math.Abs(tc.logScale-ln.Scale) > 1e-8 {
				t.Errorf("expected new scale %v, but got %v", tc.logScale, ln.Scale)
			}
		})
	}
}

func TestToLogNormalErrs(t *testing.T) {
	tt := []struct {
		name  string
		loc   float64
		scale float64
	}{
		{"zero location", 0, 1},
		{"negative location", -1, 1},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			n, err := uv.NewNormal(tc.loc, tc.scale)
			if err != nil {
				t.Fatal("unexpected error while creating Normal,", err)
			}
			ln, err := stream.ToLogNormal(n)
			if err == nil {
				t.Error("expected error, but it was nil")
			}
			if ln != nil {
				t.Error("expected nil pointer, but got", ln)
			}
		})
	}
}
