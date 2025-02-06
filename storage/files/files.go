package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/Gipohub/goTgBot/lib/e"
	"github.com/Gipohub/goTgBot/storage"
)

type Storage struct {
	basePath string
}

// у всех права на чтение изапись (0х8 система)
const defaultPerm = 0774

// ошибка в переменной чтобы её ?можно было проверить снаружи?
// например чтобы сказать пользователю что ничего не сохранено пока

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

//func (s Storage) Pars() {
//	ytClient.Pars()
//}

// сохранение файла
func (s Storage) Save(page *storage.Page) (err error) {
	//определяем способ обработки ошибок
	defer func() { e.Wrap("cant save page", err) }()

	//формируем путь до дериктории сохран файла
	fPath := filepath.Join(s.basePath, page.UserName)

	if err := os.MkdirAll(fPath, defaultPerm); err != nil {
		return err
	}

	//формируем имя файла через хэш
	fName, err := fileName(page)
	if err != nil {
		return err
	}

	//дописываем имя файла к пути
	fPath = filepath.Join(fPath, fName)

	//создаем файл
	file, err := os.Create(fPath)
	if err != nil {
		return err
	}
	// os.file.close возвращает ошибку, но тут не требуется ее обрабатывать
	// хорошей практикой будет показать что ошибка есть.
	//Конструкция defer не позволяет создать неименованую переменную,
	// поэтому оборачиваем в анонимную функцию
	defer func() { _ = file.Close() }()

	//записываем в него страницу в нужном формате
	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}

	return nil
}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() { e.Wrap("cant save page", err) }()

	path := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, storage.ErrNoSavedPages
	}

	rand.New(rand.NewSource(time.Now().UnixNano()))
	n := rand.Intn(len(files))

	file := files[n]

	return s.decodePage(filepath.Join(path, file.Name()))
}

func (s Storage) Remove(p *storage.Page) error {
	fileName, err := fileName(p)
	if err != nil {
		return e.Wrap("cant remove file", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	if err := os.Remove(path); err != nil {
		msg := fmt.Sprintf("cant remove file %s", path)
		return e.Wrap(msg, err)
	}

	return nil
}

func (s Storage) IsExists(p *storage.Page) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		return false, e.Wrap("cant check if file exsists", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("cant check if file %s exsists", path)

		return false, e.Wrap(msg, err)
	}

	return true, nil
}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, e.Wrap("cant decode page", err)
	}
	defer func() { _ = f.Close() }()

	var p storage.Page

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, e.Wrap("cant decode page", err)
	}

	return &p, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
