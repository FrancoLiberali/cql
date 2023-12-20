package conditions

// TODO que solo sea gen
type CodeGenerator[T any] interface {
	Into(file *File) error
	ForEachField(file *File, fields []Field) []T
}
