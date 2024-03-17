package cmd

import (
	"fmt"
	"net/http"

	"github.com/simon-watiau/firefly-iii-qonto-importer/importer"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var fireflyBaseUrl string
var fireflyToken string
var fireflyAssetAccountId string
var fireflyRevenueAccountId string
var fireflyExpenseAccountId string
var qontoBaseUrl string
var qontoLogin string
var qontoPassword string
var qontoIban string

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import Qonto transactions and documents into Firefly III",
	RunE: func(cmd *cobra.Command, args []string) error {
		config := zap.NewDevelopmentConfig()

		config.EncoderConfig = zapcore.EncoderConfig{
			TimeKey:       "",
			LevelKey:      "",
			NameKey:       "",
			CallerKey:     "",
			MessageKey:    "M",
			StacktraceKey: "",
		}

		logger, err := config.Build()

		if err != nil {
			return fmt.Errorf("failed to create logger: %w", err)
		}

		importerClient, err := importer.NewImporter(
			logger,
			&http.Client{},
			importer.ImporterConfig{
				Firefly: importer.FireflyConfig{
					BaseUrl:          fireflyBaseUrl + "/api",
					Token:            fireflyToken,
					AssetAccountId:   fireflyAssetAccountId,
					RevenueAccountId: fireflyRevenueAccountId,
					ExpenseAccountId: fireflyExpenseAccountId,
				},
				Qonto: importer.QontoConfig{
					BaseUrl:  qontoBaseUrl,
					Login:    qontoLogin,
					Password: qontoPassword,
					Iban:     qontoIban,
				},
			},
		)

		if err != nil {
			return fmt.Errorf("failed to create importer: %w", err)
		}

		return importerClient.Import()
	},
}

func init() {
	importCmd.PersistentFlags().StringVar(&fireflyBaseUrl, "firefly-base-url", "", "Firefly base url (example: http://localhost:8080)")
	importCmd.PersistentFlags().StringVar(&fireflyToken, "firefly-token", "", "Firefly token")
	importCmd.PersistentFlags().StringVar(&fireflyAssetAccountId, "firefly-asset-account-id", "", "Firefly asset account id")
	importCmd.PersistentFlags().StringVar(&fireflyExpenseAccountId, "firefly-expense-account-id", "", "Firefly expense account id")
	importCmd.PersistentFlags().StringVar(&fireflyRevenueAccountId, "firefly-revenue-account-id", "", "Firefly revenue account id")

	importCmd.PersistentFlags().StringVar(&qontoBaseUrl, "qonto-base-url", "https://thirdparty.qonto.com", "Qonto base url")
	importCmd.PersistentFlags().StringVar(&qontoLogin, "qonto-login", "", "Qonto API login")
	importCmd.PersistentFlags().StringVar(&qontoPassword, "qonto-password", "", "Qonto API password")
	importCmd.PersistentFlags().StringVar(&qontoIban, "qonto-iban", "", "Qonto iban")

	rootCmd.AddCommand(importCmd)

}
