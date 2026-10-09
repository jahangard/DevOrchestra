package codex

import (
 "bufio"
 "context"
 "encoding/json"
 "errors"
 "fmt"
 "os"
 "os/exec"
 "runtime"
 "strings"
 "sync"
 "time"

 "github.com/jahangard/DevOrchestra/internal/store"
)

type Status struct {
 CodexInstalled bool
 CodexVersion string
 NodeVersion string
 NPMVersion string
 Authenticated bool
 ActiveRunID int64
}
type Manager struct {
 Store *store.DB
 mu sync.Mutex
 running int64
 cancel context.CancelFunc
 installing bool
}

func New(s *store.DB)*Manager {return &Manager{Store:s}}

func executable(name string) string {
 if runtime.GOOS=="windows" && name=="npm" {
  if p,err:=exec.LookPath("npm.cmd");err==nil{return p}
 }
 p,err:=exec.LookPath(name)
 if err!=nil{return ""}
 return p
}

func version(program string,args ...string)string{
 path:=executable(program)
 if path==""{return ""}
 ctx,cancel:=context.WithTimeout(context.Background(),4*time.Second);defer cancel()
 data,err:=exec.CommandContext(ctx,path,args...).CombinedOutput()
 if err!=nil{return ""}
 txt:=strings.TrimSpace(string(data))
 if len(txt)>200{txt=txt[:200]}
 return txt
}
func (m *Manager) Status() Status {
 s:=Status{CodexVersion:version("codex","--version"), NodeVersion:version("node","--version"),NPMVersion:version("npm","--version")}
 s.CodexInstalled=s.CodexVersion!=""
 if s.CodexInstalled {
  path:=executable("codex")
  ctx,cancel:=context.WithTimeout(context.Background(),4*time.Second)
  cmd:=exec.CommandContext(ctx,path,"login","status")
  output,err:=cmd.CombinedOutput()
  cancel()
  s.Authenticated=err==nil && !strings.Contains(strings.ToLower(string(output)),"not logged")
 }
 m.mu.Lock();s.ActiveRunID=m.running;m.mu.Unlock()
 return s
}
func (m *Manager) InstallCLI() (string,error) {
 m.mu.Lock()
 if m.running!=0 || m.installing {m.mu.Unlock();return "",errors.New("another operation is in progress")}
 m.installing=true
 m.mu.Unlock()
 defer func(){m.mu.Lock();m.installing=false;m.mu.Unlock()}()
 if executable("npm")=="" {return "",errors.New("Node.js and npm must be installed first: https://nodejs.org/")}
 ctx,cancel:=context.WithTimeout(context.Background(),6*time.Minute);defer cancel()
 cmd:=exec.CommandContext(ctx,executable("npm"),"install","-g","@openai/codex")
 data,err:=cmd.CombinedOutput()
 text:=string(data)
 if len(text)>4000{text=text[len(text)-4000:]}
 if err!=nil{return text,fmt.Errorf("npm install failed: %w",err)}
 return text,nil
}
func (m *Manager) Start(projectID int64,prompt,sandbox string)(int64,error){
 prompt=strings.TrimSpace(prompt)
 if prompt=="" || len(prompt)>64000{return 0,errors.New("prompt must contain 1..64000 characters")}
 if sandbox!="read-only" && sandbox!="workspace-write" {return 0,errors.New("invalid sandbox policy")}
 if executable("codex")==""{return 0,errors.New("Codex CLI is not installed")}
 project,err:=m.Store.Project(projectID);if err!=nil{return 0,errors.New("project not found")}
 st,err:=os.Stat(project.Path);if err!=nil || !st.IsDir(){return 0,errors.New("project folder no longer exists")}
 m.mu.Lock()
 defer m.mu.Unlock()
 if m.running!=0 || m.installing {return 0,errors.New("another operation is in progress")}
 id,err:=m.Store.StartRun(projectID,prompt,sandbox);if err!=nil{return 0,err}
 ctx,cancel:=context.WithCancel(context.Background())
 m.running=id;m.cancel=cancel
 go m.execute(ctx,id,project.Path,prompt,sandbox)
 return id,nil
}
func (m *Manager) Cancel(id int64)error {
 m.mu.Lock();defer m.mu.Unlock()
 if id==0 || id!=m.running || m.cancel==nil{return errors.New("run is not active")}
 m.cancel()
 return nil
}
func (m *Manager) execute(ctx context.Context,id int64,dir,prompt,sandbox string){
 defer func(){
  m.mu.Lock()
  if m.running==id {m.running=0;m.cancel=nil}
  m.mu.Unlock()
 }()
 cmd:=exec.CommandContext(ctx,executable("codex"),"exec","--json","--sandbox",sandbox,"-")
 cmd.Dir=dir
 cmd.Stdin=strings.NewReader(prompt)
 stdout,err:=cmd.StdoutPipe()
 if err!=nil{m.Store.EndRun(id,"failed",err.Error());return}
 stderr,err:=cmd.StderrPipe()
 if err!=nil{m.Store.EndRun(id,"failed",err.Error());return}
 if err=cmd.Start();err!=nil{m.Store.EndRun(id,"failed",err.Error());return}
 var wg sync.WaitGroup
 wg.Add(2)
 go func(){
  defer wg.Done()
  scanner:=bufio.NewScanner(stdout)
  scanner.Buffer(make([]byte,65536),4*1024*1024)
  for scanner.Scan(){
   line:=scanner.Text()
   m.Store.AddEvent(id,"stdout",line)
   var evt struct{
    Type string
    ThreadID string
    Item struct{Type string;Text string}
   }
   if json.Unmarshal([]byte(line),&evt)==nil {
    if evt.Type=="thread.started" && evt.ThreadID!=""{m.Store.SetThread(id,evt.ThreadID)}
    if evt.Type=="item.completed" && evt.Item.Type=="agent_message" {m.Store.SetSummary(id,evt.Item.Text)}
   }
  }
  if err:=scanner.Err();err!=nil{m.Store.AddEvent(id,"system",err.Error())}
 }()
 go func(){
  defer wg.Done()
  scanner:=bufio.NewScanner(stderr)
  scanner.Buffer(make([]byte,4096),1024*1024)
  for scanner.Scan(){m.Store.AddEvent(id,"stderr",scanner.Text())}
  if err:=scanner.Err();err!=nil{m.Store.AddEvent(id,"system",err.Error())}
 }()
 wg.Wait()
 err=cmd.Wait()
 if ctx.Err()!=nil {m.Store.EndRun(id,"cancelled","cancelled by user");return}
 if err!=nil {m.Store.EndRun(id,"failed",err.Error());return}
 m.Store.EndRun(id,"completed","")
}
