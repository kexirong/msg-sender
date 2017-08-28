package http
import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
    "log"
    "os"
    "strings"
    "github.com/kexirong/msg-sender/email"
    "github.com/kexirong/msg-sender/wechat"
    
    
    "github.com/bitly/go-simplejson"
)
var ( 
    Info *log.Logger 
    Warning *log.Logger 
    Error *log.Logger 
   
    ) 

func init(){
    Info = log.New( os.Stdout , "INFO: " , log.Ldate|log.Ltime|log.Lshortfile )
    Warning = log.New(os.Stdout,"WARNING: ",log.Ldate|log.Ltime|log.Lshortfile)
    Error = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}


func SrvStart( cfg *simplejson.Json) {
    //fmt.Println(cfg)
    
    httpAddr:=cfg.Get("http").Get("listen").MustString("")
    smtpCfg:=cfg.Get("smtp")
    smtpAddr:=smtpCfg.Get("address").MustString("")
    smtpUser:=smtpCfg.Get("username").MustString("")
    smtpPass:=smtpCfg.Get("password").MustString("")
    smtpAuthtype:=smtpCfg.Get("authtype").MustString("")
    wxCfg:=cfg.Get("wechat")
    wxCorpID:=wxCfg.Get("CorpID").MustString("")
    wxAgentId:=wxCfg.Get("AgentId").MustInt(0) 
    wxSecret:=wxCfg.Get("Secret").MustString("") 
    
    Info.Println(fmt.Sprintf("httpAddr:%s",httpAddr))
    Info.Println(fmt.Sprintf("smtpAddr:%s,smtpUser:%s,smtpPass:%s" ,smtpAddr, smtpUser, smtpPass))
    
    
    s := email.New(smtpAddr,smtpUser,smtpPass,smtpAuthtype)
    wx:=wechat.New(wxCorpID, wxAgentId, wxSecret)
    Info.Println(wx.GetAccToken())
    
    
    
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
        err := r.ParseForm()
		if err != nil {
			panic(err)
		}
        //fmt.Println(r.PostForm)
        tos:=r.PostFormValue("to")
        to := strings.Split(tos,",")
        subject := r.PostFormValue("subject")
        content := r.PostFormValue("content")
        Info.Println("tos:",tos, "subject:",subject, "content:",content)
        err = s.SendMail( to, subject, content)
        if err != nil {
            Error.Println(err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
        } else {
            http.Error(w, "success", http.StatusOK)
        }
    })
    
    http.HandleFunc("/sender/wechat", func(w http.ResponseWriter, r *http.Request) {
        err := r.ParseForm()
		if err != nil {
			panic(err)
		}
        //fmt.Println(r.PostForm)
        tos:=r.PostFormValue("to")
        content := r.PostFormValue("content")
        Info.Println("tos:",tos, "content:",content)
        
        resp,err := wx.SendMsg(tos, "", content)
        if err != nil {
            Error.Println(err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
        } else {
            Info.Println(resp)
            http.Error(w, resp , http.StatusOK)
        }
    })
    
    
    log.Fatal(http.ListenAndServe(httpAddr, nil))
    
}