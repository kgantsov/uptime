
import { useState, useEffect } from 'react';

import {
  Tracking,
  TrackingBlock,
} from "@tremor/react";
import { Outlet, NavLink, Link } from 'react-router-dom';
import { BellIcon, Square3Stack3DIcon, ChevronLeftIcon, ChevronRightIcon } from '@heroicons/react/24/solid'
import { FaHeartbeat } from 'react-icons/fa';
import { BsChevronRight, BsChevronLeft } from 'react-icons/bs';
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
    <div className='dark2'>
        <nav className={styles.sidebar}>
            <div className={styles.logo}>
                <span className="image">
                    <Icon
                        icon={FaHeartbeat}
                        variant="simple"
                        tooltip=""
                        size="lg"
                        color="blue"
                        marginTop="mt-0"
                    />
                </span>

                <div className="header-text">
                    <span className="name">Uptime</span>
                </div>
                <span className={styles.toggle}>
                    <Icon
                        icon={BsChevronRight}
                        variant="simple"
                        tooltip=""
                        size="md"
                        color="blue"
                        marginTop="mt-0"
                    />
                </span>
            </div>
            <div className={styles.sidebarMenu}>
                <ul className={styles.menuLinks}>
                    <li className={styles.navLink}>
                        <NavLink 
                            className={(navData) => (navData.isActive ? styles.active : '')}
                            to='/monitors/'
                        >
                            <Icon
                                icon={Square3Stack3DIcon}
                                variant="simple"
                                tooltip=""
                                size="md"
                                color="blue"
                                marginTop="mt-0"
                            />
                            <span>Services</span>
                        </NavLink>
                    </li>
                    <li className={styles.navLink}>
                        <NavLink 
                            className={(navData) => (navData.isActive ? styles.active : '')}
                            to='/notifications/'
                        >
                            <Icon
                                icon={BellIcon}
                                variant="simple"
                                tooltip=""
                                size="md"
                                color="blue"
                                marginTop="mt-0"
                            />
                            <span>Notifications</span>
                        </NavLink>
                    </li>
                </ul>
            </div>
        </nav>
        <main>
            <div className={styles.main}>
                <Outlet />
            </div>
        </main>
      </div>
  );
}