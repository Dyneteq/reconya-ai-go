package eventlog

import (
	"context"
	"fmt"
	"log"
	"reconya-ai/db"
	"reconya-ai/internal/device"
	"reconya-ai/models"
	"time"
)

type EventLogService struct {
	repository    db.EventLogRepository
	DeviceService *device.DeviceService
}

func NewEventLogService(repository db.EventLogRepository, deviceService *device.DeviceService) *EventLogService {
	return &EventLogService{
		repository:    repository,
		DeviceService: deviceService,
	}
}

func (s *EventLogService) GetAll(limitSize int64) ([]models.EventLog, error) {
	ctx := context.Background()
	eventLogPtrs, err := s.repository.FindLatest(ctx, int(limitSize))
	if err != nil {
		return nil, err
	}

	eventLogs := make([]models.EventLog, len(eventLogPtrs))
	for i, logPtr := range eventLogPtrs {
		eventLogs[i] = *logPtr
		eventLogs[i].Description = s.generateDescription(eventLogs[i])
	}
	
	return eventLogs, nil
}

func (s *EventLogService) generateDescription(eventLog models.EventLog) string {
	deviceInfo := "unknown device"
	if eventLog.DeviceID != nil {
		device, err := s.DeviceService.FindByID(*eventLog.DeviceID)
		if err != nil {
			log.Printf("Error fetching device information: %v", err)
		} else if device != nil && device.IPv4 != "" {
			deviceInfo = device.IPv4
		}
	}

	switch eventLog.Type {
	case models.PingSweep:
		return "Ping sweep performed"
	case models.PortScanStarted:
		return fmt.Sprintf("Port scan started for [%s]", deviceInfo)
	case models.PortScanCompleted:
		return fmt.Sprintf("Port scan completed [%s]", deviceInfo)
	case models.DeviceOnline:
		return fmt.Sprintf("Live device [%s] found", deviceInfo)
	case models.DeviceIdle:
		return fmt.Sprintf("Device [%s] became idle", deviceInfo)
	case models.DeviceOffline:
		return fmt.Sprintf("Device [%s] is now offline", deviceInfo)
	case models.LocalIPFound:
		return fmt.Sprintf("Local IPv4 address found [%s]", deviceInfo)
	case models.LocalNetworkFound:
		return "Local network found"
	case models.Warning:
		return "Warning event occurred"
	case models.Alert:
		return "Alert event occurred"
	default:
		return "Event occurred"
	}
}

func (s *EventLogService) GetAllByDeviceId(deviceId string, limitSize int64) ([]models.EventLog, error) {
	ctx := context.Background()
	eventLogPtrs, err := s.repository.FindAllByDeviceID(ctx, deviceId)
	if err != nil {
		return nil, err
	}

	// Convert pointers to values and respect the limit
	count := int(limitSize)
	if count > len(eventLogPtrs) {
		count = len(eventLogPtrs)
	}
	
	eventLogs := make([]models.EventLog, count)
	for i := 0; i < count; i++ {
		eventLogs[i] = *eventLogPtrs[i]
	}
	
	return eventLogs, nil
}

func (s *EventLogService) CreateOne(eventLog *models.EventLog) error {
	now := time.Now()
	eventLog.CreatedAt = &now
	eventLog.UpdatedAt = &now

	return s.repository.Create(context.Background(), eventLog)
}
