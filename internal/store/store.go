package store

import (
 "database/sql"
 "errors"
 "fmt"
 "os"
 "path/filepath"
 "strings"
 "time"

 _ "modernc.org/sqlite"
)

type DB struct { Conn *sql.DB }

type Project struct {
 ID int64
 Name string
 Path string
 CreatedAt string
}
type Template struct {
 ID int64
 Title string
 Category string
 Body string
 SourceURL string
 CreatedAt string
}
type Run struct {
 ID int64
 ProjectID int64
 ProjectName string
 Prompt string
 Sandbox string
 Status string
 ThreadID string
 Summary string
 Error string
 StartedAt string
 EndedAt string
}
type Event struct {
 ID int64
 RunID int64
 Source string
 Body string
 CreatedAt string
}

func timestamp() string { return time.Now().UTC().Format(time.RFC3339Nano) }

func Open() (*DB, string, error) {
 base, err := os.UserConfigDir()
 if err != nil { return nil, "", err }
 dir := filepath.Join(base, "DevOrchestra")
 if err = os.MkdirAll(dir, 0700); err != nil { return nil, "", err }
 path := filepath.Join(dir, "devorchestra.db")
 db, err := sql.Open("sqlite", path)
 if err != nil { return nil, "", err }
 db.SetMaxOpenConns(1)
 statements := []string{
  "PRAGMA busy_timeout = 5000",
  "PRAGMA journal_mode = WAL",
  "PRAGMA foreign_keys = ON",
  "CREATE TABLE IF NOT EXISTS projects (id INTEGER PRIMARY KEY, name TEXT NOT NULL, path TEXT NOT NULL UNIQUE, created_at TEXT NOT NULL)",
  "CREATE TABLE IF NOT EXISTS templates (id INTEGER PRIMARY KEY, title TEXT NOT NULL, category TEXT NOT NULL, body TEXT NOT NULL, source_url TEXT NOT NULL DEFAULT '', created_at TEXT NOT NULL)",
  "CREATE TABLE IF NOT EXISTS runs (id INTEGER PRIMARY KEY, project_id INTEGER NOT NULL REFERENCES projects(id), prompt TEXT NOT NULL, sandbox TEXT NOT NULL, status TEXT NOT NULL, thread_id TEXT NOT NULL DEFAULT '', summary TEXT NOT NULL DEFAULT '', error TEXT NOT NULL DEFAULT '', started_at TEXT NOT NULL, ended_at TEXT NOT NULL DEFAULT '')",
  "CREATE TABLE IF NOT EXISTS run_events (id INTEGER PRIMARY KEY, run_id INTEGER NOT NULL REFERENCES runs(id) ON DELETE CASCADE, source TEXT NOT NULL, body TEXT NOT NULL, created_at TEXT NOT NULL)",
  "CREATE INDEX IF NOT EXISTS idx_run_events ON run_events(run_id,id)",
  "CREATE INDEX IF NOT EXISTS idx_runs_started ON runs(started_at DESC)",
 }
 for _, stmt := range statements {
  if _, err = db.Exec(stmt); err != nil { db.Close(); return nil, "", fmt.Errorf("init db: %w",err) }
 }
 if _, err = db.Exec("UPDATE runs SET status='interrupted', ended_at=? WHERE status='running'", timestamp()); err != nil { db.Close(); return nil,"",err }
 samples := []struct{Title,Category,Body string}{
  {"Review architecture","Review","Review this repository architecture. Identify coupling, scalability risks and technical debt. Cite file paths. Recommend the three highest-value changes. Do not edit files."},
  {"Fix a bug","Debug","Investigate the following bug: [describe symptoms]. Reproduce it, identify its root cause, implement the smallest safe fix, add regression tests, and summarize the changed files."},
  {"Add a feature","Feature","Implement [feature]. First inspect existing patterns and AGENTS.md. Preserve project architecture, add tests, run relevant checks, and summarize design trade-offs."},
  {"Improve test coverage","Testing","Identify untested critical paths in the project. Propose a test plan, add focused tests for the highest-risk behavior, execute them and report results."},
  {"Security audit","Security","Perform a read-only security audit. Look for secret leaks, command injection, unsafe file handling and authorization gaps. Provide actionable findings with severity and file paths."},
 }
 for i,p := range samples {
  if _,err = db.Exec("INSERT OR IGNORE INTO templates(id,title,category,body,created_at) VALUES(?,?,?,?,?)", i+1,p.Title,p.Category,p.Body,timestamp()); err != nil {db.Close();return nil,"",err}
 }
 return &DB{Conn:db}, path, nil
}

func (d *DB) Close() error { return d.Conn.Close() }

