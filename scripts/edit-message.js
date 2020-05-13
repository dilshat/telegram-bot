
bot = {
  onMessage: function (message) {
    process(message, null)
  },
  onCallback: function (callback) {
    process(callback.message, callback)
  },
  onTimer: function () {
  },
  onInit: function () {
  }
}

function process(message) {

  var id = send("Original message", [{ "One": "option-1", "Two": "option-2", "Three": "option-3" }])

  sleep(1000)

  editMessage(message.Chat.ID, id, "Edited message 1", [{}])

  sleep(3000)

  editMessage(message.Chat.ID, id, "Edited message 2", [{ "Three": "option-3", "Four": "option-4" }])
  
}