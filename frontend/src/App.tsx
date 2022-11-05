import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import './App.css';
import { HomePage } from './pages/HomePage/HomePage';
import { Layout } from './pages/Layout/Layout';
import { MonitorPage } from './pages/MonitorPage/MonitorPage';
import '@tremor/react/dist/esm/tremor.css';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route element={<Layout />}>
          <Route path="/" element={<HomePage />} />
          <Route path="/monitors/:monitorId" element={<MonitorPage />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}

export default App;