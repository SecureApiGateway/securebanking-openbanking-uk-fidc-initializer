package securebanking

import (
	"encoding/json"
	"io/ioutil"
	"secure-banking-uk-initializer/pkg/common"
	"secure-banking-uk-initializer/pkg/httprest"
	"strings"

	"go.uber.org/zap"
)

//objectNames - retrieve filenames from a path. the .json extension will be trimed and
//  a list of filenames will be returned
func objectNames(relativePath string) []string {
	files, err := ioutil.ReadDir(relativePath)
	if err != nil {
		zap.L().Fatal(err.Error())
	}

	var names []string
	for _, f := range files {
		name := strings.TrimSuffix(f.Name(), ".json")
		names = append(names, name)
	}
	return names
}

//MissingObjects - return a list of missing managed object names in idm.
//  supply an array of managed object names to query against.
func missingObjects(objectNames []string) []string {
	path := "/openidm/config/managed"
	result := &OBManagedObjects{}
	b, _ := httprest.Client.Get(path, map[string]string{
		"Accept":           "application/json",
		"X-Requested-With": "ForgeRock Identity Cloud Postman Collection",
	})

	err := json.Unmarshal(b, result)
	if err != nil {
		panic(err)
	}

	var missObjects = objectNames
	for _, o := range result.Objects {
		for i, objectName := range missObjects {
			zap.S().Infow("checking", "object", o)
			if strings.Contains(o.Name, objectName) {
				zap.S().Infow("ManagedObject found", "name", objectName)
				missObjects = append(missObjects[:i], missObjects[i+1:]...)
				break
			}
		}
	}
	return missObjects
}

// OBManagedObjects model
type OBManagedObjects struct {
	ID      string `json:"_id"`
	Objects []struct {
		Name string `json:"name"`
	} `json:"objects"`
}

// AddOBManagedObjects - Add managed objects to IDM. This will look for json in the managed objects OB config directory
//  and add them to IDM if they dont already exist.
func AddOBManagedObjects() {
	configPath := common.Config.Environment.Paths.ConfigSecureBanking + "managed-objects/"
	managedObjectFilenames := objectNames(configPath)
	mObjects := missingObjects(managedObjectFilenames)

	zap.S().Infow("Attempting to add Managed Object definitions", "definitions", mObjects)
	patches := make([]map[string]interface{}, 0)
	for _, o := range mObjects {
		patches = append(patches, unmarshallManagedObjectPatch(o, configPath)...)

	}
	patchManagedObjects(patches)
}

func patchManagedObjects(managedObjectPatches []map[string]interface{}) {
	zap.S().Infow("Patching Managed Object definitions", "patches", managedObjectPatches)

	path := "/openidm/config/managed"
	s := httprest.Client.Patch(path, managedObjectPatches, map[string]string{
		"Accept":       "*/*",
		"Content-Type": "application/json",
		"Connection":   "keep-alive",
	})

	zap.S().Infow("Managed object created", "statusCode", s)
}

func unmarshallManagedObjectPatch(name string, objectFolderPath string) []map[string]interface{} {
	b, err := ioutil.ReadFile(objectFolderPath + name + ".json")
	if err != nil {
		panic(err)
	}

	patch := make([]map[string]interface{}, 0)
	err = json.Unmarshal(b, &patch)
	if err != nil {
		panic(err)
	}
	return patch
}

func CreateApiJwksEndpoint() {
	zap.L().Info("Creating API JWKS Endpoint")
	b, err := ioutil.ReadFile(common.Config.Environment.Paths.ConfigSecureBanking + "create-jwks-endpoint.json")
	if err != nil {
		panic(err)
	}

	path := "/openidm/config/endpoint/apiclientjwks"
	s := httprest.Client.Put(path, b, map[string]string{
		"Accept":       "*/*",
		"Content-Type": "application/json",
		"Connection":   "keep-alive",
	})

	zap.S().Infow("JWKS endpoint", "statusCode", s)
}
