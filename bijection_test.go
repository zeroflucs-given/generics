package generics

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type bijectionTestDef struct {
	Name     string
	Left     []string
	Right    []string
	MappingL map[string]string
	MappingR map[string]string
	SuccessL bool
	SuccessR bool
}

var bijectionEqualityTests = []bijectionTestDef{
	{
		Name:     "simple",
		Left:     []string{"a", "b"},
		Right:    []string{"a", "b"},
		MappingL: map[string]string{"a": "a", "b": "b"},
		MappingR: map[string]string{"a": "a", "b": "b"},
		SuccessL: true,
		SuccessR: true,
	},
	{
		Name:     "permuted_right_order",
		Left:     []string{"a", "b"},
		Right:    []string{"b", "a"},
		MappingL: map[string]string{"a": "a", "b": "b"},
		MappingR: map[string]string{"a": "a", "b": "b"},
		SuccessL: true,
		SuccessR: true,
	},
	{
		Name:     "unequal_lengths_fail",
		Left:     []string{"a"},
		Right:    []string{"a", "b"},
		MappingL: nil,
		MappingR: nil,
		SuccessL: false,
		SuccessR: false,
	},
	{
		Name:     "no_match_for_one_left_fail",
		Left:     []string{"x", "b"},
		Right:    []string{"a", "b"},
		MappingL: nil,
		MappingR: nil,
		SuccessL: false,
		SuccessR: false,
	},
	{
		Name:     "injectivity_violation_same_right_taken",
		Left:     []string{"a", "a"},
		Right:    []string{"a", "b"},
		MappingL: nil,
		MappingR: nil,
		SuccessL: false,
		SuccessR: false,
	},
	{
		Name:     "empty_inputs_ok_empty_mapping",
		Left:     []string{},
		Right:    []string{},
		MappingL: map[string]string{},
		MappingR: map[string]string{},
		SuccessL: true,
		SuccessR: true,
	},
	{
		Name:     "single_empty_on_both_sides_ok",
		Left:     []string{""},
		Right:    []string{""},
		MappingL: map[string]string{"": ""},
		MappingR: map[string]string{"": ""},
		SuccessL: true,
		SuccessR: true,
	},
	{
		Name:     "mixed_with_empty_ok",
		Left:     []string{"", "a"},
		Right:    []string{"a", ""},
		MappingL: map[string]string{"": "", "a": "a"},
		MappingR: map[string]string{"": "", "a": "a"},
		SuccessL: true,
		SuccessR: true,
	},
	{
		Name:     "empty_mismatch_fail",
		Left:     []string{""},
		Right:    []string{"a"},
		MappingL: nil,
		MappingR: nil,
		SuccessL: false,
		SuccessR: false,
	},
	{
		Name:     "duplicate_empty_lefts_injectivity_fail",
		Left:     []string{"", ""},
		Right:    []string{"", "x"},
		MappingL: nil, // both lefts equal "" can only map to the single "" on right -> injectivity violation
		MappingR: nil,
		SuccessL: false,
		SuccessR: false,
	},
}

var bijectionContainsTests = []bijectionTestDef{
	{
		Name:     "simple",
		Left:     []string{"Elon", "Duke"},
		Right:    []string{"Elon Phoenix", "Duke Blue Devils"},
		MappingL: map[string]string{"Elon": "Elon Phoenix", "Duke": "Duke Blue Devils"},
		SuccessL: true,
		SuccessR: false,
	},
	{
		Name:     "one_left_matches_multiple_rights_fail",
		Left:     []string{"art", "chem"},
		Right:    []string{"intro-to-art", "art-and-design"},
		MappingL: nil, // "art" matches two rights -> fail (1L -> 2R)
		MappingR: nil, // nothing matches
		SuccessL: false,
		SuccessR: false,
	},
	{
		Name:     "two_lefts_target_same_right_fail",
		Left:     []string{"Elon", "lon"},
		Right:    []string{"Elon", "Musk"},
		MappingL: nil, // "Elon" and "lon" on the Left both match "Elon" on the Right (2L -> 1R)
		MappingR: nil, // "Musk" is unmatched
		SuccessL: false,
		SuccessR: false,
	},
	{
		Name:     "both_sides_empty_string",
		Left:     []string{""},
		Right:    []string{""},
		MappingL: map[string]string{"": ""},
		MappingR: map[string]string{"": ""},
		SuccessL: true,
		SuccessR: true,
	},
	{
		Name:     "empty_left_matches_everything_fail",
		Left:     []string{"", "x"},
		Right:    []string{"foo", "barx"},
		MappingL: nil, // "" matches both "foo" and "barx" -> violates uniqueness
		MappingR: nil, // No matches
		SuccessL: false,
		SuccessR: false,
	},
	{
		Name:     "empty_right_string_only",
		Left:     []string{"a"},
		Right:    []string{""},
		MappingL: nil,                        // "a" is not contained in ""
		MappingR: map[string]string{"": "a"}, // everything contains the empty string
		SuccessL: false,
		SuccessR: true,
	},
}

