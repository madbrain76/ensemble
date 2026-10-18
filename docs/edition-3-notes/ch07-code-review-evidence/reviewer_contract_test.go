package gui

import (
 "context"
 "encoding/json"
 "ensemble"
 "fmt"
 "io"
 "log"
 "net/http"
 "net/http/httptest"
 "os"
 "path/filepath"
 "sync/atomic"
 "testing"
 "time"
)

func TestReviewerTruthfulRefusedToolProjection(t *testing.T) {
 for _, explained := range []bool{true,false} {
  t.Run(fmt.Sprintf("explained_%v",explained),func(t *testing.T){
   work:=t.TempDir();marker:=filepath.Join(work,"real-command-started");var requests, admitted atomic.Int32
   srv:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
    var content []any
    if requests.Add(1)==1 {
     if explained {content=append(content,map[string]any{"type":"text","text":"Create the local observation marker."})}
     content=append(content,map[string]any{"type":"tool_use","id":"review-tool","name":"run_command","input":map[string]any{"command":fmt.Sprintf("printf admitted > %q",marker),"ai_callback_delay":0.1}})
    } else {content=[]any{map[string]any{"type":"text","text":"done"}}}
    json.NewEncoder(w).Encode(map[string]any{"content":content,"usage":map[string]int{"input_tokens":0,"output_tokens":0}})
   }));defer srv.Close()
   root:=ensemble.New(log.New(io.Discard,"",0));defer root.Shutdown()
   a,err:=root.NewAgent(ensemble.Config{Vendor:ensemble.Anthropic,Model:"claude-sonnet-4-6",APIKey:"fixture",BaseURL:srv.URL,System:"Local contract check",MaxTokens:256,DisableStreaming:true,EnableTools:true,RequireToolExplanation:true,DataDir:work,Workspace:work});if err!=nil{t.Fatal(err)}
   root.Observe(func(o ensemble.AgentObservation){if o.Agent==a {if _,ok:=o.Event.(ensemble.ToolDispatched);ok{admitted.Add(1)}}})
   h:=newHub(root,a);p:=fixturePeer(256);h.subscribe(p);<-p.send
   ctx,cancel:=context.WithTimeout(context.Background(),5*time.Second);defer cancel()
   var result ensemble.Result
   select {case result= <-a.Submit(ctx,"perform the local marker call"):case <-ctx.Done():t.Fatal("local fixture did not complete")}
   if result.Err!=nil{t.Fatal(result.Err)}
   liveStarts, liveFinishes:=0,0;var errors []string
   for len(p.send)>0 {for _,m:=range <-p.send {if m["type"]=="tool_dispatched"{liveStarts++};if m["type"]=="tool_finished"{liveFinishes++;if m["is_error"]==true{errors=append(errors,fmt.Sprint(m["result"]))}}}}
   q:=fixturePeer(1);h.subscribe(q);replayStarts:=len(messages(<-q.send,"tool_dispatched"))
   _,statErr:=os.Stat(marker);effect:=statErr==nil
   facts:=map[string]any{"explained":explained,"provider_requests":requests.Load(),"core_tool_dispatched":admitted.Load(),"actual_command_effect":effect,"live_gui_tool_dispatched":liveStarts,"replay_gui_tool_dispatched":replayStarts,"live_gui_tool_finished":liveFinishes,"errors":errors}
   b,_:=json.Marshal(facts);fmt.Printf("REVIEW_FACTS %s\n",b)
   if explained {if admitted.Load()!=1||!effect||liveStarts!=1||replayStarts!=1{t.Errorf("positive admission/effect/projection mismatch: %s",b)}} else {
    if admitted.Load()!=0||effect{t.Fatalf("core explanation refusal did not hold: %s",b)}
    if liveFinishes!=1||len(errors)!=1{t.Errorf("missing truthful refused result: %s",b)}
    if liveStarts!=0||replayStarts!=0{t.Errorf("GUI reports a dispatch for a batch refused before admission: %s",b)}
   }
  })
 }
}
