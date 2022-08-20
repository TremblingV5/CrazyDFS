package DNService

import "strconv"

func (id BlockID) ToString() string {
	str := strconv.FormatInt(int64(id), 10)
	res := ""

	for i := len(str); i <= 32-len(str); i++ {
		res += "0"
	}

	res += str
	return res
}

func (id BlockID) Next() BlockID {
	id += 1
	return id
}
