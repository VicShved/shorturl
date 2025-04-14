package repository

var storage = make(map[string]string)
var userStorage = make(map[string]*map[string]string)

func GetStorage() (*map[string]string, *map[string]*map[string]string) {
	return &storage, &userStorage
}

type MemRepiository struct {
	mp      *map[string]string
	userMap *map[string]*map[string]string
}

func GetMemRepository() *MemRepiository {
	memstorage, userStorage := GetStorage()
	return &MemRepiository{mp: memstorage, userMap: userStorage}
}

func (s MemRepiository) Save(key string, value string, userID string) error {
	(*s.mp)[key] = value
	urlMap, ok := (*s.userMap)[userID]
	if !ok {
		var localMap = make(map[string]string)
		urlMap = &localMap
		(*s.userMap)[userID] = urlMap
	}
	(*urlMap)[key] = value
	return nil
}

func (s MemRepiository) Read(key string, userID string) (string, bool) {
	result, ok := (*s.mp)[key]
	return result, ok
}

func (s MemRepiository) Len() int {
	return len(*s.mp)
}

func (s MemRepiository) Ping() error {
	return nil
}

func (s MemRepiository) Batch(data *[]KeyLongURLStr, userID string) error {
	for _, element := range *data {
		err := s.Save(element.Key, element.LongURL, userID)
		if err != nil {
			return err
		}
	}
	return nil
}
