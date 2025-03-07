package handler

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log/slog"
	"remnawave-tg-shop-bot/internal/config"
	"remnawave-tg-shop-bot/internal/cryptopay"
	"remnawave-tg-shop-bot/internal/database"
	"remnawave-tg-shop-bot/internal/yookasa"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Handler struct {
	customerRepository *database.CustomerRepository
	purchaseRepository *database.PurchaseRepository
	cryptoPayClient    *cryptopay.Client
	yookasaClient      *yookasa.Client
}

func NewHandler(
	customerRepository *database.CustomerRepository,
	purchaseRepository *database.PurchaseRepository,
	cryptoPayClient *cryptopay.Client,
	yookasaClient *yookasa.Client) *Handler {
	return &Handler{
		customerRepository: customerRepository,
		purchaseRepository: purchaseRepository,
		cryptoPayClient:    cryptoPayClient,
		yookasaClient:      yookasaClient,
	}
}

const (
	CallbackBuy     = "buy"
	CallbackSell    = "sell"
	CallbackStart   = "start"
	CallbackCrypto  = "crypto"
	CallbackCard    = "card"
	CallbackConnect = "connect"
)

func (h Handler) StartCommandHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	ctxWithTime, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	existingCustomer, err := h.customerRepository.FindByTelegramId(ctx, update.Message.Chat.ID)
	if err != nil {
		slog.Error("error finding customer by telegram id", err)
	}

	if existingCustomer == nil {
		err := h.customerRepository.Create(ctxWithTime, &database.Customer{
			TelegramID: update.Message.Chat.ID,
		})
		if err != nil {
			slog.Error("error creating customer", err)
			return
		}
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID,
		ParseMode: models.ParseModeMarkdown,
		ReplyMarkup: models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{{Text: "Купить", CallbackData: "buy"}},
				{{Text: "Подключиться", CallbackData: "connect"}},
			},
		},
		Text: fmt.Sprintf(`
👋🏻 *Привет*
Это бот для подключения к *VPN*🛡️

Доступны локации:
%s

*Как подключиться:*
• нажмите кнопку *"Подключиться"*
• следуйте короткой инструкции`, bot.EscapeMarkdown(buildAvailableCountriesLists()),
		),
	},
	)
	if err != nil {
		slog.Error("Error sending /start message", err)
	}
}

func buildAvailableCountriesLists() string {
	var locationsText strings.Builder
	countries := config.Countries()

	keys := make([]string, 0, len(countries))
	for k := range countries {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i, k := range keys {
		if i == len(keys)-1 {
			locationsText.WriteString(fmt.Sprintf("└ %s\n", countries[k]))
		} else {
			locationsText.WriteString(fmt.Sprintf("├ %s\n", countries[k]))
		}
	}

	return locationsText.String()
}

func (h Handler) StartCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	callback := update.CallbackQuery.Message.Message
	_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{ChatID: callback.Chat.ID,
		MessageID: callback.ID,
		ParseMode: models.ParseModeMarkdown,
		ReplyMarkup: models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{{Text: "Купить", CallbackData: "buy"}},
				{{Text: "Подключиться", CallbackData: "connect"}},
			},
		},
		Text: fmt.Sprintf(`
👋🏻 *Привет*
Это бот для подключения к *VPN*🛡️

Доступны локации:
%s

*Как подключиться:*
• нажмите кнопку *"Подключиться"*
• следуйте короткой инструкции`, bot.EscapeMarkdown(buildAvailableCountriesLists()),
		)})
	if err != nil {
		slog.Error("Error sending /start message", err)
	}
}

