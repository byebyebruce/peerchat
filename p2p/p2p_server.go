package p2p

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p-core/host"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"
)

//const service = "manishmeganathan/peerchat"

// A structure that represents a P2P Host
type P2PBoot struct {
	// Represents the host context layer
	Ctx context.Context

	// Represents the libp2p host
	Host host.Host

	// Represents the DHT routing table
	KadDHT *dht.IpfsDHT

	// Represents the peer discovery service
	Discovery *discovery.RoutingDiscovery

	service string
}

/*
A constructor function that generates and returns a P2P object.

Constructs a libp2p host with TLS encrypted secure transportation that works over a TCP
transport connection using a Yamux Stream Multiplexer and uses UPnP for the NAT traversal.

A Kademlia DHT is then bootstrapped on this host using the default peers offered by libp2p
and a Peer Discovery service is created from this Kademlia DHT. The PubSub handler is then
created on the host using the peer discovery service created prior.
*/
func NewP2PBoot(port int, service string, pkFile string, bootstrapPeers []multiaddr.Multiaddr) *P2PBoot {
	// Setup a background context
	ctx := context.Background()

	// Setup a P2P Host Node
	nodehost, kaddht := setupHost(ctx, port, pkFile)
	// Debug log
	logrus.Debugln("Created the P2P Host and the Kademlia DHT.")

	// Bootstrap the Kad DHT
	//bootstrapDHT(ctx, nodehost, kaddht, bootstrapPeers)
	// Debug log
	logrus.Debugln("Bootstrapped the Kademlia DHT and Connected to Bootstrap Peers")
	// Bootstrap the Kad DHT
	bootstrapDHT(ctx, nodehost, kaddht, bootstrapPeers)

	// Create a peer discovery service using the Kad DHT
	routingdiscovery := discovery.NewRoutingDiscovery(kaddht)
	// Debug log
	logrus.Debugln("Created the Peer Discovery Service.")

	// Create a PubSub handler with the routing discovery
	//pubsubhandler := setupPubSub(ctx, nodehost, routingdiscovery)
	// Debug log
	logrus.Debugln("Created the PubSub Handler.")

	// Return the P2P object
	p2phost := &P2PBoot{
		Ctx:       ctx,
		Host:      nodehost,
		KadDHT:    kaddht,
		Discovery: routingdiscovery,
		service:   service,
	}
	//p2phost.AdvertiseConnect()
	return p2phost
}

func (p2p *P2PBoot) AdvertiseConnect() {
	// Advertise the availabilty of the service on this node
	ttl, err := p2p.Discovery.Advertise(p2p.Ctx, p2p.service)
	// Debug log
	logrus.Debugln("Advertised the PeerChat Service.")
	// Sleep to give time for the advertisment to propogate
	time.Sleep(time.Second * 5)
	// Debug log
	logrus.Debugf("Service Time-to-Live is %s", ttl)

	// Find all peers advertising the same service
	peerchan, err := p2p.Discovery.FindPeers(p2p.Ctx, p2p.service)
	// Handle any potential error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatalln("P2P Peer Discovery Failed!")
	}
	// Trace log
	logrus.Traceln("Discovered PeerChat Service Peers.")

	// Connect to peers as they are discovered
	go handlePeerDiscovery(p2p.Host, peerchan)
	// Trace log
	logrus.Traceln("Started Peer Connection Handler.")
}
