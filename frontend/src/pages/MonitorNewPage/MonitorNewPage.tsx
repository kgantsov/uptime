import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { API } from '../../API';
import Async from 'react-select/async';
import { Button } from "@tremor/react";

import styles from './MonitorNewPage.module.css';

export function MonitorNewPage() {
  let navigate = useNavigate();

  const [values, setValues] = useState({
    name: '',
    url: '',
    check_interval: '',
    timeout: '',
    accepted_status_code: '',
  });
  const [notifications, setNotifications] = useState<any[]>([]);

  const handleChange = (e: { target: { name: any; value: any; }; }) => {
    setValues((oldValues) => ({
      ...oldValues,
      [e.target.name]: e.target.value,
    }));
  };

  const promiseOptions = (inputValue: string): Promise<any> => {
    return new Promise((resolve) => {
      API.fetch('GET', `/API/v1/notifications?q=${inputValue}`).then((data) => {
        resolve(
          data
            .map((notification: { name: any; }) => {
              return { value: notification.name, label: notification.name, ...notification };
            }),
        );
      });
    });
  }

  const handleSubmit = (e?: { preventDefault: () => void; }) => {
    if (e !== undefined) {
      e.preventDefault();
    }

    API.fetch('POST', `/API/v1/services`, null, {
      name: values.name,
      url: values.url,
      check_interval: Number.parseInt(values.check_interval),
      accepted_status_code: Number.parseInt(values.accepted_status_code),
      timeout: Number.parseInt(values.timeout),
      notifications: notifications,
      enabled: true,
    }).then((data) => {
      navigate(`/monitors/${data['id']}`);
    });
  };

    return (
      <>
        <div className='block'>
          <h1>Add New Monitor</h1>
          <form method="post" onSubmit={handleSubmit}>
            <div className="form-element">
              <label htmlFor="name">Name</label>
              <input
                className=""
                id="name"
                name="name"
                type="text"
                onChange={handleChange}
                required
              />
            </div>

            <div className="form-element">
              <label htmlFor="url">URL</label>
              <input
                className=""
                id="url"
                name="url"
                type="text"
                onChange={handleChange}
                required
              />
            </div>

            <div className="form-element">
              <label htmlFor="check_interval">check interval</label>
              <input
                className=""
                id="check_interval"
                name="check_interval"
                type="number"
                onChange={handleChange}
                required
              />
            </div>

            <div className="form-element">
              <label htmlFor="timeout">Timeout</label>
              <input
                className=""
                id="timeout"
                name="timeout"
                type="number"
                onChange={handleChange}
                required
              />
            </div>

            <div className="form-element">
              <label htmlFor="accepted_status_code">Accepted status code</label>
              <input
                className=""
                id="accepted_status_code"
                name="accepted_status_code"
                type="number"
                min={100}
                max={600}
                value={values.accepted_status_code}
                onChange={handleChange}
                required
              />
            </div>

            <div className="form-element">
              <label htmlFor="callback_chat_id">Notifications</label>
              <Async
                className="react-select-container"
                classNamePrefix="react-select"
                cacheOptions
                defaultOptions
                isMulti={true}
                loadOptions={promiseOptions}
                placeholder='Select'
                name="narratives"
                value={notifications}
                onChange={(option: readonly any[]) => {
                  setNotifications([...option])
                }}
              />
            </div>

            <div className="form-element">
              <div className="submit-wrapper">
                <input type="submit" value="Create" />
              </div>
            </div>
          </form>
        </div>
      </>
    )
  }