func (h Handler) BuyCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	callback := update.CallbackQuery.Message.Message
	_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:    callback.Chat.ID,
		MessageID: callback.ID,
		ParseMode: models.ParseModeMarkdown,
		ReplyMarkup: models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "1 месяц", CallbackData: buildSellCallbackData(1)},
					{Text: "3 месяц", CallbackData: buildSellCallbackData(3)},
					{Text: "6 месяц", CallbackData: buildSellCallbackData(6)},
				},
				{
					{Text: "Назад", CallbackData: CallbackStart},
				},
			},
		},
		Text: fmt.Sprintf(`
📅 *1 мес*  %d₽
📅 *3 мес*  %d₽
📅 *6 мес*  %d₽

К оплате принимаются карты российских банков и криптовалюта`,
			calculatePrice(1),
			calculatePrice(2),
			calculatePrice(3)),
	})
	if err != nil {
		slog.Error("Error sending buy message", err)
	}
}

func (h Handler) SellCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	callback := update.CallbackQuery.Message.Message
	callbackQuery := parseCallbackData(update.CallbackQuery.Data)
	month := callbackQuery["month"]
	_, err := b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
		ChatID:    callback.Chat.ID,
		MessageID: callback.ID,
		ReplyMarkup: models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "₿ Криптовалютой", CallbackData: fmt.Sprintf("%s?month=%s", CallbackCrypto, month)},
					{Text: "💳 Картой банка", CallbackData: fmt.Sprintf("%s?month=%s", CallbackCard, month)},
				},
				{
					{Text: "Назад", CallbackData: CallbackBuy},
				},
			},
		},
	})

	if err != nil {
		slog.Error("Error sending sell message", err)
	}
}

func (h Handler) CryptoCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	callback := update.CallbackQuery.Message.Message
	callbackQuery := parseCallbackData(update.CallbackQuery.Data)
	month, err := strconv.Atoi(callbackQuery["month"])
	if err != nil {
		slog.Error("Error getting month from query", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	customer, err := h.customerRepository.FindByTelegramId(ctx, callback.Chat.ID)
	if err != nil {
		slog.Error("Error finding customer", err)
	}
	if customer == nil {
		slog.Error("customer not exist", "chatID", callback.Chat.ID, "error", err)
		return
	}

	price := calculatePrice(month)

	purchaseId, err := h.purchaseRepository.Create(ctx, &database.Purchase{
		InvoiceType: database.InvoiceTypeCrypto,
		Status:      database.PurchaseStatusNew,
		Amount:      float64(price),
		Currency:    "RUB",
		CustomerID:  customer.ID,
		Month:       month,
	})
	if err != nil {
		slog.Error("Error creating purchase", err)
		return
	}

	invoice, err := h.cryptoPayClient.CreateInvoice(&cryptopay.InvoiceRequest{
		CurrencyType:   "fiat",
		Fiat:           "RUB",
		Amount:         fmt.Sprintf("%d", price),
		AcceptedAssets: "USDT",
		Payload:        fmt.Sprintf("customerId=%d&purchaseId=%d", callback.Chat.ID, purchaseId),
		Description:    fmt.Sprintf("Subscription on %d month", month),
		PaidBtnName:    "callback",
		PaidBtnUrl:     config.BotURL(),
	})
	if err != nil {
		slog.Error("Error creating invoice", err)
		return
	}

	updates := map[string]interface{}{
		"crypto_invoice_url": invoice.BotInvoiceUrl,
		"crypto_invoice_id":  invoice.InvoiceID,
		"status":             database.PurchaseStatusPending,
	}

	err = h.purchaseRepository.UpdateFields(ctx, purchaseId, updates)
	if err != nil {
		slog.Error("Error updating purchase", err)
		return
	}

	_, err = b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
		ChatID:    callback.Chat.ID,
		MessageID: callback.ID,
		ReplyMarkup: models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "Оплатить", URL: invoice.BotInvoiceUrl},
					{Text: "Назад", CallbackData: buildSellCallbackData(month)},
				},
			},
		},
	})
	if err != nil {
		slog.Error("Error updating sell message", err)
	}

}

