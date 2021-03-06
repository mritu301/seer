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

package kalman_test

import (
	"testing"

	"github.com/cshenton/seer/kalman"
	"gonum.org/v1/gonum/mat"
)

func TestPredict(t *testing.T) {
	tt := []struct {
		name   string
		k      int
		a      []float64
		b      []float64
		c      []float64
		q      []float64
		r      []float64
		locIn  []float64
		covIn  []float64
		locOut []float64
		covOut []float64
	}{
		{
			"Identity 1x1", 1, []float64{1}, []float64{1}, []float64{1}, []float64{1}, []float64{1},
			[]float64{1}, []float64{1}, []float64{1}, []float64{2},
		},
		{
			"Identity 2x2", 2, []float64{1, 0, 0, 1}, []float64{1, 0, 0, 1}, []float64{1, 1}, []float64{1, 0, 0, 1}, []float64{1},
			[]float64{1, 1}, []float64{1, 0, 0, 1}, []float64{1, 1}, []float64{2, 0, 0, 2},
		},
		{
			"Non-trivial 2x2", 2, []float64{1, 1, 0, 1}, []float64{1, 0, 0, 1}, []float64{1, 0}, []float64{0.5, 0.1, 0.1, 1.0}, []float64{0.5},
			[]float64{1, 1}, []float64{1, 0, 0, 1}, []float64{2, 1}, []float64{2.5, 1.1, 1.1, 2},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.k*tc.k != len(tc.a) {
				t.Fatalf("Expected %v process datapoints, got %v", tc.k*tc.k, len(tc.a))
			}
			a := mat.NewDense(tc.k, tc.k, tc.a)
			b := mat.NewDense(tc.k, tc.k, tc.b)
			c := mat.NewDense(1, tc.k, tc.c)
			q := mat.NewDense(tc.k, tc.k, tc.q)
			r := mat.NewDense(1, 1, tc.r)
			m, err := kalman.NewSystem(a, b, c, q, r)
			if err != nil {
				t.Fatal("failed to create kalman.System", err)
			}
			prev, err := kalman.NewState(mat.NewDense(tc.k, 1, tc.locIn), mat.NewDense(tc.k, tc.k, tc.covIn))
			if err != nil {
				t.Fatal("failed to create kalman.State", err)
			}
			next, err := kalman.Predict(prev, m)
			if err != nil {
				t.Fatal(err)
			}
			if !mat.Equal(next.Loc, mat.NewDense(tc.k, 1, tc.locOut)) {
				t.Errorf("Expected location vals %v, got %v", mat.NewDense(tc.k, 1, tc.locOut), next.Loc)
			}
			if !mat.Equal(next.Cov, mat.NewDense(tc.k, tc.k, tc.covOut)) {
				t.Errorf("Expected covariance vals %v, got %v", mat.NewDense(tc.k, tc.k, tc.covOut), next.Cov)
			}
		})
	}
}

