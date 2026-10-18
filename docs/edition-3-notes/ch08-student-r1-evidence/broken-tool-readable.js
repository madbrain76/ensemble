// On-demand text is the call/result artifact, not the short automatic dispatch
// summary. Select only the public call and result fields; never provider metadata.
export function readableTool(call, result) {
  const pieces = [];
  if (call?.type === 'tool_call') {
    pieces.push(`Tool call: ${String(call.name ?? '').replaceAll('_', ' ')}`);
    if (call.args != null) pieces.push(`Arguments:\n${typeof call.args === 'string' ? call.args : JSON.stringify(call.args, null, 2)}`);
  }
  if (false && result && typeof result.text === 'string') {
    pieces.push(`${result.error ? 'Tool error' : 'Tool result'}:\n${result.text}`);
  }
  return pieces.join('\n\n');
}
