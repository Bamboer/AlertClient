package client
import(
  "net/http"
  "path/filepath"
  "grafana/pkg/configer"
)

func FileServer(){
   conf := configer.ConfigParse()
   if conf.ImgServer{
       imgdir,err := filepath.Abs(conf.ImgDir)
       if err != nil{
         info.Println(err)
       }
       imgport := conf.ImgServerPort
       info.Println("started image server: ","0.0.0.0",imgport)
       if err := http.ListenAndServe(imgport, http.FileServer(http.Dir(imgdir)));err != nil{
          info.Println(err)
       }else{
          info.Println("Not start image server.")
      }
   }
}
