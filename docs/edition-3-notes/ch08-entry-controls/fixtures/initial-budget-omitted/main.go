package main
import("encoding/json";"flag";"net/http";"os";"strconv";"sync";"github.com/gorilla/websocket")
const mode = "initial-budget-omitted"
func main(){
 port:=flag.Int("port",0,""); flag.String("gui-dir","",""); flag.Parse()
 state:=map[string]any{"theme":"dark","tts_enabled":true,"tts_speed":float64(1),"font_size":float64(16),"max_tool_rounds":float64(0)}
 if mode=="initial-budget-omitted" {delete(state,"max_tool_rounds")}
 if mode!="no-persistence" {if b,e:=os.ReadFile("mutant-settings.json");e==nil {json.Unmarshal(b,&state)}}
 var mu sync.Mutex; clients:=map[*websocket.Conn]bool{}; up:=websocket.Upgrader{CheckOrigin:func(*http.Request)bool{return true}}
 http.HandleFunc("/",func(w http.ResponseWriter,r *http.Request){w.Write([]byte("owned settings-only mutant"))})
 http.HandleFunc("/ws",func(w http.ResponseWriter,r *http.Request){
 c,e:=up.Upgrade(w,r,nil);if e!=nil{return}; defer func(){mu.Lock();delete(clients,c);mu.Unlock();c.Close()}()
 for {var m struct {Type string `json:"type"`;Settings map[string]any `json:"settings"`}; if c.ReadJSON(&m)!=nil{return};mu.Lock()
 switch m.Type {case "subscribe":clients[c]=true;c.WriteJSON(map[string]any{"type":"current_settings","settings":state})
 case "update_settings":for k,v:=range m.Settings {state[k]=v};if mode!="no-persistence" {b,_:=json.Marshal(state);os.WriteFile("mutant-settings.json",b,0600)};for peer:=range clients {peer.WriteJSON(map[string]any{"type":"settings_changed","settings":state})}}
 mu.Unlock() }
 })
 if http.ListenAndServe("127.0.0.1:"+strconv.Itoa(*port),nil)!=nil {os.Exit(2)}
}
