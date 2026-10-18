package main
import "os"
func main() { f,e:=os.OpenFile(os.Getenv("REVIEW_SENTINEL_MARKER"),os.O_CREATE|os.O_WRONLY|os.O_APPEND,0600);if e!=nil {panic(e)}; f.WriteString("EXECUTED\n");f.Close();os.Exit(99) }
