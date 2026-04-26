package protocol

import (
	"strings"

	"github.com/faiyaz032/NotRedis/internal/store"
)

func Execute(line string, st *store.Store) string {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return "ERR: Empty Command\n"
	}

	switch strings.ToUpper(fields[0]) {
	case "PUT":
		if len(fields) < 3 {
			return "ERR: Usage PUT <key> <value>\n"
		}
		st.Set(fields[1], fields[2])
		return "OK\n"
	case "GET":
		if len(fields) < 2 {
			return "ERR: Usage GET <key>\n"
		}
		value, ok := st.Get(fields[1])
		if !ok {
			return "(nil)\n"
		}
		return value + "\n"
	default:
		return "ERR: Unknown Command\n"
	}
}
