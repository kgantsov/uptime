import { useState } from 'react';
import { useNavigate } from 'react-router-dom';

import styles from './NotificationsPage.module.css';
import { API } from '../../API';
import { Button } from "@tremor/react";
import { FaPlus } from 'react-icons/fa';

export function NotificationNewPage() {
  let navigate = useNavigate();

  const [values, setValues] = useState({
    name: '',
    callback: '',
    callback_chat_id: '',
  });

  const handleChange = (e: { target: { name: any; value: any; }; }) => {
    setValues((oldValues) => ({
      ...oldValues,
      [e.target.name]: e.target.value,
    }));
  };

  const handleSubmit = (e?: { preventDefault: () => void; }) => {
    if (e !== undefined) {
      e.preventDefault();
    }

    API.fetch('POST', '/API/v1/notifications', null, {
      name: values.name,
      callback: values.callback,
      callback_chat_id: values.callback_chat_id,
      callback_type: 'TELEGRAM',
    }).then((data) => {
      navigate(`/notifications`);
    });
  };

  return (
    <>
        <>
        <div className='block'>
          <h1>Add Notifications</h1>
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
                <Button
                    text="Add"
                    icon={FaPlus}
                    iconPosition="left"
                    size="sm"
                    color="green"
                    importance="primary"
                    handleClick={handleSubmit}
                    marginTop="mt-0"
                />
              </div>
            </div>
          </form>
        </div>
      </>
    </>
  );
}