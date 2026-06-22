import { describe, it, expect, beforeEach, vi } from 'vitest';
import { setActivePinia, createPinia } from 'pinia';
import { usePromptStore } from './prompt.js';

describe('Prompt Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  describe('初始状态', () => {
    it('should have correct initial templates', () => {
      const store = usePromptStore();
      expect(store.templates).toHaveLength(3);
      expect(store.templates[0].id).toBe('template_001');
      expect(store.templates[0].name).toBe('NPC基础人设');
      expect(store.templates[0].category).toBe('system');
    });

    it('should have correct initial variables', () => {
      const store = usePromptStore();
      expect(store.variables).toHaveLength(10);
      expect(store.variables[0].name).toBe('npc_name');
      expect(store.variables[0].type).toBe('string');
      expect(store.variables[0].source).toBe('npc');
    });
  });

  describe('addTemplate', () => {
    it('should add a new template with generated id', () => {
      const store = usePromptStore();
      const initialLength = store.templates.length;
      
      const newTemplate = {
        name: '测试模板',
        category: 'custom',
        content: '测试内容',
        variables: []
      };
      
      store.addTemplate(newTemplate);
      
      expect(store.templates).toHaveLength(initialLength + 1);
      expect(store.templates[store.templates.length - 1].name).toBe('测试模板');
      expect(store.templates[store.templates.length - 1].id).toMatch(/^template_\d+$/);
    });

    it('should preserve template properties when adding', () => {
      const store = usePromptStore();
      const newTemplate = {
        name: '测试模板',
        category: 'custom',
        content: '测试内容{{variable}}',
        variables: [{ name: 'variable', type: 'string', description: '变量' }]
      };
      
      store.addTemplate(newTemplate);
      
      const addedTemplate = store.templates[store.templates.length - 1];
      expect(addedTemplate.content).toBe('测试内容{{variable}}');
      expect(addedTemplate.variables).toHaveLength(1);
    });
  });

  describe('updateTemplate', () => {
    it('should update existing template', () => {
      const store = usePromptStore();
      const templateId = 'template_001';
      const newData = {
        name: '更新后的NPC人设',
        content: '更新后的内容'
      };
      
      store.updateTemplate(templateId, newData);
      
      const updatedTemplate = store.templates.find(t => t.id === templateId);
      expect(updatedTemplate.name).toBe('更新后的NPC人设');
      expect(updatedTemplate.content).toBe('更新后的内容');
    });

    it('should not modify template when id does not exist', () => {
      const store = usePromptStore();
      const initialTemplates = [...store.templates];
      
      store.updateTemplate('nonexistent_id', { name: '不存在的模板' });
      
      expect(store.templates).toEqual(initialTemplates);
    });

    it('should merge new data with existing template properties', () => {
      const store = usePromptStore();
      const templateId = 'template_001';
      const originalTemplate = store.templates.find(t => t.id === templateId);
      
      store.updateTemplate(templateId, { category: 'updated' });
      
      const updatedTemplate = store.templates.find(t => t.id === templateId);
      expect(updatedTemplate.category).toBe('updated');
      expect(updatedTemplate.name).toBe(originalTemplate.name);
      expect(updatedTemplate.content).toBe(originalTemplate.content);
    });
  });

  describe('deleteTemplate', () => {
    it('should delete existing template', () => {
      const store = usePromptStore();
      const templateId = 'template_001';
      const initialLength = store.templates.length;
      
      store.deleteTemplate(templateId);
      
      expect(store.templates).toHaveLength(initialLength - 1);
      expect(store.templates.find(t => t.id === templateId)).toBeUndefined();
    });

    it('should not change templates when id does not exist', () => {
      const store = usePromptStore();
      const initialLength = store.templates.length;
      
      store.deleteTemplate('nonexistent_id');
      
      expect(store.templates).toHaveLength(initialLength);
    });
  });

  describe('addVariable', () => {
    it('should add a new variable', () => {
      const store = usePromptStore();
      const initialLength = store.variables.length;
      
      const newVariable = {
        name: 'new_variable',
        type: 'string',
        description: '新变量',
        source: 'custom'
      };
      
      store.addVariable(newVariable);
      
      expect(store.variables).toHaveLength(initialLength + 1);
      expect(store.variables[store.variables.length - 1]).toEqual(newVariable);
    });

    it('should allow duplicate variable names', () => {
      const store = usePromptStore();
      const duplicateVariable = { ...store.variables[0] };
      
      store.addVariable(duplicateVariable);
      
      const variablesWithName = store.variables.filter(v => v.name === duplicateVariable.name);
      expect(variablesWithName).toHaveLength(2);
    });
  });

  describe('updateVariable', () => {
    it('should update existing variable', () => {
      const store = usePromptStore();
      const variableName = 'npc_name';
      const newData = {
        description: 'NPC名字（更新后）',
        source: 'npc_updated'
      };
      
      store.updateVariable(variableName, newData);
      
      const updatedVariable = store.variables.find(v => v.name === variableName);
      expect(updatedVariable.description).toBe('NPC名字（更新后）');
      expect(updatedVariable.source).toBe('npc_updated');
    });

    it('should not modify variable when name does not exist', () => {
      const store = usePromptStore();
      const initialVariables = [...store.variables];
      
      store.updateVariable('nonexistent_variable', { description: '不存在的变量' });
      
      expect(store.variables).toEqual(initialVariables);
    });

    it('should merge new data with existing variable properties', () => {
      const store = usePromptStore();
      const variableName = 'npc_name';
      const originalVariable = store.variables.find(v => v.name === variableName);
      
      store.updateVariable(variableName, { type: 'updated_type' });
      
      const updatedVariable = store.variables.find(v => v.name === variableName);
      expect(updatedVariable.type).toBe('updated_type');
      expect(updatedVariable.description).toBe(originalVariable.description);
    });
  });

  describe('deleteVariable', () => {
    it('should delete existing variable', () => {
      const store = usePromptStore();
      const variableName = 'npc_name';
      const initialLength = store.variables.length;
      
      store.deleteVariable(variableName);
      
      expect(store.variables).toHaveLength(initialLength - 1);
      expect(store.variables.find(v => v.name === variableName)).toBeUndefined();
    });

    it('should not change variables when name does not exist', () => {
      const store = usePromptStore();
      const initialLength = store.variables.length;
      
      store.deleteVariable('nonexistent_variable');
      
      expect(store.variables).toHaveLength(initialLength);
    });
  });
});