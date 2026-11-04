using System.Net;
using System.Net.Sockets;
using System.Text;
using System.Text.Json;
using System.Text.Json.Nodes;
using Ensemble;

// What the endpoint says about a requested reasoning effort, not what a table
// entry guessed. A refusal stops the launch; silence never does.
static class ReasoningEffortTests
{
    static void Check(bool value, string message)
    {
        if (!value)
        {
            throw new Exception(message);
        }
    }

    // Records the request bodies it saw and answers the two endpoints the check
    // may use. advertise null means the model entry publishes no effort field at
    // all, which is "the server said nothing", not evidence of refusal.
    sealed class Provider : IDisposable
    {
        readonly HttpListener listener = new();
        public readonly List<string> Paths = [];
        public readonly List<string> Bodies = [];
        public string[]? Advertise = ["low", "medium", "high"];
        public Func<string, bool> Accept = _ => true;
        public string Refusal = "reasoning_effort must be one of none, minimal, low, medium, high, xhigh, max, off; got '{0}'";
        public bool ServeModels = true;

        public Provider(string model)
        {
            int port = Free();
            Prefix = $"http://127.0.0.1:{port}/";
            listener.Prefixes.Add(Prefix);
            listener.Start();
            for (int index = 0; index < 32; index++)
            {
                _ = Task.Run(() => Loop(model));
            }
        }

        public string Prefix { get; }

        int Free()
        {
            var probe = new TcpListener(IPAddress.Loopback, 0);
            probe.Start();
            int port = ((IPEndPoint)probe.LocalEndpoint).Port;
            probe.Stop();
            return port;
        }

        async Task Loop(string model)
        {
            while (listener.IsListening)
            {
                HttpListenerContext context;
                try
                {
                    context = await listener.GetContextAsync();
                }
                catch (Exception)
                {
                    return;
                }
                string path = context.Request.Url!.AbsolutePath;
                string body = new StreamReader(context.Request.InputStream,
                    context.Request.ContentEncoding).ReadToEnd();
                lock (Paths)
                {
                    Paths.Add(path);
                    Bodies.Add(body);
                }
                int status = 200;
                string payload;
                if (path.EndsWith("/models"))
                {
                    if (!ServeModels)
                    {
                        status = 404;
                        payload = "{\"error\":{\"message\":\"no route\"}}";
                    }
                    else
                    {
                        payload = Advertise is null
                            ? "{\"data\":[]}"
                            : "{\"data\":[{\"id\":\"" + model
                              + "\",\"supported_reasoning_efforts\":[\""
                              + string.Join("\",\"", Advertise) + "\"]}]}";
                    }
                }
                else
                {
                    string requested = "";
                    try
                    {
                        using var parsed = JsonDocument.Parse(body);
                        // Responses nests the level, chat-completions carries it flat.
                        if (parsed.RootElement.TryGetProperty("reasoning_effort", out var flat)
                            && flat.ValueKind == JsonValueKind.String)
                        {
                            requested = flat.GetString() ?? "";
                        }
                        else if (parsed.RootElement.TryGetProperty("reasoning", out var reasoning)
                            && reasoning.TryGetProperty("effort", out var nested)
                            && nested.ValueKind == JsonValueKind.String)
                        {
                            requested = nested.GetString() ?? "";
                        }
                    }
                    catch (JsonException)
                    {
                    }
                    if (Accept(requested))
                    {
                        payload = "{\"choices\":[{\"message\":{\"content\":\"ok\"}}]}";
                    }
                    else
                    {
                        status = 400;
                        payload = "{\"error\":{\"message\":\""
                            + string.Format(Refusal, requested)
                            + "\",\"type\":\"invalid_request_error\"}}";
                    }
                }
                byte[] bytes = Encoding.UTF8.GetBytes(payload);
                context.Response.StatusCode = status;
                context.Response.ContentType = "application/json";
                context.Response.ContentLength64 = bytes.Length;
                await context.Response.OutputStream.WriteAsync(bytes);
                context.Response.Close();
            }
        }

        public void Dispose()
        {
            try
            {
                listener.Stop();
                listener.Close();
            }
            catch (HttpListenerException)
            {
            }
        }
    }

