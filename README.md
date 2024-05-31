# go-waifu

## API spec

GET /requests/status
-> 200 returns status of all transactions

POST /requests
-> 202 returns a transaction ID

GET /requests/{transactionId}/status
-> 404 if missing
-> 200 with body explaining status

GET /requests/{transactionId}/output
-> 404 if missing
-> 200 with image

DELETE /requests/{transactionId} (TODO in the future)
-> 404 if missing
-> 204 if successful
