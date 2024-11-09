# Lardi-Trans API Client

Go-клиент для работы с API Lardi-Trans - платформы грузоперевозок.

## Установка

```bash
go get github.com/fentezi/lardi-trans-api
```

## Использование

### Инициализация клиента

```go
import "github.com/fentezi/lardi-trans-api"

config := larditrans.Config{
    APIKey:   "ваш-api-ключ",
    Language: "ru", // или "uk" для украинского языка
}

client := larditrans.NewClient(config)
```

### Создание заявки на перевозку груза

```go
ctx := context.Background()

request := &larditrans.CargoRequest{
    DateFrom:          "2024-11-10",
    DateTo:            "2024-11-11",
    PaymentValue:      1000,
    PaymentCurrencyID: 1, // ID валюты
    PaymentUnitID:     1, // ID единицы измерения
    ContentName:       "Электроника",
    SizeMass:          1000, // кг
    WaypointListSource: []larditrans.LoadParams{
        {
            TownID:      1,
            TownName:    "Киев",
            CountrySign: "UA",
        },
    },
    WaypointListTarget: []larditrans.LoadParams{
        {
            TownID:      2,
            TownName:    "Львов",
            CountrySign: "UA",
        },
    },
}

response, err := client.CreateCargo(ctx, request)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Создана заявка с ID: %d\n", response.ID)
```

### Получение справочных данных

```go
// Получение списка валют
currencies, err := client.GetCurrencies(ctx)

// Получение типов кузова
bodyTypes, err := client.GetBodyTypes(ctx)

// Получение типов загрузки
loadTypes, err := client.GetLoadTypes(ctx)
```

## Доступные методы

- `CreateCargo` - создание заявки на перевозку груза
- `GetAreas` - получение списка регионов
- `GetLoadTypes` - получение типов загрузки
- `GetPaymentTypes` - получение типов оплаты
- `GetPackageTypes` - получение типов упаковки
- `GetBodyTypes` - получение типов кузова
- `GetPaymentMoments` - получение моментов оплаты
- `GetCurrencies` - получение списка валют
- `GetUnits` - получение единиц измерения

## Конфигурация

Параметры конфигурации клиента:

- `BaseURL` - базовый URL API (по умолчанию "https://api.lardi-trans.com")
- `APIKey` - ваш API ключ (обязательный параметр)
- `Timeout` - таймаут запросов (по умолчанию 30 секунд)
- `Language` - язык ответов API ("ru" или "uk", по умолчанию "uk")

## Обработка ошибок

Клиент возвращает ошибки в формате `APIError` со следующими полями:
- `Status` - HTTP статус код
- `Error` - код ошибки
- `Message` - описание ошибки

```go
if err != nil {
    if apiErr, ok := err.(*larditrans.APIError); ok {
        fmt.Printf("API вернул ошибку: %s\n", apiErr.Message)
    }
}
```

## Лицензия

Этот проект распространяется под [лицензией MIT](./LICENSE). Подробности можно найти в файле LICENSE.
