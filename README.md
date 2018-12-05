# msg-sender

- ![邮件API](https://github.com/kexirong/msg-sender/blob/master/email/README.md)

- ![企业微信API](https://github.com/kexirong/msg-sender/blob/master/wechat/README.md)

### cfg.json 是纯json格式


```
{
    "debug": true,
    "http": {
        "listen": "0.0.0.0:4000", //监听ip端口
        "allow":["*"],// 填写ip，"*" 代表允许全部
        "deny":[]
    },
    "smtp": {//邮件配置
        "address": "smtp.exmail.qq.com:25",//邮件发送服务器地址
        "username": "kexirong@example.com",
        "password": "123456"
    },
    "wechat":{//企业微信配置
        "CorpID":"ww2085a342", //企业ID
        "AgentId":1000002,//应用id，通过新建企业微信应用获取
        "Secret":"5WsjwD2DqyR4PMTWnJJp_qvyOothRjDAZsaKc"//密串，企业微信应用中可以得到
    }
}
- curl -d "to=test@qq.com,test@sina.com&subject=test&content=test测试..." "http://10.1.1.202:4000/sender/mail"
- curl -d "to=kexirong&content=test测试..." "http://10.1.1.202:4000/sender/wechat"

```

- curl -d "to=test@qq.com,test@sina.com&subject=test&content=test测试..." "http://10.1.1.202:4000/sender/mail"
- curl -d "to=kexirong&&content=test测试..." "http://10.1.1.202:4000/sender/wechat"

```


echo "# msgsender" >> README.md
git init
git add README.md
git commit -m "first commit"
git remote add origin https://github.com/kexirong/msg-sender.git
git push -u origin master

```
