package http
import (
	"fmt"
    "github.com/kexirong/msg-sender/email"
	"net/http"
    "github.com/bitly/go-simplejson"
	_ "net/http/pprof"
    "log"
   "strings"
)


func SrvStart( cfg *simplejson.Json) {
    //fmt.Println(cfg)
    
    httpAddr:=cfg.Get("http").Get("listen").MustString("")
    smtpCfg:=cfg.Get("smtp")
    smtpAddr:=smtpCfg.Get("address").MustString("")
    smtpUser:=smtpCfg.Get("username").MustString("")
    smtpPass:=smtpCfg.Get("password").MustString("")
    smtpAuthtype:=smtpCfg.Get("authtype").MustString("")
    
    fmt.Println(fmt.Sprintf("httpAddr:%s,smtpAddr:%s,smtpUser:%s,smtpPass:%s" ,httpAddr,smtpAddr,smtpUser,smtpPass))
    
    
    s := email.New(smtpAddr,smtpUser,smtpPass,smtpAuthtype)
    fmt.Println(s)
    
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
       /* fmt.Println(r.URL.Query())
        fmt.Println(r.Method)
        err := r.ParseForm()
		if err != nil {
			panic(err)
		}
        fmt.Println(r.Form.Get("name"))
        fmt.Println(r.Form["name"])*/
		w.Write([]byte("ok"))
	})
    
    http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
       
		w.Write([]byte(testStr))
	})
    
    
    http.HandleFunc("/sender/mail", func(w http.ResponseWriter, r *http.Request) {
    
       /*
        err := r.ParseForm()
		if err != nil {
			panic(err)
		}*/
       //  fmt.Println(r.PostForm)
        tos:=r.PostFormValue("to")
        to := strings.Split(tos,";")
        subject := r.PostFormValue("subject")
        content := r.PostFormValue("content")
        fmt.Println("tos:",tos, "subject:",subject, "content:",content)
        err := s.SendMail( to, subject, content)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        } else {
            http.Error(w, "success", http.StatusOK)
        }
    })
    
    log.Fatal(http.ListenAndServe(httpAddr, nil))
    
}