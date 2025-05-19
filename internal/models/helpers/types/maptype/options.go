package maptype

type MapNestedObjectOfOption[T any] func(*MapNestedObjectOfOptions[T])

type MapNestedObjectOfOptions[T any] struct {
	SemanticEqualityFunc mapSemanticEqualityFunc[T]
}

func WithSemanticEqualityFunc[T any](f mapSemanticEqualityFunc[T]) MapNestedObjectOfOption[T] {
	return func(o *MapNestedObjectOfOptions[T]) {
		o.SemanticEqualityFunc = f
	}
}

func newMapNestedObjectOfOptions[T any](options ...MapNestedObjectOfOption[T]) *MapNestedObjectOfOptions[T] {
	opts := &MapNestedObjectOfOptions[T]{}
	for _, opt := range options {
		opt(opts)
	}
	return opts
}
