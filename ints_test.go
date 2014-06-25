package ints

import (
	//"math"
	"math/rand"
	"strconv"
	"testing"
)

const (
	SMALL  = 10
	MEDIUM = 1000
	LARGE  = 100000
	HUGE   = 10000000
)

func Panics(fun func()) (b bool) {
	defer func() {
		err := recover()
		if err != nil {
			b = true
		}
	}()
	fun()
	return
}

func AreSlicesEqual(t *testing.T, truth, comp []int, str string) {
	if !Equal(comp, truth) {
		t.Errorf(str+". Expected %v, returned %v", truth, comp)
	}
}

func TestAdd(t *testing.T) {
	a := []int{1, 2, 3}
	b := []int{4, 5, 6}
	c := []int{7, 8, 9}
	truth := []int{12, 15, 18}
	n := make([]int, len(a))
	Add(n, a, b, c)
	AreSlicesEqual(t, truth, n, "Wrong addition of slices new receiver")
	Add(a, b, c)
	AreSlicesEqual(t, truth, n, "Wrong addition of slices for no new receiver")
	// Test that it panics
	if !Panics(func() { Add(make([]int, 2), make([]int, 3)) }) {
		t.Errorf("Did not panic with length mismatch")
	}
}

func TestAddConst(t *testing.T) {
	s := []int{3, 4, 1, 7, 5}
	c := 6
	truth := []int{9, 10, 7, 13, 11}
	AddConst(c, s)
	AreSlicesEqual(t, truth, s, "Wrong addition of constant")
}

func TestAddScaled(t *testing.T) {
	s := []int{3, 4, 1, 7, 5}
	alpha := 6
	dst := []int{1, 2, 3, 4, 5}
	ans := []int{19, 26, 9, 46, 35}
	AddScaled(dst, alpha, s)
	if !Equal(dst, ans) {
		t.Errorf("Adding scaled did not match")
	}
	short := []int{1}
	if !Panics(func() { AddScaled(dst, alpha, short) }) {
		t.Errorf("Doesn't panic if s is smaller than dst")
	}
	if !Panics(func() { AddScaled(short, alpha, s) }) {
		t.Errorf("Doesn't panic if dst is smaller than s")
	}
}

func TestAddScaledTo(t *testing.T) {
	s := []int{3, 4, 1, 7, 5}
	alpha := 6
	y := []int{1, 2, 3, 4, 5}
	dst := make([]int, 5)
	ans := []int{19, 26, 9, 46, 35}
	AddScaledTo(dst, y, alpha, s)
	if !Equal(dst, ans) {
		t.Errorf("AddScaledTo did not match")
	}
	AddScaledTo(dst, y, alpha, s)
	if !Equal(dst, ans) {
		t.Errorf("Reusing dst did not match")
	}
	short := []int{1}
	if !Panics(func() { AddScaledTo(dst, y, alpha, short) }) {
		t.Errorf("Doesn't panic if s is smaller than dst")
	}
	if !Panics(func() { AddScaledTo(short, y, alpha, s) }) {
		t.Errorf("Doesn't panic if dst is smaller than s")
	}
	if !Panics(func() { AddScaledTo(dst, short, alpha, s) }) {
		t.Errorf("Doesn't panic if y is smaller than dst")
	}
}

func TestApply(t *testing.T) {
	s := []int{3, 4, 1, 7, 5}
	f := func(val int) int {
		return val * 2
	}
	truth := make([]int, len(s))
	for i, val := range s {
		truth[i] = val * 2
	}
	Apply(f, s)
	AreSlicesEqual(t, truth, s, "Wrong application of function")
}

func TestArgsort(t *testing.T) {
	s := []int{3, 4, 1, 7, 5}
	inds := make([]int, len(s))

	Argsort(s, inds)

	sortedS := []int{1, 3, 4, 5, 7}
	trueInds := []int{2, 0, 1, 4, 3}

	if !Equal(s, sortedS) {
		t.Error("elements not sorted correctly")
	}
	for i := range trueInds {
		if trueInds[i] != inds[i] {
			t.Error("inds not correct")
		}
	}

	inds = []int{1, 2}
	if !Panics(func() { Argsort(s, inds) }) {
		t.Error("does not panic if lengths do not match")
	}
}

