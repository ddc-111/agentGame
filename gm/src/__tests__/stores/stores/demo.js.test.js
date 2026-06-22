import { describe, it, expect, vi, beforeEach } from 'vitest';
import { createPinia, setActivePinia } from 'pinia';
import { useDemoStore } from './stores/demo.js';

describe('useDemoStore', () => {
  let store;

  beforeEach(() => {
    setActivePinia(createPinia());
    store = useDemoStore();
  });

  describe('initial state', () => {
    it('should have correct initial demos', () => {
      expect(store.demos).toHaveLength(3);
      expect(store.demos[0].id).toBe('beginner_village');
      expect(store.demos[1].id).toBe('combat_demo');
      expect(store.demos[2].id).toBe('npc_ai_demo');
    });

    it('should have null activeDemo initially', () => {
      expect(store.activeDemo).toBeNull();
    });

    it('should have zero currentStep initially', () => {
      expect(store.currentStep).toBe(0);
    });

    it('should have isPlaying false initially', () => {
      expect(store.isPlaying).toBe(false);
    });

    it('should have playbackSpeed 1 initially', () => {
      expect(store.playbackSpeed).toBe(1);
    });

    it('should have empty demoResults initially', () => {
      expect(store.demoResults).toEqual([]);
    });
  });

  describe('selectDemo', () => {
    it('should select demo by id', () => {
      store.selectDemo('combat_demo');
      expect(store.activeDemo).toEqual(store.demos[1]);
    });

    it('should reset currentStep to 0', () => {
      store.currentStep = 5;
      store.selectDemo('combat_demo');
      expect(store.currentStep).toBe(0);
    });

    it('should clear demoResults', () => {
      store.demoResults = [{ step: 1 }];
      store.selectDemo('combat_demo');
      expect(store.demoResults).toEqual([]);
    });

    it('should set isPlaying to false', () => {
      store.isPlaying = true;
      store.selectDemo('combat_demo');
      expect(store.isPlaying).toBe(false);
    });

    it('should handle invalid demoId gracefully', () => {
      store.selectDemo('invalid_id');
      expect(store.activeDemo).toBeUndefined();
    });
  });

  describe('nextStep', () => {
    beforeEach(() => {
      store.selectDemo('combat_demo');
    });

    it('should increment currentStep', () => {
      store.nextStep();
      expect(store.currentStep).toBe(1);
    });

    it('should return true when there are more steps', () => {
      expect(store.nextStep()).toBe(true);
    });

    it('should not go beyond last step', () => {
      store.currentStep = store.activeDemo.steps - 1;
      expect(store.nextStep()).toBe(false);
      expect(store.currentStep).toBe(store.activeDemo.steps - 1);
    });

    it('should do nothing when no activeDemo', () => {
      store.activeDemo = null;
      const result = store.nextStep();
      expect(result).toBe(false);
      expect(store.currentStep).toBe(0);
    });
  });

  describe('prevStep', () => {
    beforeEach(() => {
      store.selectDemo('combat_demo');
      store.currentStep = 5;
    });

    it('should decrement currentStep', () => {
      store.prevStep();
      expect(store.currentStep).toBe(4);
    });

    it('should return true when not at first step', () => {
      expect(store.prevStep()).toBe(true);
    });

    it('should not go below zero', () => {
      store.currentStep = 0;
      expect(store.prevStep()).toBe(false);
      expect(store.currentStep).toBe(0);
    });
  });

  describe('goToStep', () => {
    beforeEach(() => {
      store.selectDemo('combat_demo');
    });

    it('should go to specified step', () => {
      store.goToStep(5);
      expect(store.currentStep).toBe(5);
    });

    it('should handle step 0', () => {
      store.currentStep = 5;
      store.goToStep(0);
      expect(store.currentStep).toBe(0);
    });

    it('should not go below 0', () => {
      store.goToStep(-1);
      expect(store.currentStep).toBe(0);
    });

    it('should not exceed max steps', () => {
      store.goToStep(store.activeDemo.steps);
      expect(store.currentStep).toBe(0);
    });

    it('should do nothing when no activeDemo', () => {
      store.activeDemo = null;
      store.goToStep(5);
      expect(store.currentStep).toBe(0);
    });
  });

  describe('togglePlay', () => {
    it('should toggle isPlaying from false to true', () => {
      store.isPlaying = false;
      store.togglePlay();
      expect(store.isPlaying).toBe(true);
    });

    it('should toggle isPlaying from true to false', () => {
      store.isPlaying = true;
      store.togglePlay();
      expect(store.isPlaying).toBe(false);
    });
  });

  describe('stop', () => {
    beforeEach(() => {
      store.selectDemo('combat_demo');
      store.isPlaying = true;
      store.currentStep = 5;
    });

    it('should set isPlaying to false', () => {
      store.stop();
      expect(store.isPlaying).toBe(false);
    });

    it('should reset currentStep to 0', () => {
      store.stop();
      expect(store.currentStep).toBe(0);
    });
  });

  describe('setSpeed', () => {
    it('should set playbackSpeed', () => {
      store.setSpeed(2);
      expect(store.playbackSpeed).toBe(2);
    });

    it('should handle decimal speeds', () => {
      store.setSpeed(0.5);
      expect(store.playbackSpeed).toBe(0.5);
    });
  });

  describe('addResult', () => {
    beforeEach(() => {
      vi.spyOn(Date, 'now').mockReturnValue(1234567890);
    });

    afterEach(() => {
      vi.restoreAllMocks();
    });

    it('should add result with current step and timestamp', () => {
      store.currentStep = 5;
      store.addResult({ score: 100, success: true });
      
      expect(store.demoResults).toHaveLength(1);
      expect(store.demoResults[0]).toEqual({
        step: 5,
        timestamp: 1234567890,
        score: 100,
        success: true
      });
    });

    it('should preserve existing results', () => {
      store.demoResults = [{ step: 1 }];
      store.addResult({ test: 'data' });
      
      expect(store.demoResults).toHaveLength(2);
    });
  });

  describe('resetResults', () => {
    it('should clear demoResults', () => {
      store.demoResults = [{ step: 1 }, { step: 2 }];
      store.resetResults();
      expect(store.demoResults).toEqual([]);
    });
  });
});