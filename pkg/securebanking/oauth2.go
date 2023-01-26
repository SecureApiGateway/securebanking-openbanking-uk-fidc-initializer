package securebanking

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"secure-banking-uk-initializer/pkg/httprest"
	"secure-banking-uk-initializer/pkg/types"

	"secure-banking-uk-initializer/pkg/common"

	"go.uber.org/zap"
)

func CreateSecureBankingRemoteConsentService() {
	remoteConsentId := common.Config.Identity.RemoteConsentId
	if remoteConsentExists(remoteConsentId) {
		zap.L().Info("Remote consent exists. skipping")
		return
	}
	zap.L().Info("Creating remote consent service")
	rc := &types.RemoteConsent{
		RemoteConsentRequestEncryptionAlgorithm: types.InheritedValueString{
			Inherited: false,
			Value:     "RSA-OAEP-256",
		},
		PublicKeyLocation: types.InheritedValueString{
			Inherited: false,
			Value:     "jwks",
		},
		JwksCacheTimeout: types.InheritedValueInt{
			Inherited: false,
			Value:     0,
		},
		RemoteConsentResponseSigningAlg: types.InheritedValueString{
			Inherited: false,
			Value:     "PS256",
		},
		RemoteConsentRequestSigningAlgorithm: types.InheritedValueString{
			Inherited: false,
			Value:     "RS256",
		},
		JwkSet: types.JwkSet{
			Inherited: false,
			Value:     "{\"keys\":[{\"kty\":\"RSA\",\"kid\":\"jwt-signer\",\"use\":\"sig\",\"x5t\":\"jRGHM61NuJP8vPKqGVxyupgxmdc\",\"x5c\":[\"MIIC+zCCAeOgAwIBAgIIMw5JR8AULgwwDQYJKoZIhvcNAQELBQAwLDEqMCgGA1UEAxMhVGVzdCBUcnVzdGVkIERpcmVjdG9yeSBKV1QgU2lnbmVyMB4XDTIzMDEyNTA5MjUyNloXDTMzMDEyMjA5MjUyNlowLDEqMCgGA1UEAxMhVGVzdCBUcnVzdGVkIERpcmVjdG9yeSBKV1QgU2lnbmVyMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAnkxmao72diRjgIO0Y3+dVqC1dWbbkES1n7gOAhYjACKv2GkyW6gTq5tqA96kgd9EiC40Bi1Eqg3KGoqF6oWrH0+X9PFZMOIsBDGvgZAapcqWNMjKMSl0t4EdbCmm7l1dBp52XsEY9n2nrCtvCPkXiiP8t0qV3iljEOsDpIdiSkTt4tOGJ+td317QWJof6u2kykGyZTNHAqJcvhO+AEdFqbS3McGYhiakUZ87bCaUJLABpt4llVKU7kW1lNdWeUXAw4/J9iq0yd2y9aX/JUxmVZbyS1bV72GhhpmSPJYY44riVGkeABM/nF17a5x+jPj4GlIkxPeqeQBtzQ8qMftrUwIDAQABoyEwHzAdBgNVHQ4EFgQUP5HpE/iGKyKTJ7WzBOAYF2LD2UwwDQYJKoZIhvcNAQELBQADggEBAI0fAuUQbFWqoB+/gwIjLi5XMtA6v9hW3Vna54q2iw7ZzPso+bsm8+mBFuJHHZkgbIuDi79U8ZgotWKnVUMD2+4Ytw3o8rbVCXr4WcwrbyClEpxys4Cv4sresmDonlyER/4+vddfexR3gajPxRzT5+byFwjLbqX6d/n4KCjJTdqj6tRVnPJoYmeb5gZNa2dVfDrOAMie2RY0qjCNS0eYhsRQdKUPKL/B/Nt92bxFhveQg8YN02cgKbmwseXfOawDauQvFV2pwtzA7y1TxEpW2gMWz8Cwt/o6ETxwT5E40M/ByJ4Td+8JuH29d9qm2gcEVGxlXdj/ivxq94HGJAEVvXg=\"],\"n\":\"nkxmao72diRjgIO0Y3-dVqC1dWbbkES1n7gOAhYjACKv2GkyW6gTq5tqA96kgd9EiC40Bi1Eqg3KGoqF6oWrH0-X9PFZMOIsBDGvgZAapcqWNMjKMSl0t4EdbCmm7l1dBp52XsEY9n2nrCtvCPkXiiP8t0qV3iljEOsDpIdiSkTt4tOGJ-td317QWJof6u2kykGyZTNHAqJcvhO-AEdFqbS3McGYhiakUZ87bCaUJLABpt4llVKU7kW1lNdWeUXAw4_J9iq0yd2y9aX_JUxmVZbyS1bV72GhhpmSPJYY44riVGkeABM_nF17a5x-jPj4GlIkxPeqeQBtzQ8qMftrUw\",\"e\":\"AQAB\",\"d\":\"TXx4ZZC1c_879ZyCoHHHQrBIIr_GmkgH37bopHzRhS91hCl6TQNpHYdlzZ8eR7RSh5aWQK_H_LXjSDqmcrYbnagIag356sPLaAIqNvGjGaShAhWHY6k3SNwa2udIn0V0U9kdeCRtY7r-aHiaUXoc27Hh8pa_5Y-8vntLvS3IHzH0ppMrox2rRqfRPRbeggF51-TyKDvC57VbpXr16jBgg1hDTXhT7RmPKzG03U7DAusV2RuHPN0fgbbfJIrAwq6Xx4dzkfy7eeMAp8G8RU1XlsqkelqrR5i5KKTpimgUw61X4cjnzW3ARDSN3NiX0qK9z3s99EYzgbeQ2-pbgetEgQ\",\"p\":\"-rxLVsdhmb1jNvWvPol3RKAl8y7bveGG416t_yeiUGJ6j7nKhEXZKpfwgHjrX9_EvZK74gCzyj1WwYmfHuEXrhs_fPoHVvmI_Z0hnGciJAcsZGwFjaX8UAh0UMJe13gRovQITRAEa-PDyBVVCpOzuID7WzCa_1r11iEIAMl3We0\",\"q\":\"oZ9BbqgXwpDpqS3-MgEQQ53sriAkX9xZX13wsu5MxeEDMslVKzC-4sTIpiVhgYCwVcr2NhtqmWpnGBYfBZBlh9HQ6hWsAVHwZROM6H8DmLaWAN9zthL4yJunumWEw7AmKa19uSN4De2MqRyEEGu3Q4oSoFeuMcHlKVP3ZdS8Mj8\",\"dp\":\"20WJ2UUZ_JocRvcKn3UWQYSQS6BB-UdecD7fiVE-6G-G4WRIT-7JoS8o2yFkLf8CDgthlZ6pnIaR9UUGS7nrvI4FsqTxgEVPYQNmXmHvAHTphJTXMm3jPpZ2Kh4hVOui-M_S5pnIVBXmpHbLXSmYNRVPPAtAsWn5rZ5LYmzNnG0\",\"dq\":\"gT6RXKQfHAA2WovlEqe3EdtgQX6CmAXlklXU_cGCs1rU6_AEo50_iKhErFiIOL9oQ7MiYPJbtINaHfDSrehiyNIIdqkN-6BC1RFlRScNATpBikEmjxtsxz9ySaSVHsfmPL34I-0xPDISocmj8l2xF9l6O67iJfggAYSb-lq8hJE\",\"qi\":\"OAwS016BDjIsVELRkHjobgWSYAhfKind0WCKp-uFbUzSkk2XqlCZXyQerjd61tGStZC4k1dl9wedGNU31SK8bQYQOuR5UiHjuZYBP0MF46wMuN6Adpgv1R1aGahYYxX_rL_4XdbfKlsUSSoY_FN2puBFBglify5Frz08Lt8UDrk\"}]}",
		},
		JwkStoreCacheMissCacheTime: types.InheritedValueInt{
			Inherited: false,
			Value:     0,
		},
		RemoteConsentResponseEncryptionMethod: types.InheritedValueString{
			Inherited: false,
			Value:     "A128GCM",
		},
		RemoteConsentRedirectURL: types.InheritedValueString{
			Inherited: false,
			Value:     fmt.Sprintf("https://%s", common.Config.Hosts.RcsUiFQDN),
		},
		RemoteConsentRequestEncryptionEnabled: types.InheritedValueBool{
			Inherited: false,
			Value:     false,
		},
		RemoteConsentRequestEncryptionMethod: types.InheritedValueString{
			Inherited: false,
			Value:     "A128GCM",
		},
		RemoteConsentResponseEncryptionAlgorithm: types.InheritedValueString{
			Inherited: false,
			Value:     "RSA-OAEP-256",
		},
		RequestTimeLimit: types.InheritedValueInt{
			Inherited: false,
			Value:     0,
		},
		JwksURI: types.InheritedValueString{
			Inherited: false,
			//Value:     "http://securebanking-openbanking-uk-rcs:8080/api/rcs/consent/jwk_pub",
		},
		Type: types.Type{
			ID:         "RemoteConsentAgent",
			Name:       "OAuth2 Remote Consent Service",
			Collection: true,
		},
		Userpassword: common.Config.Ig.IgRcsSecret,
	}
	path := "/am/json/realms/root/realms/" + common.Config.Identity.AmRealm + "/realm-config/agents/RemoteConsentAgent/" + remoteConsentId

	s := httprest.Client.Put(path, rc, map[string]string{
		"Accept":             "*/*",
		"Connection":         "keep-alive",
		"Accept-API-Version": "protocol=2.0,resource=1.0",
	})

	zap.S().Infow("Remote Consent Service", "statusCode", s)
}

