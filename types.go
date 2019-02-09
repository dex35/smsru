package smsru

import (
	"net/http"
	"time"
)

// Клиент
type SmsClient struct {
	ApiId string       `json:"api_id"`
	Http  *http.Client `json:"-"`
	Debug bool         `json:"-"`
}

// Отправляемое сообщение
type Sms struct {
	To        string            `json:"to"`
	Text      string            `json:"text"`
	Translit  bool              `json:"translit"`
	Multi     map[string]string `json:"multi"`
	From      string            `json:"from"`
	Time      time.Time         `json:"time"`
	Test      bool              `json:"test"`
	PartnerId int               `json:"partner_id"`
}

// Структуры для генерации ответов от api
type SendedCostStruct struct {
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	SmsId      string `json:"sms_id"`
	StatusText string `json:"status_text"`
}

type SendedSms struct {
	Status     string                   `json:"status"`
	StatusCode int                      `json:"status_code"`
	Sms        map[string]SmsCostStruct `json:"sms"`
	Balance    float32                  `json:balance`
}

type Free struct {
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	TotalFree  int    `json:"total_free"`
	UsedToday  int    `json:"used_today"`
}

type SmsCostStruct struct {
	Status     string  `json:"status"`
	StatusCode int     `json:"status_code"`
	Cost       float32 `json:"cost"`
	Sms        int     `json:"sms"`
}

type Cost struct {
	Status     string                   `json:"status"`
	StatusCode int                      `json:"status_code"`
	Sms        map[string]SmsCostStruct `json:"sms"`
	TotalCost  float32                  `json:"total_cost"`
	TotalSms   int                      `json:"total_sms"`
}

type Balance struct {
	Status     string  `json:"status"`
	StatusCode int     `json:"status_code"`
	Balance    float32 `json:"balance"`
}

type Limit struct {
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	TotalLimit string `json:"total_limit"`
	UsedToday  int    `json:"used_today"`
}
