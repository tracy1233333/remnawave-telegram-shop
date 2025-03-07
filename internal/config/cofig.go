package config

import (
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"strconv"
)

type config struct {
	telegramToken    string
	price            int
	remnawaveUrl     string
	remnawaveToken   string
	databaseURL      string
	cryptoPayURL     string
	cryptoPayToken   string
	botURL           string
	yookasaURL       string
	yookasaShopId    string
	yookasaSecretKey string
	yookasaEmail     string
	countries        map[string]string
}

var conf config

func YookasaEmail() string {
	return conf.yookasaEmail
}
func Price() int {
	return conf.price
}
func TelegramToken() string {
	return conf.telegramToken
}
func RemnawaveUrl() string {
	return conf.remnawaveUrl
}
func DadaBaseUrl() string {
	return conf.databaseURL
}
func RemnawaveToken() string {
	return conf.remnawaveToken

}
func CryptoPayUrl() string {
	return conf.cryptoPayURL
}
func CryptoPayToken() string {
	return conf.cryptoPayToken
}
func BotURL() string {
	return conf.botURL
}
func SetBotURL(botURL string) {
	conf.botURL = botURL
}
func YookasaUrl() string {
	return conf.yookasaURL
}
func YookasaShopId() string {
	return conf.yookasaShopId
}
func YookasaSecretKey() string {
	return conf.yookasaSecretKey
}
func SetCountries(countries map[string]string) {
	conf.countries = countries
}
func Countries() map[string]string {
	return conf.countries
}

func InitConfig() {
	err := godotenv.Load(".env")
	if err != nil {
		slog.Warn("Env file not found")
	}

	conf.telegramToken = os.Getenv("TELEGRAM_TOKEN")
	if conf.telegramToken == "" {
		panic("TELEGRAM_TOKEN .env variable not set")
	}

	strPrice := os.Getenv("PRICE")
	if strPrice == "" {
		panic("PRICE .env variable not set")
	}
	price, err := strconv.Atoi(strPrice)
	if err != nil {
		panic(err)
	}
	conf.price = price

	conf.remnawaveUrl = os.Getenv("REMNAWAVE_URL")
	if conf.remnawaveUrl == "" {
		panic("REMNAWAVE_URL .env variable not set")
	}

	conf.remnawaveToken = os.Getenv("REMNAWAVE_TOKEN")
	if conf.remnawaveToken == "" {
		panic("REMNAWAVE_TOKEN .env variable not set")
	}

	conf.databaseURL = os.Getenv("DATABASE_URL")
	if conf.databaseURL == "" {
		panic("DADA_BASE_URL .env variable not set")
	}

	conf.cryptoPayURL = os.Getenv("CRYPTO_PAY_URL")
	if conf.cryptoPayURL == "" {
		panic("CRYPTO_PAY_URL .env variable not set")
	}
	conf.cryptoPayToken = os.Getenv("CRYPTO_PAY_TOKEN")
	if conf.cryptoPayToken == "" {
		panic("CRYPTO_PAY_TOKEN .env variable not set")
	}

	conf.yookasaURL = os.Getenv("YOOKASA_URL")
	if conf.yookasaURL == "" {
		panic("YOOKASA_URL .env variable not set")
	}

	conf.yookasaShopId = os.Getenv("YOOKASA_SHOP_ID")
	if conf.yookasaShopId == "" {
		panic("YOOKASA_SHOP_ID .env variable not set")
	}

	conf.yookasaSecretKey = os.Getenv("YOOKASA_SECRET_KEY")
	if conf.yookasaSecretKey == "" {
		panic("YOOKASA_SECRET_KEY .env variable not set")
	}

	conf.yookasaEmail = os.Getenv("YOOKASA_EMAIL")
	if conf.yookasaEmail == "" {
		panic("YOOKASA_EMAIL .env variable not set")
	}
}
