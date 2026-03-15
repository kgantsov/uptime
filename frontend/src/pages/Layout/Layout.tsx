import { Outlet, NavLink, Link } from "react-router-dom";
import {
  FaHeartbeat,
  FaChevronRight,
  FaBell,
  FaBars,
  FaTimes,
} from "react-icons/fa";
import { BsStack } from "react-icons/bs";
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
      {/* Overlay for mobile when sidebar is open */}
      {isMobile && sidebarOpen && (
        <div className={styles.overlay} onClick={() => setSidebarOpen(false)} />
      )}

      {/* Sticky top bar shown only on mobile — contains the hamburger */}
      {isMobile && (
        <div className={styles.mobileHeader}>
          <span
            className={styles.mobileMenuBtn}
            onClick={handleToggle}
            aria-label="Open menu"
          >
            <FaBars size={"18px"} />
          </span>
          <Link to="/" className={styles.mobileHeaderBrand}>
            <span className={styles.mobileHeaderLogo}>
              <FaHeartbeat size={"22px"} />
            </span>
            <span className={styles.mobileHeaderTitle}>Uptime</span>
          </Link>
        </div>
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

          {/* Desktop: chevron toggle */}
          {!isMobile && (
            <span className={styles.toggle} onClick={handleToggle}>
              <FaChevronRight size={"20px"} />
            </span>
          )}

          {/* Mobile: close button inside the sidebar logo row */}
          {isMobile && sidebarOpen && (
            <span
              className={styles.mobileCloseBtn}
              onClick={() => setSidebarOpen(false)}
              aria-label="Close menu"
            >
              <FaTimes size={"18px"} />
            </span>
          )}
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
