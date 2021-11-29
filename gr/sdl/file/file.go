package file

type File interface {
	ReadFloat64() float64
	ReadInt() int
	ReadInt64() int64
	WriteFloat64(float64)
	WriteInt(int)
	WriteInt64(int64)
	Open(string) error
	IsOpen() bool
	Close()
	Create(string)
}

func New() File {
	return nil
}
