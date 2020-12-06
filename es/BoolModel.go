package es

type H map[string]interface{}

type ESBool struct {
	Must   []H `json:"must"`
	Should []H `json:"should"`
}

func CreateByQueryMap(queryMap H) ESBool {
	var result ESBool
	for k, v := range queryMap {
		result.AddMustMatch(k, v)
	}
	return result
}

func (b *ESBool) AddMustMatch(key string, value interface{}) {
	b.AddMust("match", key, value)
}

func (b *ESBool) AddMust(condition string, key string, value interface{}) {
	b.Must = append(b.Must, createItem(condition, key, value))
}

func (b *ESBool) AddShouldMatch(key string, value interface{}) {
	b.AddShould("match", key, value)
}

func (b *ESBool) AddShould(condition string, key string, value interface{}) {
	b.Should = append(b.Should, createItem(condition, key, value))
}

func createItem(condition string, key string, value interface{}) H {
	return H{
		condition: H{
			key: value,
		},
	}
}