    // A table row permits a level and is a real negative when it omits one. A
    // model the table does not name is silence, so "qwen-local-not-in-the-table"
    // asks and the endpoint answers. An unrequested effort still sends no traffic.
    const string Listed = "gpt-6-astra";
    const string Unlisted = "qwen-local-not-in-the-table";

    static Agent Agent(Ensemble.Ensemble root, string url, string effort, string name)
        => new(root, new Configuration(Vendor.OpenAI, Listed, url, "test-key",
            reasoningEffort: effort), name: name);

    static string Root(string temp, string name) => Path.Combine(temp, "effort-" + name);

    public static async Task Run(string temp)
    {
        // An empty request sends no traffic at all: ordinary launches and every
        // offline fixture stay free of extra probes.
        using (var quiet = new Provider(Listed))
        {
            var root = new Ensemble.Ensemble(TextWriter.Null, dataDirectory: Root(temp, "quiet"));
            var agent = new Agent(root, new Configuration(Vendor.OpenAI, Unlisted,
                quiet.Prefix, "test-key"), name: "quiet");
            await agent.Engine().ValidateReasoningEffortAsync(CancellationToken.None);
            lock (quiet.Paths)
            {
                Check(quiet.Paths.Count == 0,
                    "an unrequested effort still probed the endpoint: " + string.Join(",", quiet.Paths));
            }
        }

        // A row is a real negative and keeps refusing. Silence about a model is
        // not: the endpoint gets to answer for itself. A summary is a different
        // provider field with no probe behind it, so it stays table-gated.
        var tableRoot = new Ensemble.Ensemble(TextWriter.Null, dataDirectory: Root(temp, "table"));
        try
        {
            _ = new Agent(tableRoot, new Configuration(Vendor.OpenAI, Listed,
                "http://127.0.0.1:1/v1", "test-key", reasoningEffort: "not-a-level"),
                name: "table-value");
            Check(false, "a level outside the table row was accepted");
        }
        catch (ArgumentException)
        {
        }
        try
        {
            _ = new Agent(tableRoot, new Configuration(Vendor.OpenAI, Unlisted,
                "http://127.0.0.1:1/v1", "test-key", reasoningSummary: "auto"), name: "table-summary");
            Check(false, "an unlisted model was allowed to request a summary");
        }
        catch (ArgumentException)
        {
        }

        // An advertised level needs no probe: one request, no completion.
        using (var known = new Provider(Listed))
        {
            var root = new Ensemble.Ensemble(TextWriter.Null, dataDirectory: Root(temp, "known"));
            var agent = Agent(root, known.Prefix, "medium", "known");
            await agent.Engine().ValidateReasoningEffortAsync(CancellationToken.None);
            lock (known.Paths)
            {
                Check(known.Paths.Count == 1 && known.Paths[0].EndsWith("/models"),
                    "an advertised effort cost more than the listing: " + string.Join(",", known.Paths));
            }
        }

        // An advertisement is only a hint: at least one server publishes a strict
        // subset of what it accepts, so the value still gets to prove itself.
        using (var partial = new Provider(Listed))
        {
            partial.Advertise = ["low"];
            var root = new Ensemble.Ensemble(TextWriter.Null, dataDirectory: Root(temp, "partial"));
            var agent = Agent(root, partial.Prefix, "high", "partial");
            await agent.Engine().ValidateReasoningEffortAsync(CancellationToken.None);
            lock (partial.Paths)
            {
                Check(partial.Paths.Count == 2,
                    "an unadvertised but accepted effort was not given its chance");
            }
        }

        // An explicit refusal of this parameter stops the launch, with the server's
        // own words, because continuing would advertise a denied capability.
        using (var refused = new Provider(Listed))
        {
            refused.Accept = value => value != "xhigh";
            var root = new Ensemble.Ensemble(TextWriter.Null, dataDirectory: Root(temp, "refused"));
            var agent = Agent(root, refused.Prefix, "xhigh", "refused");
            string message = "";
            try
            {
                await agent.Engine().ValidateReasoningEffortAsync(CancellationToken.None);
            }
            catch (InvalidOperationException failure)
            {
                message = failure.Message;
            }
            Check(message.Contains("reasoning_effort", StringComparison.OrdinalIgnoreCase)
                && message.Contains("xhigh", StringComparison.Ordinal),
                "the provider's refusal did not reach the operator: " + message);
        }

        // A refusal of something else is not evidence about effort levels.
        using (var unrelated = new Provider(Listed))
        {
            unrelated.Accept = _ => false;
            unrelated.Refusal = "Insufficient credits, please top up your account";
            var root = new Ensemble.Ensemble(TextWriter.Null, dataDirectory: Root(temp, "unrelated"));
            var agent = Agent(root, unrelated.Prefix, "high", "unrelated");
            await agent.Engine().ValidateReasoningEffortAsync(CancellationToken.None);
        }

        // Silence from the endpoint is accepted with a note, never treated as a
        // denial: ensemble already starts against an offline provider.
        var quietNotes = new StringWriter();
        using (var silent = new Provider(Listed))
        {
            silent.ServeModels = false;
            silent.Accept = _ => false;
            silent.Refusal = "The requested endpoint is not available on this deployment";
            var root = new Ensemble.Ensemble(quietNotes, dataDirectory: Root(temp, "silent"));
            var agent = Agent(root, silent.Prefix, "xhigh", "silent");
            await agent.Engine().ValidateReasoningEffortAsync(CancellationToken.None);
        }
        var deadNotes = new StringWriter();
        var deadRoot = new Ensemble.Ensemble(deadNotes, dataDirectory: Root(temp, "dead"));
        var deadAgent = Agent(deadRoot, "http://127.0.0.1:1/v1", "high", "dead");
        await deadAgent.Engine().ValidateReasoningEffortAsync(CancellationToken.None);

        // A configured level reaches the wire, so the check is not theatre.
        var renderRoot = new Ensemble.Ensemble(TextWriter.Null, dataDirectory: Root(temp, "render"));
        var requested = Agent(renderRoot, "http://127.0.0.1:1/v1", "high", "render");
        Check(requested.Engine().Render(new Context())["reasoning"]?["effort"]
            ?.GetValue<string>() == "high", "the requested effort never reached the request body");

        // The case the table cannot speak about: an unlisted model on the
        // chat-completions dialect must be able to ask, and an accepted level has
        // to reach the wire in that dialect's own shape. A server that answers
        // nothing about effort is silence, which must not stop the run.
        using (var qwen = new Provider(Unlisted))
        {
            qwen.Advertise = null;
            qwen.Accept = value => value == "high";
            var root = new Ensemble.Ensemble(TextWriter.Null, dataDirectory: Root(temp, "qwen"));
            var asked = new Agent(root, new Configuration(Vendor.OpenAI, Unlisted, qwen.Prefix,
                "test-key", reasoningEffort: "high"), name: "qwen");
            await asked.Engine().ValidateReasoningEffortAsync(CancellationToken.None);
            Check(asked.Engine().Render(new Context())["reasoning_effort"]?.GetValue<string>() == "high",
                "an accepted effort never reached the chat-completions body");
            asked.Destroy();
        }
        // The same endpoint refusing this exact parameter stops the launch with its
        // own words, which is the whole reason the table may stay silent.
        using (var denying = new Provider(Unlisted))
        {
            denying.Advertise = null;
            denying.Accept = value => value != "ultra";
            var root = new Ensemble.Ensemble(TextWriter.Null, dataDirectory: Root(temp, "denying"));
            var asked = new Agent(root, new Configuration(Vendor.OpenAI, Unlisted, denying.Prefix,
                "test-key", reasoningEffort: "ultra"), name: "denying");
            string message = "";
            try
            {
                await asked.Engine().ValidateReasoningEffortAsync(CancellationToken.None);
            }
            catch (InvalidOperationException failure)
            {
                message = failure.Message;
            }
            Check(message.Contains("ultra", StringComparison.Ordinal)
                && message.Contains("reasoning_effort", StringComparison.OrdinalIgnoreCase),
                "an unlisted endpoint's refusal did not stop the launch: " + message);
            asked.Destroy();
        }
    }
}
