package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ReneKroon/ttlcache"
	"github.com/dilshat/telegram-bot/mocks"
	"github.com/robertkrimen/otto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yanzay/tbot/v2"
)

var (
	fileID         = "123"
	chatID         = "123"
	msgID          = 123
	userID         = "123"
	text           = "Hello world"
	attachmentsDir = "attachments"
	err            = errors.New("error")
)

func TestSetCacheItem(t *testing.T) {
	cache := ttlcache.NewCache()
	a := &application{cache: cache}
	a.setCacheItem("test", "value")

	if _, ok := cache.Get("test"); !ok {
		t.Errorf("Expected value but got cache miss")
	}
}

func TestGetCacheItem(t *testing.T) {
	cache := ttlcache.NewCache()
	a := &application{cache: cache}
	cache.Set("test", "value")

	if val := a.getCacheItem("test"); val == nil {
		t.Errorf("Expected value but got nil")
	}
}

func TestDelCacheItem(t *testing.T) {
	cache := ttlcache.NewCache()
	a := &application{cache: cache}
	cache.Set("test", "value")
	a.delCacheItem("test")

	if _, ok := cache.Get("test"); ok {
		t.Errorf("Expected cache miss but got value")
	}
}

func TestGetFileLink(t *testing.T) {
	telebot := &mocks.Telebot{}
	filePath := "path"
	fileInfo := &tbot.File{FilePath: filePath}
	telebot.On("GetFileInfo", fileID).Return(fileInfo, nil)
	a := &application{tgClient: telebot}

	a.getFileLink(fileID)

	telebot.AssertExpectations(t)

	//negative test

	telebot = &mocks.Telebot{}
	filePath = "path"
	fileInfo = &tbot.File{FilePath: filePath}
	telebot.On("GetFileInfo", fileID).Return(fileInfo, err)
	a = &application{tgClient: telebot}

	link := a.getFileLink(fileID)

	if link != "" {
		t.Errorf("Exptected empty link but got %s", link)
	}
}

type VmFactoryStub struct {
	vm Vm
}

func (v VmFactoryStub) GetVm() Vm {
	return v.vm
}

type VmStub struct {
}

func (VmStub) Set(name string, value interface{}) error {
	return nil
}
func (VmStub) Run(src interface{}) (otto.Value, error) {
	return otto.Value{}, nil
}

func (VmStub) Call(source string, argumentList ...interface{}) (otto.Value, error) {
	return otto.Value{}, nil
}

func (VmStub) Object(source string) (*otto.Object, error) {
	return &otto.Object{}, nil
}

func (VmStub) Copy() Vm {
	return &VmStub{}
}

func TestReplaceInlineOptions(t *testing.T) {
	telebot := &mocks.Telebot{}

	telebot.On("EditInlineMarkup", chatID, msgID, mock.Anything).Return(0, err)

	a := &application{tgClient: telebot}
	inlineOptions := []map[string]interface{}{}

	a.replaceInlineOptions(chatID, msgID, inlineOptions)

	telebot.AssertExpectations(t)
}

func TestDeleteMsg(t *testing.T) {
	telebot := &mocks.Telebot{}

	telebot.On("DeleteMsg", chatID, msgID).Return(err)

	a := &application{tgClient: telebot}

	a.deleteMessage(chatID, msgID)

	telebot.AssertExpectations(t)
}

func TestEditMsg(t *testing.T) {
	telebot := &mocks.Telebot{}

	telebot.On("EditMsg", chatID, msgID, text, mock.Anything).Return(err)

	a := &application{tgClient: telebot}
	inlineOptions := []map[string]interface{}{}

	a.editMessage(chatID, msgID, text, inlineOptions)

	telebot.AssertExpectations(t)
}

func TestDoGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("OK"))
	}))
	// Close the server when test finishes
	defer server.Close()

	a := &application{}

	res := a.doGet(server.URL, map[string]interface{}{}, map[string]interface{}{}, 10)

	if res != "OK" {
		t.Errorf("Exptected OK but got %s", res)
	}

	//negaive test
	res = a.doGet("", map[string]interface{}{}, map[string]interface{}{}, 10)

	if res != "" {
		t.Errorf("Exptected empty response but got %s", res)
	}
}

func TestDoPost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("OK"))
	}))
	// Close the server when test finishes
	defer server.Close()

	a := &application{}

	res := a.doPOST(server.URL, map[string]interface{}{}, map[string]interface{}{}, 10)

	if res != "OK" {
		t.Errorf("Exptected OK but got %s", res)
	}

	//negaive test
	res = a.doPOST("", map[string]interface{}{}, map[string]interface{}{}, 10)

	if res != "" {
		t.Errorf("Exptected empty response but got %s", res)
	}
}

