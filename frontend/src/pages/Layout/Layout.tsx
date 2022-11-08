
import { useState, useEffect } from 'react';

import {
  Tracking,
  TrackingBlock,
} from "@tremor/react";
import { Outlet, NavLink, Link } from 'react-router-dom';
import { PlusIcon, HeartIcon } from '@heroicons/react/24/solid'
import { Icon } from "@tremor/react";
import { Badge } from "@tremor/react";

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
        <header>
            <Link to="/">Uptime</Link>
        </header>
        <main>
            <div className={styles.main}>
                <Outlet />
            </div>
        </main>
      </>
  );
}