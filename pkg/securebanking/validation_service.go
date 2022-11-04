package securebanking

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"secure-banking-uk-initializer/pkg/common"
)

func CreateAmValidationService(cookie *http.Cookie) {
	zap.L().Info("Creating Validation Service")

	var validationServiceConfig map[string]interface{}
	common.Unmarshal(common.Config.Environment.Paths.ConfigSecureBanking+"create-validation-service.json", &common.Config, &validationServiceConfig)
	path := fmt.Sprintf("https://%s/am/json/realms/root/realms/%s/realm-config/services/validation?_action=create",
		common.Config.Hosts.IdentityPlatformFQDN, common.Config.Identity.AmRealm)

	resp, err := restClient.R().
		SetHeader("Accept", "*/*").
		SetHeader("Content-Type", "application/json").
		SetHeader("Connection", "keep-alive").
		SetHeader("X-Requested-With", "XMLHttpRequest").
		SetHeader("Accept-API-Version", "protocol=1.0,resource=1.0").
		SetContentLength(true).
		SetCookie(cookie).
		SetBody(validationServiceConfig).
		Post(path)

	zap.S().Info("resp is " + resp.String())
	if resp != nil && resp.StatusCode() == 409 {
		zap.S().Info("Validation Service configuration already exists")
	} else {
		common.RaiseForStatus(err, resp.Error(), resp.StatusCode())
		zap.S().Info("Created Validation Service")
	}
}
