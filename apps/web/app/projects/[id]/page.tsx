'use client';

import { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import Navigation from '@/components/Navigation';
import { apiClient } from '@/lib/api-client';
import type { Project, Task, LogEntry } from '@/lib/types';

export default function ProjectPage() {
  const params = useParams();
  const router = useRouter();
  const projectId = parseInt(params.id as string);

  const [project, setProject] = useState<Project | null>(null);
  const [tasks, setTasks] = useState<Task[]>([]);
  const [logEntries, setLogEntries] = useState<LogEntry[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  // New task form
  const [showNewTask, setShowNewTask] = useState(false);
  const [newTaskTitle, setNewTaskTitle] = useState('');
  const [newTaskDescription, setNewTaskDescription] = useState('');
  const [newTaskStatus, setNewTaskStatus] = useState('todo');

  // New log entry form
  const [showNewLog, setShowNewLog] = useState(false);
  const [newLogContent, setNewLogContent] = useState('');
  const [selectedTaskId, setSelectedTaskId] = useState<number | undefined>();

  useEffect(() => {
    loadProjectData();
  }, [projectId]);

  const loadProjectData = async () => {
    try {
      const [projectData, tasksData, logsData] = await Promise.all([
        apiClient.getProject(projectId),
        apiClient.getTasks(projectId),
        apiClient.getLogEntries(projectId),
      ]);
      setProject(projectData);
      setTasks(tasksData);
      setLogEntries(logsData);
    } catch (err) {
      setError('Failed to load project data');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleCreateTask = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    try {
      await apiClient.createTask({
        project_id: projectId,
        title: newTaskTitle,
        description: newTaskDescription,
        status: newTaskStatus,
      });
      setNewTaskTitle('');
      setNewTaskDescription('');
      setNewTaskStatus('todo');
      setShowNewTask(false);
      loadProjectData();
    } catch (err) {
      setError('Failed to create task');
    }
  };

  const handleCreateLog = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    try {
      await apiClient.createLogEntry({
        project_id: projectId,
        task_id: selectedTaskId,
        content: newLogContent,
        log_date: new Date().toISOString().split('T')[0],
      });
      setNewLogContent('');
      setSelectedTaskId(undefined);
      setShowNewLog(false);
      loadProjectData();
    } catch (err) {
      setError('Failed to create log entry');
    }
  };

  const handleUpdateTaskStatus = async (taskId: number, newStatus: string) => {
    try {
      const task = tasks.find((t) => t.id === taskId);
      if (task) {
        await apiClient.updateTask(taskId, {
          title: task.title,
          description: task.description,
          status: newStatus,
        });
        loadProjectData();
      }
    } catch (err) {
      setError('Failed to update task status');
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

  if (!project) {
    return (
      <div className="min-h-screen bg-gray-50">
        <Navigation />
        <div className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
          <div className="text-center">
            <h1 className="text-2xl font-bold text-gray-900 mb-4">Project not found</h1>
            <button
              onClick={() => router.push('/')}
              className="text-blue-600 hover:text-blue-500"
            >
              Go back to projects
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <Navigation />
      <div className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <div className="px-4 py-6 sm:px-0">
          <div className="mb-6">
            <button
              onClick={() => router.push('/')}
              className="text-blue-600 hover:text-blue-500 mb-4"
            >
              ← Back to projects
            </button>
            <h1 className="text-3xl font-bold text-gray-900">{project.name}</h1>
            <p className="mt-2 text-gray-600">{project.description}</p>
          </div>

          {error && (
            <div className="rounded-md bg-red-50 p-4 mb-4">
              <div className="text-sm text-red-800">{error}</div>
            </div>
          )}

          {/* Tasks Section */}
          <div className="mb-8">
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-2xl font-semibold text-gray-900">Tasks</h2>
              <button
                onClick={() => setShowNewTask(!showNewTask)}
                className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
              >
                New Task
              </button>
            </div>

            {showNewTask && (
              <div className="bg-white shadow rounded-lg p-6 mb-6">
                <h3 className="text-xl font-semibold mb-4">Create New Task</h3>
                <form onSubmit={handleCreateTask}>
                  <div className="space-y-4">
                    <div>
                      <label htmlFor="task-title" className="block text-sm font-medium text-gray-700">
                        Title
                      </label>
                      <input
                        type="text"
                        id="task-title"
                        required
                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                        value={newTaskTitle}
                        onChange={(e) => setNewTaskTitle(e.target.value)}
                      />
                    </div>
                    <div>
                      <label htmlFor="task-description" className="block text-sm font-medium text-gray-700">
                        Description
                      </label>
                      <textarea
                        id="task-description"
                        rows={3}
                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                        value={newTaskDescription}
                        onChange={(e) => setNewTaskDescription(e.target.value)}
                      />
                    </div>
                    <div>
                      <label htmlFor="task-status" className="block text-sm font-medium text-gray-700">
                        Status
                      </label>
                      <select
                        id="task-status"
                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                        value={newTaskStatus}
                        onChange={(e) => setNewTaskStatus(e.target.value)}
                      >
                        <option value="todo">To Do</option>
                        <option value="in_progress">In Progress</option>
                        <option value="done">Done</option>
                      </select>
                    </div>
                    <div className="flex justify-end space-x-3">
                      <button
                        type="button"
                        onClick={() => setShowNewTask(false)}
                        className="px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                      >
                        Cancel
                      </button>
                      <button
                        type="submit"
                        className="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                      >
                        Create
                      </button>
                    </div>
                  </div>
                </form>
              </div>
            )}

            <div className="bg-white shadow rounded-lg divide-y">
              {tasks.length === 0 ? (
                <div className="p-6 text-center text-gray-500">No tasks yet.</div>
              ) : (
                tasks.map((task) => (
                  <div key={task.id} className="p-6">
                    <div className="flex items-start justify-between">
                      <div className="flex-1">
                        <h3 className="text-lg font-medium text-gray-900">{task.title}</h3>
                        <p className="mt-1 text-sm text-gray-600">{task.description}</p>
                      </div>
                      <select
                        value={task.status}
                        onChange={(e) => handleUpdateTaskStatus(task.id, e.target.value)}
                        className="ml-4 border border-gray-300 rounded-md shadow-sm py-1 px-2 text-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                      >
                        <option value="todo">To Do</option>
                        <option value="in_progress">In Progress</option>
                        <option value="done">Done</option>
                      </select>
                    </div>
                  </div>
                ))
              )}
            </div>
          </div>

          {/* Log Entries Section */}
          <div>
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-2xl font-semibold text-gray-900">Log Entries</h2>
              <button
                onClick={() => setShowNewLog(!showNewLog)}
                className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500"
              >
                New Log Entry
              </button>
            </div>

            {showNewLog && (
              <div className="bg-white shadow rounded-lg p-6 mb-6">
                <h3 className="text-xl font-semibold mb-4">Create New Log Entry</h3>
                <form onSubmit={handleCreateLog}>
                  <div className="space-y-4">
                    <div>
                      <label htmlFor="log-task" className="block text-sm font-medium text-gray-700">
                        Task (Optional)
                      </label>
                      <select
                        id="log-task"
                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                        value={selectedTaskId || ''}
                        onChange={(e) => setSelectedTaskId(e.target.value ? parseInt(e.target.value) : undefined)}
                      >
                        <option value="">No task</option>
                        {tasks.map((task) => (
                          <option key={task.id} value={task.id}>
                            {task.title}
                          </option>
                        ))}
                      </select>
                    </div>
                    <div>
                      <label htmlFor="log-content" className="block text-sm font-medium text-gray-700">
                        Content
                      </label>
                      <textarea
                        id="log-content"
                        required
                        rows={4}
                        className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                        value={newLogContent}
                        onChange={(e) => setNewLogContent(e.target.value)}
                        placeholder="What did you work on?"
                      />
                    </div>
                    <div className="flex justify-end space-x-3">
                      <button
                        type="button"
                        onClick={() => setShowNewLog(false)}
                        className="px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                      >
                        Cancel
                      </button>
                      <button
                        type="submit"
                        className="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500"
                      >
                        Create
                      </button>
                    </div>
                  </div>
                </form>
              </div>
            )}

            <div className="bg-white shadow rounded-lg divide-y">
              {logEntries.length === 0 ? (
                <div className="p-6 text-center text-gray-500">No log entries yet.</div>
              ) : (
                logEntries.map((log) => (
                  <div key={log.id} className="p-6">
                    <div className="flex items-start justify-between">
                      <div className="flex-1">
                        <p className="text-gray-900">{log.content}</p>
                        {log.task_id && (
                          <p className="mt-1 text-sm text-gray-500">
                            Task: {tasks.find((t) => t.id === log.task_id)?.title || 'Unknown'}
                          </p>
                        )}
                      </div>
                      <div className="ml-4 text-sm text-gray-500">
                        {new Date(log.log_date).toLocaleDateString()}
                      </div>
                    </div>
                  </div>
                ))
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
