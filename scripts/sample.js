
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


function process(message, callback) {

    var userID = getUserID(message, callback)
    var userName = getUserName(message, callback)

    var step = get("step")
    if (!step) {
        step = { id: 0 }
        set("step", step)
    }

    if (step.id == 0) {
        send("Hello " + userName + "! Do you like this picture?", [["Yes", "No"], ["Maybe", "Hard to say"]], "smile.jpg")
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
            console.log("User " + userName + " ID " + userID + " replied to question " + message.ReplyToMessage.Text + " with " + message.Text)
        }
        send("Bye!")
        del("step")
    }

    //log
    if (!callback) {
        console.log("User " + userName + " ID " + userID + " sent " + message.Text)
    } else {
        console.log("User " + userName + " ID " + userID + " selected option " + callback.Data)
    }

}