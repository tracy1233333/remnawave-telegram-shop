# Remnawave Telegram Shop

## Description

A Telegram bot for selling subscriptions with integration to Remnawave (https://remna.st/). This service allows users to
purchase and manage subscriptions through Telegram with multiple payment system options.

## Admin commands

- `/sync` - Poll users from remnawave and synchronize them with the database

## Features

- Purchase VPN subscriptions with different payment methods (bank cards, cryptocurrency)
- Multiple subscription plans (1, 3, and 6 months)
- Automated subscription management
- **Subscription Notifications**: The bot automatically sends notifications to users 3 days before their subscription
  expires, helping them avoid service interruption
- Multi-language support (Russian and English)

## Environment Variables

The application requires the following environment variables to be set:

| Variable                 | Description                                                                                                          |
| ------------------------ | -------------------------------------------------------------------------------------------------------------------- | 
| `PRICE`                  | The base price for the service (default: 1)                                                                          |
| `TELEGRAM_TOKEN`         | Telegram Bot API token for bot functionality                                                                         |
| `DATABASE_URL`           | PostgreSQL connection string                                                                                         |
| `POSTGRES_USER`          | PostgreSQL username                                                                                                  |
| `POSTGRES_PASSWORD`      | PostgreSQL password                                                                                                  |
| `POSTGRES_DB`            | PostgreSQL database name                                                                                             |
| `REMNAWAVE_URL`          | Remnawave API URL                                                                                                    |
| `REMNAWAVE_MODE`         | Remnawave mode (remote/local), default is remote. If local set – you can pass http://remnawave:3000 to REMNAWAVE_URL |
| `REMNAWAVE_TOKEN`        | Authentication token for Remnawave API                                                                               |
| `CRYPTO_PAY_ENABLED`     | Enable/disable CryptoPay payment method (true/false)                                                                 |
| `CRYPTO_PAY_TOKEN`       | CryptoPay API token                                                                                                  |
| `CRYPTO_PAY_URL`         | CryptoPay API URL                                                                                                    |
| `YOOKASA_ENABLED`        | Enable/disable YooKassa payment method (true/false)                                                                  |
| `YOOKASA_SECRET_KEY`     | YooKassa API secret key                                                                                              |
| `YOOKASA_SHOP_ID`        | YooKassa shop identifier                                                                                             |
| `YOOKASA_URL`            | YooKassa API URL                                                                                                     |
| `YOOKASA_EMAIL`          | Email address associated with YooKassa account                                                                       |
| `TRAFFIC_LIMIT`          | Maximum allowed traffic in gb (0 to set unlimited)                                                                   |
| `TELEGRAM_STARS_ENABLED` | Enable/disable Telegram Stars payment method (true/false)                                                            |
| `SERVER_STATUS_URL`      | URL to server status page (optional) - if not set, button will not be displayed                                      |
| `SUPPORT_URL`            | URL to support chat or page (optional) - if not set, button will not be displayed                                    |
| `FEEDBACK_URL`           | URL to feedback/reviews page (optional) - if not set, button will not be displayed                                   |
| `CHANNEL_URL`            | URL to Telegram channel (optional) - if not set, button will not be displayed                                        |
| `ADMIN_TELEGRAM_ID`      | Admin telegram id                                                                                                    |
| `TRIAL_TRAFFIC_LIMIT`    | Maximum allowed traffic in gb for trial subscriptions                                                                |     
| `TRIAL_DAYS`             | Number of days for trial subscriptions                                                                               |

## User Interface

The bot dynamically creates buttons based on available environment variables:

- Main buttons for purchasing and connecting to the VPN are always shown
- Additional buttons for Server Status, Support, Feedback, and Channel are only displayed if their corresponding URL
  environment variables are set

## Automated Notifications

The bot includes a notification system that runs daily at 16:00 UTC to check for expiring subscriptions:

- Users receive a notification 3 days before their subscription expires
- The notification includes the exact expiration date and a convenient button to renew the subscription
- Notifications are sent in the user's preferred language

## Plugins and Dependencies

### Telegram Bot

- [Telegram Bot API](https://core.telegram.org/bots/api)
- [Go Telegram Bot API](https://github.com/go-telegram/bot)

### Database

- [PostgreSQL](https://www.postgresql.org/)
- [pgx - PostgreSQL Driver](https://github.com/jackc/pgx)

### Payment Systems

- [YooKassa API](https://yookassa.ru/developers/api)
- [CryptoPay API](https://help.crypt.bot/crypto-pay-api)
- Telegram Stars

## Setup Instructions

1. Clone the repository
2. Create a `.env` file in the root directory with all the environment variables listed above
3. Run the application:

   ```bash
   docker compose up -d
   ```

## Update Instructions

1. Pull the latest Docker image:

   ```bash
   docker compose pull
   ```

2. Restart the containers:
   ```bash
   docker compose down && docker compose up -d
   ```

## Donations

If you appreciate this project and want to help keep it running (and fuel those caffeine-fueled coding marathons),
consider donating. Your support helps drive future updates and improvements.

**Donation Methods:**

- **Ethereum:** `0xd6d35119f8EE2a54Df344E4812A47e1C348ADE1c`
