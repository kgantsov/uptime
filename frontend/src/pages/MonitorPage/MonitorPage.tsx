import { useState, useEffect, SetStateAction } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { AreaChart, Tracking, TrackingBlock } from "@tremor/react";
import { Button } from "@tremor/react";
import { FaTrashAlt, FaPencilAlt } from 'react-icons/fa';
import { BsPlay, BsStop } from 'react-icons/bs';
import { format } from 'date-fns';
import { Service } from '../../types/services';
import { Heartbeat, STATUS_COLORS_MAP } from '../../types/heartbeats';

import styles from './MonitorPage.module.css';
import { API } from '../../API';
import { MonitorCards } from './MonitorCards';

export function MonitorPage() {
    let navigate = useNavigate();
    const { monitorId } = useParams();
    const [service, setService] = useState<Service | null>(null);
    const [latencies, setLatencies] = useState<Heartbeat[]>([]);

    async function fetchData(monitorId: number) {
      try {
        const resp = await API.fetch('GET', `/API/v1/services/${monitorId}`);
        const data = await resp.json();
        setService(data)
      } catch(e) {
          console.log(e);
      }
    }

    const size = 100

    async function fetchLatencies(monitorId: number) {
      try {
        const resp = await API.fetch('GET', `/API/v1/heartbeats/latencies?service_id=${monitorId}&size=${size}`);
        const data = await resp.json();
        setLatencies(data.reverse())
      } catch (e) {
          console.log(e);
      }
    }

    async function handleServiceEnableToggle() {
      try {
        await API.fetch(
          'PATCH', `/API/v1/services/${monitorId}`, null, {enabled: !service?.enabled}
        );
        fetchData(Number(monitorId));
      } catch(e) {
        console.log(e);
      }
    }

    async function handleServiceEdit() {
      navigate(`/monitors/${monitorId}/edit`);
    }

    async function handleServiceDelete() {
      try {
        const response = await API.fetch('DELETE', `/API/v1/services/${monitorId}`);
        if (response.status === 204) {
          navigate(`/`);
        }
      } catch(e) {
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

    let heartbeats : Heartbeat[] = [...new Array(size)].map((_, i) => ({
      created_at: new Date("2022-11-05T18:19:14.843001+01:00"),
      id: -size + i,
      response_time: 0,
      service_id: service.id,
      status: "UNKNOWN",
      status_code: 0
    }))

    for (let i = 0; i < latencies.length; i++) {
      heartbeats[(size - latencies.length) + i] = latencies[i]
    }

    return (
      <>
        <div className={styles.header}>
          <div>
            <h1 className={styles.title}>{service.name}</h1>
            <h4><a href={service.url} target='_blank' rel="noreferrer">{service.url}</a></h4>
          </div>
          <div className={styles.controls}>
            <Button
                text={(service?.enabled) ? 'Stop' : 'Start'}
                icon={(service?.enabled) ? BsStop : BsPlay}
                iconPosition="left"
                size="sm"
                color={(service?.enabled) ? 'gray' : 'green'}
                importance="primary"
                handleClick={handleServiceEnableToggle}
                marginTop="mt-0"
            />
            <Button
                text="Edit"
                icon={FaPencilAlt}
                iconPosition="left"
                size="sm"
                color="orange"
                importance="primary"
                handleClick={handleServiceEdit}
                marginTop="mt-0"
            />
            <Button
                text="Delete"
                icon={FaTrashAlt}
                iconPosition="left"
                size="sm"
                color="red"
                importance="primary"
                handleClick={handleServiceDelete}
                marginTop="mt-0"
            />
          </div>
        </div>
        <MonitorCards monitorId={monitorId} />
        <div className='block'>
        <Tracking marginTop="mt-2">
            {heartbeats.map(heartbeat => {
                return (
                    <TrackingBlock
                    key={heartbeat.id}
                    color={STATUS_COLORS_MAP[heartbeat.status]}
                    tooltip={`Response time: ${heartbeat.response_time} ms`}
                    />
                    );
                })}
        </Tracking>
        </div>
        <div className='block'>
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
            colors={["yellow"]}
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
        </div>
      </>
    )
  }