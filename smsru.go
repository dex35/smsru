package smsru

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Api URL
const API_URL = "https://sms.ru"

var codeStatus map[int]string = map[int]string{
	-1:  "Сообщение не найдено",
	100: "Запрос выполнен или сообщение находится в нашей очереди",
	101: "Сообщение передается оператору",
	102: "Сообщение отправлено (в пути)",
	103: "Сообщение доставлено",
	104: "Не может быть доставлено: время жизни истекло",
	105: "Не может быть доставлено: удалено оператором",
	106: "Не может быть доставлено: сбой в телефоне",
	107: "Не может быть доставлено: неизвестная причина",
	108: "Не может быть доставлено: отклонено",
	110: "Сообщение прочитано",
	150: "Не может быть доставлено: не найден маршрут на данный номер",
	200: "Неправильный api_id",
	201: "Не хватает средств на лицевом счету",
	202: "Неправильно указан номер телефона получателя, либо на него нет маршрута",
	203: "Нет текста сообщения",
	204: "Имя отправителя не согласовано с администрацией",
	205: "Сообщение слишком длинное (превышает 8 СМС)",
	206: "Будет превышен или уже превышен дневной лимит на отправку сообщений",
	207: "На этот номер нет маршрута для доставки сообщений",
	208: "Параметр time указан неправильно",
	209: "Вы добавили этот номер (или один из номеров) в стоп-лист",
	210: "Используется GET, где необходимо использовать POST",
	211: "Метод не найден",
	212: "Текст сообщения необходимо передать в кодировке UTF-8 (вы передали в другой кодировке)",
	213: "Указано более 100 номеров в списке получателей",
	220: "Сервис временно недоступен, попробуйте чуть позже",
	230: "Превышен общий лимит количества сообщений на этот номер в день",
	231: "Превышен лимит одинаковых сообщений на этот номер в минуту",
	232: "Превышен лимит одинаковых сообщений на этот номер в день",
	300: "Неправильный token (возможно истек срок действия, либо ваш IP изменился)",
	301: "Неправильный api_id, либо логин/пароль",
	302: "Пользователь авторизован, но аккаунт не подтвержден (пользователь не ввел код, присланный в регистрационной смс)",
	303: "Код подтверждения неверен",
	304: "Отправлено слишком много кодов подтверждения. Пожалуйста, повторите запрос позднее",
	305: "Слишком много неверных вводов кода, повторите попытку позднее",
	500: "Ошибка на сервере. Повторите запрос.",
	901: "Callback: URL неверный (не начинается на http://)",
	902: "Callback: Обработчик не найден (возможно был удален ранее)",
}

func GetErrorByCode(error_id int) string {
	return codeStatus[error_id]
}

// Создаем API клиента
func CreateClient(id string) *SmsClient {
	return CreateHTTPClient(id, &http.Client{})
}

// Создаем HTTP клиента для API
func CreateHTTPClient(id string, client *http.Client) *SmsClient {
	c := &SmsClient{
		ApiId: id,
		Http:  client,
	}
	return c
}

// Создаем сообщение
func CreateSMS(to string, text string) *Sms {
	return &Sms{
		To:   to,
		Text: text,
	}
}

// Создаем множество сообщений из ранее созданных, с помощью CreateSMS
func CreateMultipleSMS(sms ...*Sms) *Sms {
	arr := make(map[string]string)
	for _, o := range sms {
		arr[o.To] = o.Text
	}

	return &Sms{
		Multi: arr,
	}
}

// Сборщик запроса
func (c *SmsClient) makeRequest(endpoint string, params url.Values) ([]byte, error) {
	params.Set("api_id", c.ApiId)
	params.Set("json", "1")
	url := API_URL + endpoint + "?" + params.Encode()

	response, err := c.Http.Get(url)
	if err != nil {
		// return false, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		panic(err)
	}

	return []byte(body), err
}

// Отправка СМС сообщения
func (c *SmsClient) SmsSend(p *Sms) (SendedSms, error) {
	var params = url.Values{}

	if len(p.Multi) > 0 {
		for to, text := range p.Multi {
			key := fmt.Sprintf("multi[%s]", to)
			params.Add(key, text)
		}
	} else {
		params.Set("to", p.To)
		params.Set("text", p.Text)
	}

	if len(p.From) > 0 {
		params.Set("from", p.From)
	}

	if p.PartnerId > 0 {
		val := strconv.Itoa(p.PartnerId)
		params.Set("partner_id", val)
	}

	if p.Test {
		params.Set("test", "1")
	}

	if p.Time.After(time.Now()) {
		val := strconv.FormatInt(p.Time.Unix(), 10)
		params.Set("time", val)
	}

	if p.Translit {
		params.Set("translit", "1")
	}

	// params.Set("test", "1")

	body, err := c.makeRequest("/sms/send", params)
	if err != nil {
		return SendedSms{}, err
	}

	sendedSms := SendedSms{}
	err = json.Unmarshal(body, &sendedSms)
	if err != nil {
		return SendedSms{}, err
	}

	return sendedSms, nil
}

