package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnv(t *testing.T) {
	os.Setenv("key", "val")

	val := GetEnv("key", "")

	assert.Equal(t, "val", val)

	os.Unsetenv("key")
	val = GetEnv("ke2", "default")

	assert.Equal(t, "default", val)
}

func TestGetEnvAsInt(t *testing.T) {
	os.Setenv("key", "123")

	val := GetEnvAsInt("key", 0)

	assert.Equal(t, 123, val)

	os.Unsetenv("key")

	val = GetEnvAsInt("key", 0)

	assert.Equal(t, 0, val)

	os.Setenv("key", "blabla")

	val = GetEnvAsInt("key", 1)

	assert.Equal(t, 1, val)
}

func TestFileExists(t *testing.T) {
	f, err := ioutil.TempFile(os.TempDir(), "util_test")
	if err != nil {
		t.Error(err)
	}
	defer func() {
		f.Close()
	}()

	assert.True(t, FileExists(f.Name()))

	assert.False(t, FileExists(""))
}

func TestGetFileType(t *testing.T) {
	assert.Equal(t, PHOTO, GetFileType("attachments/smile.jpg"))

	assert.Equal(t, VIDEO, GetFileType("attachments/puppy.mp4"))

	assert.Equal(t, AUDIO, GetFileType("attachments/music.mp3"))

	assert.Equal(t, OTHER, GetFileType("attachments/document.txt"))
}

func TestParseFileType(t *testing.T) {
	assert.Equal(t, PHOTO, ParseFileType("photo"))

	assert.Equal(t, VIDEO, ParseFileType("video"))

	assert.Equal(t, AUDIO, ParseFileType("audio"))

	assert.Equal(t, OTHER, ParseFileType("docx"))
}

func TestIsValidUrl(t *testing.T) {
	assert.False(t, isValidUrl(""))

	assert.True(t, isValidUrl("http://www.kg"))
}

func TestReadFile(t *testing.T) {
	s, err := ReadFile("attachments/document.txt")

	assert.NotEmpty(t, s)
	assert.Nil(t, err)

	s, err = ReadFile("attachments/blabla")

	assert.Empty(t, s)
	assert.NotNil(t, err)

}

func TestDoGet2(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("OK"))
	}))

	defer server.Close()

	res, err := doGET(server.URL, map[string]interface{}{"a": "b"}, map[string]interface{}{"c": "d"}, 10)

	assert.NotEmpty(t, res)
	assert.Nil(t, err)
}

func TestDoPost2(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("OK"))
	}))

	defer server.Close()

	res, err := doPOST(server.URL, map[string]interface{}{"a": "b"}, map[string]interface{}{"c": "d"}, 10)

	assert.NotEmpty(t, res)
	assert.Nil(t, err)
}
