```javascript
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { createPinia, setActivePinia } from 'pinia';
import { useDemoStore } from './demo.js';

describe('useDemoStore', () => {
  let store;

  beforeEach(() => {
    setActivePinia(createPinia());
    store = useDemoStore();
    vi.spyOn(Date, 'now').mockReturnValue(1234567890);
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('should initialize with default values', () => {
    expect(store.demos.value).toHaveLength(3);
    expect(store.activeDemo.value).toBeNull();
    expect(store.currentStep.value).toBe(0);
    expect(store.isPlaying.value).toBe(false);
    expect(store.playbackSpeed.value).toBe(1);
    expect(store.demoResults.value).toEqual([]);
  });

  describe('selectDemo', () => {
    it('should set the active demo and reset state', () => {
      store.selectDemo('beginner_village');
      const expectedDemo = store.demos.value.find(d => d.id === 'beginner_village');
      expect(store.activeDemo.value).toEqual(expectedDemo);
      expect(store.currentStep.value).toBe(0);
      expect(store.demoResults.value).toEqual([]);
      expect(store.isPlaying.value).toBe(false);
    });

    it('should set activeDemo to undefined if demo not found', () => {
      store.selectDemo('non_existent');
      expect(store.activeDemo.value).toBeUndefined();
      expect(store.currentStep.value).toBe(0);
      expect(store.demoResults.value).toEqual([]);
      expect(store.isPlaying.value).toBe(false);
    });
  });

  describe('nextStep', () => {
    it('should increment currentStep if not at the end', () => {
      store.selectDemo('beginner_village');
      store.currentStep.value = 5;
      const result = store.nextStep();
      expect(result).toBe(true);
      expect(store.currentStep.value).toBe(6);
    });

    it('should return false and not increment if at the end', () => {
      store.selectDemo('beginner_village');
      store.currentStep.value = 13;
      const result = store.nextStep();
      expect(result).toBe(false);
      expect(store.currentStep.value).toBe(13);
    });

    it('should return false if no activeDemo', () => {
      const result = store.nextStep();
      expect(result).toBe(false);
    });
  });

  describe('prevStep', () => {
    it('should decrement currentStep if greater than 0', () => {
      store.selectDemo('beginner_village');
      store.currentStep.value = 5;
      const result = store.prevStep();
      expect(result).toBe(true);
      expect(store.currentStep.value).toBe(4);
    });

    it('should return false and not decrement if at 0', () => {
      store.selectDemo('beginner_village');
      store.currentStep.value = 0;
      const result = store.prevStep();
      expect(result).toBe(false);
      expect(store.currentStep.value).toBe(0);
    });
  });

  describe('goToStep', () => {
    it('should set currentStep to valid step', () => {
      store.selectDemo('beginner_village');
      store.goToStep(7);
      expect(store.currentStep.value).toBe(7);
    });

    it('should not change currentStep for invalid step (negative)', () => {
      store.selectDemo('beginner_village');
      store.currentStep.value = 5;
      store.goToStep(-1);
      expect(store.currentStep.value).toBe(5);
    });

    it('should not change currentStep for invalid step (out of bounds)', () => {
      store.selectDemo('beginner_village');
      store.currentStep.value = 5;
      store.goToStep(14);
      expect(store.currentStep.value).toBe(5);
    });

    it('should not change if no activeDemo', () => {
      store.currentStep.value = 5;
      store.goToStep(7);
      expect(store.currentStep.value).toBe(5);
    });
  });

  describe('togglePlay', () => {
    it('should toggle isPlaying', () => {