func (h Handler) YookasaCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	callback := update.CallbackQuery.Message.Message
	callbackQuery := parseCallbackData(update.CallbackQuery.Data)
	month, err := strconv.Atoi(callbackQuery["month"])
	if err != nil {
		slog.Error("Error getting month from query", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	customer, err := h.customerRepository.FindByTelegramId(ctx, callback.Chat.ID)
	if err != nil {
		slog.Error("Error finding customer", err)
	}
	if customer == nil {
		slog.Error("customer not exist", "chatID", callback.Chat.ID, "error", err)
		return
	}

	price := calculatePrice(month)
	purchaseId, err := h.purchaseRepository.Create(ctx, &database.Purchase{
		InvoiceType: database.InvoiceTypeYookasa,
		Status:      database.PurchaseStatusNew,
		Amount:      float64(price),
		Currency:    "RUB",
		CustomerID:  customer.ID,
		Month:       month,
	})
	if err != nil {
		slog.Error("Error creating purchase", err)
		return
	}

	invoice, err := h.yookasaClient.CreateInvoice(ctx, price, month, customer.ID, purchaseId)
	if err != nil {
		slog.Error("Error creating invoice", err)
		return
	}

	updates := map[string]interface{}{
		"yookasa_url": invoice.Confirmation.ConfirmationURL,
		"yookasa_id":  invoice.ID,
		"status":      database.PurchaseStatusPending,
	}

	err = h.purchaseRepository.UpdateFields(ctx, purchaseId, updates)
	if err != nil {
		slog.Error("Error updating purchase", err)
		return
	}

	_, err = b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
		ChatID:    callback.Chat.ID,
		MessageID: callback.ID,
		ReplyMarkup: models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "Оплатить", URL: invoice.Confirmation.ConfirmationURL},
					{Text: "Назад", CallbackData: buildSellCallbackData(month)},
				},
			},
		},
	})
	if err != nil {
		slog.Error("Error updating sell message", err)
	}
}

func (h Handler) ConnectCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	callback := update.CallbackQuery.Message.Message

	customer, err := h.customerRepository.FindByTelegramId(ctx, callback.Chat.ID)
	if err != nil {
		slog.Error("Error finding customer", err)
	}
	if customer == nil {
		slog.Error("customer not exist", "chatID", callback.Chat.ID, "error", err)
	}

	_, err = b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:    callback.Chat.ID,
		MessageID: callback.ID,
		Text:      buildConnectText(customer),
		ReplyMarkup: models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{{Text: "Назад", CallbackData: CallbackStart}},
			},
		},
	})

	if err != nil {
		slog.Error("Error sending connect message", err)
	}
}

func buildConnectText(customer *database.Customer) string {
	var info strings.Builder

	if customer.ExpireAt != nil {
		currentTime := time.Now()

		if currentTime.Before(*customer.ExpireAt) {
			formattedDate := customer.ExpireAt.Format("02.01.2006 15:04")
			info.WriteString(fmt.Sprintf("Ваша подписка действует до: %s", formattedDate))

			if customer.SubscriptionLink != nil && *customer.SubscriptionLink != "" {
				info.WriteString(fmt.Sprintf("\n\nСсылка на подписку: %s", *customer.SubscriptionLink))
			}
		} else {
			info.WriteString("У вас нет активной подписки")
		}
	} else {
		info.WriteString("У вас нет активной подписки")
	}

	return info.String()
}

func parseCallbackData(data string) map[string]string {
	result := make(map[string]string)

	parts := strings.Split(data, "?")
	if len(parts) < 2 {
		return result
	}

	params := strings.Split(parts[1], "&")
	for _, param := range params {
		kv := strings.SplitN(param, "=", 2)
		if len(kv) == 2 {
			result[kv[0]] = kv[1]
		}
	}

	return result
}

func buildSellCallbackData(month int) string {
	return fmt.Sprintf("%s?month=%d", CallbackSell, month)
}

func calculatePrice(month int) int {
	return config.Price() * month
}
