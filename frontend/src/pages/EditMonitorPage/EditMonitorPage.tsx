import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { API } from '../../API';
import { Button } from "@tremor/react";

import styles from './MonitorPage.module.css';

export function EditMonitorPage() {
  let navigate = useNavigate();
  const { monitorId } = useParams();

  const [values, setValues] = useState({
    name: '',
    url: '',
    check_interval: '',
    timeout: '',
  });

  const handleChange = (e: { target: { name: any; value: any; }; }) => {
    setValues((oldValues) => ({
      ...oldValues,
      [e.target.name]: e.target.value,
    }));
  };

  useEffect(() => {
    API.fetch('GET', `/API/v1/services/${monitorId}`).then((data) => {
      setValues(data)
    })
  }, [monitorId])

  const handleSubmit = (e?: { preventDefault: () => void; }) => {
    if (e !== undefined) {
      e.preventDefault();
    }

    API.fetch('PATCH', `/API/v1/services/${monitorId}`, null, {
      name: values.name,
      url: values.url,
      check_interval: Number.parseInt(values.check_interval),
      timeout: Number.parseInt(values.timeout),
    }).then((data) => {
      navigate(`/monitors/${monitorId}`);
    });
  };

    return (
      <>
        <div className='block'>
          <h1>Edit Monitor</h1>
          <form method="post" onSubmit={handleSubmit}>
            <div className="form-element">
              <label htmlFor="name">Name</label>
              <input
                className=""
                id="name"
                name="name"
                type="text"
                value={values.name}
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
                value={values.url}
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
                value={values.check_interval}
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
                value={values.timeout}
                onChange={handleChange}
                required
              />
            </div>

            <div className="form-element">
              <div className="submit-wrapper">
                <input type="submit" value="Save" />
              </div>
            </div>
          </form>
        </div>
      </>
    )
  }