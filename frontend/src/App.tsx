import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { HomePage } from './pages/HomePage/HomePage';
import { NotificationsPage } from './pages/NotificationsPage/NotificationsPage';
import { Layout } from './pages/Layout/Layout';
import { MonitorPage } from './pages/MonitorPage/MonitorPage';
import { NewMonitorPage } from './pages/NewMonitorPage/NewMonitorPage';
import '@tremor/react/dist/esm/tremor.css';
import './App.css';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route element={<Layout />}>
          <Route path="/" element={<HomePage />} />
          <Route path="/monitors/" element={<HomePage />} />
          <Route path="/notifications/" element={<NotificationsPage />} />
          <Route path="/monitors/:monitorId" element={<MonitorPage />} />
          <Route path="/monitors/new" element={<NewMonitorPage />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}

export default App;
