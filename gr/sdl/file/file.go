package file

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

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
	return &file{}
}

type file struct {
	fd      *os.File
	scanner *bufio.Scanner
}

func (f *file) readString() string {
	f.scanner.Scan()
	t := f.scanner.Text()
	if t == "" {
		panic("nothing to read from file")
	}
	return t
}

func (f *file) ReadFloat64() float64 {

	i, err := strconv.ParseFloat(f.readString(), 64)
	if err != nil {
		panic(err)
	}
	return i
}

func (f *file) ReadInt() int {
	f.scanner.Scan()
	i, err := strconv.Atoi(f.readString())
	if err != nil {
		panic(err)
	}
	return i
}

func (f *file) ReadInt64() int64 {
	f.scanner.Scan()
	i, err := strconv.ParseInt(f.readString(), 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}

func (f *file) WriteFloat64(d float64) {
	f.fd.WriteString(fmt.Sprintf("%f\n", d))
}

func (f *file) WriteInt(i int) {
	f.fd.WriteString(fmt.Sprintf("%d\n", i))
}

func (f *file) WriteInt64(i int64) {
	f.fd.WriteString(fmt.Sprintf("%d\n", i))
}

func (f *file) Open(name string) error {
	fd, err := os.OpenFile(name, os.O_RDONLY, 0)
	f.fd = fd
	f.scanner = bufio.NewScanner(f.fd)
	return err
}
func (f *file) IsOpen() bool { return f.fd != nil }

func (f *file) Close() { f.fd.Close() }

func (f *file) Create(name string) {
	fd, err := os.Create(name)
	f.fd = fd
	if err != nil {
		panic(err)
	}
}