func TestCount(t *testing.T) {
	s := []int{3, 4, 1, 7, 5}
	f := func(v int) bool { return v > 3 }
	truth := 3
	n := Count(f, s)
	if n != truth {
		t.Errorf("Wrong number of elements counted")
	}
}

func TestCumProd(t *testing.T) {
	s := []int{3, 4, 1, 7, 5}
	receiver := make([]int, len(s))
	CumProd(receiver, s)
	truth := []int{3, 12, 12, 84, 420}
	AreSlicesEqual(t, truth, receiver, "Wrong cumprod returned with new receiver")
	CumProd(receiver, s)
	AreSlicesEqual(t, truth, receiver, "Wrong cumprod returned with reused receiver")
	// Test that it panics
	if !Panics(func() { CumProd(make([]int, 2), make([]int, 3)) }) {
		t.Errorf("Did not panic with length mismatch")
	}
}

func TestCumSum(t *testing.T) {
	s := []int{3, 4, 1, 7, 5}
	receiver := make([]int, len(s))
	CumSum(receiver, s)
	truth := []int{3, 7, 8, 15, 20}
	AreSlicesEqual(t, truth, receiver, "Wrong cumsum returned with new receiver")
	CumSum(receiver, s)
	AreSlicesEqual(t, truth, receiver, "Wrong cumsum returned with reused receiver")

	// Test that it panics
	if !Panics(func() { CumSum(make([]int, 2), make([]int, 3)) }) {
		t.Errorf("Did not panic with length mismatch")
	}
}

func TestDiv(t *testing.T) {
	s1 := []int{5, 12, 27}
	s2 := []int{1, 2, 3}
	ans := []int{5, 6, 9}
	Div(s1, s2)
	if !Equal(s1, ans) {
		t.Errorf("Mul doesn't give correct answer")
	}
	s1short := []int{1}
	if !Panics(func() { Div(s1short, s2) }) {
		t.Errorf("Did not panic with unequal lengths")
	}
	s2short := []int{1}
	if !Panics(func() { Div(s1, s2short) }) {
		t.Errorf("Did not panic with unequal lengths")
	}
}

func TestDivTo(t *testing.T) {
	s1 := []int{5, 12, 27}
	s1orig := []int{5, 12, 27}
	s2 := []int{1, 2, 3}
	s2orig := []int{1, 2, 3}
	dst := make([]int, 3)
	ans := []int{5, 6, 9}
	DivTo(dst, s1, s2)
	if !Equal(dst, ans) {
		t.Errorf("DivTo doesn't give correct answer")
	}
	if !Equal(s1, s1orig) {
		t.Errorf("S1 changes during multo")
	}
	if !Equal(s2, s2orig) {
		t.Errorf("s2 changes during multo")
	}
	DivTo(dst, s1, s2)
	if !Equal(dst, ans) {
		t.Errorf("DivTo doesn't give correct answer reusing dst")
	}
	dstShort := []int{1}
	if !Panics(func() { DivTo(dstShort, s1, s2) }) {
		t.Errorf("Did not panic with s1 wrong length")
	}
	s1short := []int{1}
	if !Panics(func() { DivTo(dst, s1short, s2) }) {
		t.Errorf("Did not panic with s1 wrong length")
	}
	s2short := []int{1}
	if !Panics(func() { DivTo(dst, s1, s2short) }) {
		t.Errorf("Did not panic with s2 wrong length")
	}
}

func TestDot(t *testing.T) {
	s1 := []int{1, 2, 3, 4}
	s2 := []int{-3, 4, 5, -6}
	truth := -4
	ans := Dot(s1, s2)
	if ans != truth {
		t.Errorf("Dot product computed incorrectly")
	}

	// Test that it panics
	if !Panics(func() { Dot(make([]int, 2), make([]int, 3)) }) {
		t.Errorf("Did not panic with length mismatch")
	}
}

