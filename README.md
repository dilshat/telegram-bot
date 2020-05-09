[![Build Status](https://travis-ci.com/Dilshat/telegram-bot.svg?branch=master)](https://travis-ci.com/Dilshat/telegram-bot)
[![Go Report Card](https://goreportcard.com/badge/github.com/Dilshat/telegram-bot)](https://goreportcard.com/report/github.com/Dilshat/telegram-bot)

## A simple javascript telegram bot

Sometimes a simple bot is needed to be built quickly. This bot allows to write logic of bot in javascript. The script can use embedded methods to implement various functionality: sending messages with attachments, custom keyboards and inline options. You can store objects in cache and in database (mysql and postgres are supported). You can code conversational bot using session like approach. See examples in `scripts/*.js`.

### Embedded functions:

**set(key, value)** - puts a json object to cache indexed by key
```
var sess = {step:0}
set("session", sess)
```


**get(key)** - returns a json object from cache by key
```
var sess = get("session")
if (sess) {
  console.log(sess.step) 
}
```


**del(key)** - removes a json object from cache by key
```
del("session")
```


**send(...)** - sends a message to user
```
send("Hi") // sends a simple text message

send("Choose one", [["option1", "option2"]]) // sends a text message with custom keyboard

send("Choose one", [["option1"], ["option2"]]) // sends a text message with custom keyboard with button groups each on a new line

send("Download this file", null, "the-file.ext") // sends a text message with a file from attachments directory

send("Forwarded this file", null, "{fileID}:photo") // forwards a file, that is already uploaded to the chat by its fileID and type (photo, video, audio) separated by colon

send("Forwarded this file", null, "{fileID}") // forwards a file, that is already uploaded to the chat by its fileID as a generic document

send("hi Admin", null, null, adminId) // sends a message to the specified user (by telegram id)
```

**send("Test", [{ "One": "option-1", "Two": "option-2", "Three": "option-3" }])** - sends a message with inline keyboard.
 When user presses a button, script will have access to a callback object

```
send("Test", [{ "One": "option-1", "Two": "option-2" }])

send("Test", [{ "One": "option-1"}, {"Two": "option-2" }]) // button are displayed in groups, each on a new line

 ...

//process user click on a button
console.log("Option " + callback.Data + " selected")

```
 _if button data is a valid URL, clicking the button will not trigger callback but rather attempt to navigate the specified url_


**getFileLink(fileID)** - returns link to file download by its fileID. It is no recommended to share the link with users, since it contains bot token. Is is supposed to be used by bot admins
```
if (message.Photo && message.Photo.length > 0) {
  console.log("Photo received:" + getFileLink(message.Photo[message.Photo.length - 1].FileID))
}
```


**replaceOptions** - replaces inline keyboard
```
replaceOptions(message.Chat.ID, message.MessageID, [{ "Three": "option-3", "Four": "option-4" }] )
```


**editMessage** - updates message text and inline keyboard. _Currently it is possible to edit messages with inline keyboard only_
```
var id = send("Original message", [{ "One": "option-1", "Two": "option-2", "Three": "option-3" }])
sleep(1000)
editMessage(message.Chat.ID, id, "Edited message", [{ "Three": "option-3", "Four": "option-4" }])

```


**deleteMessage** - deletes message
```
deleteMessage(message.Chat.ID, message.MessageID)
```

**prompt(text, attachment, userId)** - sends message prompting user to reply to it (force reply)
```
prompt("What is your phone number") 

...

//process user reply
if (message.ReplyToMessage) {
  console.log(message.ReplyToMessage.Text + " -> " + message.Text)
}
```



### How to use:

+ In file `.env` set TELEGRAM_TOKEN to your bot token
+ Implement logic in `scripts/*logic*.js`.
+ Build bot `go build` and run `./telegram-bot`
