using System.Text.Json;
using Ensemble;
using EnsembleGui;

// A fresh GUI must follow the operating system, so the served theme is "system".
// The Go settings tests pinned that default; the C# port regressed it to "dark",
// and theme.js then stored the server's guess as the browser's own preference,
// so an OS-light visitor was locked into dark and never recovered.
static class ThemeDefaultTests
{
    static void Check(bool value, string message)
    {
        if (!value) throw new InvalidOperationException(message);
    }

    public static async Task RunAsync()
    {
        string work = Path.Combine(Path.GetTempPath(), "gui-theme-" + Guid.NewGuid().ToString("N"));
        Directory.CreateDirectory(work);
        var fleet = new Ensemble.Ensemble(TextWriter.Null, dataDirectory: Path.Combine(work, "fleet"));
        var agent = new Agent(fleet, new Configuration(Vendor.OpenAI, "offline", "http://127.0.0.1:1"), name: "theming");
        await Task.CompletedTask;
        try
        {
            Check(new Preferences().Theme == "system", "The Preferences record no longer defaults to the system theme");

            // The path a running GUI actually reads: no file yet, so defaults stand.
            string settings = Path.Combine(work, "settings.json");
            Check(Gui.ReadPreferences(agent, settings).Theme == "system",
                "A fresh GUI without a settings file does not offer the system theme");

            // An explicit choice still wins, so this default is not a refusal to obey.
            Directory.CreateDirectory(Path.GetDirectoryName(settings)!);
            File.WriteAllText(settings, JsonSerializer.Serialize(new { theme = "dark" }));
            Check(Gui.ReadPreferences(agent, settings).Theme == "dark",
                "An explicitly stored dark theme was overridden by the new default");

            File.WriteAllText(settings, JsonSerializer.Serialize(new { theme = "nonsense" }));
            Check(Gui.ReadPreferences(agent, settings).Theme == "system",
                "An unrecognised stored theme did not fall back to the system theme");
        }
        finally
        {
            agent.Destroy();
            fleet.Destroy();
        }
        Console.WriteLine("PASS theme defaults: system default, disk precedence, invalid fallback");
    }
}
