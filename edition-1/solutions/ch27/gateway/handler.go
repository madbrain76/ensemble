package gateway

// Handler is implemented by the application that uses the gateway client.
//
// This interface is what makes this package a library rather than a
// framework. The client owns the connection, the handshake and the
// delivery bookkeeping; it owns no opinion about what a message means.
// Ensemble implements Handler to route envelopes into its actor system, a
// standalone interviewer implements it to collect answers, and a home
// automation program implements it to switch on a light.
//
// The client calls these methods from its receive goroutine, one at a time,
// in the order traffic arrives. A Handler therefore does not need its own
// locking for ordering, but it must not block for long: the connection
// makes no further progress until the call returns. Work that takes real
// time belongs on a goroutine the Handler owns.
type Handler interface {
	// OnDeliver is called for each message the gateway delivers to this
	// agent. Returning from this call is what acknowledges the message, so
	// the client only acks work the Handler has actually accepted.
	OnDeliver(env Envelope)

	// OnConnected is called after a successful Noise-KK handshake, with
	// the gateway's welcome.
	OnConnected(welcome WelcomePayload)

	// OnDisconnected is called when the connection drops, with the reason.
	OnDisconnected(reason string)

	// OnError is called on protocol errors that do not end the connection.
	OnError(err error)
}

// NopHandler implements Handler and does nothing.
//
// Embed it to implement only the methods you care about, so that adding a
// method to Handler later does not break every implementation.
type NopHandler struct{}

func (NopHandler) OnDeliver(Envelope)         {}
func (NopHandler) OnConnected(WelcomePayload) {}
func (NopHandler) OnDisconnected(string)      {}
func (NopHandler) OnError(error)              {}
