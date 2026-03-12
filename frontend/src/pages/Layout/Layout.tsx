import { Outlet, NavLink } from "react-router-dom";
import {
  FaHeartbeat,
  FaChevronRight,
  FaBell,
  FaBars,
  FaTimes,
} from "react-icons/fa";
import { BsStack } from "react-icons/bs";
import { Icon } from "@tremor/react";
import { useState, useEffect } from "react";

import styles from "./Layout.module.css";

export function Layout() {
  const [sidebarOpen, setSidebarOpen] = useState(true);
  const [isMobile, setIsMobile] = useState(false);

  useEffect(() => {
    const checkMobile = () => {
      const mobile = window.innerWidth <= 768;
      setIsMobile(mobile);
      if (mobile) {
        setSidebarOpen(false);
      } else {
        setSidebarOpen(true);
      }
    };

    checkMobile();
    window.addEventListener("resize", checkMobile);
    return () => window.removeEventListener("resize", checkMobile);
  }, []);

  const handleToggle = () => {
    setSidebarOpen(!sidebarOpen);
  };

  const handleLinkClick = () => {
    if (isMobile) {
      setSidebarOpen(false);
    }
  };

  return (
    <div className="dark2">
      {isMobile && sidebarOpen && (
        <div className={styles.overlay} onClick={() => setSidebarOpen(false)} />
      )}
      <nav
        className={
          sidebarOpen ? `${styles.sidebar} ${styles.open}` : styles.sidebar
        }
      >
        <div className={styles.logo}>
          <span className="image">
            <FaHeartbeat size={"40px"} />
          </span>

          <div className={styles.logoText}>
            <span className="name">Uptime</span>
          </div>
          <span className={styles.toggle} onClick={handleToggle}>
            {isMobile ? (
              sidebarOpen ? (
                <FaTimes size={"20px"} />
              ) : (
                <FaBars size={"20px"} />
              )
            ) : (
              <FaChevronRight size={"20px"} />
            )}
          </span>
        </div>
        <div className={styles.sidebarMenu}>
          <ul className={styles.menuLinks}>
            <li className={styles.navLink}>
              <NavLink
                className={(navData) => (navData.isActive ? styles.active : "")}
                to="/monitors/"
                onClick={handleLinkClick}
              >
                <BsStack size={"25px"} />
                <span className={styles.navLinkTitle}>Services</span>
              </NavLink>
            </li>
            <li className={styles.navLink}>
              <NavLink
                className={(navData) => (navData.isActive ? styles.active : "")}
                to="/notifications/"
                onClick={handleLinkClick}
              >
                <FaBell size={"25px"} />
                <span className={styles.navLinkTitle}>Notifications</span>
              </NavLink>
            </li>
          </ul>
        </div>
      </nav>
      <main className={sidebarOpen && !isMobile ? styles.open : ""}>
        <div className={styles.main}>
          <Outlet />
        </div>
      </main>
    </div>
  );
}
