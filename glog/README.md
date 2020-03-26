# glog

simple, goroutine safe logger

## log level

1. info
2. warning
3. error
4. debug

## glog.Config

```
{
    "type": <"console" | "file" | "tcp" | "udp">,
    "format": <"json" | "text">,
    "datetime_format": <"rfc3339">,
    "log_level": <"info" | "warning" | "error" | "debug">, // default: debug
    "file": {                                              // if using file type
        "path": <"file path">,
        "file_name": <"file name">,
        "rotation": <"every_hour" | "every_day">
    },
    "tcp": {                                               // if using tcp type
        "host":"127.0.0.1:1234"
    },
    "udp": {                                               // if using udp type
        "host":"127.0.0.1:1234"
    },
    "sub_logger":[global.Config...]
}
```