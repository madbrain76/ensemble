from pathlib import Path
R=Path('/Users/bill/projects/ensemble');B=Path('/tmp/cadence-ch08-adapter-before');B.mkdir(exist_ok=False)
paths=['cmd/grade-edition3/main.go','cmd/grade-edition3/main_test.go','internal/gradeedition2/ch08.go']
for n in paths:
 p=B/n;p.parent.mkdir(parents=True,exist_ok=True);p.write_bytes((R/n).read_bytes())
def replace(s,a,b,n=1):
 assert s.count(a)==n,(a,s.count(a));return s.replace(a,b)
p=R/paths[0];s=p.read_text();s=replace(s,'1, 4, 5, 6, or 7','1, 4, 5, 6, 7, or 8',2);s=replace(s,'*ch != 6 && *ch != 7','*ch != 6 && *ch != 7 && *ch != 8');s=replace(s,'*ch != 7 && (seen["gui"]','*ch != 7 && *ch != 8 && (seen["gui"]');s=replace(s,'only for Chapter 7','only for Chapters 7 and 8');s=replace(s,'*ch == 7 && (strings.TrimSpace','(*ch == 7 || *ch == 8) && (strings.TrimSpace');s=replace(s,'"Chapter 7 requires -gui and -gui-model"','"Chapter", *ch, "requires -gui and -gui-model"');s=replace(s,'*ch == 6 || *ch == 7 {','*ch == 6 || *ch == 7 || *ch == 8 {');s=replace(s,'if *ch == 7 {','if *ch == 7 || *ch == 8 {');a='case 7:\n\t\trep, err = gradeedition2.Chapter7(bin, *core, *gui, *guiModel, *consumer, *workflow, *streamExercise, *streamModel, example, models)';s=replace(s,a,a+'\n\tcase 8:\n\t\trep, err = gradeedition2.Chapter8(bin, *core, *gui, *guiModel, *consumer, *workflow, *streamExercise, *streamModel, example, models)');p.write_text(s)
p=R/paths[1];s=p.read_text();s=replace(s,'for _, ch := range []string{"2", "8"}','for _, ch := range []string{"2", "9"}');s=replace(s,'1, 4, 5, 6, or 7','1, 4, 5, 6, 7, or 8');s=replace(s,'func TestEdition3Chapter7GuardsDoNotLaunchTarget(t *testing.T) {','''func TestEdition3Chapter7GuardsDoNotLaunchTarget(t *testing.T) {
 testGUIChapterGuardsDoNotLaunchTarget(t, "7")
}

func TestEdition3Chapter8GuardsDoNotLaunchTarget(t *testing.T) {
 testGUIChapterGuardsDoNotLaunchTarget(t, "8")
}

func testGUIChapterGuardsDoNotLaunchTarget(t *testing.T, chapter string) {''');s=replace(s,'[]string{"-ch", "7"}, args...','[]string{"-ch", chapter}, args...');p.write_text(s)
p=R/paths[2];s=p.read_text();s=replace(s,'[]string{"theme", "font_size", "tts_enabled", "tts_speed"}','[]string{"theme", "font_size", "tts_enabled", "tts_speed", "max_tool_rounds"}');p.write_text(s)
