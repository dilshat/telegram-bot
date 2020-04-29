/*-------------
Embedded functions:

* set(key, value) - puts a json object to cache indexed by key

ex:
  var sess = {step:0}
  set("session", sess)

 ----------
* get(key) - returns a json object from cache by key

ex:
  var sess = get("session")
  if (sess) {
     console.log(sess.step) 
  }

 ----------
* del(key) - remvoes a json object from cache by key
ex:
  del("session")
-----------
* send - sends message to user

ex:
  send("Hi") - a simple text message
  send("Choose one", [["option1", "option2"]]) - a text message with custom keyboard
    или send("Choose one", [["option1"], ["option2"]]) - a text message with custom keyboard with button groups each on a new line

  send("Download this file", null, "the-file.ext") - a text message with a file from attachments directory
  send("Forwarded this file", null, "{fileID}:photo") - forwards a file, that is already uploaded to the chat by its fileID and type (photo, video, audio) separated by colon
  send("Forwarded this file", null, "{fileID}") - forwards a file, that is already uploaded to the chat by its fileID as a generic document

  send("hi Admin", null, null, adminId) - a message to the specified user (telegram id)

----------
* getFileLink - returns link to file download by its fileID. It is no recommended to share the link with users, since it contains bot token. Is is supposed to be used by bot admins

ex:
  if (message.Photo && message.Photo.length > 0) {
      console.log("Photo received:" + getFileLink(message.Photo[message.Photo.length - 1].FileID))
  }

 ----------
* send("Test", [{ "One": "option-1", "Two": "option-2", "Three": "option-3" }]) - sends a message with inline keyboard.
 When user presses a buttons, script will have access to a callback object

ex:
  send("Test", [{ "One": "option-1", "Two": "option-2" }])
    or  send("Test", [{ "One": "option-1"}, {"Two": "option-2" }]) // button groups each on a new line
 ...
  if (callback) {
      console.log("Option " + callback.Data + " selected")
  }

 if button data is a valid URL, clicking the button will not trigger callback but rather attempt to navigate the specified url
 ----------
* replaceOptions - replaces inline keyboard

ex
  if (callback) {
      replaceOptions(message.Chat.ID, message.MessageID, [{ "Three": "option-3", "Four": "option-4" }] )
  }

 ----------
 * prompt(text, attachment, userId) - sends message prompting user to reply to it (force reply)

 ex:
  prompt("What is your phone number") 
  ...
  if (message.ReplyToMessage) {
    console.log(message.ReplyToMessage.Text + " -> " + message.Text)
  }

-----------------*/

//*** bot business logic

var step = get("step")

if (!step) {
  step = { id: 0 }
  set("step", step)
}

var name = message.From.FirstName ? message.From.FirstName : message.From.Username
var userID = message.Chat.ID
if (!callback) {
  console.log("User " + name + " ID " + userID + " sent " + message.Text)
} else {
  name = callback.From.FirstName ? callback.From.FirstName : callback.From.Username
  userID = callback.From.ID
  console.log("User " + name + " ID " + userID + " selected option " + callback.Data)
}

if (step.id == 0) {
  send("Hello " + name + "! Do you like this picture?", [["Yes", "No"], ["Maybe", "Hard to say"]], "smile.jpg")
  step.id++
  set("step", step)
} else if (step.id == 1) {
  send("how much will be 2 + 2?", [{ "3": "answer_three", "4": "answer_four" }, { "5": "answer_five", "Don't know": "answer_dunno" }])
  step.id++
  set("step", step)
} else if (step.id == 2) {
  prompt("Thanks, will we meet again?")
  step.id++
  set("step", step)
} else {
  if (message.ReplyToMessage) {
    console.log("User " + name + " ID " + userID + " replied to question " + message.ReplyToMessage.Text + " with " + message.Text)
  }
  send("Bye!")
  del("step")
}
