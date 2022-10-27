
import { useState, useEffect } from 'react';

import {
  Tracking,
  TrackingBlock,
  Flex,
  Text,
 } from "@tremor/react";
import { Outlet, NavLink } from 'react-router-dom';

import styles from './Layout.module.css';
import { Service } from '../../types/services';
import { Heartbeat } from '../../types/heartbeats';


export function Layout() {
    const [services, setServices] = useState<Service[]>([]);
    const [stats, setStats] = useState<Heartbeat[]>([]);

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
            const response = await fetch('/API/v1/heartbeats/stats?size=25');
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
                    <h3>Monotors</h3>
                </div>
                <ul className={styles.monitorList}>
                    {services.map(service => {
                        const heartbeats = stats.filter(i => i.service_id === service.id).reverse()
                        const success = heartbeats.reduce(
                            (prev, cur) => (cur.is_success === true) ? prev + 1 : prev, 0
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
                                        <Text>Uptime {success * 100 / heartbeats.length}%</Text>
                                    </Flex>
                                    <Tracking marginTop="mt-2">
                                        {heartbeats.map(heartbeat => {
                                                return (
                                                    <TrackingBlock
                                                        key={heartbeat.id}
                                                        color={(heartbeat.is_success === true) ? "emerald" : "rose"}
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