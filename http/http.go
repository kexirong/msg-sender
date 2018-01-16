package http

import (
	"errors"
	"strconv"

	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/kexirong/msg-sender/email"
	"github.com/kexirong/msg-sender/wechat"
)

var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

type logfile struct {
	file        *os.File
	isOpen      bool
	rolte       <-chan time.Time
	preFileName string
}

func (lf *logfile) Write(b []byte) (int, error) {
	select {
	case <-lf.rolte:
		if lf.isOpen {
			lf.Close()
		}
		err := lf.SetFile(lf.preFileName + time.Now().Format("2006-01-02") + ".log")
		if err != nil {
			return 0, err
		}

	default:
		if !lf.isOpen {
			return 0, errors.New("logfile is not config")
		}
	}

	return lf.file.Write(b)
}

func NewLogFile(preFileName string) *logfile {
	var lf logfile
	lf.preFileName = preFileName
	err := lf.SetFile(lf.preFileName + time.Now().Format("2006-01-02") + ".log")
	if err != nil {
		panic(err)
	}

	go func() {
		t := time.Now()
		nx := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).AddDate(0, 0, 1)
		<-time.After(nx.Sub(t))
		lf.rolte = time.Tick(time.Hour * 24)
	}()

	return &lf
}

func (lf *logfile) SetFile(filename string) error {
	var filepath string
	wd, err := os.Getwd()
	if err == nil {
		filepath = path.Join(wd, filename)
	} else {
		panic(err)
	}

	_, err = os.Stat(filepath)

	if err == nil || os.IsExist(err) {
		err := os.Rename(filepath, filepath+strconv.FormatInt(time.Now().Unix(), 10))
		if err != nil {
			return err
		}
	}

	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0664)
	if err != nil {
		return err
	}
	lf.file = file
	lf.isOpen = true
	return nil
}

func (lf *logfile) Close() {
	lf.file.Close()
}

func init() {
	fli := NewLogFile("msg-sender")
	Info = log.New(fli, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(fli, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(fli, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func SrvStart(cfg *simplejson.Json) {
	//fmt.Println(cfg)

	httpAddr := cfg.Get("http").Get("listen").MustString("")
	smtpCfg := cfg.Get("smtp")
	smtpAddr := smtpCfg.Get("address").MustString("")
	smtpUser := smtpCfg.Get("username").MustString("")
	smtpPass := smtpCfg.Get("password").MustString("")
	smtpAuthtype := smtpCfg.Get("authtype").MustString("")
	wxCfg := cfg.Get("wechat")
	wxCorpID := wxCfg.Get("CorpID").MustString("")
	wxAgentID := wxCfg.Get("AgentId").MustInt(0)
	wxSecret := wxCfg.Get("Secret").MustString("")

	Info.Println(fmt.Sprintf("httpAddr:%s", httpAddr))
	Info.Println(fmt.Sprintf("smtpAddr:%s,smtpUser:%s,smtpPass:%s", smtpAddr, smtpUser, smtpPass))

	s := email.New(smtpAddr, smtpUser, smtpPass, smtpAuthtype)
	wx := wechat.New(wxCorpID, wxAgentID, wxSecret)
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
		tos := r.PostFormValue("to")
		to := strings.Split(tos, ",")
		subject := r.PostFormValue("subject")
		content := r.PostFormValue("content")
		Info.Println("#sendMail# ", "client: ", r.RemoteAddr, "tos:", tos, "subject:", subject, "content:", content)
		err = s.SendMail(to, subject, content)
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
		tos := r.PostFormValue("to")
		content := r.PostFormValue("content")
		Info.Println("#sendWechat# ", "client: ", r.RemoteAddr, "tos:", tos, "content:", content)

		resp, err := wx.SendMsg(tos, "", content)
		if err != nil {
			Error.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			Info.Println(resp)
			http.Error(w, resp, http.StatusOK)
		}
	})

	log.Fatal(http.ListenAndServe(httpAddr, nil))

}
