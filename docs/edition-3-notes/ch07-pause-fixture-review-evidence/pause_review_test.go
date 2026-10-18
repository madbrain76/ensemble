package gradeedition2
import("encoding/json";"fmt";"os";"testing";"github.com/waywardgeek/ensemble/internal/grade")
func TestOwnedPauseDiagnostic(t *testing.T){
 messages,err:=driveGUIPause(os.Getenv("DIAG_GUI_BIN"),os.Getenv("DIAG_GUI_ASSETS"),"claude-sonnet-4-6")
 r:=&grade.Ch8Result{BuildOK:true,PauseOK:err==nil,PauseMessages:oldMessages(messages)}
 if err!=nil{r.PauseErr=err.Error()}
 c:=grade.Ch8Evaluate(r)[3:4];report:=grade.NewTitledReport("Repaired pause evaluator; review only",c,"");b,_:=json.Marshal(report);fmt.Printf("PAUSE_REPORT %s\n",b)
 if !report.Passed{t.Fail()}
}
