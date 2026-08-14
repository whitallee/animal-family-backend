package enclosure

import (
	"reflect"
	"testing"
)

// dedupe is what makes comparing RowsAffected against len(wanted) meaningful.
// Without it, a repeated id inflates the expected count and a valid request
// would be rejected as unowned.
func TestDedupe(t *testing.T) {
	cases := []struct {
		name string
		in   []int
		want []int
	}{
		{"already unique", []int{3, 1, 2}, []int{3, 1, 2}},
		{"repeats collapse", []int{5, 5, 5}, []int{5}},
		{"keeps first-seen order", []int{4, 2, 4, 9, 2}, []int{4, 2, 9}},
		{"empty", []int{}, []int{}},
		{"nil", nil, []int{}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := dedupe(tc.in)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("dedupe(%v) = %v, want %v", tc.in, got, tc.want)
			}
		})
	}
}

// parseCascade replaces v1's separate /withtasks and /withanimalsandtasks
// delete routes. Getting it wrong silently deletes the wrong amount of data, so
// each accepted spelling is pinned here.
func TestParseCascadeAcceptedValues(t *testing.T) {
	cases := map[string]cascadeMode{
		"":                    cascadeNone,
		"tasks":               cascadeTasks,
		"animals,tasks":       cascadeAnimalsAndTasks,
		"tasks,animals":       cascadeAnimalsAndTasks,
		" tasks , animals ":   cascadeAnimalsAndTasks,
		"tasks,":              cascadeTasks,
		"animals,tasks,tasks": cascadeAnimalsAndTasks,
	}

	for input, want := range cases {
		t.Run(input, func(t *testing.T) {
			got, err := parseCascade(input)
			if err != nil {
				t.Fatalf("unexpected error for %q: %v", input, err)
			}
			if got != want {
				t.Errorf("parseCascade(%q) = %v, want %v", input, got, want)
			}
		})
	}
}

func TestParseCascadeRejectsUnknownValues(t *testing.T) {
	for _, input := range []string{"bogus", "animal", "task", "tasks;animals", "all"} {
		t.Run(input, func(t *testing.T) {
			if _, err := parseCascade(input); err == nil {
				t.Errorf("expected an error for %q", input)
			}
		})
	}
}

// Deleting an enclosure's animals while keeping their tasks would leave those
// tasks pointing at animals that no longer exist, and no store method does it.
func TestParseCascadeRejectsAnimalsWithoutTasks(t *testing.T) {
	if _, err := parseCascade("animals"); err == nil {
		t.Error("expected cascade=animals alone to be rejected")
	}
}
