package file

type File interface {
	ReadFloat64() float64
	ReadInt() int
	WriteFloat64(float64)
	WriteInt(int)
	Open(string) error
	IsOpen() bool
	Close()
	Create(string)
}

func New() File {
	return nil
}
