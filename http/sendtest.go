package http
const testStr = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>test</title>
</head>

<body>

<div>


<form action="/sender/mail" method="post">
收件人:
<input type="text" name="to" value="kexirong@outlook.com">
<br>
标题:
<input type="text" name="subject" value="test">
<br>

内容:
<br>
  <textarea name="content" cols="30" rows="4"> 
        正文
  </textarea>  

<br><br>
<input type="submit" value="提交">
</form> 

<div>


</body>
</html>
  `