func TestCheckBijective_Equality(t *testing.T) {
	// Simple bijection, should work both ways.
	predEqual := func(a, b string) bool { return a == b }

	for _, test := range bijectionEqualityTests {
		t.Run(test.Name, func(t *testing.T) {
			t.Run("left_right", func(t *testing.T) {
				mapping, success := CheckBijection(test.Left, test.Right, predEqual)
				assert.Equal(t, test.SuccessL, success)
				assert.Equal(t, test.MappingL, mapping)
			})

			t.Run("right_left", func(t *testing.T) {
				mapping, success := CheckBijection(test.Right, test.Left, predEqual)
				assert.Equal(t, test.SuccessR, success)
				assert.Equal(t, test.MappingR, mapping)
			})
		})
	}
}

func TestCheckBijective_Substring(t *testing.T) {
	predContains := func(a, b string) bool { return strings.Contains(b, a) }
	for _, test := range bijectionContainsTests {
		t.Run(test.Name, func(t *testing.T) {

			t.Run("left_right", func(t *testing.T) {
				mapping, success := CheckBijection(test.Left, test.Right, predContains)
				assert.Equal(t, test.SuccessL, success)
				assert.Equal(t, test.MappingL, mapping)
			})

			t.Run("right_left", func(t *testing.T) {
				mapping, success := CheckBijection(test.Right, test.Left, predContains)
				assert.Equal(t, test.SuccessR, success)
				assert.Equal(t, test.MappingR, mapping)
			})
		})
	}
}

func TestCheckPartialBijection_Equality(t *testing.T) {
	// Simple bijection, should work both ways.
	predEqual := func(a, b string) bool { return a == b }

	tests := []bijectionTestDef{
		{
			Name:     "simple",
			Left:     []string{"a", "b"},
			Right:    []string{"a", "b"},
			MappingL: map[string]string{"a": "a", "b": "b"},
			MappingR: map[string]string{"a": "a", "b": "b"},
			SuccessL: true,
			SuccessR: true,
		},
		{
			Name:     "permuted_right_order",
			Left:     []string{"a", "b"},
			Right:    []string{"b", "a"},
			MappingL: map[string]string{"a": "a", "b": "b"},
			MappingR: map[string]string{"a": "a", "b": "b"},
			SuccessL: true,
			SuccessR: true,
		},
		{
			Name:     "no_match_for_one_left_fail",
			Left:     []string{"x", "b"},
			Right:    []string{"a", "b"},
			MappingL: nil,
			MappingR: nil,
			SuccessL: false,
			SuccessR: false,
		},
		{
			Name:     "injectivity_violation_same_right_taken",
			Left:     []string{"a", "a"},
			Right:    []string{"a", "b"},
			MappingL: nil,
			MappingR: nil,
			SuccessL: false,
			SuccessR: false,
		},
		{
			Name:     "empty_inputs_ok_empty_mapping",
			Left:     []string{},
			Right:    []string{},
			MappingL: map[string]string{},
			MappingR: map[string]string{},
			SuccessL: true,
			SuccessR: true,
		},
		{
			Name:     "single_empty_on_both_sides_ok",
			Left:     []string{""},
			Right:    []string{""},
			MappingL: map[string]string{"": ""},
			MappingR: map[string]string{"": ""},
			SuccessL: true,
			SuccessR: true,
		},
		{
			Name:     "mixed_with_empty_ok",
			Left:     []string{"", "a"},
			Right:    []string{"a", ""},
			MappingL: map[string]string{"": "", "a": "a"},
			MappingR: map[string]string{"": "", "a": "a"},
			SuccessL: true,
			SuccessR: true,
		},
		{
			Name:     "empty_mismatch_fail",
			Left:     []string{""},
			Right:    []string{"a"},
			MappingL: nil,
			MappingR: nil,
			SuccessL: false,
			SuccessR: false,
		},
		{
			Name:     "duplicate_empty_lefts_injectivity_fail",
			Left:     []string{"", ""},
			Right:    []string{"", "x"},
			MappingL: nil, // both lefts equal "" can only map to the single "" on right -> injectivity violation
			MappingR: nil,
			SuccessL: false,
			SuccessR: false,
		},
		{
			Name:     "simple_subset",
			Left:     []string{"a", "b"},
			Right:    []string{"a"},
			MappingL: map[string]string{"a": "a"}, // {a} is a subject of {a,b}, and no two {a,b} maps to {a}
			MappingR: nil,
			SuccessL: true,
			SuccessR: false,
		},
		{
			Name:     "completely_different_subset",
			Left:     []string{"a", "b"},
			Right:    []string{"c"},
			MappingL: nil,
			MappingR: nil,
			SuccessL: false,
			SuccessR: false,
		},
		{
			Name:     "larger_subset",
			Left:     []string{"a", "b"},
			Right:    []string{"a", "b", "c"},
			MappingL: nil,
			MappingR: map[string]string{"a": "a", "b": "b"},
			SuccessL: false,
			SuccessR: true,
		},
		{
			Name:     "partial_subset",
			Left:     []string{"a", "b"},
			Right:    []string{"b", "c"},
			MappingL: nil,
			MappingR: nil,
			SuccessL: false,
			SuccessR: false,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			t.Run("left_right", func(t *testing.T) {
				mapping, success := CheckPartialBijection(test.Left, test.Right, predEqual)
				assert.Equal(t, test.SuccessL, success)
				assert.Equal(t, test.MappingL, mapping)
			})

			t.Run("right_left", func(t *testing.T) {
				mapping, success := CheckPartialBijection(test.Right, test.Left, predEqual)
				assert.Equal(t, test.SuccessR, success)
				assert.Equal(t, test.MappingR, mapping)
			})
		})
	}
}
