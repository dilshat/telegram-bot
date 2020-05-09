
bot = {
    onMessage: function (message) {
        //This method is called when user sends message. User is available thru `message.From` property
        console.log(message.From.ID)
        send("Please press a button", [{ "One": "option-1", "Two": "option-2"}, {"This is a link, not option": "http://www.kg" }])
    },
    onCallback: function (callback) {
        //This method is called when user presses inline keyboard button. User is available thru `callback.From` property
        //Original message is available thru `callback.Messsage` property
        console.log(callback.From.ID)
        console.log("Option " + callback.Data + " selected")
        console.log("Original message: " + callback.Message.Text)
    },
    onTimer: function () {
        //This method is called regularly as configured by TIMER env var. When using send/prompt methods, you have to explicilty specify target user id
        console.log("OnTimer method " + new Date())
    },
    onInit: function () {
        //This method is called once in the beginning. It is supposed to be used for initialization purposes like a db setup. You can not use send/prompt methods here
        console.log("OnInit method")
    }
}
