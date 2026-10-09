package server

import(
 "errors"
 "fmt"
 "io"
 "net/http"
 "net/url"
 "strings"
 "time"
)

func(s *Server)importPrompt(w http.ResponseWriter,r *http.Request){
 if !method(w,r,"POST"){return}
 var data struct{URL string}
 if err:=input(w,r,&data);err!=nil{bad(w,err);return}
 parsed,err:=url.Parse(strings.TrimSpace(data.URL))
 if err!=nil||parsed==nil||parsed.Scheme!="https"||parsed.Host!="raw.githubusercontent.com"||parsed.User!=nil||
 (!strings.HasSuffix(strings.ToLower(parsed.Path),".md")&&!strings.HasSuffix(strings.ToLower(parsed.Path),".txt")){
  bad(w,errors.New("only HTTPS Markdown/text URLs on raw.githubusercontent.com are supported"));return
 }
 client:=&http.Client{
  Timeout:15*time.Second,
  CheckRedirect:func(req *http.Request,via []*http.Request)error{return errors.New("redirects disabled")},
 }
 req,err:=http.NewRequestWithContext(r.Context(),"GET",parsed.String(),nil)
 if err!=nil{bad(w,err);return}
 response,err:=client.Do(req)
 if err!=nil{bad(w,err);return}
 defer response.Body.Close()
 if response.StatusCode!=200{bad(w,fmt.Errorf("source returned HTTP %d",response.StatusCode));return}
 raw,err:=io.ReadAll(io.LimitReader(response.Body,40001))
 if err!=nil{bad(w,err);return}
 if len(raw)>40000{bad(w,errors.New("source is too large; maximum 40KB"));return}
 title:="Imported reference"
 for _,line:=range strings.Split(string(raw),"\n"){
  t:=strings.TrimSpace(line)
  if strings.HasPrefix(t,"# "){title=strings.TrimSpace(strings.TrimPrefix(t,"# "));break}
 }
 if len(title)>150{title=title[:150]}
 item,err:=s.DB.AddTemplate(title,"Imported - review before use",string(raw),parsed.String())
 if err!=nil{bad(w,err);return}
 output(w,item)
}
