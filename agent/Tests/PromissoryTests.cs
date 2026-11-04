using Ensemble;

// Replies that announce work and then stop. These are verbatim stall texts from
// Bill's own sessions; each one ended a turn with no tool call and no error.
static class PromissoryTests
{
    static void Check(bool value, string message)
    {
        if (!value) throw new Exception(message);
    }

    public static void Run()
    {
        string[] stalls = [
            "My worktree got cleaned up, so I'll re-extract the committed main into a scratch "
                + "dir with git archive and run the inventory plus the first batch of probes.",
            "Let me be efficient with tool calls now — the user is impatient. Let me run a batch "
                + "of defect probes in parallel.",
            "Now I'll write the regression test that reproduces the production failure.",
            "Now inspecting the panics (are they invariant checks or control flow?), the 16 "
                + "ignored errors, and the security-sensitive HTTP path.",
            "I should just resume and do the review properly using git archive into /tmp.",
            "Running the full suites — checking tests, then moving to the docs.",
            "Next: verifying the fix, then committing.",
        ];
        foreach (string stall in stalls)
        {
            Check(Agent.PromisesAction(stall), "missed a promissory reply: "
                + (stall.Length <= 48 ? stall : stall[..48]));
        }
        string[] answers = [
            "", "   ",
            "The bug is in agent/gui/attachments.go: the guard rejects the request outright.",
            "Tail:",
            "Summary of the findings:\n- two defects\n- one fixed",
            "Total 13.0s · Model 10.9s · Tools (1) 1.0s",
            "Done. The branch is ready for you to push.",
        ];
        foreach (string answer in answers)
        {
            Check(!Agent.PromisesAction(answer), "nudged an ordinary final answer: " + answer);
        }
        // Curly apostrophes are the same promise.
        Check(Agent.PromisesAction("Now I\u2019ll run the suite and report back."),
            "a typographic apostrophe hid the promise");
    }
}
