package journal

import (
	"fmt"
	"slices"

	"github.com/qlustra/intel/logic"
)

const (
	VisitNotRecordedCode logic.CheckCode = "visit_not_recorded"
	VisitOutOfOrderCode  logic.CheckCode = "visit_out_of_order"
)

type explorationState int

const (
	visiting = iota
	visited
)

type Journal[I comparable, V any] struct {
	visited map[I]explorationState
	catalog map[I]V
	order   []V
	trace   []I
	failure error
}

// Err reports the first failure recorded for this session.
func (j *Journal[I, V]) Err() error {
	return j.failure
}

// Fail records a session failure. Callers must not reuse its partial results.
func (j *Journal[I, V]) Fail(err error) {
	if j.failure == nil {
		j.failure = err
	}
}

func New[I comparable, V any]() *Journal[I, V] {
	return &Journal[I, V]{}
}

func (j *Journal[I, V]) ensure() {
	if j.visited == nil {
		j.visited = make(map[I]explorationState)
	}
	if j.catalog == nil {
		j.catalog = make(map[I]V)
	}
	if j.order == nil {
		j.order = make([]V, 0)
	}
	if j.trace == nil {
		j.trace = make([]I, 0)
	}
}

func (j *Journal[I, V]) Visit(identity I, value V) {
	if j.Seen(identity) {
		return
	}
	j.ensure()
	j.visited[identity] = visiting
	j.trace = append(j.trace, identity)
	j.order = append(j.order, value)
	j.catalog[identity] = value
}

func (j *Journal[I, V]) Existing(identity I) (V, bool) {
	value, ok := j.catalog[identity]
	return value, ok
}

func (j *Journal[I, V]) Finish(identity I) logic.Check {
	if !j.Recorded(identity) {
		return logic.Nope(VisitNotRecordedCode, fmt.Sprintf("cannot finish %v: no visit was recorded", identity))
	}
	if j.Visited(identity) {
		return logic.Ok()
	}
	if len(j.trace) == 0 || j.trace[len(j.trace)-1] != identity {
		return logic.Nope(VisitOutOfOrderCode, fmt.Sprintf("cannot finish %v before its nested visits", identity))
	}
	j.visited[identity] = visited
	j.trace = j.trace[:len(j.trace)-1]
	return logic.Ok()
}

func (j *Journal[I, V]) Traced(identity I) bool {
	return slices.Contains(j.trace, identity)
}

func (j *Journal[I, V]) Recorded(identity I) bool {
	_, exists := j.visited[identity]
	return exists
}

func (j *Journal[I, V]) Visiting(identity I) bool {
	state, exists := j.visited[identity]
	return exists && state == visiting
}

func (j *Journal[I, V]) Seen(identity I) bool {
	state, exists := j.visited[identity]
	return exists && (state == visiting || state == visited)
}

func (j *Journal[I, V]) Visited(identity I) bool {
	state, exists := j.visited[identity]
	return exists && state == visited
}
