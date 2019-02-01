package network

import(
	"os"
	"os/exec"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
	"net"
	log "github.com/Sirupsen/logrus"
	"encoding/json"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"path"
	"fmt"
	"mydocker/container"
	"mydocker/util"
	"runtime"
)

var (
	//defaultNetworkPath = "/var/run/mydocker/network/network/"
	defaultNetworkPath = util.GetNetworkPath()
	defaultEndpointPath = util.GetEndPointPath()
	drivers = map[string]NetworkDriver{}
	networks = map[string]*Network{}
	endpoints = map[string]*Endpoint{}
)

type Endpoint struct {
	ID string `json:"id"`
	Device netlink.Veth `json:"dev"`
	IPAddress net.IP `json:"ip"`
	MacAddress net.HardwareAddr `json:"mac"`
	Network    *Network
	PortMapping []string
}

type Network struct {
	Name string
	IpRange *net.IPNet
	Driver string
}

type NetworkDriver interface {
	Name() string
	Create(subnet string, name string) (*Network, error)
	Delete(network Network) error
	Connect(network *Network, endpoint *Endpoint) error
	Disconnect(network Network, endpoint *Endpoint) error
}

func CreateNetwork(driver, subnet, name string) error {
	_, cidr, _ := net.ParseCIDR(subnet)
	gatewayIp, err := ipAllocator.Allocate(cidr)
	if err != nil {
		return nil
	}
	cidr.IP = gatewayIp

	nw, err := drivers[driver].Create(cidr.String(), name)
	//return nw.dump(defaultNetworkPath)
	return dump(nw, defaultNetworkPath, nw.Name)
}

func dump(obj interface{}, dumpPath string, containerName string) error {
        if _, err := os.Stat(dumpPath); os.IsNotExist(err) {
                if err = os.MkdirAll(dumpPath, 0644); err != nil {
                        return err
                }
        }
        nwPath := path.Join(dumpPath, containerName)
        nwFile, err := os.OpenFile(nwPath, os.O_TRUNC | os.O_WRONLY | os.O_CREATE, 0644)
        if err != nil {
                log.Errorf("dump os.OpenFile error: %v", err)
                return err
        }
        defer nwFile.Close()

        nsJson, err := json.Marshal(obj)
        if err != nil {
                log.Errorf("dump json.Marshal error: %v", err)
                return err
        }
        _, err = nwFile.Write(nsJson)
        if err != nil {
                log.Errorf("Network dump nwFile.Write error: %v", err)
                return err
        }
        return nil
}

func load(dumpPath string, obj interface{}) (interface{}, error) {
        nwConfigFile, err := os.Open(dumpPath)
        if err != nil {
                return nil, err
        }
        defer nwConfigFile.Close()

        nwJson := make([]byte, 2000)
        n, err := nwConfigFile.Read(nwJson)
        if err != nil {
                return nil, err
        }

        err = json.Unmarshal(nwJson[:n], obj)
        //err = json.Unmarshal(nwJson[:n], ep)
        if err != nil {
                return nil,err
        }
        return obj, nil
}

func Init() error {
	var bridgeDriver = BridgeNetworkDriver{}
	drivers[bridgeDriver.Name()] = &bridgeDriver
	if _, err := os.Stat(defaultNetworkPath); os.IsNotExist(err) {
                if err = os.MkdirAll(defaultNetworkPath, 0644); err != nil {
                        return err
                }
        }

	//search defaultNetworkPath to find network definition file
	filepath.Walk(defaultNetworkPath, func(nwPath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		_, nwName := path.Split(nwPath)
		obj, err := load(nwPath, &Network{})
                if err != nil {
                        return fmt.Errorf("error load network: %s", err)
                }

                nw, _ := obj.(*Network)
		networks[nwName] = nw

		return nil
	})

	if _, err := os.Stat(defaultEndpointPath); os.IsNotExist(err) {
                if err = os.MkdirAll(defaultEndpointPath, 0644); err != nil {
                        return err
                }
        }
	//search defaultEndpointPath to find container endpoint definition file
	filepath.Walk(defaultEndpointPath, func(nwPath string, info os.FileInfo, err error) error {
		if info.IsDir() {
                        return nil
                }
		_, containerName := path.Split(nwPath)
		obj, err := load(nwPath, &Endpoint{})
		//ep, err := load(nwPath)
		if err != nil {
			return fmt.Errorf("error load endpoint: %s", err)
		}

		ep, _ := obj.(*Endpoint)
		endpoints[containerName] = ep
		return nil
	})

	return nil
}

func ListNetwork() {
	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	fmt.Fprint(w, "NAME\tIpRange\tDriver\n")
	for _, nw := range networks {
		fmt.Fprintf(w, "%s\t%s\t%s\n",
			nw.Name,
			nw.IpRange.String(),
			nw.Driver,
		)
	}
	if err := w.Flush(); err != nil {
		log.Errorf("Flush error %v", err)
		return
	}
}

func DeleteNetwork(networkName string) error {
	nw, ok := networks[networkName]
	if !ok {
		return fmt.Errorf("no such network: %s", networkName)
	}
	if err := ipAllocator.Release(nw.IpRange, &nw.IpRange.IP); err != nil {
		return fmt.Errorf("Error Remove Network gateway ip: %s", err)
	}
	if err := drivers[nw.Driver].Delete(*nw); err != nil {
		return fmt.Errorf("Error Remove Network driver: %v", err)
	}
	return nw.remove(defaultNetworkPath)
}

