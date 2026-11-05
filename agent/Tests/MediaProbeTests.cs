using System.Net;
using System.Text;
using System.Text.Json;
using System.Text.Json.Nodes;
using Ensemble;
using Application = Ensemble.Ensemble;

// A curated capability table cannot name every local or newly released model.
// These cases run a real Engine against a fake provider and then ask the real
// admission validator the same question the GUI image route asks it.
internal static class MediaProbeTests
{
    private const string RedPng = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAIAAACQd1PeAAAAEElEQVR4nGL6z8AACAAA//8DCQECWLbVUAAAAABJRU5ErkJggg==";
    private const string GreenPng = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAIAAACQd1PeAAAAEElEQVR4nGJi+M8ACAAA//8CCgECX2t3gQAAAABJRU5ErkJggg==";
    private const string BluePng = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAIAAACQd1PeAAAAEElEQVR4nGJiYPgPCAAA//8BCwEC1pVDqAAAAABJRU5ErkJggg==";

    private static void Check(bool condition, string message)
    {
        if (!condition) throw new Exception(message);
    }

    private static bool Refused(Action action)
    {
        try { action(); } catch (NotSupportedException) { return true; }
        return false;
    }

    private sealed class Provider : HttpMessageHandler
    {
        public readonly List<string> Bodies = [];
        public string Name = "provider";
        public Func<string, string> Answer = _ => "I cannot see any image.";
        public string Reasoning = "";
        public HttpStatusCode Status = HttpStatusCode.OK;

        protected override async Task<HttpResponseMessage> SendAsync(HttpRequestMessage request, CancellationToken cancellation)
        {
            string body = await request.Content!.ReadAsStringAsync(cancellation);
            Bodies.Add(body);
            if (Status != HttpStatusCode.OK)
            {
                return new HttpResponseMessage(Status) { Content = new StringContent("{}") };
            }
            Check(request.Method == HttpMethod.Post
                && request.RequestUri!.AbsolutePath == "/v1/chat/completions", "probe used another route " + Name);
            // Compare parsed JSON: an escaped base64 '+' is the same bytes on the
            // wire, so a raw substring match would be a flaky assertion.
            using var document = JsonDocument.Parse(body);
            var root = document.RootElement;
            Check(root.GetProperty("max_tokens").GetInt32() >= 64, "probe left no answer budget " + Name);
            string url = root.GetProperty("messages")[0].GetProperty("content")[1]
                .GetProperty("image_url").GetProperty("url").GetString() ?? "";
            Check(url.StartsWith("data:image/png;base64,", StringComparison.Ordinal),
                "probe sent no inline image " + Name);
            string png = url["data:image/png;base64,".Length..];
            string asked = png == RedPng ? "red" : png == GreenPng ? "green" : png == BluePng ? "blue" : "";
            Check(asked.Length > 0, "probe sent an unknown solid-color image " + Name);
            var message = new JsonObject { ["role"] = "assistant" };
            if (Reasoning.Length > 0)
            {
                // A reasoning-only reply is the case the probe must refuse: the
                // color was produced, but never delivered as the answer.
                message["content"] = null;
                message["reasoning_content"] = Reasoning;
            }
            else
            {
                message["content"] = Answer(asked);
            }
            var payload = new JsonObject
            {
                ["id"] = "probe-1", ["object"] = "chat.completion", ["model"] = "unlisted-probe-model",
                ["choices"] = new JsonArray(new JsonObject
                    {
                        ["index"] = 0, ["message"] = message, ["finish_reason"] = "stop"
                    }),
                ["usage"] = new JsonObject
                {
                    ["prompt_tokens"] = 24, ["completion_tokens"] = 1,
                    ["prompt_tokens_details"] = new JsonObject { ["cached_tokens"] = 0 }
                }
            };
            return new HttpResponseMessage(HttpStatusCode.OK)
            {
                Content = new StringContent(payload.ToJsonString(), Encoding.UTF8, "application/json")
            };
        }
    }

    // Each proof asks two distinct colors, in parallel.
    private const int ProbeRequests = 2;

    private static BlobPart Image()
        => new("image/png", new BlobRef(RefKind.Uri, "data:image/png;base64," + RedPng));

    private static async Task<(Agent Agent, Provider Provider)> Connected(Application app, string model)
    {
        var config = new Configuration(Vendor.OpenAI, model, "http://127.0.0.1:9", apiKey: "key");
        var agent = new Agent(app, config);
        var provider = new Provider { Name = model };
        agent.Engine().Destroy();
        _ = new Engine(agent, provider);
        return (agent, provider);
    }

