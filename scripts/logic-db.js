
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

  dbExec("insert into users(name, phone) values('John Black', '123123123');")

  var res = dbQuery("select * from users")

  var users = JSON.stringify(res)

  console.log(users[0].name)

  send("Hello " + users[0].name)

  // dbReport("select name, phone, to_char(birth_date, 'DD-MM-YYYY HH24:MI:SS') from users", "users", "Here is a list of users")

}