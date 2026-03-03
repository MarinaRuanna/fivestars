package domain

type Builder[T any] struct {
	Value T
}

func (b *Builder[T]) Build() T {
	return b.Value
}

func (b *Builder[T]) BuilderPointer() *T {
	build := b.Build()
	return &build
}

func NewBuilder[T any](inicialValue T) *Builder[T] {
	return &Builder[T]{
		Value: inicialValue,
	}
}