func TestEquals(t *testing.T) {
	s1 := []int{1, 2, 3, 4}
	s2 := []int{1, 2, 3, 4}
	if !Equal(s1, s2) {
		t.Errorf("Equal slices returned as unequal")
	}
	s2 = []int{1, 2, 3, 7}
	if Equal(s1, s2) {
		t.Errorf("Unequal slices returned as equal")
	}
}

func TestEqualLengths(t *testing.T) {
	s1 := []int{1, 2, 3, 4}
	s2 := []int{1, 2, 3, 4}
	s3 := []int{1, 2, 3}
	if !EqualLengths(s1, s2) {
		t.Errorf("Equal lengths returned as unequal")
	}
	if EqualLengths(s1, s3) {
		t.Errorf("Unequal lengths returned as equal")
	}
	if !EqualLengths(s1) {
		t.Errorf("Single slice returned as unequal")
	}
	if !EqualLengths() {
		t.Errorf("No slices returned as unequal")
	}
}

func eqIntSlice(one, two []int) string {
	if len(one) != len(two) {
		return "Length mismatch"
	}
	for i, val := range one {
		if val != two[i] {
			return "Index " + strconv.Itoa(i) + " mismatch"
		}
	}
	return ""
}

func TestFind(t *testing.T) {
	s := []int{3, 4, 1, 7, 5}
	f := func(v int) bool { return v > 3 }
	allTrueInds := []int{1, 3, 4}

	// Test finding first two elements
	inds, err := Find(nil, f, s, 2)
	if err != nil {
		t.Errorf("Find first two: Improper error return")
	}
	trueInds := allTrueInds[:2]
	str := eqIntSlice(inds, trueInds)
	if str != "" {
		t.Errorf("Find first two: " + str)
	}

	// Test finding first two elements with non nil slice
	inds = []int{1, 2, 3, 4, 5, 6}
	inds, err = Find(inds, f, s, 2)
	if err != nil {
		t.Errorf("Find first two non-nil: Improper error return")
	}
	str = eqIntSlice(inds, trueInds)
	if str != "" {
		t.Errorf("Find first two non-nil: " + str)
	}

	// Test finding too many elements
	inds, err = Find(inds, f, s, 4)
	if err == nil {
		t.Errorf("Request too many: No error returned")
	}
	str = eqIntSlice(inds, allTrueInds)
	if str != "" {
		t.Errorf("Request too many: Does not match all of the inds: " + str)
	}

	// Test finding all elements
	inds, err = Find(nil, f, s, -1)
	if err != nil {
		t.Errorf("Find all: Improper error returned")
	}
	str = eqIntSlice(inds, allTrueInds)
	if str != "" {
		t.Errorf("Find all: Does not match all of the inds: " + str)
	}
}

func TestMax(t *testing.T) {
	s := []int{3, 4, 1, 7, 5}
	val, ind := Max(s)
	if val != 7 {
		t.Errorf("Wrong value returned")
	}
	if ind != 3 {
		t.Errorf("Wrong index returned")
	}
}

func TestMin(t *testing.T) {
	s := []int{3, 4, 1, 7, 5}
	val, ind := Min(s)
	if val != 1 {
		t.Errorf("Wrong value returned")
	}
	if ind != 2 {
		t.Errorf("Wrong index returned")
	}
}

func TestMul(t *testing.T) {
	s1 := []int{1, 2, 3}
	s2 := []int{1, 2, 3}
	ans := []int{1, 4, 9}
	Mul(s1, s2)
	if !Equal(s1, ans) {
		t.Errorf("Mul doesn't give correct answer")
	}
	s1short := []int{1}
	if !Panics(func() { Mul(s1short, s2) }) {
		t.Errorf("Did not panic with unequal lengths")
	}
	s2short := []int{1}
	if !Panics(func() { Mul(s1, s2short) }) {
		t.Errorf("Did not panic with unequal lengths")
	}
}