func remoteConsentExists(name string) bool {
	path := "/am/json/realms/root/realms/" + common.Config.Identity.AmRealm + "/realm-config/agents/RemoteConsentAgent?_queryFilter=true&_pageSize=10&_fields=agentgroup"
	consent := &types.AmResult{}
	b, _ := httprest.Client.Get(path, map[string]string{
		"Accept":             "application/json",
		"X-Requested-With":   "ForgeRock Identity Cloud Postman Collection",
		"Accept-Api-Version": "protocol=2.0,resource=1.0",
	})

	err := json.Unmarshal(b, consent)
	if err != nil {
		panic(err)
	}

	return common.Find(name, consent, func(r *types.Result) string {
		return r.ID
	})
}

func CreateSoftwarePublisherAgentOBRI() {
	if softwarePublisherAgentExists(common.Config.Identity.ObriSoftwarePublisherAgent) {
		zap.L().Info("Skipping creation of Software publisher agent")
		return
	}

	zap.L().Info("Creating software publisher agent")
	pa := types.PublisherAgent{
		PublicKeyLocation: types.InheritedValueString{
			Inherited: false,
			Value:     "jwks_uri",
		},
		JwksCacheTimeout: types.InheritedValueInt{
			Inherited: false,
			Value:     3600000,
		},
		SoftwareStatementSigningAlgorithm: types.InheritedValueString{
			Inherited: false,
			Value:     "PS256",
		},
		JwkSet: types.JwkSet{
			Inherited: false,
		},
		Issuer: types.InheritedValueString{
			Inherited: false,
			Value:     "ForgeRock",
		},
		JwkStoreCacheMissCacheTime: types.InheritedValueInt{
			Inherited: false,
			Value:     60000,
		},
		JwksURI: types.InheritedValueString{
			Inherited: false,
			Value:     "https://service.directory.ob.forgerock.financial/api/directory/keys/jwk_uri",
		},
	}
	path := "/am/json/realms/root/realms/" + common.Config.Identity.AmRealm + "/realm-config/agents/SoftwarePublisher/" + common.Config.Identity.ObriSoftwarePublisherAgent
	s := httprest.Client.Put(path, pa, map[string]string{
		"Accept":             "*/*",
		"Connection":         "keep-alive",
		"Accept-API-Version": "protocol=2.0,resource=1.0",
	})

	zap.S().Infow("Software Publisher Agent", "statusCode", s)
}

