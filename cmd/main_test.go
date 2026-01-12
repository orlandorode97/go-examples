package main

import "testing"

var classificationTestCases = []struct {
	description string
	input       int64
	ok          bool
	expected    Classification
}{
	{
		description: "Smallest perfect number is classified correctly",
		input:       6,
		ok:          true,
		expected:    ClassificationPerfect,
	},
	{
		description: "Medium perfect number is classified correctly",
		input:       28,
		ok:          true,
		expected:    ClassificationPerfect,
	},
	{
		description: "Large perfect number is classified correctly",
		input:       33550336,
		ok:          true,
		expected:    ClassificationPerfect,
	},
	{
		description: "Smallest abundant number is classified correctly",
		input:       12,
		ok:          true,
		expected:    ClassificationAbundant,
	},
	{
		description: "Medium abundant number is classified correctly",
		input:       30,
		ok:          true,
		expected:    ClassificationAbundant,
	},
	{
		description: "Large abundant number is classified correctly",
		input:       33550335,
		ok:          true,
		expected:    ClassificationAbundant,
	},
	{
		description: "Smallest prime deficient number is classified correctly",
		input:       2,
		ok:          true,
		expected:    ClassificationDeficient,
	},
	{
		description: "Smallest non-prime deficient number is classified correctly",
		input:       4,
		ok:          true,
		expected:    ClassificationDeficient,
	},
	{
		description: "Medium deficient number is classified correctly",
		input:       32,
		ok:          true,
		expected:    ClassificationDeficient,
	},
	{
		description: "Large deficient number is classified correctly",
		input:       33550337,
		ok:          true,
		expected:    ClassificationDeficient,
	},
	{
		description: "Edge case (no factors other than itself) is classified correctly",
		input:       1,
		ok:          true,
		expected:    ClassificationDeficient,
	},
	{
		description: "Zero is rejected (as it is not a positive integer)",
		input:       0,
		ok:          false,
	},
	{
		description: "Negative integer is rejected (as it is not a positive integer)",
		input:       -1,
		ok:          false,
	},
}

func TestGivesPositiveRequiredError(t *testing.T) {

	if _, err := Classify(0); err != ErrOnlyPositive {

		t.Fatalf("FAIL GivesPositiveRequiredError Expected error %q but got %q", ErrOnlyPositive, err)

	}

	t.Logf("PASS GivesPositiveRequiredError")

}

func TestClassifiesCorrectly(t *testing.T) {

	for _, c := range classificationTestCases {

		cat, err := Classify(c.input)

		switch {

		case err != nil:

			if c.ok {

				t.Fatalf("FAIL %s\nClassify(%d)\nExpected no error but got error %q", c.description, c.input, err)

			}

		case !c.ok:

			t.Fatalf("FAIL %s\nClassify(%d)\nExpected error but got %q", c.description, c.input, cat)

		case cat != c.expected:

			t.Fatalf("FAIL %s\nClassify(%d)\nExpected %q, got %q", c.description, c.input, c.expected, cat)

		}

		t.Logf("PASS %s", c.description)

	}

}

// Test that the classifications are not equal to each other.

// If they are equal, then the tests will return false positives.

func TestClassificationsNotEqual(t *testing.T) {

	classifications := []struct {
		class Classification

		name string
	}{

		{ClassificationAbundant, "ClassificationAbundant"},

		{ClassificationDeficient, "ClassificationDeficient"},

		{ClassificationPerfect, "ClassificationPerfect"},
	}

	for i, pair1 := range classifications {

		for j := i + 1; j < len(classifications); j++ {

			pair2 := classifications[j]

			if pair1.class == pair2.class {

				t.Fatalf("%s should not be equal to %s", pair1.name, pair2.name)

			}

		}

	}

}

func BenchmarkClassify(b *testing.B) {

	if testing.Short() {

		b.Skip("skipping benchmark in short mode.")

	}

	for i := 0; i < b.N; i++ {

		for _, c := range classificationTestCases {

			Classify(c.input)

		}

	}

}
