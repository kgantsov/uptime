import { Divider } from '@tremor/react';
import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Link } from 'react-router-dom';
import { FaTelegramPlane } from 'react-icons/fa';
import { Button } from "@tremor/react";
import { FaTrashAlt, FaPencilAlt } from 'react-icons/fa';
import { useForm, SubmitHandler } from "react-hook-form";

import styles from './NotificationEditPage.module.css';
import { API } from '../../API';

type Inputs = {
  name: string,
  callback: string,
  callback_chat_id: string,
};

export function NotificationEditPage() {
  let navigate = useNavigate();
  const { notificationName } = useParams();

  const { register, setValue, handleSubmit, watch, formState: { errors } } = useForm<Inputs>();

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
    API.fetch('GET', `/API/v1/notifications/${notificationName}`)
    .then(resp => resp.json())
    .then((data) => {
      setValue('name', data.name);
      setValue('callback', data.callback);
      setValue('callback_chat_id', data.callback_chat_id);
    })
  }, [notificationName])

  async function handleNotificationDelete() {

    try {
      const response = await API.fetch('DELETE', `/API/v1/notifications/${notificationName}`);
      if (response.status === 204) {
        navigate(`/notifications/`);
      }
    } catch(e) {
      console.log(e);
    }
  }

  const onSubmit: SubmitHandler<Inputs> = data => {
    API.fetch('PATCH', `/API/v1/notifications/${notificationName}`, null, {
      name: data.name,
      callback: data.callback,
      callback_chat_id: data.callback_chat_id,
      callback_type: 'TELEGRAM',
      // callback_type: values.callback_type,
    }).then(resp => resp.json()).then((data) => {
      navigate(`/notifications`);
    });
  };

  return (
    <>
        <>
        <div className='block'>
          <h1>Edit Notification</h1>
          <form onSubmit={handleSubmit(onSubmit)}>
          <div className={(errors.name) ? "form-element error" : "form-element"}>
              <label htmlFor="name">Name</label>
              <input
                type="text"
                {...register("name", { required: 'Name is required' })}
                disabled
              />
              <div className="error-message">{errors.name?.message}</div>
            </div>

            <div className={(errors.callback) ? "form-element error" : "form-element"}>
              <label htmlFor="callback">Callback</label>
              <input
                type="text"
                {...register("callback", { required: 'Callback is required' })}
              />
              <div className="error-message">{errors.callback?.message}</div>
            </div>

            <div className={(errors.callback_chat_id) ? "form-element error" : "form-element"}>
              <label htmlFor="callback_chat_id">Callback chat ID</label>
              <input
                type="text"
                {...register("callback_chat_id", { required: 'Callback chat ID is required' })}
              />
              <div className="error-message">{errors.callback_chat_id?.message}</div>
            </div>

            <div className="form-element">
              <div className="submit-wrapper">
                <Button
                    text="Save"
                    icon={undefined}
                    iconPosition="left"
                    size="sm"
                    color="green"
                    importance="primary"
                    handleClick={handleSubmit(onSubmit)}
                    marginTop="mt-0"
                />
                <Button
                    text="Delete"
                    icon={FaTrashAlt}
                    iconPosition="left"
                    size="sm"
                    color="red"
                    importance="primary"
                    handleClick={handleNotificationDelete}
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