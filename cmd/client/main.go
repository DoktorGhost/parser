package main

import (
	"fmt"
	"github.com/tebeka/selenium"
	"log"
	"parser/internal/config"
	"parser/internal/pars"
	"parser/internal/storage"
	"parser/internal/storage/csvRecord"
	"parser/internal/usecase"
	"strings"
	"time"
)

func main() {
	start := time.Now()

	configElements, err := config.ReadConfig("../../config.json")
	if err != nil {
		log.Fatal("Error reading config:", err)
	}

	urlSeller := configElements.UrlSeller
	city := configElements.DeliveryAddress.City
	street := configElements.DeliveryAddress.Street
	categorysName := configElements.Categories
	proxy := configElements.Proxy
	headless := configElements.Headless

	// Инициализация сервиса драйвера Chrome на порту 4444
	service, err := pars.Service("./chromedriver", 4444)
	if err != nil {
		log.Fatal(err)
	}
	defer service.Stop()

	driver, err := pars.NewDriver(proxy, headless)
	if err != nil {
		log.Fatal(err)
	}
	defer driver.Quit()

	//открываем окно максимально
	err = driver.MaximizeWindow("")
	if err != nil {
		log.Fatal("Ошибка развертывания окна на максимум:", err)
	}

	// Посещение целевой страницы
	err = driver.Get(urlSeller)
	if err != nil {
		log.Fatal("Ошибка открытия целевого сайта:", err)
	}

	// Создать карту для хранения названий категорий и их ссылок
	categoryMap := make(map[string]string)

	//Ожидание загрузки страницы
	err = pars.WhiteElement(driver, ".CatalogTree_root__IHeCU", "selector")
	if err != nil {
		log.Fatal("Error:", err)
	}

	//все категории
	category, err := driver.FindElement(selenium.ByCSSSelector, ".CatalogTree_root__IHeCU")
	if err != nil {
		log.Fatal("Error:", err)
	}

	// Найти все подкатегории
	subCat, err := category.FindElements(selenium.ByCSSSelector, ".CatalogTreeSectionCard_categories__4uYFm.CatalogTreeSectionCard_categories_hidden__k0Mib")
	if err != nil {
		log.Fatal("Error finding divs:", err)
	}

	// Запишем в мапу все категории и ссылки
	for _, cat := range subCat {
		urls, err := cat.FindElements(selenium.ByTagName, "a")
		if err != nil {
			log.Fatal("Error finding divs:", err)
		}
		for _, url := range urls {
			// Получить значение атрибута 'href' ссылки
			href, err := url.GetAttribute("href")
			if err != nil {
				log.Fatal("Error getting href attribute:", err)
			}

			// Найти вложенный тег 'span'
			spanElement, err := url.FindElement(selenium.ByCSSSelector, "span._text_7xv2z_4._text--type_p2SemiBold_7xv2z_97.CatalogTreeSectionCard_categoryName__DQfpi")
			if err != nil {
				log.Fatal("Error finding span element:", err)
			}

			// Получить текстовое содержимое элемента span с помощью GetAttribute("textContent")
			textContent, err := spanElement.GetAttribute("textContent")
			if err != nil {
				log.Fatal("Error getting textContent attribute:", err)
			}

			categoryMap[textContent] = href
		}
	}

	//////////////////////////////////////////////////////////////////////

	//ВВОД АДРЕСА

	//клик по кнопке
	element, err := driver.FindElement(selenium.ByXPATH, "//span[text()='Нет, другой']")
	if err != nil {
		log.Fatalf("Error finding the button: %v", err)
	}
	err = element.Click()
	if err != nil {
		log.Fatalf("Error clicking the button: %v", err)
	}

	// Ввести текст в поле "Город"
	cityInput, err := driver.FindElement(selenium.ByCSSSelector, "input[placeholder='Город']")
	if err != nil {
		log.Fatalf("Error finding the input field: %v", err)
	}

	err = cityInput.SendKeys(city)
	if err != nil {
		log.Fatalf("Error entering text: %v", err)
	}

	//выбор из списка
	elements, err := driver.FindElements(selenium.ByCSSSelector, ".Suggest_suggestItem__hOaW9 span._text_7xv2z_4._text--type_p1SemiBold_7xv2z_109")
	if err != nil {
		log.Fatal("Error:", err)
	}

	var targetElement selenium.WebElement
	for _, element := range elements {
		text, err := element.Text()
		if err != nil {
			log.Fatal("Error:", err)
		}
		if text == city {
			targetElement = element
			break
		}
	}

	if targetElement == nil {
		log.Fatalf("Error: Element with text %v not found", city)
	}

	// Клик по найденному городу
	if err := targetElement.Click(); err != nil {
		log.Fatal("Error:", err)
	}

	//ввод улицы
	streetInput, err := driver.FindElement(selenium.ByCSSSelector, "input[placeholder='Улица и дом']")
	if err != nil {
		log.Fatalf("Error finding the input field: %v", err)
	}

	// Ввести текст в поле "Улица и дом"
	err = streetInput.SendKeys(street)
	if err != nil {
		log.Fatalf("Error entering text: %v", err)
	}

	//Ожидание загрузки страницы
	err = pars.WhiteElement(driver, ".Suggest_suggestItems__wnlQV", "selector")
	if err != nil {
		log.Fatal("Error:", err)
	}

	//выбор из списка
	elements, err = driver.FindElements(selenium.ByCSSSelector, ".Suggest_suggestItems__wnlQV span._text_7xv2z_4._text--type_p1SemiBold_7xv2z_109")
	if err != nil {
		log.Fatal("Error:", err)
	}

	var streetElement selenium.WebElement

	for _, element := range elements {
		text, err := element.Text()
		if err != nil {
			log.Fatal("Error:", err)
		}
		if text == street {
			streetElement = element
			break
		}
	}

	if streetElement == nil {
		log.Fatalf("Error: Element with text %v not found", street)
	}

	// Клик по найденной улице
	if err := streetElement.Click(); err != nil {
		log.Fatal("Error:", err)
	}

	//Ожидание загрузки страницы и появления кнопок
	err = pars.WhiteElement(driver, "._button--size_m_10nio_88._button--theme_primary_10nio_56.AddressCreation_button__ow_WB", "selector")
	if err != nil {
		log.Fatal("Error:", err)
	}

	button, err := driver.FindElement(selenium.ByCSSSelector, "._button--size_m_10nio_88._button--theme_primary_10nio_56.AddressCreation_button__ow_WB button._control_10nio_4")
	if err != nil {
		log.Fatalf("Error finding the button: %v", err)
	}
	err = button.Click()
	if err != nil {
		log.Fatalf("Error clicking the button: %v", err)
	}

	err = pars.WhiteElement(driver, "._icon_1c7va_1", "selector")
	if err != nil {
		log.Fatal("Error:", err)
	}

	cardArr := []storage.Card{}

	for i, cat := range categorysName {
		card := storage.Card{}
		url := categoryMap[cat]
		if len(url) < 1 {
			continue
		}
		err = driver.Get(url)
		if err != nil {
			log.Fatal("Ошибка открытия целевого сайта:", i, err)
		}

		//ждем загрузки сайта
		err = pars.WhiteElement(driver, ".InlineSearch_content__9_P60", "selector")
		if err != nil {
			log.Fatal("Error:", err)
		}

		//обновление страницы
		err = driver.Refresh()
		if err != nil {
			log.Fatal("Error refreshing page:", err)
		}

		//все товары
		productsElements, err := driver.FindElement(selenium.ByCSSSelector, ".CategoryPage_container__fWbVL")
		if err != nil {
			log.Fatal("Error:", err)
		}
		//имя категории
		nameCat, err := productsElements.FindElement(selenium.ByCSSSelector, ".CategoryPage_categoryNameContainer__C35DT")
		if err != nil {
			log.Fatal("Error:", err)
		}
		nameText, err := nameCat.GetAttribute("textContent")
		if err != nil {
			log.Fatal("Error getting textContent attribute:", err)
		}
		card.Category = nameText

		//одна подкатегория с именем и карточками товаров
		subCats, err := productsElements.FindElements(selenium.ByCSSSelector, ".CategorySection_root__6Ai7Z")
		if err != nil {
			log.Fatal("Error:", err)
		}

		for _, subCat := range subCats {
			//имя подкатегории
			nameSubCat, err := subCat.FindElement(selenium.ByCSSSelector, ".CategorySection_header___0aG6 span")
			if err != nil {
				log.Fatal("Error:", err)
			}
			nameText, err := nameSubCat.Text()
			if err != nil {
				log.Fatal(err)
			}
			card.SubCategory = nameText

			cards, err := subCat.FindElements(selenium.ByCSSSelector, ".ProductsList_productList__XIJx_ a")
			if err != nil {
				log.Fatal("Error:", err)
			}
			//итерация по карточкам внутри субкатегории
			for _, oneCard := range cards {
				//ссылка
				href, err := oneCard.GetAttribute("href")
				if err != nil {
					log.Fatal("Error getting href attribute:", err)
				}
				card.Url = href

				//карточка с ценой и картинкой
				cardImgPrice, err := oneCard.FindElement(selenium.ByCSSSelector, ".ProductCard_root__OBGd_")
				if err != nil {
					log.Fatal("Error:", err)
				}

				cradImg, err := cardImgPrice.FindElement(selenium.ByCSSSelector, ".ProductCardImage_root__b96bY img")
				if err != nil {
					log.Fatal("Error:", err)
				}

				imgUrl, err := cradImg.GetAttribute("src")
				if err != nil {
					log.Fatal("Error getting imgUrl attribute:", err)
				}
				card.UrlImage = imgUrl

				name, err := cradImg.GetAttribute("alt")
				if err != nil {
					log.Fatal("Error getting name attribute:", err)
				}
				card.Name = name

				//цена
				pricesCard, err := cardImgPrice.FindElement(selenium.ByCSSSelector, ".ProductCard_actions__2AbGZ")
				if err != nil {
					log.Fatal("Error:", err)
				}

				//ProductCardActions_text__rHfOY
				priceCard, err := pricesCard.FindElements(selenium.ByCSSSelector, ".ProductCardActions_text__rHfOY span")
				if err != nil {
					log.Fatal("Error:", err)
				}

				if len(priceCard) == 3 {
					priceWithoutDiscont, err := priceCard[1].Text()
					if err != nil {
						log.Fatal(err)
					}
					card.PriceWithoutDiscount = priceWithoutDiscont

					price, err := priceCard[2].Text()
					if err != nil {
						log.Fatal(err)
					}
					parts := strings.Split(price, " ")

					card.Price = parts[0]
				} else {
					price, err := priceCard[0].Text()
					if err != nil {
						log.Fatal(err)
					}
					parts := strings.Split(price, " ")
					card.Price = parts[0]
					card.PriceWithoutDiscount = parts[0]
				}
				card.Address = fmt.Sprintf("%s, %s", city, street)
				cardArr = append(cardArr, card)
			}
		}
	}

	//экземпляр БД
	stor := csvRecord.NewCsvRecord()
	parser := usecase.NewUseCaseParser(stor)

	//запись в бд

	for _, card := range cardArr {
		parser.Add("../../EXPORT.csv", card)
	}

	//подсчет времени выполнения
	end := time.Now()
	duration := end.Sub(start)
	fmt.Println("Время выполнения:", duration)

}
