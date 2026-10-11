package common

// CacheLens compares main-conversation requests, never judge/compressor prompts.
// Engine supplies sanitized wire bytes and normalized usage at the send boundary.
type CacheLens interface {
	// Engine exposes the required owner, including its Agent's data directory.
	Engine() Engine
	// ObserveRequest records one actual request before the HTTP exchange.
	ObserveRequest(name string, body []byte)
	// ObserveUsage compares provider-granted tokens with the byte prediction.
	ObserveUsage(Usage)
}
