package main

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/robertkrimen/otto"
	_ "github.com/robertkrimen/otto/underscore"
	"github.com/yanzay/tbot/v2"
)

func (a *application) setCacheItem(key string, val interface{}) {
	a.cache.Set(key, val)
}

func (a *application) getCacheItem(key string) interface{} {
	val, _ := a.cache.Get(key)
	return val
}

func (a *application) delCacheItem(key string) {
	a.cache.Remove(key)
}

func (a *application) getFileLink(fileID string) string {
	file, err := a.tgClient.GetFileInfo(fileID)
	if err != nil {
		log.Error("Error getting file link ", err)
		return ""
	}
	return fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", a.token, file.FilePath)
}

func (a *application) messageHandler(m *tbot.Message) {
	a.handleMessage(m, nil)
}

func (a *application) callbackHandler(cq *tbot.CallbackQuery) {
	a.tgClient.AnswerCallback(cq.ID)

	a.handleMessage(cq.Message, cq)
}

func (a *application) replaceInlineOptions(chatID string, msgID int, inlineOptions []map[string]interface{}) int {
	id, err := a.tgClient.EditInlineMarkup(chatID, msgID, buildInlineOptions(inlineOptions))
	if err != nil {
		log.Error("Error replacing inline options ", err)
	}
	return id
}

func (a *application) deleteMessage(chatID string, msgID int) {
	err := a.tgClient.DeleteMsg(chatID, msgID)
	if err != nil {
		log.Error("Error deleting message ", err)
	}
}

func (a *application) editMessage(chatID string, msgID int, text string, inlineOptions []map[string]interface{}) {
	err := a.tgClient.EditMsg(chatID, msgID, text, buildInlineOptions(inlineOptions))
	if err != nil {
		log.Error("Error editing message ", err)
	}
}

func (a *application) doGet(aURL string, params map[string]interface{}, headers map[string]interface{}) string {
	resp, err := doGET(aURL, params, headers)
	if err != nil {
		log.Error("Error performing GET request ", err)
	}

	return resp
}

func (a *application) doPOST(aURL string, params map[string]interface{}, headers map[string]interface{}) string {
	resp, err := doPOST(aURL, params, headers)
	if err != nil {
		log.Error("Error performing POST request ", err)
	}

	return resp
}
func (a *application) QueryDB(query string) []map[string]interface{} {
	result := []map[string]interface{}{}
	if a.useDB {
		rows, err := a.dbClient.Query(query)
		if err != nil {
			log.Error("Error querying db ", err)
			return result
		}
		defer rows.Close()

		columnTypes, err := rows.ColumnTypes()

		if err != nil {
			log.Error("Error querying db ", err)
			return result
		}

		count := len(columnTypes)

		for rows.Next() {

			scanArgs := make([]interface{}, count)

			for i, v := range columnTypes {

				switch v.DatabaseTypeName() {
				case "VARCHAR", "TEXT", "UUID", "TIMESTAMP":
					scanArgs[i] = new(sql.NullString)
					break
				case "BOOL":
					scanArgs[i] = new(sql.NullBool)
					break
				case "INT4":
					scanArgs[i] = new(sql.NullInt64)
					break
				default:
					scanArgs[i] = new(sql.NullString)
				}
			}

			err := rows.Scan(scanArgs...)

			if err != nil {
				log.Error("Error querying db ", err)
				return result
			}

			masterData := map[string]interface{}{}

			for i, v := range columnTypes {

				if z, ok := (scanArgs[i]).(*sql.NullBool); ok {
					masterData[v.Name()] = z.Bool
					continue
				}

				if z, ok := (scanArgs[i]).(*sql.NullString); ok {
					masterData[v.Name()] = z.String
					continue
				}

				if z, ok := (scanArgs[i]).(*sql.NullInt64); ok {
					masterData[v.Name()] = z.Int64
					continue
				}

				if z, ok := (scanArgs[i]).(*sql.NullFloat64); ok {
					masterData[v.Name()] = z.Float64
					continue
				}

				if z, ok := (scanArgs[i]).(*sql.NullInt32); ok {
					masterData[v.Name()] = z.Int32
					continue
				}

				masterData[v.Name()] = scanArgs[i]
			}

			result = append(result, masterData)
		}

		// z, err := json.Marshal(finalRows)
	}
	return result
}

func (a *application) ExecDB(query string) {
	if a.useDB {
		if _, err := a.dbClient.Exec(query); err != nil {
			log.Error("Error executing db query ", err)
		}
	}
}

