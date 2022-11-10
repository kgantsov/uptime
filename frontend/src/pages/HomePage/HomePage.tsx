
import { useState, useEffect } from 'react';

import {
  Tracking,
  TrackingBlock,
} from "@tremor/react";
import { Outlet, NavLink, Link, useNavigate } from 'react-router-dom';
import { FaHeart, FaCheckCircle, FaExclamationCircle } from 'react-icons/fa';
import { FaPlus } from 'react-icons/fa';
import { Button } from "@tremor/react";
import { Icon } from "@tremor/react";
import { Badge } from "@tremor/react";

import styles from './HomePage.module.css';
import { Service } from '../../types/services';
import { Heartbeat, STATUS_COLORS_MAP } from '../../types/heartbeats';


export function HomePage() {
    let navigate = useNavigate();
    const [services, setServices] = useState<Service[]>([]);
    const [stats, setStats] = useState<Heartbeat[]>([]);
    const size = 20;

    async function fetchServices() {
        try {
            const response = await fetch('/API/v1/services');
            const data = await response.json();
            return data
        } catch(e) {
            console.log(e);
        }
    }

    async function fetchStats() {
        try {
            const response = await fetch(`/API/v1/heartbeats/stats?size=${size}`);
            const data = await response.json();
            return data
        } catch (e) {
            console.log(e);
        }
    }

    async function fetchData() {
        const servicesData = await fetchServices()

        const statsData = await fetchStats()

        if (statsData) {
            setStats(statsData)
        }
        if (servicesData) {
            setServices(servicesData)
        }
    }

    useEffect(() => {
        fetchData();
        const interval = setInterval(() => {
            fetchData();
        }, 5000);
        return () => clearInterval(interval);
    }, []);

  return (
    <>
        <div className={styles.controls}>
            <Button
                text="new"
                icon={FaPlus}
                iconPosition="left"
                size="sm"
                color="green"
                importance="primary"
                handleClick={() => navigate('/monitors/new')}
                // disabled={false}
                marginTop="mt-0"
            />
        </div>
        <table className={styles.monitorList}>
            <thead>
                <tr>
                    <td>Status</td>
                    <td>Service</td>
                    <td>Health</td>
                    <td>Heartbeats</td>
                    <td>Response&nbsp;time</td>
                </tr>
            </thead>
            <tbody>
                {services.map(service => {
                    let _heartbeats = stats.filter(i => i.service_id === service.id);

                    let heartbeats = [...new Array(size)].map((_, i) => {
                        return _heartbeats[i] || {
                            created_at: "2022-11-05T18:19:14.843001+01:00",
                            id: -size + i,
                            response_time: 0,
                            service_id: service.id,
                            status: "UNKNOWN",
                            status_code: 0
                        }
                    });
                    heartbeats = heartbeats.reverse()

                    const success = heartbeats.reduce(
                        (prev, cur) => (cur.status === "UP" || cur.status === "UNKNOWN") ? prev + 1 : prev, 0
                    )
                    const lastHeartbeats = heartbeats[heartbeats.length - 1]
                    const isUp = (
                        lastHeartbeats.status === 'UP' || lastHeartbeats.status === 'UNKNOWN'
                    )
                    return (
                        <tr key={service.id}>
                            <td>
                                <Badge
                                    text={isUp ? 'Up' : 'Down'}
                                    color={isUp ? "green" : "rose"}
                                    size="sm"
                                    icon={isUp ? FaCheckCircle : FaExclamationCircle}
                                    tooltip=""
                                    marginTop="mt-0" />
                            </td>
                            <td className={styles.title}>
                                <NavLink to={`/monitors/${service.id}`}>
                                    <div className={styles.name}>{service.name}</div>
                                    <div className={styles.url}>{service.url}</div>
                                </NavLink>
                            </td>
                            <td>
                                <Badge
                                    text={`${(success * 100 / heartbeats.length) || 0}%`}
                                    color={"green"}
                                    size="sm"
                                    icon={FaHeart}
                                    tooltip=""
                                    marginTop="mt-0" />
                            </td>
                            <td className={styles.heartbeats}>
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
                            </td>
                            <td>{lastHeartbeats.response_time} ms</td>
                        </tr>
                    )
                })}
            </tbody>
        </table>
    </>
  );
}