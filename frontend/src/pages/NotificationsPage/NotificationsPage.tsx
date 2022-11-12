import { Divider } from '@tremor/react';
import { useState, useEffect } from 'react';
import { Notification } from '../../types/services';
import { Button } from "@tremor/react";
import { Link, useNavigate } from 'react-router-dom';
import { FaPlus, FaTelegramPlane } from 'react-icons/fa';

import styles from './NotificationsPage.module.css';


async function fetchNotifications() {
  try {
      const response = await fetch('/API/v1/notifications');
      const data = await response.json();
      return data
  } catch(e) {
      console.log(e);
  }
}


export function NotificationsPage() {
  let navigate = useNavigate();
  const [notifications, setNotifications] = useState<Notification[]>([]);

  async function fetchData() {
  
    const notificationsData = await fetchNotifications()
    if (notificationsData) {
      setNotifications(notificationsData)
    }
  }

  useEffect(() => {
    fetchData()
  }, []);

  return (
    <>
        <div className='pageHeader'>
            <div>
                <h2>Notifications</h2>
            </div>
            <div className='pageControls'>
              <Button
                  text="new"
                  icon={FaPlus}
                  iconPosition="left"
                  size="sm"
                  color="green"
                  importance="primary"
                  handleClick={() => navigate('/notifications/new')}
                  // disabled={false}
                  marginTop="mt-0"
              />
            </div>
        </div>
        <div className='block'>
        <table className={styles.monitorList}>
            <thead>
                <tr>
                    <td>Name</td>
                </tr>
            </thead>
            <tbody>
              {notifications.map(notification => {
                return (
                  <tr key={notification.id}>
                    <td>
                      <Link to={`/notifications/${notification.name}/edit`}>
                        <FaTelegramPlane size={'25px'}/> {notification.name}
                      </Link>
                    </td>
                  </tr>
                )
              })}
            </tbody>
          </table>
        </div>
    </>
  );
}