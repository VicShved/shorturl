package app

var storage = make(map[string]string)

func GetStorage() *map[string]string {
	return &storage
}
