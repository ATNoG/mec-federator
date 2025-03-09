package models

type InterfaceDetails struct {
	// Each Port and corresponding traffic protocol exposed by the component is identified by a name. Application client on user device requires this to uniquely identify the interface.
	InterfaceId string `json:"interfaceId"`
	// Defines the IP transport communication protocol i.e., TCP, UDP or HTTP
	CommProtocol string `json:"commProtocol"`
	// Port number exposed by the component. OP may generate a dynamic port towards the UCs corresponding to this internal port and forward the client traffic from dynamic port to container Port.
	CommPort int32 `json:"commPort"`
	// Defines whether the interface is exposed to outer world or not i.e., external, or internal. If this is set to \"external\", then it is exposed to external applications otherwise it is exposed internally to edge application components within edge cloud. When exposed to external world, an external dynamic port is assigned for UC traffic and mapped to the internal container Port
	VisibilityType string `json:"visibilityType"`
	// Name of the network.  In case the application has to be associated with more than 1 network then app provider must define the name of the network on which this interface has to be exposed.  This parameter is required only if the port has to be exposed on a specific network other than default.
	Network string `json:"network,omitempty"`
	// Interface Name. Required only if application has to be attached to a network other than default.
	InterfaceName string `json:"InterfaceName,omitempty"`
}
