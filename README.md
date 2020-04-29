## A simple javascript telegram bot

Sometimes a simple bot is needed to be built quickly. This bot allows to write logic of bot in javascript. The script can use embedded methods to implement various functionality: sending messages with attachments, custom keyboards and inline options. You can store objecs in cache and in database (mysql and postgres are supported). You can code conversational bot using session like approach. See examples in `scripts/*.js`. File `lib.js` is used to store some helper functionality.

### How to use:

+ In file `.env` set TELEGRAM_TOKEN to your bot token
+ Implement bot logic in `attacments/logic.js`. It contains an example script to be used as reference
+ Build bot and run `./scriptable-bot`