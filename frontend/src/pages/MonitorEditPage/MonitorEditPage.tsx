import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { API } from '../../API';
import { Button } from "@tremor/react";
import Async from 'react-select/async';
import { useForm, SubmitHandler } from "react-hook-form";
import { yupResolver } from '@hookform/resolvers/yup';
import * as yup from "yup";

import styles from './MonitorPage.module.css';

type Inputs = {
  name: string,
  url: string,
  check_interval: number,
  accepted_status_code: number,
  timeout: number,
};

const schema = yup.object({
  name: yup.string().required(),
  url: yup.string().required(),
  check_interval: yup.number().integer().min(5).max(3600).required(),
  accepted_status_code: yup.number().min(100).max(600).integer().required(),
  timeout: yup.number().integer().min(0).max(120).required(),
}).required();


export function MonitorEditPage() {
  let navigate = useNavigate();
  const { monitorId } = useParams();
  const { register, setValue, handleSubmit, watch, formState: { errors } } = useForm<Inputs>({
    resolver: yupResolver(schema)
  });

  const [notifications, setNotifications] = useState<any[]>([]);

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

  useEffect(() => {
    API.fetch('GET', `/API/v1/services/${monitorId}`).then((data) => {
      setValue('name', data.name);
      setValue('url', data.url);
      setValue('check_interval', data.check_interval);
      setValue('accepted_status_code', data.accepted_status_code);
      setValue('timeout', data.timeout);
      setNotifications(data.notifications.map((notification: { name: any; }) => {
        return { value: notification.name, label: notification.name, ...notification };
      }))
    })
  }, [monitorId])

  const onSubmit: SubmitHandler<Inputs> = data => {

    API.fetch('PATCH', `/API/v1/services/${monitorId}`, null, {
      name: data.name,
      url: data.url,
      check_interval: data.check_interval,
      timeout: data.timeout,
      accepted_status_code: data.accepted_status_code,
      notifications: notifications,
    }).then((data) => {
      navigate(`/monitors/${monitorId}`);
    });
  };

  return (
    <>
      <div className='block'>
        <h1>Edit Monitor</h1>
        <form onSubmit={handleSubmit(onSubmit)}>
          <div className={(errors.name) ? "form-element error" : "form-element"}>
            <label htmlFor="name">Name</label>
            <input
              type="text"
              {...register("name", { required: 'Name is required' })}
            />
            <div className="error-message">{errors.name?.message}</div>
          </div>

          <div className={(errors.url) ? "form-element error" : "form-element"}>
            <label htmlFor="url">URL</label>
            <input
              type="text"
              {...register("url", { required: 'URL is required' })}
            />
            <div className="error-message">{errors.url?.message}</div>
          </div>

          <div className={(errors.check_interval) ? "form-element error" : "form-element"}>
            <label htmlFor="check_interval">check interval</label>
            <input
              type="number"
              {...register("check_interval", { required: 'Check interval is required' })}
            />
            <div className="error-message">{errors.check_interval?.message}</div>
          </div>

          <div className={(errors.timeout) ? "form-element error" : "form-element"}>
            <label htmlFor="timeout">Timeout</label>
            <input
              type="number"
              {...register("timeout", { required: 'Timeout interval is required' })}
            />
            <div className="error-message">{errors.timeout?.message}</div>
          </div>

          <div className={(errors.accepted_status_code) ? "form-element error" : "form-element"}>
            <label htmlFor="accepted_status_code">Accepted status code</label>
            <input
              type="number"
              {...register("accepted_status_code", { required: 'Accepted status code is required' })}
            />
            <div className="error-message">{errors.accepted_status_code?.message}</div>
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
            </div>
          </div>
        </form>
      </div>
    </>
  )
}