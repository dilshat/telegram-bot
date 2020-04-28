package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadScript(t *testing.T) {
	s, err := loadScript()

	assert.NotEmpty(t, s)
	assert.Nil(t, err)

	//negative test
	os.Setenv("SCRIPTS_PATH", "")

	s, err = loadScript()

	assert.Empty(t, s)
	assert.NotNil(t, err)

}
