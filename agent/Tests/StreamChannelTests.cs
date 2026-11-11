using Ensemble;

// Per-channel delivery selection: what may stream live, without ever hiding
// content from the settled answer.
static class StreamChannelTests
{
    static void Check(bool value, string message)
    {
        if (!value) throw new Exception(message);
    }

    public static void Run()
    {
        Check(Engine.ParseStream("all") == StreamChannel.All, "'all' did not select every channel");
        Check(Engine.ParseStream("text,thinking,tools") == StreamChannel.All,
            "the explicit list did not equal all");
        Check(Engine.ParseStream("text") == StreamChannel.Text, "one channel was not honoured");
        Check(Engine.ParseStream("thinking, tools") == (StreamChannel.Thinking | StreamChannel.Tools),
            "whitespace in a list was not tolerated");
        Check(Engine.ParseStream(" THINKING ") == StreamChannel.Thinking,
            "case and padding were not tolerated");
        Check(Engine.ParseStream("all,text") == StreamChannel.All,
            "all combined with a channel should stay all");
        foreach (string bad in new[] { "", "   ", "text,,thinking", "banana", "text,banana" })
        {
            try
            {
                _ = Engine.ParseStream(bad);
                Check(false, $"accepted the streaming selection \"{bad}\"");
            }
            catch (ArgumentException error)
            {
                Check(bad.Trim().Length == 0
                        ? error.Message.Contains("empty streaming selection")
                        : error.Message.Contains("invalid streaming channel"),
                    $"the refusal for \"{bad}\" did not explain the accepted values: {error.Message}");
            }
        }
        // The predicate the streaming path consults.
        Check(Engine.StreamSelected(StreamChannel.All, ChunkKind.Reply)
            && Engine.StreamSelected(StreamChannel.All, ChunkKind.Thinking)
            && Engine.StreamSelected(StreamChannel.All, ChunkKind.ToolCall),
            "'all' did not select every kind");
        Check(Engine.StreamSelected(StreamChannel.Text, ChunkKind.Reply)
            && !Engine.StreamSelected(StreamChannel.Text, ChunkKind.Thinking)
            && !Engine.StreamSelected(StreamChannel.Text, ChunkKind.ToolCall),
            "a text-only selection leaked another channel");
        Check(!Engine.StreamSelected(StreamChannel.Thinking, ChunkKind.Reply)
            && Engine.StreamSelected(StreamChannel.Thinking, ChunkKind.Thinking),
            "a thinking-only selection was not symmetric");
        Check(Engine.StreamSelected(StreamChannel.None, ChunkKind.Reply) == false,
            "no selection should stream nothing");
        // Configuration carries the choice, and the default streams everything.
        var defaultConfig = new Configuration(Vendor.OpenAI, "gpt-4o", "http://127.0.0.1:1");
        Check(defaultConfig.Streams == StreamChannel.All,
            "the default configuration stopped streaming some channel");
        var chosen = new Configuration(Vendor.OpenAI, "gpt-4o", "http://127.0.0.1:1",
            stream: "text,tools");
        Check(chosen.Streams == (StreamChannel.Text | StreamChannel.Tools),
            "the configuration lost the selected channels");
        Check(chosen.WithModel("gpt-4o-mini").Streams == chosen.Streams,
            "a model switch dropped the delivery selection");
        // Selecting a configured offering changes the route, not what the operator
        // asked to watch. An offering projection that copies route facts one by
        // once once omitted this one, and because an unset selection reads as "no
        // channel", every configured Agent stopped streaming live while still
        // settling correctly — a stall that looked like a slow model.
        var root = new global::Ensemble.Ensemble(TextWriter.Null);
        var carrierVendor = new ModelVendor(root, "carrier", "Carrier fixture");
        var carrierAuth = new ModelAuthentication(root, "none");
        var carrierEndpoint = new ModelEndpoint(carrierVendor, "chat",
            "http://127.0.0.1:1", Surface.ChatCompletions, carrierAuth);
        var carrierModel = new Model(carrierEndpoint, "small", "carrier-small", "Small",
            new ModelFeatures(Media.None, Stream: true, ContextWindow: 262144), 4096);
        Check(new Configuration(defaultConfig, carrierModel).Streams == StreamChannel.All,
            "a configured offering lost the default delivery selection");
        Check(new Configuration(chosen, carrierModel).Streams
                == (StreamChannel.Text | StreamChannel.Tools),
            "a configured offering lost the operator's delivery selection");
    }
}
