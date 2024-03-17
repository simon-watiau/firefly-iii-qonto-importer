# firefly-iii-qonto-importer

Imports transactions and documents from Qonto into [Firefly III](https://www.firefly-iii.org/)

## Workflow

- For each Qonto transactions
- Check if a transaction exist in Firefly using `external_id_is:TRANSACTION_ID` or `notes_contains:"[TRANSACTION_ID]"`
  - If not, create a new transaction using the provided revenue, expense and asset accounts and sets its external_id
- Check if the Qonto documents exist in Firefly (using the file name) on the found transaction(s)
  - If not, upload the missing documents on all matching transactions


## Configuration

| Command line flag | Docker environment variable | Example | Description | 
| -- | -- | -- | -- |
| `--firefly-base-url` | `FIREFLY_BASE_URL`| `http://localhost:8080` | URL of your Firefly III instance |
| `--firefly-token` | `FIREFLY_TOKEN` | | Firefly III [Personal access token](https://docs.firefly-iii.org/how-to/firefly-iii/features/api/#personal-access-tokens) |
|` --firefly-asset-account-id` | `FIREFLY_ASSET_ACCOUNT_ID` | `1` | Id of the asset account used for revenue and withdawal |
|  `--firefly-revenue-account-id` | `FIREFLY_REVENUE_ACCOUNT_ID` | `2` | Id of the revenue account for revenue |
| `--firefly-expense-account-id` | `FIREFLY_EXPENSE_ACCOUNT_ID` | `3` | Id of the expense account for withdrawal |
| `--qonto-login` | `QONTO_LOGIN` | `ABCD-4600` | [Qonto API](https://api-doc.qonto.com/docs/business-api/ZG9jOjQ2NDA2-introduction) login |
| `--qonto-password` | `QONTO_PASSWORD` | --- | [Qonto API](https://api-doc.qonto.com/docs/business-api/ZG9jOjQ2NDA2-introduction) password |
| `--qonto-iban` | `QONTO_IBAN` | --- | Qonto account IBAN to fetch transactions from | 
## Usage

### Manual

Use `GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o importer-cli .` to build the importer as `importer-cli`

```
./importer-cli import \
  --firefly-base-url $FIREFLY_BASE_URL \
  --firefly-token $FIREFLY_TOKEN \
  --firefly-asset-account-id $FIREFLY_ASSET_ACCOUNT_ID \ 
  --firefly-revenue-account-id $FIREFLY_REVENUE_ACCOUNT_ID \
  --firefly-expense-account-id $FIREFLY_EXPENSE_ACCOUNT_ID \
  --qonto-login $QONTO_LOGIN \
  --qonto-password $QONTO_PASSWORD \
  --qonto-iban $QONTO_IBAN
```

### Docker 

The `main` branch is build as `simonwatiau/firefly-iii-qonto-importer:latest`

Using `docker run` :
##### 
```
docker run \
  -e FIREFLY_BASE_URL=FIREFLY_BASE_URL \
  -e FIREFLY_TOKEN=FIREFLY_TOKEN \
  -e FIREFLY_ASSET_ACCOUNT_ID=FIREFLY_ASSET_ACCOUNT_ID \
  -e FIREFLY_REVENUE_ACCOUNT_ID=FIREFLY_REVENUE_ACCOUNT_ID \
  -e FIREFLY_EXPENSE_ACCOUNT_ID=FIREFLY_EXPENSE_ACCOUNT_ID \
  -e QONTO_LOGIN=QONTO_LOGIN \
  -e QONTO_PASSWORD=QONTO_PASSWORD \
  -e QONTTO_IBAN=QONTTO_IBAN \
  --rm simonwatiau/firefly-iii-qonto-importer:latest 
```

Using `compose`/`swarm`:
```
  qonto-import:
    image: simonwatiau/firefly-iii-qonto-importer:latest
    environment:
      FIREFLY_BASE_URL: "---"
      FIREFLY_TOKEN: "---"
      FIREFLY_ASSET_ACCOUNT_ID: "---"
      FIREFLY_REVENUE_ACCOUNT_ID: "---"
      FIREFLY_EXPENSE_ACCOUNT_ID: "---"
      QONTO_LOGIN: "---"
      QONTO_PASSWORD: "---"
      QONTTO_IBAN: "---"
```