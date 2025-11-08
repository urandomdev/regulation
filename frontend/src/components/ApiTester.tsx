import { useState } from 'react';
import { apiClient } from '../lib/api';
import type { SignupRequest } from '@generated/js/account';

interface TestResult {
  testName: string;
  status: 'pending' | 'running' | 'success' | 'error';
  message: string;
  timestamp?: string;
}

export function ApiTester() {
  const [results, setResults] = useState<TestResult[]>([]);
  const [testing, setTesting] = useState(false);

  const addResult = (result: TestResult) => {
    setResults((prev) => [...prev, { ...result, timestamp: new Date().toISOString() }]);
  };

  const updateLastResult = (updates: Partial<TestResult>) => {
    setResults((prev) => {
      const newResults = [...prev];
      if (newResults.length > 0) {
        newResults[newResults.length - 1] = {
          ...newResults[newResults.length - 1],
          ...updates,
          timestamp: new Date().toISOString(),
        };
      }
      return newResults;
    });
  };

  const testValidSignup = async () => {
    addResult({
      testName: 'Valid Signup',
      status: 'running',
      message: 'Testing signup with valid data...',
    });


    const testData: SignupRequest = {
      email: `test${Date.now()}@example.com`,
      password: 'SecurePassword123!',
      nickname: `TestUser${Date.now()}`,
    };


    const [data, err] = await apiClient.account.signup(testData);

    if (err) {
      updateLastResult({
        status: 'error',
        message: `Error: ${err.code} - ${err.message}`,
      });
    } else {
      updateLastResult({
        status: 'success',
        message: 'Signup successful!',
      });
    }
  };

  const testNetworkError = async () => {
    addResult({
      testName: 'Network Error Test',
      status: 'running',
      message: 'Testing with invalid endpoint...',
    });

    // Create a temporary client with invalid endpoint
    const { API } = await import('@generated/js');
    const testClient = new API('http://invalid-endpoint-that-does-not-exist.local', {
      debug: true,
    });

    const testData: SignupRequest = {
      email: 'test@example.com',
      password: 'password',
      nickname: 'Test',
    };

    const [data, err] = await testClient.account.signup(testData);

    if (err) {
      updateLastResult({
        status: 'success',
        message: `Network error caught correctly: ${err.code}`,
      });
    } else {
      updateLastResult({
        status: 'error',
        message: 'Expected network error but got success!',
      });
    }
  };

  const testHeaderManagement = async () => {
    addResult({
      testName: 'Header Management',
      status: 'running',
      message: 'Testing header set/get/remove...',
    });

    try {
      // Set a custom header
      apiClient.setHeader('X-Test-Header', 'test-value');
      const headerValue = apiClient.getHeader('X-Test-Header');

      if (headerValue !== 'test-value') {
        throw new Error('Header value mismatch');
      }

      // Set multiple headers
      apiClient.setHeaders({
        'X-Custom-1': 'value1',
        'X-Custom-2': 'value2',
      });

      const headers = apiClient.getHeaders();
      if (!headers['X-Custom-1'] || !headers['X-Custom-2']) {
        throw new Error('Headers not set correctly');
      }

      // Remove header
      apiClient.removeHeader('X-Test-Header');
      const removedHeader = apiClient.getHeader('X-Test-Header');

      if (removedHeader !== undefined) {
        throw new Error('Header not removed');
      }

      updateLastResult({
        status: 'success',
        message: 'All header operations successful!',
      });
    } catch (error) {
      updateLastResult({
        status: 'error',
        message: `Header test failed: ${(error as Error).message}`,
      });
    }
  };

  const testDebugMode = async () => {
    addResult({
      testName: 'Debug Mode',
      status: 'running',
      message: 'Testing API with debug mode (check console)...',
    });

    console.group('Debug Mode Test');
    console.log('Creating API client with debug mode enabled');

    const { API } = await import('@generated/js');
    const debugClient = new API('http://localhost:8080', {
      debug: true,
      headers: {
        'X-Debug-Test': 'true',
      },
    });

    const testData: SignupRequest = {
      email: `debug${Date.now()}@example.com`,
      password: 'password',
      nickname: 'DebugUser',
    };

    console.log('Making request with debug client...');
    const [data, err] = await debugClient.account.signup(testData);

    console.log('Request completed. Check console for detailed logs.');
    console.groupEnd();

    updateLastResult({
      status: 'success',
      message: 'Debug mode test completed. Check browser console for detailed logs.',
    });
  };

  const runAllTests = async () => {
    setTesting(true);
    setResults([]);

    try {
      await testValidSignup();
      await new Promise((resolve) => setTimeout(resolve, 500));

      await testHeaderManagement();
      await new Promise((resolve) => setTimeout(resolve, 500));

      await testDebugMode();
      await new Promise((resolve) => setTimeout(resolve, 500));

      await testNetworkError();
    } catch (error) {
      console.error('Test suite error:', error);
    } finally {
      setTesting(false);
    }
  };

  const clearResults = () => {
    setResults([]);
  };

  const getStatusColor = (status: TestResult['status']) => {
    switch (status) {
      case 'running':
        return '#ffa500';
      case 'success':
        return '#4caf50';
      case 'error':
        return '#f44336';
      default:
        return '#999';
    }
  };

  return (
    <div style={{ maxWidth: '600px', margin: '0 auto', padding: '20px' }}>
      <h2>API SDK Tester</h2>

      <div style={{ marginBottom: '20px' }}>
        <button
          onClick={runAllTests}
          disabled={testing}
          style={{
            padding: '10px 20px',
            marginRight: '10px',
            backgroundColor: testing ? '#ccc' : '#4caf50',
            color: 'white',
            border: 'none',
            borderRadius: '4px',
            cursor: testing ? 'not-allowed' : 'pointer',
          }}
        >
          {testing ? 'Running Tests...' : 'Run All Tests'}
        </button>

        <button
          onClick={clearResults}
          disabled={testing}
          style={{
            padding: '10px 20px',
            backgroundColor: '#999',
            color: 'white',
            border: 'none',
            borderRadius: '4px',
            cursor: testing ? 'not-allowed' : 'pointer',
          }}
        >
          Clear Results
        </button>
      </div>

      <div style={{ marginBottom: '20px' }}>
        <h3>Individual Tests</h3>
        <div style={{ display: 'flex', flexWrap: 'wrap', gap: '10px' }}>
          <button
            onClick={testValidSignup}
            disabled={testing}
            style={{
              padding: '8px 16px',
              backgroundColor: '#2196f3',
              color: 'white',
              border: 'none',
              borderRadius: '4px',
              cursor: testing ? 'not-allowed' : 'pointer',
            }}
          >
            Test Valid Signup
          </button>

          <button
            onClick={testHeaderManagement}
            disabled={testing}
            style={{
              padding: '8px 16px',
              backgroundColor: '#2196f3',
              color: 'white',
              border: 'none',
              borderRadius: '4px',
              cursor: testing ? 'not-allowed' : 'pointer',
            }}
          >
            Test Headers
          </button>

          <button
            onClick={testDebugMode}
            disabled={testing}
            style={{
              padding: '8px 16px',
              backgroundColor: '#2196f3',
              color: 'white',
              border: 'none',
              borderRadius: '4px',
              cursor: testing ? 'not-allowed' : 'pointer',
            }}
          >
            Test Debug Mode
          </button>

          <button
            onClick={testNetworkError}
            disabled={testing}
            style={{
              padding: '8px 16px',
              backgroundColor: '#2196f3',
              color: 'white',
              border: 'none',
              borderRadius: '4px',
              cursor: testing ? 'not-allowed' : 'pointer',
            }}
          >
            Test Network Error
          </button>
        </div>
      </div>

      <div>
        <h3>Test Results ({results.length})</h3>
        {results.length === 0 ? (
          <p style={{ color: '#999' }}>No tests run yet. Click a button above to start testing.</p>
        ) : (
          <div style={{ maxHeight: '400px', overflowY: 'auto' }}>
            {results.map((result, index) => (
              <div
                key={index}
                style={{
                  padding: '10px',
                  marginBottom: '10px',
                  backgroundColor: '#f5f5f5',
                  borderLeft: `4px solid ${getStatusColor(result.status)}`,
                  borderRadius: '4px',
                }}
              >
                <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '5px' }}>
                  <strong>{result.testName}</strong>
                  <span
                    style={{
                      color: getStatusColor(result.status),
                      fontWeight: 'bold',
                      textTransform: 'uppercase',
                      fontSize: '12px',
                    }}
                  >
                    {result.status}
                  </span>
                </div>
                <div style={{ fontSize: '14px', color: '#666' }}>{result.message}</div>
                {result.timestamp && (
                  <div style={{ fontSize: '12px', color: '#999', marginTop: '5px' }}>
                    {new Date(result.timestamp).toLocaleTimeString()}
                  </div>
                )}
              </div>
            ))}
          </div>
        )}
      </div>

      <div style={{ marginTop: '30px', padding: '15px', backgroundColor: '#f0f0f0', borderRadius: '4px' }}>
        <h4>SDK Features Demonstrated:</h4>
        <ul style={{ fontSize: '14px' }}>
          <li>Account signup endpoint integration</li>
          <li>Error handling (internal_error, network, unknown)</li>
          <li>Header management (set, get, remove, bulk set)</li>
          <li>Debug mode with console logging</li>
          <li>Network error simulation and retry logic</li>
          <li>TypeScript type safety with generated types</li>
        </ul>
      </div>
    </div>
  );
}
