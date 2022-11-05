
import { useState, useEffect } from 'react';

import {
  Tracking,
  TrackingBlock,
  Flex,
  Text,
 } from "@tremor/react";
import { Outlet, NavLink, Link } from 'react-router-dom';
import { PlusIcon } from '@heroicons/react/24/solid'
import { Icon } from "@tremor/react";

import styles from './Layout.module.css';
import { Service } from '../../types/services';
import { Heartbeat, STATUS_COLORS_MAP } from '../../types/heartbeats';


export function Layout() {
    const [services, setServices] = useState<Service[]>([]);
    const [stats, setStats] = useState<Heartbeat[]>([]);
    const size = 20;

    async function fetchServices() {
        try {
            const response = await fetch('/API/v1/services');
            const data = await response.json();
            setServices(data)
        } catch(e) {
            console.log(e);
        }
    }

    async function fetchStats() {
        try {
            const response = await fetch(`/API/v1/heartbeats/stats?size=${size}`);
            const data = await response.json();
            setStats(data)
        } catch (e) {
            console.log(e);
        }
    }

    async function fetchData() {
        await fetchServices()

        await fetchStats()
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
        <header>Uptime</header>
        <main>
            <div className={styles.sidebar}>
                <div className={styles.monitors}>
                <div className={styles.monitorHeader}>
                    <div>
                        <h3>Monotors</h3>
                    </div>
                    <div>
                        <Link to='/monitors/new'>
                            <Icon
                                icon={ PlusIcon }
                                variant="simple"
                                tooltip=""
                                size="sm"
                                color="blue"
                                marginTop="mt-0"
                            />
                        </Link>
                    </div>
                </div>
                <ul className={styles.monitorList}>
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
                        return (
                            <li key={service.id}>
                                <NavLink
                                    className={(navData) => (
                                        navData.isActive ? `${styles.monitorItem} ${styles.active}` : styles.monitorItem
                                    )}
                                    to={`/monitors/${service.id}`}
                                >
                                    {service.name}
                                    <Flex justifyContent="justify-end" marginTop="mt-4">
                                        <Text>Uptime {(success * 100 / heartbeats.length) || 0}%</Text>
                                    </Flex>
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
                                </NavLink>
                            </li>
                        )
                    })}
                </ul>
                </div>
            </div>
            <div className={styles.main}>
                <Outlet />
            </div>
        </main>
      </>
  );
}