    public static async Task UnlistedEndpointEarnsImages()
    {
        var app = new Application(TextWriter.Null);
        try
        {
            var (agent, provider) = await Connected(app, "qwen-probe-capable");
            try
            {
                Check(Refused(() => agent.Engine().ValidateMedia([Image()], Actor.Human)),
                    "unlisted model admitted images before any proof");
                // Case-insensitive on purpose: the probe must not depend on how
                // the endpoint capitalizes a one-word answer.
                provider.Answer = asked => asked.Length == 0 ? "I cannot see any image."
                    : char.ToUpperInvariant(asked[0]) + asked[1..];
                Check(await agent.Engine().ProbeImagesAsync(default), "capable endpoint was refused");
                Check(provider.Bodies.Count == ProbeRequests, "capable endpoint took an unexpected probe count");
                agent.Engine().ValidateMedia([Image()], Actor.Human);
                Check(await agent.Engine().ProbeImagesAsync(default), "proven route lost its capability");
                Check(provider.Bodies.Count == ProbeRequests, "proven route spent another request");
                // The two probes must ask different colors, or a memorized
                // single answer would satisfy the proof.
                Check(provider.Bodies[0] != provider.Bodies[1], "both probes asked the same color");
            }
            finally { agent.Destroy(); }
        }
        finally { app.Destroy(); }
    }

    public static async Task UnprovenEndpointsStayRefused()
    {
        var app = new Application(TextWriter.Null);
        try
        {
            // A server that never received the bytes cannot name the color.
            var (blind, blindProvider) = await Connected(app, "qwen-probe-blind");
            try
            {
                Check(!await blind.Engine().ProbeImagesAsync(default), "text-only endpoint was credited with vision");
                Check(Refused(() => blind.Engine().ValidateMedia([Image()], Actor.Human)),
                    "refused probe still admitted image bytes");
                Check(!await blind.Engine().ProbeImagesAsync(default), "negative result expired immediately");
                Check(blindProvider.Bodies.Count == ProbeRequests, "refused endpoint took an unexpected count");
                Check(blindProvider.Bodies.All(body => body.Contains("ensemble-media-probe",
                    StringComparison.Ordinal)), "a probe omitted the trace marker");
            }
            finally { blind.Destroy(); }
            // The color named only inside reasoning is not a delivered answer.
            var (hidden, hiddenProvider) = await Connected(app, "qwen-probe-reasoning");
            try
            {
                hiddenProvider.Reasoning = "red";
                Check(!await hidden.Engine().ProbeImagesAsync(default), "reasoning was accepted as an answer");
                Check(Refused(() => hidden.Engine().ValidateMedia([Image()], Actor.Human)),
                    "reasoning-only probe admitted image bytes");
            }
            finally { hidden.Destroy(); }
            // A reply that hedges over every color is not an observation.
            var (hedged, hedgedProvider) = await Connected(app, "qwen-probe-hedged");
            try
            {
                hedgedProvider.Answer = _ => "red, green or blue";
                Check(!await hedged.Engine().ProbeImagesAsync(default), "an all-color hedge was credited with vision");
                Check(Refused(() => hedged.Engine().ValidateMedia([Image()], Actor.Human)),
                    "hedged probe admitted image bytes");
            }
            finally { hedged.Destroy(); }
            // A fixed answer cannot survive two different questions. This is
            // the case a single-probe design would pass one third of the time.
            var (stubborn, fixedProvider) = await Connected(app, "qwen-probe-fixed");
            try
            {
                fixedProvider.Answer = _ => "red";
                Check(!await stubborn.Engine().ProbeImagesAsync(default),
                    "a fixed color answer earned image capability");
                Check(fixedProvider.Bodies.Count == ProbeRequests, "fixed answer took an unexpected count");
            }
            finally { stubborn.Destroy(); }
            // A refused route is not punished again for a transport failure.
            var (down, downProvider) = await Connected(app, "qwen-probe-down");
            try
            {
                downProvider.Status = HttpStatusCode.ServiceUnavailable;
                Check(!await down.Engine().ProbeImagesAsync(default), "provider failure was read as capability");
                Check(downProvider.Bodies.Count == ProbeRequests, "failed route took an unexpected count");
            }
            finally { down.Destroy(); }
        }
        finally { app.Destroy(); }
    }

    public static async Task CuratedTableStaysAuthoritative()
    {
        var app = new Application(TextWriter.Null);
        try
        {
            // A table row is never re-probed, in either direction.
            var (listed, listedProvider) = await Connected(app, "gpt-4o");
            try
            {
                Check(await listed.Engine().ProbeImagesAsync(default), "table row lost its declared images");
                Check(listedProvider.Bodies.Count == 0, "curated model was probed");
                listed.Engine().ValidateMedia([Image()], Actor.Human);
            }
            finally { listed.Destroy(); }
            var (textOnly, textProvider) = await Connected(app, "gpt-4o-mini");
            try
            {
                Check(await textOnly.Engine().ProbeImagesAsync(default), "second table row lost its images");
                Check(textProvider.Bodies.Count == 0, "curated model was probed");
            }
            finally { textOnly.Destroy(); }
        }
        finally { app.Destroy(); }
    }
}
