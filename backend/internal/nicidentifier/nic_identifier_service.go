package nicidentifier

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"reconya-ai/internal/device"
	"reconya-ai/internal/eventlog"
	"reconya-ai/internal/network"
	"reconya-ai/internal/systemstatus"
	"reconya-ai/models"
)

// NicIdentifierService struct
type NicIdentifierService struct {
	NetworkService      *network.NetworkService
	SystemStatusService *systemstatus.SystemStatusService
	EventLogService     *eventlog.EventLogService
	DeviceService       *device.DeviceService
}

// NewNicIdentifierService creates a new instance of NicIdentifierService
func NewNicIdentifierService(
	networkService *network.NetworkService,
	systemStatusService *systemstatus.SystemStatusService,
	eventLogService *eventlog.EventLogService,
	deviceService *device.DeviceService) *NicIdentifierService {
	return &NicIdentifierService{
		NetworkService:      networkService,
		SystemStatusService: systemStatusService,
		EventLogService:     eventLogService,
		DeviceService:       deviceService,
	}
}

// Identify performs the NIC identification process
func (s *NicIdentifierService) Identify() {
	log.Printf("Attempting to identify network")
	nic := s.getLocalNic()
	fmt.Printf("NIC: %v\n", nic)
	cidr := extractCIDR(nic.IPv4)
	publicIP, err := s.getPublicIp()
	if err != nil {
		log.Printf("Failed to get public IP: %v", err)
		return
	}

	networkEntity, err := s.NetworkService.FindOrCreate(cidr)
	if err != nil {
		log.Printf("Failed to find or create network: %v", err)
		return
	}

	localDevice := models.Device{
		Name:   nic.Name,
		IPv4:   nic.IPv4,
		Status: models.DeviceStatusOnline,
	}

	// Save or update the device in the database.
	savedDevice, err := s.DeviceService.CreateOrUpdate(&localDevice)
	if err != nil {
		log.Printf("Failed to save or update local device: %v", err)
		return
	}

	// Use the potentially updated or newly created system status,
	// including the saved device.
	systemStatus := models.SystemStatus{
		LocalDevice: *savedDevice, // Use savedDevice, which is now updated with ID
		Network:     networkEntity,
		PublicIP:    &publicIP,
	}

	_, err = s.SystemStatusService.CreateOrUpdate(&systemStatus)
	if err != nil {
		log.Printf("Failed to create or update system status: %v", err)
		return
	}

	s.EventLogService.CreateOne(&models.EventLog{
		Type:     models.LocalIPFound,
		DeviceID: &savedDevice.ID,
	})

	s.EventLogService.CreateOne(&models.EventLog{
		Type: models.LocalNetworkFound,
	})
}

func (s *NicIdentifierService) getLocalNic() models.NIC {
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Println("Error getting network interfaces:", err)
		return models.NIC{}
	}

	for _, iface := range interfaces {
		fmt.Printf("Checking interface: %s\n", iface.Name)
		if iface.Flags&net.FlagUp == 0 {
			fmt.Printf("Skipping %s: interface is down\n", iface.Name)
			continue
		}
		if iface.Flags&net.FlagLoopback != 0 {
			fmt.Printf("Skipping %s: interface is loopback\n", iface.Name)
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			fmt.Printf("Skipping %s: error getting addresses: %v\n", iface.Name, err)
			continue
		}

		for _, addr := range addrs {
			ip, _, err := net.ParseCIDR(addr.String())
			if err != nil || ip.To4() == nil {
				fmt.Printf("Skipping address %s on %s: not a valid IPv4\n", addr.String(), iface.Name)
				continue
			}

			if !ip.IsLoopback() {
				fmt.Printf("Found matching interface: %s with IPv4: %s\n", iface.Name, ip.String())
				return models.NIC{Name: iface.Name, IPv4: ip.String()}
			}
		}
	}

	return models.NIC{}
}

func (s *NicIdentifierService) getPublicIp() (string, error) {
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(ip), nil
}
