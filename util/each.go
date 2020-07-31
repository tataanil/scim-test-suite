package util

func (suite *Suite) ForEachMap(i interface{}, f func(map[string]interface{})) {
	switch i := i.(type) {
	case []interface{}:
		for _, v := range i {
			suite.ForEachMap(v, f)
		}
	case map[string]interface{}:
		f(i)
		for _, v := range i {
			suite.ForEachMap(v, f)
		}
	}
}