func CreateSoftwarePublisherAgentTestPublisher() {
	if softwarePublisherAgentExists(common.Config.Identity.TestSoftwarePublisherAgent) {
		zap.L().Info("Skipping creation of Software publisher agent")
		return
	}

	zap.L().Info("Creating software publisher agent")
	pa := types.PublisherAgent{
		Userpassword: common.Config.Ig.IgSsaSecret,
		PublicKeyLocation: types.InheritedValueString{
			Inherited: false,
			Value:     "jwks_uri",
		},
		JwksCacheTimeout: types.InheritedValueInt{
			Inherited: false,
			Value:     3600000,
		},
		SoftwareStatementSigningAlgorithm: types.InheritedValueString{
			Inherited: false,
			Value:     "HS256",
		},
		JwkSet: types.JwkSet{
			Inherited: false,
		},
		Issuer: types.InheritedValueString{
			Inherited: false,
			Value:     "test-publisher",
		},
		JwkStoreCacheMissCacheTime: types.InheritedValueInt{
			Inherited: false,
			Value:     60000,
		},
		JwksURI: types.InheritedValueString{
			Inherited: false,
		},
	}
	path := "/am/json/realms/root/realms/" + common.Config.Identity.AmRealm + "/realm-config/agents/SoftwarePublisher/" + common.Config.Identity.TestSoftwarePublisherAgent
	s := httprest.Client.Put(path, pa, map[string]string{
		"Accept":             "*/*",
		"Connection":         "keep-alive",
		"Accept-API-Version": "protocol=2.0,resource=1.0",
	})

	zap.S().Infow("Software Publisher Agent", "statusCode", s)
}

