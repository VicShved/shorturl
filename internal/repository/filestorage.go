package repository

import (
	"bufio"
	"encoding/json"
	"os"
	"strconv"

	"github.com/VicShved/shorturl/internal/logger"
	"go.uber.org/zap"
)

// Element struct
type Element struct {
	ID       string `json:"id,omitempty"`
	Short    string `json:"short_url"`
	Original string `json:"original_url"`
	UserID   string `json:"user_id"`
}

// Consumer struct
type Consumer struct {
	file    *os.File
	scanner *bufio.Scanner
}

// NewConsumer(filename string)
func NewConsumer(filename string) (*Consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return &Consumer{file: file, scanner: bufio.NewScanner(file)}, nil
}

// ReadElement()
func (c *Consumer) ReadElement() (*Element, error) {
	if !c.scanner.Scan() {
		return nil, c.scanner.Err()
	}
	data := c.scanner.Bytes()
	element := new(Element)
	err := json.Unmarshal(data, element)
	if err != nil {
		return nil, err
	}
	return element, nil
}

// Close() error
func (c *Consumer) Close() error {
	return c.file.Close()
}

// InitFromFile(filename string, storage *RepoInterface)
func InitFromFile(filename string, storage *RepoInterface) error {
	logger.Log.Info("InitFromFile", zap.String("filename", filename))
	consumer, err := NewConsumer(filename)
	if err != nil {
		return err
	}

	elem, err := consumer.ReadElement()
	if err != nil {
		return err
	}
	for elem != nil {
		err = (*storage).Save(elem.Short, elem.Original, elem.UserID)
		if err != nil {
			return err
		}
		elem, err = consumer.ReadElement()
		if err != nil {
			return err
		}
	}
	return nil
}

// Producer struct
type Producer struct {
	file *os.File
}

// NewProducer(filename string)
func NewProducer(filename string) (*Producer, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return &Producer{file: file}, nil
}

// WriteElement(elem *Element)
func (p *Producer) WriteElement(elem *Element) error {
	data, err := json.Marshal(*elem)
	logger.Log.Debug("WriteElement", zap.String("data", string(data)))
	if err != nil {
		return err
	}
	data = append(data, '\n')
	_, err = p.file.Write(data)
	return err
}

// FileRepository struct
type FileRepository struct {
	sr       *MemRepiository
	Filename string
	Producer *Producer
}

// GetFileRepository(filenme string)
func GetFileRepository(filenme string) *FileRepository {
	repo := &FileRepository{
		sr:       GetMemRepository(),
		Filename: filenme,
	}
	repo.InitFromFile()
	repo.InitSaveFile()
	return repo
}

// InitSaveFile()
func (r *FileRepository) InitSaveFile() error {

	err := error(nil)
	r.Producer, err = NewProducer(r.Filename)
	if err != nil {
		return err
	}
	return nil
}

// Save(short, original string, userID string)
func (r FileRepository) Save(short, original string, userID string) error {
	_, ok, _ := r.sr.ReadWithUser(short, userID)
	if !ok {
		if r.Producer != nil {
			logger.Log.Debug("Save to FILE", zap.String("short", short), zap.String("original", original), zap.String("userID", userID))
			id := r.sr.Len() + 1
			err := r.Producer.WriteElement(&Element{ID: strconv.Itoa(id), Short: short, Original: original, UserID: userID})
			if err != nil {
				return err
			}
		}
	}
	return r.sr.Save(short, original, userID)
}

// Read(short string, userID string)
func (r FileRepository) Read(short string, userID string) (string, bool, bool) {
	return (*r.sr).Read(short, userID)
}

// Ping()
func (r FileRepository) Ping() error {
	return r.sr.Ping()
}

// Len()
func (r FileRepository) Len() int {
	return r.sr.Len()
}

// Batch(data *[]KeyLongURLStr, userID string)
func (r FileRepository) Batch(data *[]KeyLongURLStr, userID string) error {
	for _, element := range *data {
		err := r.Save(element.Key, element.LongURL, userID)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetUserUrls(userID string)
func (r FileRepository) GetUserUrls(userID string) (*[]KeyOriginalURL, error) {
	return r.sr.GetUserUrls(userID)
}

// InitFromFile()
func (r FileRepository) InitFromFile() error {
	logger.Log.Debug("InitFromFile", zap.String("filename", r.Filename))
	consumer, err := NewConsumer(r.Filename)
	if err != nil {
		return err
	}

	elem, err := consumer.ReadElement()
	if err != nil {
		return err
	}
	for elem != nil {
		err = r.Save(elem.Short, elem.Original, elem.UserID)
		if err != nil {
			return err
		}
		elem, err = consumer.ReadElement()
		if err != nil {
			return err
		}
	}
	return nil
}

// DelUserUrls(shortURLs *[]string, userID string)
func (r FileRepository) DelUserUrls(shortURLs *[]string, userID string) error {
	return nil // TODO need realizaion
}

// Close file rep
func (r FileRepository) Close() {

}

// UsersCount()
func (r FileRepository) UsersCount() (int, error) {
	return 0, nil
}

// UrlsCount
func (r FileRepository) UrlsCount() (int, error) {
	return 0, nil
}
