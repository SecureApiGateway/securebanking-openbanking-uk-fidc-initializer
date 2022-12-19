package securebanking

import (
	"go.uber.org/zap"
	"io/ioutil"
	"secure-banking-uk-initializer/pkg/common"
	"secure-banking-uk-initializer/pkg/httprest"
)

const FOLDER = "PSD2-authentication-trees-config/"

// CreateSecureBankingPSD2AuthenticationTrees will attempt to create the PSD2 Authentication trees
func CreateSecureBankingPSD2AuthenticationTrees() {
	zap.L().Info("Attempt to create PSD2CustomerAuthentication and PSD2SecureCustomerAuthentication trees")
	createPSD2CustomerAuthenticationUsernameNode()
	createPSD2CustomerAuthenticationPasswordNode()
	createPSD2CustomerAuthenticationTree()
	createPSD2SecureCustomerAuthenticationUsernameNode()
	createPSD2SecureCustomerAuthenticationPasswordNode()
	createPSD2SecureCustomerAuthenticationTree()
}

func createPSD2CustomerAuthenticationUsernameNode() {
	zap.L().Info("Creating PSD2CustomerAuthentication Username Node")
	b, err := ioutil.ReadFile(common.Config.Environment.Paths.ConfigSecureBanking +
		FOLDER + "PSD2CustomerAuthentication-username-node.json")
	if err != nil {
		panic(err)
	}

	zap.S().Debugw("PSD2CustomerAuthentication username node", "body", string(b))
	path := "/am/json/realms/root/realms/" + common.Config.Identity.AmRealm + "/realm-config/authentication/authenticationtrees/nodes/UsernameCollectorNode/ada9ef86-d550-4591-b9dc-5751e7adbb62"
	status := httprest.Client.Put(path, b, map[string]string{
		"Accept":             "*/*",
		"Content-Type":       "application/json",
		"Accept-Api-Version": "protocol=2.0, resource=1.0",
	})

	zap.S().Infow("PSD2CustomerAuthentication Node Username", "statusCode", status)
}

func createPSD2CustomerAuthenticationPasswordNode() {
	zap.L().Info("Creating PSD2CustomerAuthentication Password Node")
	b, err := ioutil.ReadFile(common.Config.Environment.Paths.ConfigSecureBanking +
		FOLDER + "PSD2CustomerAuthentication-password-node.json")
	if err != nil {
		panic(err)
	}

	zap.S().Debugw("PSD2CustomerAuthentication Password node", "body", string(b))
	path := "/am/json/realms/root/realms/" + common.Config.Identity.AmRealm + "/realm-config/authentication/authenticationtrees/nodes/PasswordCollectorNode/1db869b1-09de-4a8e-b340-e0563891c3bf"
	status := httprest.Client.Put(path, b, map[string]string{
		"Accept":             "*/*",
		"Content-Type":       "application/json",
		"Accept-Api-Version": "protocol=2.0, resource=1.0",
	})

	zap.S().Infow("PSD2CustomerAuthentication Node Password", "statusCode", status)
}

func createPSD2CustomerAuthenticationTree() {
	zap.L().Info("Creating PSD2CustomerAuthentication tree")
	b, err := ioutil.ReadFile(common.Config.Environment.Paths.ConfigSecureBanking +
		FOLDER + "PSD2CustomerAuthentication-tree.json")
	if err != nil {
		panic(err)
	}

	zap.S().Debugw("PSD2CustomerAuthentication tree", "body", string(b))
	path := "/am/json/realms/root/realms/" + common.Config.Identity.AmRealm + "/realm-config/authentication/authenticationtrees/trees/PSD2CustomerAuthentication"
	status := httprest.Client.Put(path, b, map[string]string{
		"Accept":             "*/*",
		"Content-Type":       "application/json",
		"Accept-Api-Version": "protocol=2.0, resource=1.0",
	})

	zap.S().Infow("PSD2CustomerAuthentication tree", "statusCode", status)
}