func softwarePublisherAgentExists(name string) bool {
	path := "/am/json/realms/root/realms/" + common.Config.Identity.AmRealm + "/realm-config/agents/SoftwarePublisher?_queryFilter=true&_pageSize=10&_fields=agentgroup"
	agent := &types.AmResult{}
	b, _ := httprest.Client.Get(path, map[string]string{
		"Accept":             "application/json",
		"X-Requested-With":   "ForgeRock Identity Cloud Postman Collection",
		"Accept-Api-Version": "protocol=2.0,resource=1.0",
	})

	err := json.Unmarshal(b, agent)
	if err != nil {
		panic(err)
	}

	return common.Find(name, agent, func(r *types.Result) string {
		return r.ID
	})
}

// CreateOIDCClaimsScript -
func CreateOIDCClaimsScript(cookie *http.Cookie) string {

	zap.L().Info("Creating OIDC claims script")
	b, err := ioutil.ReadFile(common.Config.Environment.Paths.ConfigSecureBanking + "oidc-script.json")
	if err != nil {
		panic(err)
	}

	path := fmt.Sprintf("https://%s/am/json/"+common.Config.Identity.AmRealm+"/scripts/?_action=create", common.Config.Hosts.IdentityPlatformFQDN)

	claimsScript := &types.RequestScript{}

	err = json.Unmarshal(b, claimsScript)
	if err != nil {
		panic(err)
	}

	id := httprest.GetScriptIdByName(claimsScript.Name)
	if id != "" {
		zap.L().Info("Script exists")
		return id
	}

	resp, err := restClient.R().
		SetHeader("Accept", "*/*").
		SetHeader("Content-Type", "application/json").
		SetHeader("Connection", "keep-alive").
		SetHeader("Accept-API-Version", "protocol=2.0,resource=1.0").
		SetContentLength(true).
		SetCookie(cookie).
		SetResult(claimsScript).
		SetBody(b).
		Post(path)

	common.RaiseForStatus(err, resp.Error(), resp.StatusCode())

	zap.S().Infow("OIDC claims script", "statusCode", resp.StatusCode(), "claimsScriptID", claimsScript.ID, "createdBy", claimsScript.CreatedBy)
	return claimsScript.ID
}