func (nw *Network) remove(dumpPath string) error {
	if _, err := os.Stat(path.Join(dumpPath, nw.Name)); err != nil {
		if os.IsNotExist(err) {
                        return nil
                } else {
			return err
		}
        } else {
		return os.Remove(path.Join(dumpPath, nw.Name))
	}
}

func Disconnect(networkName string, info *container.ContainerInfo) error {
	nw, ok := networks[networkName]
	if !ok {
		return fmt.Errorf("No Such Network: %s", networkName)
	}

	ep, ok := endpoints[info.Name]
	if !ok {
                return fmt.Errorf("No Such container %s endpoint", info.Name)
        }
	//release the container ip
	if err := ipAllocator.Release(nw.IpRange, &ep.IPAddress); err != nil {
		return fmt.Errorf("Error Remove container %s ip %s: %v", info.Name, ep.IPAddress, err)
	}
	//delete container endpoint info file
	deleteEndpointInfo(info.Name)
	return nil
}

func deleteEndpointInfo(containerName string) error {
	epfile := defaultEndpointPath + containerName
	if err := os.RemoveAll(epfile); err != nil {
                return fmt.Errorf("deleteEndpointInfo Remove dir %s error %v", epfile, err)
        }
	return nil
}

func Connect(networkName string, info *container.ContainerInfo) error {
	network, ok := networks[networkName]
	if !ok {
		return fmt.Errorf("network %s not exist", networkName)
	}
	ip, err := ipAllocator.Allocate(network.IpRange)
        if err != nil {
                return fmt.Errorf("error in ipAllocator.Allocate of Connect: %v", err)
        }
	ep := &Endpoint{
		ID: fmt.Sprintf("%s-%s", info.Id, networkName),
		IPAddress: ip,
		Network: network,
		PortMapping: info.PortMapping,
	}
	if err := drivers[network.Driver].Connect(network, ep); err != nil {
		return fmt.Errorf("error in driver Connect : %v", err)
	}

	if err = configEndpointIpAddressAndRoute(ep, info); err != nil {
		return fmt.Errorf("error in configEndpointIpAddressAndRoute of Connect: %v", err)
	}

	//write endpoint definition to file, file name is container name
	dump(ep, defaultEndpointPath, info.Name)
	return configPortMapping(ep, info)
}

func configEndpointIpAddressAndRoute(ep *Endpoint, cinfo *container.ContainerInfo) error {
	peerLink, err := netlink.LinkByName(ep.Device.PeerName)
	if err != nil {
		return fmt.Errorf("fail config endpoint: %v", err)
	}

	defer enterContainerNetns(&peerLink, cinfo)()

	interfaceIP := *ep.Network.IpRange
	interfaceIP.IP = ep.IPAddress

	if err = setInterfaceIP(ep.Device.PeerName, interfaceIP.String()); err != nil {
		return fmt.Errorf("%v,%s", ep.Network, err)
	}

	if err = setInterfaceUP(ep.Device.PeerName); err != nil {
		return err
	}

	if err = setInterfaceUP("lo"); err != nil {
		return err
	}

	_, cidr, _ := net.ParseCIDR("0.0.0.0/0")

	defaultRoute := &netlink.Route{
		LinkIndex: peerLink.Attrs().Index,
		Gw: ep.Network.IpRange.IP,
		Dst: cidr,
	}

	if err = netlink.RouteAdd(defaultRoute); err != nil {
		return err
	}

	return nil
}

func enterContainerNetns(enLink *netlink.Link, cinfo *container.ContainerInfo) func() {
	f, err := os.OpenFile(fmt.Sprintf("/proc/%s/ns/net", cinfo.Pid), os.O_RDONLY, 0)
	if err != nil {
		log.Errorf("error get container net namespace, %v", err)
	}

	nsFD := f.Fd()
	runtime.LockOSThread()

	if err = netlink.LinkSetNsFd(*enLink, int(nsFD)); err != nil {
		log.Errorf("error set link netns , %v", err)
	}

	origns, err := netns.Get()
	if err != nil {
		log.Errorf("error get current netns, %v", err)
	}

	if err = netns.Set(netns.NsHandle(nsFD)); err != nil {
		log.Errorf("error set netns, %v", err)
	}
	return func () {
		netns.Set(origns)
		origns.Close()
		runtime.UnlockOSThread()
		f.Close()
	}
}

func configPortMapping(ep *Endpoint, cinfo *container.ContainerInfo) error {
	for _, pm := range ep.PortMapping {
		portMapping :=strings.Split(pm, ":")
		if len(portMapping) != 2 {
			log.Errorf("port mapping format error, %v", pm)
			continue
		}
		iptablesCmd := fmt.Sprintf("-t nat -A PREROUTING -p tcp -m tcp --dport %s -j DNAT --to-destination %s:%s",
			portMapping[0], ep.IPAddress.String(), portMapping[1])
		cmd := exec.Command("iptables", strings.Split(iptablesCmd, " ")...)
		//err := cmd.Run()
		output, err := cmd.Output()
		if err != nil {
			log.Errorf("iptables Output, %v", output)
			continue
		}
	}
	return nil
}

