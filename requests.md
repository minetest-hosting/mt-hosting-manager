
# zahlsch webhook callback

```sh
export BASEURL=http://localhost:8080
export ZAHLSCH_WEBHOOK_KEY=
export json='{"transaction":{"id":"","amount":0,"pageUuid":"","status":"confirmed","invoice":{"custom_fields":[{"name":"custom_user_id","value":""},{"name":"custom_transaction_id","value":""}]}}}'

curl --data "${json}" -H "Content-Type: application/json" ${BASEURL}/api/webhook/zahlsch?key=${ZAHLSCH_WEBHOOK_KEY}
```