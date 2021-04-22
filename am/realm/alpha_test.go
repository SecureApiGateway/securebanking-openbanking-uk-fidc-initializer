package realm

import (
	"io/ioutil"
	"testing"

	"github.com/secureBankingAcceleratorToolkit/securebanking-openbanking-uk-fidc-initialiszer/am"
	mocks "github.com/secureBankingAcceleratorToolkit/securebanking-openbanking-uk-fidc-initialiszer/mocks/am"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFindExistingAlphaClient(t *testing.T) {
	mockRestReaderWriter := &mocks.RestReaderWriter{}
	am.Client = mockRestReaderWriter
	buffer, _ := ioutil.ReadFile("client-check-test.json")
	mockRestReaderWriter.On("Get", mock.Anything, mock.Anything).
		Return(buffer)

	b := AlphaClientsExist("Doesnt existy")
	assert.False(t, b)
	mockRestReaderWriter.AssertCalled(t, "Get", mock.Anything, mock.Anything)

	b = AlphaClientsExist("ig-client")
	assert.True(t, b)
}

func TestMangedObjectExists(t *testing.T) {
	mockRestReaderWriter := &mocks.RestReaderWriter{}
	am.Client = mockRestReaderWriter
	buffer, _ := ioutil.ReadFile("managed-objects-test.json")
	mockRestReaderWriter.On("Get", mock.Anything, mock.Anything).
		Return(buffer)

	b := ManagedObjectExists("api_client")
	assert.True(t, b)
	mockRestReaderWriter.AssertCalled(t, "Get", mock.Anything, mock.Anything)

	b = ManagedObjectExists("xyz")
	assert.False(t, b)
}
