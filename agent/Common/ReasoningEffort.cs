namespace Ensemble;

// A requested reasoning effort that the endpoint will not honour is a launch
// mistake, not a runtime condition: the operator believes they are getting
// extended thinking while the server quietly ignores or rejects the field. The
// table row only states what this application knows, so the endpoint itself gets
// the last word before the first turn.
public sealed partial class Engine
{
    public static readonly TimeSpan EffortProbeTimeout = TimeSpan.FromSeconds(10);

    // Confirm the configured effort is one this endpoint will accept. An empty
    // request sends no traffic at all, which keeps ordinary launches and offline
    // fixtures free of extra probes. Anything short of an explicit refusal of
    // this parameter launches with a note, because ensemble already starts
    // against an offline provider.
    public partial Task ValidateReasoningEffortAsync(CancellationToken cancellation);
}
