FROM golang:1.22-alpine AS builder

COPY . .

RUN go get -d -v

RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /importer

FROM alpine

COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /importer /importer

CMD ["sh", "-c",  "/importer import --firefly-base-url $FIREFLY_BASE_URL --firefly-token $FIREFLY_TOKEN --firefly-asset-account-id $FIREFLY_ASSET_ACCOUNT_ID --firefly-revenue-account-id $FIREFLY_REVENUE_ACCOUNT_ID --firefly-expense-account-id $FIREFLY_EXPENSE_ACCOUNT_ID --qonto-login $QONTO_LOGIN --qonto-password $QONTO_PASSWORD --qonto-iban $QONTTO_IBAN"]
