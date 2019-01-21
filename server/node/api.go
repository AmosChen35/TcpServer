package node

import(
)

// PrivateAdminAPI is the collection of administrative API methods exposed only
// over a secure RPC channel.
type PrivateAdminAPI struct {
    node *Node // Node interfaced by this API
}

// NewPrivateAdminAPI creates a new API definition for the private admin methods
// of the node itself.
func NewPrivateAdminAPI(node *Node) *PrivateAdminAPI {
    return &PrivateAdminAPI{node: node}
}

// AddPeer requests connecting to a remote node, and also maintaining the new
// connection at all times, even reconnecting if it is lost.
func (api *PrivateAdminAPI) HelloAdmin() (string, error) {
    return "HelloAdmin", nil
}
