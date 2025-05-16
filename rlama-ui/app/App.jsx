import React, { useState } from 'react';
import { Layout, Menu, Button, Typography } from 'antd';
import {
  HomeOutlined,
  PlusOutlined,
  DatabaseOutlined,
  QuestionCircleOutlined,
  BookOutlined
} from '@ant-design/icons';
import { createHashRouter, RouterProvider, Link, Outlet, useLocation } from 'react-router-dom';

// Importation des composants de pages
import Home from './views/Home';
import Dashboard from './views/Dashboard';
import CreateRag from './views/CreateRag';
import RagDetail from './views/RagDetail';

const { Header, Sider, Content } = Layout;
const { Title } = Typography;

// Composant de layout principal
const MainLayout = () => {
  const [collapsed, setCollapsed] = useState(false);
  const location = useLocation();

  return (
    <Layout style={{ minHeight: '100vh', position: 'relative', overflow: 'hidden' }}>
      <Sider 
        collapsible 
        collapsed={collapsed} 
        onCollapse={setCollapsed}
        breakpoint="lg"
        collapsedWidth={0}
        style={{ boxShadow: 'var(--shadow-md)' }}
      >
        <div className="logo">
          {!collapsed && "RLAMA UI"}
        </div>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={[location.pathname]}
          items={[
            {
              key: '/',
              icon: <HomeOutlined />,
              label: <Link to="/">Home</Link>,
            },
            {
              key: '/systems',
              icon: <DatabaseOutlined />,
              label: <Link to="/systems">My RAG Systems</Link>,
            },
            {
              key: '/create',
              icon: <PlusOutlined />,
              label: <Link to="/create">New RAG</Link>,
            },
            {
              type: 'divider',
            },
            {
              key: 'docs',
              icon: <BookOutlined />,
              label: <a href="https://github.com/DonTizi/rlama" target="_blank" rel="noopener noreferrer">Documentation</a>,
            },
          ]}
        />
      </Sider>
      <Layout className="site-layout">
        <Header style={{ padding: 0, background: '#fff', boxShadow: 'var(--shadow-sm)', position: 'relative', overflow: 'hidden' }}>
          <div className="flex justify-between items-center px-4" style={{ height: '100%', overflow: 'hidden' }}>
            <div className="flex items-center flex-shrink-0">
              <Title level={4} style={{ margin: 0, color: 'var(--primary-700)', display: 'flex', alignItems: 'center' }}>
                <span style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                  RLAMA
                </span>
              </Title>
              <div className="connection-status ml-4">
                <span className="status-dot"></span>
                <span className="status-text">Connected</span>
              </div>
            </div>
            <Button type="text" icon={<QuestionCircleOutlined />} onClick={() => window.open('https://github.com/DonTizi/rlama', '_blank')}>
              Help
            </Button>
          </div>
        </Header>
        <Content className="content-container">
          <div className="fade-in">
            <Outlet />
          </div>
        </Content>
      </Layout>
    </Layout>
  );
};

// Définition des routes de l'application
const router = createHashRouter([
  {
    path: "/",
    element: <MainLayout />,
    children: [
      {
        index: true,
        element: <Home />,
      },
      {
        path: "systems",
        element: <Dashboard />,
      },
      {
        path: "create",
        element: <CreateRag />,
      },
      {
        path: "rag/:ragName",
        element: <RagDetail />,
      },
    ],
  },
]);

// Composant App qui intègre le routeur
const App = () => {
  return <RouterProvider router={router} />;
};

export default App; 