func createPSD2SecureCustomerAuthenticationUsernameNode() {
	zap.L().Info("Creating PSD2SecureCustomerAuthentication Username Node")
	b, err := ioutil.ReadFile(common.Config.Environment.Paths.ConfigSecureBanking +
		FOLDER + "PSD2SecureCustomerAuthentication-username-node.json")
	if err != nil {
		panic(err)
	}

	zap.S().Debugw("PSD2CustomerAuthentication username node", "body", string(b))
	path := "/am/json/realms/root/realms/" + common.Config.Identity.AmRealm + "/realm-config/authentication/authenticationtrees/nodes/UsernameCollectorNode/ee0efdc1-9fba-4323-95ef-ec468f6ad30c"
	status := httprest.Client.Put(path, b, map[string]string{
		"Accept":             "*/*",
		"Content-Type":       "application/json",
		"Accept-Api-Version": "protocol=2.0, resource=1.0",
	})

	zap.S().Infow("PSD2SecureCustomerAuthentication Node Username", "statusCode", status)
}

func createPSD2SecureCustomerAuthenticationPasswordNode() {
	zap.L().Info("Creating PSD2SecureCustomerAuthentication Password Node")
	b, err := ioutil.ReadFile(common.Config.Environment.Paths.ConfigSecureBanking +
		FOLDER + "PSD2SecureCustomerAuthentication-password-node.json")
	if err != nil {
		panic(err)
	}

	zap.S().Debugw("PSD2CustomerAuthentication Password node", "body", string(b))
	path := "/am/json/realms/root/realms/" + common.Config.Identity.AmRealm + "/realm-config/authentication/authenticationtrees/nodes/PasswordCollectorNode/4785b3c1-5dc9-4883-b01e-2f1b6bfda50e"
	status := httprest.Client.Put(path, b, map[string]string{
		"Accept":             "*/*",
		"Content-Type":       "application/json",
		"Accept-Api-Version": "protocol=2.0, resource=1.0",
	})

	zap.S().Infow("PSD2SecureCustomerAuthentication Node Password", "statusCode", status)
}

func createPSD2SecureCustomerAuthenticationTree() {
	zap.L().Info("Creating PSD2SecureCustomerAuthentication tree")
	b, err := ioutil.ReadFile(common.Config.Environment.Paths.ConfigSecureBanking +
		FOLDER + "PSD2SecureCustomerAuthentication-tree.json")
	if err != nil {
		panic(err)
	}

	zap.S().Debugw("Login tree PSD2CustomerAuthentication", "body", string(b))
	path := "/am/json/realms/root/realms/" + common.Config.Identity.AmRealm + "/realm-config/authentication/authenticationtrees/trees/PSD2SecureCustomerAuthentication"
	status := httprest.Client.Put(path, b, map[string]string{
		"Accept":             "*/*",
		"Content-Type":       "application/json",
		"Accept-Api-Version": "protocol=2.0, resource=1.0",
	})

	zap.S().Infow("PSD2SecureCustomerAuthentication tree", "statusCode", status)
}

// ConfigureRealmDefaultUserAuthenticationService
// This configures the default user authentication service to use for a particular realm. This service is used as the
// fallback auth service when a more specific service isn't configured/specified.
//
// This is driven off the IDENTITY.DEFAULT_USER_AUTHENTICATION_SERVICE configuration value
//
// For FIDC environments: this configuration is mandatory and will cause the program to exit if it is missing.
// For CDK environments: configuring this is optional, the CDK comes preconfigured with a sensible default.
func ConfigureRealmDefaultUserAuthenticationService() {
	if common.Config.Identity.DefaultUserAuthenticationService == "" {
		if common.Config.Environment.Type != "FIDC" {
			zap.L().Info("No DefaultUserAuthenticationService configuration found, nothing to do")
			return
		} else {
			panic("Configuration: DEFAULT_USER_AUTHENTICATION_SERVICE is required for FIDC environments")
		}
	}

	zap.S().Infow("Configuring Default Authentication Service for realm",
		"realm", common.Config.Identity.AmRealm, "authService", common.Config.Identity.DefaultUserAuthenticationService)
	path := "/am/json/realms/root/realms/" + common.Config.Identity.AmRealm + "/realm-config/authentication"

	status := httprest.Client.Put(path, map[string]string{
		"orgConfig": common.Config.Identity.DefaultUserAuthenticationService,
	}, map[string]string{
		"Accept":             "application/json",
		"Content-Type":       "application/json",
		"Accept-Api-Version": "protocol=2.0, resource=1.0",
	})

	zap.S().Infow("Configure Default Authentication Service response", "statusCode", status)
}
