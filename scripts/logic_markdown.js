
var testMsg = "<b>bold</b>, <strong>bold</strong>\n" +
    "<i>italic</i>, <em>italic</em>\n" +
    "<u>underline</u>, <ins>underline</ins>\n" +
    "<s>strikethrough</s>, <strike>strikethrough</strike>, <del>strikethrough</del>\n" +
    "<b>bold <i>italic bold <s>italic bold strikethrough</s> <u>underline italic bold</u></i> bold</b>\n" +
    "<a href='http://www.example.com/'>inline URL</a>\n" +
    "<a href='tg://user?id=123456789'>inline mention of a user</a>\n" +
    "<code>inline fixed-width code</code>\n" +
    "<pre>pre-formatted fixed-width code block</pre>\n" +
    "<pre><code class='language-python'>pre-formatted fixed-width code block written in the Python programming language</code></pre>";

send(testMsg)