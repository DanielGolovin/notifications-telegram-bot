# notifications-telegram-bot

This is a simple Telegram bot that broadcasts data in a form of stringified json to all subscribers, when receives a POST request to /notify endpoint.

example of POST request:

```json
{
  "data": {
    "message": "Hello World!"
  }
}
```

X-Secret header is required to authenticate the request.

Secret is a string that you can set in environment variable `NOTIFICATION_SECRET`

TODO: rest of the README
