package protocol

import (
	"fmt"
	"strings"

	"github.com/faiyaz032/NotRedis/internal/store"
)

const (
	SimpleString = "+"
	Error        = "-"
	Integer      = ":"
	BulkString   = "$"
	Array        = "*"
)

type Value struct {
	Type  string
	Str   string
	Num   int
	Bulk  string
	Array []Value
	Null  bool
}

func (v Value) Marshal() []byte {
	if v.Null {
		switch v.Type {
		case BulkString:
			return []byte("$-1\r\n")
		case Array:
			return []byte("*-1\r\n")
		}
	}

	switch v.Type {
	case SimpleString:
		return []byte(fmt.Sprintf("+%s\r\n", v.Str))
	case Error:
		return []byte(fmt.Sprintf("-%s\r\n", v.Str))
	case Integer:
		return []byte(fmt.Sprintf(":%d\r\n", v.Num))
	case BulkString:
		return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(v.Bulk), v.Bulk))
	case Array:
		res := []byte(fmt.Sprintf("*%d\r\n", len(v.Array)))
		for _, val := range v.Array {
			res = append(res, val.Marshal()...)
		}
		return res
	default:
		return []byte{}
	}
}

func Execute(args []Value, st *store.Store) Value {
	if len(args) == 0 {
		return Value{Type: Error, Str: "ERR empty command"}
	}

	command := strings.ToUpper(args[0].Bulk)
	switch command {
	case "SET":
		if len(args) < 3 {
			return Value{Type: Error, Str: "ERR wrong number of arguments for 'set' command"}
		}
		st.Set(args[1].Bulk, args[2].Bulk)
		return Value{Type: SimpleString, Str: "OK"}
	case "GET":
		if len(args) < 2 {
			return Value{Type: Error, Str: "ERR wrong number of arguments for 'get' command"}
		}
		value, ok := st.Get(args[1].Bulk)
		if !ok {
			return Value{Type: BulkString, Null: true}
		}
		return Value{Type: BulkString, Bulk: value}
	default:
		return Value{Type: Error, Str: fmt.Sprintf("ERR unknown command '%s'", command)}
	}
}
