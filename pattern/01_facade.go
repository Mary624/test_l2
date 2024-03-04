package pattern

import (
	"encoding/json"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit"
)

// Простой интерфейс для более сложной системы
// Плюсы: изолирует клиентов от сложной подсистемы
// Минусы: рискует взять на себя слишком много работы

func ExampleFacade() {
	sender := NewRandomOrderSender()
	sender.SendOrder(os.Stdout)
}

type RandomOrderSender struct {
}

func NewRandomOrderSender() *RandomOrderSender {
	return &RandomOrderSender{}
}

type order struct {
	OrderUid          string    `json:"order_uid"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Delivery          delivery  `json:"delivery"`
	Payment           payment   `json:"payment"`
	Items             []Item    `json:"items"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerId        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	Shardkey          string    `json:"shardkey"`
	SmId              int64     `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OofShard          string    `json:"oof_shard"`
}

type delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type payment struct {
	Transaction  string `json:"transaction"`
	RequestId    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int64  `json:"amount"`
	PaymentDt    int64  `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int64  `json:"delivery_cost"`
	GoodsTotal   int64  `json:"goods_total"`
	CustomFee    int64  `json:"custom_fee"`
}

type Item struct {
	ChrtId      int64  `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int64  `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int64  `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int64  `json:"total_price"`
	NmId        int64  `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int64  `json:"status"`
}

func (sender *RandomOrderSender) SendOrder(writer io.Writer) error {
	data, err := json.Marshal(sender.randomNormal())
	if err != nil {
		return err
	}
	_, err = writer.Write(data)
	return err
}

func (sender *RandomOrderSender) randomNormal() order {
	gofakeit.Seed(time.Now().UnixMilli())
	return order{
		OrderUid:    gofakeit.UUID(),
		TrackNumber: sender.randomStr(),
		Entry:       sender.randomStr(),
		Delivery: delivery{
			Name:    gofakeit.Name(),
			Phone:   gofakeit.Phone(),
			Zip:     strconv.Itoa(int(gofakeit.Uint32())),
			City:    gofakeit.City(),
			Address: gofakeit.Address().Address,
			Region:  sender.randomStr(),
			Email:   gofakeit.Email(),
		},
		Payment: payment{
			Transaction:  gofakeit.UUID(),
			RequestId:    sender.randomStr(),
			Currency:     sender.randomStr(),
			Provider:     sender.randomStr(),
			Amount:       int64(gofakeit.Uint32()),
			PaymentDt:    int64(gofakeit.Uint32()),
			Bank:         sender.randomStr(),
			DeliveryCost: int64(gofakeit.Uint32()),
			GoodsTotal:   int64(gofakeit.Uint32()),
			CustomFee:    int64(gofakeit.Uint32()),
		},
		Items:             sender.randomNormalItems(),
		Locale:            sender.randomStr(),
		InternalSignature: sender.randomStr(),
		CustomerId:        sender.randomStr(),
		DeliveryService:   sender.randomStr(),
		Shardkey:          strconv.Itoa(int(gofakeit.Uint32())),
		SmId:              int64(gofakeit.Uint32()),
		DateCreated:       time.Now(),
		OofShard:          strconv.Itoa(int(gofakeit.Uint32())),
	}
}

func (sender *RandomOrderSender) randomNormalItems() []Item {
	min, max := 1, 20
	c := rand.Intn(max-min) + min
	items := make([]Item, 0, c)
	for i := 0; i < c; i++ {
		items = append(items, Item{
			ChrtId:      int64(gofakeit.Uint32()),
			TrackNumber: sender.randomStr(),
			Price:       int64(gofakeit.Uint32()),
			Rid:         gofakeit.UUID(),
			Name:        gofakeit.Word(),
			Sale:        int64(gofakeit.Uint32()),
			Size:        sender.randomStr(),
			TotalPrice:  int64(gofakeit.Uint32()),
			NmId:        int64(gofakeit.Uint32()),
			Brand:       sender.randomStr(),
			Status:      int64(gofakeit.Uint32()),
		})
	}
	return items
}

func (sender *RandomOrderSender) randomStr() string {
	var b strings.Builder
	min, max := 1, 20
	c := rand.Intn(max-min) + min
	chs := "abcdefghijklmnopqrstuvwxyz"
	for i := 0; i < c; i++ {
		ch := ([]rune)(chs)[rand.Intn(len(chs)-0)+0]
		b.WriteRune(ch)
	}
	return b.String()
}
