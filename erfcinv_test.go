package main

import (
	"math"
	"testing"
)

func tolerance(a, b, e float64) bool {
	// Multiplying by e here can underflow denormal values to zero.
	// Check a==b so that at least if a and b are small and identical
	// we say they match.
	if a == b {
		return true
	}
	d := a - b
	if d < 0 {
		d = -d
	}

	// note: b is correct (expected) value, a is actual value.
	// make error tolerance a fraction of b, not a.
	if b != 0 {
		e = e * b
		if e < 0 {
			e = -e
		}
	}
	return d < e
}

func close(a, b float64) bool     { return tolerance(a, b, 1e-14) }
func veryclose(a, b float64) bool { return tolerance(a, b, 4e-16) }

// func soclose(a, b, e float64) bool { return tolerance(a, b, e) }
func alike(a, b float64) bool {
	switch {
	case math.IsNaN(a) && math.IsNaN(b):
		return true
	case a == b:
		return math.Signbit(a) == math.Signbit(b)
	}
	return false
}

var vf = []float64{
	4.9790119248836735e+00,
	7.7388724745781045e+00,
	-2.7688005719200159e-01,
	-5.0106036182710749e+00,
	9.6362937071984173e+00,
	2.9263772392439646e+00,
	5.2290834314593066e+00,
	2.7279399104360102e+00,
	1.8253080916808550e+00,
	-8.6859247685756013e+00,
}
var erfinv = []float64{
	4.746037673358033586786350696e-01,
	8.559054432692110956388764172e-01,
	-2.45427830571707336251331946e-02,
	-4.78116683518973366268905506e-01,
	1.479804430319470983648120853e+00,
	2.654485787128896161882650211e-01,
	5.027444534221520197823192493e-01,
	2.466703532707627818954585670e-01,
	1.632011465103005426240343116e-01,
	-1.06672334642196900710000389e+00,
}

var erfcinvSC = []float64{
	math.Inf(+1),
	math.Inf(-1),
	0,
	math.NaN(),
	math.NaN(),
	math.NaN(),
}

var vferfcinvSC = []float64{
	0,
	2,
	1,
	math.Inf(1),
	math.Inf(-1),
	math.NaN(),
}

func TestErfcinv(t *testing.T) {
	for i := 0; i < len(vf); i++ {
		a := 1.0 - (vf[i] / 10)
		if f := Erfcinv(a); !veryclose(erfinv[i], f) {
			t.Errorf("Erfcinv(%g) = %g, want %g", a, f, erfinv[i])
		}
	}
	for i := 0; i < len(vferfcinvSC); i++ {
		if f := Erfcinv(vferfcinvSC[i]); !alike(erfcinvSC[i], f) {
			t.Errorf("Erfcinv(%g) = %g, want %g", vferfcinvSC[i], f, erfcinvSC[i])
		}
	}
	for x := 0.1; x <= 1.9; x += 1e-2 {
		if f := math.Erfc(Erfcinv(x)); !close(x, f) {
			t.Errorf("Erfc(Erfcinv(%g)) = %g, want %g", x, f, x)
		}
	}
	for x := 0.1; x <= 1.9; x += 1e-2 {
		if f := Erfcinv(math.Erfc(x)); !close(x, f) {
			t.Errorf("Erfcinv(Erfc(%g)) = %g, want %g", x, f, x)
		}
	}
}
