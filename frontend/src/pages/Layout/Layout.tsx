
import { Outlet, NavLink } from 'react-router-dom';
import { FaHeartbeat, FaChevronRight, FaBell } from 'react-icons/fa';
import { BsStack } from 'react-icons/bs';
import { Icon } from "@tremor/react";
import { useState } from 'react';

import styles from './Layout.module.css';


export function Layout() {
  const [sidebarOpen, setSidebarOpen] = useState(true)
  return (
    <div className='dark2'>
        <nav className={(sidebarOpen) ? `${styles.sidebar} ${styles.open}` : styles.sidebar}>
            <div className={styles.logo}>
                <span className="image">
                    <FaHeartbeat size={'40px'} />
                </span>

                <div className={styles.logoText}>
                    <span className="name">Uptime</span>
                </div>
                <span className={styles.toggle} onClick={() => {console.log('111'); setSidebarOpen(!sidebarOpen)}}>
                    <FaChevronRight size={'20px'} />
                </span>
            </div>
            <div className={styles.sidebarMenu}>
                <ul className={styles.menuLinks}>
                    <li className={styles.navLink}>
                        <NavLink 
                            className={(navData) => (navData.isActive ? styles.active : '')}
                            to='/monitors/'
                        >
                            <BsStack size={'25px'} />
                            <span className={styles.navLinkTitle}>Services</span>
                        </NavLink>
                    </li>
                    <li className={styles.navLink}>
                        <NavLink 
                            className={(navData) => (navData.isActive ? styles.active : '')}
                            to='/notifications/'
                        >
                            <FaBell size={'25px'} />
                            <span className={styles.navLinkTitle}>Notifications</span>
                        </NavLink>
                    </li>
                </ul>
            </div>
        </nav>
        <main className={(sidebarOpen) ? styles.open : ''}>
            <div className={styles.main}>
                <Outlet />
            </div>
        </main>
      </div>
  );
}