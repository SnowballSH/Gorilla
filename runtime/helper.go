package runtime

import "fmt"

func getElement(k []BaseObject, index int) (BaseObject, error) {
	if index >= len(k) {
		return nil, fmt.Errorf("missing argument #%d", index+1)
	}

	i := k[index]
	return i, nil
}
