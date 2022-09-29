# uptime
Simple uptime monitoring for http services


### Run monitoring with the config path
```bash
./uptime -config config.json
```

### Example config.json
```json
{
    "services": [
        {
            "name": "Example website",
            "url": "http://example.com",
            "timeout": 1,
            "check_interval": 5,
            "notifications": [
                {
                    "callback_type": "TELEGRAM",
                    "callback_chat_id": "1232131",
                    "callback": "https://api.telegram.org/botTELEGRAM_TOKEN_HERE/sendMessage"
                }
            ]
        }
    ]
}
```
