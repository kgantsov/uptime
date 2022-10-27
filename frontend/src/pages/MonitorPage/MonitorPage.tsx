import { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { AreaChart } from "@tremor/react";
import { format } from 'date-fns';
import { Service } from '../../types/services';
import { Heartbeat } from '../../types/heartbeats';

import styles from './MonitorPage.module.css';

export function MonitorPage() {
    const { monitorId } = useParams();
    const [service, setService] = useState<Service | null>(null);
    const [latencies, setLatencies] = useState<Heartbeat[]>([]);

    async function fetchData(monitorId: number) {
      try {
        const response = await fetch(`/API/v1/services/${monitorId}`);
        const data = await response.json();
        setService(data)
      } catch(e) {
          console.log(e);
      }
    }

    async function fetchLatencies(monitorId: number) {
      try {
        const response = await fetch(`/API/v1/heartbeats/latencies?service_id=${monitorId}&size=10`);
        const data = await response.json();
        setLatencies(data.reverse())
      } catch (e) {
          console.log(e);
      }
    }

    useEffect(() => {
      fetchData(Number(monitorId));
      fetchLatencies(Number(monitorId));

      const interval = setInterval(() => {
        fetchLatencies(Number(monitorId));
      }, 5000);
      return () => clearInterval(interval);
    }, [monitorId]);

    if (!service) {
      return <></>
    }

    return (
      <>
        <div className={styles.header}>
          <h1 className={styles.title}>{service.name}</h1>
          <h4><a href={service.url} target='_blank' rel="noreferrer">{service.url}</a></h4>
        </div>

        <AreaChart
          data={latencies.map(item => {
            const createdAt = new Date(item.created_at)
            return {
              time: format(createdAt, 'H:mm'),
              response_time: item.response_time,
            }
          })}
          categories={['response_time',]}
          dataKey="time"
          colors={["cyan"]}
          valueFormatter={undefined}
          startEndOnly={false}
          showXAxis={true}
          showYAxis={true}
          yAxisWidth="w-14"
          showTooltip={true}
          showLegend={false}
          showGridLines={false}
          showAnimation={true}
          height="h-80"
          marginTop="mt-0"
        />
      </>
    )
  }