using System.Runtime.CompilerServices;

// Isolated lease lifecycle fixtures; public GUI regression consumers stay public.
[assembly: InternalsVisibleTo("VoiceLeaseTests")]

// Digest parsing and verification are the security-critical parts, so the GUI
// regression suite exercises them directly as well as through a live listener.
[assembly: InternalsVisibleTo("GuiTests")]
