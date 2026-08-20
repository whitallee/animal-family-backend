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
		"animals":             cascadeAnimalsAndTasks,
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

// "animals" means animals and their tasks. Keeping the tasks would leave them
// pointing at animals that no longer exist, so there is no other reading, and
// spelling it "animals,tasks" must mean the same thing.
func TestParseCascadeAnimalsImpliesTasks(t *testing.T) {
	withoutTasks, err := parseCascade("animals")
	if err != nil {
		t.Fatalf("cascade=animals should be accepted: %v", err)
	}

	spelledOut, err := parseCascade("animals,tasks")
	if err != nil {
		t.Fatalf("cascade=animals,tasks should be accepted: %v", err)
	}

	if withoutTasks != spelledOut || withoutTasks != cascadeAnimalsAndTasks {
		t.Errorf("both spellings should delete animals and tasks, got %v and %v", withoutTasks, spelledOut)
	}
}
