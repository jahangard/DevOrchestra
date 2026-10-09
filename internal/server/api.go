package server

import (
 "errors"
 "fmt"
 "net/http"
 "strconv"
 "strings"
 "github.com/jahangard/DevOrchestra/internal/codex"
 "github.com/jahangard/DevOrchestra/internal/store"
)

func(s *Server)status(w http.ResponseWriter,r *http.Request){
 if !method(w,r,"GET"){return}
 output(w,struct{Engine codex.Status;DBPath string;Version string}{Engine:s.Codex.Status(),DBPath:s.DBPath,Version:"0.1.0"})
}
func(s *Server)projects(w http.ResponseWriter,r *http.Request){
 switch r.Method{
 case "GET":
  p,err:=s.DB.Projects();if err!=nil{bad(w,err);return};output(w,p)
 case "POST":
  var data struct{Name string;Path string}
  if err:=input(w,r,&data);err!=nil{bad(w,err);return}
  p,err:=s.DB.AddProject(data.Name,data.Path);if err!=nil{bad(w,err);return}
  output(w,p)
 default:method(w,r,"GET","POST")
 }
}
func(s *Server)templates(w http.ResponseWriter,r *http.Request){
 switch r.Method{
 case "GET":
  p,err:=s.DB.Templates();if err!=nil{bad(w,err);return};output(w,p)
 case "POST":
  var data struct{Title string;Category string;Body string}
  if err:=input(w,r,&data);err!=nil{bad(w,err);return}
  p,err:=s.DB.AddTemplate(data.Title,data.Category,data.Body,"");if err!=nil{bad(w,err);return}
  output(w,p)
 default:method(w,r,"GET","POST")
 }
}
func(s *Server)installCodex(w http.ResponseWriter,r *http.Request){
 if !method(w,r,"POST"){return}
 message,err:=s.Codex.InstallCLI()
 if err!=nil{bad(w,fmt.Errorf("%v: %s",err,message));return}
 output(w,map[string]string{"message":message})
}
func(s *Server)runs(w http.ResponseWriter,r *http.Request){
 switch r.Method{
 case "GET":
  p,err:=s.DB.Runs();if err!=nil{bad(w,err);return};output(w,p)
 case "POST":
  var data struct{ProjectID int64;Prompt string;Sandbox string}
  if err:=input(w,r,&data);err!=nil{bad(w,err);return}
  id,err:=s.Codex.Start(data.ProjectID,data.Prompt,data.Sandbox)
  if err!=nil{bad(w,err);return}
  output(w,map[string]int64{"id":id})
 default:method(w,r,"GET","POST")
 }
}
func(s *Server)runDetail(w http.ResponseWriter,r *http.Request){
 parts:=strings.Split(strings.Trim(strings.TrimPrefix(r.URL.Path,"/api/runs/"),"/"),"/")
 if len(parts)==0||parts[0]==""{http.NotFound(w,r);return}
 id,err:=strconv.ParseInt(parts[0],10,64)
 if err!=nil||id<=0{bad(w,errors.New("invalid run ID"));return}
 if len(parts)==2&&parts[1]=="cancel"{
  if !method(w,r,"POST"){return}
  if err:=s.Codex.Cancel(id);err!=nil{bad(w,err);return}
  output(w,map[string]bool{"ok":true});return
 }
 if len(parts)==2&&parts[1]=="events"{
  if !method(w,r,"GET"){return}
  after,_:=strconv.ParseInt(r.URL.Query().Get("after"),10,64)
  if after<0{after=0}
  run,err:=s.DB.Run(id);if err!=nil{bad(w,err);return}
  events,err:=s.DB.Events(id,after);if err!=nil{bad(w,err);return}
  output(w,struct{Run store.Run;Events []store.Event}{Run:run,Events:events});return
 }
 http.NotFound(w,r)
}
