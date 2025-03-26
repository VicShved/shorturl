package repository

type SaverReaderMem struct {
	mp *map[string]string
}

func NewSaverReaderMem(mp *map[string]string) *SaverReaderMem {
	return &SaverReaderMem{mp: mp}
}

func (s *SaverReaderMem) Save(key string, value string) error {
	(*s.mp)[key] = value
	return nil
}

func (s *SaverReaderMem) Read(key string) (string, bool) {
	result, ok := (*s.mp)[key]
	return result, ok
}

func (s *SaverReaderMem) Len() int {
	return len(*s.mp)
}
