package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

// Define RESP protocol types
// RESP (REdis Serialization Protocol) uses specific characters to denote different types of values.
const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

// Value represents a RESP value.
type Value struct {
	typ   string
	str   string
	num   int
	bulk  string
	array []Value
}

type Resp struct {
	reader *bufio.Reader
}

// NewResp creates and returns a new Resp instance that reads from the provided io.Reader.
// It initializes a buffered reader to efficiently read RESP formatted data.
func NewResp(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}

// readLine reads a line from the RESP protocol.
// It reads bytes until it encounters a CRLF sequence,
// and returns the line without the CRLF, the number of bytes read,
// and an error if any occurs.
func (r *Resp) readLine() (line []byte, n int, err error) {
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		n += 1
		line = append(line, b)
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}
	return line[:len(line)-2], n, nil
}

// readInteger reads an integer from the RESP protocol.
// It reads a line from the reader, parses it as an integer, and returns the integer
// along with the number of bytes read and an error if any occurs.
func (r *Resp) readInteger() (x int, n int, err error) {
	line, n, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, n, err
	}
	return int(i64), n, nil
}

// Read reads a RESP value from the reader.
// It reads the type of the value and then calls the appropriate method
// to read the value based on its type.
// It returns a Value containing the parsed value or an error if any occurs.
func (r *Resp) Read() (Value, error) {
	_type, err := r.reader.ReadByte()

	if err != nil {
		return Value{}, err
	}

	switch _type {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	default:
		fmt.Printf("Unknown type: %v", string(_type))
		return Value{}, nil
	}
}

// readArray reads an array from the RESP protocol.
// It reads the length of the array, then reads each element in the array,
// and returns a Value containing the array or an error if any occurs.
func (r *Resp) readArray() (Value, error) {
	v := Value{}
	v.typ = "array"

	// read length of array
	len, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	// foreach line, parse and read the value
	v.array = make([]Value, 0)
	for i := 0; i < len; i++ {
		val, err := r.Read()
		if err != nil {
			return v, err
		}

		// append parsed value to array
		v.array = append(v.array, val)
	}

	return v, nil
}

// readBulk reads a bulk string from the RESP protocol.
// It reads the length of the bulk string, then reads the bulk data,
// and finally reads the trailing CRLF.
// It returns a Value containing the bulk string or an error if any occurs.
func (r *Resp) readBulk() (Value, error) {
	v := Value{}

	v.typ = "bulk"

	len, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	bulk := make([]byte, len)

	r.reader.Read(bulk)

	v.bulk = string(bulk)

	// Read the trailing CRLF
	r.readLine()

	return v, nil
}

// Marshal value to bytes
// Marshal converts a Value to its RESP byte representation.
// It checks the type of the value and calls the appropriate marshal method
// based on the type (array, bulk, string, null, or error).
// It returns the marshaled bytes.
// If the type is unknown, it returns an empty byte slice.
// This function is used to prepare the value for transmission over the network or for storage.
func (v Value) Marshal() []byte {
	switch v.typ {
	case "array":
		return v.marshalArray()
	case "bulk":
		return v.marshalBulk()
	default:
		return []byte{}
	case "string":
		return v.marshalString()
	case "null":
		return v.marshalNull()
	case "error":
		return v.marshalError()
	}
}

// marshalString marshals string array to RESP format.
// It prepends the type identifier for a string, appends the string value,
// and adds the CRLF sequence at the end.
// It returns the marshaled bytes.
func (v Value) marshalString() []byte {
	var bytes []byte
	bytes = append(bytes, STRING)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

// marshalBulk marshals a bulk string to RESP format.
// It prepends the type identifier for a bulk string, appends the length of the bulk
// string, adds the bulk string itself, and appends the CRLF sequence at the end.
// It returns the marshaled bytes.
// The length is represented as a string followed by CRLF.
func (v Value) marshalBulk() []byte {
	var bytes []byte
	bytes = append(bytes, BULK)
	bytes = append(bytes, []byte(strconv.Itoa(len(v.bulk)))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, v.bulk...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

// marshalArray marshals an array to RESP format.
// It prepends the type identifier for an array, appends the length of the array,
// and then marshals each element in the array, appending them to the result.
// It returns the marshaled bytes.
func (v Value) marshalArray() []byte {
	var bytes []byte
	bytes = append(bytes, ARRAY)
	bytes = append(bytes, []byte(strconv.Itoa(len(v.array)))...)
	bytes = append(bytes, '\r', '\n')
	for _, v := range v.array {
		bytes = append(bytes, v.Marshal()...)
	}
	return bytes
}

// marshalError marshals an error message to RESP format.
// It prepends the type identifier for an error, appends the error message,
// and adds the CRLF sequence at the end.
// It returns the marshaled bytes.
func (v Value) marshalError() []byte {
	var bytes []byte
	bytes = append(bytes, ERROR)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

// marshalNull marshals a null value to RESP format.
// It uses the RESP representation for null, which is a special case.
func (v Value) marshalNull() []byte {
	return []byte("$-1\r\n")
}

// Writer is a struct that provides methods to write RESP formatted data to an io.Writer.
// It encapsulates an io.Writer and provides a method to write RESP values.
type Writer struct {
	writer io.Writer
}

// NewWriter creates and returns a new Writer that writes RESP (REdis Serialization Protocol)
// formatted data to the provided io.Writer.
func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: w}
}

// Write writes a RESP value to the underlying io.Writer.
// It marshals the Value to its RESP byte representation and writes it to the writer.
func (w *Writer) Write(v Value) error {
	var bytes = v.Marshal()

	_, err := w.writer.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}
