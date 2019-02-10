# Golang клиент для сервиса sms.ru #

[![Build Status](https://travis-ci.org/dex35/smsru.svg?branch=dev)](https://travis-ci.org/dex35/smsru)

Поддерживаемые методы:
- sms/send, sms/status, sms/cost, sms/status
- my/balance, my/limit, my/free, my/senders
- stoplist/add, stoplist/del, stoplist/get
- callback/add, callback/del, callback/get

## Использование ##
Установка:
```go
go get github.com/dex35/smsru
```
Импорт:
```go
import "github.com/dex35/smsru"
```

## Пример использования ##
```go
package main

import (
	"log"
	"github.com/dex35/smsru"
)

func main() {
	smsClient := smsru.CreateClient("API_KEY")

	// Отправка сообщения
	phone := "номер телеофна в формате 79991112233"
	sms := smsru.CreateSMS(phone, "тестовое сообщение")
	sendedsms, err := smsClient.SmsSend(sms)
	if err != nil {
		log.Panic(err)
	} else {
		log.Printf("Статус запроса: %s, Статус-код выполнения: %d (%s), Баланс: %f", sendedsms.Status, sendedsms.StatusCode, smsru.GetErrorByCode(sendedsms.StatusCode), sendedsms.Balance)
		log.Printf("Сообщение: %s, Статус-код выполнения: %d (%s), Идентификатор: %s, Описание ошибки: %s", sendedsms.Sms[phone].Status, sendedsms.Sms[phone].StatusCode, smsru.GetErrorByCode(sendedsms.StatusCode), sendedsms.Sms[phone].SmsId, sendedsms.Sms[phone].StatusText)
	}

	// Запрос баланса
	balance, err := smsClient.MyBalance()
	if err != nil {
		log.Panic(err)
	} else {
		log.Printf("Статус запроса: %s, Статус-код выполнения: %d (%s), Баланс: %f", balance.Status, balance.StatusCode, smsru.GetErrorByCode(balance.StatusCode), balance.Balance)
	}
}
```