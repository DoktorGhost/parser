<p align="left">
      <img src="https://i.ibb.co/cYzQsPG/logoza-ru.png" alt="Project Logo" width="726">
</p>

# Мой первый парсер

*Программа парсит онлайн-магазин Самокат* 


1. Необходимо скачать chromedriver под установленную систему (по умолчанию стоит драйвер под Windows). 
Заходим на сайт https://googlechromelabs.github.io/chrome-for-testing/ и выбираем драйвер. Помещаем драйвер в директорию chromedriver в корневом каталоге.
2. Открываем файл *config.json* из корневого каталога и редактируем необходимые поля:


```json
config.json

{
  "urlSeller": "https://samokat.ru/",  //программа только под этот сайт, так что это менять не нужно 
  "deliveryAddress": {
    "city": "Ростов-на-Дону",  //так как в самокате нет магазинов, а только доставка, указываем город
    "street": "Орбитальная, 74"  //и улицу с номером дома через запятую
  },
  "categories": [  //пишем, какие категории нам нужны
    "Уход и макияж",  
    "Электроника",
    "Фрукты и ягоды"
  ],
  "proxy": "",  //если есть прокси-сервер - вписываем его сюда
  "headless": false  //оконный режим false-выключен, true-включен
  "port": 4444, //порт, на котором будет работать драйвер
  "chromedriver": "../../chromedriver/chromedriver", //путь к драйверу
  "export": "../../EXPORT.csv" //путь к файлу, в который будет происходить экспорт
}
```
2. Запускаем файл cmd/client/main.go, парсер начинает работу.
