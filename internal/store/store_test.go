package store

import(
 "os"
 "path/filepath"
 "runtime"
 "testing"
)

func TestSQLiteLifecycle(t *testing.T){
 if runtime.GOOS=="windows"{t.Setenv("APPDATA",t.TempDir())}else{t.Setenv("XDG_CONFIG_HOME",t.TempDir())}
 db,path,err:=Open()
 if err!=nil{t.Fatal(err)}
 defer db.Close()
 if filepath.Base(path)!="devorchestra.db"{t.Fatal(path)}
 project,err:=db.AddProject("Example",t.TempDir())
 if err!=nil{t.Fatal(err)}
 projects,err:=db.Projects()
 if err!=nil||len(projects)!=1{t.Fatalf("projects: %v %v",projects,err)}
 id,err:=db.StartRun(project.ID,"Review code","read-only")
 if err!=nil{t.Fatal(err)}
 if err:=db.AddEvent(id,"stdout","test");err!=nil{t.Fatal(err)}
 events,err:=db.Events(id,0)
 if err!=nil||len(events)!=1||events[0].Body!="test"{t.Fatalf("events: %v %v",events,err)}
 db.EndRun(id,"completed","")
 r,err:=db.Run(id)
 if err!=nil||r.Status!="completed"{t.Fatalf("run: %v %v",r,err)}
 if _,err:=os.Stat(path);err!=nil{t.Fatal(err)}
}
