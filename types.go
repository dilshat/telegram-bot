package main

import (
	"net/url"

	"github.com/ReneKroon/ttlcache"
	"github.com/robertkrimen/otto"
	"github.com/yanzay/tbot/v2"
)

type application struct {
	client         Telebot
	cache          *ttlcache.Cache
	attachmentsDir string
	logicScript    string
	token          string
	vmFactory      VmFactory
}

type Vm interface {
	Set(name string, value interface{}) error
	Run(src interface{}) (otto.Value, error)
}

type VmWrapper struct {
	vm *otto.Otto
}

func (v VmWrapper) Set(name string, value interface{}) error {
	return v.vm.Set(name, value)
}

func (v VmWrapper) Run(src interface{}) (otto.Value, error) {
	return v.vm.Run(src)
}

type VmFactory interface {
	GetVm() Vm
}

type VmFactoryImpl struct {
}

func (v VmFactoryImpl) GetVm() Vm {
	return &VmWrapper{vm: otto.New()}
}

type Telebot interface {
	GetFileInfo(fileID string) (*tbot.File, error)
	AnswerCallback(callbackQueryID string) error
	EditInlineMarkup(chatID string, messageID int, markup *tbot.InlineKeyboardMarkup) error
	AttachPhoto(chatID string, filename string, text string, option func(r url.Values)) error
	AttachVideo(chatID string, filename string, text string, option func(r url.Values)) error
	AttachAudio(chatID string, filename string, text string, option func(r url.Values)) error
	AttachFile(chatID string, filename string, text string, option func(r url.Values)) error
	ForwardPhoto(chatID string, fileID string, text string, option func(r url.Values)) error
	ForwardVideo(chatID string, fileID string, text string, option func(r url.Values)) error
	ForwardAudio(chatID string, fileID string, text string, option func(r url.Values)) error
	ForwardFile(chatID string, fileID string, text string, option func(r url.Values)) error
	SendText(chatID string, text string, option func(r url.Values)) error
}

type TbotWrapper struct {
	*tbot.Client
}

func (t *TbotWrapper) AnswerCallback(callbackQueryID string) error {
	return t.AnswerCallbackQuery(callbackQueryID)
}

func (t *TbotWrapper) GetFileInfo(fileID string) (*tbot.File, error) {
	return t.GetFile(fileID)
}

func (t *TbotWrapper) EditInlineMarkup(chatID string, messageID int, markup *tbot.InlineKeyboardMarkup) error {
	_, err := t.EditMessageReplyMarkup(chatID, messageID, tbot.OptInlineKeyboardMarkup(markup))
	return err
}

func (t *TbotWrapper) AttachPhoto(chatID string, filename string, text string, option func(r url.Values)) error {
	_, err := t.SendPhotoFile(chatID, filename, tbot.OptCaption(text), tbot.OptParseModeHTML, option)
	return err
}

func (t *TbotWrapper) AttachVideo(chatID string, filename string, text string, option func(r url.Values)) error {
	_, err := t.SendVideoFile(chatID, filename, tbot.OptCaption(text), tbot.OptParseModeHTML, option)
	return err
}

func (t *TbotWrapper) AttachAudio(chatID string, filename string, text string, option func(r url.Values)) error {
	_, err := t.SendAudioFile(chatID, filename, tbot.OptCaption(text), tbot.OptParseModeHTML, option)
	return err
}

func (t *TbotWrapper) AttachFile(chatID string, filename string, text string, option func(r url.Values)) error {
	_, err := t.SendDocumentFile(chatID, filename, tbot.OptCaption(text), tbot.OptParseModeHTML, option)
	return err
}

func (t *TbotWrapper) ForwardPhoto(chatID string, fileID string, text string, option func(r url.Values)) error {
	_, err := t.SendPhoto(chatID, fileID, tbot.OptCaption(text), tbot.OptParseModeHTML, option)
	return err
}

func (t *TbotWrapper) ForwardVideo(chatID string, fileID string, text string, option func(r url.Values)) error {
	_, err := t.SendVideo(chatID, fileID, tbot.OptCaption(text), tbot.OptParseModeHTML, option)
	return err
}

func (t *TbotWrapper) ForwardAudio(chatID string, fileID string, text string, option func(r url.Values)) error {
	_, err := t.SendAudio(chatID, fileID, tbot.OptCaption(text), tbot.OptParseModeHTML, option)
	return err
}

func (t *TbotWrapper) ForwardFile(chatID string, fileID string, text string, option func(r url.Values)) error {
	_, err := t.SendDocument(chatID, fileID, tbot.OptCaption(text), tbot.OptParseModeHTML, option)
	return err
}

func (t *TbotWrapper) SendText(chatID string, text string, option func(r url.Values)) error {
	_, err := t.SendMessage(chatID, text, tbot.OptParseModeHTML, option)
	return err
}
