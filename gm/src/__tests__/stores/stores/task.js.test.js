```javascript
import { describe, it, expect, vi, beforeEach } from 'vitest';

vi.mock('pinia', () => ({
  defineStore: vi.fn((name, setup) => {
    return () => {
      const store = setup();
      return store;
    };
  })
}));

vi.mock('vue', () => ({
  ref: vi.fn((initialValue) => ({
    value: initialValue
  }))
}));

import { useTaskStore } from './task.js';

describe('useTaskStore', () => {
  let store;

  beforeEach(() => {
    store = useTaskStore();
  });

  describe('initial state', () => {
    it('should have tasks with correct initial data', () => {
      expect(store.tasks.value).toBeInstanceOf(Array);
      expect(store.tasks.value.length).toBe(2);
      expect(store.tasks.value[0].id).toBe('task_001');
      expect(store.tasks.value[0].name).toBe('初来乍到');
      expect(store.tasks.value[0].status).toBe('active');
    });

    it('should have flows with correct initial data', () => {
      expect(store.flows.value).toBeInstanceOf(Array);
      expect(store.flows.value.length).toBe(1);
      expect(store.flows.value[0].id).toBe('flow_001');
      expect(store.flows.value[0].name).toBe('NPC出门购物流程');
    });
  });

  describe('addTask', () => {
    it('should add a new task with generated id', () => {
      const newTask = { name: 'New Task', type: 'side' };
      store.addTask(newTask);

      expect(store.tasks.value.length).toBe(3);
      const addedTask = store.tasks.value[store.tasks.value.length - 1];
      expect(addedTask.id).toMatch(/^task_/);
      expect(addedTask.name).toBe('New Task');
      expect(addedTask.type).toBe('side');
    });
  });

  describe('updateTask', () => {
    it('should update existing task data', () => {
      const taskId = 'task_001';
      const updateData = { name: 'Updated Name' };
      store.updateTask(taskId, updateData);

      const updatedTask = store.getTaskById(taskId);
      expect(updatedTask.name).toBe('Updated Name');
      expect(updatedTask.status).toBe('active');
    });

    it('should not update if task id not found', () => {
      const initialLength = store.tasks.value.length;
      store.updateTask('nonexistent_id', { name: 'test' });
      expect(store.tasks.value.length).toBe(initialLength);
    });
  });

  describe('deleteTask', () => {
    it('should delete task by id', () => {
      const taskId = 'task_001';
      store.deleteTask(taskId);

      expect(store.tasks.value.length).toBe(1);
      expect(store.getTaskById(taskId)).toBeUndefined();
    });

    it('should not delete if task id not found', () => {
      const initialLength = store.tasks.value.length;
      store.deleteTask('nonexistent_id');
      expect(store.tasks.value.length).toBe(initialLength);
    });
  });

  describe('getTaskById', () => {
    it('should return task by id', () => {
      const task = store.getTaskById('task_001');
      expect(task).toBeDefined();
      expect(task.id).toBe('task_001');
      expect(task.name).toBe('初来乍到');
    });

    it('should return undefined for non-existent id', () => {
      const task = store.getTaskById('nonexistent_id');
      expect(task).toBeUndefined();
    });
  });

  describe('addFlow', () => {
    it('should add a new flow with generated id', () => {
      const newFlow = { name: 'New Flow', description: 'A new flow' };
      store.addFlow(newFlow);

      expect(store.flows.value.length).toBe(2);
      const addedFlow = store.flows.value[store.flows.value.length - 1];
      expect(addedFlow.id).toMatch(/^flow_/);
      expect(addedFlow.name).toBe('New Flow');
      expect(addedFlow.description).toBe('A new flow');
    });
  });

  describe('updateFlow', () => {
    it('should update existing flow data', () => {
      const flowId = 'flow_001';
      const updateData = { name: 'Updated Flow Name' };
      store.updateFlow(flowId, updateData);

      const updatedFlow = store.flows.value.find(f => f.id === flowId);
      expect(updatedFlow.name).toBe('Updated Flow Name');
      expect(updatedFlow.description).toBe('NPC从家里出发，前往商店购买物品的完整流程');
    });

    it('should not update if flow id not found', () => {
      const initialLength = store.flows.value.length;
      store.updateFlow('nonexistent_id', { name: 'test' });
      expect(store.flows.value.length).toBe(initialLength);
    });
  });

  describe('deleteFlow', () => {
    it('should delete flow by id', () => {
      const flowId = 'flow_001';
      store.deleteFlow(flowId);

      expect(store.flows.value.length).toBe(0);
    });

    it('should not delete if flow id not found', () => {
      const initialLength = store.flows.value.length;
      store.deleteFlow('nonexistent_id');
      expect(store.flows.value.length).toBe(initialLength);
    });
  });
});
```