func (a *application) handleMessage(m *tbot.Message, cq *tbot.CallbackQuery) {

	vm := a.vmFactory.GetVm()

	vm.Set("doGet", a.getDoGetFunc())

	vm.Set("doPost", a.getDoPostFunc())

	vm.Set("dbQuery", a.getQueryDBFunc())

	vm.Set("dbExec", a.getExecDBFunc())

	vm.Set("getFileLink", a.getGetFileLinkFunc())

	vm.Set("replaceOptions", a.getReplaceOptionsFunc())

	vm.Set("deleteMessage", a.getDeleteMessageFunc())

	vm.Set("editMessage", a.getEditMessageFunc())

	vm.Set("sleep", a.getSleepFunc())

	id := ""
	if m != nil {
		id = m.Chat.ID
		vm.Set("timer", false)
	} else {
		vm.Set("timer", true)
	}

	vm.Set("send", a.getSendFunc(id))

	vm.Set("prompt", a.getPromptFunc(id))

	vm.Set("set", a.getSetFunc(id))

	vm.Set("get", a.getGetFunc(id))

	vm.Set("del", a.getDelFunc(id))

	vm.Set("callback", cq)

	vm.Set("message", m)

	_, err := vm.Run(a.logicScript)

	if err != nil {
		log.Error("Error executing script ", err)
	}
}

func (a *application) getQueryDBFunc() func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		result := otto.Value{}

		if query, err := call.Argument(0).ToString(); err == nil {
			result, _ = otto.New().ToValue(a.QueryDB(query))
		}

		return result
	}
}

func (a *application) getExecDBFunc() func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		if query, err := call.Argument(0).ToString(); err == nil {
			a.ExecDB(query)
		}

		return otto.Value{}
	}
}

func (a *application) getSleepFunc() func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		if ms, err := call.Argument(0).ToInteger(); err == nil {
			time.Sleep(time.Duration(ms) * time.Millisecond)
		}

		return otto.Value{}
	}
}

func (a *application) getDoGetFunc() func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		result := otto.Value{}

		if aURL, err := call.Argument(0).ToString(); err == nil {
			if paramsInterface, err := call.Argument(1).Export(); err == nil {
				if params, ok := paramsInterface.(map[string]interface{}); ok {
					var headers map[string]interface{}
					if headersInterface, err := call.Argument(2).Export(); err == nil {
						headers, _ = headersInterface.(map[string]interface{})
					}
					result, _ = otto.ToValue(a.doGet(aURL, params, headers))
				}
			}
		}

		return result
	}
}

func (a *application) getDoPostFunc() func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		result := otto.Value{}

		if aURL, err := call.Argument(0).ToString(); err == nil {
			if paramsInterface, err := call.Argument(1).Export(); err == nil {
				if params, ok := paramsInterface.(map[string]interface{}); ok {
					var headers map[string]interface{}
					if headersInterface, err := call.Argument(2).Export(); err == nil {
						headers, _ = headersInterface.(map[string]interface{})
					}
					result, _ = otto.ToValue(a.doPOST(aURL, params, headers))
				}
			}
		}

		return result
	}
}

func (a *application) getReplaceOptionsFunc() func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		if chatID, err := call.Argument(0).ToString(); err == nil {
			if msgID, err := call.Argument(1).ToInteger(); err == nil {
				if optionsInterface, err := call.Argument(2).Export(); err == nil {
					if inlineOptions, ok := optionsInterface.([]map[string]interface{}); ok {
						a.replaceInlineOptions(chatID, int(msgID), inlineOptions)
					}
				}
			}
		}

		return otto.Value{}
	}
}

func (a *application) getDeleteMessageFunc() func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		if chatID, err := call.Argument(0).ToString(); err == nil {
			if msgID, err := call.Argument(1).ToInteger(); err == nil {
				a.deleteMessage(chatID, int(msgID))
			}
		}

		return otto.Value{}
	}
}

func (a *application) getEditMessageFunc() func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		if chatID, err := call.Argument(0).ToString(); err == nil {
			if msgID, err := call.Argument(1).ToInteger(); err == nil {
				if text, err := call.Argument(2).ToString(); err == nil {
					if optionsInterface, err := call.Argument(3).Export(); err == nil {
						if inlineOptions, ok := optionsInterface.([]map[string]interface{}); ok {
							a.editMessage(chatID, int(msgID), text, inlineOptions)
						}
					}
				}
			}
		}

		return otto.Value{}
	}
}

