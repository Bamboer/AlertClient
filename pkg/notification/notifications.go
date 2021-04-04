package notification
import(
  "io"
  "os"
  "log"
  "flag"
  "path"
  "strings"
  "grafana/pkg/client"
)


var (
   info    *log.Logger
   SNS   =  make(map[string] func(state string,msg client.SimpleInfo,b []byte)err )
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
   Send(state string,msg client.SimpleInfo},b []byte) error
}

func Emit(state string ,msg client.SimpleInfo,b []byte){
   for k,send := range(SNS){
     if err := send(state,msg,b);err !=nil{
         info.Println(k,"send err: ",err)
     }else{
         info.Println(k,"send message: ")
     }
   }
}
