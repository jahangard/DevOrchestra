package main

import(
 "context"
 "fmt"
 "log"
 "os"
 "os/signal"
 "syscall"
 "time"
 "github.com/jahangard/DevOrchestra/internal/codex"
 "github.com/jahangard/DevOrchestra/internal/server"
 "github.com/jahangard/DevOrchestra/internal/store"
)

func main(){
 db,path,err:=store.Open()
 if err!=nil{log.Fatal(err)}
 defer db.Close()
 runner:=codex.New(db)
 app,err:=server.New(db,runner,path)
 if err!=nil{log.Fatal(err)}
 done:=make(chan error,1)
 go func(){done<-app.Serve()}()
 fmt.Println("DevOrchestra listening at",app.URL())
 fmt.Println("SQLite:",path)
 if err:=server.OpenBrowser(app.URL());err!=nil{
  fmt.Fprintln(os.Stderr,"Could not open browser:",err)
  fmt.Println("Open the URL above in a browser.")
 }
 signals:=make(chan os.Signal,1)
 signal.Notify(signals,os.Interrupt,syscall.SIGTERM)
 select{
 case err:=<-done:
  if err!=nil{log.Println(err)}
 case <-signals:
 }
 ctx,cancel:=context.WithTimeout(context.Background(),5*time.Second)
 defer cancel()
 if err:=app.Close(ctx);err!=nil{log.Println(err)}
}
