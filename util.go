package array

func toSnake(name string) string {
	var convert []byte
	for i, asc := range name {
		if asc >= 65 && asc <= 90 {
			asc += 32
			if i > 0 {
				convert = append(convert, 95)
			}
		}
		convert = append(convert, uint8(asc))
	}
	return string(convert)
}
