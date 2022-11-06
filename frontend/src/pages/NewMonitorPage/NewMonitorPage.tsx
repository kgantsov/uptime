import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { API } from '../../API';

import styles from './MonitorPage.module.css';

export function NewMonitorPage() {
  let navigate = useNavigate();

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

  const handleSubmit = (e: { preventDefault: () => void; }) => {
    e.preventDefault();

    console.log('=====>', values)

    API.fetch('POST', `/API/v1/services`, null, {
      name: values.name,
      url: values.url,
      check_interval: Number.parseInt(values.check_interval),
      timeout: Number.parseInt(values.timeout),
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
              <div className="submit-wrapper">
                <input type="submit" value="Create" />
              </div>
            </div>
          </form>
        </div>
      </>
    )
  }