import { describe, it, expect, vi, beforeEach } from 'vitest';
import { setActivePinia, createPinia } from 'pinia';
import { useTaskStore } from '@/stores/task';

describe('Task Store', () => {
  let store;

  beforeEach(() => {
    setActivePinia(createPinia());
    store = useTaskStore();
  });

  describe('Initial State', () => {
    it('should have initial tasks with two items', () => {
      expect(store.tasks).toHaveLength(2);
      expect(store.tasks[0].id).toBe('task_001');
      expect(store.tasks[1].id).toBe('task_002');
    });

    it('should have initial flows with one item', () => {
      expect(store.flows).toHaveLength(1);
      expect(store.flows[0].id).toBe('flow_001');
    });

    it('should have correct task structure', () => {
      const task = store.tasks[0];
      expect(task).toHaveProperty('id');
      expect(task).toHaveProperty('name');
      expect(task).toHaveProperty('type');
      expect(task).toHaveProperty('description');
      expect(task).toHaveProperty('status');
      expect(task).toHaveProperty('trigger');
      expect(task).toHaveProperty('objectives');
      expect(task).toHaveProperty('rewards');
      expect(task).toHaveProperty('nextTask');
      expect(task).toHaveProperty('dialogue');
    });

    it('should have correct flow structure', () => {
      const flow = store.flows[0];
      expect(flow).toHaveProperty('id');
      expect(flow).toHaveProperty('name');
      expect(flow).toHaveProperty('description');
      expect(flow).toHaveProperty('nodes');
      expect(flow).toHaveProperty('edges');
    });
  });

  describe('addTask', () => {
    it('should add a new task with generated id', () => {
      const initialLength = store.tasks.length;
      const newTask = {
        name: '新任务',
        type: 'side',
        description: '这是一个新任务',
        status: 'active'
      };

      const dateNowMock = vi.spyOn(Date, 'now').mockReturnValue(1234567890);
      store.addTask(newTask);

      expect(store.tasks).toHaveLength(initialLength + 1);
      expect(store.tasks[store.tasks.length - 1].id).toBe('task_1234567890');
      expect(store.tasks[store.tasks.length - 1].name).toBe('新任务');
      dateNowMock.mockRestore();
    });

    it('should preserve all provided properties', () => {
      const newTask = {
        name: '测试任务',
        type: 'daily',
        description: '测试描述',
        status: 'locked',
        rewards: { exp: 50, gold: 100 }
      };

      vi.spyOn(Date, 'now').mockReturnValue(9876543210);
      store.addTask(newTask);

      const addedTask = store.tasks.find(t => t.id === 'task_9876543210');
      expect(addedTask.name).toBe('测试任务');
      expect(addedTask.type).toBe('daily');
      expect(addedTask.description).toBe('测试描述');
      expect(addedTask.status).toBe('locked');
      expect(addedTask.rewards).toEqual({ exp: 50, gold: 100 });
      
      Date.now.mockRestore();
    });
  });

  describe('updateTask', () => {
    it('should update an existing task', () => {
      const taskId = 'task_001';
      const updateData = {
        name: '更新后的任务名',
        status: 'completed'
      };

      store.updateTask(taskId, updateData);
      const updatedTask = store.tasks.find(t => t.id === taskId);

      expect(updatedTask.name).toBe('更新后的任务名');
      expect(updatedTask.status).toBe('completed');
      // Original data should be preserved
      expect(updatedTask.type).toBe('main');
      expect(updatedTask.description).toBe('新来的冒险者，先去杂货铺买些必需品吧');
    });

    it('should not update if task id does not exist', () => {
      const initialTasks = [...store.tasks];
      store.updateTask('non_existent_task', { name: '更新' });
      
      expect(store.tasks).toEqual(initialTasks);
    });

    it('should handle partial updates', () => {
      const taskId = 'task_001';
      const originalTask = store.tasks.find(t => t.id === taskId);
      const originalName = originalTask.name;
      const originalDescription = originalTask.description;

      store.updateTask(taskId, { status: 'inactive' });
      const updatedTask = store.tasks.find(t => t.id === taskId);

      expect(updatedTask.status).toBe('inactive');
      expect(updatedTask.name).toBe(originalName);
      expect(updatedTask.description).toBe(originalDescription);
    });
  });

  describe('deleteTask', () => {
    it('should delete a task by id', () => {
      const taskIdToDelete = 'task_001';
      const initialLength = store.tasks.length;

      store.deleteTask(taskIdToDelete);

      expect(store.tasks).toHaveLength(initialLength - 1);
      expect(store.tasks.find(t => t.id === taskIdToDelete)).toBeUndefined();
    });

    it('should not delete if task id does not exist', () => {
      const initialTasks = [...store.tasks];
      store.deleteTask('non_existent_task');

      expect(store.tasks).toEqual(initialTasks);
    });

    it('should preserve other tasks when deleting', () => {
      const taskIdToDelete = 'task_001';
      const otherTaskId = 'task_002';

      store.deleteTask(taskIdToDelete);

      expect(store.tasks).toHaveLength(1);
      expect(store.tasks[0].id).toBe(otherTaskId);
    });
  });

  describe('getTaskById', () => {
    it('should return task by id', () => {
      const taskId = 'task_001';
      const task = store.getTaskById(taskId);

      expect(task).toBeDefined();
      expect(task.id).toBe(taskId);
      expect(task.name).toBe('初来乍到');
    });

    it('should return undefined for non-existent task id', () => {
      const task = store.getTaskById('non_existent_task');
      expect(task).toBeUndefined();
    });

    it('should return the correct task among multiple', () => {
      const task2 = store.getTaskById('task_002');
      expect(task2).toBeDefined();
      expect(task2.id).toBe('task_002');
      expect(task2.name).toBe('装备自己');
    });
  });

  describe('addFlow', () => {
    it('should add a new flow with generated id', () => {
      const initialLength = store.flows.length;
      const newFlow = {
        name: '新流程',
        description: '测试流程',
        nodes: [{ id: 'node_1', type: 'start', position: { x: 0, y: 0 }, data: { label: '开始' } }],
        edges: []
      };

      const dateNowMock = vi.spyOn(Date, 'now').mockReturnValue(111222333);
      store.addFlow(newFlow);

      expect(store.flows).toHaveLength(initialLength + 1);
      expect(store.flows[store.flows.length - 1].id).toBe('flow_111222333');
      expect(store.flows[store.flows.length - 1].name).toBe('新流程');
      dateNowMock.mockRestore();
    });

    it('should preserve all flow properties', () => {
      const newFlow = {
        name: '完整流程',
        description: '包含节点和边的流程',
        nodes: [
          { id: 'start', type: 'start' },
          { id: 'end', type: 'end' }
        ],
        edges: [{ id: 'e1', source: 'start', target: 'end' }]
      };

      vi.spyOn(Date, 'now').mockReturnValue(444555666);
      store.addFlow(newFlow);

      const addedFlow = store.flows.find(f => f.id === 'flow_444555666');
      expect(addedFlow.name).toBe('完整流程');
      expect(addedFlow.description).toBe('包含节点和边的流程');
      expect(addedFlow.nodes).toHaveLength(2);
      expect(addedFlow.edges).toHaveLength(1);
      
      Date.now.mockRestore();
    });
  });

  describe('updateFlow', () => {
    it('should update an existing flow', () => {
      const flowId = 'flow_001';
      const updateData = {
        name: '更新后的流程',
        description: '更新后的描述'
      };

      store.updateFlow(flowId, updateData);
      const updatedFlow = store.flows.find(f => f.id === flowId);

      expect(updatedFlow.name).toBe('更新后的流程');
      expect(updatedFlow.description).toBe('更新后的描述');
      expect(updatedFlow.nodes).toHaveLength(11);
    });

    it('should not update if flow id does not exist', () => {
      const initialFlows = [...store.flows];
      store.updateFlow('non_existent_flow', { name: '更新' });

      expect(store.flows).toEqual(initialFlows);
    });

    it('should handle partial updates', () => {
      const flowId = 'flow_001';
      const originalFlow = store.flows.find(f => f.id === flowId);
      const originalNodes = originalFlow.nodes;

      store.updateFlow(flowId, { description: '新描述' });
      const updatedFlow = store.flows.find(f => f.id === flowId);

      expect(updatedFlow.description).toBe('新描述');
      expect(updatedFlow.name).toBe(originalFlow.name);
      expect(updatedFlow.nodes).toEqual(originalNodes);
    });
  });

  describe('deleteFlow', () => {
    it('should delete a flow by id', () => {
      const flowIdToDelete = 'flow_001';
      const initialLength = store.flows.length;

      store.deleteFlow(flowIdToDelete);

      expect(store.flows).toHaveLength(initialLength - 1);
      expect(store.flows.find(f => f.id === flowIdToDelete)).toBeUndefined();
    });

    it('should not delete if flow id does not exist', () => {
      const initialFlows = [...store.flows];
      store.deleteFlow('non_existent_flow');

      expect(store.flows).toEqual(initialFlows);
    });
  });

  describe('Integration Tests', () => {
    it('should handle complete task workflow', () => {
      // Add a new task
      const newTask = {
        name: '集成测试任务',
        type: 'integration',
        status: 'active'
      };
      vi.spyOn(Date, 'now').mockReturnValue(999888777);
      store.addTask(newTask);
      
      // Get the task
      const addedTask = store.getTaskById('task_999888777');
      expect(addedTask).toBeDefined();
      expect(addedTask.name).toBe('集成测试任务');

      // Update the task
      store.updateTask('task_999888777', { status: 'completed', name: '完成的任务' });
      const updatedTask = store.getTaskById('task_999888777');
      expect(updatedTask.status).toBe('completed');
      expect(updatedTask.name).toBe('完成的任务');

      // Delete the task
      store.deleteTask('task_999888777');
      expect(store.getTaskById('task_999888777')).toBeUndefined();
      
      Date.now.mockRestore();
    });

    it('should handle complete flow workflow', () => {
      // Add a new flow
      const newFlow = {
        name: '测试流程',
        nodes: [{ id: 'node_1' }],
        edges: []
      };
      vi.spyOn(Date, 'now').mockReturnValue(777666555);
      store.addFlow(newFlow);
      
      // Get the flow
      const addedFlow = store.flows.find(f => f.id === 'flow_777666555');
      expect(addedFlow).toBeDefined();
      expect(addedFlow.name).toBe('测试流程');

      // Update the flow
      store.updateFlow('flow_777666555', { name: '更新后的流程' });
      const updatedFlow = store.flows.find(f => f.id === 'flow_777666555');
      expect(updatedFlow.name).toBe('更新后的流程');

      // Delete the flow
      store.deleteFlow('flow_777666555');
      expect(store.flows.find(f => f.id === 'flow_777666555')).toBeUndefined();
      
      Date.now.mockRestore();
    });
  });
});