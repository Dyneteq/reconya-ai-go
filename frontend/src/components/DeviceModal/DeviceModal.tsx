import React from "react";
import { Device } from "../../models/device.model";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCircle, faExternalLinkAlt } from "@fortawesome/free-solid-svg-icons";

interface DeviceModalProps {
  device: Device | null;
  onClose: () => void;
}

const DeviceModal: React.FC<DeviceModalProps> = ({ device, onClose }) => {
  if (!device) return null;

  const calcTimeElapsed = (dateString: string | undefined) => {
    if (!dateString) return "N/A";
    return new Date(dateString).toLocaleDateString("en-US", {
      year: "numeric",
      month: "long",
      day: "numeric",
    });
  };

  const getStatusIcon = (status: string | undefined) => {
    if (!status) return null;

    let colorClass = "";
    switch (status.toLowerCase()) {
      case "online":
        colorClass = "text-success";
        break;
      case "offline":
        colorClass = "text-danger";
        break;
      case "idle":
        colorClass = "text-warning";
        break;
      default:
        colorClass = "text-muted";
        break;
    }

    return <FontAwesomeIcon icon={faCircle} className={`${colorClass} me-2`} />;
  };

  const getPortStateIcon = (state: string) => {
    if (state === "open") {
      return (
        <FontAwesomeIcon
          icon={faCircle}
          className="text-danger me-2"
          style={{ fontSize: "0.6rem" }}
        />
      );
    } else if (state === "filtered") {
      return (
        <FontAwesomeIcon
          icon={faCircle}
          className="text-warning me-2"
          style={{ fontSize: "0.6rem" }}
        />
      );
    }
    return null;
  };

  const getPortLink = (portNumber: number, ipAddress: string) => {
    const httpPorts = [80, 8080, 8000];
    const httpsPorts = [443, 8443];

    if (httpPorts.includes(portNumber)) {
      return `http://${ipAddress}:${portNumber}`;
    } else if (httpsPorts.includes(portNumber)) {
      return `https://${ipAddress}:${portNumber}`;
    }
    return null;
  };

  const renderPortIcons = () => {
    if (device?.Ports?.some((port) => port.state === "open")) {
      return (
        <FontAwesomeIcon
          icon={faCircle}
          className="text-danger me-2"
          style={{ fontSize: "1.2rem" }}
        />
      );
    }

    if (device?.Ports?.some((port) => port.state === "filtered")) {
      return (
        <FontAwesomeIcon
          icon={faCircle}
          className="text-warning me-2"
          style={{ fontSize: "1.2rem" }}
        />
      );
    }

    return (
      <FontAwesomeIcon
        icon={faCircle}
        className="text-success me-2"
        style={{ fontSize: "1.2rem" }}
      />
    );
  };

  const handleModalContentClick = (e: React.MouseEvent) => {
    e.stopPropagation();
  };

  return (
    <div
      className="modal show d-block"
      tabIndex={-1}
      style={{ backgroundColor: "rgba(0,0,0,0.5)", marginTop: "170px" }}
      onClick={onClose}
    >
      <div
        className="modal-dialog modal-lg"
        onClick={handleModalContentClick}
        style={{ minHeight: "400px" }}
      >
        <div className="modal-content">
          <div
            className="modal-body bg-black border border-success border-radius-0 text-success p-5"
            style={{ minHeight: "400px" }}
          >
            <div className="mb-3">
              <div className="border-bottom border-success pb-2 mb-3 d-flex justify-content-between align-items-center">
                <div className="d-flex align-items-center">
                  <span className="orbitron fw-bold fs-2">{device.IPv4}</span>
                </div>
                <div className="d-flex align-items-center">
                  {renderPortIcons()}
                </div>
              </div>

              <table className="text-success w-100 p-2 mb-4">
                <tbody className="p-2">
                  <tr>
                    <td className="w-25 ps-2 fw-bold">Hostname</td>
                    <td>{device.Hostname || "Unknown"}</td>
                  </tr>
                  <tr>
                    <td className="w-25 ps-2 fw-bold">H/W vendor</td>
                    <td>{device.Vendor}</td>
                  </tr>
                  <tr>
                    <td className="w-25 ps-2 fw-bold">MAC Address</td>
                    <td>{device.MAC}</td>
                  </tr>
                  <tr>
                    <td className="w-25 ps-2 fw-bold">Status</td>
                    <td>
                      {getStatusIcon(device.Status)}
                      {device.Status}
                    </td>
                  </tr>
                  <tr>
                    <td className="w-25 ps-2 fw-bold">First appeared</td>
                    <td>{calcTimeElapsed(device.CreatedAt)}</td>
                  </tr>
                  <tr>
                    <td className="w-25 ps-2 fw-bold">Last seen online</td>
                    <td>{calcTimeElapsed(device.LastSeenOnlineAt)}</td>
                  </tr>
                </tbody>
              </table>

              <h6>[ PORTS ]</h6>
              <table className="text-success w-100 p-2 mb-4">
                <tbody className="p-2">
                  {device.Ports?.map((port, index) => {
                    const portNumber = Number(port.number); // Ensure port number is treated as a number
                    const portLink = getPortLink(portNumber, device.IPv4 || "");
                    return (
                      <tr key={index}>
                        <td className="ps-2 fw-bold" style={{ width: "15%" }}>
                          {portNumber}
                        </td>
                        <td style={{ width: "15%" }}>
                          <span className="badge bg-black border border-dark text-success">
                            {getPortStateIcon(port.state)}
                            {port.state}
                          </span>
                        </td>
                        <td style={{ width: "20%" }}>
                          {port.protocol.toUpperCase()}
                        </td>
                        <td>{port.service || "Unknown"}</td>
                        <td>
                          {portLink && (
                            <a
                              href={portLink}
                              target="_blank"
                              rel="noopener noreferrer"
                              className="ms-2 text-light"
                            >
                              <FontAwesomeIcon
                                icon={faExternalLinkAlt}
                                className="text-success"
                              />
                            </a>
                          )}
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>

              {/* Event Log section can be added here based on your data and requirements */}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default DeviceModal;
