package log
import (
  "os"
  "io"
  "log"
  "path"
  "strings"
)

var (
  debug *log.Logger
  infor *log.Logger
  warn *log.Logger
  broken *log.Logger
)

func init(){
  arg := path.Base(os.Args[0])
  logfile := strings.ToLower(arg + ".log")
  file,err := os.OpenFile(logfile,os.O_CREATE|os.O_WRONLY|os.O_APPEND,0666)
  if err != nil{
     log.Fatalln("Failed to open file: ",err)
  }
  debug = log.New(os.Stdout,"Trace",log.Ldate|log.Ltime|log.Lshortfile)
  infor = log.New(file,"Info",log.Ldate|log.Ltime|log.Lshortfile)
  warn = log.New(file,"Warn",log.Ldate|log.Ltime|log.Lshortfile)
  broken = log.New(io.MultiWriter(os.Stdout,file),"Error",log.Ldate|log.Ltime|log.Lshortfile)
}

func Infotf(format string,info ...interface{}){
  infor.Printf(format,info...)
}

func Debugtf(format string,info ...interface{}){
  debug.Printf(format,info...)
}

func Warntf(format string,info ...interface{}){
  warn.Printf(format,info...)
}

func Errortf(format string,info ...interface{}){
  broken.Printf(format,info...)
}

func Infoln(info ...interface{}){
  infor.Println(info...)
}

func Debugln(info...interface{}){
  debug.Println(info...)
}

func Warnln(info...interface{}){
  warn.Println(info...)
}

func Errorln(info...interface{}){
  broken.Fatalln(info...)
}