func (a *application) getGetFileLinkFunc() func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		result := otto.Value{}
		if call.Argument(0).IsString() {
			fileID, _ := call.Argument(0).ToString()
			result, _ = otto.ToValue(a.getFileLink(fileID))
		}

		return result
	}
}

func (a *application) getSetFunc(userID string) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		if call.Argument(0).IsString() {
			key, _ := call.Argument(0).ToString()
			if call.Argument(1).IsObject() {
				val := call.Argument(1).Object()
				a.setCacheItem(fmt.Sprintf("%s_%s", userID, key), val)
			}
		}

		return otto.Value{}
	}
}

func (a *application) getGetFunc(userID string) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		result := otto.Value{}
		if call.Argument(0).IsString() {
			key, _ := call.Argument(0).ToString()
			result, _ = otto.ToValue(a.getCacheItem(fmt.Sprintf("%s_%s", userID, key)))
		}

		return result
	}
}

func (a *application) getDelFunc(userID string) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		if call.Argument(0).IsString() {
			key, _ := call.Argument(0).ToString()
			a.delCacheItem(fmt.Sprintf("%s_%s", userID, key))
		}

		return otto.Value{}
	}
}

func (a *application) getPromptFunc(userID string) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		var text, attachment, targetUser string

		if call.Argument(0).IsString() {
			text, _ = call.Argument(0).ToString()
			text = strings.TrimSpace(text)
		}
		if call.Argument(1).IsString() {
			attachment, _ = call.Argument(1).ToString()
			attachment = strings.TrimSpace(attachment)
		}

		targetUser = userID
		if call.Argument(2).IsDefined() {
			targetUser, _ = call.Argument(2).ToString()
		}

		id := a.promptUser(targetUser, text, attachment)

		result, _ := otto.ToValue(id)

		return result
	}
}

func (a *application) getSendFunc(userID string) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		var text, attachment, targetUser string
		var options [][]string
		var inlineOptions []map[string]interface{}

		if call.Argument(0).IsString() {
			text, _ = call.Argument(0).ToString()
			text = strings.TrimSpace(text)
		}
		optionsInterface, err := call.Argument(1).Export()

		if err == nil {
			var ok bool
			if options, ok = optionsInterface.([][]string); !ok {
				inlineOptions, _ = optionsInterface.([]map[string]interface{})
			}
		}

		if call.Argument(2).IsString() {
			attachment, _ = call.Argument(2).ToString()
			attachment = strings.TrimSpace(attachment)
		}

		targetUser = userID
		if call.Argument(3).IsDefined() {
			targetUser, _ = call.Argument(3).ToString()
		}

		id := a.sendMessage(targetUser, text, options, inlineOptions, attachment)

		result, _ := otto.ToValue(id)

		return result
	}
}

func (a *application) promptUser(userID string, text string, attachment string) int {

	defer func() {
		if r := recover(); r != nil {
			log.Error("Recovered in promptUser ", r)
		}
	}()

	attachmentFile := filepath.Join(a.attachmentsDir, attachment)
	hasAttachment := attachment != "" && FileExists(attachmentFile)

	var id int
	var err error

	if hasAttachment {
		fileType := GetFileType(attachmentFile)
		if fileType == PHOTO {
			id, err = a.tgClient.AttachPhoto(userID, attachmentFile, text, tbot.OptForceReply)
		} else if fileType == VIDEO {
			id, err = a.tgClient.AttachVideo(userID, attachmentFile, text, tbot.OptForceReply)
		} else if fileType == AUDIO {
			id, err = a.tgClient.AttachAudio(userID, attachmentFile, text, tbot.OptForceReply)
		} else {
			id, err = a.tgClient.AttachFile(userID, attachmentFile, text, tbot.OptForceReply)
		}
	} else if attachment != "" {
		fileParts := strings.Split(attachment, ":")
		if len(fileParts) == 2 {
			fileType := ParseFileType(fileParts[1])
			if fileType == PHOTO {
				id, err = a.tgClient.ForwardPhoto(userID, fileParts[0], text, tbot.OptForceReply)
			} else if fileType == VIDEO {
				id, err = a.tgClient.ForwardVideo(userID, fileParts[0], text, tbot.OptForceReply)
			} else if fileType == AUDIO {
				id, err = a.tgClient.ForwardAudio(userID, fileParts[0], text, tbot.OptForceReply)
			} else {
				id, err = a.tgClient.ForwardFile(userID, fileParts[0], text, tbot.OptForceReply)
			}
		} else {
			id, err = a.tgClient.ForwardFile(userID, attachment, text, tbot.OptForceReply)
		}
	} else if strings.TrimSpace(text) != "" {
		id, err = a.tgClient.SendText(userID, text, tbot.OptForceReply)
	} else {
		log.Warn("Ignoring empty response")
	}

	if err != nil {
		log.Error("Error prompting user ", err)
	}

	return id
}

