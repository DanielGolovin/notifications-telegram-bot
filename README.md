# notifications-telegram-bot

This is a simple Telegram bot that broadcasts data in a form of stringified json to all subscribers, when receives a POST request to /notify endpoint.

example of POST request:

```json
{
  "data": {
    "message": "Hello World!"
  },
  "secret": "my_secret"
}
```

Secret is a string that you can set in environment variable `NOTIFICATION_SECRET`