// Проверка статуса сообщения
func (c *SmsClient) SmsStatus(sms_id string) (SmsStatus, error) {
	var params = url.Values{}

	params.Set("sms_id", sms_id)

	body, err := c.makeRequest("/sms/status", params)

	smsstatuslist := SmsStatus{}
	err = json.Unmarshal(body, &smsstatuslist)
	if err != nil {
		return SmsStatus{}, err
	}

	return smsstatuslist, err
}

// Проверка стоимости сообщения перед отправкой
func (c *SmsClient) SmsCost(p *Sms) (Cost, error) {
	var params = url.Values{}
	params.Set("to", p.To)
	params.Set("text", p.Text)
	if p.Translit {
		params.Set("translit", "1")
	}

	body, err := c.makeRequest("/sms/cost", params)
	if err != nil {
		return Cost{}, err
	}

	cost := Cost{}
	err = json.Unmarshal(body, &cost)
	if err != nil {
		return Cost{}, err
	}

	return cost, nil
}

// Остаток на балансе
func (c *SmsClient) MyBalance() (Balance, error) {
	body, err := c.makeRequest("/my/balance", url.Values{})
	if err != nil {
		return Balance{}, err
	}

	balance := Balance{}
	err = json.Unmarshal(body, &balance)
	if err != nil {
		return Balance{}, err
	}

	return balance, nil
}

// Получить информацию о бесплатных сообщениях и их использовании
func (c *SmsClient) MyFree() (Free, error) {
	body, err := c.makeRequest("/my/free", url.Values{})
	if err != nil {
		return Free{}, err
	}

	free := Free{}
	err = json.Unmarshal(body, &free)
	if err != nil {
		return Free{}, err
	}

	return free, err
}

// Проверка информации о дневном лимите
func (c *SmsClient) MyLimit() (Limit, error) {
	body, err := c.makeRequest("/my/limit", url.Values{})

	limit := Limit{}
	err = json.Unmarshal(body, &limit)
	if err != nil {
		return Limit{}, err
	}

	return limit, err
}

// Получение списка одобренных отправителей
func (c *SmsClient) MySenders() (Senders, error) {
	body, err := c.makeRequest("/my/senders", url.Values{})

	senders := Senders{}
	err = json.Unmarshal(body, &senders)
	if err != nil {
		return Senders{}, err
	}

	return senders, err
}

// Добавление номера в стоплист
func (c *SmsClient) StoplistAdd(phone string, text string) (StopList, error) {
	var params = url.Values{}

	params.Set("stoplist_phone", phone)
	params.Set("stoplist_text", text)

	body, err := c.makeRequest("/stoplist/add", params)

	stoplist := StopList{}
	err = json.Unmarshal(body, &stoplist)
	if err != nil {
		return StopList{}, err
	}

	return stoplist, err
}

// Удаление номера из стоплиста
func (c *SmsClient) StoplistDel(phone string) (StopList, error) {
	var params = url.Values{}

	params.Set("stoplist_phone", phone)

	body, err := c.makeRequest("/stoplist/del", params)

	stoplist := StopList{}
	err = json.Unmarshal(body, &stoplist)
	if err != nil {
		return StopList{}, err
	}

	return stoplist, err
}

// Выгрузка всего стоплиста
func (c *SmsClient) StoplistGet() (StopList, error) {
	body, err := c.makeRequest("/stoplist/get", url.Values{})

	stoplist := StopList{}
	err = json.Unmarshal(body, &stoplist)
	if err != nil {
		return StopList{}, err
	}

	return stoplist, err
}

// Добавление callback обработчика
func (c *SmsClient) CallbackAdd(callback_url string) (Callback, error) {
	var params = url.Values{}

	params.Set("url", callback_url)

	body, err := c.makeRequest("/callback/add", params)

	callback := Callback{}
	err = json.Unmarshal(body, &callback)
	if err != nil {
		return Callback{}, err
	}

	return callback, err
}

// Удаление callback обработчика
func (c *SmsClient) CallbackDel(callback_url string) (Callback, error) {
	var params = url.Values{}

	params.Set("url", callback_url)

	body, err := c.makeRequest("/callback/del", params)

	callback := Callback{}
	err = json.Unmarshal(body, &callback)
	if err != nil {
		return Callback{}, err
	}

	return callback, err
}

// Выгрузка всех callback обработчиков
func (c *SmsClient) CallbackGet() (Callback, error) {
	body, err := c.makeRequest("/callback/get", url.Values{})

	callback := Callback{}
	err = json.Unmarshal(body, &callback)
	if err != nil {
		return Callback{}, err
	}

	return callback, err
}
