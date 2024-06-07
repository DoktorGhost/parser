package parser

import (
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"log"
)

func Driver() selenium.WebDriver {
	service, err := selenium.NewChromeDriverService("./chromedriver", 4444)
	if err != nil {
		log.Fatal("Ошибка инициализации selenium:", err)
	}
	defer service.Stop()

	caps := selenium.Capabilities{}
	caps.AddChrome(chrome.Capabilities{
		Args: []string{
			//proxy,  //работа через прокси
			//"--headless", // безоконный режим
		},
	})

	// Создание нового удаленного клиента с указанными опциями
	driver, err := selenium.NewRemote(caps, "")
	if err != nil {
		log.Fatal("Ошибка создания клиента:", err)
	}

	return driver
}
