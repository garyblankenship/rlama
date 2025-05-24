import React, { useState, useEffect } from 'react';
import { Layout, Menu, Button, Typography, Tooltip } from 'antd';
import {
  HomeOutlined,
  PlusOutlined,
  DatabaseOutlined,
  BookOutlined,
  RobotOutlined,
  MenuOutlined,
  GlobalOutlined,
  TwitterOutlined,
  LinkedinOutlined,
  GithubOutlined,
  SettingOutlined
} from '@ant-design/icons';
import { createHashRouter, RouterProvider, Link, Outlet, useLocation } from 'react-router-dom';

// Importation des composants de pages
import Home from './views/Home';
import Dashboard from './views/Dashboard';
import CreateRag from './views/CreateRag';
import RagDetail from './views/RagDetail';
import AgentsView from './views/AgentsView';
import Settings from './views/Settings';
import TitleBar from './components/TitleBar';

const { Header, Sider, Content } = Layout;
const { Title } = Typography;

// Composant de layout principal
const MainLayout = () => {
  const [collapsed, setCollapsed] = useState(false);
  const location = useLocation();

  // Raccourci clavier pour toggle la sidebar
  useEffect(() => {
    const handleKeyDown = (event) => {
      if ((event.ctrlKey || event.metaKey) && event.key === 'b') {
        event.preventDefault();
        setCollapsed(prev => !prev);
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, []);

  return (
    <div style={{ minHeight: '100vh', display: 'flex', flexDirection: 'column' }}>
      <TitleBar />
      <Layout style={{ flex: 1, position: 'relative', overflow: 'hidden' }}>
        <Sider 
          collapsible 
          collapsed={collapsed} 
          onCollapse={setCollapsed}
          trigger={null} // Désactive le carré trigger par défaut
          breakpoint="lg"
          collapsedWidth={0}
          style={{ boxShadow: 'var(--shadow-md)' }}
        >
          <div className="logo">
            {!collapsed && "RLAMA"}
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
                key: '/agents',
                icon: <RobotOutlined />,
                label: <Link to="/agents">Agents</Link>,
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
                key: '/settings',
                icon: <SettingOutlined />,
                label: <Link to="/settings">Settings</Link>,
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
          <Header style={{ padding: 0, background: 'var(--bg-secondary)', boxShadow: 'var(--shadow-sm)', position: 'relative', overflow: 'hidden' }}>
            <div className="flex justify-between items-center px-4" style={{ height: '100%', overflow: 'hidden' }}>
              <div className="flex items-center flex-shrink-0">
                {/* Nouveau bouton hamburger élégant */}
                <Tooltip title={`${collapsed ? 'Expand' : 'Collapse'} sidebar (Ctrl+B)`} placement="bottom">
                  <Button
                    type="text"
                    icon={<MenuOutlined />}
                    onClick={() => setCollapsed(!collapsed)}
                    className="sidebar-toggle-btn"
                    style={{ 
                      marginRight: '16px',
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                      width: '40px',
                      height: '40px'
                    }}
                  />
                </Tooltip>

              </div>
              <div className="flex items-center gap-2" style={{ marginRight: '16px' }}>
                <Tooltip title="Discord Community" placement="bottom">
                  <Button
                    type="text"
                    icon={<svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                      <path d="M20.317 4.492c-1.53-.69-3.17-1.2-4.885-1.49a.075.075 0 0 0-.079.036c-.21.369-.444.85-.608 1.23a18.566 18.566 0 0 0-5.487 0 12.36 12.36 0 0 0-.617-1.23A.077.077 0 0 0 8.562 3c-1.714.29-3.354.8-4.885 1.491a.07.07 0 0 0-.032.027C.533 9.093-.32 13.555.099 17.961a.08.08 0 0 0 .031.055 20.03 20.03 0 0 0 5.993 2.98.078.078 0 0 0 .084-.026 13.83 13.83 0 0 0 1.226-1.963.074.074 0 0 0-.041-.104 13.201 13.201 0 0 1-1.872-.878.075.075 0 0 1-.008-.125c.126-.093.252-.19.372-.287a.075.075 0 0 1 .078-.01c3.927 1.764 8.18 1.764 12.061 0a.075.075 0 0 1 .079.009c.12.098.246.195.372.288a.075.075 0 0 1-.006.125c-.598.344-1.22.635-1.873.877a.075.075 0 0 0-.041.105c.36.687.772 1.341 1.225 1.962a.077.077 0 0 0 .084.028 19.963 19.963 0 0 0 6.002-2.981.076.076 0 0 0 .032-.054c.5-5.094-.838-9.52-3.549-13.442a.06.06 0 0 0-.031-.028zM8.02 15.278c-1.182 0-2.157-1.069-2.157-2.38 0-1.312.956-2.38 2.157-2.38 1.21 0 2.176 1.077 2.157 2.38 0 1.311-.956 2.38-2.157 2.38zm7.975 0c-1.183 0-2.157-1.069-2.157-2.38 0-1.312.955-2.38 2.157-2.38 1.21 0 2.176 1.077 2.157 2.38 0 1.311-.946 2.38-2.157 2.38z"/>
                    </svg>}
                    onClick={() => window.open('https://discord.gg/5D4prvVF4k', '_blank')}
                    style={{
                      color: 'var(--text-secondary)',
                      transition: 'var(--transition-fast)'
                    }}
                    onMouseEnter={(e) => e.target.style.color = '#5865F2'}
                    onMouseLeave={(e) => e.target.style.color = 'var(--text-secondary)'}
                  />
                </Tooltip>
                
                <Tooltip title="Website" placement="bottom">
                  <Button
                    type="text"
                    icon={<GlobalOutlined />}
                    onClick={() => window.open('https://rlama.dev/', '_blank')}
                    style={{
                      color: 'var(--text-secondary)',
                      transition: 'var(--transition-fast)'
                    }}
                    onMouseEnter={(e) => e.target.style.color = 'var(--accent-primary)'}
                    onMouseLeave={(e) => e.target.style.color = 'var(--text-secondary)'}
                  />
                </Tooltip>
                
                <Tooltip title="Follow us on X" placement="bottom">
                  <Button
                    type="text"
                    icon={<TwitterOutlined />}
                    onClick={() => window.open('https://x.com/Nuviastudio', '_blank')}
                    style={{
                      color: 'var(--text-secondary)',
                      transition: 'var(--transition-fast)'
                    }}
                    onMouseEnter={(e) => e.target.style.color = '#1DA1F2'}
                    onMouseLeave={(e) => e.target.style.color = 'var(--text-secondary)'}
                  />
                </Tooltip>
                
                <Tooltip title="LinkedIn Company" placement="bottom">
                  <Button
                    type="text"
                    icon={<LinkedinOutlined />}
                    onClick={() => window.open('https://www.linkedin.com/company/rlama', '_blank')}
                    style={{
                      color: 'var(--text-secondary)',
                      transition: 'var(--transition-fast)'
                    }}
                    onMouseEnter={(e) => e.target.style.color = '#0077B5'}
                    onMouseLeave={(e) => e.target.style.color = 'var(--text-secondary)'}
                  />
                </Tooltip>
                
                <Tooltip title="GitHub Repository" placement="bottom">
                  <Button
                    type="text"
                    icon={<GithubOutlined />}
                    onClick={() => window.open('https://github.com/DonTizi/rlama', '_blank')}
                    style={{
                      color: 'var(--text-secondary)',
                      transition: 'var(--transition-fast)'
                    }}
                    onMouseEnter={(e) => e.target.style.color = '#ffffff'}
                    onMouseLeave={(e) => e.target.style.color = 'var(--text-secondary)'}
                  />
                </Tooltip>
              </div>
            </div>
          </Header>
          <Content className="content-container">
            <div className="fade-in">
              <Outlet />
            </div>
          </Content>
        </Layout>
      </Layout>
    </div>
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
        path: "agents",
        element: <AgentsView />,
      },
      {
        path: "create",
        element: <CreateRag />,
      },
      {
        path: "settings",
        element: <Settings />,
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