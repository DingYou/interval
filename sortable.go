package interval

// baseSortable is the set of basic type which sortable
type baseSortable interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64 | ~string
}

// SortComparable defines a function to compare self to another object with same type
type SortComparable[T any] interface {
	// CompareTo returns
	CompareTo(other T) int
}

func Compare[T SortComparable[T]](e1, e2 T) int {
	return e1.CompareTo(e2)
}