func TestObserve(t *testing.T) {
	tt := []struct {
		name   string
		k      int
		a      []float64
		b      []float64
		c      []float64
		q      []float64
		r      []float64
		locIn  []float64
		covIn  []float64
		locOut []float64
		covOut []float64
	}{
		{
			"Identity 1x1", 1, []float64{1}, []float64{1}, []float64{1}, []float64{1}, []float64{1},
			[]float64{1}, []float64{1}, []float64{1}, []float64{2},
		},
		{
			"Identity 2x2", 2, []float64{1, 0, 0, 1}, []float64{1, 0, 0, 1}, []float64{1, 1}, []float64{1, 0, 0, 1}, []float64{1},
			[]float64{1, 1}, []float64{1, 0, 0, 1}, []float64{2}, []float64{3},
		},
		{
			"Non-trivial 2x2", 2, []float64{1, 1, 0, 1}, []float64{1, 0, 0, 1}, []float64{1, 0}, []float64{0.5, 0.1, 0.1, 1.0}, []float64{0.5},
			[]float64{1, 1}, []float64{1, 0, 0, 1}, []float64{1}, []float64{1.5},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.k*tc.k != len(tc.a) {
				t.Fatalf("Expected %v process datapoints, got %v", tc.k*tc.k, len(tc.a))
			}
			a := mat.NewDense(tc.k, tc.k, tc.a)
			b := mat.NewDense(tc.k, tc.k, tc.b)
			c := mat.NewDense(1, tc.k, tc.c)
			q := mat.NewDense(tc.k, tc.k, tc.q)
			r := mat.NewDense(1, 1, tc.r)
			m, err := kalman.NewSystem(a, b, c, q, r)
			prev, err := kalman.NewState(mat.NewDense(tc.k, 1, tc.locIn), mat.NewDense(tc.k, tc.k, tc.covIn))
			if err != nil {
				t.Fatal("failed to create System", err)
			}
			obs, err := kalman.Observe(prev, m)
			if err != nil {
				t.Fatal("failed to create observation", err)
			}

			if !mat.Equal(obs.Loc, mat.NewDense(1, 1, tc.locOut)) {
				t.Errorf("Expected Loc vals %v, got %v", mat.NewDense(1, 1, tc.locOut), obs.Loc)
			}
			if !mat.Equal(obs.Cov, mat.NewDense(1, 1, tc.covOut)) {
				t.Errorf("Expected Cov vals %v, got %v", mat.NewDense(1, 1, tc.covOut), obs.Cov)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	tt := []struct {
		name   string
		k      int
		a      []float64
		b      []float64
		c      []float64
		q      []float64
		r      []float64
		locIn  []float64
		covIn  []float64
		locOut []float64
		covOut []float64
	}{
		{
			"Identity 1x1", 1, []float64{1}, []float64{1}, []float64{1}, []float64{1}, []float64{1},
			[]float64{1}, []float64{1}, []float64{0.66666}, []float64{0.5},
		},
		{
			"Identity 2x2", 2, []float64{1, 0, 0, 1}, []float64{1, 0, 0, 1}, []float64{1, 1}, []float64{1, 0, 0, 1}, []float64{1},
			[]float64{1, 1}, []float64{1, 0, 0, 1}, []float64{0.44444, 0.44444}, []float64{0.66666, -0.33333, -0.33333, 0.66666},
		},
		{
			"Non-trivial 2x2", 2, []float64{1, 1, 0, 1}, []float64{1, 0, 0, 1}, []float64{1, 0}, []float64{0.5, 0.1, 0.1, 1.0}, []float64{0.5},
			[]float64{1, 1}, []float64{1, 0, 0, 1}, []float64{0.55555, 1}, []float64{0.33333, 0, 0, 1},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.k*tc.k != len(tc.a) {
				t.Fatalf("Expected %v process datapoints, got %v", tc.k*tc.k, len(tc.a))
			}
			a := mat.NewDense(tc.k, tc.k, tc.a)
			b := mat.NewDense(tc.k, tc.k, tc.b)
			c := mat.NewDense(1, tc.k, tc.c)
			q := mat.NewDense(tc.k, tc.k, tc.q)
			r := mat.NewDense(1, 1, tc.r)
			m, err := kalman.NewSystem(a, b, c, q, r)
			pre, err := kalman.NewState(mat.NewDense(tc.k, 1, tc.locIn), mat.NewDense(tc.k, tc.k, tc.covIn))
			if err != nil {
				t.Fatal("failed to create System", err)
			}
			post, _, err := kalman.Update(pre, m, 1.0/3.0)
			if err != nil {
				t.Fatal(err)
			}
			if !mat.EqualApprox(post.Loc, mat.NewDense(tc.k, 1, tc.locOut), 1e-3) {
				t.Errorf("Expected Loc vals %v, got %v", mat.NewDense(tc.k, 1, tc.locOut), post.Loc)
			}
			if !mat.EqualApprox(post.Cov, mat.NewDense(tc.k, tc.k, tc.covOut), 1e-3) {
				t.Errorf("Expected Cov vals %v, got %v", mat.NewDense(tc.k, tc.k, tc.covOut), post.Cov)
			}
		})
	}
}
