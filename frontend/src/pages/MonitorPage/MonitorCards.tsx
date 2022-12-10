import { useState, useEffect, SetStateAction } from 'react';
import { HeartbeatStats, STATUS_COLORS_MAP } from '../../types/heartbeats';
import { Card, Metric, Title, Text, Flex, CategoryBar, DonutChart, Divider } from "@tremor/react";

import styles from './MonitorPage.module.css';
import { API } from '../../API';


interface Props {
    monitorId: string | undefined;
}

interface StatsData {
    UP: HeartbeatStats;
    DOWN: HeartbeatStats;
    TIMEOUT: HeartbeatStats;
    FAILED: HeartbeatStats;
    UNKNOWN: HeartbeatStats;
}


export function MonitorCards({ monitorId }: Props): JSX.Element {

    const [statsDay, setStatsDay] = useState<StatsData|undefined>(undefined);
    const [statsWeek, setStatsWeek] = useState<StatsData|undefined>(undefined);
    const [statsMonth, setStatsMonth] = useState<StatsData|undefined>(undefined);

    async function fetchStatsData(days: number, setDataFunc: { (value: SetStateAction<StatsData | undefined>): void; (value: SetStateAction<StatsData | undefined>): void; (arg0: any): void; }) {
        try {
          const response = await fetch(`/API/v1/heartbeats/stats/${days}`);
          const data = await response.json();
          setDataFunc(
            data
              .filter((x: HeartbeatStats) => x.service_id === Number(monitorId))
              .reduce(
                (a: any, v: HeartbeatStats) => ({ ...a, [v.status]: v}),
                {}
              )
          )
        } catch(e) {
            console.log(e);
        }
      }


    useEffect(() => {
      fetchStatsData(1, setStatsDay);
      fetchStatsData(7, setStatsWeek);
      fetchStatsData(30, setStatsMonth);

      const interval = setInterval(() => {
        fetchStatsData(1, setStatsDay);
        fetchStatsData(7, setStatsWeek);
        fetchStatsData(30, setStatsMonth);
      }, 5000);
      return () => clearInterval(interval);
    }, [monitorId]);

    if (!statsDay || !statsWeek || !statsMonth) {
      return <></>
    }

    const calcPercentage = (data: StatsData): number => {
        const success = data?.UP?.counter | 0;
        const total = Object.values(data).reduce((a: number, v: HeartbeatStats) => a + v.counter, 0);
        
        return (success > 0) ? 100 / total * success : 0
    }

    return (
        <div className='cards'>
            <Card maxWidth="max-w-xs" decoration="top" decorationColor="orange">
                <Metric>24 hours</Metric>
                <Text>Avg. Response</Text>
                <Title>{Math.round(statsDay['UP']?.average_response_time || 0)} ms</Title>
                <Text>Uptime</Text>
                <Title>{calcPercentage(statsDay).toFixed(2)}%</Title>
            </Card>
            <Card maxWidth="max-w-xs" decoration="top" decorationColor="orange">
                <Metric>7 days</Metric>
                <Text>Avg. Response</Text>
                <Title>{Math.round(statsWeek['UP']?.average_response_time || 0)} ms</Title>
                <Text>Uptime</Text>
                <Title>{calcPercentage(statsWeek).toFixed(2)}%</Title>
            </Card>
            <Card maxWidth="max-w-xs" decoration="top" decorationColor="orange">
                <Metric>30 days</Metric>
                <Text>Avg. Response</Text>
                <Title>{Math.round(statsMonth['UP']?.average_response_time || 0)} ms</Title>
                <Text>Uptime</Text>
                <Title>{calcPercentage(statsMonth).toFixed(2)}%</Title>
            </Card>
        </div>
    )
  }