from pathlib import Path
R=Path('/Users/bill/projects/ensemble');B=Path('/tmp/cadence-ch08-adapter-before')
def rep(s,a,b,n=1):
 assert s.count(a)==n,(a,s.count(a));return s.replace(a,b)
p=R/'cmd/grade-edition3/main.go';s=p.read_text();assert s==(B/'cmd/grade-edition3/main.go').read_text()
s=rep(s,'(1, 4, 5, 6 or 7)','(1, 4, 5, 6, 7 or 8)')
for a,b,n in [('Chapter 4/5/6/7','Chapter 4/5/6/7/8',3),('Chapter 5/6/7','Chapter 5/6/7/8',2),('Chapter 6/7','Chapter 6/7/8',2),('Chapter 7 separate','Chapter 7/8 separate',1),('Chapter 7 declared','Chapter 7/8 declared',1),('-ch 4|5|6|7 ','-ch 4|5|6|7|8 ',1)]:s=rep(s,a,b,n)
s=rep(s,'*chapter != 7','*chapter != 7 && *chapter != 8',3)
for a in ['if *chapter == 4 || *chapter == 5 || *chapter == 6 || *chapter == 7 {','if *chapter == 5 || *chapter == 6 || *chapter == 7 {','if *chapter == 6 || *chapter == 7 {']:
 n=2 if a.startswith('if *chapter == 5') else 1
 s=rep(s,a,a[:-2]+' || *chapter == 8 {',n)
s=rep(s,'\t\tif *chapter == 7 {','\t\tif *chapter == 7 || *chapter == 8 {')
s=rep(s,'Chapters 6/7','Chapters 6/7/8',2);s=rep(s,'Chapters 5/6/7','Chapters 5/6/7/8');s=rep(s,'required for Chapter 7','required for Chapters 7/8');s=rep(s,'gui for Chapter 7)','gui for Chapters 7/8)')
a='\t} else {\n\t\treport, err = gradeedition3.Chapter1(bin)';s=rep(s,a,'\t} else if *chapter == 8 {\n\t\treport, err = gradeedition2.Chapter8(bin, *core, *gui, *guiModel, *consumer, *workflow, *streamExercise, *streamModel, example, models)\n'+a);p.write_text(s)
p=R/'cmd/grade-edition3/main_test.go';s=p.read_text();s=rep(s,'{}, {""}, {"-ch", "2", target},','{}, {""}, {"-ch", "2", target}, {"-ch", "9", target},');s=rep(s,'{"-ch", "7", "-h", target}}','{"-ch", "7", "-h", target}, {"-ch", "8", "-h", target}}');s=rep(s,'func TestChapter7ValidationBeforeLaunch(t *testing.T) {','''func TestChapter7ValidationBeforeLaunch(t *testing.T) {
 testGUIChapterValidationBeforeLaunch(t, "7")
}

func TestChapter8ValidationBeforeLaunch(t *testing.T) {
 testGUIChapterValidationBeforeLaunch(t, "8")
}

func testGUIChapterValidationBeforeLaunch(t *testing.T, chapter string) {''');s=rep(s,'base := []string{"-ch", "7",','base := []string{"-ch", chapter,');p.write_text(s)
p=R/'internal/gradeedition2/ch08.go';s=p.read_text();s=rep(s,'[]string{"theme", "font_size", "tts_enabled", "tts_speed"}','[]string{"theme", "font_size", "tts_enabled", "tts_speed", "max_tool_rounds"}');p.write_text(s)
