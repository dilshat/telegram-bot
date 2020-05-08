if (timer) {
    console.log("Hi " + new Date())
} else {

    var step = get("step")

    if (!step) {
        step = { id: 0 }
        set("step", step)
    }

    var name = message.From.FirstName + ' ' + message.From.LastName
    if (name.length == 1) {
        name = message.From.Username
    }
    var userID = message.Chat.ID
    if (!callback) {
        console.log("User " + name + " ID " + userID + " sent " + message.Text)
    } else {
        name = callback.From.FirstName + ' ' + callback.From.LastName
        if (name.length == 1) {
            name = callback.From.Username
        }
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

}