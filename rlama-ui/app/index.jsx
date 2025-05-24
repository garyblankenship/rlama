import React, { useEffect } from 'react';
import { createRoot } from 'react-dom/client';
import { ConfigProvider } from 'antd';
import frFR from 'antd/lib/locale/fr_FR';
import App from './App';
import './styles/global.css';

// Wrapper component to clean up stray elements
const AppWrapper = () => {
  useEffect(() => {
    // Function to remove stray logo images
    const cleanupStrayElements = () => {
      const strayImages = document.querySelectorAll('body > img, #root > img');
      strayImages.forEach(img => {
        if (img.src && img.src.includes('rlama.svg')) {
          img.style.display = 'none';
          console.log('Removed stray logo image');
        }
      });
    };
    
    // Run initial cleanup
    cleanupStrayElements();
    
    // Set interval to periodically check and remove stray elements
    const interval = setInterval(cleanupStrayElements, 1000);
    
    return () => clearInterval(interval);
  }, []);
  
  return <App />;
};

// Add error boundary component
class ErrorBoundary extends React.Component {
  constructor(props) {
    super(props);
    this.state = { hasError: false, error: null, errorInfo: null };
  }

  static getDerivedStateFromError(error) {
    return { hasError: true };
  }

  componentDidCatch(error, errorInfo) {
    this.setState({ error, errorInfo });
    console.error("React Error:", error, errorInfo);
  }

  render() {
    if (this.state.hasError) {
      return (
        <div style={{ padding: 20, color: 'red' }}>
          <h1>Something went wrong.</h1>
          <pre>{this.state.error?.toString()}</pre>
          <pre>{this.state.errorInfo?.componentStack}</pre>
        </div>
      );
    }
    return this.props.children;
  }
}

// Add debugging log
console.log("React app initializing...");
const rootElement = document.getElementById('root');

if (!rootElement) {
  console.error("Root element not found! DOM may not be ready.");
  document.addEventListener('DOMContentLoaded', () => {
    console.log("DOM is now ready, attempting to mount app");
    mount();
  });
} else {
  mount();
}

function mount() {
  try {
    const root = createRoot(document.getElementById('root'));
    console.log("Root created, rendering app");
    
    root.render(
      <ErrorBoundary>
        <React.StrictMode>
          <ConfigProvider locale={frFR}>
            <AppWrapper />
          </ConfigProvider>
        </React.StrictMode>
      </ErrorBoundary>
    );
    console.log("App rendered successfully");
  } catch (error) {
    console.error("Failed to mount React app:", error);
    document.body.innerHTML = `
      <div style="padding: 20px; color: red;">
        <h1>Fatal Error</h1>
        <p>${error.toString()}</p>
      </div>
    `;
  }
} 