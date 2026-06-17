import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useDemoStore } from '@/stores/demo'

describe('useDemoStore', () => {
  let store

  beforeEach(() => {
    setActivePinia(createPinia())
    store = useDemoStore()
  })

  describe('initial state', () => {
    it('should have demos array', () => {
      expect(Array.isArray(store.demos)).toBe(true)
    })

    it('should have default reactive values', () => {
      expect(store.activeDemo).toBeNull()
      expect(store.currentStep).toBe(0)
      expect(store.isPlaying).toBe(false)
      expect(store.playbackSpeed).toBe(1)
    })
  })

  describe('selectDemo', () => {
    it('should select a demo by id', () => {
      if (store.demos.length > 0) {
        const demoId = store.demos[0].id
        store.selectDemo(demoId)
        expect(store.activeDemo).toBeDefined()
        expect(store.currentStep).toBe(0)
      }
    })

    it('should set activeDemo to null/undefined for invalid id', () => {
      store.selectDemo('invalid_id')
      expect(store.activeDemo).toBeFalsy()
    })
  })

  describe('nextStep', () => {
    it('should increment current step', () => {
      store.currentStep = 0
      store.nextStep()
      expect(store.currentStep).toBeGreaterThanOrEqual(0)
    })
  })

  describe('prevStep', () => {
    it('should decrement current step when not at first step', () => {
      store.currentStep = 5
      store.prevStep()
      expect(store.currentStep).toBeLessThanOrEqual(5)
    })
  })

  describe('togglePlay', () => {
    it('should toggle isPlaying state', () => {
      const initialState = store.isPlaying
      store.togglePlay()
      expect(store.isPlaying).toBe(!initialState)
    })
  })

  describe('stop', () => {
    it('should stop playback and reset current step', () => {
      store.isPlaying = true
      store.currentStep = 5
      
      store.stop()
      expect(store.isPlaying).toBe(false)
      expect(store.currentStep).toBe(0)
    })
  })

  describe('setSpeed', () => {
    it('should set playback speed', () => {
      store.setSpeed(2)
      expect(store.playbackSpeed).toBe(2)
    })
  })
})
