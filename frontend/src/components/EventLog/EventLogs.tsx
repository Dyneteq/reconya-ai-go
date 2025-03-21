import React from 'react';
import useEventLogs from '../../hooks/useEventLogs'; // Adjust path as needed
import { EventLog } from '../../models/eventLog.model';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { EventLogIcons } from '../../models/eventLogIcons.model';

const EventLogs = () => {
  const eventLogs: EventLog[] = useEventLogs();

  const formatDate = (date: string | Date) => {
    const parsedDate = typeof date === 'string' ? new Date(date) : date;
    return `${parsedDate.toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
    })} ${parsedDate.toLocaleTimeString('en-US', {
      hour: '2-digit',
      minute: '2-digit',
      hour12: false,
    })}`;
  };

  return (
    <div className="mt-5">
      <h6 className="text-success d-block w-100">[ EVENT LOG ]</h6>
      <table className="table table-dark table-sm table-compact border-dark border-bottom text-success" style={{ fontSize: '13px' }}>
        <tbody>
          {eventLogs.length > 0 ? (
            eventLogs.map((log: EventLog, index: React.Key | null | undefined) => {
              const icon = EventLogIcons[log.Type];

              return (
                <tr key={index}>
                  <td className="bg-transparent text-success px-3">
                    {icon ? (
                      <FontAwesomeIcon icon={icon} className="text-success" />
                    ) : (
                      <span>??</span>
                    )}
                  </td>
                  <td className="bg-transparent text-success px-3">{log.Description}</td>
                  <td className="bg-transparent text-success px-3 text-end">{formatDate(log.CreatedAt)}</td>
                </tr>
              );
            })
          ) : (
            <tr>
              <td className="bg-transparent text-success px-3" colSpan={3}>No event logs available</td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  );
};

export default EventLogs;
