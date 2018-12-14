# airview

## Server

* Start server by `go run bin/main.go` .
* Access to http://127.0.0.1:8000/ .

## FlashAir

### Setup WiFi Connection

Please see [Documentation of FlashAir](https://www.flashair-developers.com/ja/documents/api/config/) .

### Setup Application

* Edit `PHOTO_TARGET_PATH` and `NOTIFY_URL` in `notify.lua` .
* Copy `lua/notify.lua` to FlashAir `/` .
* Add these lines in `/SD_WLAN/CONFIG` .

```
LUA_SD_EVENT=/notify.lua
```