func (a *application) sendMessage(userID string, text string, options [][]string, inlineOptions []map[string]interface{}, attachment string) int {

	defer func() {
		if r := recover(); r != nil {
			log.Error("Recovered in sendMessage ", r)
		}
	}()

	attachmentFile := filepath.Join(a.attachmentsDir, attachment)
	hasAttachment := attachment != "" && FileExists(attachmentFile)
	hasOptions := len(options) > 0
	hasInlineOptions := len(inlineOptions) > 0

	var id int
	var err error

	if hasAttachment {
		//file uploading
		fileType := GetFileType(attachmentFile)
		if hasOptions {
			if fileType == PHOTO {
				id, err = a.tgClient.AttachPhoto(userID, attachmentFile, text, tbot.OptReplyKeyboardMarkup(
					buildReplyOptions(options),
				))
			} else if fileType == VIDEO {
				id, err = a.tgClient.AttachVideo(userID, attachmentFile, text, tbot.OptReplyKeyboardMarkup(
					buildReplyOptions(options),
				))
			} else if fileType == AUDIO {
				id, err = a.tgClient.AttachAudio(userID, attachmentFile, text, tbot.OptReplyKeyboardMarkup(
					buildReplyOptions(options),
				))
			} else {
				id, err = a.tgClient.AttachFile(userID, attachmentFile, text, tbot.OptReplyKeyboardMarkup(
					buildReplyOptions(options),
				))
			}

		} else if hasInlineOptions {
			if fileType == PHOTO {
				id, err = a.tgClient.AttachPhoto(userID, attachmentFile, text, tbot.OptInlineKeyboardMarkup(
					buildInlineOptions(inlineOptions),
				))
			} else if fileType == VIDEO {
				id, err = a.tgClient.AttachVideo(userID, attachmentFile, text, tbot.OptInlineKeyboardMarkup(
					buildInlineOptions(inlineOptions),
				))
			} else if fileType == AUDIO {
				id, err = a.tgClient.AttachAudio(userID, attachmentFile, text, tbot.OptInlineKeyboardMarkup(
					buildInlineOptions(inlineOptions),
				))
			} else {
				id, err = a.tgClient.AttachFile(userID, attachmentFile, text, tbot.OptInlineKeyboardMarkup(
					buildInlineOptions(inlineOptions),
				))
			}
		} else {
			if fileType == PHOTO {
				id, err = a.tgClient.AttachPhoto(userID, attachmentFile, text, tbot.OptReplyKeyboardRemove)
			} else if fileType == VIDEO {
				id, err = a.tgClient.AttachVideo(userID, attachmentFile, text, tbot.OptReplyKeyboardRemove)
			} else if fileType == AUDIO {
				id, err = a.tgClient.AttachAudio(userID, attachmentFile, text, tbot.OptReplyKeyboardRemove)
			} else {
				id, err = a.tgClient.AttachFile(userID, attachmentFile, text, tbot.OptReplyKeyboardRemove)
			}
		}
	} else if attachment != "" {
		//file forwarding
		fileParts := strings.Split(attachment, ":")
		if len(fileParts) == 2 {
			//file type is specified
			fileType := ParseFileType(fileParts[1])
			if hasOptions {
				if fileType == PHOTO {
					id, err = a.tgClient.ForwardPhoto(userID, fileParts[0], text, tbot.OptReplyKeyboardMarkup(
						buildReplyOptions(options),
					))
				} else if fileType == VIDEO {
					id, err = a.tgClient.ForwardVideo(userID, fileParts[0], text, tbot.OptReplyKeyboardMarkup(
						buildReplyOptions(options),
					))
				} else if fileType == AUDIO {
					id, err = a.tgClient.ForwardAudio(userID, fileParts[0], text, tbot.OptReplyKeyboardMarkup(
						buildReplyOptions(options),
					))
				} else {
					id, err = a.tgClient.ForwardFile(userID, fileParts[0], text, tbot.OptReplyKeyboardMarkup(
						buildReplyOptions(options),
					))
				}
			} else if hasInlineOptions {
				if fileType == PHOTO {
					id, err = a.tgClient.ForwardPhoto(userID, fileParts[0], text, tbot.OptInlineKeyboardMarkup(
						buildInlineOptions(inlineOptions),
					))
				} else if fileType == VIDEO {
					id, err = a.tgClient.ForwardVideo(userID, fileParts[0], text, tbot.OptInlineKeyboardMarkup(
						buildInlineOptions(inlineOptions),
					))
				} else if fileType == AUDIO {
					id, err = a.tgClient.ForwardAudio(userID, fileParts[0], text, tbot.OptInlineKeyboardMarkup(
						buildInlineOptions(inlineOptions),
					))
				} else {
					id, err = a.tgClient.ForwardFile(userID, fileParts[0], text, tbot.OptInlineKeyboardMarkup(
						buildInlineOptions(inlineOptions),
					))
				}
			} else {
				if fileType == PHOTO {
					id, err = a.tgClient.ForwardPhoto(userID, fileParts[0], text, tbot.OptReplyKeyboardRemove)
				} else if fileType == VIDEO {
					id, err = a.tgClient.ForwardVideo(userID, fileParts[0], text, tbot.OptReplyKeyboardRemove)
				} else if fileType == AUDIO {
					id, err = a.tgClient.ForwardAudio(userID, fileParts[0], text, tbot.OptReplyKeyboardRemove)
				} else {
					id, err = a.tgClient.ForwardFile(userID, fileParts[0], text, tbot.OptReplyKeyboardRemove)
				}
			}
		} else {
			//send generic document
			if hasOptions {
				id, err = a.tgClient.ForwardFile(userID, attachment, text, tbot.OptReplyKeyboardMarkup(
					buildReplyOptions(options),
				))
			} else if hasInlineOptions {
				id, err = a.tgClient.ForwardFile(userID, attachment, text, tbot.OptInlineKeyboardMarkup(
					buildInlineOptions(inlineOptions),
				))
			} else {
				id, err = a.tgClient.ForwardFile(userID, attachment, text, tbot.OptReplyKeyboardRemove)
			}
		}
	} else if hasOptions {
		id, err = a.tgClient.SendText(
			userID,
			text,
			tbot.OptReplyKeyboardMarkup(
				buildReplyOptions(options),
			),
		)
	} else if hasInlineOptions {
		id, err = a.tgClient.SendText(
			userID,
			text,
			tbot.OptInlineKeyboardMarkup(
				buildInlineOptions(inlineOptions),
			),
		)
	} else if strings.TrimSpace(text) != "" {
		id, err = a.tgClient.SendText(userID, text, tbot.OptReplyKeyboardRemove)
	} else {
		log.Warn("Ignoring empty response")
	}

	if err != nil {
		log.Error("Error sending message ", err)
	}

	return id
}

