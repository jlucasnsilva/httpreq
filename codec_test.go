package httpreq

import (
	"fmt"
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestCreateFileAndParse(t *testing.T) {
	password, passwordSet := os.LookupEnv("PASSWORD")
	if !passwordSet {
		password = "123456"
		os.Setenv("PASSWORD", password)
		defer func() {
			os.Unsetenv("PASSWORD")
		}()
	}

	afs := afero.NewMemMapFs()
	fileName := "text"
	expected := example

	expected.Macros["password"] = password
	expected.EndpointCalls[0].Payload = fmt.Sprintf("{ \"email\": \"mr.frodo@lotr.com\", \"password\": \"%v\" }", password)
	createFileErr := createFile(afs, fileName)
	g, parseErr := parse(afs, fileName)

	assert.Nil(t, createFileErr)
	assert.Nil(t, parseErr)
	assert.Equal(t, expected, *g)
}
