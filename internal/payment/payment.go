package payment

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log/slog"
	"remnawave-tg-shop-bot/internal/config"
	"remnawave-tg-shop-bot/internal/cryptopay"
	"remnawave-tg-shop-bot/internal/database"
	"remnawave-tg-shop-bot/internal/remnawave"
	"remnawave-tg-shop-bot/internal/translation"
	"remnawave-tg-shop-bot/internal/yookasa"
	"time"
)

type PaymentService struct {
	purchaseRepository *database.PurchaseRepository
	remnawaveClient    *remnawave.Client
	customerRepository *database.CustomerRepository
	telegramBot        *bot.Bot
	translation        *translation.Manager
	cryptoPayClient    *cryptopay.Client
	yookasaClient      *yookasa.Client
}

func NewPaymentService(
	translation *translation.Manager,
	purchaseRepository *database.PurchaseRepository,
	remnawaveClient *remnawave.Client,
	customerRepository *database.CustomerRepository,
	telegramBot *bot.Bot,
	cryptoPayClient *cryptopay.Client,
	yookasaClient *yookasa.Client,
) *PaymentService {
	return &PaymentService{
		purchaseRepository: purchaseRepository,
		remnawaveClient:    remnawaveClient,
		customerRepository: customerRepository,
		telegramBot:        telegramBot,
		translation:        translation,
		cryptoPayClient:    cryptoPayClient,
		yookasaClient:      yookasaClient,
	}
}

func (s PaymentService) ProcessPurchaseById(purchaseId int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()
	purchase, err := s.purchaseRepository.FindById(ctx, purchaseId)
	if err != nil {
		return err
	}
	if purchase == nil {
		return fmt.Errorf("purchase with crypto invoice id %d not found", purchaseId)
	}

	customer, err := s.customerRepository.FindById(ctx, purchase.CustomerID)
	if err != nil {
		return err
	}
	if customer == nil {
		return fmt.Errorf("customer %s not found", purchase.CustomerID)
	}

	username := fmt.Sprintf("%d_%d", customer.ID, customer.TelegramID)

	user, err := s.remnawaveClient.CreateOrUpdateUser(ctx, username, purchase.Month)
	if err != nil {
		return err
	}

	err = s.purchaseRepository.MarkAsPaid(ctx, purchase.ID)
	if err != nil {
		return err
	}

	customerFilesToUpdate := map[string]interface{}{
		"subscription_link": user.SubscriptionURL,
		"expire_at":         user.ExpireAt,
	}

	err = s.customerRepository.UpdateFields(ctx, customer.ID, customerFilesToUpdate)
	if err != nil {
		return err
	}

	_, err = s.telegramBot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: customer.TelegramID,
		Text:   s.translation.GetText(customer.Language, "subscription_activated"),
		ReplyMarkup: models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: s.translation.GetText(customer.Language, "connect_button"), URL: user.SubscriptionURL},
				},
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (s PaymentService) CreatePurchase(ctx context.Context, amount int, months int, customerID int64, invoiceType database.InvoiceType) (string, error) {
	switch invoiceType {
	case database.InvoiceTypeCrypto:
		return s.createCryptoInvoice(ctx, amount, months, customerID)
	case database.InvoiceTypeYookasa:
		return s.createYookasaInvoice(ctx, amount, months, customerID)
	default:
		return "", fmt.Errorf("unknown invoice type: %s", invoiceType)
	}
}

func (s PaymentService) createCryptoInvoice(ctx context.Context, amount int, months int, customerID int64) (string, error) {
	purchaseId, err := s.purchaseRepository.Create(ctx, &database.Purchase{
		InvoiceType: database.InvoiceTypeCrypto,
		Status:      database.PurchaseStatusNew,
		Amount:      float64(amount),
		Currency:    "RUB",
		CustomerID:  customerID,
		Month:       months,
	})
	if err != nil {
		slog.Error("Error creating purchase", err)
		return "", err
	}

	invoice, err := s.cryptoPayClient.CreateInvoice(&cryptopay.InvoiceRequest{
		CurrencyType:   "fiat",
		Fiat:           "RUB",
		Amount:         fmt.Sprintf("%d", amount),
		AcceptedAssets: "USDT",
		Payload:        fmt.Sprintf("purchaseId=%d", purchaseId),
		Description:    fmt.Sprintf("Subscription on %d month", months),
		PaidBtnName:    "callback",
		PaidBtnUrl:     config.BotURL(),
	})
	if err != nil {
		slog.Error("Error creating invoice", err)
		return "", err
	}

	updates := map[string]interface{}{
		"crypto_invoice_url": invoice.BotInvoiceUrl,
		"crypto_invoice_id":  invoice.InvoiceID,
		"status":             database.PurchaseStatusPending,
	}

	err = s.purchaseRepository.UpdateFields(ctx, purchaseId, updates)
	if err != nil {
		slog.Error("Error updating purchase", err)
		return "", err
	}

	return invoice.BotInvoiceUrl, nil
}

func (s PaymentService) createYookasaInvoice(ctx context.Context, amount int, months int, customerId int64) (string, error) {
	purchaseId, err := s.purchaseRepository.Create(ctx, &database.Purchase{
		InvoiceType: database.InvoiceTypeYookasa,
		Status:      database.PurchaseStatusNew,
		Amount:      float64(amount),
		Currency:    "RUB",
		CustomerID:  customerId,
		Month:       months,
	})
	if err != nil {
		slog.Error("Error creating purchase", err)
		return "", err
	}

	invoice, err := s.yookasaClient.CreateInvoice(ctx, amount, months, customerId, purchaseId)
	if err != nil {
		slog.Error("Error creating invoice", err)
		return "", err
	}

	updates := map[string]interface{}{
		"yookasa_url": invoice.Confirmation.ConfirmationURL,
		"yookasa_id":  invoice.ID,
		"status":      database.PurchaseStatusPending,
	}

	err = s.purchaseRepository.UpdateFields(ctx, purchaseId, updates)
	if err != nil {
		slog.Error("Error updating purchase", err)
		return "", err
	}

	return invoice.Confirmation.ConfirmationURL, nil
}
