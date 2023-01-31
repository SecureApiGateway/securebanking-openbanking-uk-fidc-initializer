package securebanking

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-jose/go-jose/v3"
	"go.uber.org/zap"
)

// CreateRcsJwks Produces a JWK Set representing the RCS signing public key
// this can be used to configured AM to trust JWTs signed by RCS.
func CreateRcsJwks(rsaPublicKey string, keyId string) string {
	zap.S().Infow("Creating JWKS", "rsaPublicKey", rsaPublicKey, "keyId", keyId)

	rsaPubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(rsaPublicKey))
	if err != nil {
		panic(err)
	}
	jwk := jose.JSONWebKey{
		KeyID: keyId,
		Key:   rsaPubKey,
		Use:   "sig",
	}

	// API provides no way to build a JWKS and marshall it, so we need to do that manually
	jwkRawJson, err := jwk.MarshalJSON()
	if err != nil {
		panic(err)
	}

	// Unmarshall the JWK back to a map
	var jwkJson map[string]interface{}
	json.Unmarshal(jwkRawJson, &jwkJson)

	// Wrap our key as a single item in a slice
	jwkSlice := []map[string]interface{}{
		jwkJson,
	}

	// Produce the JWKS containing the slice
	jwks := map[string]interface{}{
		"keys": jwkSlice,
	}

	jwksMarshalled, err := json.Marshal(jwks)
	if err != nil {
		panic(err)
	}
	jwksStr := string(jwksMarshalled)
	zap.L().Info("RCS JWKS: " + jwksStr)
	return jwksStr
}