// UpdateOAuth2Provider - update the oauth 2 provider, must supply the claimScript ID
func UpdateOAuth2Provider(claimsScriptID string) {
	zap.S().Info("UpdateOAuth2Provider() Creating OAuth2Provider service in the " + common.Config.Identity.AmRealm + " realm")

	oauth2Provider := &types.OAuth2Provider{}
	err := common.Unmarshal(common.Config.Environment.Paths.ConfigSecureBanking+"oauth2provider-update.json", &common.Config, oauth2Provider)
	if err != nil {
		panic(err)
	}

	if oauth2ProviderExists(oauth2Provider.Type.ID) {
		zap.L().Info("UpdateOAuth2Provider() OAuth2 provider exists")
		return
	}

	zap.S().Infof("Pushing the following config %+v", oauth2Provider)

	oauth2Provider.PluginsConfig.OidcClaimsScript = claimsScriptID
	zap.S().Infow("UpdateOAuth2Provider() Updating OAuth2 provider", "claimScriptId", claimsScriptID)
	path := "/am/json/" + common.Config.Identity.AmRealm + "/realm-config/services/oauth-oidc"
	zap.S().Info("UpdateOAuth2Provider() Updating OAuth2Provider via the following path {}", path)
	s := httprest.Client.Put(path, oauth2Provider, map[string]string{
		"Accept":           "*/*",
		"Content-Type":     "application/json",
		"Connection":       "keep-alive",
		"X-Requested-With": "ForgeRock Identity Cloud Postman Collection",
	})

	zap.S().Infow("UpdateOAuth2Provider() OAuth2 provider", "statusCode", s)
}

func oauth2ProviderExists(id string) bool {
	path := "/am/json/realms/root/realms/" + common.Config.Identity.AmRealm + "/realm-config/services?_queryFilter=true"
	r := &types.AmResult{}
	b, _ := httprest.Client.Get(path, map[string]string{
		"Accept":             "application/json",
		"X-Requested-With":   "ForgeRock Identity Cloud Postman Collection",
		"Accept-Api-Version": "protocol=1.0,resource=1.0",
	})

	err := json.Unmarshal(b, r)
	if err != nil {
		panic(err)
	}

	return common.Find(id, r, func(r *types.Result) string {
		return r.ID
	})
}

func CreateBaseURLSourceService(cookie *http.Cookie) {
	zap.S().Info("Creating BaseURLSource service in the " + common.Config.Identity.AmRealm + " realm")

	s := &types.Source{}
	err := common.Unmarshal(common.Config.Environment.Paths.ConfigSecureBanking+"create-base-url-source.json", &common.Config, s)
	if err != nil {
		panic(err)
	}

	path := fmt.Sprintf("https://%s/am/json/realms/root/realms/"+common.Config.Identity.AmRealm+"/realm-config/services/baseurl?_action=create",
		common.Config.Hosts.IdentityPlatformFQDN)

	resp, err := restClient.R().
		SetHeader("Accept", "application/json").
		SetHeader("Accept-API-Version", "protocol=1.0,resource=1.0").
		SetHeader("Content-Type", "application/json").
		SetContentLength(true).
		SetCookie(cookie).
		SetBody(s).
		Post(path)

	zap.S().Info("resp is " + resp.String())
	if resp != nil && resp.StatusCode() == 409 {
		zap.S().Info("Did not create BaseURLSource service in " + common.Config.Identity.AmRealm + " realm. It already exists.")
	} else {
		common.RaiseForStatus(err, resp.Error(), resp.StatusCode())
		zap.S().Info("Created Base URL Service in AM's " + common.Config.Identity.AmRealm + " realm")
	}
}
