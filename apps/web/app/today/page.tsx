'use client';

import { useEffect, useState } from 'react';
import Navigation from '@/components/Navigation';
import { apiClient } from '@/lib/api-client';
import type { LogEntry } from '@/lib/types';

export default function TodayPage() {
  const [logEntries, setLogEntries] = useState<LogEntry[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showNewLog, setShowNewLog] = useState(false);
  const [newLogContent, setNewLogContent] = useState('');

  const loadTodayLogs = async () => {
    try {
      const data = await apiClient.getTodayLogEntries();
      setLogEntries(data);
    } catch {
      setError("Failed to load today's log entries");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadTodayLogs();
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const handleCreateLog = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    try {
      await apiClient.createLogEntry({
        content: newLogContent,
        log_date: new Date().toISOString().split('T')[0],
      });
      setNewLogContent('');
      setShowNewLog(false);
      loadTodayLogs();
    } catch {
      setError('Failed to create log entry');
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50">
        <Navigation />
        <div className="flex items-center justify-center h-[calc(100vh-4rem)]">
          <div className="text-xl">Loading...</div>
        </div>
      </div>
    );
  }

  const today = new Date().toLocaleDateString('en-US', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  });

  return (
    <div className="min-h-screen bg-gray-50">
      <Navigation />
      <div className="max-w-4xl mx-auto py-6 sm:px-6 lg:px-8">
        <div className="px-4 py-6 sm:px-0">
          <div className="mb-6">
            <h1 className="text-3xl font-bold text-gray-900">Today</h1>
            <p className="mt-2 text-gray-600">{today}</p>
          </div>

          {error && (
            <div className="rounded-md bg-red-50 p-4 mb-4">
              <div className="text-sm text-red-800">{error}</div>
            </div>
          )}

          <div className="mb-6">
            <button
              onClick={() => setShowNewLog(!showNewLog)}
              className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500"
            >
              {showNewLog ? 'Cancel' : 'Log Your Progress'}
            </button>
          </div>

          {showNewLog && (
            <div className="bg-white shadow rounded-lg p-6 mb-6">
              <h2 className="text-xl font-semibold mb-4">What did you work on today?</h2>
              <form onSubmit={handleCreateLog}>
                <div className="space-y-4">
                  <div>
                    <textarea
                      required
                      rows={4}
                      className="block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                      value={newLogContent}
                      onChange={(e) => setNewLogContent(e.target.value)}
                      placeholder="Describe what you accomplished today..."
                    />
                  </div>
                  <div className="flex justify-end space-x-3">
                    <button
                      type="button"
                      onClick={() => {
                        setShowNewLog(false);
                        setNewLogContent('');
                      }}
                      className="px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                    >
                      Cancel
                    </button>
                    <button
                      type="submit"
                      className="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500"
                    >
                      Save Log
                    </button>
                  </div>
                </div>
              </form>
            </div>
          )}

          <div className="bg-white shadow rounded-lg">
            <div className="px-6 py-4 border-b border-gray-200">
              <h2 className="text-lg font-semibold text-gray-900">
                Today&apos;s Logs ({logEntries.length})
              </h2>
            </div>
            <div className="divide-y divide-gray-200">
              {logEntries.length === 0 ? (
                <div className="p-6 text-center text-gray-500">
                  <p className="mb-2">No logs yet for today.</p>
                  <p className="text-sm">Start logging your progress to track what you accomplish!</p>
                </div>
              ) : (
                logEntries.map((log) => (
                  <div key={log.id} className="p-6">
                    <p className="text-gray-900 whitespace-pre-wrap">{log.content}</p>
                    <div className="mt-2 flex items-center text-sm text-gray-500">
                      <span>
                        {new Date(log.created_at).toLocaleTimeString('en-US', {
                          hour: '2-digit',
                          minute: '2-digit',
                        })}
                      </span>
                      {log.project_id && (
                        <>
                          <span className="mx-2">•</span>
                          <span>Project ID: {log.project_id}</span>
                        </>
                      )}
                      {log.task_id && (
                        <>
                          <span className="mx-2">•</span>
                          <span>Task ID: {log.task_id}</span>
                        </>
                      )}
                    </div>
                  </div>
                ))
              )}
            </div>
          </div>

          <div className="mt-6 bg-blue-50 rounded-lg p-6">
            <h3 className="text-lg font-semibold text-blue-900 mb-2">💡 Tips</h3>
            <ul className="list-disc list-inside text-blue-800 space-y-1 text-sm">
              <li>Log your progress throughout the day to track your work</li>
              <li>Use the project pages to link logs to specific projects and tasks</li>
              <li>Regular logging helps you stay accountable and see your progress over time</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  );
}
