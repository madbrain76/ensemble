using Ensemble;
using Application = Ensemble.Ensemble;

// An unconfined Agent with no caller-supplied directory invents its own scratch
// beside Path.GetTempPath(). OutputDirectory is readonly and the next Jobs draws a
// fresh GUID beside it, so after a restart nothing in the program can name the old
// directory again: it is unreachable by construction, not merely unused. Left in
// place, every Agent ever built leaves one behind, which reached 27,000 entries in
// /tmp on a machine that runs the agent continuously. Destruction is the last point
// at which the owner exists, so that is where the scratch goes with it.
internal static partial class JobScratchTests
{
    public static void TemporaryScratchIsReclaimed()
    {
        string home = Home("reclaim");
        string previous = Directory.GetCurrentDirectory();
        Directory.SetCurrentDirectory(home);
        try
        {
            var app = new Application(TextWriter.Null, dataDirectory: home);
            // No name and no data directory: this is the one construction whose Jobs
            // path is invented here rather than derived from a home.
            var agent = new Agent(app, Offline());
            string directory = agent.Jobs().OutputDirectory;
            Check(Directory.Exists(directory), "the Agent never created its scratch directory");
            Check(directory.StartsWith(Path.GetTempPath(), StringComparison.Ordinal)
                && Path.GetFileName(directory).StartsWith("ensemble-jobs-", StringComparison.Ordinal),
                "the scratch was not the invented temporary default, so this case would not "
                    + "have tested reclamation at all: " + directory);
            // Retained output as a completed job would have left it.
            File.WriteAllText(Path.Combine(directory, "1"), "retained bytes");
            app.Destroy();
            Check(!Directory.Exists(directory),
                "destroying the Agent left its scratch directory and its output behind: "
                    + directory);
        }
        finally
        {
            Directory.SetCurrentDirectory(previous);
            Delete(home);
        }
    }

    // The rule is ownership, not location. A directory somebody named is retained
    // output, exactly as a confined Agent's own root keeps its own; deleting either
    // would destroy evidence the operator asked for.
    public static void SuppliedDirectoryIsRetained()
    {
        string home = Home("supplied");
        string supplied = Path.Combine(home, "kept-jobs");
        string previous = Directory.GetCurrentDirectory();
        Directory.SetCurrentDirectory(home);
        try
        {
            var app = new Application(TextWriter.Null, dataDirectory: home);
            var agent = new Agent(app, Offline(), name: "keeper", outputDirectory: supplied);
            string directory = agent.Jobs().OutputDirectory;
            Check(Directory.Exists(directory), "a supplied output directory was not created");
            File.WriteAllText(Path.Combine(directory, "1"), "keep me");
            app.Destroy();
            Check(Directory.Exists(directory)
                && File.ReadAllText(Path.Combine(directory, "1")) == "keep me",
                "destroying the Agent deleted output the caller had chosen to keep");
        }
        finally
        {
            Directory.SetCurrentDirectory(previous);
            Delete(home);
        }
    }

    public static void ConfinedScratchIsRetained()
    {
        string home = Home("confined");
        string root = Path.Combine(home, "sandbox");
        Directory.CreateDirectory(root);
        string previous = Directory.GetCurrentDirectory();
        Directory.SetCurrentDirectory(home);
        try
        {
            var app = new Application(TextWriter.Null);
            var configuration = new Configuration(Vendor.OpenAI, "gpt-4.1",
                "http://127.0.0.1:1", sandboxRoot: root);
            var agent = new Agent(app, configuration);
            string directory = agent.Jobs().OutputDirectory;
            Check(directory.StartsWith(root + Path.DirectorySeparatorChar,
                StringComparison.Ordinal),
                "a confined Agent's scratch escaped its root: " + directory);
            File.WriteAllText(Path.Combine(directory, "1"), "confined bytes");
            app.Destroy();
            Check(Directory.Exists(directory)
                && File.ReadAllText(Path.Combine(directory, "1")) == "confined bytes",
                "destroying a confined Agent removed the retained output under its root");
        }
        finally
        {
            Directory.SetCurrentDirectory(previous);
            Delete(home);
        }
    }

    private static Configuration Offline() =>
        new(Vendor.OpenAI, "gpt-4.1", "http://127.0.0.1:1");

    private static string Home(string kind)
    {
        string home = Path.Combine(Path.GetTempPath(),
            "ensemble-scratch-" + kind + "-" + Guid.NewGuid().ToString("N"));
        Directory.CreateDirectory(home);
        return home;
    }

    private static void Delete(string path)
    {
        try
        {
            Directory.Delete(path, recursive: true);
        }
        catch (IOException)
        {
        }
    }

    private static void Check(bool value, string message)
    {
        if (!value)
        {
            throw new InvalidOperationException(message);
        }
    }
}
