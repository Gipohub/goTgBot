package ytClient

import (
	"fmt"
	"log"
	"time"

	"github.com/tebeka/selenium"
)

func Pars() {
	// Укажите путь к ChromeDriver
	const (
		chromeDriverPath = "C:/goes/goTgBot/storage/chromedriver.exe"
		port             = 8080 // Порт для WebDriver
	)
	fmt.Println("starts")
	// Запускаем Selenium WebDriver
	service, err := selenium.NewChromeDriverService(chromeDriverPath, port)
	if err != nil {
		log.Fatalf("Ошибка при запуске ChromeDriver: %v", err)
	}
	defer service.Stop()

	// Подключаемся к WebDriver
	caps := selenium.Capabilities{"browserName": "chrome"}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		log.Fatalf("Ошибка подключения к WebDriver: %v", err)
	}
	defer wd.Quit()

	// Открываем веб-страницу
	if err := wd.Get("https://www.google.com"); err != nil {
		log.Fatalf("Ошибка при загрузке страницы: %v", err)
	}

	// Ищем элемент на странице (пример: поле поиска Google)
	elem, err := wd.FindElement(selenium.ByCSSSelector, "input[name='q']")
	if err != nil {
		log.Fatalf("Ошибка при поиске элемента: %v", err)
	}

	// Вводим текст в поле поиска
	if err := elem.SendKeys("Selenium WebDriver"); err != nil {
		log.Fatalf("Ошибка при вводе текста: %v", err)
	}

	// Ждем несколько секунд и закрываем браузер
	time.Sleep(5 * time.Second)
	fmt.Println("done")
}