func TestMulTo(t *testing.T) {
	s1 := []int{1, 2, 3}
	s1orig := []int{1, 2, 3}
	s2 := []int{1, 2, 3}
	s2orig := []int{1, 2, 3}
	dst := make([]int, 3)
	ans := []int{1, 4, 9}
	MulTo(dst, s1, s2)
	if !Equal(dst, ans) {
		t.Errorf("MulTo doesn't give correct answer")
	}
	if !Equal(s1, s1orig) {
		t.Errorf("S1 changes during multo")
	}
	if !Equal(s2, s2orig) {
		t.Errorf("s2 changes during multo")
	}
	MulTo(dst, s1, s2)
	if !Equal(dst, ans) {
		t.Errorf("MulTo doesn't give correct answer reusing dst")
	}
	dstShort := []int{1}
	if !Panics(func() { MulTo(dstShort, s1, s2) }) {
		t.Errorf("Did not panic with s1 wrong length")
	}
	s1short := []int{1}
	if !Panics(func() { MulTo(dst, s1short, s2) }) {
		t.Errorf("Did not panic with s1 wrong length")
	}
	s2short := []int{1}
	if !Panics(func() { MulTo(dst, s1, s2short) }) {
		t.Errorf("Did not panic with s2 wrong length")
	}
}

// func TestNearest(t *testing.T) {
// 	s := []int{6.2, 3, 5, 6.2, 8}
// 	ind := Nearest(s, 2.0)
// 	if ind != 1 {
// 		t.Errorf("Wrong index returned when value is less than all of elements")
// 	}
// 	ind = Nearest(s, 9.0)
// 	if ind != 4 {
// 		t.Errorf("Wrong index returned when value is greater than all of elements")
// 	}
// 	ind = Nearest(s, 3.1)
// 	if ind != 1 {
// 		t.Errorf("Wrong index returned when value is greater than closest element")
// 	}
// 	ind = Nearest(s, 3.1)
// 	if ind != 1 {
// 		t.Errorf("Wrong index returned when value is greater than closest element")
// 	}
// 	ind = Nearest(s, 2.9)
// 	if ind != 1 {
// 		t.Errorf("Wrong index returned when value is less than closest element")
// 	}
// 	ind = Nearest(s, 3)
// 	if ind != 1 {
// 		t.Errorf("Wrong index returned when value is equal to element")
// 	}
// 	ind = Nearest(s, 6.2)
// 	if ind != 0 {
// 		t.Errorf("Wrong index returned when value is equal to several elements")
// 	}
// 	ind = Nearest(s, 4)
// 	if ind != 1 {
// 		t.Errorf("Wrong index returned when value is exactly between two closest elements")
// 	}
// }

// func TestNearestWithinSpan(t *testing.T) {

// 	if !Panics(func() { NearestWithinSpan(13, 7, 8.2, 10) }) {
// 		t.Errorf("Did not panic below lower bound")
// 	}
// 	if !Panics(func() { NearestWithinSpan(13, 7, 8.2, 10) }) {
// 		t.Errorf("Did not panic above upper bound")
// 	}
// 	ind := NearestWithinSpan(13, 7, 8.2, 7.19)
// 	if ind != 2 {
// 		t.Errorf("Wrong value when just below the bucket. %i found, %i expected", ind, 2)
// 	}
// 	ind = NearestWithinSpan(13, 7, 8.2, 7.21)
// 	if ind != 2 {
// 		t.Errorf("Wrong value when just above the bucket. %i found, %i expected", ind, 2)
// 	}
// 	ind = NearestWithinSpan(13, 7, 8.2, 7.2)
// 	if ind != 2 {
// 		t.Errorf("Wrong value when equal to bucket. %i found, %i expected", ind, 2)
// 	}
// 	ind = NearestWithinSpan(13, 7, 8.2, 7.151)
// 	if ind != 2 {
// 		t.Errorf("Wrong value when just above halfway point. %i found, %i expected", ind, 2)
// 	}
// 	ind = NearestWithinSpan(13, 7, 8.2, 7.249)
// 	if ind != 2 {
// 		t.Errorf("Wrong value when just below halfway point. %i found, %i expected", ind, 2)
// 	}
// }

