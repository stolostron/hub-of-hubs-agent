package consumer

import bundle "github.com/stolostron/hub-of-hubs-agent/pkg/spec/bundle"

// Transport is an interface for transport layer.
type Consumer interface {
	// Start starts the transport.
	Start()
	// Stop stops the transport.
	Stop()
	// Register registers a bundle ID to a CustomBundleRegistration. None-registered bundles are assumed to be
	// of type GenericBundle, and are handled by the GenericBundleSyncer.
	Register(msgID string, customBundleRegistration *bundle.CustomBundleRegistration)

	GetGenericBundleChan() chan *bundle.GenericBundle
}