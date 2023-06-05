package array

import (
	"fmt"
	"strings"
)

func Has[T string | int8 | int32 | int64](parent string, sub T) bool {
	list := strings.Split(parent, ",")
	s := fmt.Sprintf(`%v`, sub)
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}