func TestProd(t *testing.T) {
	s := []int{}
	val := Prod(s)
	if val != 1 {
		t.Errorf("Val not returned as default when slice length is zero")
	}
	s = []int{3, 4, 1, 7, 5}
	val = Prod(s)
	if val != 420 {
		t.Errorf("Wrong prod returned. Expected %v returned %v", 420, val)
	}
}

func TestScale(t *testing.T) {
	s := []int{3, 4, 1, 7, 5}
	c := 5
	truth := []int{15, 20, 5, 35, 25}
	Scale(c, s)
	AreSlicesEqual(t, truth, s, "Bad scaling")
}

// func TestSpan(t *testing.T) {
// 	receiver := make([]int, 5)
// 	truth := []int{1, 2, 3, 4, 5}
// 	Span(receiver, 1, 5)
// 	AreSlicesEqual(t, truth, receiver, "Improper linspace")
// 	receiver = make([]int, 6)
// 	truth = []int{0, 0.2, 0.4, 0.6, 0.8, 1.0}
// 	Span(receiver, 0, 1)
// 	AreSlicesEqual(t, truth, receiver, "Improper linspace")
// 	if !Panics(func() { Span(nil, 1, 5) }) {
// 		t.Errorf("Span accepts nil argument")
// 	}
// 	if !Panics(func() { Span(make([]int, 1), 1, 5) }) {
// 		t.Errorf("Span accepts argument of len = 1")
// 	}
// }

func TestSub(t *testing.T) {
	s := []int{3, 4, 1, 7, 5}
	v := []int{1, 2, 3, 4, 5}
	truth := []int{2, 2, -2, 3, 0}
	Sub(s, v)
	AreSlicesEqual(t, truth, s, "Bad subtract")
	// Test that it panics
	if !Panics(func() { Sub(make([]int, 2), make([]int, 3)) }) {
		t.Errorf("Did not panic with length mismatch")
	}
}

func TestSubTo(t *testing.T) {
	s := []int{3, 4, 1, 7, 5}
	v := []int{1, 2, 3, 4, 5}
	truth := []int{2, 2, -2, 3, 0}
	dst := make([]int, len(s))
	SubTo(dst, s, v)
	AreSlicesEqual(t, truth, dst, "Bad subtract")
	// Test that all mismatch combinations panic
	if !Panics(func() { SubTo(make([]int, 2), make([]int, 3), make([]int, 3)) }) {
		t.Errorf("Did not panic with dst different length")
	}
	if !Panics(func() { SubTo(make([]int, 3), make([]int, 2), make([]int, 3)) }) {
		t.Errorf("Did not panic with subtractor different length")
	}
	if !Panics(func() { SubTo(make([]int, 3), make([]int, 3), make([]int, 2)) }) {
		t.Errorf("Did not panic with subtractee different length")
	}
}

func TestSum(t *testing.T) {
	s := []int{}
	val := Sum(s)
	if val != 0 {
		t.Errorf("Val not returned as default when slice length is zero")
	}
	s = []int{3, 4, 1, 7, 5}
	val = Sum(s)
	if val != 20 {
		t.Errorf("Wrong sum returned")
	}
}

func RandomSlice(l int) []int {
	s := make([]int, l)
	for i := range s {
		s[i] = rand.Int()
	}
	return s
}

func benchmarkMin(b *testing.B, s []int) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Min(s)
	}
}

func BenchmarkMinSmall(b *testing.B) {
	s := RandomSlice(SMALL)
	benchmarkMin(b, s)
}

func BenchmarkMinMed(b *testing.B) {
	s := RandomSlice(MEDIUM)
	benchmarkMin(b, s)
}

