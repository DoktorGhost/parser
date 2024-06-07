package pars

import (
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"log"
	"time"
)

func Service(path string, port int) (*selenium.Service, error) {
	service, err := selenium.NewChromeDriverService(path, port)
	if err != nil {
		log.Fatal("Ошибка инициализации selenium: ", err)
		return nil, err
	}
	return service, nil
}

func NewDriver(proxy string, headless bool) (selenium.WebDriver, error) {

	caps := selenium.Capabilities{}
	chromeCaps := chrome.Capabilities{
		Args: []string{
			"--disable-gpu",
			"--blink-settings=imagesEnabled=false",
			"--disable-javascript",
			"--adblocker",
			"--hide-scrollbars",
			"--ignore-certificate-errors",
			"--user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		},
	}

	if proxy != "" {
		proxy = "--proxy-server=" + proxy
		chromeCaps.Args = append(chromeCaps.Args, proxy)
	}
	if !headless {
		chromeCaps.Args = append(chromeCaps.Args, "--headless=new")
	}

	caps.AddChrome(chromeCaps)

	// Создание нового удаленного клиента с указанными опциями
	driver, err := selenium.NewRemote(caps, "")
	if err != nil {
		log.Fatal("Ошибка создания клиента:", err)
		return nil, err
	}

	return driver, nil
}

func xpathName(xpath string) string {
	var by string
	switch xpath {
	case "XPATH":
		by = selenium.ByXPATH
	case "selector":
		by = selenium.ByCSSSelector
	case "TAG":
		by = selenium.ByTagName
	}
	return by
}

// Ожидание
func WhiteElement(driver selenium.WebDriver, selector string, xpath string) error {
	xpath = xpathName(xpath)
	err := driver.WaitWithTimeout(func(driver selenium.WebDriver) (bool, error) {
		lastProduct, _ := driver.FindElement(xpath, selector)
		if lastProduct != nil {
			return lastProduct.IsDisplayed()
		}
		return false, nil
	}, 10*time.Second)
	return err
}
