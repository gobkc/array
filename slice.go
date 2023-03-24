package array

func MakeSlice[T any](dest ...T) []T {
	objects := make([]T, len(dest))
	objects = append(objects, dest...)
	return objects
}