func BenchmarkMinLarge(b *testing.B) {
	s := RandomSlice(LARGE)
	benchmarkMin(b, s)
}
func BenchmarkMinHuge(b *testing.B) {
	s := RandomSlice(HUGE)
	benchmarkMin(b, s)
}

func benchmarkAdd(b *testing.B, s ...[]int) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Add(s[0], s[1:]...)
	}
}

func BenchmarkAddTwoSmall(b *testing.B) {
	i := SMALL
	s := RandomSlice(i)
	t := RandomSlice(i)
	benchmarkAdd(b, s, t)
}

func BenchmarkAddFourSmall(b *testing.B) {
	i := SMALL
	s := RandomSlice(i)
	t := RandomSlice(i)
	u := RandomSlice(i)
	v := RandomSlice(i)
	benchmarkAdd(b, s, t, u, v)
}

func BenchmarkAddTwoMed(b *testing.B) {
	i := MEDIUM
	s := RandomSlice(i)
	t := RandomSlice(i)
	benchmarkAdd(b, s, t)
}

func BenchmarkAddFourMed(b *testing.B) {
	i := MEDIUM
	s := RandomSlice(i)
	t := RandomSlice(i)
	u := RandomSlice(i)
	v := RandomSlice(i)
	benchmarkAdd(b, s, t, u, v)
}

func BenchmarkAddTwoLarge(b *testing.B) {
	i := LARGE
	s := RandomSlice(i)
	t := RandomSlice(i)
	benchmarkAdd(b, s, t)
}

func BenchmarkAddFourLarge(b *testing.B) {
	i := LARGE
	s := RandomSlice(i)
	t := RandomSlice(i)
	u := RandomSlice(i)
	v := RandomSlice(i)
	benchmarkAdd(b, s, t, u, v)
}

func BenchmarkAddTwoHuge(b *testing.B) {
	i := HUGE
	s := RandomSlice(i)
	t := RandomSlice(i)
	benchmarkAdd(b, s, t)
}

func BenchmarkAddFourHuge(b *testing.B) {
	i := HUGE
	s := RandomSlice(i)
	t := RandomSlice(i)
	u := RandomSlice(i)
	v := RandomSlice(i)
	benchmarkAdd(b, s, t, u, v)
}

// func benchmarkLogSumExp(b *testing.B, s []int) {
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		LogSumExp(s)
// 	}
// }

// func BenchmarkLogSumExpSmall(b *testing.B) {
// 	s := RandomSlice(SMALL)
// 	benchmarkLogSumExp(b, s)
// }

// func BenchmarkLogSumExpMed(b *testing.B) {
// 	s := RandomSlice(MEDIUM)
// 	benchmarkLogSumExp(b, s)
// }

// func BenchmarkLogSumExpLarge(b *testing.B) {
// 	s := RandomSlice(LARGE)
// 	benchmarkLogSumExp(b, s)
// }
// func BenchmarkLogSumExpHuge(b *testing.B) {
// 	s := RandomSlice(HUGE)
// 	benchmarkLogSumExp(b, s)
// }

func benchmarkDot(b *testing.B, s1 []int, s2 []int) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Dot(s1, s2)
	}
}

func BenchmarkDotSmall(b *testing.B) {
	s1 := RandomSlice(SMALL)
	s2 := RandomSlice(SMALL)
	benchmarkDot(b, s1, s2)
}

func BenchmarkDotMed(b *testing.B) {
	s1 := RandomSlice(MEDIUM)
	s2 := RandomSlice(MEDIUM)
	benchmarkDot(b, s1, s2)
}

func BenchmarkDotLarge(b *testing.B) {
	s1 := RandomSlice(LARGE)
	s2 := RandomSlice(LARGE)
	benchmarkDot(b, s1, s2)
}
func BenchmarkDotHuge(b *testing.B) {
	s1 := RandomSlice(HUGE)
	s2 := RandomSlice(HUGE)
	benchmarkDot(b, s1, s2)
}
