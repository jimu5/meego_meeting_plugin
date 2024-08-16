package util

func GetPointerInfo[T any](s *T) T {
	if s == nil {
		t := new(T)
		return *t
	}
	return *s
}
