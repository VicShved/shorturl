package repository

type MemRepiository struct {
	mp *map[string]string
}

func GetMemRepository(mp *map[string]string) *MemRepiository {
	return &MemRepiository{mp: mp}
}

func (s MemRepiository) Save(key string, value string) error {
	(*s.mp)[key] = value
	return nil
}

func (s MemRepiository) Read(key string) (string, bool) {
	result, ok := (*s.mp)[key]
	return result, ok
}

func (s MemRepiository) Len() int {
	return len(*s.mp)
}

func (s MemRepiository) Ping() error {
	return nil
}

func (s MemRepiository) Batch(data *[]KeyLongURLStr) error {
	for _, element := range *data {
		err := s.Save(element.Key, element.LongURL)
		if err != nil {
			return err
		}
	}
	return nil
}
