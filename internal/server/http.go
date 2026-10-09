package server

import (
 "context"
 "crypto/rand"
 "crypto/subtle"
 "encoding/hex"
 "encoding/json"
 "errors"
 "io"
 "net"
 "net/http"
 "os/exec"
 "runtime"
 "strings"
 "time"
 "github.com/jahangard/DevOrchestra/internal/codex"
 "github.com/jahangard/DevOrchestra/internal/store"
 "github.com/jahangard/DevOrchestra/web"
)

type Server struct{
 DB *store.DB
 Codex *codex.Manager
 DBPath string
 addr string
 token string
 listener net.Listener
 httpServer *http.Server
}
func New(db *store.DB,c *codex.Manager,dbPath string)(*Server,error){
 listener,err:=net.Listen("tcp","127.0.0.1:0")
 if err!=nil{return nil,err}
 token:=make([]byte,32)
 if _,err=rand.Read(token);err!=nil{listener.Close();return nil,err}
 s:=&Server{DB:db,Codex:c,DBPath:dbPath,addr:listener.Addr().String(),token:hex.EncodeToString(token),listener:listener}
 mux:=http.NewServeMux()
 mux.HandleFunc("/",s.home)
 mux.HandleFunc("/app.js",s.asset)
 mux.HandleFunc("/style.css",s.asset)
 mux.HandleFunc("/api/status",s.status)
 mux.HandleFunc("/api/projects",s.projects)
 mux.HandleFunc("/api/templates",s.templates)
 mux.HandleFunc("/api/import",s.importPrompt)
 mux.HandleFunc("/api/install-codex",s.installCodex)
 mux.HandleFunc("/api/runs",s.runs)
 mux.HandleFunc("/api/runs/",s.runDetail)
 s.httpServer=&http.Server{Handler:s.guard(mux),ReadHeaderTimeout:5*time.Second}
 return s,nil
}
func(s *Server)URL()string{return "http://"+s.addr}
func(s *Server)Serve()error{
 err:=s.httpServer.Serve(s.listener)
 if errors.Is(err,http.ErrServerClosed){return nil}
 return err
}
func(s *Server)Close(ctx context.Context)error{return s.httpServer.Shutdown(ctx)}
func OpenBrowser(url string)error{
 var cmd *exec.Cmd
 switch runtime.GOOS{
 case "windows":cmd=exec.Command("rundll32","url.dll,FileProtocolHandler",url)
 case "darwin":cmd=exec.Command("open",url)
 default:cmd=exec.Command("xdg-open",url)
 }
 return cmd.Start()
}
func(s *Server)guard(next http.Handler)http.Handler{
 return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  if r.Host!=s.addr{http.Error(w,"invalid host",403);return}
  w.Header().Set("X-Content-Type-Options","nosniff")
  w.Header().Set("X-Frame-Options","DENY")
  w.Header().Set("Cache-Control","no-store")
  w.Header().Set("Content-Security-Policy","default-src 'none'; script-src 'self'; style-src 'self'; connect-src 'self'; img-src 'self' data:; base-uri 'none'; frame-ancestors 'none'")
  if origin:=r.Header.Get("Origin");origin!="" && origin!=s.URL(){http.Error(w,"invalid origin",403);return}
  if strings.HasPrefix(r.URL.Path,"/api/"){
   if subtle.ConstantTimeCompare([]byte(r.Header.Get("X-DevOrchestra-Token")),[]byte(s.token))!=1{http.Error(w,"unauthorized",401);return}
  }
  next.ServeHTTP(w,r)
 })
}
func(s *Server)home(w http.ResponseWriter,r *http.Request){
 if r.URL.Path!="/"{http.NotFound(w,r);return}
 file,err:=web.FS.ReadFile("index.html")
 if err!=nil{http.Error(w,err.Error(),500);return}
 w.Header().Set("Content-Type","text/html; charset=utf-8")
 io.WriteString(w,strings.ReplaceAll(string(file),"__TOKEN__",s.token))
}
func(s *Server)asset(w http.ResponseWriter,r *http.Request){
 name:=strings.TrimPrefix(r.URL.Path,"/")
 if name!="app.js"&&name!="style.css"{http.NotFound(w,r);return}
 b,err:=web.FS.ReadFile(name)
 if err!=nil{http.Error(w,err.Error(),500);return}
 if name=="app.js"{w.Header().Set("Content-Type","text/javascript; charset=utf-8")}else{w.Header().Set("Content-Type","text/css; charset=utf-8")}
 w.Write(b)
}
func method(w http.ResponseWriter,r *http.Request,allowed ...string)bool{
 for _,m:=range allowed{if r.Method==m{return true}}
 w.Header().Set("Allow",strings.Join(allowed,", "))
 http.Error(w,"method not allowed",405)
 return false
}
func output(w http.ResponseWriter,v any){
 w.Header().Set("Content-Type","application/json; charset=utf-8")
 json.NewEncoder(w).Encode(v)
}
func bad(w http.ResponseWriter,err error){http.Error(w,err.Error(),400)}
func input(w http.ResponseWriter,r *http.Request,v any)error{
 r.Body=http.MaxBytesReader(w,r.Body,130000)
 decoder:=json.NewDecoder(r.Body);decoder.DisallowUnknownFields()
 return decoder.Decode(v)
}
