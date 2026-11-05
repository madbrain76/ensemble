using System.Text;
using System.Text.Json;
using System.Text.Json.Nodes;

namespace Ensemble;

public sealed partial class Memory
{
    internal partial AgentSettings? ReadSettings()
    {
        lock (gate)
        {
            // Only absence means defaults. Corruption, confinement refusal and
            // unreadable files must not silently enable a disabled band or recall.
            string settingsPath = HomePath("settings.json");
            // The launch handler reports failures by exception type only, so a
            // rejected document must name its own file in the diagnostics log or
            // "check private diagnostics" leaves the operator nowhere to look.
            InvalidDataException Reject(string reason)
            {
                Parent().Log($"Agent settings rejected at {settingsPath}: {reason}");
                return new InvalidDataException(reason);
            }
            Stream stream;
            try
            {
                stream = Parent().Sandbox().OpenFile(settingsPath, FileMode.Open, FileAccess.Read);
            }
            catch (FileNotFoundException)
            {
                return null;
            }
            catch (DirectoryNotFoundException)
            {
                return null;
            }
            using (stream)
            {
                byte[] bytes = new byte[16 * 1024 + 1];
                int length = stream.ReadAtLeast(bytes, bytes.Length, throwOnEndOfStream: false);
                if (length == bytes.Length)
                {
                    throw Reject("Agent settings exceed 16 KiB");
                }
                using var document = JsonDocument.Parse(bytes.AsMemory(0, length));
                var root = document.RootElement;
                if (root.ValueKind != JsonValueKind.Object)
                {
                    throw Reject("Invalid Agent settings document");
                }
                var names = root.EnumerateObject().Select(field => field.Name).Order(StringComparer.Ordinal).ToArray();
                // Chapter 22's complete memory-only document remains readable.
                // Supplied recall is strict; null/malformed never means defaults.
                bool legacy = names.SequenceEqual(new[] { "memory", "version" });
                bool complete = names.SequenceEqual(new[] { "memory", "recall", "version" });
                if ((!legacy && !complete) || root.GetProperty("version").GetInt32() != 1)
                {
                    throw Reject("Invalid Agent settings document");
                }
                var memory = JsonCodec.ReadMemorySettings(root.GetProperty("memory"));
                var recall = legacy ? new RecallSettings()
                    : JsonCodec.ReadRecallSettings(root.GetProperty("recall"));
                return new AgentSettings(memory, recall);
            }
        }
    }

    internal partial void WriteSettings(AgentSettings settings)
    {
        lock (gate)
        {
            var document = new JsonObject
            {
                ["version"] = 1,
                ["memory"] = JsonCodec.WriteMemorySettings(settings.Memory),
                ["recall"] = JsonCodec.WriteRecallSettings(settings.Recall)
            };
            byte[] bytes = Encoding.UTF8.GetBytes(document.ToJsonString());
            if (bytes.Length > 16 * 1024)
            {
                throw new InvalidDataException("Agent settings exceed 16 KiB");
            }
            Parent().Sandbox().ReplaceFile(HomePath("settings.json"), bytes);
        }
    }
}
