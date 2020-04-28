var resp1 = doGet("https://jsonplaceholder.typicode.com/comments", { postId: "1" })
resp1Json = JSON.parse(resp1)
console.log(resp1Json[0].email)
send("Hi " + resp1Json[0].name)


var resp2 = doPost("https://jsonplaceholder.typicode.com/posts", { title: 'foo', body: 'bar', userId: 1 })
var resp2Json = JSON.parse(resp2)
console.log(resp2)
send("Post with id " + resp2Json.id + " created")