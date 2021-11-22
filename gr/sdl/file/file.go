package file

type File interface {
	ReadFloat64() float64
	ReadInt() int
	WriteFloat64(float64)
	WriteInt(int)
}
