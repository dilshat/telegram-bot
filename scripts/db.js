
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

  dbExec("insert into users(name, phone, birth_date) values($1, $2, $3)", 'James Bond', '996222123321', '1981-04-25')

  var res = dbQuery("select * from users where name = $1", 'James Bond')

  var users = JSON.stringify(res)

  console.log(users[0].name)


  //dbReport("users", "Here is a list of users", null, "select name, phone, to_char(birth_date, 'DD-MM-YYYY HH24:MI:SS') as bd from users where id = $1", 1)

}