
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
  //var res = bExec("insert into users(name, phone, birth_date) values($1,$2,$3)", 'James Bond', '996777123456', '1981-04-25')
  //var result = JSON.parse(res)
  //console.log(result.lastInsertId) // valid only for mysql db
  
  var res = dbExec("update users set name=$1 where id=$2", 'Jason Bourne', 1)
  var result = JSON.parse(res)
  console.log(result.rowsAffected)

  var json = dbQuery("select id, name, phone, to_char(birth_date, 'DD-MM-YYYY HH24:MI:SS') as bd from users")
  console.log(json)

  var users = JSON.parse(json)
  console.log(users[0].name)

  //dbReport("users", "Here is a list of users", null, "select name, phone, to_char(birth_date, 'DD-MM-YYYY HH24:MI:SS') as bd from users where id = $1", 1)

}