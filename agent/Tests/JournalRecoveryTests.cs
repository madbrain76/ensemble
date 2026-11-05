using Ensemble;
using Application = Ensemble.Ensemble;

internal static class JournalRecoveryTests
{
    private static void Check(bool value, string message)
    {
        if (!value)
        {
            throw new Exception(message);
        }
    }

    private static string DirectoryForTest()
    {
        string directory = Path.Combine(Path.GetTempPath(), "journal-recovery-" + Guid.NewGuid().ToString("N"));
        Directory.CreateDirectory(directory);
        return directory;
    }

    private static Agent Owner(Application app)
    {
        return new Agent(app, new Configuration(Vendor.OpenAI, "offline", "http://127.0.0.1:1"));
    }

    private static Event Message(string text)
    {
        return new Event(EventKind.MessageReceived, actor: Actor.Human, parts: [new TextPart(text)]);
    }

    private static string Text(History history)
    {
        return string.Join("|", history.View().Dialogue.SelectMany(entry => entry.Parts)
            .OfType<TextPart>().Select(part => part.Text));
    }

    // Two launchers over one actor home numbered their events independently, so the
    // recovery journal stopped being contiguous and every later load failed.
    public static void RejectedTail()
    {
        string directory = DirectoryForTest();
        string path = Path.Combine(directory, "save.json");
        string journal = path + ".journal.jsonl";
        var log = new StringWriter();
        var app = new Application(log);
        try
        {
            var source = Owner(app).History();
            source.StartJournal(path);
            source.Append(Message("one"));
            source.Save(new PersistenceOptions(path, 0));
            source.Append(Message("two"));
            string[] committed = File.ReadAllLines(journal);
            Check(committed.Length > 0, "Journal fixture wrote no records");
            File.AppendAllText(journal, committed[^1] + "\n");

            var restored = Owner(app).History();
            Check(restored.LoadSave(path), "a rejected journal must still restore the saved snapshot");
            Check(Text(restored) == "one|two",
                "the snapshot plus every contiguous journal record must survive the damaged tail");

            string written = log.ToString();
            Check(written.Contains(journal, StringComparison.Ordinal)
                && written.Contains("rejected", StringComparison.Ordinal),
                "the log must name the rejected journal and the refusal");
            Check(written.Contains("continued from the saved snapshot", StringComparison.Ordinal)
                && written.Contains("2 journal record(s) before the damage", StringComparison.Ordinal),
                "the log must say what recovery salvaged instead of failing");

            string[] preserved = Directory.GetFiles(directory, "*.rejected-*");
            Check(preserved.Length == 1, "the rejected records must be preserved exactly once");
            Check(File.ReadAllLines(preserved[0]).Length == committed.Length + 1,
                "the preserved file must hold every original record, including the duplicate");
            Check(!File.Exists(journal), "the rejected file must stop being the active journal");

            restored.Append(Message("three"));
            restored.Save(new PersistenceOptions(path, 0));
            var again = Owner(app).History();
            Check(again.LoadSave(path) && Text(again) == "one|two|three",
                "later saves must stay recoverable after a rejected tail");
        }
        finally
        {
            app.Destroy();
        }
    }

    // With no usable snapshot the journal is the only surviving record, so it must
    // never be moved aside: the load refuses and every file stays exactly as found.
    public static void RefusesWithoutSnapshot()
    {
        string directory = DirectoryForTest();
        string path = Path.Combine(directory, "save.json");
        string journal = path + ".journal.jsonl";
        var log = new StringWriter();
        var app = new Application(log);
        try
        {
            var facts = new Event[] { Message("first"), Message("second"), Message("second again") };
            long seq = 0;
            using (var seed = new StreamWriter(journal))
            {
                foreach (Event fact in facts)
                {
                    seq++;
                    seed.WriteLine(JsonCodec.WriteEvent(fact.Seq == 0 ? fact.Numbered(seq, DateTimeOffset.UtcNow) : fact));
                }

                // Repeat the last sequence: the record shape a second writer produced.
                seed.WriteLine(JsonCodec.WriteEvent(facts[^1].Numbered(seq, DateTimeOffset.UtcNow)));
            }

            var history = Owner(app).History();
            try
            {
                history.LoadSave(path);
                throw new Exception("an unusable only record must be refused");
            }
            catch (InvalidDataException)
            {
            }

            Check(File.Exists(journal), "the only surviving record must be left in place");
            Check(Directory.GetFiles(directory, "*.rejected-*").Length == 0,
                "nothing may be quarantined when no snapshot remains");
            Check(log.ToString().Contains("files left untouched", StringComparison.Ordinal),
                "the log must say that the files were not modified");
        }
        finally
        {
            app.Destroy();
        }
    }
}
