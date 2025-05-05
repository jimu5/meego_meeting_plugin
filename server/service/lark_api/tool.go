package lark_api

func GetPtr[T any](t T) *T {
	return &t
}
