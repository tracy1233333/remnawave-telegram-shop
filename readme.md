# Remnawave Telegram Shop

## Description

A Telegram bot for selling subscriptions with integration to Remnawave (https://remna.st/). This service allows users to purchase and manage subscriptions through Telegram with multiple payment system options.

## Environment Variables

The application requires the following environment variables to be set:

| Variable | Description |
|----------|-------------|
| `PRICE` | The base price for the service (default: 1) |
| `TELEGRAM_TOKEN` | Telegram Bot API token for bot functionality |
| `DATABASE_URL` | PostgreSQL connection string |
| `POSTGRES_USER` | PostgreSQL username |
| `POSTGRES_PASSWORD` | PostgreSQL password |
| `POSTGRES_DB` | PostgreSQL database name |
| `REMNAWAVE_URL` | Remnawave API URL |
| `REMNAWAVE_TOKEN` | Authentication token for Remnawave API |
| `CRYPTO_PAY_TOKEN` | CryptoPay API token |
| `CRYPTO_PAY_URL` | CryptoPay API URL |
| `YOOKASA_SECRET_KEY` | YooKassa API secret key |
| `YOOKASA_SHOP_ID` | YooKassa shop identifier |
| `YOOKASA_URL` | YooKassa API URL |
| `YOOKASA_EMAIL` | Email address associated with YooKassa account |

## Plugins and Dependencies

### Telegram Bot
- [Telegram Bot API](https://core.telegram.org/bots/api)
- [Go Telegram Bot API](https://github.com/go-telegram-bot-api/telegram-bot-api)

### Database
- [PostgreSQL](https://www.postgresql.org/)
- [pgx - PostgreSQL Driver](https://github.com/jackc/pgx)

### Payment Systems
- [YooKassa API](https://yookassa.ru/developers/api)
- [CryptoPay API](https://help.crypt.bot/crypto-pay-api)


## Setup Instructions

1. Clone the repository
2. Create a `.env` file in the root directory with all the environment variables listed above
3. Run the application:

```bash
docker-compose up -d
```
