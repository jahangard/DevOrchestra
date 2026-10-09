"use strict";
const token=document.querySelector('meta[name="api-token"]').content;
const $=id=>document.getElementById(id);
const state={lang:localStorage.getItem("devorchestra-language")==="en"?"en":"fa",projects:[],templates:[],runs:[],status:null,current:0,lastEvent:0,polling:false};
const tr={
fa:{dashboard:"داشبورد",projects:"پروژه‌ها",templates:"کتابخانه پرامپت",history:"تاریخچه اجرا",settings:"تنظیمات",workspace:"فضای کاری محلی",heroTitle:"کدنویسی با Codex، ساده‌تر از همیشه",heroDesc:"پروژه‌ات را انتخاب کن، درخواستت را بنویس و تمام مراحل اجرا را همین‌جا دنبال کن.",cliStatus:"وضعیت Codex",projectCount:"پروژه‌ها",runCount:"تعداد اجراها",newRequest:"درخواست جدید",newRequestSub:"درخواست مستقیماً برای Codex CLI ارسال می‌شود.",localRun:"اجرای محلی",targetProject:"پروژه مقصد",permissions:"سطح دسترسی",readOnly:"فقط خواندنی (امن‌تر)",workspaceWrite:"اجازه ویرایش پروژه",prompt:"پرامپت",promptPlaceholder:"دقیقاً بگو چه چیزی را در پروژه تغییر دهیم...",enhance:"تکمیل ساختاریافته",saveTemplate:"ذخیره پرامپت",execute:"اجرای Codex ←",liveOutput:"خروجی و رویدادهای اجرا",noRun:"هنوز اجرایی انتخاب نشده است.",cancel:"توقف اجرا",addProject:"افزودن پروژه",projectNote:"مسیر پوشه باید روی همین کامپیوتر وجود داشته باشد.",name:"نام پروژه",folderPath:"مسیر پوشه",add:"افزودن",myProjects:"پروژه‌های من",promptLibrary:"کتابخانه پرامپت",templateNote:"الگوهای آماده، ذخیره‌شده و واردشده از منابع مرجع.",importTitle:"واردکردن از منبع مرجع",importHint:"فقط فایل‌های Markdown یا Text از raw.githubusercontent.com. محتوای واردشده قبل از استفاده باید بازبینی شود.",import:"دریافت",recentRuns:"آخرین اجراها",refresh:"بازخوانی",tools:"وضعیت ابزارها",authentication:"احراز هویت",installCodex:"نصب Codex CLI با npm",loginTitle:"ورود به Codex",loginHelp:"پس از نصب، دستور زیر را در ترمینال اجرا کن. DevOrchestra رمز یا توکن حساب را ذخیره نمی‌کند.",copy:"کپی",ready:"آماده",missing:"نصب نشده",yes:"وارد شده",no:"وارد نشده",use:"استفاده",open:"نمایش نتیجه",empty:"موردی وجود ندارد.",running:"در حال اجرا",completed:"تکمیل‌شده",failed:"ناموفق",cancelled:"متوقف‌شده",interrupted:"قطع‌شده",promptTitle:"نام پرامپت",saved:"ذخیره شد",working:"در حال انجام...",error:"خطا",reviewExternal:"پرامپت‌های دریافتی را قبل از اجرا بررسی کن.",unknown:"نامشخص"},
en:{dashboard:"Dashboard",projects:"Projects",templates:"Prompt library",history:"Run history",settings:"Settings",workspace:"LOCAL WORKSPACE",heroTitle:"Work with Codex, without the friction",heroDesc:"Choose a project, enter your request and follow every execution event here.",cliStatus:"Codex status",projectCount:"Projects",runCount:"Runs",newRequest:"New request",newRequestSub:"Your request is sent directly to Codex CLI.",localRun:"Local execution",targetProject:"Target project",permissions:"Sandbox policy",readOnly:"Read only (safer)",workspaceWrite:"Allow project writes",prompt:"Prompt",promptPlaceholder:"Describe exactly what needs changing in your project...",enhance:"Structure prompt",saveTemplate:"Save prompt",execute:"Run Codex →",liveOutput:"Execution output",noRun:"No run selected yet.",cancel:"Cancel",addProject:"Add project",projectNote:"The folder must exist on this computer.",name:"Project name",folderPath:"Folder path",add:"Add",myProjects:"My projects",promptLibrary:"Prompt library",templateNote:"Built-in, saved and imported templates.",importTitle:"Import reference",importHint:"Only Markdown or text on raw.githubusercontent.com. Review untrusted imports before use.",import:"Import",recentRuns:"Recent runs",refresh:"Refresh",tools:"Tool status",authentication:"Authentication",installCodex:"Install Codex CLI with npm",loginTitle:"Sign in to Codex",loginHelp:"After installation run this command in a terminal. DevOrchestra never stores your account password or tokens.",copy:"Copy",ready:"Ready",missing:"Not installed",yes:"Signed in",no:"Not signed in",use:"Use",open:"View result",empty:"Nothing to show.",running:"Running",completed:"Completed",failed:"Failed",cancelled:"Cancelled",interrupted:"Interrupted",promptTitle:"Prompt title",saved:"Saved",working:"Working...",error:"Error",reviewExternal:"Review imported prompts before execution.",unknown:"Unknown"}
};
function t(key){return tr[state.lang][key]||key}
function notify(message,error=false){const n=$("notice");n.hidden=false;n.className="notice"+(error?" error":"");n.textContent=message;window.setTimeout(()=>{if(n.textContent===message)n.hidden=true},7000)}
async function api(path,method="GET",body){
 const response=await fetch(path,{method,headers:{"X-DevOrchestra-Token":token,...(body?{"Content-Type":"application/json"}:{})},body:body?JSON.stringify(body):undefined,cache:"no-store"});
 if(!response.ok)throw new Error((await response.text()).trim());
 return response.json()
}
async function safe(callback){try{return await callback()}catch(e){notify(e.message||String(e),true)}}
function applyLanguage(){
 document.documentElement.lang=state.lang;document.documentElement.dir=state.lang==="fa"?"rtl":"ltr";
 document.querySelectorAll("[data-i18n]").forEach(e=>{e.textContent=t(e.dataset.i18n)});
 document.querySelectorAll("[data-i18n-placeholder]").forEach(e=>{e.placeholder=t(e.dataset.i18nPlaceholder)});
 $("language").textContent=state.lang==="fa"?"English":"فارسی";
 document.querySelector(".nav-item.active")?.click();
 draw();
}
function view(name){
 document.querySelectorAll(".view").forEach(e=>{e.hidden=e.id!=="view-"+name});
 document.querySelectorAll(".nav-item").forEach(e=>e.classList.toggle("active",e.dataset.view===name));
 $("view-title").textContent=t(name);
}
function textElement(tag,text,klass){const e=document.createElement(tag);if(klass)e.className=klass;e.textContent=text;return e}
function empty(target){target.replaceChildren(textElement("p",t("empty"),"empty"))}
function draw(){
 const status=state.status;
 $("stat-cli").textContent=status?.Engine?.CodexInstalled?t("ready"):t("missing");
 $("stat-projects").textContent=String(state.projects.length);
 $("stat-runs").textContent=String(state.runs.length);
 $("cli-dot").classList.toggle("online",!!status?.Engine?.CodexInstalled);
 $("cli-short").textContent=status?.Engine?.CodexInstalled?"Codex CLI ✓":"Codex CLI —";
 $("version-codex").textContent=status?.Engine?.CodexVersion||t("missing");
 $("version-node").textContent=status?.Engine?.NodeVersion||t("missing");
 $("version-npm").textContent=status?.Engine?.NPMVersion||t("missing");
 $("version-auth").textContent=status?.Engine?.Authenticated?t("yes"):t("no");
 $("database-path").textContent=status?.DBPath||"—";
 const select=$("project-select"),old=select.value;select.replaceChildren();
 for(const p of state.projects){const opt=document.createElement("option");opt.value=String(p.ID);opt.textContent=p.Name+" — "+p.Path;select.append(opt)}
 if(state.projects.some(p=>String(p.ID)===old))select.value=old;
 if(!state.projects.length){const opt=document.createElement("option");opt.value="";opt.textContent=t("addProject");select.append(opt)}
 drawProjects();drawTemplates();drawHistory();
}
function drawProjects(){
 const el=$("project-list");el.replaceChildren();if(!state.projects.length){empty(el);return}
 for(const p of state.projects){
  const item=textElement("div","","collection-item"),body=document.createElement("div");
  body.append(textElement("strong",p.Name),textElement("p",p.Path));item.append(body);el.append(item);
 }
}
function drawTemplates(){
 const el=$("template-list");el.replaceChildren();if(!state.templates.length){empty(el);return}
 for(const p of state.templates){
  const item=textElement("div","","collection-item"),body=document.createElement("div");
  body.append(textElement("strong",p.Title),textElement("p",p.Category+" · "+p.Body.slice(0,140).replace(/\s+/g," ")));
  if(p.SourceURL){const a=textElement("a",p.SourceURL,"muted");a.href=p.SourceURL;a.target="_blank";a.rel="noopener noreferrer";body.append(a)}
  const actions=textElement("div","","actions"),button=textElement("button",t("use"),"button secondary");
  button.type="button";button.addEventListener("click",()=>{
   $("prompt-text").value=p.Body;view("dashboard");$("prompt-text").focus();
   if(p.SourceURL)notify(t("reviewExternal"))
  });actions.append(button);item.append(body,actions);el.append(item)
 }
}
function drawHistory(){
 const el=$("history-list");el.replaceChildren();if(!state.runs.length){empty(el);return}
 for(const r of state.runs){
  const item=textElement("div","","collection-item"),body=document.createElement("div"),actions=textElement("div","","actions");
  body.append(textElement("strong",r.ProjectName+" · "+t(r.Status)),textElement("p",r.Prompt.slice(0,165).replace(/\s+/g," ")));
  const date=textElement("span",new Date(r.StartedAt).toLocaleString(state.lang==="fa"?"fa-IR":"en-US"),"muted");body.append(date);
  const b=textElement("button",t("open"),"button secondary");b.addEventListener("click",()=>{view("dashboard");showRun(r.ID)});
  actions.append(b);item.append(body,actions);el.append(item);
 }
}
async function refresh(){
 const [status,projects,templates,runs]=await Promise.all([api("/api/status"),api("/api/projects"),api("/api/templates"),api("/api/runs")]);
 state.status=status;state.projects=projects;state.templates=templates;state.runs=runs;draw();
}
function enhancePrompt(){
 const input=$("prompt-text"),current=input.value.trim();if(!current){notify(t("promptPlaceholder"),true);return}
 const missing=state.lang==="fa"?
 "\n\n## زمینه و محدوده\n- پروژه و فایل‌های مرتبط را بررسی کن.\n- اصول و معماری موجود را حفظ کن.\n\n## معیارهای پذیرش\n- تغییرات کوچک، دقیق و قابل نگهداری باشند.\n- تست‌های مرتبط اجرا شوند.\n- فایل‌های تغییر یافته و نتیجه تست‌ها گزارش شوند.\n- در صورت ابهام یا خطر حذف داده، قبل از اقدام سؤال کن.":
 "\n\n## Context and scope\n- Inspect relevant files and project conventions.\n- Preserve the current architecture.\n\n## Acceptance criteria\n- Make small, maintainable changes.\n- Run relevant tests.\n- Report changed files and test results.\n- Ask before any destructive or ambiguous action.";
 if(!current.includes("## Acceptance criteria")&&!current.includes("## معیارهای پذیرش"))input.value=current+missing;
}
function summarizeEvent(ev){
 if(ev.Source!=="stdout")return "["+ev.Source+"] "+ev.Body;
 try{
  const j=JSON.parse(ev.Body),item=j.item;
  if(item&&item.type==="agent_message")return "[assistant] "+(item.text||"");
  if(item&&item.type==="command_execution")return "[command] "+(item.command||"")+" "+(item.status||"");
  if(j.type==="turn.completed")return "[turn.completed] "+JSON.stringify(j.usage||{});
  return "["+(j.type||"event")+"] "+ev.Body.slice(0,750);
 }catch{return ev.Body}
}
async function showRun(id){
 state.current=id;state.lastEvent=0;$("event-stream").textContent="";await poll();
}
async function poll(){
 if(!state.current||state.polling)return;
 state.polling=true;
 try{
  const data=await api("/api/runs/"+state.current+"/events?after="+state.lastEvent);
  for(const e of data.Events){
   state.lastEvent=e.ID;
   $("event-stream").textContent+=summarizeEvent(e)+"\n";
  }
  const pre=$("event-stream");if(pre.textContent.length>150000)pre.textContent=pre.textContent.slice(-150000);pre.scrollTop=pre.scrollHeight;
  $("run-caption").textContent="#"+data.Run.ID+" · "+data.Run.ProjectName+" · "+t(data.Run.Status);
  $("cancel-run").hidden=data.Run.Status!=="running";
  if(data.Run.Status!=="running"){await refresh()}
  else window.setTimeout(()=>{state.polling=false;poll()},1100);
 }catch(e){notify(e.message,true)}
 finally{
  if(!state.current||$("cancel-run").hidden)state.polling=false;
 }
}
document.querySelectorAll(".nav-item").forEach(b=>b.addEventListener("click",()=>view(b.dataset.view)));
$("language").addEventListener("click",()=>{state.lang=state.lang==="fa"?"en":"fa";localStorage.setItem("devorchestra-language",state.lang);applyLanguage()});
$("enhance").addEventListener("click",enhancePrompt);
$("copy-login").addEventListener("click",()=>safe(async()=>{await navigator.clipboard.writeText("codex login");notify(t("copy"))}));
$("project-form").addEventListener("submit",e=>{e.preventDefault();safe(async()=>{
 const data=new FormData(e.currentTarget);await api("/api/projects","POST",{Name:data.get("name"),Path:data.get("path")});e.target.reset();await refresh();notify(t("saved"))
})});
$("save-template").addEventListener("click",()=>safe(async()=>{
 const body=$("prompt-text").value.trim();if(!body)return notify(t("promptPlaceholder"),true);
 const title=window.prompt(t("promptTitle"));if(!title)return;
 await api("/api/templates","POST",{Title:title,Category:"Custom",Body:body});await refresh();notify(t("saved"))
}));
$("import-form").addEventListener("submit",e=>{e.preventDefault();safe(async()=>{
 const data=new FormData(e.target);await api("/api/import","POST",{URL:data.get("url")});e.target.reset();await refresh();notify(t("reviewExternal"))
})});
$("run-form").addEventListener("submit",e=>{e.preventDefault();safe(async()=>{
 const button=$("run-button");button.disabled=true;
 try{
  const result=await api("/api/runs","POST",{ProjectID:Number($("project-select").value),Prompt:$("prompt-text").value,Sandbox:$("sandbox-select").value});
  await refresh();await showRun(result.id);
 }finally{button.disabled=false}
})});
$("cancel-run").addEventListener("click",()=>safe(async()=>{await api("/api/runs/"+state.current+"/cancel","POST",{});notify(t("cancel"))}));
$("refresh-status").addEventListener("click",()=>safe(refresh));
$("refresh-history").addEventListener("click",()=>safe(refresh));
$("install-codex").addEventListener("click",()=>safe(async()=>{
 if(!window.confirm("Install @openai/codex using npm? / نصب Codex با npm؟"))return;
 const button=$("install-codex");button.disabled=true;button.textContent=t("working");
 try{await api("/api/install-codex","POST",{});notify(t("saved"));await refresh()}
 finally{button.disabled=false;button.textContent=t("installCodex")}
}));
applyLanguage();safe(async()=>{await refresh();if(state.status?.Engine?.ActiveRunID)await showRun(state.status.Engine.ActiveRunID)});