func buildReplyOptions(replyOptions [][]string) *tbot.ReplyKeyboardMarkup {
	keyboard := make([][]tbot.KeyboardButton, len(replyOptions))
	for i := range replyOptions {
		keyboard[i] = make([]tbot.KeyboardButton, len(replyOptions[i]))
		for j := range replyOptions[i] {
			keyboard[i][j] = tbot.KeyboardButton{Text: replyOptions[i][j]}
		}
	}

	return &tbot.ReplyKeyboardMarkup{
		Keyboard:        keyboard,
		OneTimeKeyboard: true,
		ResizeKeyboard:  true,
	}

}

func buildInlineOptions(inlineOptions []map[string]interface{}) *tbot.InlineKeyboardMarkup {
	keyboard := make([][]tbot.InlineKeyboardButton, len(inlineOptions))

	for i := range inlineOptions {
		keyboard[i] = make([]tbot.InlineKeyboardButton, len(inlineOptions[i]))
		j := 0
		for key, val := range inlineOptions[i] {
			valStr := fmt.Sprintf("%s", val)
			if isValidUrl(valStr) {
				button := tbot.InlineKeyboardButton{
					Text: fmt.Sprintf("%s", key),
					URL:  valStr,
				}
				keyboard[i][j] = button
			} else {
				button := tbot.InlineKeyboardButton{
					Text:         fmt.Sprintf("%s", key),
					CallbackData: valStr,
				}
				keyboard[i][j] = button
			}
			j++
		}
	}

	return &tbot.InlineKeyboardMarkup{
		InlineKeyboard: keyboard,
	}
}
