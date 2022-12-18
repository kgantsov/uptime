import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useForm, SubmitHandler } from "react-hook-form";

import styles from './NotificationsPage.module.css';
import { API } from '../../API';
import { Button } from "@tremor/react";
import { FaPlus } from 'react-icons/fa';

type Inputs = {
  name: string,
  callback: string,
  callback_chat_id: string,
};

export function NotificationNewPage() {
  let navigate = useNavigate();
  const { register, handleSubmit, watch, formState: { errors } } = useForm<Inputs>();
  const onSubmit: SubmitHandler<Inputs> = data => {
    API.fetch('POST', '/API/v1/notifications', null, {
      name: data.name,
      callback: data.callback,
      callback_chat_id: data.callback_chat_id,
      callback_type: 'TELEGRAM',
    }).then((data) => {
      navigate(`/notifications`);
    });
  }

  return (
    <>
        <>
        <div className='block'>
          <h1>Add Notifications</h1>
          <form onSubmit={handleSubmit(onSubmit)}>
            <div className={(errors.name) ? "form-element error" : "form-element"}>
              <label htmlFor="name">Name</label>
              <input
                type="text"
                {...register("name", { required: 'Name is required' })}
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
                    text="Add"
                    icon={FaPlus}
                    iconPosition="left"
                    size="sm"
                    color="green"
                    importance="primary"
                    handleClick={handleSubmit(onSubmit)}
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