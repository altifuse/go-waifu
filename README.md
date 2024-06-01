# go-waifu
A simple, containerized web UI for [nihui/waifu2x-ncnn-vulkan](https://github.com/nihui/waifu2x-ncnn-vulkan).

## API spec

POST /requests
-> 202 returns a transaction ID
-> 503 if queue is full

GET /requests/status
-> 200 returns status of all transactions

GET /requests/{transactionId}/status
-> 404 if missing
-> 200 with body explaining status

GET /requests/{transactionId}/output
-> 404 if missing
-> 200 with image

DELETE /requests/{transactionId} (TODO in the future)
-> 404 if missing
-> 204 if successful