func TestPromptUser(t *testing.T) {

	//file attachments
	testCases := map[string]string{"smile.jpg": "AttachPhoto", "puppy.mp4": "AttachVideo", "music.mp3": "AttachAudio", "document.txt": "AttachFile"}

	for attachment, method := range testCases {
		telebot := &mocks.Telebot{}
		a := &application{tgClient: telebot, attachmentsDir: attachmentsDir}

		telebot.On(method, userID, filepath.Join(attachmentsDir, attachment), text, mock.AnythingOfType("func(url.Values)")).Return(err)

		a.promptUser(userID, text, attachment)

		telebot.AssertExpectations(t)
	}

	//file forwarding
	testCases = map[string]string{"id:photo": "ForwardPhoto", "id:video": "ForwardVideo", "id:audio": "ForwardAudio", "id:doc": "ForwardFile", "id": "ForwardFile"}

	for attachment, method := range testCases {
		telebot := &mocks.Telebot{}
		a := &application{tgClient: telebot, attachmentsDir: attachmentsDir}

		telebot.On(method, userID, strings.Split(attachment, ":")[0], text, mock.AnythingOfType("func(url.Values)")).Return(err)

		a.promptUser(userID, text, attachment)

		telebot.AssertExpectations(t)
	}

	//send text
	telebot := &mocks.Telebot{}
	a := &application{tgClient: telebot, attachmentsDir: attachmentsDir}

	telebot.On("SendText", userID, text, mock.AnythingOfType("func(url.Values)")).Return(err)

	a.promptUser(userID, text, "")

	telebot.AssertExpectations(t)

	//ignore empty message
	telebot = &mocks.Telebot{}
	a = &application{tgClient: telebot, attachmentsDir: attachmentsDir}

	a.promptUser(userID, "", "")

	if len(telebot.Calls) != 0 {
		t.Errorf("Expected 0 but got %d calls", len(telebot.Calls))
	}
}

