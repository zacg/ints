// package ints provides a set of helper routines for dealing with slices
// of int. The functions avoid allocations to allow for use within tight
// loops without garbage collection overhead.

package ints

import (
	"errors"
	//"math"
	"sort"
)

// Add returns the element-wise sum of all the slices with the
// results stored in the first slice.
// For computational efficiency, it is assumed that all of
// the variadic arguments have the same length. If this is
// in doubt, EqLen can be used.
func Add(dst []int, slices ...[]int) []int {
	if len(slices) == 0 {
		return nil
	}
	if len(dst) != len(slices[0]) {
		panic("ints: length of destination does not match length of the slices")
	}
	for _, slice := range slices {
		for j, val := range slice {
			dst[j] += val
		}
	}
	return dst
}

// AddConst adds the value c to all of the values in s.
func AddConst(c int, s []int) {
	for i := range s {
		s[i] += c
	}
}

// AddScaled performs dst = dst + alpha * s.
// It panics if the lengths of dst and s are not equal.
func AddScaled(dst []int, alpha int, s []int) {
	if len(dst) != len(s) {
		panic("ints: length of destination and source to not match")
	}
	for i, val := range s {
		dst[i] += alpha * val
	}
}

// AddScaledTo performs dst = y + alpha * s.
// It panics if the lengths of dst, y, and s are not equal.
func AddScaledTo(dst []int, y []int, alpha int, s []int) []int {
	if len(dst) != len(s) || len(dst) != len(y) {
		panic("ints: lengths of slices do not match")
	}
	for i, val := range s {
		dst[i] = y[i] + alpha*val
	}
	return dst
}

type argsort struct {
	s    []int
	inds []int
}

func (a argsort) Len() int {
	return len(a.s)
}

func (a argsort) Less(i, j int) bool {
	return a.s[i] < a.s[j]
}

func (a argsort) Swap(i, j int) {
	a.s[i], a.s[j] = a.s[j], a.s[i]
	a.inds[i], a.inds[j] = a.inds[j], a.inds[i]
}

// ApplyFunc applies a function f (math.Exp, math.Sin, etc.) to every element
// of the slice s.
func Apply(f func(int) int, s []int) {
	for i, val := range s {
		s[i] = f(val)
	}
}

// Argsort sorts the elements of s while tracking their original order.
// At the conclusion of Argsort, s will contain the original elements of s
// but sorted in increasing order, and inds will contain the original position
// of the elements in the slice such that s[i] = sOrig[inds[i]].
func Argsort(s []int, inds []int) {
	if len(s) != len(inds) {
		panic("ints: length of inds does not match length of slice")
	}
	for i := range s {
		inds[i] = i
	}

	a := argsort{s: s, inds: inds}
	sort.Sort(a)
}

// Count applies the function f to every element of s and returns the number
// of times the function returned true.
func Count(f func(int) bool, s []int) int {
	var n int
	for _, val := range s {
		if f(val) {
			n++
		}
	}
	return n
}

// Cumprod finds the cumulative product of the first i elements in
// s and puts them in place into the ith element of the
// destination. A panic will occur if lengths of do not match.
func CumProd(dst, s []int) []int {
	if len(dst) != len(s) {
		panic("ints: length of destination does not match length of the source")
	}
	dst[0] = s[0]
	for i := 1; i < len(s); i++ {
		dst[i] = dst[i-1] * s[i]
	}
	return dst
}

// Cumsum finds the cumulative sum of the first i elements in
// s and puts them in place into the ith element of the
// destination. A panic will occur if lengths of arguments do not match.
func CumSum(dst, s []int) []int {
	if len(dst) != len(s) {
		panic("ints: length of destination does not match length of the source")
	}
	dst[0] = s[0]
	for i := 1; i < len(s); i++ {
		dst[i] = dst[i-1] + s[i]
	}
	return dst
}

// Div performs element-wise division between s
// and t and stores the value in s. It panics if the
// lengths of s and t are not equal.
func Div(s []int, t []int) {
	if len(s) != len(t) {
		panic("ints: slice lengths do not match")
	}
	for i, val := range t {
		s[i] /= val
	}
}

// DivTo performs element-wise division between s
// and t and stores the value in dst. It panics if the
// lengths of s, t, and dst are not equal.
func DivTo(dst []int, s []int, t []int) []int {
	if len(s) != len(t) || len(dst) != len(t) {
		panic("ints: slice lengths do not match")
	}
	for i, val := range t {
		dst[i] = s[i] / val
	}
	return dst
}

// Dot computes the dot product of s1 and s2, i.e.
// sum_{i = 1}^N s1[i]*s2[i].
// A panic will occur if lengths of arguments do not match.
func Dot(s1, s2 []int) int {
	if len(s1) != len(s2) {
		panic("ints: lengths of the slices do not match")
	}
	var sum int
	for i, val := range s1 {
		sum += val * s2[i]
	}
	return sum
}

// Equal returns true if the slices have equal lengths and
// all elements are numerically identical.
func Equal(s1, s2 []int) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i, val := range s1 {
		if s2[i] != val {
			return false
		}
	}
	return true
}

// EqualsFunc returns true if the slices have the same lengths
// and the function returns true for all element pairs.
func EqualFunc(s1, s2 []int, f func(int, int) bool) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i, val := range s1 {
		if !f(val, s2[i]) {
			return false
		}
	}
	return true

}

func ulpDiff(a, b uint64) uint64 {
	if a > b {
		a, b = b, a
	}
	return b - a
}

