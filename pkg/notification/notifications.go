package notification
import(
  "io"
  "os"
  "log"
  "flag"
  "path"
  "strings"
)


var (
   info    *log.Logger
   SNS   =  make(map[string]Notification)
)

func init() {
        arg := path.Base(os.Args[0])
        logfile := strings.ToLower(arg + ".log")
        file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
        if err != nil {
                log.Println("Failed to open file: ", err)
        }
        info = log.New(io.MultiWriter(os.Stdout, file), "Info: ", log.Ldate|log.Ltime|log.Lshortfile)
}

type Notification interface{
//send message to the receiver
   Send(state string,alertNum int,msg interface{}) error
}

func Emit(msg interface{},b []byte){
   for k,v := range(SNS){
     if err := v.Send(msg,b);err !=nil{
         info.Println(k,"send err: ",err)
     }else{
         info.Println(k,"send message: ")
     }
   }
}
