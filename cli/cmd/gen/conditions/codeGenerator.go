package conditions

type CodeGenerator[T any] interface {
	Into(file *File) error
	ForEachField(file *File, fields []Field) []T
}
