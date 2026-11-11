namespace Ensemble;

public sealed partial class Memory
{
    internal partial AgentSettings? ReadSettings()
    {
        lock (gate)
        {
            string settingsPath = HomePath("settings.json");
            try
            {
                return AgentSettings.Read(Parent().Sandbox(), settingsPath);
            }
            catch (InvalidDataException failure)
            {
                // The launch handler reports failures by exception type only, so a
                // rejected document must name its own file in the diagnostics log or
                // "check private diagnostics" leaves the operator nowhere to look.
                Parent().Log($"Agent settings rejected at {settingsPath}: {failure.Message}");
                throw;
            }
        }
    }

    internal partial void WriteSettings(AgentSettings settings)
    {
        lock (gate)
        {
            settings.Write(Parent().Sandbox(), HomePath("settings.json"));
        }
    }
}
