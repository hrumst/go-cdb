// nolint
package tools

import (
	"fmt"

	"github.com/golang/mock/gomock"
)

type doMatch[V any] struct {
	match  func(v V) bool
	object any
}

func DoMatch[V any](m func(v V) bool) gomock.Matcher {
	return &doMatch[V]{
		match: m,
	}
}

func (o *doMatch[V]) Matches(object any) bool {
	o.object = object
	v, ok := object.(V)
	if !ok {
		return false
	}

	return o.match(v)
}

func (o *doMatch[V]) String() string {
	return fmt.Sprintf("is matched to %v", o.object)
}
