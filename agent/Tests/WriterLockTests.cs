using Ensemble;
using Application = Ensemble.Ensemble;

internal static class WriterLockTests
{
    private static void Check(bool value, string message)
    {
        if (!value)
        {
            throw new Exception(message);
        }
    }

    private static Agent Actor(Application app, string home)
    {
        return new Agent(app, new Configuration(Vendor.Anthropic, "model-a",
            "http://127.0.0.1:1", "SECRET-KEY"), dataDirectory: home);
    }

    public static void SingleWriter()
    {
        string root = Path.Combine(Path.GetTempPath(), "writer-lock-" + Guid.NewGuid().ToString("N"));
        string home = Path.Combine(root, "agents", "root");
        Directory.CreateDirectory(home);
        var app = new Application(TextWriter.Null);
        string save = Path.Combine(home, "save.json");
        string lockPath = Path.Combine(home, ".writer.lock");
        try
        {
            var first = Actor(app, home);
            first.StartPersistence(new PersistenceOptions(save));
            Check(File.Exists(lockPath), "StartPersistence must hold writer authority beside the save file");

            // The held handle is the lock itself, so a second writer over the same
            // actor home must be refused while the first Agent is alive.
            var rival = Actor(app, home);
            WriterConflictException? refusal = null;
            try
            {
                rival.StartPersistence(new PersistenceOptions(save));
            }
            catch (WriterConflictException error)
            {
                refusal = error;
            }
            catch (Exception error)
            {
                throw new Exception($"expected refusal, got {error.GetType().Name}: {error.Message}");
            }

            Check(refusal is not null, "a second writer over one actor home must be refused");
            Check(refusal!.Message.Contains(home, StringComparison.Ordinal),
                "the refusal must name the contested actor home");

            first.Destroy();
            using (var released = new FileStream(lockPath, FileMode.Open, FileAccess.ReadWrite, FileShare.None))
            {
                Check(released.CanWrite, "Destroy must release writer authority");
            }

            var next = Actor(app, home);
            next.StartPersistence(new PersistenceOptions(save));
            next.Destroy();
        }
        finally
        {
            app.Destroy();
            try
            {
                Directory.Delete(root, true);
            }
            catch (IOException)
            {
                // A held diagnostic handle must not fail an otherwise sound check.
            }
        }
    }
}
