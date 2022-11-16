import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { MonotorsPage } from './pages/MonotorsPage/MonotorsPage';
import { NotificationsPage } from './pages/NotificationsPage/NotificationsPage';
import { Layout } from './pages/Layout/Layout';
import { MonitorPage } from './pages/MonitorPage/MonitorPage';
import { NotificationNewPage } from './pages/NotificationNewPage/NotificationNewPage';
import { NotificationEditPage } from './pages/NotificationEditPage/NotificationEditPage';
import { MonitorNewPage } from './pages/MonitorNewPage/MonitorNewPage';
import { MonitorEditPage } from './pages/MonitorEditPage/MonitorEditPage';
import '@tremor/react/dist/esm/tremor.css';
import './App.css';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route element={<Layout />}>
          <Route path="/" element={<MonotorsPage />} />
          <Route path="/monitors/" element={<MonotorsPage />} />
          <Route path="/notifications/" element={<NotificationsPage />} />
          <Route path="/notifications/:notificationName/edit" element={<NotificationEditPage />} />
          <Route path="/notifications/new" element={<NotificationNewPage />} />
          <Route path="/monitors/:monitorId" element={<MonitorPage />} />
          <Route path="/monitors/:monitorId/edit" element={<MonitorEditPage />} />
          <Route path="/monitors/new" element={<MonitorNewPage />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}

export default App;