func TestSendMessage(t *testing.T) {
	//file attachments with keyboard options
	testCases := map[string]string{"smile.jpg": "AttachPhoto", "puppy.mp4": "AttachVideo", "music.mp3": "AttachAudio", "document.txt": "AttachFile"}
	keyboardOptions := [][]string{{"one", "two"}}
	inlineOptions := []map[string]interface{}{{"one": "1", "two": 2}}
	emptyKeyboardOptions := [][]string{}
	emptyInlineOptions := []map[string]interface{}{}

	for attachment, method := range testCases {
		telebot := &mocks.Telebot{}
		a := &application{tgClient: telebot, attachmentsDir: attachmentsDir}

		telebot.On(method, userID, filepath.Join(attachmentsDir, attachment), text, mock.AnythingOfType("func(url.Values)")).Return(err)

		a.sendMessage(userID, text, keyboardOptions, emptyInlineOptions, attachment)

		telebot.AssertExpectations(t)
	}

	//file attachments with inline options
	testCases = map[string]string{"smile.jpg": "AttachPhoto", "puppy.mp4": "AttachVideo", "music.mp3": "AttachAudio", "document.txt": "AttachFile"}

	for attachment, method := range testCases {
		telebot := &mocks.Telebot{}
		a := &application{tgClient: telebot, attachmentsDir: attachmentsDir}

		telebot.On(method, userID, filepath.Join(attachmentsDir, attachment), text, mock.AnythingOfType("func(url.Values)")).Return(err)

		a.sendMessage(userID, text, emptyKeyboardOptions, inlineOptions, attachment)

		telebot.AssertExpectations(t)
	}

	//file attachments
	testCases = map[string]string{"smile.jpg": "AttachPhoto", "puppy.mp4": "AttachVideo", "music.mp3": "AttachAudio", "document.txt": "AttachFile"}
	for attachment, method := range testCases {
		telebot := &mocks.Telebot{}
		a := &application{tgClient: telebot, attachmentsDir: attachmentsDir}

		telebot.On(method, userID, filepath.Join(attachmentsDir, attachment), text, mock.AnythingOfType("func(url.Values)")).Return(err)

		a.sendMessage(userID, text, emptyKeyboardOptions, emptyInlineOptions, attachment)

		telebot.AssertExpectations(t)
	}

	//file forwarding with file type and keyboard options
	testCases = map[string]string{"id:photo": "ForwardPhoto", "id:video": "ForwardVideo", "id:audio": "ForwardAudio", "id:doc": "ForwardFile"}
	for attachment, method := range testCases {
		telebot := &mocks.Telebot{}
		a := &application{tgClient: telebot, attachmentsDir: attachmentsDir}

		telebot.On(method, userID, strings.Split(attachment, ":")[0], text, mock.AnythingOfType("func(url.Values)")).Return(err)

		a.sendMessage(userID, text, keyboardOptions, emptyInlineOptions, attachment)

		telebot.AssertExpectations(t)
	}

	//file forwarding with file type and inline options
	testCases = map[string]string{"id:photo": "ForwardPhoto", "id:video": "ForwardVideo", "id:audio": "ForwardAudio", "id:doc": "ForwardFile"}
	for attachment, method := range testCases {
		telebot := &mocks.Telebot{}
		a := &application{tgClient: telebot, attachmentsDir: attachmentsDir}

		telebot.On(method, userID, strings.Split(attachment, ":")[0], text, mock.AnythingOfType("func(url.Values)")).Return(err)

		a.sendMessage(userID, text, emptyKeyboardOptions, inlineOptions, attachment)

		telebot.AssertExpectations(t)
	}

	//file forwarding with file type
	testCases = map[string]string{"id:photo": "ForwardPhoto", "id:video": "ForwardVideo", "id:audio": "ForwardAudio", "id:doc": "ForwardFile"}
	for attachment, method := range testCases {
		telebot := &mocks.Telebot{}
		a := &application{tgClient: telebot, attachmentsDir: attachmentsDir}

		telebot.On(method, userID, strings.Split(attachment, ":")[0], text, mock.AnythingOfType("func(url.Values)")).Return(err)

		a.sendMessage(userID, text, emptyKeyboardOptions, emptyInlineOptions, attachment)

		telebot.AssertExpectations(t)
	}

	//file forwarding with generic document and keyboard options
	telebot := &mocks.Telebot{}
	a := &application{tgClient: telebot, attachmentsDir: attachmentsDir}

	telebot.On("ForwardFile", userID, "id", text, mock.AnythingOfType("func(url.Values)")).Return(err)

	a.sendMessage(userID, text, keyboardOptions, emptyInlineOptions, "id")

	telebot.AssertExpectations(t)

	//file forwarding with generic document and keyboard options
	telebot = &mocks.Telebot{}
	a = &application{tgClient: telebot, attachmentsDir: attachmentsDir}

	telebot.On("ForwardFile", userID, "id", text, mock.AnythingOfType("func(url.Values)")).Return(err)

	a.sendMessage(userID, text, emptyKeyboardOptions, inlineOptions, "id")

	telebot.AssertExpectations(t)

	//file forwarding with generic document and keyboard options
	telebot = &mocks.Telebot{}
	a = &application{tgClient: telebot, attachmentsDir: attachmentsDir}

	telebot.On("ForwardFile", userID, "id", text, mock.AnythingOfType("func(url.Values)")).Return(err)

	a.sendMessage(userID, text, emptyKeyboardOptions, emptyInlineOptions, "id")

	telebot.AssertExpectations(t)

	//send text with keyboard options
	telebot = &mocks.Telebot{}
	a = &application{tgClient: telebot, attachmentsDir: attachmentsDir}

	telebot.On("SendText", userID, text, mock.AnythingOfType("func(url.Values)")).Return(err)

	a.sendMessage(userID, text, keyboardOptions, emptyInlineOptions, "")

	telebot.AssertExpectations(t)

	//send text with inline options
	telebot = &mocks.Telebot{}
	a = &application{tgClient: telebot, attachmentsDir: attachmentsDir}
	inlineOptions = []map[string]interface{}{{"a": "http://www.kg"}}

	telebot.On("SendText", userID, text, mock.AnythingOfType("func(url.Values)")).Return(err)

	a.sendMessage(userID, text, emptyKeyboardOptions, inlineOptions, "")

	telebot.AssertExpectations(t)

	//send text
	telebot = &mocks.Telebot{}
	a = &application{tgClient: telebot, attachmentsDir: attachmentsDir}

	telebot.On("SendText", userID, text, mock.AnythingOfType("func(url.Values)")).Return(err)

	a.sendMessage(userID, text, emptyKeyboardOptions, emptyInlineOptions, "")

	telebot.AssertExpectations(t)

	//ignore empty message
	telebot = &mocks.Telebot{}
	a = &application{tgClient: telebot, attachmentsDir: attachmentsDir}

	a.sendMessage(userID, "", emptyKeyboardOptions, emptyInlineOptions, "")

	if len(telebot.Calls) != 0 {
		t.Errorf("Expected 0 but got %d calls", len(telebot.Calls))
	}

}

func TestExecDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	a := &application{dbClient: db}

	mock.ExpectExec("update table set status=1").WillReturnError(err)

	a.ExecDB("update table set status=1", nil)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestQueryDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	a := &application{dbClient: db}

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "tom").
		AddRow(2, "jerry")

	mock.ExpectQuery("select id, name from user").WillReturnRows(rows)

	res := a.QueryDB("select id, name from user", nil)

	assert.Equal(t, 2, len(res))

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}
