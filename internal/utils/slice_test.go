package utils

import (
	"strings"
	"testing"
)

func TestNumSliceMoveElemTo(t *testing.T) {
	nums := []int{1, 3, 2}
	pos := 1
	target := 2
	epxNums := []int{1, 2, 3}
	nums = MoveSliceElemTo(nums, pos, func(e int) bool {
		return e == target
	})
	t.Log("new nums:", nums)
	if !EqualsSlice(epxNums, nums) {
		t.Error("two slice is different")
	}
}

func TestStrSliceMoveElemTo(t *testing.T) {
	strs := []string{"hello-world.hash", "hello-world.sh", "hello-world.zig"}
	pos := 0
	suffix := ".zig"
	expStrs := []string{"hello-world.zig", "hello-world.hash", "hello-world.sh"}
	strs = MoveSliceElemTo(strs, pos, func(s string) bool {
		return strings.HasSuffix(s, suffix)
	})
	t.Log("new string slice:", strs)
	if !EqualsSlice(expStrs, strs) {
		t.Error("two slice is different")
	}
}
