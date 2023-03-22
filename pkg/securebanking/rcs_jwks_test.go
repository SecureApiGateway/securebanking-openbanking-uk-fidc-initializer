package securebanking

import (
	"testing"
)

// This test can be used to generate a JWKS for a given rsaPublicKey PEM and keyId
func TestCreateRcsJwks(t *testing.T) {
	// keyId needs to make k8s configmap value for target env e.g. RCS_CONSENT_RESPONSE_JWT_SIGNINGKEYID: rcs-jwt-signer
	var keyId = "rcs-jwt-signer"
	// Drop the PEM file contents into this var
	var rsaPublicKey = "-----BEGIN PUBLIC KEY-----\nMIIBojANBgkqhkiG9w0BAQEFAAOCAY8AMIIBigKCAYEAyHZhzne927T5sTbLxZwg\nED4hONhx1QZFV8wKIfnLHtHOC1F4sQwfig0BKiGGf+sK4qZxGG9MENOpqYSahJeg\n6dKd+z3UC3rUgHuhu0zY4VRixZIEmiJIigsAqpwuVfMRw5BkpNl+lrRor3gQeq1q\nvBxhTDxY/pEyhSYvMceV1lGQ6P0zakGke4/aZa99if1frmKj572zBKCgWvFPDyOQ\nQr6i1LCtthwvFiUxy0whjGM/u7HAIAJ4XVcSusqOYoozUsk9QjG9ch1TNGU79haJ\nCgugpGxZB3jQ7LqE/cinbLfUX2rinMtjqz1DncsF2D0MUfiB31exTCy4XD2J0mfd\nP677nekBpfRoV6+1tEGy5K0RqLFqQdVdWFUtpIk3a3U6SwSA1ww8qGK8M8ng2YPY\nxmLZ4FPJroQkYTvhqMI/o8VRCMc74h1GxqY5+ScQoN6zggwPO/FZeWTCQeWAhjdS\nPHjGznqyN5X3qFYOx3CZxAfPCpMp9/ASUFLwlV2C0LZlAgMBAAE=\n-----END PUBLIC KEY-----"
	jwks := CreateRcsJwks(rsaPublicKey, keyId)
	println(jwks)
}
