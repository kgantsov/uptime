import { Divider } from '@tremor/react';
import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Link } from 'react-router-dom';
import { FaTelegramPlane } from 'react-icons/fa';

import styles from './NotificationEditPage.module.css';
import { API } from '../../API';

export function NotificationEditPage() {
  let navigate = useNavigate();
  const { notificationName } = useParams();

  const [values, setValues] = useState({
    name: '',
    callback: '',
    callback_chat_id: '',
    // callback_type: '',
  });

  const handleChange = (e: { target: { name: any; value: any; }; }) => {
    setValues((oldValues) => ({
      ...oldValues,
      [e.target.name]: e.target.value,
    }));
  };

  useEffect(() => {
    API.fetch('GET', `/API/v1/notifications/${notificationName}`).then((data) => {
      setValues(data)
    })
  }, [notificationName])

  const handleSubmit = (e?: { preventDefault: () => void; }) => {
    if (e !== undefined) {
      e.preventDefault();
    }

    API.fetch('PATCH', `/API/v1/notifications/${notificationName}`, null, {
      name: values.name,
      callback: values.callback,
      callback_chat_id: values.callback_chat_id,
      callback_type: 'TELEGRAM',
      // callback_type: values.callback_type,
    }).then((data) => {
      navigate(`/notifications`);
    });
  };

  return (
    <>
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
              <label htmlFor="callback">Callback</label>
              <input
                className=""
                id="callback"
                name="callback"
                type="text"
                value={values.callback}
                onChange={handleChange}
                required
              />
            </div>

            <div className="form-element">
              <label htmlFor="callback_chat_id">Callback chat ID</label>
              <input
                className=""
                id="callback_chat_id"
                name="callback_chat_id"
                type="text"
                value={values.callback_chat_id}
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
    </>
  );
}