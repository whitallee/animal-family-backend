package task

import (
	"reflect"
	"testing"

	"github.com/whitallee/animal-family-backend/types"
)

// A task belongs to exactly one subject: "taskSubject" holds two nullable
// foreign keys and the store rejects any other combination.
func TestExactlyOneSubject(t *testing.T) {
	animal, enclosure := 3, 7

	cases := []struct {
		name        string
		animalId    *int
		enclosureId *int
		wantErr     bool
	}{
		{"animal only", &animal, nil, false},
		{"enclosure only", nil, &enclosure, false},
		{"neither", nil, nil, true},
		{"both", &animal, &enclosure, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := exactlyOneSubject(tc.animalId, tc.enclosureId)
			if tc.wantErr && err == nil {
				t.Error("expected an error, got none")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func taskFor(id int, animalId *int, enclosureId *int) *types.TaskWithSubject {
	return &types.TaskWithSubject{TaskId: id, AnimalId: animalId, EnclosureId: enclosureId}
}

func TestFilterTasksBySubject(t *testing.T) {
	a3, a4, e7 := 3, 4, 7
	tasks := []*types.TaskWithSubject{
		taskFor(1, &a3, nil),
		taskFor(2, &a4, nil),
		taskFor(3, nil, &e7),
		taskFor(4, &a3, nil),
	}

	cases := []struct {
		name        string
		animalId    *int
		enclosureId *int
		wantIds     []int
	}{
		{"no filter returns everything", nil, nil, []int{1, 2, 3, 4}},
		{"by animal", &a3, nil, []int{1, 4}},
		{"by enclosure", nil, &e7, []int{3}},
		{"animal with no tasks", &[]int{99}[0], nil, []int{}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := filterTasksBySubject(tasks, tc.animalId, tc.enclosureId)

			ids := make([]int, 0, len(got))
			for _, task := range got {
				ids = append(ids, task.TaskId)
			}

			if !reflect.DeepEqual(ids, tc.wantIds) {
				t.Errorf("got task ids %v, want %v", ids, tc.wantIds)
			}
		})
	}
}

// An enclosure task must not be matched by an animal filter that shares its id,
// which a naive comparison ignoring which pointer is set would do.
func TestFilterTasksBySubjectDoesNotConfuseSubjectTypes(t *testing.T) {
	id := 5
	tasks := []*types.TaskWithSubject{taskFor(1, nil, &id)}

	if got := filterTasksBySubject(tasks, &id, nil); len(got) != 0 {
		t.Errorf("animal filter matched an enclosure task: %v", got)
	}
}

func TestZeroIfNil(t *testing.T) {
	value := 9

	if got := zeroIfNil(nil); got != 0 {
		t.Errorf("zeroIfNil(nil) = %d, want 0", got)
	}
	if got := zeroIfNil(&value); got != 9 {
		t.Errorf("zeroIfNil(&9) = %d, want 9", got)
	}
}
