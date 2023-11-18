package f_sql

type FSQL[T interface{}] struct {
	items []T
}

func Set[T interface{}](data []T) *FSQL[T] {
	return &FSQL[T]{
		items: data,
	}
}

func (f *FSQL[T]) Where(exp func(e T) bool) *FSQL[T] {
	newItem := make([]T, 0)
	for _, item := range f.items {
		if exp(item) {
			newItem = append(newItem, item)
		}
	}

	return &FSQL[T]{items: newItem}
}

func (f *FSQL[T]) FirstWhere(exp func(e T) bool, def T) T {
	for _, item := range f.items {
		if exp(item) {
			return item
		}
	}

	return def
}

func (f *FSQL[T]) LastWhere(exp func(e T) bool, def T) T {
	for _, item := range f.items {
		if exp(item) {
			def = item
		}
	}
	return def
}

func (f *FSQL[T]) indexOf(exp func(e T) bool, def T) int {
	for index, item := range f.items {
		if exp(item) {
			return index
		}
	}

	return -1
}

func (f *FSQL[T]) Count() int {
	return len(f.items)
}

func (f *FSQL[T]) All() []T {
	return f.items
}
