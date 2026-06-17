```javascript
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { setActivePinia, createPinia } from 'pinia';
import { usePromptStore } from './stores/prompt';

describe('usePromptStore', () => {
  let store;

  beforeEach(() => {
    setActivePinia(createPinia());
    store = usePromptStore();
  });

  describe('initial state', () => {
    it('should have templates array with initial templates', () => {
      expect(store.templates).toHaveLength(3);
      expect(store.templates[0].id).toBe('template_001');
      expect(store.templates[0].name).toBe('NPC基础人设');
      expect(store.templates[1].id).toBe('template_002');
      expect(store.templates[2].id).toBe('template_003');
    });

    it('should have variables array with initial variables', () => {
      expect(store.variables).toHaveLength(10);
      expect(store.variables[0].name).toBe('npc_name');
      expect(store.variables[0].type).toBe('string');
      expect(store.variables[0].source).toBe('npc');
    });
  });

  describe('addTemplate', () => {
    it('should add a new template with generated id', () => {
      const mockDate = 1234567890;
      vi.spyOn(Date, 'now').mockImplementation(() => mockDate);
      
      const newTemplate = {
        name: 'Test Template',
        content: 'Test content',
        variables: [],
        category: 'test'
      };
      
      store.addTemplate(newTemplate);
      
      expect(store.templates).toHaveLength(4);
      expect(store.templates[3].id).toBe(`template_${mockDate}`);
      expect(store.templates[3].name).toBe('Test Template');
      expect(store.templates[3].content).toBe('Test content');
      expect(store.templates[3].category).toBe('test');
      
      vi.restoreAllMocks();
    });

    it('should preserve existing template properties when adding', () => {
      const initialLength = store.templates.length;
      
      store.addTemplate({
        id: 'existing_id',
        name: 'Test Template'
      });
      
      expect(store.templates).toHaveLength(initialLength + 1);
      expect(store.templates[store.templates.length - 1].id).toBe('existing_id');
    });
  });

  describe('updateTemplate', () => {
    it('should update existing template by id', () => {
      const updates = {
        name: 'Updated Name',
        content: 'Updated content'
      };
      
      store.updateTemplate('template_001', updates);
      
      const updatedTemplate = store.templates.find(t => t.id === 'template_001');
      expect(updatedTemplate.name).toBe('Updated Name');
      expect(updatedTemplate.content).toBe('Updated content');
      expect(updatedTemplate.category).toBe('system');
    });

    it('should not modify other templates', () => {
      const originalTemplate = { ...store.templates[1] };
      
      store.updateTemplate('template_001', { name: 'Updated' });
      
      expect(store.templates[1]).toEqual(originalTemplate);
    });

    it('should do nothing if template id not found', () => {
      const originalTemplates = [...store.templates];
      
      store.updateTemplate('non_existent_id', { name: 'Updated' });
      
      expect(store.templates).toEqual(originalTemplates);
    });
  });

  describe('deleteTemplate', () => {
    it('should delete template by id', () => {
      const initialLength = store.templates.length;
      
      store.deleteTemplate('template_002');
      
      expect(store.templates).toHaveLength(initialLength - 1);
      expect(store.templates.find(t => t.id === 'template_002')).toBeUndefined();
    });

    it('should not delete other templates', () => {
      const template1Before = store.templates.find(t => t.id === 'template_001');
      const template3Before = store.templates.find(t => t.id === 'template_003');
      
      store.deleteTemplate('template_002');
      
      expect(store.templates.find(t => t.id === 'template_001')).toEqual(template1Before);
      expect(store.templates.find(t => t.id === 'template_003')).toEqual(template3Before);
    });
  });

  describe('addVariable', () => {
    it('should add a new variable to variables array', () => {
      const initialLength = store.variables.length;
      const newVariable = {
        name: 'new_var',
        type: 'string',
        description: 'New variable',
        source: 'test'
      };
      
      store.addVariable(newVariable);
      
      expect(store.variables).toHaveLength(initialLength + 1);
      expect(store.variables[store.variables.length - 1]).toEqual(newVariable);
    });
  });

  describe('updateVariable', () => {
    it('should update existing variable by name', () => {
      const updates = {
        description: 'Updated description',
        type: 'number'
      };
      
      store.updateVariable('npc_name', updates);
      
      const updatedVariable = store.variables.find(v => v.name === 'npc_name');
      expect(updatedVariable.description).toBe('Updated description');
      expect(updatedVariable.type).toBe('number');
      expect(updatedVariable.source).toBe('npc');
    });

    it('should do nothing if variable name not found', () => {
      const originalVariables = [...store.variables];
      
      store.updateVariable('non_existent', { description: 'Test' });
      
      expect(store.variables).toEqual(originalVariables);
    });
  });

  describe('deleteVariable', () => {
    it('should delete variable by name', () => {
      const initialLength = store.variables.length;
      
      store.deleteVariable('player_name');
      
      expect(store.variables).toHaveLength(initialLength - 1);
      expect(store.variables.find(v => v.name === 'player_name')).toBeUndefined();
    });

    it('should not delete other variables', () => {
      const npcNameBefore = store.variables.find(v => v.name === 'npc_name');
      const npcTitleBefore = store.variables.find(v => v.name === 'npc_title');
      
      store.deleteVariable('player_name');
      
      expect(store.variables.find(v => v.name === 'npc_name')).toEqual(npcNameBefore);
      expect(store.variables.find(v => v.name === 'npc_title')).toEqual(npcTitleBefore);
    });
  });
});
```