package repository

var storage = make(map[string]string)
var userStorage = make(map[string]*map[string]string)

// GetStorage()
func GetStorage() (*map[string]string, *map[string]*map[string]string) {
	return &storage, &userStorage
}

// MemRepiository struct
type MemRepiository struct {
	mp      *map[string]string
	userMap *map[string]*map[string]string
}

// GetMemRepository() *MemRepiository
func GetMemRepository() *MemRepiository {
	memstorage, userStorage := GetStorage()
	return &MemRepiository{mp: memstorage, userMap: userStorage}
}

// Save(key string, value string, userID string)
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

// Read(key string, userID string)
func (s MemRepiository) Read(key string, userID string) (string, bool, bool) {
	result, ok := (*s.mp)[key]
	// urlMap, ok := (*s.userMap)[userID]
	// if !ok {
	// 	return "", ok
	// }
	// result, ok := (*urlMap)[key]
	return result, ok, false
}

// ReadWithUser(key string, userID string)
func (s MemRepiository) ReadWithUser(key string, userID string) (string, bool, bool) {
	urlMap, ok := (*s.userMap)[userID]
	if !ok {
		return "", ok, false
	}
	result, ok := (*urlMap)[key]
	return result, ok, false
}

// Len()
func (s MemRepiository) Len() int {
	return len(*s.mp)
}

// Ping()
func (s MemRepiository) Ping() error {
	return nil
}

// Batch(data *[]KeyLongURLStr, userID string)
func (s MemRepiository) Batch(data *[]KeyLongURLStr, userID string) error {
	for _, element := range *data {
		err := s.Save(element.Key, element.LongURL, userID)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetUserUrls(userID string)
func (s MemRepiository) GetUserUrls(userID string) (*[]KeyOriginalURL, error) {
	var results []KeyOriginalURL
	uMap, ok := (*s.userMap)[userID]
	if !ok {
		return &results, nil
	}
	for key, original := range *uMap {
		results = append(results, KeyOriginalURL{Key: key, Original: original})
	}
	return &results, nil
}

// DelUserUrls(shortURLs *[]string, userID string)
func (s MemRepiository) DelUserUrls(shortURLs *[]string, userID string) error {
	return nil // TODO need realizaion
}

// Close memrep
func (s MemRepiository) Close() {

}

// UsersCount()
func (r MemRepiository) UsersCount() (int, error) {
	return 0, nil
}

// UrlsCount
func (r MemRepiository) UrlsCount() (int, error) {
	return 0, nil
}
