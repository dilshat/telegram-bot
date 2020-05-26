
bot = {
    onMessage: function (message) {
        process(message)
    },
    onCallback: function (callback) {
    },
    onTimer: function () {
    },
    onInit: function () {
    }
}

function process(message) {
    //works for message.Audio, message.Video, message.Photo and message.Document
    if (message.Photo && message.Photo.length > 0) {
        //photo arrives as array of files of different size
        var fileID = message.Photo[message.Photo.length - 1].FileID
        console.log("Photo received:" + getFileLink(fileID))
        console.log("Echoing fileback to sender...")
        send("You sent this file:", null, fileID + ":photo")
    } else {
        send("Please, send a photo")
    }
}