/*
  Please use inside `bot` object methods: onMessage, onCallback, onTimer, onInit
*/

dbExec("insert into users(name, phone) values('John Black', '123123123');")

var res = dbQuery("select * from users")

var users = JSON.stringify(res)

console.log(users[0].name)

send("Hello " + users[0].name)