func (d *DB) AddProject(name, dir string) (Project,error) {
 name = strings.TrimSpace(name)
 if len(name)<1 || len(name)>100 { return Project{},errors.New("project name must be 1..100 characters") }
 abs,err := filepath.Abs(strings.TrimSpace(dir))
 if err!=nil {return Project{},err}
 abs,err = filepath.EvalSymlinks(abs)
 if err!=nil {return Project{},fmt.Errorf("project folder: %w",err)}
 info,err := os.Stat(abs)
 if err!=nil || !info.IsDir() {return Project{},errors.New("project folder must exist")}
 p:=Project{Name:name,Path:abs,CreatedAt:timestamp()}
 result,err:=d.Conn.Exec("INSERT INTO projects(name,path,created_at) VALUES(?,?,?)",p.Name,p.Path,p.CreatedAt)
 if err!=nil{return Project{},err}
 p.ID,_=result.LastInsertId()
 return p,nil
}
func (d *DB) Projects() ([]Project,error) {
 rows,err:=d.Conn.Query("SELECT id,name,path,created_at FROM projects ORDER BY id DESC")
 if err!=nil{return nil,err}
 defer rows.Close()
 out:=[]Project{}
 for rows.Next(){var p Project; if err=rows.Scan(&p.ID,&p.Name,&p.Path,&p.CreatedAt);err!=nil{return nil,err};out=append(out,p)}
 return out,rows.Err()
}
func (d *DB) Project(id int64) (Project,error) {
 var p Project
 err:=d.Conn.QueryRow("SELECT id,name,path,created_at FROM projects WHERE id=?",id).Scan(&p.ID,&p.Name,&p.Path,&p.CreatedAt)
 return p,err
}
func (d *DB) Templates() ([]Template,error) {
 rows,err:=d.Conn.Query("SELECT id,title,category,body,source_url,created_at FROM templates ORDER BY id DESC LIMIT 500")
 if err!=nil{return nil,err}
 defer rows.Close()
 out:=[]Template{}
 for rows.Next(){var p Template;if err=rows.Scan(&p.ID,&p.Title,&p.Category,&p.Body,&p.SourceURL,&p.CreatedAt);err!=nil{return nil,err};out=append(out,p)}
 return out,rows.Err()
}
func (d *DB) AddTemplate(title,category,body,source string)(Template,error){
 title=strings.TrimSpace(title);body=strings.TrimSpace(body)
 if title=="" || len(title)>200 || body=="" || len(body)>40000 {return Template{},errors.New("title/body missing or too long")}
 if category==""{category="Custom"}
 t:=Template{Title:title,Category:category,Body:body,SourceURL:source,CreatedAt:timestamp()}
 result,err:=d.Conn.Exec("INSERT INTO templates(title,category,body,source_url,created_at) VALUES(?,?,?,?,?)",t.Title,t.Category,t.Body,t.SourceURL,t.CreatedAt)
 if err!=nil{return Template{},err}
 t.ID,_=result.LastInsertId()
 return t,nil
}
func (d *DB) StartRun(projectID int64,prompt,sandbox string)(int64,error){
 r,err:=d.Conn.Exec("INSERT INTO runs(project_id,prompt,sandbox,status,started_at) VALUES(?,?,?,'running',?)",projectID,prompt,sandbox,timestamp())
 if err!=nil{return 0,err}
 return r.LastInsertId()
}
func (d *DB) AddEvent(id int64, source, body string) error {
 if len(body)>200000 {body=body[:200000]+" [truncated]"}
 _,err:=d.Conn.Exec("INSERT INTO run_events(run_id,source,body,created_at) VALUES(?,?,?,?)",id,source,body,timestamp())
 return err
}
func (d *DB) SetThread(id int64,thread string) {
 d.Conn.Exec("UPDATE runs SET thread_id=? WHERE id=?",thread,id)
}
func (d *DB) SetSummary(id int64,summary string) {
 if len(summary)>12000{summary=summary[:12000]}
 d.Conn.Exec("UPDATE runs SET summary=? WHERE id=?",summary,id)
}
func (d *DB) EndRun(id int64,status,message string) {
 if len(message)>3000{message=message[:3000]}
 d.Conn.Exec("UPDATE runs SET status=?, error=?, ended_at=? WHERE id=?",status,message,timestamp(),id)
}
func (d *DB) Runs()([]Run,error){
 rows,err:=d.Conn.Query("SELECT r.id,r.project_id,p.name,r.prompt,r.sandbox,r.status,r.thread_id,r.summary,r.error,r.started_at,r.ended_at FROM runs r JOIN projects p ON p.id=r.project_id ORDER BY r.id DESC LIMIT 150")
 if err!=nil{return nil,err}
 defer rows.Close()
 out:=[]Run{}
 for rows.Next(){var r Run;if err=rows.Scan(&r.ID,&r.ProjectID,&r.ProjectName,&r.Prompt,&r.Sandbox,&r.Status,&r.ThreadID,&r.Summary,&r.Error,&r.StartedAt,&r.EndedAt);err!=nil{return nil,err};out=append(out,r)}
 return out,rows.Err()
}
func (d *DB) Run(id int64)(Run,error) {
 var r Run
 err:=d.Conn.QueryRow("SELECT r.id,r.project_id,p.name,r.prompt,r.sandbox,r.status,r.thread_id,r.summary,r.error,r.started_at,r.ended_at FROM runs r JOIN projects p ON p.id=r.project_id WHERE r.id=?",id).Scan(&r.ID,&r.ProjectID,&r.ProjectName,&r.Prompt,&r.Sandbox,&r.Status,&r.ThreadID,&r.Summary,&r.Error,&r.StartedAt,&r.EndedAt)
 return r,err
}
func (d *DB) Events(id,after int64)([]Event,error){
 rows,err:=d.Conn.Query("SELECT id,run_id,source,body,created_at FROM run_events WHERE run_id=? AND id>? ORDER BY id LIMIT 300",id,after)
 if err!=nil{return nil,err}
 defer rows.Close()
 out:=[]Event{}
 for rows.Next(){var e Event;if err=rows.Scan(&e.ID,&e.RunID,&e.Source,&e.Body,&e.CreatedAt);err!=nil{return nil,err};out=append(out,e)}
 return out,rows.Err()
}