// Eqlen returns true if all of the slices have equal length,
// and false otherwise. Returns true if there are no input slices.
func EqualLengths(slices ...[]int) bool {
	// This length check is needed: http://play.golang.org/p/sdty6YiLhM
	if len(slices) == 0 {
		return true
	}
	l := len(slices[0])
	for i := 1; i < len(slices); i++ {
		if len(slices[i]) != l {
			return false
		}
	}
	return true
}

// Fill loops over the elements of s and stores a value generated from f.
// f is called n times, where n = len(s)
func Fill(f func() int, s []int) {
	for i := range s {
		s[i] = f()
	}
}

// Find applies f to every element of s and returns the indices of the first
// k elements for which the f returns true, or all such elements
// if k < 0.
// Find will reslice inds to have 0 length, and will append
// found indices to inds.
// If k > 0 and there are fewer than k elements in s satisfying f,
// all of the found elements will be returned along with an error.
func Find(inds []int, f func(int) bool, s []int, k int) ([]int, error) {

	// inds is also returned to allow for calling with nil

	// Reslice inds to have zero length
	inds = inds[:0]

	// If zero elements requested, can just return
	if k == 0 {
		return inds, nil
	}

	// If k < 0, return all of the found indices
	if k < 0 {
		for i, val := range s {
			if f(val) {
				inds = append(inds, i)
			}
		}
		return inds, nil
	}

	// Otherwise, find the first k elements
	nFound := 0
	for i, val := range s {
		if f(val) {
			inds = append(inds, i)
			nFound++
			if nFound == k {
				return inds, nil
			}
		}
	}
	// Finished iterating over the loop, which means k elements were not found
	return inds, errors.New("ints: insufficient elements found")
}

// Max returns the maximum value in the slice and the location of
// the maximum value. If the input slice is empty, Max will panic.
func Max(s []int) (max int, ind int) {
	max = s[0]
	ind = 0
	for i, val := range s {
		if val > max {
			max = val
			ind = i
		}
	}
	return max, ind
}

// Min returns the minimum value in the slice and the index of
// the minimum value. If the input slice is empty, Min will panic.
func Min(s []int) (min int, ind int) {
	min = s[0]
	ind = 0
	for i, val := range s {
		if val < min {
			min = val
			ind = i
		}
	}
	return min, ind
}

// Mul performs element-wise multiplication between s
// and t and stores the value in s. Panics if the
// lengths of s and t are not equal.
func Mul(s []int, t []int) {
	if len(s) != len(t) {
		panic("ints: slice lengths do not match")
	}
	for i, val := range t {
		s[i] *= val
	}
}

// MulTo performs element-wise multiplication between s
// and t and stores the value in dst. Panics if the
// lengths of s, t, and dst are not equal.
func MulTo(dst []int, s []int, t []int) []int {
	if len(s) != len(t) || len(dst) != len(t) {
		panic("ints: slice lengths do not match")
	}
	for i, val := range t {
		dst[i] = val * s[i]
	}
	return dst
}

// Nearest returns the index of the element in s
// whose value is nearest to v.  If several such
// elements exist, the lowest index is returned.
// func Nearest(s []int, v int) (ind int) {
// 	dist := math.Abs(v - s[0])
// 	ind = 0
// 	for i, val := range s {
// 		newDist := math.Abs(v - val)
// 		if newDist < dist {
// 			dist = newDist
// 			ind = i
// 		}
// 	}
// 	return
// }

// NearestInSpan return the index of a hypothetical vector created
// by Span with length n and bounds l and u whose value is closest
// to v. Assumes u > l. If the value is greater than u or less than
// l, the function will panic.
// func NearestWithinSpan(n int, l, u int, v int) int {
// 	if v < l || v > u {
// 		panic("ints: value outside span bounds")
// 	}

// 	// Can't guarantee anything about exactly halfway between
// 	// because of floating point weirdness
// 	return int((int(n)-1)/(u-l)*(v-l) + 0.5)
// }

// Prod returns the product of the elements of the slice
// Returns 1 if len(s) = 0.
func Prod(s []int) (prod int) {
	prod = 1
	for _, val := range s {
		prod *= val
	}
	return prod
}

// Scale multiplies every element in s by c.
func Scale(c int, s []int) {
	for i := range s {
		s[i] *= c
	}
}

// Span returns a set of N equally spaced points between l and u, where N
// is equal to the length of the destination. The first element of the destination
// is l, the final element of the destination is u.
// Panics if len(dst) < 2.
func Span(dst []int, l, u int) []int {
	n := len(dst)
	if n < 2 {
		panic("ints: destination must have length >1")
	}
	step := (u - l) / int(n-1)
	for i := range dst {
		dst[i] = l + step*int(i)
	}
	return dst
}

// Sub subtracts, element-wise, the first argument from the second. Assumes
// the lengths of s and t match (can be tested with EqLen).
func Sub(s, t []int) {
	if len(s) != len(t) {
		panic("ints: length of the slices do not match")
	}
	for i, val := range t {
		s[i] -= val
	}
}

// SubTo subtracts, element-wise, the first argument from the second and
// stores the result in dest. Panics if the lengths of s and t do not match.
func SubTo(dst, s, t []int) []int {
	if len(s) != len(t) {
		panic("ints: length of subtractor and subtractee do not match")
	}
	if len(dst) != len(s) {
		panic("ints: length of destination does not match length of subtractor")
	}
	for i, val := range t {
		dst[i] = s[i] - val
	}
	return dst
}

// Sum returns the sum of the elements of the slice.
func Sum(s []int) (sum int) {
	for _, val := range s {
		sum += val
	}
	return
}
