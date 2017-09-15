#### 请求URL:

- http://ip:port/sender/mail

#### 请求方式：

- POST

#### 请求头：

|参数名|是否必须|类型|说明|
|:----    |:---|:----- |-----   |
|Content-Type |否  |string |请求类型： text/plain   |



#### 请求参数:

|参数名|是否必须|类型|说明|
|:----    |:---|:----- |----- ---  |
|to|是  |string | 收件人地址，多个收件人用(,)分隔|
|subject|是  |string | 邮件标题|
|content|是  |string | 邮件内容|
|contentType|否  |string |填 html|

